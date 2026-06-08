package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SSEProducer 负责向事件通道写入业务事件，并在 context 结束后退出。
type SSEProducer[T any] func(ctx context.Context, events chan<- T)

// SSEHandler 创建一个通用 SSE Handler，负责响应头、事件写出和 flush。
func SSEHandler[T any](eventType string, produce SSEProducer[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		events := make(chan T)
		done := make(chan struct{})
		go func() {
			defer close(done)
			defer close(events)
			produce(ctx, events)
		}()

		for {
			select {
			case <-ctx.Done():
				<-done
				return
			case event, ok := <-events:
				if !ok {
					<-done
					return
				}
				if err := WriteSSEEvent(w, eventType, event); err != nil {
					cancel()
					<-done
					return
				}
				flusher.Flush()
			}
		}
	}
}

// WriteSSEEvent 按 SSE 格式写出事件类型和 JSON 数据。
func WriteSSEEvent(w http.ResponseWriter, eventType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
	return err
}
