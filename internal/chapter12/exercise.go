package chapter12

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "12",
		Title:  "算法题与面试验收",
		Source: "大模型应用后端教材/第12章-算法题与面试验收.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
