package chapter01

import (
	"fmt"
	"log/slog"
)

/*
*
拦截器中进行 auth 处理
*/
func auth() error {
	return nil
}

/*
*
拦截器中进行限流
*/
func ratelimit() error {
	return nil
}

func readSession(sessionId string) {

}

func buildPrompt() (string, error) {
	return "你是一个 xxx", nil
}

func callModel(prompt string) (string, error) {
	return "", nil
}

func Exercise2() {
	readSession("xxx")
	prompt, err := buildPrompt()
	if err != nil {
		fmt.Errorf(err.Error())

	}
	slog.Info("xxx")
	callModel(prompt)
}
