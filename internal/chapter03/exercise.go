package chapter03

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "03",
		Title:  "Prompt 与上下文工程",
		Source: "大模型应用后端教材/第03章-Prompt与上下文工程.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
