package chapter06

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "06",
		Title:  "文件管理与知识库工程",
		Source: "大模型应用后端教材/第06章-文件管理与知识库工程.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
