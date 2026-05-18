package chapter04

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "04",
		Title:  "SSE 流式输出与 Channel 通信",
		Source: "大模型应用后端教材/第04章-SSE流式输出与Channel通信.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
