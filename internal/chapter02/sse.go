package chapter02

import (
	"context"
	"net/http"
	"strings"
	"time"

	"agent-plan/internal/shared"
)

const (
	defaultTimeSSETimeout  = 10 * time.Second
	defaultTimeSSEInterval = time.Second
)

type timeSSEEvent struct {
	ClientID string `json:"client_id"`
	Time     string `json:"time"`
}

// TimeSSEHandler 返回一个 SSE 接口，每隔 1 秒向客户端发送一次当前时间。
//
// 请求必须携带 client_id 查询参数，例如：
// /chapter02/time-sse?client_id=demo
func TimeSSEHandler() http.HandlerFunc {
	return newTimeSSEHandler(defaultTimeSSETimeout, defaultTimeSSEInterval, time.Now)
}

// newTimeSSEHandler 创建可配置超时时间和发送间隔的 SSE Handler，方便测试时缩短等待时间。
func newTimeSSEHandler(timeout, interval time.Duration, now func() time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := strings.TrimSpace(r.URL.Query().Get("client_id"))
		if clientID == "" {
			http.Error(w, "missing client_id query parameter", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		handler := shared.SSEHandler("time", func(ctx context.Context, events chan<- timeSSEEvent) {
			streamTimeEvents(ctx, events, clientID, interval, now)
		})
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

// streamTimeEvents 在独立协程中定时生成当前时间事件，直到 context 超时或客户端断开连接。
func streamTimeEvents(ctx context.Context, events chan<- timeSSEEvent, clientID string, interval time.Duration, now func() time.Time) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			event := timeSSEEvent{
				ClientID: clientID,
				Time:     now().Format(time.RFC3339),
			}
			select {
			case <-ctx.Done():
				return
			case events <- event:
			}
		}
	}
}
