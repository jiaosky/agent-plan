package chapter07

import "agent-plan/internal/shared"

func Chapter() shared.Chapter {
	return shared.Chapter{
		Number: "07",
		Title:  "Tool Calling 与 Function Calling",
		Source: "大模型应用后端教材/第07章-ToolCalling与FunctionCalling.md",
	}
}

func Exercises() []shared.Exercise {
	return nil
}
