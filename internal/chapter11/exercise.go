package chapter11

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "11",
		Title:  "Go 在 AI 后端中的典型场景",
		Source: "大模型应用后端教材/第11章-Go在AI后端中的典型场景.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
