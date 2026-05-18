package chapter09

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "09",
		Title:  "Agent 范式与工作流编排",
		Source: "大模型应用后端教材/第09章-Agent范式与工作流编排.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
