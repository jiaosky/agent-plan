package chapter05

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "05",
		Title:  "RAG 与向量检索",
		Source: "大模型应用后端教材/第05章-RAG与向量检索.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
