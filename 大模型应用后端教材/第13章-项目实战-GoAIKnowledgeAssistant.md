# 第 13 章：项目实战 - Go AI Knowledge Assistant

## 1. 学习什么

本章把前面 12 章的知识收束成一个可以真正做、可以写进简历、可以在面试中讲清楚的项目：

```text
Go AI Knowledge Assistant
```

它是一个基于 Go 的 AI 知识库问答助手，覆盖 AI 应用后端最核心的能力：

- 模型 API 调用。
- SSE 流式输出。
- 文件上传。
- 文本解析和切片。
- 简易 RAG 检索。
- Prompt 构造。
- 会话管理。
- 模型调用日志。
- token 和耗时记录。
- 基础限流。
- 项目讲述和面试表达。

本章不是要求你一个月内做成商业级系统，而是做一个“面试能讲、代码能跑、链路完整”的项目。

本章学完后，你应该能说清楚：

> 我做了一个基于 Go 的 AI 知识库问答系统，支持大模型对话、SSE 流式输出、文件上传、文档切片、简易检索增强生成、模型调用日志和基础限流。这个项目用来证明我能把大模型能力工程化接入后端系统。

## 2. 为什么学习

只背概念是不够的。

面试官听到你说：

```text
我了解 RAG、SSE、Tool Calling、MCP、Agent。
```

可能会继续问：

- 你做过吗？
- SSE 怎么实现？
- RAG 链路怎么设计？
- 文件上传后怎么入库？
- Prompt 怎么拼？
- 模型调用日志记什么？
- Go 在里面负责什么？

如果你没有项目，只能背答案，很容易露怯。

一个小项目的价值是：

- 把概念串起来。
- 给简历增加 AI 相关经历。
- 面试时有真实链路可讲。
- 让你知道工程细节在哪里。
- 证明你不是只看过文章。

这个项目不追求功能大，而追求链路完整：

```text
用户提问
  -> 后端接收
  -> 构造 Prompt
  -> 调模型
  -> SSE 流式返回
```

再加上：

```text
文件上传
  -> 文本切片
  -> 简易检索
  -> 拼入 Prompt
  -> 模型基于资料回答
```

这已经足够覆盖 AI 应用后端的主干。

## 3. 知识详细内容

### 3.1 项目定位

项目名称：

```text
Go AI Knowledge Assistant
```

项目一句话：

> 一个基于 Go 的 AI 知识库问答后端，支持模型对话、SSE 流式输出、文件上传、简易 RAG 检索和模型调用日志。

项目目标：

- 让传统后端简历有 AI 项目。
- 让你能讲清 AI 应用后端链路。
- 覆盖面试高频点。
- 不陷入过度复杂。

项目不做：

- 不训练模型。
- 不做复杂前端。
- 不做复杂多租户。
- 不做完整向量数据库集群。
- 不做复杂 Agent 自动执行。
- 不做生产级权限系统。

你要避免“项目过大做不完”。这个项目的意义是打通链路。

### 3.2 功能范围

第一版必须做：

```text
1. 普通聊天接口
2. SSE 流式聊天接口
3. 文件上传接口
4. 文本文件解析
5. 文档切片
6. 简易检索
7. RAG Prompt 构造
8. 模型调用日志
9. 基础限流
10. 项目 README 和面试讲稿
```

可以后续增强：

```text
1. pgvector / Milvus / Qdrant
2. Embedding 模型
3. PDF / Word 解析
4. 用户登录
5. 多租户
6. Tool Calling
7. MCP Server
8. Agent 工作流
```

建议第一版先用关键词检索或简单内存检索模拟 RAG。这样你可以先把完整链路跑通，再替换成向量检索。

### 3.3 技术栈

建议技术栈：

```text
语言：Go
HTTP 框架：标准库 net/http 或 Gin
配置：环境变量 + config
存储：SQLite / MySQL / PostgreSQL 均可
缓存和限流：Redis 可选
模型：OpenAI compatible API / DeepSeek / 通义 / 本地兼容接口
文件存储：本地目录即可
检索：第一版关键词检索
前端：可选，curl 或简单 HTML 页面即可
```

如果你时间紧：

- 用 `net/http`。
- 文件存在本地。
- 元数据用 SQLite。
- 检索用关键词匹配。
- 模型 API 用 OpenAI compatible 格式。

面试时可以解释：

> 第一版为了控制复杂度，用关键词检索模拟 RAG 链路，后续可以把 Retriever 接口替换成 pgvector 或 Milvus。

这句话很重要，说明你知道工程演进。

### 3.4 项目模块设计

推荐模块：

```text
cmd/server：启动入口
internal/config：配置读取
internal/llm：模型客户端
internal/chat：聊天和 SSE
internal/file：文件上传和文档管理
internal/rag：切片和检索
internal/prompt：Prompt 构造
internal/usage：模型调用日志
internal/ratelimit：限流
```

核心接口：

```go
type LLMClient interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}

type Retriever interface {
    Search(ctx context.Context, req SearchRequest) ([]SearchResult, error)
}

type PromptBuilder interface {
    BuildRAGPrompt(question string, chunks []SearchResult) []Message
}
```

模块关系：

```text
HTTP Handler
  -> Chat Service
  -> Retriever
  -> PromptBuilder
  -> LLMClient
  -> Usage Logger
```

### 3.5 API 设计

建议接口：

```text
POST /api/chat
POST /api/chat/stream
POST /api/files
GET  /api/files
GET  /api/files/{id}
DELETE /api/files/{id}
GET  /api/usage
```

#### 普通聊天

```http
POST /api/chat
Content-Type: application/json

{
  "session_id": "s1",
  "question": "什么是 RAG？"
}
```

返回：

```json
{
  "answer": "RAG 是检索增强生成...",
  "usage": {
    "input_tokens": 100,
    "output_tokens": 200
  }
}
```

#### 流式聊天

```http
POST /api/chat/stream
Content-Type: application/json

{
  "session_id": "s1",
  "question": "解释一下 SSE"
}
```

SSE 事件：

```text
event: message
data: {"delta":"SSE"}

event: done
data: {"finish_reason":"stop"}
```

#### 文件上传

```http
POST /api/files
Content-Type: multipart/form-data

file=@knowledge.md
```

返回：

```json
{
  "file_id": "f1",
  "status": "READY",
  "chunk_count": 12
}
```

### 3.6 数据模型

文件表：

```text
file_id
original_name
storage_path
mime_type
size_bytes
status
created_at
deleted_at
```

chunk 表：

```text
chunk_id
file_id
chunk_index
content
created_at
```

会话表：

```text
session_id
created_at
updated_at
```

消息表：

```text
message_id
session_id
role
content
created_at
```

模型调用日志表：

```text
trace_id
session_id
model
request_type
input_tokens
output_tokens
latency_ms
first_token_ms
status
error_code
created_at
```

第一版可以简化，但字段要能支撑面试讲述。

### 3.7 RAG 第一版实现

第一版不要一上来接复杂向量库。

可以这样做：

```text
上传 txt / md
  -> 按段落切片
  -> chunk 存数据库
  -> 用户提问时做关键词匹配
  -> 取 top_k chunk
  -> 拼 Prompt
  -> 调模型回答
```

简易检索思路：

```text
把问题分词或按关键字拆分
遍历 chunk
计算命中次数
按分数排序
取 top_k
```

后续演进：

```text
Retriever interface
  -> KeywordRetriever
  -> VectorRetriever
  -> HybridRetriever
```

面试表达：

> 第一版我用关键词检索快速打通 RAG 链路，同时抽象了 Retriever 接口。后续可以无缝替换成 pgvector 或 Milvus，并接入 embedding 和 rerank。

### 3.8 Prompt 构造

RAG Prompt 模板：

```text
你是企业知识库助手。
请只基于提供的资料回答用户问题。
如果资料不足，请回答“资料不足，无法确认”。
回答后列出引用资料编号。

资料：
[资料1]
来源：fileA.md
内容：...

[资料2]
来源：fileB.md
内容：...

用户问题：
...
```

后端要记录：

- 本次用了哪些 chunk。
- chunk 对应哪个文件。
- Prompt 版本。
- 模型和 token。

### 3.9 SSE 实现重点

流式接口要做到：

- 设置 `Content-Type: text/event-stream`。
- 每个 chunk 写 `event: message`。
- 每次写后 `Flush()`。
- 监听 `r.Context().Done()`。
- 结束时发 `done`。
- 出错时发 `error`。
- 累计完整回答并落库。

面试讲点：

> 流式不是简单循环写字符串，还要处理断连、取消、flush、done/error 事件和完整内容保存。

### 3.10 限流和日志

基础限流：

```text
按 IP 或用户每分钟限制请求数
```

如果接 Redis：

```text
key = rate_limit:{user_id}:{minute}
INCR + EXPIRE
```

模型日志：

```text
trace_id
model
request_type
latency_ms
input_tokens
output_tokens
status
error_code
```

你要能说：

> 大模型接口有成本和限额，所以项目里记录了 usage，并预留了按用户和模型做限流、成本统计的扩展点。

### 3.11 项目分阶段实现

第一阶段：模型聊天。

```text
实现 /api/chat
封装 LLMClient
记录日志
```

第二阶段：SSE。

```text
实现 /api/chat/stream
支持 chunk 转发
处理断连取消
```

第三阶段：文件和切片。

```text
实现文件上传
保存文件 metadata
解析 txt/md
按段落切片
```

第四阶段：简易 RAG。

```text
实现 Retriever
检索 top_k chunk
拼 RAG Prompt
回答并返回引用
```

第五阶段：工程化。

```text
限流
错误处理
README
接口文档
项目讲稿
```

## 4. 考题、编程题与验收

### 4.1 概念题

1. 这个项目解决什么问题？
2. 为什么第一版可以用关键词检索模拟 RAG？
3. 为什么要抽象 LLMClient？
4. 为什么 SSE 接口要监听 context？
5. 文件上传后为什么要切片？
6. RAG Prompt 里为什么要放引用编号？
7. 模型调用日志要记录哪些字段？
8. 如何把关键词检索升级成向量检索？
9. 这个项目和真实企业知识库还有哪些差距？
10. 这个项目怎么写进简历？

### 4.2 面试问答题

问题一：介绍一下你的 AI 项目。

参考回答：

> 我做了一个 Go AI Knowledge Assistant，是一个基于 Go 的知识库问答后端。它支持普通聊天、SSE 流式输出、文件上传、文档切片、简易 RAG 检索、Prompt 构造、模型调用日志和基础限流。用户上传 txt 或 md 文件后，系统会切片入库，提问时先检索相关片段，再拼入 Prompt 调模型回答，并返回引用来源。

问题二：为什么你第一版没有直接上向量数据库？

参考回答：

> 我第一版的目标是先打通 RAG 完整链路，所以用关键词检索实现 Retriever，降低复杂度。代码上保留了 Retriever 接口，后续可以把实现替换成 pgvector、Milvus 或 Qdrant，再接 embedding 和 rerank。这样演进路径比较清晰。

问题三：SSE 流式输出你怎么做的？

参考回答：

> 模型客户端把 stream chunk 转成 channel，handler 设置 text/event-stream 响应头，循环读取 chunk，写成 message event 并 flush。handler 同时监听 request context，用户断开或取消时停止上游请求。结束时发送 done event，并把完整回答保存到会话。

问题四：项目里怎么记录模型调用？

参考回答：

> 每次模型调用生成 trace_id，记录 session、model、request_type、input_tokens、output_tokens、latency_ms、first_token_ms、status 和 error_code。这样可以排查慢请求、统计成本，也方便后续按用户或租户做限流和预算。

### 4.3 编程题

编程题一：实现 `LLMClient` 接口。

要求：

- 支持 Chat。
- 支持 Stream。
- 使用 context。
- 响应包含 usage。

验收标准：

- 能替换模型供应商。
- 能处理超时和取消。
- 能返回普通响应和流式 chunk。

编程题二：实现 `KeywordRetriever`。

要求：

- 输入 query。
- 遍历 chunks。
- 计算简单命中分数。
- 返回 top_k。

验收标准：

- 能按分数排序。
- 能返回 chunk 来源。
- 能替换为向量检索实现。

编程题三：实现 SSE handler。

要求：

- 设置响应头。
- 获取 Flusher。
- 读取 channel。
- 写 message event。
- 写 done event。
- 监听 context。

验收标准：

- curl 能看到流式输出。
- 用户断连后后端停止。
- 完整回答能保存。

编程题四：设计 usage 日志表。

要求字段：

```text
trace_id
session_id
model
request_type
input_tokens
output_tokens
latency_ms
status
error_code
created_at
```

验收标准：

- 能统计调用次数。
- 能统计 token 成本。
- 能排查失败请求。

### 4.4 项目验收清单

功能验收：

- [ ] `/api/chat` 能返回模型回答。
- [ ] `/api/chat/stream` 能流式返回。
- [ ] 上传 txt/md 文件后能生成 chunk。
- [ ] 提问时能检索相关 chunk。
- [ ] 回答能基于检索资料。
- [ ] 回答能返回引用来源。
- [ ] 模型调用有日志。
- [ ] 断连或取消能停止流式请求。
- [ ] README 能说明启动方式和接口。

面试验收：

- [ ] 能 2 分钟讲清项目背景。
- [ ] 能画出系统架构图。
- [ ] 能讲清一次 RAG 问答链路。
- [ ] 能讲清 SSE 实现。
- [ ] 能讲清 LLMClient 抽象。
- [ ] 能讲清后续如何升级向量库。
- [ ] 能讲清项目不足和改进方向。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能设计 Go AI Knowledge Assistant 的模块。
- 能写出核心接口。
- 能解释项目为什么适合后端转 AI。
- 能把项目写进简历。
- 能回答项目里的 RAG、SSE、日志、限流问题。
- 能说出项目后续演进方向。

## 5. 和其他知识点的相关性

和 LLM 调用的关系：

- 项目用 LLMClient 接模型 API。
- 记录 usage、latency、error。

和 Prompt 的关系：

- 项目需要构造普通聊天 Prompt 和 RAG Prompt。

和 SSE 的关系：

- 项目通过 SSE 实现流式输出。

和 RAG 的关系：

- 项目用简易检索打通 RAG 链路。
- 后续可升级向量库和 rerank。

和文件管理的关系：

- 项目支持文件上传、解析、切片和删除。

和 Go 的关系：

- 项目用 Go 实现 HTTP API、context、channel、interface、SSE 和异步任务扩展点。

和算法题的关系：

- 简易检索涉及字符串匹配、排序、top_k。

和面试表达的关系：

- 这是你从传统后端转 AI 应用后端最核心的项目讲述材料。

## 6. 演变过程与后续方向

### 6.1 从 Demo 到项目

最简单 demo：

```text
用户问题 -> 调模型 -> 返回回答
```

项目版：

```text
聊天
SSE
文件
切片
检索
Prompt
日志
限流
```

这说明你开始考虑工程化，而不是只调 API。

### 6.2 从关键词检索到向量检索

第一版：

```text
KeywordRetriever
```

升级：

```text
Embedding -> VectorRetriever -> HybridRetriever -> Rerank
```

这样项目可以持续演进。

### 6.3 从单用户到多租户

第一版可以不做复杂多租户。

后续可以加：

- user_id。
- tenant_id。
- 文件权限。
- chunk metadata。
- 检索 filter。

### 6.4 从知识库问答到 Agent

后续可以增加：

- Tool Calling。
- MCP Server。
- Agent 工作流。
- 沙盒文件处理。

例如：

```text
让 AI 查询知识库后自动生成总结报告
让 AI 调工具查看文件状态
让 AI 通过 MCP 暴露知识库搜索能力
```

### 6.5 后续方向

项目后续可以升级：

- 接 pgvector。
- 接真实 embedding。
- 支持 PDF / Word。
- 支持用户登录和权限。
- 支持多租户。
- 支持 Tool Calling。
- 实现 MCP Server。
- 增加 Agent 报告生成。
- 增加评测集。
- 增加成本统计 dashboard。

## 7. 工作中典型场景

### 场景一：面试官要求讲项目架构

你可以这样讲：

> 项目分成 API 层、LLMClient、RAG Retriever、PromptBuilder、FileService、UsageLogger 和 RateLimiter。用户提问后，ChatService 根据是否启用知识库决定直接调模型还是先检索 chunk，再构造 Prompt。普通接口一次返回，流式接口通过 SSE 返回。

### 场景二：面试官问 RAG 怎么做

回答：

> 上传文件后解析为文本，按段落切片，chunk 入库。提问时 Retriever 检索 top_k chunk，把 chunk 带来源编号拼进 Prompt，要求模型只能基于资料回答并返回引用。第一版用关键词检索，后续可替换为向量检索。

### 场景三：面试官问项目难点

可以说：

- SSE 断连和取消处理。
- RAG Prompt 构造。
- 文件切片和引用来源。
- LLMClient 抽象。
- usage 日志和限流。
- 后续向量检索演进。

### 场景四：面试官问项目不足

不要硬撑完美。

可以说：

> 当前版本是面试和学习项目，主要打通 AI 后端链路。它的不足是检索还只是关键词，没有接真实 embedding 和向量库；文件解析只支持 txt/md；权限和多租户较简单。后续我会用 Retriever 接口替换成 pgvector，并补充文件权限、PDF 解析和评测集。

### 场景五：简历写法

可以写：

```text
Go AI Knowledge Assistant：基于 Go 实现 AI 知识库问答后端，支持大模型对话、SSE 流式输出、文件上传、文档切片、简易 RAG 检索、Prompt 构造、模型调用日志和基础限流。通过 LLMClient、Retriever、PromptBuilder 等接口抽象模型调用和检索能力，支持后续扩展向量数据库和多模型供应商。
```

### 场景六：面试官问你为什么做这个项目

回答：

> 我从传统后端转 AI 应用后端，想证明自己不只是了解概念，而是能把模型能力接入业务系统。所以我选了知识库问答这个典型场景，覆盖模型调用、SSE、RAG、文件管理、日志和限流，这些都是 AI 应用后端的高频能力。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

