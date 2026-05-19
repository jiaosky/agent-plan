package chapter01

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func answer(req *ChatRequest) {

	fmt.Println("Model:", req.Model)
	fmt.Println("Temperature:", req.Temperature)
	fmt.Println("MaxTokens:", req.MaxTokens)

	for _, message := range req.Messages {
		fmt.Println("Message:", message.Content)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := question(ctx)

	go func() {

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case chunk, ok := <-c:
				if !ok {
					return
				}
				fmt.Println("Chunk:", chunk.Delta)
			case <-ticker.C:
				cancel()
			}
		}
	}()

	time.Sleep(2 * time.Minute)
}

func question(ctx context.Context) <-chan ChatChunk {
	chunks := make(chan ChatChunk)

	go func() {
		defer close(chunks)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				chunks <- ChatChunk{Delta: fmt.Sprintf("Question %d", i)}
				i++
			}
		}
	}()
	return chunks
}

func TestDemo(t *testing.T) {
	answer(&ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
		Temperature: 0.5,
		MaxTokens:   100,
	})
}
