# 第 11 章：Go 在 AI 后端中的典型场景

## 1. 学习什么

本章专门解决一个现实问题：你会写 Go，但面试时担心说不出 Go 在大模型公司里到底怎么用。

这章不按语法书讲 Go，而是按 AI 应用后端的典型场景讲：

- Go 如何做模型网关。
- Go 如何封装多模型 LLMClient。
- `context` 如何处理模型超时、取消和链路传递。
- `channel` 如何转发模型流式 chunk。
- `goroutine` 如何做并发检索和工具调用。
- `interface` 如何抽象模型供应商、向量库、工具系统。
- Go 如何实现 SSE。
- Go 如何做 RAG 后端、文件任务、embedding 队列。
- Go 如何做 Agent 编排和 MCP Server。
- Redis、MQ、HTTP Client 在 AI 后端中的用法。
- 面试时如何把“我会 Go”说成“我知道 Go 在 AI 后端怎么落地”。

本章学完后，你应该能说清楚：

> Go 在 AI 后端里适合做稳定的服务层和编排层，比如模型网关、SSE 流式服务、RAG 检索编排、Agent 工具调度、MCP Server、异步任务和限流观测。

## 2. 为什么学习

很多后端工程师说自己会 Go，但面试时容易停留在：

```text
我用 Go 写过接口。
我用过 Gin。
我会 goroutine 和 channel。
我用过 Redis 和 MySQL。
```

这些当然有用，但在大模型公司面试里，你需要进一步说出 AI 场景：

```text
goroutine 可以并发执行多个检索源。
channel 可以承接模型流式 chunk。
context 可以取消模型长请求。
interface 可以抽象 OpenAI、DeepSeek、通义等模型供应商。
Redis 可以做模型调用限流和会话缓存。
MQ 可以做文件解析、embedding、评测异步任务。
```

这样面试官会感觉你不是只会传统 CRUD，而是已经把 Go 能力迁移到了 AI 应用后端。

Go 在 AI 后端里的优势：

- 并发模型简单。
- HTTP 服务稳定。
- 部署方便。
- 资源占用较低。
- 适合做网关和编排。
- 适合连接各种外部服务。
- 适合写长连接和流式服务。
- 适合做任务 worker。

你不需要把 Go 讲成模型训练语言。Go 的定位是：

> AI 应用后端的稳定服务层、网关层、编排层和工程化治理层。

## 3. 知识详细内容

### 3.1 Go 在 AI 后端的整体定位

Go 通常不负责训练大模型。

Go 更适合做：

```text
API 服务
模型网关
SSE 流式输出
RAG 后端
文件处理任务编排
Agent 状态机
MCP Server
工具调用系统
限流和计费
日志和观测
```

可以这样表达：

> Python 生态更偏模型训练、数据处理和算法实验；Go 更适合承接线上服务、并发请求、长连接、网关、任务调度和稳定性治理。

这不是说 Go 不能做 AI，而是说 Go 在大模型公司里更偏工程化落地。

### 3.2 场景一：模型网关

模型网关负责统一接入多个模型供应商。

例如：

```text
业务服务 -> 模型网关 -> OpenAI / DeepSeek / 通义 / Claude / 私有模型
```

Go 适合做模型网关，因为：

- HTTP Client 能力成熟。
- 并发处理简单。
- 容易做统一中间件。
- 部署成独立服务方便。

模型网关要做：

- 统一请求结构。
- 统一响应结构。
- 供应商适配。
- 超时控制。
- 重试。
- fallback。
- usage 统计。
- 日志。
- 限流。

接口抽象：

```go
type LLMClient interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}
```

面试表达：

> 我会用 Go 做一层模型网关，业务层只依赖统一的 LLMClient 接口，不直接绑定某个供应商。不同模型做 adapter，网关统一处理 context 超时、重试、限流、usage 统计和日志。

### 3.3 场景二：context 控制超时和取消

大模型调用常常慢。

用户可能：

- 关闭页面。
- 点击停止生成。
- 请求超时。
- 网关断开。

Go 的 `context.Context` 很适合处理这种链路取消。

示例：

```go
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()

resp, err := llm.Chat(ctx, req)
```

使用原则：

- HTTP handler 使用 `r.Context()`。
- 下游模型调用必须传 ctx。
- RAG 检索也传 ctx。
- 工具调用也传 ctx。
- 用户取消后及时释放资源。

面试表达：

> 模型调用比普通接口更容易慢，所以 Go 后端要把 context 贯穿整个链路。用户断连或点击停止后，context cancel 能传递到模型请求、检索请求和工具调用，避免 goroutine 和连接泄漏，也能控制成本。

### 3.4 场景三：channel 转发流式 chunk

大模型流式输出可以抽象成一个数据流。

Go 的 channel 很适合表达：

```text
模型供应商 stream -> channel -> SSE handler -> 前端
```

结构：

```go
type ChatChunk struct {
    Delta        string
    FinishReason string
}
```

流式接口：

```go
Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
```

handler 中：

```go
for {
    select {
    case <-ctx.Done():
        return
    case chunk, ok := <-chunks:
        if !ok {
            return
        }
        writeSSE(w, "message", chunk)
        flusher.Flush()
    }
}
```

面试表达：

> 我会把模型流式返回统一成只读 channel，handler 用 select 同时监听 chunk 和 ctx.Done。收到 chunk 就写 SSE 并 flush，ctx 取消就停止上游请求，避免用户断开后模型还继续生成。

### 3.5 场景四：goroutine 并发检索和工具调用

AI 后端里经常要并发：

- 同时查向量库和关键词索引。
- 同时查多个知识库。
- 同时调用多个工具。
- 同时执行多个评测用例。
- 文件批量 embedding。

例如 RAG 混合检索：

```text
向量检索
关键词检索
权限过滤
合并去重
rerank
```

Go 可以用 goroutine + errgroup：

```go
g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
    vectorResults, err = vectorSearch(ctx, query)
    return err
})

g.Go(func() error {
    keywordResults, err = keywordSearch(ctx, query)
    return err
})

if err := g.Wait(); err != nil {
    return err
}
```

注意：

- 并发不是越多越好。
- 要限制并发数。
- 要传 context。
- 要处理部分失败。
- 要记录每个子任务耗时。

### 3.6 场景五：interface 抽象外部能力

AI 后端会接很多外部系统：

- 模型供应商。
- 向量数据库。
- Embedding 服务。
- Rerank 服务。
- 工具系统。
- 文件解析服务。
- 沙盒执行服务。

Go 的 interface 很适合做边界抽象。

例如向量检索：

```go
type Retriever interface {
    Search(ctx context.Context, req SearchRequest) ([]SearchResult, error)
}
```

工具执行：

```go
type Tool interface {
    Definition() ToolDefinition
    Execute(ctx context.Context, user User, args json.RawMessage) (ToolResult, error)
}
```

沙盒执行：

```go
type SandboxRunner interface {
    Run(ctx context.Context, req SandboxRequest) (*SandboxResult, error)
}
```

面试表达：

> 我会用 interface 把模型、检索、工具、沙盒这些外部能力抽象出来，业务编排层依赖接口，不依赖具体供应商。这样方便替换模型、切换向量库，也方便写 mock 做测试。

### 3.7 场景六：SSE 服务

Go 实现 SSE 的关键：

- 设置响应头。
- 使用 `http.Flusher`。
- 循环写 event。
- 每次写后 flush。
- 监听 context。

场景：

- 聊天流式输出。
- Agent 执行步骤。
- 文件解析进度。
- 沙盒运行日志摘要。

事件示例：

```text
event: message
data: {"delta":"你好"}

event: done
data: {"finish_reason":"stop"}
```

面试表达：

> Go 标准库就可以实现 SSE。重点不是循环写字符串，而是处理响应头、flush、断连、context 取消、done/error event 和完整内容落库。

### 3.8 场景七：RAG 后端

Go 在 RAG 中做编排：

```text
接收用户问题
  -> 鉴权
  -> 构造权限 filter
  -> 调 embedding
  -> 查向量库
  -> 查关键词索引
  -> rerank
  -> 构造 Prompt
  -> 调模型
  -> SSE 返回
  -> 记录引用和日志
```

Go 不一定负责计算 embedding，但可以调用 embedding 服务。

Go 适合负责：

- 权限过滤。
- 检索编排。
- top_k 控制。
- Prompt 构造。
- 日志。
- API 服务。

面试表达：

> RAG 里 Go 更适合做服务编排层，负责文件 metadata、权限过滤、检索调用、rerank 调度、Prompt 拼接、模型调用和引用日志。embedding 和解析可以调用专门服务。

### 3.9 场景八：文件解析和 embedding 异步任务

文件上传后不适合同步解析和 embedding。

Go 可以做：

- 上传接口。
- 创建文件记录。
- 投递 MQ。
- worker 消费任务。
- 调解析服务。
- 切片。
- 调 embedding。
- 写向量库。
- 更新状态。

任务状态：

```text
UPLOADED -> PARSING -> INDEXING -> READY
```

Go worker 要注意：

- 幂等。
- 重试。
- 死信。
- 任务状态。
- 失败原因。
- context 超时。

面试表达：

> 文件入库通常是异步管线。Go 服务上传后投递任务，worker 解析、切片、embedding、写向量库，并更新文件状态。任务要做幂等和失败重试，避免重复索引或删除后又写回。

### 3.10 场景九：Agent 编排

Agent 后端需要：

- 创建任务。
- 保存状态。
- 生成计划。
- 执行步骤。
- 调工具。
- 等待用户确认。
- 失败恢复。
- SSE 推送进度。

Go 适合做状态机和编排器。

任务结构：

```go
type AgentTask struct {
    ID          string
    Goal        string
    Status      string
    CurrentStep int
}
```

编排器：

```go
type AgentOrchestrator interface {
    Start(ctx context.Context, taskID string) error
    Cancel(ctx context.Context, taskID string) error
    Resume(ctx context.Context, taskID string) error
}
```

面试表达：

> Agent 在后端不是一个无限循环，而是任务状态机。Go 负责任务状态、步骤执行、工具调度、context 取消、SSE 进度、人工确认和审计日志。

### 3.11 场景十：MCP Server

Go 可以实现 MCP Server，把内部能力暴露给 AI 应用。

适合暴露：

- 知识库搜索。
- 订单查询。
- 工单查询。
- 文件读取。
- Git diff。
- 监控查询。

Go 后端要做：

- 工具注册。
- 参数 schema。
- JSON-RPC 消息处理。
- 内部 API 调用。
- 权限校验。
- 审计日志。

面试表达：

> 如果公司要把内部系统接入 Agent，我可以用 Go 实现 MCP Server，把有限、安全的业务能力包装成 tools/resources。Server 负责 tools/list、tools/call、参数校验、权限控制、内部 API 适配和审计。

### 3.12 场景十一：限流和成本控制

AI 后端限流不只看 QPS，还要看 token。

Go 服务可以结合 Redis 做：

- 用户级限流。
- 租户级限流。
- 模型级限流。
- RPM 限流。
- TPM 限流。
- 日预算。
- 月预算。

限流维度：

```text
user_id
tenant_id
model
request_type
estimated_tokens
```

面试表达：

> 大模型接口成本高，不能只做普通 QPS 限流。我会用 Redis 按用户、租户、模型做 RPM 和 TPM 控制，同时记录 input_tokens、output_tokens 和成本，超预算时降级或拒绝。

### 3.13 场景十二：日志和可观测

AI 后端日志要记录：

```text
trace_id
user_id
tenant_id
model
request_type
prompt_version
input_tokens
output_tokens
latency_ms
first_token_ms
retrieved_chunk_ids
tool_calls
status
error_code
```

Go 服务可以通过中间件统一注入：

- trace_id。
- user_id。
- request_id。
- latency。
- error。

面试表达：

> AI 后端排障依赖 trace。一次请求可能经过 RAG、模型、工具、SSE、Agent 多个环节，所以 Go 服务要把 trace_id 贯穿下去，记录 token、耗时、召回 chunk、工具调用和错误码。

### 3.14 常见误区

误区一：Go 在 AI 公司没用。

不对。Go 不一定训练模型，但非常适合做线上 AI 应用后端。

误区二：会 goroutine/channel 就够了。

不够。要能说出它们在 SSE、并发检索、任务 worker 里的具体场景。

误区三：interface 是八股。

在 AI 后端里，interface 对多模型、多向量库、多工具抽象很实用。

误区四：模型调用就是普通 HTTP。

还要考虑 context、流式、重试、限流、token、日志。

误区五：并发越多越好。

要控制并发、超时、错误和资源占用。

## 4. 考题、编程题与验收

### 4.1 概念题

1. Go 在 AI 后端里适合做什么？
2. 为什么模型调用必须传 `context.Context`？
3. channel 在流式输出中怎么用？
4. goroutine 在 RAG 中有哪些使用场景？
5. interface 如何抽象不同模型供应商？
6. Go 如何实现 SSE？
7. 文件 embedding 为什么适合异步 worker？
8. Go 如何做 Agent 状态机？
9. Go 如何实现 MCP Server？
10. AI 后端限流为什么要看 token？
11. 日志里为什么要记录 first_token_ms？
12. 为什么业务代码不应该直接依赖某个模型 SDK？

### 4.2 面试问答题

问题一：Go 在大模型应用后端里有什么典型使用场景？

参考回答：

> Go 适合做模型网关、SSE 流式服务、RAG 后端、文件解析和 embedding 异步任务编排、Agent 状态机、MCP Server、工具调用系统、限流计费和日志观测。它的并发、HTTP 服务和部署能力很适合 AI 应用的服务层和编排层。

问题二：你怎么用 Go 处理大模型流式输出？

参考回答：

> 我会让模型客户端把 stream chunk 转成只读 channel，HTTP handler 用 SSE 返回给前端。handler 用 select 同时监听 chunk 和 request context，收到 chunk 就写 event 并 flush，context done 就取消上游模型请求。结束时发送 done event，并保存完整回答和 usage。

问题三：如何用 Go 抽象多个模型供应商？

参考回答：

> 我会定义统一的 LLMClient 接口，包括 Chat 和 Stream 方法。OpenAI、DeepSeek、通义、私有模型分别实现 adapter。业务层只依赖接口，模型网关统一处理超时、重试、限流、usage 统计和日志。这样可以切换模型或做 fallback。

问题四：Go 在 RAG 系统里负责什么？

参考回答：

> Go 可以负责 RAG 的 API 和编排层，包括鉴权、权限 filter、调用 embedding、向量检索、关键词检索、rerank、Prompt 拼接、模型调用、SSE 返回、引用记录和日志。文件解析和 embedding 可以是异步任务或外部服务。

### 4.3 编程题

编程题一：设计统一模型客户端接口。

要求：

- 支持普通聊天。
- 支持流式聊天。
- 带 context。
- 响应包含 usage。

参考方向：

```go
type LLMClient interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}
```

验收标准：

- 能解释为什么需要 Stream。
- 能解释为什么 usage 要放响应里。
- 能解释为什么业务层依赖接口。

编程题二：写并发检索伪代码。

要求：

- 并发执行向量检索和关键词检索。
- 使用 context。
- 合并结果。
- 处理错误。

参考方向：

```go
g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
    vectorResults, err = vector.Search(ctx, req)
    return err
})

g.Go(func() error {
    keywordResults, err = keyword.Search(ctx, req)
    return err
})

if err := g.Wait(); err != nil {
    return nil, err
}

results := mergeAndDeduplicate(vectorResults, keywordResults)
```

编程题三：设计模型调用日志结构。

要求字段：

```text
trace_id
tenant_id
user_id
model
request_type
input_tokens
output_tokens
first_token_ms
latency_ms
status
error_code
created_at
```

验收标准：

- 能查成本。
- 能查慢请求。
- 能查首 token 慢。
- 能按租户统计。

### 4.4 系统设计题

题目：用 Go 设计一个 AI 应用后端服务。

需求：

```text
支持聊天、SSE 流式输出、RAG 检索、文件上传入库、工具调用和模型调用日志。
```

要求说明：

- 服务模块如何拆分。
- LLMClient 如何设计。
- SSE 如何实现。
- RAG 如何编排。
- 文件任务如何异步。
- 工具系统如何抽象。
- context 如何贯穿链路。
- Redis / MQ 用在哪里。
- 日志和限流怎么做。

回答必须包含：

- LLMClient interface。
- Retriever interface。
- Tool interface。
- SSE handler。
- context。
- channel。
- Redis 限流。
- MQ 异步任务。
- trace 日志。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能说出 Go 在 AI 后端中的 8 个典型场景。
- 能解释 context、channel、goroutine、interface 在 AI 后端中的用途。
- 能设计 LLMClient、Retriever、Tool 三个接口。
- 能讲清 Go 如何实现 SSE。
- 能讲清 Go 如何做 RAG 编排。
- 能回答“你 Go 只是会用，具体场景是什么”这个问题。

## 5. 和其他知识点的相关性

和 LLM 调用的关系：

- Go 的 LLMClient 是模型调用工程化的载体。
- context、HTTP Client、usage 日志都在这里用到。

和 SSE 的关系：

- Go 的 channel、Flusher、context 是 SSE 流式输出的关键。

和 RAG 的关系：

- Go 负责检索编排、权限过滤、Prompt 拼接和模型调用。

和文件管理的关系：

- Go 负责文件状态、异步任务、删除同步和索引编排。

和 Tool Calling 的关系：

- Go 的 Tool interface 和 registry 能承接工具系统。

和 MCP 的关系：

- Go 可以实现 MCP Server，把企业内部能力暴露给 AI Host。

和 Agent 的关系：

- Go 适合做 Agent 状态机、任务编排、工具调度和 SSE 进度。

和沙盒的关系：

- Go 负责创建沙盒任务、传递 context、收集结果和审计。

## 6. 演变过程与后续方向

### 6.1 从传统 Go Web 后端到 AI 应用后端

传统 Go 后端：

```text
HTTP API -> 业务逻辑 -> 数据库/缓存 -> 返回 JSON
```

AI Go 后端：

```text
HTTP API
  -> Prompt 构造
  -> RAG 检索
  -> 模型调用
  -> 工具执行
  -> SSE 返回
  -> token 成本和日志
```

底层工程能力没变，但接入对象和交互范式变了。

### 6.2 从单模型调用到模型网关

早期：

```text
业务代码直接调某个模型 API
```

后续：

```text
业务代码 -> Go 模型网关 -> 多模型供应商
```

模型网关会成为 AI 应用公司的基础设施。

### 6.3 从同步 API 到流式和长任务

普通接口一次返回。

AI 应用里会更多出现：

- SSE 长连接。
- Agent 长任务。
- 文件异步处理。
- 沙盒执行。
- 任务进度推送。

Go 的并发和 context 能力会更重要。

### 6.4 从业务接口到工具生态

传统后端提供 API 给前端和其他服务。

AI 后端还要把能力包装成：

- Tool。
- MCP Server。
- Agent Step。
- Resource。

后端接口会从“给人点按钮用”变成“给模型和 Agent 调用”。

### 6.5 后续方向

Go 在 AI 后端里的方向：

- 模型网关。
- AI API Gateway。
- Agent Runtime。
- MCP Server / MCP Gateway。
- RAG 服务。
- 文件知识库服务。
- 沙盒任务编排。
- AI Observability。
- 成本和计费系统。

## 7. 工作中典型场景

### 场景一：面试官问 Go 的 channel 用在哪里

可以回答：

> 在 AI 后端里，channel 很适合承接模型流式 chunk。模型客户端读取供应商 stream 后写入 channel，SSE handler 从 channel 读出并 flush 给前端，同时监听 context 取消。也可以用于任务进度事件，但要注意关闭 channel 和避免生产者阻塞。

### 场景二：面试官问 context 用在哪里

可以回答：

> context 会贯穿 HTTP 请求、模型调用、RAG 检索、工具调用和沙盒执行。用户取消、请求超时或服务关闭时，可以通过 context 让下游停止执行，避免资源泄漏和成本继续增长。

### 场景三：面试官问 goroutine 用在哪里

可以回答：

> RAG 混合检索里可以并发查向量库和关键词索引；Agent 可以并发执行无依赖工具；文件 embedding worker 可以并发处理任务。但要配合 errgroup、context 和并发限制，不能无限开 goroutine。

### 场景四：模型调用经常超时

Go 处理：

- 使用 context timeout。
- HTTP Client 设置 timeout。
- 限制 Prompt 长度。
- 对可重试错误做有限重试。
- fallback 到备用模型。
- 记录 latency 和 error_code。

### 场景五：RAG 查询很慢

Go 处理：

- 并发向量检索和关键词检索。
- 给检索设置超时。
- 记录每个阶段耗时。
- 控制 top_k。
- 缓存热门问题。
- rerank 超时降级。

### 场景六：面试官问“你 Go 只是会用，具体 AI 场景能说吗”

可以这样回答：

> 可以。比如我会用 Go 做模型网关，定义 LLMClient 抽象多个模型；用 context 控制模型超时和取消；用 channel 接模型 stream chunk，再通过 SSE 推给前端；用 goroutine 并发做 RAG 混合检索；用 interface 抽象 Retriever、Tool、SandboxRunner；用 MQ 做文件解析和 embedding 异步任务；用 Redis 做 RPM/TPM 限流；用 Go 状态机编排 Agent 任务。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

