# 第 07 章：Tool Calling 与 Function Calling

## 1. 学习什么

本章学习大模型应用从“只会回答问题”走向“能调用外部能力”的关键机制：Tool Calling 与 Function Calling。

你需要掌握以下内容：

- Tool Calling / Function Calling 是什么。
- 它和普通模型文本回答有什么区别。
- 模型如何根据用户意图选择工具。
- 工具定义、参数 schema、返回结果如何设计。
- 后端在工具调用中负责什么。
- 工具执行结果如何回填给模型。
- 工具调用失败、参数错误、权限不足时怎么处理。
- 为什么不能让模型直接操作数据库或执行命令。
- 工具白名单、参数校验、权限校验、审计日志如何设计。
- Tool Calling 和 RAG、MCP、Agent 的关系。

本章学完后，你应该能说清楚：

> Tool Calling 的本质是让模型提出“我需要调用哪个工具以及参数是什么”，真正执行工具的是后端。后端必须负责工具注册、参数校验、权限控制、执行、错误处理、结果回填和审计。

## 2. 为什么学习

如果大模型只能回答文本，它能做的事情有限。

用户问：

```text
帮我查一下订单 12345 的物流状态。
```

模型自己并不知道订单系统里的实时数据。它需要调用后端工具：

```text
get_order_shipping_status(order_id=12345)
```

用户说：

```text
帮我把这个客户的工单状态改成已解决。
```

模型不能凭空修改业务系统。它需要后端提供工具，并且必须经过权限校验：

```text
update_ticket_status(ticket_id, status)
```

Tool Calling 让模型从“聊天机器人”变成“能使用工具的助手”。

典型场景：

- 查订单。
- 查库存。
- 查用户信息。
- 生成 SQL。
- 调内部 API。
- 创建工单。
- 发送通知。
- 查询知识库。
- 读取文件。
- 执行代码。
- 调用搜索。
- 操作日历。
- 生成报表。

面试里，Tool Calling 会和 Agent、MCP、安全、权限一起被问。你需要清楚：模型只是决策者之一，不是执行者，更不是权限系统。

## 3. 知识详细内容

### 3.1 Tool Calling 是什么

Tool Calling 是模型在生成回答时，不直接输出最终答案，而是输出一个“工具调用请求”。

普通问答：

```text
用户：订单 12345 到哪了？
模型：我不知道你的订单信息。
```

Tool Calling：

```text
用户：订单 12345 到哪了？
模型：需要调用 get_order_status，参数 order_id=12345
后端：执行 get_order_status
工具结果：已发货，预计明天送达
模型：你的订单已发货，预计明天送达。
```

完整链路：

```text
用户问题
  -> 后端把可用工具列表发给模型
  -> 模型判断是否需要工具
  -> 模型输出 tool_call
  -> 后端校验工具名和参数
  -> 后端校验用户权限
  -> 后端执行工具
  -> 后端把工具结果回填给模型
  -> 模型生成最终回答
```

### 3.2 Function Calling 和 Tool Calling 的区别

很多资料里会混用 Function Calling 和 Tool Calling。

可以这样理解：

| 概念 | 理解方式 |
| --- | --- |
| Function Calling | 早期常见说法，强调模型调用函数 |
| Tool Calling | 更通用说法，工具可以是函数、API、检索器、数据库、文件系统、浏览器等 |

Function Calling 更像：

```text
模型选择一个函数并生成参数
```

Tool Calling 更像：

```text
模型选择一种外部能力并生成调用请求
```

对后端来说，重点不是纠结名字，而是理解模式：

> 模型只生成调用意图和参数，后端负责真实执行和安全控制。

### 3.3 工具定义包含什么

一个工具定义通常包含：

```text
name：工具名
description：工具用途说明
parameters：参数 schema
required：必填参数
return：返回结果说明
permission：权限要求
timeout：超时时间
side_effect：是否有副作用
```

示例：

```json
{
  "name": "get_order_status",
  "description": "查询订单当前物流状态。只用于用户询问订单状态时。",
  "parameters": {
    "type": "object",
    "properties": {
      "order_id": {
        "type": "string",
        "description": "订单 ID"
      }
    },
    "required": ["order_id"]
  }
}
```

description 很重要。模型会根据工具描述判断什么时候使用工具。

写得太模糊，模型可能乱用；写得太窄，模型可能不用。

### 3.4 参数 schema 为什么重要

模型生成的参数不一定可靠。

可能出现：

- 缺少必填字段。
- 字段类型错误。
- 枚举值不存在。
- ID 格式不合法。
- 参数越权。
- 参数被 prompt injection 影响。

所以工具参数必须有 schema，并且后端要二次校验。

例如：

```text
ticket_id 必须是当前租户下的工单
status 只能是 open / processing / resolved
用户必须有工单编辑权限
```

不要因为模型生成了参数，就直接执行。

面试表达：

> Tool schema 只是让模型更容易生成正确参数，真正的参数校验和权限判断必须在后端完成。

### 3.5 后端在 Tool Calling 中负责什么

后端职责包括：

- 注册工具。
- 生成模型可理解的工具描述。
- 接收模型 tool_call。
- 校验工具是否存在。
- 校验参数 schema。
- 校验用户权限。
- 执行工具。
- 处理超时和错误。
- 脱敏工具结果。
- 把结果回填给模型。
- 记录审计日志。
- 控制高危操作是否需要人工确认。

后端不是“模型说什么就做什么”。

正确原则：

```text
模型负责提出调用建议
后端负责判断是否允许执行
后端负责真实执行
后端负责记录和兜底
```

### 3.6 工具调用结果如何回填

工具执行完成后，后端需要把结果作为 tool message 或上下文回填给模型。

例如：

```text
tool: get_order_status
result:
{
  "order_id": "12345",
  "status": "shipped",
  "eta": "2026-04-28"
}
```

然后模型基于工具结果生成自然语言回答：

```text
你的订单 12345 已发货，预计 2026-04-28 送达。
```

为什么不直接把工具结果返回给用户？

有些场景可以直接返回，但大多数时候模型需要：

- 把技术字段转成人话。
- 总结多个工具结果。
- 结合用户问题回答。
- 解释下一步建议。

但也要注意：工具结果如果包含敏感字段，回填给模型前要脱敏。

### 3.7 工具有无副作用

工具分为无副作用和有副作用。

无副作用工具：

- 查询订单。
- 查询天气。
- 搜索知识库。
- 读取文件。
- 查询库存。

有副作用工具：

- 修改工单。
- 删除文件。
- 发送邮件。
- 扣款。
- 创建订单。
- 执行命令。

对有副作用工具要更谨慎：

- 严格权限校验。
- 参数白名单。
- 二次确认。
- 审计日志。
- 幂等设计。
- 必要时人工审批。

面试表达：

> 查询类工具可以相对自动化，但写操作、删除操作、发消息、扣款、执行命令这类有副作用工具必须做权限、确认和审计，不能让模型直接执行。

### 3.8 工具错误处理

工具可能失败：

- 参数错误。
- 权限不足。
- 业务数据不存在。
- 下游接口超时。
- 下游服务异常。
- 工具返回数据过大。
- 工具结果无法解析。

处理方式：

| 错误 | 处理 |
| --- | --- |
| 参数错误 | 不执行，要求模型或用户补充 |
| 权限不足 | 返回明确拒绝，不让模型绕过 |
| 数据不存在 | 告诉模型未查询到 |
| 超时 | 可重试或提示稍后再试 |
| 服务异常 | fallback 或返回失败 |
| 高危操作 | 要求用户确认 |

工具错误也可以回填给模型，让模型用自然语言解释：

```text
工具 get_order_status 返回：订单不存在。
```

模型回答：

```text
没有查询到订单 12345，请确认订单号是否正确。
```

### 3.9 工具调用循环

有些任务需要多次工具调用。

例如：

```text
用户：帮我看看这个客户最近的订单和工单，判断是否需要重点跟进。
```

可能链路：

```text
查客户信息
  -> 查订单
  -> 查工单
  -> 查最近沟通记录
  -> 模型总结
```

这已经接近 Agent 工作流。

要控制：

- 最大工具调用次数。
- 最大执行时间。
- 是否允许并发工具。
- 是否允许调用高危工具。
- 每一步是否记录 trace。
- 工具之间的数据如何传递。

否则模型可能陷入循环：

```text
调用工具 -> 看结果 -> 又调用工具 -> 一直不结束
```

### 3.10 工具权限设计

工具权限至少包括三层：

#### 用户权限

用户是否有资格执行这个动作。

例如：

```text
普通客服只能查询订单，不能退款。
```

#### 数据权限

用户是否能访问这条数据。

例如：

```text
只能查自己负责客户的工单。
```

#### 工具权限

当前应用或 Agent 是否开放这个工具。

例如：

```text
知识库助手不允许调用 delete_file。
```

工具执行前要检查：

```text
tool 是否存在
tool 是否对当前场景开放
用户是否有工具权限
用户是否有数据权限
参数是否合法
是否需要二次确认
```

### 3.11 审计日志

Tool Calling 必须可审计，尤其是有副作用工具。

建议记录：

```text
trace_id
user_id
tenant_id
session_id
tool_name
tool_call_id
arguments_summary
permission_result
execution_status
latency_ms
error_code
side_effect
confirmation_required
created_at
```

对敏感参数要脱敏：

```text
不要明文记录密码、token、身份证、银行卡等敏感信息。
```

审计日志的作用：

- 排查模型为什么调用工具。
- 追踪谁触发了操作。
- 证明权限校验是否执行。
- 回滚或补偿错误操作。
- 满足企业合规要求。

### 3.12 Go 后端工具系统设计

可以设计一个工具接口：

```go
type ToolDefinition struct {
    Name        string
    Description string
    Parameters  JSONSchema
    SideEffect  bool
}

type ToolCall struct {
    ID        string
    Name      string
    Arguments json.RawMessage
}

type ToolResult struct {
    ToolCallID string
    Content    string
    Error      string
}

type Tool interface {
    Definition() ToolDefinition
    Execute(ctx context.Context, user User, args json.RawMessage) (ToolResult, error)
}
```

工具注册表：

```go
type ToolRegistry interface {
    Register(tool Tool)
    Get(name string) (Tool, bool)
    List() []ToolDefinition
}
```

执行流程：

```text
模型返回 tool_call
  -> registry 找工具
  -> schema 校验参数
  -> 权限校验
  -> Execute
  -> 记录审计
  -> 返回 ToolResult
```

### 3.13 常见误区

误区一：Tool Calling 是模型直接调用后端接口。

不对。模型只生成调用意图，后端才是真正执行者。

误区二：schema 有了就不需要校验。

不对。schema 帮模型生成参数，后端仍要严格校验。

误区三：工具越多越好。

工具太多会让模型选择困难，也增加安全风险。工具要按场景开放。

误区四：查询工具和写操作工具一样处理。

不对。有副作用工具必须更严格。

误区五：工具错误直接返回 500。

工具错误需要区分参数错误、权限错误、业务不存在、下游异常，并以模型能理解的方式回填或返回。

## 4. 考题、编程题与验收

### 4.1 概念题

1. Tool Calling 是什么？
2. Function Calling 和 Tool Calling 有什么区别？
3. 模型在工具调用里负责什么，后端负责什么？
4. 工具定义通常包含哪些内容？
5. 为什么工具参数需要 schema？
6. 为什么后端仍然要二次校验参数？
7. 什么是有副作用工具？
8. 高危工具为什么要二次确认？
9. 工具执行结果如何回填给模型？
10. 工具审计日志应该记录哪些字段？
11. 为什么不能让模型直接执行 SQL 或 shell 命令？
12. 如何避免模型无限循环调用工具？

### 4.2 面试问答题

问题一：你怎么理解 Tool Calling？

参考回答：

> Tool Calling 是让模型在需要外部能力时输出工具调用请求，包括工具名和参数。模型负责判断是否需要工具，后端负责校验工具名、参数、权限并执行真实工具。工具结果再回填给模型，由模型生成最终回答。生产系统里关键是工具白名单、参数校验、权限控制、错误处理和审计。

问题二：如果模型要调用删除文件工具，你会怎么处理？

参考回答：

> 删除文件属于有副作用且高危操作。后端不能因为模型请求就直接执行。需要确认该工具是否在当前场景开放，用户是否有删除权限，文件是否属于当前租户和用户可操作范围，并要求用户二次确认。执行后要删除文件记录、chunk、向量索引和缓存，并记录审计日志。

问题三：工具参数是模型生成的，可以直接信任吗？

参考回答：

> 不能。模型生成的参数可能缺字段、类型错误、枚举越界，也可能被用户诱导生成越权参数。schema 只能帮助模型生成更规范的参数，后端必须做 JSON schema 校验、业务校验和权限校验。

问题四：工具调用失败后怎么处理？

参考回答：

> 要区分失败类型。参数错误可以让用户补充或让模型修正；权限不足要明确拒绝；业务数据不存在可以回填给模型说明未找到；下游超时或异常可以重试、降级或提示稍后再试。所有失败都要记录 trace 和 tool_call_id。

### 4.3 编程题

编程题一：设计工具接口。

要求：

- 工具能返回定义。
- 工具能执行。
- 执行时带 context 和用户信息。
- 参数使用 JSON。

参考方向：

```go
type Tool interface {
    Definition() ToolDefinition
    Execute(ctx context.Context, user User, args json.RawMessage) (ToolResult, error)
}
```

验收标准：

- 能解释为什么 Execute 要带 user。
- 能解释为什么要有 Definition。
- 能解释为什么 args 不能直接信任。

编程题二：写一个工具执行伪代码。

要求包含：

- 查找工具。
- 参数 schema 校验。
- 权限校验。
- 执行工具。
- 记录审计。
- 返回结果。

伪代码方向：

```text
executeToolCall(user, toolCall):
  tool = registry.get(toolCall.name)
  if tool not found:
    return error("unknown tool")
  validateSchema(tool.definition.parameters, toolCall.arguments)
  checkToolPermission(user, tool)
  checkDataPermission(user, toolCall.arguments)
  result = tool.execute(user, toolCall.arguments)
  auditLog(user, toolCall, result)
  return result
```

编程题三：设计工具审计表。

要求字段：

```text
id
trace_id
tenant_id
user_id
session_id
tool_call_id
tool_name
arguments_summary
permission_result
status
error_code
latency_ms
side_effect
created_at
```

验收标准：

- 能按 trace 查工具调用链。
- 能查某个用户执行过哪些工具。
- 能查失败率和耗时。
- 能审计高危操作。

### 4.4 系统设计题

题目：设计一个支持工具调用的 AI 客服助手。

要求说明：

- 有哪些工具。
- 工具如何注册。
- 工具描述如何提供给模型。
- 模型返回 tool_call 后后端如何处理。
- 如何校验参数。
- 如何控制用户和数据权限。
- 如何处理查询类工具和写操作工具。
- 如何把工具结果回填给模型。
- 如何记录审计日志。
- 如何处理工具失败。

回答必须包含：

- Tool Registry。
- Tool Definition。
- JSON Schema。
- 权限校验。
- Side Effect 标记。
- Tool Result 回填。
- 审计日志。
- 最大工具调用次数。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能解释 Tool Calling 的完整链路。
- 能区分模型负责什么、后端负责什么。
- 能设计一个 Tool 接口和 Tool Registry。
- 能说明工具参数为什么必须校验。
- 能说明高危工具为什么要二次确认。
- 能设计工具审计日志。
- 能回答支持工具调用的 AI 客服系统设计题。

## 5. 和其他知识点的相关性

和 LLM 调用的关系：

- Tool Calling 是模型调用的一种输出形态。
- 模型可能输出文本，也可能输出工具调用请求。

和 Prompt 的关系：

- Prompt 和工具描述会影响模型何时选择工具。
- 工具结果回填后，也需要 Prompt 指导模型如何总结。

和 SSE 的关系：

- 工具调用过程可以通过 SSE 展示给前端。
- 例如正在查询订单、正在更新工单、工具执行完成。

和 RAG 的关系：

- RAG 可以被封装成一个检索工具。
- Tool Calling 更偏动作，RAG 更偏知识召回，但两者经常组合使用。

和文件管理的关系：

- 读取文件、删除文件、生成文件都可以是工具。
- 文件工具必须遵守文件权限和沙盒边界。

和 MCP 的关系：

- MCP 是把工具、资源和 Prompt 标准化暴露给模型应用的协议。
- Tool Calling 是模型使用这些工具的交互机制。

和 Agent 的关系：

- Agent 本质上经常是多轮模型调用 + 多次工具调用 + 状态管理。
- Tool Calling 是 Agent 的核心执行能力。

和沙盒的关系：

- 执行代码、运行命令、读写文件这类工具通常需要沙盒。
- 沙盒提供工具执行的安全边界。

## 6. 演变过程与后续方向

### 6.1 从文本回答到函数调用

早期模型主要输出文本。

后来为了让模型接入外部系统，出现 Function Calling：

```text
模型输出函数名和参数
后端执行函数
模型基于结果回答
```

这让模型可以查实时数据、调用业务系统，而不是只靠训练知识。

### 6.2 从 Function Calling 到 Tool Calling

Function Calling 更像调用函数。

Tool Calling 范围更广：

- API。
- 数据库查询。
- 文件系统。
- 搜索。
- RAG 检索。
- 浏览器。
- 代码执行。
- 内部平台。

模型不只是调用函数，而是在使用外部工具完成任务。

### 6.3 从单工具到工具编排

简单场景只调用一个工具。

复杂场景需要多工具：

```text
查客户 -> 查订单 -> 查工单 -> 生成跟进建议
```

这会逐渐演变成 Agent Workflow。

### 6.4 从私有工具系统到 MCP

不同框架最初都有自己的工具定义和调用方式。

MCP 的方向是把工具、资源和 Prompt 标准化，让不同 AI 应用更容易接入外部能力。

Tool Calling 是模型“要用工具”，MCP 是工具“如何被标准化提供”。

### 6.5 后续方向

Tool Calling 后续会继续发展：

- 更稳定的结构化参数。
- 更细粒度权限控制。
- 工具市场。
- 自动工具发现。
- 多工具并发执行。
- 人工审批工作流。
- 工具调用评测。
- 沙盒化执行。
- MCP 生态集成。

## 7. 工作中典型场景

### 场景一：AI 客服要查询订单

工具：

```text
get_order_status(order_id)
```

后端要做：

- 校验 order_id。
- 校验用户是否能查该订单。
- 调订单系统。
- 脱敏返回结果。
- 回填模型。
- 记录审计。

### 场景二：AI 助手要修改工单状态

这是有副作用操作。

后端要做：

- 校验工单归属。
- 校验用户编辑权限。
- 校验 status 枚举。
- 可能要求用户确认。
- 执行更新。
- 记录操作日志。

### 场景三：模型生成了错误参数

例如：

```text
status = "done"
```

但系统只支持：

```text
open / processing / resolved
```

处理方式：

- schema 校验失败。
- 不执行工具。
- 返回参数错误。
- 让模型询问用户或修正参数。

### 场景四：模型想调用不存在的工具

处理方式：

- registry 查不到工具。
- 返回 unknown_tool。
- 不执行任何操作。
- 记录日志。
- 必要时让模型重新选择可用工具。

### 场景五：工具返回大量数据

例如查询用户最近 1000 条订单。

处理方式：

- 工具层分页。
- 只返回摘要。
- 限制最大结果数。
- 敏感字段脱敏。
- 必要时让模型要求用户缩小范围。

### 场景六：面试官问“怎么保证工具调用安全”

可以这样回答：

> 我会把工具调用当成后端受控执行，而不是模型直接执行。工具必须注册在白名单里，参数按 schema 校验，再做业务校验和权限校验。对写操作、删除、发消息、执行命令这类有副作用工具，要求二次确认或人工审批。所有工具调用记录 trace、用户、租户、参数摘要、权限结果、执行状态和耗时，便于审计。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

