package chapter02

import "log/slog"

// Message 表示发送给聊天模型或从聊天模型返回的一条消息。
//
// Role 使用 string，是因为不同模型供应商支持的角色集合可能略有差异；
// role 字段常见取值及含义：
// "system"：系统指令，用于为对话设定初始背景、规则或风格；
// "user"：用户输入，代表实际提问、任务说明或需求；
// "assistant"：AI 助手的回复，模型根据上下文生成的内容；
// "tool"：工具调用或插件响应，部分高级场景用于描述“函数调用”之类结构化交互。
// 这里保持结构简单，是为了让第一章先聚焦接口契约，而不是提前陷入具体供应商细节。
type Message struct {
	Role    string
	Content string
}

// ChatRequest 描述一次模型调用所需的最小信息。
//
// Model 让调用方可以选择具体模型，同时不需要修改业务代码。
// Messages 承载对话上下文，这是模型真正看到的输入。
// Temperature 控制回答随机性，MaxTokens 限制最大输出长度，便于后端控制延迟和成本。
type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float64
	MaxTokens   int
}

// ChatResponse 表示一次非流式聊天调用的结果。
//
// Content 是模型最终生成的完整回答。
// Token 用量放在响应里，是因为 AI 应用后端通常需要记录成本、审计用量，
// 并排查某些请求为什么消耗异常高。
type ChatResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

type LLMRequest struct {
	company string
	ChatRequest
}

type LLMResponse struct {
	ChatResponse
}

// 请求接口
type Chat interface {
	stream(request LLMRequest) (LLMRequest, error)
	general(requset LLMRequest) (LLMRequest, error)
}

type Qwen struct {
}

func (*Qwen) stream(request LLMRequest) (LLMRequest, error) {

	return LLMRequest{}, nil
}

func (*Qwen) general(request LLMRequest) (LLMRequest, error) {

	return LLMRequest{}, nil
}

type Kimi struct {
}

func (*Kimi) stream(request LLMRequest) (LLMRequest, error) {

	return LLMRequest{}, nil
}

func (*Kimi) general(request LLMRequest) (LLMRequest, error) {

	return LLMRequest{}, nil
}

func switchCompany(company string) Chat {
	switch company {
	case "qwen":
		return &Qwen{}
	case "kimi":
		return &Kimi{}
	default:
		slog.Info("无企业命中")
		return nil
	}
}
