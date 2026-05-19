package chapter01

import (
	"context"
	"fmt"
	"io"
	"strings"

	"agent-plan/internal/shared"
)

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

// ChatChunk 表示流式模型响应中的一个增量片段。
//
// Delta 是本次新生成的文本片段。
// FinishReason 是可选的结束原因，用来告诉调用方流为什么结束，例如 "stop"、
// "length"，或者某个供应商自定义的原因。
type ChatChunk struct {
	Delta        string
	FinishReason string
}

// LLMClient 抽象了对大语言模型供应商的访问。
//
// 业务层应该依赖这个接口，而不是直接依赖 OpenAI、DeepSeek、Claude、通义，
// 或某个私有模型 SDK。这样后续切换供应商、用 fake client 写测试、
// 或者建设统一模型网关时，都不需要重写上层业务逻辑。
//
// 每个方法都接收 context.Context，是因为模型调用通常是远程、慢、且有成本的操作。
// Context 可以让调用方设置超时，在用户断开连接时取消请求，并在后端调用链路中传递
// 请求级别的 trace 信息。
type LLMClient interface {
	// Chat 等待模型供应商生成完成后，返回一个完整答案。
	// 它适合短任务，或者不需要实时展示中间输出的内部任务。
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// Stream 返回一个只读 channel，让调用方可以在模型生成内容时逐块接收，
	// 并通常通过 SSE 转发给前端。
	// 使用只读 channel 可以避免消费方误写入模型流，也能把流生产的所有权留在客户端实现内部。
	Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
}

// MockLLMClient 是一个用于学习和本地调试的假模型客户端。
//
// 它的价值不是生成真实智能回答，而是帮助业务代码先依赖 LLMClient 接口跑起来。
// 这样在还没有接入真实供应商之前，也可以练习请求构造、普通返回、流式返回、
// context 取消和后续 Web Handler 调用流程。
type MockLLMClient struct {
	Reply string
}

// Chat 返回一个固定回答，用来模拟非流式模型调用。
//
// 这里先检查 ctx.Done()，是为了体现真实模型调用必须支持取消：
// 如果用户已经断开连接，或者请求已经超时，后端就不应该继续生成和计费。
func (c MockLLMClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	content := c.reply()

	return &ChatResponse{
		Content:      content,
		InputTokens:  estimateTokensFromMessages(req.Messages),
		OutputTokens: estimateTokens(content),
	}, nil
}

// Stream 把固定回答拆成多个片段，通过只读 channel 返回。
//
// 真实供应商的流式接口通常是一边生成一边返回 chunk。Mock 这里用字符串切片模拟
// 这个过程，方便后续 Web 服务把每个 ChatChunk 转成 SSE 事件。
func (c MockLLMClient) Stream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
	chunks := make(chan ChatChunk)
	parts := strings.Fields(c.reply())

	go func() {
		defer close(chunks)

		for _, part := range parts {
			select {
			case <-ctx.Done():
				return
			case chunks <- ChatChunk{Delta: part + " "}:
			}
		}

		select {
		case <-ctx.Done():
		case chunks <- ChatChunk{FinishReason: "stop"}:
		}
	}()

	return chunks, nil
}

func (c MockLLMClient) reply() string {
	if strings.TrimSpace(c.Reply) == "" {
		return "这是一个 mock 模型回答，用来演示 LLMClient 接口。"
	}
	return c.Reply
}

func estimateTokensFromMessages(messages []Message) int {
	total := 0
	for _, message := range messages {
		total += estimateTokens(message.Role)
		total += estimateTokens(message.Content)
	}
	return total
}

func estimateTokens(text string) int {
	return len([]rune(strings.TrimSpace(text)))
}

// RunMockDemo 运行第一章模型客户端接口的本地演示。
//
// demo 放在 chapter01 包里，是因为这里最了解 Message、ChatRequest 和 MockLLMClient
// 的教学含义；cmd/chapter 只负责命令行调度，不应该堆放每一章的业务演示细节。
func RunMockDemo(ctx context.Context, out io.Writer) error {
	client := MockLLMClient{}
	req := ChatRequest{
		Model: "mock-model",
		Messages: []Message{
			{Role: "system", Content: "你是一个帮助学习 AI 后端的助手。"},
			{Role: "user", Content: "请演示一下模型客户端接口。"},
		},
		Temperature: 0.7,
		MaxTokens:   128,
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Mock Chat:")

	resp, err := client.Chat(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "content: %s\n", resp.Content)
	fmt.Fprintf(out, "input_tokens: %d, output_tokens: %d\n", resp.InputTokens, resp.OutputTokens)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Mock Stream:")

	stream, err := client.Stream(ctx, req)
	if err != nil {
		return err
	}
	for chunk := range stream {
		if chunk.Delta != "" {
			fmt.Fprint(out, chunk.Delta)
		}
		if chunk.FinishReason != "" {
			fmt.Fprintf(out, "\nfinish_reason: %s\n", chunk.FinishReason)
		}
	}

	return nil
}

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "01",
		Title:  "大模型应用后端全景",
		Source: "大模型应用后端教材/第01章-大模型应用后端全景.md",
	}
}

func Exercises() []shared.Exercise {
	return []shared.Exercise{
		{
			Name:        "4.3.1 设计模型客户端接口",
			Description: "定义 LLMClient，支持普通聊天、流式聊天、context 控制和基础模型参数。",
			Done:        true,
		},
	}
}c
