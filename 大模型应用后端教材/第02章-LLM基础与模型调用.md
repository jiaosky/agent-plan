# 第 02 章：LLM 基础与模型调用

## 1. 学习什么

本章学习大模型应用后端最基础、最高频的能力：如何理解 LLM，如何调用模型 API，如何把一次自然语言请求变成可控、可记录、可重试、可排查的后端调用。

你需要掌握以下内容：

- LLM 是什么，和传统 NLP 服务有什么不同。
- 模型调用的基本输入和输出。
- message、role、context window、token 的含义。
- temperature、top_p、max_tokens、stop 等常见参数。
- 普通响应和流式响应的区别。
- 后端如何封装模型客户端。
- 模型调用中的超时、重试、限流、降级。
- token 成本、调用日志、trace_id 如何记录。
- 为什么生产系统不能只写一个“调 API 的 demo”。

本章学完后，你应该能说清楚：

> LLM 调用不是简单的 HTTP 请求。生产系统需要处理消息结构、上下文长度、token 成本、超时取消、错误重试、限流降级、日志观测和多模型兼容。

## 2. 为什么学习

大模型应用后端的所有能力，最后都会落到一次或多次模型调用上。

无论你做的是：

- 普通 AI 聊天。
- 企业知识库问答。
- 智能客服。
- 代码助手。
- 数据分析 Agent。
- MCP 工具调用。
- 文档总结。
- 简历优化。
- SQL 生成。
- 自动工单处理。

底层都绕不开模型调用。

如果你只会写：

```text
把用户问题发给模型 -> 拿到回答 -> 返回前端
```

那只能说明你做过 demo，还不能说明你能做生产系统。

生产系统里，面试官和同事会关心：

- 模型超时怎么办？
- 用户点“停止生成”后，后端能不能取消请求？
- 模型返回 429 限流怎么办？
- 模型返回格式不稳定怎么办？
- 多轮对话历史太长怎么办？
- 一次请求用了多少 token？
- 哪个用户、哪个租户、哪个模型成本最高？
- 模型供应商挂了，是否有 fallback？
- 如何兼容 OpenAI、DeepSeek、通义、Claude 或私有模型？
- 日志里能不能查到这次回答的完整链路？

这些问题就是后端工程师的价值所在。

你不需要一开始懂模型内部每一层神经网络，但必须先懂模型作为一个外部智能服务时，后端如何接入、治理和观测它。

## 3. 知识详细内容

### 3.1 LLM 是什么

LLM 是 Large Language Model，大语言模型。它的核心能力是根据输入上下文预测和生成后续内容。

从应用后端角度看，可以先把 LLM 理解成一个特殊的远程服务：

```text
输入：一组消息、一些参数、可选工具、可选上下文
输出：文本、结构化 JSON、工具调用请求、流式增量、错误信息、token 用量
```

它和普通业务接口不同：

- 普通接口的输出通常确定。
- LLM 的输出具有概率性。
- 普通接口字段结构稳定。
- LLM 可能输出多余解释、格式错误、幻觉内容。
- 普通接口主要关注 QPS、延迟、错误码。
- LLM 还要关注 token、上下文长度、模型成本、回答质量。

所以 AI 后端不是把 LLM 当普通 HTTP 接口随便调一下，而是要在它外面加一层工程治理。

### 3.2 一次模型调用的基本结构

一次典型的模型调用包含：

```text
model：使用哪个模型
messages：消息列表
parameters：生成参数
tools：可选工具定义
stream：是否流式返回
metadata：业务侧追踪信息
```

可以简化成：

```json
{
  "model": "some-model",
  "messages": [
    {"role": "system", "content": "你是一个严谨的后端面试辅导老师"},
    {"role": "user", "content": "什么是 RAG？"}
  ],
  "temperature": 0.2,
  "max_tokens": 800,
  "stream": false
}
```

后端不要把用户输入直接裸传给模型，而要经过：

```text
用户输入
  -> 参数校验
  -> 敏感词或安全检查
  -> 读取会话历史
  -> 追加系统提示词
  -> 可选检索知识库
  -> 构造 messages
  -> 调用模型
  -> 校验输出
  -> 记录日志
  -> 返回结果
```

### 3.3 message 和 role

大模型聊天接口通常不是传一个字符串，而是传一个消息数组。

常见 role：

| Role | 含义 | 后端关注点 |
| --- | --- | --- |
| system | 系统级指令，定义模型行为边界 | 通常由平台控制，用户不能随便改 |
| developer | 开发者指令，定义应用规则 | 有些模型支持，用于放业务约束 |
| user | 用户输入 | 需要校验、审计、保存 |
| assistant | 模型历史回答 | 多轮对话时作为上下文 |
| tool | 工具执行结果 | Tool Calling 后回填给模型 |

一个多轮对话可能长这样：

```text
system: 你是企业知识库助手，只能基于资料回答
user: 我们公司的年假制度是什么？
assistant: 根据资料，公司年假...
user: 那试用期员工有吗？
```

模型没有真正“永久记忆”。你每次希望它知道什么，就要通过 messages、检索结果、摘要或外部状态传进去。

### 3.4 Context Window 上下文窗口

Context Window 是模型一次请求最多能处理的上下文长度。

它包括：

- system message。
- user message。
- assistant 历史消息。
- tool 结果。
- RAG 检索片段。
- 输出预留空间。

常见误区：

> 上下文窗口越大，就一定越好。

实际并不是。上下文越大，可能带来：

- 成本更高。
- 首 token 更慢。
- 无关信息干扰回答。
- 检索片段太多导致重点不清。
- 日志存储和排查更困难。

后端需要做上下文管理：

```text
保留关键系统提示词
  -> 裁剪过长历史
  -> 摘要旧对话
  -> 控制 RAG chunk 数量
  -> 给输出预留 token
```

### 3.5 Token 是什么

Token 是模型处理文本的基本计量单位。它不完全等于中文字符，也不完全等于英文单词。

后端需要关心 token，因为它影响：

- 是否超过上下文窗口。
- 输入成本。
- 输出成本。
- 响应延迟。
- 限流额度。
- 日志统计。

一次调用的成本通常可以粗略理解为：

```text
成本 = input_tokens * 输入单价 + output_tokens * 输出单价
```

在系统设计里，建议记录：

```text
trace_id
user_id
tenant_id
request_id
session_id
model
input_tokens
output_tokens
total_tokens
latency_ms
status
error_code
created_at
```

这些字段能回答老板和排障时最常见的问题：

- 哪个用户最费钱？
- 哪个租户调用最多？
- 哪类请求最慢？
- 哪个模型错误率最高？
- token 成本为什么突然上涨？

### 3.6 常见生成参数

| 参数 | 作用 | 使用建议 |
| --- | --- | --- |
| temperature | 控制随机性 | 面试、问答、代码类任务建议低一些 |
| top_p | 控制采样范围 | 一般不要和 temperature 同时乱调 |
| max_tokens | 限制最大输出长度 | 防止输出过长和成本失控 |
| stop | 遇到特定内容停止 | 可用于协议化输出场景 |
| stream | 是否流式返回 | 聊天产品通常开启 |

简单理解：

- 想要稳定、可复现、少跑偏：temperature 低一点。
- 想要创意、文案、头脑风暴：temperature 可以高一点。
- 想控制成本：设置 max_tokens。
- 想提升体验：使用 stream。

面试时不要说“temperature 越高越聪明”。它控制的是随机性，不是智商。

### 3.7 普通响应和流式响应

普通响应：

```text
请求 -> 等模型生成完整答案 -> 一次性返回
```

优点：

- 实现简单。
- 方便做完整结果校验。
- 适合短文本、分类、结构化输出。

缺点：

- 用户等待时间长。
- 长回答体验差。
- 中途取消不直观。

流式响应：

```text
请求 -> 模型边生成 -> 后端边推送 -> 前端边显示
```

优点：

- 首屏体验好。
- 适合聊天、长文生成、代码生成。
- 用户可以中途停止。

缺点：

- 后端要处理连接断开。
- 要处理 chunk 拼接。
- 要处理工具调用流。
- 日志要保存最终完整内容。

第 04 章会专门讲 SSE 和 Channel，这里先记住：聊天类大模型产品里，流式输出是基础能力。

### 3.8 后端如何封装模型客户端

不要在业务代码里到处直接写模型供应商的 HTTP 调用。

推荐抽象：

```go
type Role string

const (
    RoleSystem    Role = "system"
    RoleUser      Role = "user"
    RoleAssistant Role = "assistant"
    RoleTool      Role = "tool"
)

type Message struct {
    Role    Role
    Content string
}

type ChatRequest struct {
    Model       string
    Messages    []Message
    Temperature float64
    MaxTokens   int
    Stream      bool
    Metadata    map[string]string
}

type Usage struct {
    InputTokens  int
    OutputTokens int
    TotalTokens  int
}

type ChatResponse struct {
    Content      string
    FinishReason string
    Usage        Usage
}

type ChatChunk struct {
    Delta        string
    FinishReason string
}

type LLMClient interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}
```

这样做的好处：

- 业务代码不绑定某个供应商。
- 方便切换模型。
- 方便做统一日志。
- 方便做统一重试和限流。
- 方便后续加 RAG、Tool Calling、MCP。

### 3.9 context、超时和取消

模型调用可能慢，也可能卡住。所以 Go 后端必须使用 `context.Context`。

典型场景：

- 用户关闭页面，后端取消模型请求。
- 用户点击“停止生成”，后端中断流式调用。
- 模型超过 30 秒未返回，后端超时。
- 上游请求取消，下游模型调用也要停止。

伪代码：

```go
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()

resp, err := llm.Chat(ctx, req)
if err != nil {
    return err
}
```

面试表达：

> 模型调用比普通接口更容易慢或超时，所以请求链路必须传 context。用户取消、网关超时、服务关闭时，要能及时释放资源，避免 goroutine 和连接堆积。

### 3.10 错误处理、重试和降级

常见错误：

| 错误 | 含义 | 处理方式 |
| --- | --- | --- |
| 400 | 参数错误、消息格式错误 | 不重试，修正请求 |
| 401/403 | API Key 或权限问题 | 告警，不重试 |
| 408/504 | 超时 | 可重试或降级 |
| 429 | 限流 | 等待、退避、换模型 |
| 500/503 | 模型服务异常 | 重试、fallback |
| context canceled | 用户取消 | 不算系统错误 |

重试不能无脑做。

建议：

- 对 400、401、403 不重试。
- 对 429、500、503 可做有限重试。
- 使用指数退避。
- 设置最大重试次数。
- 记录每次重试日志。
- 流式输出已经开始后，重试要非常谨慎。

降级策略：

```text
主模型失败 -> 快速模型兜底
长上下文失败 -> 缩短上下文再试
RAG 失败 -> 提示暂时无法检索知识库
结构化输出失败 -> 重新要求模型按 schema 输出
```

### 3.11 限流和成本控制

模型接口比普通接口更贵，所以限流很重要。

常见维度：

- 按用户限流。
- 按租户限流。
- 按 API Key 限流。
- 按模型限流。
- 按 RPM 限流。
- 按 TPM 限流。
- 按日/月预算限额。

这里有两个常见缩写：

- RPM：Requests Per Minute，每分钟请求数。
- TPM：Tokens Per Minute，每分钟 token 数。

传统限流只看请求次数可能不够。因为一个请求可能很短，也可能带着几万 token 上下文。

更合理的是同时控制：

```text
请求次数 + token 数 + 预算金额
```

### 3.12 日志和观测

模型调用日志不要只记录“成功/失败”。

建议至少记录：

```text
trace_id：链路追踪 ID
user_id：用户
tenant_id：租户
session_id：会话
request_type：chat/rag/tool/agent
model：模型名
provider：供应商
input_tokens：输入 token
output_tokens：输出 token
latency_ms：耗时
first_token_ms：首 token 时间
status：成功/失败/取消
error_code：错误码
retry_count：重试次数
finish_reason：结束原因
```

如果涉及 RAG，还要记录：

```text
retrieved_doc_ids
retrieved_chunk_ids
top_k
rerank_used
```

如果涉及工具调用，还要记录：

```text
tool_name
tool_args_summary
tool_latency_ms
tool_status
```

这些日志的价值是：

- 排查慢请求。
- 统计成本。
- 发现模型错误率。
- 追踪幻觉来源。
- 证明回答引用了哪些资料。
- 审计工具调用是否越权。

### 3.13 生产调用链路

一个比较完整的模型调用链路：

```text
HTTP 请求进入
  -> 生成 trace_id
  -> 鉴权
  -> 限流
  -> 参数校验
  -> 读取会话历史
  -> 构造 messages
  -> 估算 token
  -> 选择模型
  -> 调用 LLMClient
  -> 处理响应或流式 chunk
  -> 保存会话
  -> 记录 usage 和日志
  -> 返回结果
```

面试时你只要能把这条链路讲清楚，就已经比“我调过 API”强很多。

## 4. 考题、编程题与验收

### 4.1 概念题

1. LLM 调用和普通 HTTP 接口调用有什么不同？
2. message 里的 system、user、assistant、tool 分别是什么？
3. 为什么说模型没有真正的长期记忆？
4. context window 和 token 有什么关系？
5. temperature 是控制什么的？是不是越高越聪明？
6. 普通响应和流式响应分别适合什么场景？
7. 为什么生产系统要记录 input_tokens 和 output_tokens？
8. 哪些模型错误适合重试，哪些不适合？
9. 为什么模型调用要使用 `context.Context`？
10. RPM 和 TPM 有什么区别？

### 4.2 面试问答题

问题一：你如何设计一个多模型兼容的调用层？

参考回答：

> 我会先定义统一的 LLMClient 接口，业务层只依赖 Chat 和 Stream 能力，不直接依赖某个供应商的 SDK。不同供应商各自实现适配器，在模型网关层统一处理鉴权、超时、重试、限流、日志和 usage 统计。这样后续可以按业务场景切换模型，也可以在主模型失败时做 fallback。

问题二：模型调用为什么要记录 token？

参考回答：

> token 影响上下文长度、接口成本、响应延迟和供应商限流。生产系统需要记录 input_tokens、output_tokens、total_tokens，按用户、租户、模型、请求类型统计成本。否则成本上涨时无法定位是历史消息太长、RAG 召回太多，还是某个用户异常调用。

问题三：模型返回 429 你怎么处理？

参考回答：

> 429 通常是限流。后端可以根据业务重要性做有限重试和指数退避，也可以排队、切换备用模型或提示用户稍后再试。同时要记录 provider、model、retry_count 和错误码。如果是租户自己的额度耗尽，还要返回明确的业务错误，而不是无限重试。

问题四：用户点击停止生成，后端应该怎么做？

参考回答：

> 前端停止后，后端要取消当前请求上下文，让下游模型连接尽快关闭。Go 里应该用 request context 或 context.WithCancel，把 ctx 传给模型客户端。服务端还要保存已经生成的部分内容或标记为 canceled，避免 goroutine、HTTP 连接和流式 channel 泄漏。

### 4.3 编程题

编程题一：实现一个模型客户端接口定义。

要求：

- 用 Go 定义 `LLMClient`。
- 支持普通聊天和流式聊天。
- 请求结构包含 model、messages、temperature、max_tokens。
- 响应结构包含 content、finish_reason、usage。
- 所有方法必须带 `context.Context`。

验收标准：

- 能解释为什么需要接口抽象。
- 能解释为什么要返回 usage。
- 能解释为什么 stream 返回的是 chunk。

编程题二：写一个模型调用包装函数。

要求：

- 生成 trace_id。
- 设置 30 秒超时。
- 调用 `LLMClient.Chat`。
- 记录耗时、token、错误。
- 对用户取消不做系统错误告警。

伪代码方向：

```go
func ChatWithLog(parent context.Context, client LLMClient, req ChatRequest) (*ChatResponse, error) {
    traceID := newTraceID()
    start := time.Now()

    ctx, cancel := context.WithTimeout(parent, 30*time.Second)
    defer cancel()

    resp, err := client.Chat(ctx, req)
    latency := time.Since(start)

    if err != nil {
        if errors.Is(ctx.Err(), context.Canceled) {
            logCanceled(traceID, latency)
            return nil, err
        }
        logError(traceID, err, latency)
        return nil, err
    }

    logSuccess(traceID, resp.Usage, latency)
    return resp, nil
}
```

编程题三：设计 token 成本统计表。

要求字段：

- 用户 ID。
- 租户 ID。
- 模型。
- 输入 token。
- 输出 token。
- 总 token。
- 耗时。
- 状态。
- 错误码。
- 创建时间。

验收标准：

- 能按用户查成本。
- 能按租户查成本。
- 能按模型查成本。
- 能统计错误率和平均耗时。

### 4.4 系统设计题

题目：设计一个模型网关服务。

要求说明：

- 如何统一接入多个模型供应商。
- 如何支持普通响应和流式响应。
- 如何做鉴权、限流、超时和重试。
- 如何记录 token 成本。
- 如何支持 fallback 模型。
- 如何让业务方不感知底层供应商差异。

回答必须包含：

- 统一请求和响应结构。
- `LLMClient` 或 Provider Adapter。
- usage 统计。
- trace 日志。
- 错误码映射。
- 限流策略。
- fallback 策略。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能画出一次模型调用的完整后端链路。
- 能解释 message、role、token、context window。
- 能说出 temperature、max_tokens、stream 的作用。
- 能设计一个简单的 `LLMClient` 接口。
- 能解释为什么要做超时、取消、重试、限流。
- 能说明模型调用日志需要记录哪些字段。
- 能回答“模型调用为什么不是简单调 API”。

## 5. 和其他知识点的相关性

和 Prompt 的关系：

- 模型调用的核心输入是 messages。
- 第 03 章会讲如何设计 system prompt、用户输入、上下文和输出格式。

和 SSE 的关系：

- 本章讲普通响应和流式响应的差异。
- 第 04 章会讲如何用 SSE 把流式 chunk 推给前端。

和 RAG 的关系：

- RAG 最终也要把检索结果拼进模型调用。
- 第 05 章会讲如何控制检索片段数量，避免上下文过长。

和文件管理的关系：

- 文件解析后的文本会成为模型上下文。
- 文件越多，越需要 token 控制、权限过滤和引用来源。

和 Tool Calling 的关系：

- Tool Calling 也是模型调用的一种输出形态。
- 模型可能不直接输出最终答案，而是输出“我要调用某个工具”。

和 MCP 的关系：

- MCP 可以把外部工具和资源接入模型应用。
- 但最后仍然要通过模型调用判断何时使用这些工具。

和 Agent 的关系：

- Agent 通常是多次模型调用 + 多次工具执行 + 状态管理。
- 模型调用层如果不稳定，Agent 会更不稳定。

和 Go 的关系：

- Go 的 `context`、`http.Client`、`channel`、`interface` 都会直接用于模型调用层。
- 第 11 章会进一步把这些能力场景化。

## 6. 演变过程与后续方向

### 6.1 从传统 NLP API 到 LLM API

传统 NLP API 多数是单点能力：

```text
分词、情感分析、关键词提取、文本分类、机器翻译
```

输入输出通常比较固定。

LLM API 更像通用语言推理引擎：

```text
聊天、总结、分类、改写、问答、代码生成、工具调用、规划任务
```

它能力更强，但输出也更不稳定，所以后端治理更重要。

### 6.2 从单模型调用到模型网关

早期 demo 通常直接调用一个模型：

```text
业务代码 -> OpenAI API
```

生产系统会演变成：

```text
业务代码 -> 模型网关 -> 多个模型供应商
```

模型网关负责：

- 统一协议。
- 统一鉴权。
- 统一日志。
- 统一限流。
- 统一成本统计。
- 统一错误处理。
- 模型路由和 fallback。

### 6.3 从同步返回到流式交互

早期接口常常等完整结果返回。

聊天产品普及后，流式输出成为标配：

```text
模型生成一点 -> 后端推一点 -> 前端显示一点
```

后续更复杂的是：

- 工具调用过程也流式展示。
- Agent 执行步骤也流式展示。
- 多个模型或工具并行执行后合并输出。

### 6.4 从只看 QPS 到同时看 token

传统后端主要关注 QPS、RT、错误率。

AI 后端还要关注：

- 输入 token。
- 输出 token。
- TPM。
- 上下文长度。
- 首 token 时间。
- 每次回答成本。
- 每个租户预算。

这会让后端的计费、限流和监控模型发生变化。

### 6.5 后续方向

模型调用层后续会继续演变：

- 多模型路由：按场景选择快模型、强模型、便宜模型。
- 语义缓存：相似问题复用回答或检索结果。
- 批处理：离线任务合并调用降低成本。
- 结构化输出：让模型输出更稳定的 JSON。
- 工具调用原生化：模型输出不只是文本，而是动作。
- 私有模型接入：企业内部模型和公有模型混合使用。
- 可观测增强：把模型输入、输出、检索、工具、成本串成完整 trace。

## 7. 工作中典型场景

### 场景一：业务方说“帮我接一个模型”

你不能只问 API Key。

你要追问：

- 用于什么场景？
- 是否要流式输出？
- 是否要结构化 JSON？
- 最大延迟能接受多少？
- 失败是否要重试？
- 是否有备用模型？
- 是否需要记录 token 成本？
- 用户输入和模型输出是否要入库？

### 场景二：模型调用越来越贵

排查方向：

- 历史消息是否无限追加。
- RAG 是否召回过多 chunk。
- max_tokens 是否设置过大。
- 是否有失败重试导致重复调用。
- 是否存在异常用户或租户。
- 是否所有场景都用了最贵模型。

处理方式：

- 裁剪或摘要历史上下文。
- 控制 RAG top_k。
- 设置 max_tokens。
- 按场景路由不同模型。
- 增加预算告警。
- 做语义缓存。

### 场景三：用户反馈回答等很久才出来

需要看：

- first_token_ms。
- total_latency_ms。
- 检索耗时。
- 模型供应商耗时。
- Prompt token 数。
- 是否启用了流式输出。

处理方式：

- 使用 SSE 流式返回。
- 缩短上下文。
- 并发执行检索。
- 换更快模型。
- 做超时降级。

### 场景四：模型偶尔返回格式错误

常见原因：

- Prompt 约束不够。
- 使用普通文本输出承载结构化数据。
- 没有做服务端校验。
- temperature 太高。
- 输出太长被截断。

处理方式：

- 使用 JSON Schema 或结构化输出能力。
- 服务端解析和校验。
- 失败后进行一次修复重试。
- 对下游返回明确错误。

### 场景五：供应商接口不稳定

处理方式：

- 设置合理超时。
- 有限重试。
- 熔断。
- fallback 到备用模型。
- 区分用户取消和系统错误。
- 记录 provider 级别错误率。

### 场景六：面试官问“你怎么封装大模型调用”

可以这样回答：

> 我会做一层模型网关或 LLM Client 抽象。业务层只传统一的 ChatRequest，不直接感知供应商。网关层负责 provider adapter、超时、取消、重试、限流、日志、token usage、错误码映射和 fallback。这样后续接 OpenAI、DeepSeek、通义或私有模型时，不需要改业务主流程。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

