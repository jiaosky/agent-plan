package chapter10

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "10",
		Title:  "沙盒执行与安全隔离",
		Source: "大模型应用后端教材/第10章-沙盒执行与安全隔离.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
