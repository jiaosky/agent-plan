# 第 09 章：Agent 范式与工作流编排

## 1. 学习什么

本章学习大模型应用公司里非常常见、但也最容易被说玄的概念：Agent。

Agent 不是“套一个很酷的名词”，也不是简单写一个 `while` 循环反复调模型。对后端工程师来说，Agent 更像一个由模型、工具、状态、规则和流程组成的任务执行系统。

你需要掌握以下内容：

- Agent 是什么，和普通 Chatbot 有什么区别。
- 单轮 LLM 调用、RAG、Tool Agent、Workflow Agent 的区别。
- Planner、Executor、Reviewer 分别是什么。
- 单 Agent 和多 Agent 的适用场景。
- Agent 如何使用工具、RAG、MCP、沙盒。
- Agent 状态如何保存和恢复。
- 如何控制最大步骤数、超时、成本和失败。
- 人类确认 Human-in-the-loop 为什么重要。
- Agent 后端如何做日志、审计、回放和可观测。
- 工作中如何设计一个可控的 Agent，而不是不可预测的自动化黑盒。

本章学完后，你应该能说清楚：

> Agent 是围绕目标进行多步推理和工具执行的系统。后端要负责状态管理、工具编排、权限控制、执行边界、失败恢复、日志观测和人工确认，而不是让模型无限自由行动。

## 2. 为什么学习

大模型应用早期多是聊天：

```text
用户问 -> 模型答
```

后来加入 RAG：

```text
用户问 -> 检索知识库 -> 模型基于资料答
```

再后来加入 Tool Calling：

```text
用户问 -> 模型调用工具 -> 模型基于工具结果答
```

Agent 更进一步：

```text
用户给目标
  -> Agent 拆解任务
  -> 选择工具
  -> 执行步骤
  -> 观察结果
  -> 调整计划
  -> 继续执行
  -> 输出最终结果
```

例如用户说：

```text
帮我分析最近一周客户投诉，找出前三类问题，并生成一份处理建议。
```

这可能需要：

- 查询工单系统。
- 聚类投诉内容。
- 查知识库。
- 生成统计表。
- 生成处理建议。
- 输出报告。

这就不是一次模型调用能稳定完成的，而是一个多步骤工作流。

面试和工作中，Agent 常被用于：

- 代码修复助手。
- 数据分析助手。
- 运维排障助手。
- 企业办公助手。
- 客服工单处理。
- 智能报表生成。
- 自动化测试。
- 文档生成。
- 浏览器操作。
- 沙盒代码执行。

你需要理解：Agent 强大，但风险也高。真正可落地的 Agent，核心不是“自主性越强越好”，而是“在可控边界内完成任务”。

## 3. 知识详细内容

### 3.1 Agent 是什么

可以先用一句话理解：

> Agent 是能围绕一个目标，结合模型推理、工具调用、状态记忆和执行反馈，完成多步骤任务的系统。

普通 Chatbot：

```text
输入一句话 -> 输出一句话
```

Agent：

```text
输入一个目标
  -> 计划
  -> 执行
  -> 观察
  -> 调整
  -> 再执行
  -> 完成
```

一个典型 Agent 循环：

```text
Goal
  -> Plan
  -> Act
  -> Observe
  -> Reflect
  -> Next Act
  -> Final Answer
```

但生产系统里不能让它无限循环，所以需要后端控制：

- 最大步骤数。
- 最大执行时间。
- 最大 token 成本。
- 可用工具范围。
- 高危操作确认。
- 失败终止条件。

### 3.2 Agent 和普通 Chatbot 的区别

| 对比项 | Chatbot | Agent |
| --- | --- | --- |
| 目标 | 回答问题 | 完成任务 |
| 调用次数 | 通常一次或少量 | 多次模型调用和工具调用 |
| 状态 | 对话历史为主 | 任务状态、步骤状态、工具结果 |
| 工具使用 | 可选 | 通常是核心能力 |
| 风险 | 回答错误 | 可能执行错误操作 |
| 后端重点 | 模型调用、SSE、会话 | 工作流、工具、权限、审计、恢复 |

面试表达：

> Chatbot 偏问答，Agent 偏任务执行。Agent 不只是模型回答，还要有工具、状态、计划、执行反馈和边界控制。

### 3.3 几种常见 Agent 范式

#### Tool-using Agent

模型根据用户问题调用工具。

```text
用户问订单 -> 调订单工具 -> 回答
```

适合：

- 查询业务系统。
- 简单操作。
- 客服助手。

#### Workflow Agent

流程比较固定，模型只在部分节点参与。

```text
分类 -> 检索 -> 生成 -> 校验 -> 返回
```

适合：

- 生产系统。
- 稳定业务流程。
- 需要可控结果。

#### Planner-Executor Agent

Planner 负责拆任务，Executor 负责执行。

```text
Planner: 先查数据，再分析，再生成报告
Executor: 按步骤调用工具完成
```

适合：

- 多步骤任务。
- 数据分析。
- 代码修改。

#### Reviewer Agent

执行后再由 Reviewer 检查质量。

```text
生成答案 -> 检查是否符合资料 -> 修正
```

适合：

- 代码审查。
- 文档生成。
- 高质量问答。

#### Multi-agent

多个 Agent 扮演不同角色。

```text
Planner + Executor + Reviewer + Reporter
```

适合复杂任务，但工程成本更高。

初学和面试建议重点掌握：

```text
Tool-using Agent
Workflow Agent
Planner-Executor
Human-in-the-loop
```

### 3.4 Planner、Executor、Reviewer

#### Planner

负责拆解目标。

输入：

```text
用户目标、可用工具、约束
```

输出：

```text
步骤列表、每步目的、需要工具
```

例如：

```text
1. 查询最近一周投诉工单
2. 按问题类型聚类
3. 查询处理知识库
4. 生成改进建议
```

#### Executor

负责执行步骤。

它可以：

- 调工具。
- 查 RAG。
- 读文件。
- 写文件。
- 调模型总结。

#### Reviewer

负责检查结果。

检查：

- 是否完成目标。
- 是否遗漏步骤。
- 是否违反规则。
- 是否需要补充工具调用。
- 输出格式是否正确。

后端可以把三者做成不同 Prompt，也可以做成同一个 Agent 的不同阶段。

### 3.5 Workflow Agent 更适合生产

很多人一提 Agent，就想到“完全自主”。

但生产系统通常更喜欢 Workflow Agent。

例如智能客服处理工单：

```text
用户输入
  -> 意图分类
  -> 是否需要查订单
  -> 调订单工具
  -> 是否需要查知识库
  -> 生成回答
  -> 敏感操作拦截
  -> 返回用户
```

这里流程是后端控制的，模型只负责某些判断和生成。

优点：

- 可控。
- 易排查。
- 容易测试。
- 成本可控。
- 安全边界清楚。

面试表达：

> 生产 Agent 不一定追求完全自主。很多业务更适合 Workflow Agent，由后端控制主流程，模型负责分类、生成、工具选择等局部智能，这样稳定性和安全性更好。

### 3.6 Agent 状态管理

Agent 不只是聊天历史，还需要任务状态。

状态可能包括：

```text
task_id
user_id
tenant_id
goal
current_step
step_list
tool_results
intermediate_files
status
cost
created_at
updated_at
```

状态枚举：

```text
CREATED
PLANNING
RUNNING
WAITING_USER_CONFIRMATION
FAILED
CANCELED
COMPLETED
```

为什么状态重要？

- 长任务需要恢复。
- 用户刷新页面后要继续查看。
- 人工确认需要暂停。
- 失败后要知道卡在哪一步。
- 审计要能回放完整链路。

### 3.7 工具编排

Agent 的能力很大程度来自工具。

后端要管理：

- 当前 Agent 可用工具。
- 每个工具的权限。
- 每步可调用哪些工具。
- 工具调用顺序。
- 工具超时和重试。
- 工具结果如何进入下一步。

一个工具编排流程：

```text
模型提出 tool_call
  -> 后端检查工具是否允许
  -> 校验参数
  -> 校验权限
  -> 执行工具
  -> 保存 tool_result
  -> 回填模型或进入下一步
```

Agent 不能拥有所有工具。应该按场景配置工具集合：

```text
客服 Agent：订单查询、工单查询、知识库搜索
代码 Agent：读文件、写文件、运行测试、git diff
数据 Agent：查表、生成 SQL、执行只读查询、生成图表
```

### 3.8 Human-in-the-loop 人类确认

不是所有操作都能自动执行。

需要人工确认的场景：

- 删除文件。
- 修改数据库。
- 发送邮件。
- 提交代码。
- 退款。
- 关闭工单。
- 执行命令。
- 对外发布内容。

流程：

```text
Agent 生成操作计划
  -> 后端判断为高危操作
  -> 暂停任务
  -> 展示给用户确认
  -> 用户确认后继续
  -> 用户拒绝则取消或修改计划
```

状态：

```text
WAITING_USER_CONFIRMATION
```

面试表达：

> Agent 不是越自动越好。对有副作用或高风险操作，应该引入 human-in-the-loop，让模型提出建议，用户确认后由后端执行。

### 3.9 失败恢复

Agent 执行中可能失败：

- 模型调用失败。
- 工具超时。
- 参数错误。
- 权限不足。
- 检索不到资料。
- 沙盒命令失败。
- 输出格式错误。
- 达到最大步骤数。

处理策略：

| 失败 | 处理 |
| --- | --- |
| 模型超时 | 重试或降级模型 |
| 工具超时 | 重试、跳过或终止 |
| 参数错误 | 要求模型修正或问用户 |
| 权限不足 | 明确拒绝 |
| 资料不足 | 返回无法确认 |
| 达到步骤上限 | 总结已完成内容并停止 |
| 高危操作 | 暂停等待确认 |

关键是不能让 Agent 无限制自我修复，否则成本和风险不可控。

### 3.10 Agent 可观测

Agent 的日志要比普通聊天更详细。

建议记录：

```text
trace_id
task_id
user_id
tenant_id
agent_type
goal
step_index
step_name
model
prompt_version
tool_name
tool_args_summary
tool_status
input_tokens
output_tokens
latency_ms
status
error_code
cost
```

前端也可以展示执行过程：

```text
正在分析任务
正在查询工单
正在检索知识库
正在生成报告
等待用户确认
任务完成
```

这通常通过 SSE 事件流实现。

### 3.11 Agent 后端架构

一个可落地的 Agent 后端可以拆成：

```text
Agent API：创建任务、查询状态、取消任务
Agent Orchestrator：主编排器
Planner：生成计划
Executor：执行步骤
Tool Registry：工具注册和权限
State Store：保存任务状态
Event Stream：SSE 推送步骤进度
Audit Log：审计工具和高危操作
Sandbox：隔离执行代码或命令
```

核心流程：

```text
createTask(goal)
  -> save task
  -> planner generates steps
  -> executor runs step by step
  -> each step may call tools
  -> update state
  -> emit events
  -> complete or fail
```

### 3.12 Go 后端如何实现 Agent 编排

可以设计任务结构：

```go
type AgentTask struct {
    ID          string
    UserID      string
    TenantID    string
    Goal        string
    Status      string
    CurrentStep int
    Steps       []AgentStep
}

type AgentStep struct {
    Index  int
    Name   string
    Status string
    Tool   string
    Result string
}
```

编排器接口：

```go
type AgentOrchestrator interface {
    Start(ctx context.Context, taskID string) error
    Cancel(ctx context.Context, taskID string) error
    Resume(ctx context.Context, taskID string) error
}
```

Go 适合做：

- 任务状态机。
- 工具调度。
- context 取消。
- SSE 事件推送。
- 异步任务。
- 超时控制。
- 审计日志。

模型负责智能判断，Go 后端负责控制边界和执行可靠性。

### 3.13 常见误区

误区一：Agent 就是 while 循环调模型。

不对。生产 Agent 要有状态、工具、边界、失败处理、日志和安全控制。

误区二：Agent 越自主越好。

不一定。生产系统更看重可控、稳定、可审计。

误区三：多 Agent 一定比单 Agent 强。

不一定。多 Agent 成本高、延迟高、排查难。简单任务用 Workflow Agent 更合适。

误区四：让模型决定所有流程。

核心业务流程最好由后端控制，模型负责局部智能。

误区五：Agent 错了只是回答错。

不对。如果 Agent 能调用工具，错误可能变成真实业务操作事故。

## 4. 考题、编程题与验收

### 4.1 概念题

1. Agent 是什么？
2. Agent 和普通 Chatbot 有什么区别？
3. Tool-using Agent 和 Workflow Agent 有什么区别？
4. Planner、Executor、Reviewer 分别负责什么？
5. 为什么生产系统更偏好 Workflow Agent？
6. Agent 状态需要保存哪些信息？
7. Human-in-the-loop 适合哪些场景？
8. Agent 如何避免无限循环？
9. Agent 日志为什么比普通聊天更复杂？
10. 多 Agent 有什么优缺点？
11. Agent 和 MCP 有什么关系？
12. Agent 和沙盒有什么关系？

### 4.2 面试问答题

问题一：你怎么理解 Agent？

参考回答：

> Agent 是围绕目标进行多步骤执行的系统。它会使用模型进行计划和判断，调用工具获取信息或执行动作，保存任务状态，根据工具结果调整下一步，最终完成任务。后端要负责工作流编排、状态管理、工具权限、失败恢复、人工确认、日志审计和成本控制。

问题二：生产系统里你会怎么设计 Agent？

参考回答：

> 我倾向先做 Workflow Agent，而不是完全开放式 Agent。后端控制主流程，比如分类、检索、工具调用、生成、校验，每个节点明确输入输出。模型负责局部判断和生成，工具调用经过白名单、参数校验和权限校验。高危操作进入人工确认，所有步骤记录 trace，便于回放和排查。

问题三：如何防止 Agent 无限调用工具？

参考回答：

> 可以设置最大步骤数、最大工具调用次数、最大执行时间、最大 token 成本和最大重试次数。每一步要保存状态和结果，失败时按错误类型终止、重试或等待用户输入。不能让模型无限自我修复。

问题四：哪些场景需要 human-in-the-loop？

参考回答：

> 有副作用或高风险操作都需要人工确认，比如删除文件、修改数据库、退款、关闭工单、发送邮件、提交代码、执行命令、对外发布内容。模型可以生成计划和参数，但后端暂停任务，用户确认后再执行。

### 4.3 编程题

编程题一：设计 Agent 任务状态。

要求包含：

- CREATED。
- PLANNING。
- RUNNING。
- WAITING_USER_CONFIRMATION。
- FAILED。
- CANCELED。
- COMPLETED。

参考方向：

```go
type AgentTaskStatus string

const (
    TaskCreated      AgentTaskStatus = "CREATED"
    TaskPlanning     AgentTaskStatus = "PLANNING"
    TaskRunning      AgentTaskStatus = "RUNNING"
    TaskWaitingUser  AgentTaskStatus = "WAITING_USER_CONFIRMATION"
    TaskFailed       AgentTaskStatus = "FAILED"
    TaskCanceled     AgentTaskStatus = "CANCELED"
    TaskCompleted    AgentTaskStatus = "COMPLETED"
)
```

验收标准：

- 能解释为什么需要 WAITING_USER_CONFIRMATION。
- 能解释 CANCELED 和 FAILED 的区别。
- 能解释为什么长任务必须保存状态。

编程题二：写一个 Agent 执行伪代码。

要求包含：

- 创建任务。
- 生成计划。
- 遍历步骤。
- 调用工具。
- 保存状态。
- 发送事件。
- 处理失败。
- 处理人工确认。

伪代码方向：

```text
runAgent(task):
  task.status = PLANNING
  steps = planner.plan(task.goal)
  saveSteps(task, steps)
  task.status = RUNNING

  for step in steps:
    if exceedsLimit(task):
      fail(task, "limit_exceeded")
      return
    if step.requiresConfirmation:
      task.status = WAITING_USER_CONFIRMATION
      emitEvent("waiting_confirmation", step)
      return
    result = executor.execute(step)
    saveStepResult(step, result)
    emitEvent("step_completed", step)

  task.status = COMPLETED
  emitEvent("completed", task)
```

编程题三：设计 Agent 日志表。

要求字段：

```text
id
trace_id
task_id
user_id
tenant_id
step_index
step_name
model
tool_name
status
input_tokens
output_tokens
latency_ms
error_code
created_at
```

验收标准：

- 能按 task_id 回放执行过程。
- 能统计失败步骤。
- 能统计工具耗时。
- 能统计 token 成本。

### 4.4 系统设计题

题目：设计一个企业工单分析 Agent。

需求：

```text
用户输入：分析最近一周客户投诉，找出主要问题并生成处理建议。
```

要求说明：

- Agent 如何生成计划。
- 需要哪些工具。
- 如何查询工单系统。
- 如何调用知识库 RAG。
- 如何保存任务状态。
- 如何通过 SSE 返回进度。
- 哪些步骤需要人工确认。
- 如何处理工具失败。
- 如何记录审计日志。
- 如何控制成本和最大步骤数。

回答必须包含：

- Planner。
- Executor。
- Tool Registry。
- State Store。
- SSE Event。
- Human-in-the-loop。
- 最大步骤数。
- trace 日志。

### 4.5 自测验收标准

学完本章后，你应该能做到：

- 能解释 Agent 和 Chatbot 的区别。
- 能说清 Tool Agent、Workflow Agent、Planner-Executor。
- 能设计 Agent 任务状态机。
- 能说明为什么要 human-in-the-loop。
- 能设计 Agent 执行日志。
- 能回答企业工单分析 Agent 的系统设计题。
- 能说出 Agent 的风险和控制手段。

## 5. 和其他知识点的相关性

和 LLM 调用的关系：

- Agent 通常包含多次 LLM 调用。
- 每次调用都要考虑 token、超时、重试和日志。

和 Prompt 的关系：

- Planner、Executor、Reviewer 往往使用不同 Prompt。
- Agent 的每一步都需要构造上下文。

和 SSE 的关系：

- Agent 执行过程适合用 SSE 展示步骤事件。
- 用户需要看到“正在查工单”“正在生成报告”等进度。

和 RAG 的关系：

- Agent 经常需要检索知识库辅助决策。
- RAG 可以作为 Agent 的一个工具或上下文来源。

和 Tool Calling 的关系：

- Tool Calling 是 Agent 执行动作的基础。
- Agent 可以多轮选择和调用工具。

和 MCP 的关系：

- MCP 可以给 Agent 提供标准化工具和资源。
- Agent 平台可以连接多个 MCP Server 扩展能力。

和文件管理的关系：

- Agent 可能读取、生成、修改文件。
- 文件权限和版本管理仍然要由后端控制。

和沙盒的关系：

- Code Agent、文件处理 Agent、命令执行 Agent 都需要沙盒隔离。
- 沙盒控制 Agent 能接触哪些文件、命令、网络和资源。

## 6. 演变过程与后续方向

### 6.1 从 Chatbot 到 Agent

最早的大模型应用多是聊天问答。

随着工具调用、RAG、MCP 和沙盒能力成熟，应用开始从“回答问题”变成“完成任务”。

这就是 Agent 兴起的基础。

### 6.2 从开放式 Agent 到 Workflow Agent

早期 Agent 常强调自主规划和自动执行。

生产落地后，大家发现完全开放式 Agent 风险高、成本高、不可控。

所以很多场景转向 Workflow Agent：

```text
后端控制主流程
模型负责局部智能
工具执行受控
高危操作人工确认
```

### 6.3 从单 Agent 到多 Agent

复杂任务可能拆成多个角色：

- Planner。
- Researcher。
- Executor。
- Reviewer。
- Reporter。

多 Agent 可以提升分工清晰度，但也带来：

- 更多 token 成本。
- 更长延迟。
- 更难调试。
- 更复杂状态管理。

所以不是所有任务都需要多 Agent。

### 6.4 从短任务到长任务

Agent 会从短对话逐渐走向长任务：

- 代码重构。
- 数据分析报告。
- 自动测试。
- 运维排障。
- 文档整理。

长任务要求：

- 可暂停。
- 可恢复。
- 可查看进度。
- 可人工确认。
- 可失败重试。
- 可审计回放。

### 6.5 后续方向

Agent 后续会继续发展：

- 更可靠的工具编排。
- 更成熟的 Agent 状态机。
- 多 Agent 协作平台。
- Agent 评测体系。
- Agent 安全沙盒。
- 人机协同审批流。
- 企业内部 Agent 平台。
- Agent 与 MCP 工具生态深度结合。

## 7. 工作中典型场景

### 场景一：产品说“我们要做一个 Agent”

你需要追问：

- Agent 要完成什么任务？
- 是否必须自动执行？
- 需要哪些工具？
- 哪些操作有副作用？
- 是否需要人工确认？
- 任务是否可能超过一分钟？
- 是否需要保存中间状态？
- 是否需要用户看到执行过程？

不要一上来就做完全自主 Agent。先判断能否用 Workflow Agent。

### 场景二：Agent 一直循环不结束

排查：

- 是否设置最大步骤数。
- 是否设置最大工具调用次数。
- 是否工具结果无法满足模型。
- Prompt 是否要求模型必须结束。
- 是否缺少失败终止条件。

处理：

- 加步骤上限。
- 加成本预算。
- 加明确停止条件。
- 失败后总结已完成内容。

### 场景三：Agent 调错工具

可能原因：

- 工具描述不清楚。
- 工具太多。
- Prompt 没限制工具范围。
- 用户意图分类错误。

处理：

- 按场景裁剪工具集。
- 优化工具 name 和 description。
- 增加工具选择前的意图分类。
- 对高危工具做确认。

### 场景四：Agent 任务中途失败

处理：

- 保存失败步骤。
- 保存错误码。
- 展示给用户。
- 支持重试该步骤。
- 或从上一个稳定状态恢复。

不要让用户只看到“任务失败”，却不知道失败在哪一步。

### 场景五：Agent 要执行命令或改文件

必须考虑：

- 沙盒。
- 文件范围限制。
- 命令白名单。
- 超时。
- 日志。
- 人工确认。
- 结果回滚。

这会在第 10 章继续展开。

### 场景六：面试官问“Agent 后端怎么设计”

可以这样回答：

> 我会把 Agent 设计成可控的任务编排系统。入口创建 task，Planner 生成步骤，Executor 按步骤调用工具，Tool Registry 控制工具和权限，State Store 保存任务状态，SSE 推送执行进度。对高危操作进入 human-in-the-loop，所有步骤记录 trace、tool、token、耗时和错误。系统设置最大步骤数、最大时间和成本预算，避免无限循环。

## 本章完成标记

- [ ] 我能用自己的话解释本章核心概念。
- [ ] 我能回答本章面试题。
- [ ] 我能完成本章编程题或设计题。
- [ ] 我能说出这个知识点在工作中的典型场景。

