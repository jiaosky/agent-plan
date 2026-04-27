# 第 04 章：SSE 流式输出与 Channel 通信

## 1. 学习什么

本章学习大模型应用后端中非常高频、也很容易在面试和工作中被问到的能力：SSE 流式输出与 Channel 通信。

传统后端常见模式是：

```text
请求 -> 处理 -> 一次性返回 JSON
```

但大模型生成文本可能需要几秒到几十秒。如果等完整答案生成完再返回，用户体验会很差。所以 ChatGPT 类产品通常是：

```text
请求 -> 模型边生成 -> 后端边推送 -> 前端边显示
```

你需要掌握以下内容：

- SSE 是什么，适合解决什么问题。
- SSE 和 WebSocket、普通 HTTP 响应的区别。
- token、chunk、delta、finish_reason 的含义。
- Go 后端如何实现 SSE。
- Go 的 `channel` 在流式转发中的作用。
- 如何处理用户取消、客户端断连、超时。
- 如何处理模型流式错误和部分输出。
- 什么是背压，为什么流式系统要关注消费速度。
- 如何记录首 token 时间、总耗时、输出 token、取消状态。
- 前后端如何约定事件格式。

本章学完后，你应该能说清楚：

> 大模型流式输出通常用 SSE 实现，后端从模型供应商读取 chunk，再通过 Go channel 或回调转发给 HTTP 响应，同时要处理 flush、断连、取消、超时、错误、日志和最终内容落库。

## 2. 为什么学习

SSE 是 AI 应用后端非常实用的工程能力。

在传统业务系统里，很多接口耗时较短，用户能接受等待一个完整 JSON 返回。但大模型场景不同：

- 模型首 token 可能需要 500ms 到数秒。
- 长文本生成可能持续几十秒。
- 代码生成、报告生成、总结文档更慢。
- 中途可能发生工具调用。
- 用户经常需要点击“停止生成”。
- 前端需要实时显示模型正在工作。

如果没有流式输出，用户看到的是：

```text
点击发送 -> 页面一直转圈 -> 很久后突然出现一大段回答
```

体验很差，也不利于用户判断模型是否跑偏。

有了流式输出后，体验变成：

```text
点击发送 -> 很快出现第一个字 -> 回答持续展开 -> 用户可随时停止
```

面试里，SSE 常常和这些问题一起出现：

- 你怎么实现 ChatGPT 那种打字机效果？
- SSE 和 WebSocket 怎么选？
- Go 里怎么把模型 chunk 推给前端？
- 客户端断开后，后端怎么停止模型调用？
- 流式输出过程中报错怎么办？
- 如何统计首 token 时间？
- 如何保存完整回答？

这章学好后，你会明显不像“只调过 API 的人”，而更像能做 AI 产品后端的人。

## 3. 知识详细内容

### 3.1 什么是 SSE

SSE 是 Server-Sent Events，服务端发送事件。

它基于 HTTP，特点是：

- 客户端发起一个 HTTP 请求。
- 服务端保持连接不关闭。
- 服务端不断向客户端写入事件数据。
- 客户端通过 EventSource 或 fetch stream 接收。
- 通信方向主要是服务端到客户端。

典型响应头：

```http
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

SSE 数据格式类似：

```text
event: message
data: {"delta":"你好"}

event: message
data: {"delta":"，我是"}

event: done
data: {"finish_reason":"stop"}
```

注意：SSE 每个事件之间通常用空行分隔。

### 3.2 SSE、WebSocket、普通 HTTP 的区别

| 方式 | 特点 | 适合场景 |
| --- | --- | --- |
| 普通 HTTP | 一次请求，一次完整响应 | 短任务、结构化结果、非实时 |
| SSE | 客户端请求，服务端持续推送 | 大模型聊天、进度通知、日志流 |
| WebSocket | 双向长连接 | 实时协作、语音、游戏、复杂双向交互 |

大模型文本聊天一般优先考虑 SSE，因为：

- 实现比 WebSocket 简单。
- 基于 HTTP，容易接入网关和鉴权。
- 适合单向持续输出。
- 浏览器原生支持 EventSource。
- 对 ChatGPT 式文本流很合适。

什么时候选 WebSocket？

- 前后端需要频繁双向通信。
- 语音实时对话。
- 多人协作。
- 复杂 Agent 控制台。
- 同一个连接里要持续发送用户操作和服务端事件。

面试表达：

> 普通聊天流式输出用 SSE 就够了，因为主要是服务端持续推送模型生成内容。如果需要复杂双向交互，比如实时语音、协作编辑、持续控制 Agent，再考虑 WebSocket。

### 3.3 token、chunk、delta、finish_reason

流式输出里常见几个词：

| 术语 | 含义 |
| --- | --- |
| token | 模型处理文本的基本单位 |
| chunk | 模型流式返回的一小段数据 |
| delta | 本次 chunk 新增的内容 |
| finish_reason | 本次生成结束原因 |

模型不一定每次返回一个完整词，也不一定每次返回一个完整句子。它可能返回：

```text
chunk1: "你"
chunk2: "好"
chunk3: "，"
chunk4: "我"
chunk5: "是"
```

后端需要把这些 delta 转发给前端，前端再拼成完整内容。

常见 finish_reason：

| 值 | 含义 |
| --- | --- |
| stop | 正常结束 |
| length | 达到最大输出长度 |
| tool_calls | 模型要调用工具 |
| content_filter | 内容被安全策略拦截 |
| error | 业务侧可自定义的错误结束 |

不同供应商字段可能不同，模型网关可以统一映射成自己的结构。

### 3.4 大模型流式链路

一个常见链路：

```text
前端发送问题
  -> 后端鉴权
  -> 后端创建 trace_id
  -> 后端构造 messages
  -> 后端调用模型 stream API
  -> 模型返回 chunk
  -> 后端转换成 SSE event
  -> 前端实时渲染 delta
  -> 后端累计完整回答
  -> 生成结束后保存会话和 usage
```

如果用户中途停止：

```text
用户点击停止
  -> 前端关闭连接或发取消请求
  -> 后端 context cancel
  -> 模型请求停止
  -> 后端保存已生成部分或标记 canceled
```

如果客户端断网：

```text
连接断开
  -> r.Context().Done()
  -> 后端停止读取模型流
  -> 释放 goroutine 和连接
  -> 记录 canceled 或 client_closed
```

### 3.5 Go 后端实现 SSE 的基本要点

Go 里实现 SSE 的关键点：

- 设置正确响应头。
- 判断 `http.Flusher` 是否可用。
- 每次写入事件后调用 `Flush()`。
- 监听 `r.Context().Done()`。
- 不要在连接断开后继续写。
- 结束时发送 done event。

示例骨架：

```go
func streamHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }

    ctx := r.Context()
    chunks, err := llm.Stream(ctx, req)
    if err != nil {
        writeSSE(w, "error", `{"message":"stream start failed"}`)
        flusher.Flush()
        return
    }

    for {
        select {
        case <-ctx.Done():
            return
        case chunk, ok := <-chunks:
            if !ok {
                writeSSE(w, "done", `{"finish_reason":"stop"}`)
                flusher.Flush()
                return
            }
            writeSSE(w, "message", toJSON(chunk))
            flusher.Flush()
        }
    }
}
```

这里的重点不是背代码，而是理解：

```text
模型流 -> 后端 channel -> SSE event -> 前端渲染
```

### 3.6 Channel 在流式通信中的作用

Go 的 channel 很适合表达“持续产生的数据流”。

模型客户端可以把供应商返回的 chunk 转成 channel：

```go
type ChatChunk struct {
    Delta        string
    FinishReason string
}

type LLMClient interface {
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}
```

调用方只需要：

```go
for chunk := range chunks {
    // 写给 SSE
}
```

channel 的价值：

- 解耦模型读取和 HTTP 写入。
- 方便在中间做日志、过滤、累计。
- 方便 select 监听 ctx.Done。
- 方便把不同供应商统一成一种流。

但 channel 也要小心：

- 生产者不能在消费者退出后一直阻塞。
- 出错时要能关闭 channel。
- 不要忘记释放底层 HTTP 响应体。
- 如果使用带缓冲 channel，要控制缓冲大小。

### 3.7 事件格式设计

后端和前端要约定 SSE event 格式。

建议至少包含：

```text
message：正常增量内容
error：错误事件
done：结束事件
metadata：可选元信息
tool_call：工具调用事件
```

示例：

```text
event: message
data: {"delta":"你好","index":1}

event: metadata
data: {"trace_id":"abc123","model":"deepseek-chat"}

event: done
data: {"finish_reason":"stop","input_tokens":123,"output_tokens":456}
```

不要只传裸文本。结构化 event 更方便：

- 前端区分消息、错误、结束。
- 后端记录 trace。
- 后续扩展工具调用事件。
- 排查问题时更清楚。

### 3.8 断连、取消和超时

流式接口最容易漏掉的是取消。

需要处理三类结束：

| 类型 | 说明 | 后端处理 |
| --- | --- | --- |
| 正常结束 | 模型返回 stop | 发送 done，保存完整回答 |
| 用户取消 | 用户点停止 | cancel context，保存部分内容或标记取消 |
| 客户端断连 | 浏览器关闭或网络断开 | 停止模型调用，释放资源 |
| 服务超时 | 超过最大生成时间 | cancel context，返回 error 或 done |

Go 里关键是 `r.Context()`：

```go
select {
case <-r.Context().Done():
    // 客户端断开或请求被取消
    return
case chunk := <-chunks:
    // 正常写出
}
```

如果你没有处理取消，可能出现：

- 用户页面关了，模型还在继续生成。
- API 调用成本继续增加。
- goroutine 堆积。
- HTTP 连接泄漏。
- 会话状态不一致。

### 3.9 背压是什么

背压可以简单理解为：

> 生产者生成数据的速度超过消费者处理数据的速度。

在流式模型里：

```text
模型供应商快速返回 chunk
  -> 后端读取
  -> 前端或网络消费慢
  -> 写响应阻塞
```

可能原因：

- 用户网络慢。
- 前端渲染慢。
- 代理层缓冲。
- 后端 channel 缓冲太大。
- 多个中间处理步骤阻塞。

处理思路：

- channel 缓冲不要无限大。
- 写入超时或监听 context。
- 客户端断连及时取消上游。
- 不要在流式循环里做重 CPU 工作。
- 大日志异步处理，不阻塞主流。

面试时知道这个概念就很加分，不需要一开始实现复杂背压系统。

### 3.10 代理和网关注意事项

SSE 经过 Nginx、API 网关、负载均衡时可能遇到问题：

- 响应被缓冲，前端不能实时看到。
- 网关超时断开长连接。
- gzip 压缩影响 flush。
- idle timeout 太短。
- 某些代理不支持长连接。

常见处理：

- 关闭代理缓冲。
- 设置合适的 read timeout。
- 定期发送 heartbeat。
- 确保响应头正确。
- 避免中间层把 SSE 当普通响应缓存。

heartbeat 示例：

```text
event: ping
data: {}
```

或者 SSE 注释行：

```text
: ping
```

### 3.11 日志和观测

流式接口要记录的不只是总耗时。

建议记录：

```text
trace_id
user_id
tenant_id
session_id
model
stream=true
first_token_ms
total_latency_ms
input_tokens
output_tokens
chunk_count
finish_reason
status
error_code
canceled_by_user
client_disconnected
```

其中 `first_token_ms` 很重要：

- 用户体感通常取决于多久看到第一个字。
- 总耗时长不一定糟糕，如果首 token 很快，体验可能还行。

还要保存完整回答：

```text
chunk1 + chunk2 + ... + chunkN = full_content
```

保存时要注意：

- 正常结束保存为 completed。
- 用户取消保存为 canceled。
- 出错保存为 failed 或 partial_failed。
- 可以保留部分内容，供前端恢复或审计。

### 3.12 前端协作要点

后端需要和前端约定：

- 请求接口路径。
- 鉴权方式。
- SSE event 类型。
- message 事件字段。
- error 事件字段。
- done 事件字段。
- 用户停止生成如何触发。
- 断线后是否重连。
- 是否展示 trace_id。

前端通常需要：

- 收到 delta 后追加到当前消息。
- 收到 done 后结束 loading。
- 收到 error 后展示错误。
- 用户停止时关闭连接或调用 cancel 接口。
- 处理 Markdown 渲染和代码块未闭合问题。

后端不要把所有责任推给前端。流式协议设计不清楚，前后端都会痛苦。

### 3.13 常见误区

误区一：SSE 就是循环写字符串。

实际还要处理响应头、flush、event 格式、断连、取消、错误和日志。

误区二：用户关页面后模型会自动停止。

不一定。后端必须监听 context，并把取消传给模型客户端。

误区三：流式接口不需要保存完整回答。

会话系统通常仍然要保存完整 assistant 消息，否则刷新页面后历史丢失。

误区四：WebSocket 一定比 SSE 高级。

不是。普通大模型文本流 SSE 更简单、更合适。

误区五：流式开始后还能随便重试。

如果已经给用户输出了一半，再重试可能产生重复内容或上下文不一致。流式重试要谨慎。

## 4. 考题、编程题与验收

### 4.1 概念题

1. SSE 是什么？为什么适合大模型聊天？
2. SSE 和 WebSocket 有什么区别？
3. chunk、delta、finish_reason 分别是什么意思？
4. 为什么每次写 SSE 后要 flush？
5. 用户关闭页面后，后端应该如何停止模型调用？
6. 为什么流式输出仍然要保存完整回答？
7. first_token_ms 和 total_latency_ms 有什么区别？
8. 什么是背压？在流式输出里可能怎么出现？
9. SSE 经过网关时可能遇到哪些问题？
10. 为什么流式输出开始后不适合无脑重试？

### 4.2 面试问答题

问题一：如何实现 ChatGPT 那种逐字输出效果？

参考回答：

> 后端可以使用模型供应商的 stream API，持续读取模型返回的 chunk。Go 服务端把 chunk 转成 SSE event 写给前端，每次写完 flush。前端收到 delta 后追加到当前消息。后端同时累计完整内容，结束时保存会话和 usage。整个过程中要监听 request context，处理用户取消、客户端断连和超时。

问题二：SSE 和 WebSocket 你怎么选？

参考回答：

> 如果只是大模型文本生成，主要是服务端持续推送给前端，SSE 更简单，基于 HTTP，也更容易接入现有网关和鉴权。如果是实时语音、多端协作、复杂双向控制，前后端都需要频繁发送消息，就更适合 WebSocket。

问题三：客户端断开连接后，后端怎么处理？

参考回答：

> Go 里可以监听 `r.Context().Done()`。一旦客户端断开或请求取消，就停止 SSE 写入，并把 context 取消传递给模型客户端，让上游模型请求中断。同时记录状态为 canceled 或 client_disconnected，释放 goroutine、HTTP 响应体和相关资源。

问题四：流式输出过程报错怎么办？

参考回答：

> 如果还没开始输出，可以直接返回错误事件或 HTTP 错误。如果已经输出了部分内容，通常发送 error event 或 done event 标记异常结束，保存部分内容和错误状态。不能无脑重试，因为用户已经看到一部分内容，重试可能导致重复或不一致。

### 4.3 编程题

编程题一：写一个 SSE 写入函数。

要求：

- 支持 event name。
- data 是 JSON 字符串。
- 每个事件后空行分隔。

参考方向：

```go
func writeSSE(w io.Writer, event string, data string) error {
    if event != "" {
        if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
            return err
        }
    }
    _, err := fmt.Fprintf(w, "data: %s\n\n", data)
    return err
}
```

验收标准：

- 能解释为什么要有空行。
- 能解释 event 和 data 的作用。
- 能解释为什么 data 最好是 JSON。

编程题二：实现一个最小 SSE handler。

要求：

- 设置 SSE 响应头。
- 获取 `http.Flusher`。
- 从 `chan ChatChunk` 读取数据。
- 每个 chunk 写 `message` event。
- channel 关闭后写 `done` event。
- 监听 `r.Context().Done()`。

伪代码方向：

```go
func handleStream(w http.ResponseWriter, r *http.Request) {
    flusher := setupSSE(w)
    chunks := make(chan ChatChunk)

    go produceChunks(r.Context(), chunks)

    for {
        select {
        case <-r.Context().Done():
            return
        case chunk, ok := <-chunks:
            if !ok {
                writeSSE(w, "done", `{"finish_reason":"stop"}`)
                flusher.Flush()
                return
            }
            writeSSE(w, "message", toJSON(chunk))
            flusher.Flush()
        }
    }
}
```

编程题三：设计流式调用日志字段。

要求字段：

- trace_id。
- user_id。
- session_id。
- model。
- first_token_ms。
- total_latency_ms。
- chunk_count。
- input_tokens。
- output_tokens。
- finish_reason。
- status。
- canceled_by_user。
- client_disconnected。

验收标准：

- 能查首 token 慢的问题。
- 能查用户取消比例。
- 能查模型输出长度。
- 能查错误和断连情况。

### 4.4 系统设计题

题目：设计一个支持流式输出的大模型聊天接口。

要求说明：

- 接口路径和请求参数。
- 如何鉴权和限流。
- 如何调用模型 stream API。
- 如何把模型 chunk 转成 SSE event。
- 如何处理断连、取消、超时。
- 如何保存完整回答。
- 如何记录 token 和日志。
- 如何和前端约定事件格式。

回答必须包含：

- SSE 响应头。
- `http.Flusher`。
- `context.Context`。
- channel 或回调式 chunk 转发。
- done/error event。
- first_token_ms。
- 会话保存。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能解释 SSE 为什么适合大模型聊天。
- 能说清 SSE 和 WebSocket 的选择标准。
- 能写出最小 SSE handler。
- 能用 Go channel 表达模型 chunk 流。
- 能解释用户取消和客户端断连如何处理。
- 能设计流式 event 格式。
- 能说出流式日志需要记录哪些字段。

## 5. 和其他知识点的相关性

和 LLM 调用的关系：

- 第 02 章讲模型普通调用和 stream 调用。
- 本章讲 stream 调用如何通过 SSE 返回给前端。

和 Prompt 的关系：

- Prompt 决定模型输出内容。
- SSE 决定输出内容如何实时传输。

和 RAG 的关系：

- RAG 问答通常也需要流式输出。
- 检索阶段可能先等待，模型生成阶段再流式返回。
- 后续可以把“正在检索”“正在生成”也作为事件推给前端。

和文件管理的关系：

- 文档总结、文件问答、报告生成都可能耗时较长，适合流式返回。
- 文件处理进度也可以用 SSE 推送。

和 Tool Calling 的关系：

- 模型流式输出中可能出现 tool_call 事件。
- 后端需要暂停文本输出、执行工具、再继续流式返回结果。

和 MCP 的关系：

- MCP 工具执行过程可以通过 SSE 展示给前端。
- Agent 调用 MCP Server 时，前端可以看到步骤流。

和 Agent 的关系：

- Agent 通常不是只输出最终答案，而是输出计划、工具调用、执行结果和总结。
- SSE 可以承载这些步骤事件。

和 Go 的关系：

- `channel` 适合表示 chunk 流。
- `context` 适合取消和超时。
- `http.Flusher` 是 SSE 的关键。
- goroutine 用于读取上游流并转发。

## 6. 演变过程与后续方向

### 6.1 从一次性响应到流式响应

传统接口一般是：

```text
请求 -> 完整处理 -> 返回 JSON
```

大模型聊天演变成：

```text
请求 -> 生成一点 -> 返回一点 -> 持续更新
```

这本质上是用户体验驱动的变化。模型生成慢，但只要首 token 快，用户就会觉得系统在响应。

### 6.2 从文本流到事件流

最初可以只流式返回文本：

```text
你 好 ， 我 是 ...
```

但产品复杂后，需要返回更多事件：

- message：文本增量。
- tool_call：工具调用。
- tool_result：工具结果。
- progress：任务进度。
- error：错误。
- done：结束。

所以流式接口会从“文本流”演变成“事件流”。

### 6.3 从单模型流到 Agent 执行流

普通聊天只需要模型输出。

Agent 场景需要展示：

```text
开始分析
  -> 选择工具
  -> 调用工具
  -> 工具返回
  -> 再次思考
  -> 输出最终答案
```

这会让 SSE 承载更多结构化事件，而不仅是 delta 文本。

### 6.4 从简单连接到可观测流式系统

成熟系统会关注：

- 首 token 时间。
- chunk 数量。
- 每段耗时。
- 断连率。
- 取消率。
- 网关超时。
- 上游模型错误率。
- 前端渲染异常。

流式输出不是“能跑就行”，而是需要可观测、可排障。

### 6.5 后续方向

SSE 和流式通信后续会继续发展：

- Agent step streaming。
- Tool call streaming。
- 多模态流式输出。
- 语音实时流。
- 前端可中断、可编辑、可恢复的生成。
- 长任务后台运行 + 事件订阅。
- 多模型并行输出和合并。
- 更统一的事件协议。

## 7. 工作中典型场景

### 场景一：产品要求实现“打字机效果”

你需要做：

- 后端接模型 stream API。
- 用 SSE 推送 delta。
- 前端逐步追加内容。
- 后端保存完整回答。
- done event 结束 loading。

你需要追问：

- 是否支持用户停止生成？
- 是否需要展示引用来源？
- 是否需要展示工具调用过程？
- 是否需要保存未完成内容？

### 场景二：用户点停止后，账单还在增加

可能原因：

- 前端只停止渲染，没有通知后端。
- 后端没有监听 context。
- 后端没有取消模型请求。
- 模型客户端没有关闭响应体。

处理方式：

- 前端关闭 SSE 连接或调用 cancel API。
- 后端监听 `r.Context().Done()`。
- 模型请求使用同一个 context。
- 记录 canceled 状态。

### 场景三：前端说不是流式，还是一次性出现

排查方向：

- 后端是否调用了 `Flush()`。
- 响应头是否是 `text/event-stream`。
- Nginx 是否开启了 buffering。
- 网关是否缓冲响应。
- 前端是否按流读取。
- 是否被 gzip 或代理影响。

### 场景四：流式输出中间报错

处理方式：

- 如果未输出内容，返回 error event。
- 如果已输出部分内容，发送 error 或 done 标记异常结束。
- 保存 partial 内容。
- 日志记录错误发生时的 chunk_count。
- 不要无脑重试。

### 场景五：首 token 很慢

排查方向：

- Prompt 是否太长。
- RAG 检索是否慢。
- 模型本身是否慢。
- 是否排队或限流。
- 是否网络慢。

优化方向：

- 并发检索。
- 缩短上下文。
- 使用更快模型。
- 预先返回 progress event。
- 缓存常见问题。

### 场景六：面试官问“Go 里 channel 在 AI 流式输出中怎么用”

可以这样回答：

> 我会让模型客户端把供应商返回的 stream chunk 转成只读 channel，handler 用 select 同时监听 chunk 和 request context。收到 chunk 就写成 SSE event 并 flush，context done 就停止上游请求并释放资源。channel 让模型读取和 HTTP 输出解耦，也方便统一不同模型供应商的流式协议。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

