# 第 08 章：MCP 协议与工具生态

## 1. 学习什么

本章学习 MCP，Model Context Protocol，模型上下文协议。它是当前大模型应用、Agent、IDE AI 助手、工具生态里非常重要的协议方向。

你需要掌握以下内容：

- MCP 是什么，解决什么问题。
- MCP 和普通 API、Tool Calling、Agent 框架有什么区别。
- Host、Client、Server 分别是什么。
- Tools、Resources、Prompts 三类能力分别是什么。
- MCP 的工具发现和工具调用流程。
- JSON-RPC 在 MCP 中扮演什么角色。
- Stdio、HTTP / Streamable HTTP 等传输方式如何理解。
- 后端工程师如何实现 MCP Server。
- MCP Server 如何暴露公司内部系统能力。
- MCP 的权限、安全、审计和沙盒风险。
- MCP 和 RAG、文件管理、Tool Calling、Agent 的关系。

本章学完后，你应该能说清楚：

> MCP 是让大模型应用以标准方式连接外部工具、数据源和上下文资源的协议。它不是模型本身，也不是简单 API，而是把工具、资源、Prompt 用统一协议暴露给 AI 应用。

## 2. 为什么学习

前面第 07 章讲了 Tool Calling：模型可以提出“我要调用某个工具，参数是什么”。

但如果每个产品、每个 Agent 框架、每个 IDE 都自己定义一套工具接入方式，就会出现问题：

- 工具不能复用。
- 每个平台都要重复接一遍。
- 工具描述格式不统一。
- 权限和审计方式分散。
- 本地文件、数据库、浏览器、Git、内部 API 难以标准化接入。

MCP 想解决的是：

```text
AI 应用如何标准化连接外部工具和上下文？
```

可以类比：

```text
HTTP 让系统之间能标准化通信
LSP 让编辑器和语言服务能标准化协作
MCP 让 AI 应用和工具/资源能标准化连接
```

在大模型公司或 AI 应用公司里，MCP 常见场景包括：

- IDE AI 助手读取项目文件。
- Agent 调用 Git、终端、数据库。
- 企业内部系统暴露成工具。
- 知识库检索暴露成 MCP Resource 或 Tool。
- 浏览器、文件系统、任务系统接入 Agent。
- 公司内部沉淀一批 MCP Server，供不同 AI 应用复用。

你不一定一开始就要写复杂 MCP 框架，但至少要能讲清它的架构、价值、调用流程和安全边界。

## 3. 知识详细内容

### 3.1 MCP 是什么

MCP 是 Model Context Protocol，模型上下文协议。

一句话理解：

> MCP 是 AI 应用连接外部工具、数据源和 Prompt 模板的一套标准协议。

它解决的不是“模型如何生成文本”，而是：

- 模型应用如何发现可用工具。
- 模型应用如何调用工具。
- 模型应用如何读取外部资源。
- 模型应用如何获取可复用 Prompt。
- 工具服务如何以统一方式暴露能力。

传统后端可能这样接能力：

```text
业务代码 -> HTTP Client -> 订单系统 API
```

MCP 场景更像：

```text
AI Host
  -> MCP Client
  -> MCP Server
  -> Tools / Resources / Prompts
  -> 公司内部系统
```

### 3.2 MCP 的核心架构

MCP 是 client-server 架构。

核心角色：

| 角色 | 说明 |
| --- | --- |
| Host | AI 应用本体，例如 IDE、桌面助手、Agent 平台 |
| Client | Host 内部负责连接某个 MCP Server 的客户端 |
| Server | 暴露 tools、resources、prompts 的服务 |

可以画成：

```text
AI Host
  ├─ MCP Client A -> MCP Server A -> 文件系统工具
  ├─ MCP Client B -> MCP Server B -> Git 工具
  └─ MCP Client C -> MCP Server C -> 企业知识库工具
```

重要点：

- Host 是用户直接使用的 AI 应用。
- Client 通常在 Host 内部。
- Server 是能力提供方。
- 一个 Host 可以连接多个 MCP Server。
- 一个 MCP Server 可以暴露多个工具和资源。

面试表达：

> MCP 里不是模型直接连数据库，而是 Host 通过 MCP Client 连接 MCP Server，由 Server 暴露受控工具和资源，后端仍然负责权限、参数校验和安全边界。

### 3.3 Tools、Resources、Prompts

MCP Server 主要暴露三类能力。

#### Tools

Tools 是可执行动作。

例如：

- 查询订单。
- 搜索知识库。
- 读取文件。
- 执行 Git diff。
- 调用内部 API。
- 查询数据库。
- 创建工单。

Tools 通常是 model-controlled，也就是模型可以根据上下文决定是否调用。

#### Resources

Resources 是可读取的上下文数据。

例如：

- 文件内容。
- 数据库记录。
- API 响应。
- 配置文档。
- 项目信息。
- 知识库资料。

Resources 更偏“给模型看什么”。

#### Prompts

Prompts 是可复用的提示词模板。

例如：

- 代码审查 Prompt。
- SQL 分析 Prompt。
- 故障排查 Prompt。
- 文档总结 Prompt。

Prompts 更偏“让模型如何做”。

三者关系：

```text
Tools：能做什么
Resources：能看什么
Prompts：怎么引导模型做
```

### 3.4 MCP 和 Tool Calling 的关系

Tool Calling 是模型调用工具的一种交互机制。

MCP 是工具、资源、Prompt 的标准暴露协议。

可以这样理解：

```text
Tool Calling：模型说“我要调用工具”
MCP：工具以什么标准方式被发现和调用
```

例如：

```text
MCP Server 暴露 get_order_status 工具
Host 通过 tools/list 发现这个工具
模型根据用户问题选择调用
Host 通过 tools/call 调用工具
Server 执行后返回结果
模型基于结果回答用户
```

所以 MCP 不是替代 Tool Calling，而是让 Tool Calling 的工具来源更标准化。

### 3.5 MCP 和普通 HTTP API 的区别

普通 HTTP API：

```text
调用方知道接口路径、参数、返回格式
业务代码主动调用接口
```

MCP：

```text
AI Host 可以发现 Server 暴露的工具和资源
工具描述给模型理解
模型可参与选择工具
消息遵循 MCP 协议
```

对比：

| 对比项 | 普通 HTTP API | MCP |
| --- | --- | --- |
| 调用者 | 业务代码 | AI Host / Agent |
| 能力发现 | 代码里写死或查文档 | tools/list、resources/list、prompts/list |
| 参数说明 | API 文档 | schema 和 description |
| 面向对象 | 程序员和服务 | AI 应用和工具生态 |
| 标准化目标 | 系统间接口 | 模型应用接工具和上下文 |

MCP Server 内部仍然可能调用 HTTP API。MCP 并不取代你的业务接口，而是在 AI 应用和业务系统之间加了一层标准适配。

### 3.6 JSON-RPC 是什么角色

MCP 消息使用 JSON-RPC 风格。

JSON-RPC 是一种远程过程调用格式，通常包含：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

响应：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": []
  }
}
```

通知没有 `id`，通常不需要响应：

```json
{
  "jsonrpc": "2.0",
  "method": "notifications/tools/list_changed"
}
```

后端要理解：

- request 有 id，需要 response。
- notification 没有 id，不需要 response。
- method 表示要执行的协议动作。
- params 是参数。
- result 是成功结果。
- error 是错误结果。

### 3.7 工具发现 tools/list

MCP Client 可以向 Server 请求工具列表。

请求：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
```

响应里包含工具数组：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "get_order_status",
        "title": "查询订单状态",
        "description": "根据订单 ID 查询物流和状态",
        "inputSchema": {
          "type": "object",
          "properties": {
            "order_id": {"type": "string"}
          },
          "required": ["order_id"]
        }
      }
    ]
  }
}
```

工具发现很关键，因为它让 Host 知道：

- Server 提供哪些能力。
- 每个工具什么时候用。
- 参数 schema 是什么。
- 模型可以如何选择工具。

### 3.8 工具调用 tools/call

当模型决定使用工具后，Host 通过 MCP Client 调用 Server。

请求：

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_order_status",
    "arguments": {
      "order_id": "12345"
    }
  }
}
```

Server 执行后返回：

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "订单 12345 已发货，预计明天送达。"
      }
    ]
  }
}
```

后端实现 MCP Server 时，重点是：

- 工具名必须匹配。
- 参数必须校验。
- 用户和租户权限必须校验。
- 执行结果要可被模型理解。
- 错误要结构化返回。

### 3.9 Resources

Resources 用来暴露可读取的数据。

例如：

```text
file:///project/README.md
database://customers/123
kb://policy/annual-leave
```

Resources 适合：

- 文件内容。
- 文档列表。
- 数据库记录。
- 当前项目上下文。
- 知识库材料。

资源的核心是“读取上下文”，而不是执行动作。

区别：

```text
Resource：给模型看资料
Tool：让模型触发动作
```

例如：

- 读取员工手册内容，可以是 Resource。
- 搜索员工手册，可以是 Tool。
- 删除员工手册，必须是高危 Tool。

### 3.10 Prompts

Prompts 是可复用提示词模板。

例如：

```text
review_code
summarize_document
generate_sql_explanation
incident_analysis
```

MCP Server 可以暴露 Prompt，让 Host 或用户选择使用。

对后端来说，Prompts 可以和前面第 03 章的 Prompt 模板管理结合：

- prompt_id。
- version。
- arguments。
- description。
- scenario。

企业内部可以沉淀一些标准 Prompt：

- 代码审查。
- 故障复盘。
- 合同风险分析。
- 客服话术总结。
- SQL 慢查询分析。

### 3.11 Transport 传输方式

MCP 是协议，不绑定单一传输方式。

常见传输：

| 传输 | 适合场景 |
| --- | --- |
| stdio | 本地进程，IDE 插件，本机工具 |
| HTTP / Streamable HTTP | 远程服务，企业内部工具平台 |

stdio 可以理解为：

```text
Host 启动一个本地 MCP Server 进程
通过标准输入输出交换 JSON-RPC 消息
```

远程 HTTP 可以理解为：

```text
Host 通过网络连接远程 MCP Server
Server 暴露公司内部能力
```

选型思路：

- 本地文件、Git、终端工具：常用 stdio。
- 企业知识库、订单系统、CRM：更适合远程 HTTP。
- 涉及敏感系统：必须考虑认证、授权、TLS、审计。

### 3.12 后端如何实现 MCP Server

实现 MCP Server，本质上是把已有后端能力包装成 MCP 能发现和调用的能力。

步骤：

```text
选择要暴露的能力
  -> 定义 Tool / Resource / Prompt
  -> 编写 name、description、inputSchema
  -> 实现 handler
  -> 参数校验
  -> 权限校验
  -> 调内部 API 或数据库
  -> 返回 MCP 格式结果
  -> 记录审计日志
```

例如暴露企业知识库搜索：

```text
Tool: search_knowledge_base
参数: query, top_k
执行: 调 RAG 检索服务
返回: 文档片段、来源、chunk_id
```

暴露订单查询：

```text
Tool: get_order_status
参数: order_id
执行: 调订单系统
返回: 状态、物流、预计送达
```

暴露文件资源：

```text
Resource: kb://doc/{doc_id}
执行: 读取有权限的文档内容
返回: 文本和 metadata
```

### 3.13 MCP 安全边界

MCP 很强，也很危险。

因为它可能让 AI 应用接触：

- 文件系统。
- 数据库。
- 内部 API。
- 终端命令。
- Git 仓库。
- 用户隐私数据。
- 企业敏感资料。

安全原则：

```text
最小权限
工具白名单
参数校验
用户授权
租户隔离
敏感字段脱敏
高危操作确认
审计日志
沙盒隔离
```

不要把整个系统能力无脑暴露给 MCP。

危险做法：

```text
暴露 execute_sql(sql string)
暴露 run_shell(command string)
暴露 read_file(path string) 且不限制目录
暴露 delete_file(path string) 且不确认
```

更安全的做法：

```text
暴露 query_order(order_id)
暴露 search_kb(query, top_k)
暴露 read_allowed_file(file_id)
暴露 run_test(test_name) 而不是任意 shell
```

### 3.14 MCP 和企业内部系统

传统企业内部系统很多：

- CRM。
- ERP。
- 工单系统。
- 订单系统。
- 知识库。
- 文档系统。
- 监控系统。
- CI/CD。
- Git 平台。

MCP Server 可以作为这些系统的 AI 适配层：

```text
AI Host -> MCP Server -> 内部系统 API
```

后端工程师的价值是：

- 知道内部系统 API。
- 知道权限边界。
- 知道数据模型。
- 知道哪些操作危险。
- 能把能力包装成模型可理解的工具。
- 能加审计、限流、告警。

### 3.15 常见误区

误区一：MCP 就是 Tool Calling。

不对。Tool Calling 是模型使用工具的交互机制，MCP 是工具、资源、Prompt 的标准连接协议。

误区二：MCP Server 就是普通 HTTP Server。

不完全对。MCP Server 要遵守协议消息、能力发现、工具调用、资源和 Prompt 的约定。

误区三：有了 MCP 就不用做权限。

大错。MCP 只定义连接方式，不替你完成业务权限和安全策略。

误区四：暴露越多工具越好。

工具越多，模型越容易选错，风险也越大。应该按场景暴露最小必要工具。

误区五：本地 MCP Server 就一定安全。

本地工具可能读文件、执行命令、访问密钥，同样需要限制范围和审计。

## 4. 考题、编程题与验收

### 4.1 概念题

1. MCP 是什么？解决什么问题？
2. Host、Client、Server 分别是什么？
3. Tools、Resources、Prompts 分别是什么？
4. MCP 和普通 HTTP API 有什么区别？
5. MCP 和 Tool Calling 有什么关系？
6. MCP 为什么使用 JSON-RPC 风格消息？
7. `tools/list` 和 `tools/call` 分别做什么？
8. stdio 传输适合什么场景？
9. 远程 MCP Server 要考虑哪些安全问题？
10. 为什么不能暴露任意 SQL 或任意 shell 工具？
11. MCP Server 如何接入企业知识库？
12. MCP 和 Agent 有什么关系？

### 4.2 面试问答题

问题一：你怎么理解 MCP？

参考回答：

> MCP 是模型上下文协议，用来标准化 AI 应用和外部工具、资源、Prompt 的连接方式。Host 是 AI 应用，Client 在 Host 内负责连接 Server，Server 暴露 Tools、Resources、Prompts。模型可以通过 Host 发现工具并请求调用，但真正执行仍由 MCP Server 和后端系统完成，后端要负责权限、参数校验、安全和审计。

问题二：MCP 和 Tool Calling 有什么区别？

参考回答：

> Tool Calling 是模型调用工具的机制，模型输出工具名和参数。MCP 是工具如何被 AI 应用发现、描述和调用的协议。可以理解为 Tool Calling 解决“模型要用工具”，MCP 解决“工具如何标准化接入 AI 应用”。

问题三：如果让你把公司订单系统接入 MCP，你怎么做？

参考回答：

> 我不会直接暴露任意订单 API，而是先选择安全的能力，比如 get_order_status。MCP Server 定义工具 name、description 和 inputSchema，执行时校验 order_id、用户身份、租户和数据权限，再调用订单系统 API。返回结果要脱敏，并记录 trace_id、user_id、tool_name、arguments_summary、permission_result、status 和 latency。

问题四：MCP 有哪些安全风险？

参考回答：

> MCP 可能暴露文件、数据库、内部 API 和命令执行能力，所以风险包括越权访问、Prompt Injection 诱导工具调用、敏感信息泄露、高危操作误执行、工具伪装和审计缺失。后端要做最小权限、工具白名单、参数校验、用户授权、租户隔离、敏感字段脱敏、高危操作确认、沙盒隔离和审计日志。

### 4.3 编程题

编程题一：设计一个 MCP Tool 定义。

要求：

- 工具名是 `search_knowledge_base`。
- 参数包含 `query` 和 `top_k`。
- `query` 必填。
- `top_k` 是整数，默认可由后端处理。

参考方向：

```json
{
  "name": "search_knowledge_base",
  "title": "搜索企业知识库",
  "description": "根据用户问题检索企业知识库，返回相关文档片段和来源。",
  "inputSchema": {
    "type": "object",
    "properties": {
      "query": {
        "type": "string",
        "description": "用户要检索的问题"
      },
      "top_k": {
        "type": "integer",
        "description": "返回的最大片段数量"
      }
    },
    "required": ["query"]
  }
}
```

验收标准：

- 能解释 name、description、inputSchema 的作用。
- 能解释为什么 query 必填。
- 能解释 top_k 为什么要限制范围。

编程题二：写一个 MCP 工具调用伪代码。

要求包含：

- 解析 `tools/call`。
- 找到工具。
- 校验参数。
- 校验权限。
- 执行内部服务。
- 返回 MCP result。
- 记录审计日志。

伪代码方向：

```text
handleToolsCall(request):
  name = request.params.name
  args = request.params.arguments
  tool = registry.get(name)
  if tool not found:
    return jsonrpc_error(MethodNotFound)
  validate(args, tool.inputSchema)
  user = auth(request)
  checkPermission(user, tool, args)
  result = tool.execute(user, args)
  auditLog(user, name, args, result)
  return jsonrpc_result(content=result)
```

编程题三：设计 MCP Server 暴露内部知识库的工具清单。

至少包含：

```text
search_knowledge_base(query, top_k)
get_document(doc_id)
list_user_documents(keyword)
```

验收标准：

- 每个工具都有明确用途。
- 每个工具都有参数 schema。
- 每个工具都说明权限边界。
- 不暴露任意 SQL 或任意文件路径。

### 4.4 系统设计题

题目：设计一个企业内部 MCP Server，让 AI 助手可以查询知识库和工单系统。

要求说明：

- MCP Server 暴露哪些 Tools。
- 是否暴露 Resources。
- 是否暴露 Prompts。
- 如何做工具发现。
- 如何处理工具调用。
- 如何接内部系统 API。
- 如何鉴权和租户隔离。
- 如何处理高危操作。
- 如何记录审计日志。
- 如何处理工具超时和错误。

回答必须包含：

- Host / Client / Server 架构。
- Tools / Resources / Prompts。
- JSON-RPC 请求响应。
- `tools/list`。
- `tools/call`。
- 参数 schema。
- 权限校验。
- 审计日志。
- 最小权限原则。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能解释 MCP 的核心价值。
- 能画出 Host、Client、Server 架构。
- 能区分 Tools、Resources、Prompts。
- 能说明 MCP 和 Tool Calling 的关系。
- 能写出 `tools/list` 和 `tools/call` 的基本结构。
- 能设计一个简单 MCP Server 的工具清单。
- 能说明 MCP 的主要安全风险和防护措施。

## 5. 和其他知识点的相关性

和 Tool Calling 的关系：

- Tool Calling 是模型调用工具的机制。
- MCP 是工具被发现和调用的标准协议。

和 Agent 的关系：

- Agent 需要使用工具和资源完成任务。
- MCP 可以成为 Agent 工具生态的标准连接层。

和 RAG 的关系：

- RAG 检索能力可以包装成 MCP Tool。
- 知识库文档可以包装成 MCP Resource。

和文件管理的关系：

- 文件系统可以通过 MCP 暴露资源和工具。
- 文件权限和路径限制必须在 MCP Server 中实现。

和 Prompt 的关系：

- MCP Server 可以暴露可复用 Prompt。
- Host 可以把 Prompt 与工具、资源组合给模型。

和 SSE 的关系：

- MCP 工具调用过程可以通过上层应用的 SSE 展示给前端。
- MCP 本身关注协议通信，上层产品仍要考虑用户体验流。

和沙盒的关系：

- 执行命令、读写文件、运行代码的 MCP Server 应该受沙盒限制。
- 沙盒是 MCP 工具执行的安全边界之一。

和 Go 的关系：

- Go 适合实现远程 MCP Server、工具注册、权限校验、内部 API 适配和审计日志。
- 本地 stdio MCP Server 也可以用 Go 写。

## 6. 演变过程与后续方向

### 6.1 从私有工具接入到标准协议

早期 AI 应用接工具，通常每个平台自己定义：

```text
工具描述格式
参数格式
调用方式
错误格式
权限策略
```

这导致工具难以复用。

MCP 的方向是：

```text
用统一协议连接 AI 应用和工具生态
```

### 6.2 从 API 调用到模型可发现工具

传统 API 是开发者写代码调用。

MCP 里的工具是 AI Host 可以发现、模型可以理解、后端可以执行的能力。

这要求工具描述不仅给人看，也要让模型能判断什么时候使用。

### 6.3 从单应用工具到工具生态

如果一个企业内部沉淀了多个 MCP Server：

- Git MCP Server。
- 知识库 MCP Server。
- 工单 MCP Server。
- 数据库 MCP Server。
- 监控 MCP Server。

不同 AI 助手就可以复用这些能力，而不是每个助手重写一遍。

### 6.4 从本地工具到远程企业服务

早期很多 MCP 用于本地 IDE 和文件系统。

企业落地会逐渐走向：

- 远程 MCP Server。
- 统一认证。
- 统一权限。
- 统一审计。
- 工具市场。
- 企业内部工具治理。

### 6.5 后续方向

MCP 后续会继续发展：

- 更成熟的远程认证和授权。
- 更丰富的工具生态。
- MCP Server 市场。
- 企业内部 MCP 网关。
- 工具调用安全扫描。
- 工具权限分级。
- 与 Agent 平台深度集成。
- 与沙盒、浏览器、IDE、数据库结合。

## 7. 工作中典型场景

### 场景一：把企业知识库接入 AI 助手

可以做一个 MCP Server：

```text
Tool: search_knowledge_base
Resource: kb://doc/{doc_id}
Prompt: summarize_policy
```

后端要做：

- 查询用户权限。
- 调 RAG 检索服务。
- 返回 chunk、source、doc_id。
- 限制 top_k。
- 记录审计。

### 场景二：IDE AI 助手需要读取项目文件

本地 MCP Server 可以暴露：

```text
list_project_files
read_file
search_code
git_diff
```

安全重点：

- 限制工作目录。
- 不读取密钥文件。
- 不允许任意路径读取。
- 记录敏感文件访问。

### 场景三：AI 助手要执行数据库查询

危险做法：

```text
execute_sql(sql)
```

更安全做法：

```text
get_customer_summary(customer_id)
search_orders(customer_id, date_range)
get_ticket_stats(filters)
```

原则：

- 不暴露任意 SQL。
- 只暴露业务语义工具。
- 参数 schema 固定。
- 做数据权限控制。

### 场景四：MCP 工具返回了敏感字段

处理方式：

- 工具层脱敏。
- 返回前过滤字段。
- 日志中只记录参数摘要。
- 对敏感工具加权限。
- 必要时要求用户确认。

### 场景五：模型选错 MCP 工具

可能原因：

- 工具描述太模糊。
- 工具太多。
- 工具命名相似。
- Prompt 没有说明工具边界。

处理方式：

- 优化 name 和 description。
- 按场景只开放必要工具。
- 给工具加更明确的参数 schema。
- 在系统 Prompt 里说明工具使用规则。

### 场景六：面试官问“你会怎么实现 MCP Server”

可以这样回答：

> 我会先选择有限、安全、明确的业务能力暴露成 MCP Tools，比如知识库搜索、订单查询，而不是任意 SQL 或 shell。Server 实现 tools/list 返回工具定义和 inputSchema，实现 tools/call 做参数校验、用户鉴权、租户和数据权限校验，然后调用内部 API，返回结构化结果。所有调用记录 trace、user、tool、参数摘要、权限结果、耗时和错误，涉及文件或命令执行的工具要放到沙盒或受限环境中。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

