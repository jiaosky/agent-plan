package chapter02

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "02",
		Title:  "LLM 基础与模型调用",
		Source: "大模型应用后端教材/第02章-LLM基础与模型调用.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
