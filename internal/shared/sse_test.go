package shared

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type testSSEEvent struct {
	Message string `json:"message"`
}

func TestSSEHandlerStreamsJSONEvents(t *testing.T) {
	handler := SSEHandler("message", func(ctx context.Context, events chan<- testSSEEvent) {
		select {
		case <-ctx.Done():
		case events <- testSSEEvent{Message: "hello"}:
		}
	})
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "text/event-stream; charset=utf-8" {
		t.Fatalf("content type = %q, want SSE content type", got)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"event: message\n",
		`data: {"message":"hello"}`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body %q should contain %q", body, want)
		}
	}
}

func TestSSEHandlerStopsWhenRequestContextIsCanceled(t *testing.T) {
	started := make(chan struct{})
	stopped := make(chan struct{})
	handler := SSEHandler("message", func(ctx context.Context, events chan<- testSSEEvent) {
		close(started)
		<-ctx.Done()
		close(stopped)
	})
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		handler.ServeHTTP(rec, req)
	}()

	<-started
	cancel()

	select {
	case <-stopped:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("producer did not receive request cancellation")
	}

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("handler did not stop after request context was canceled")
	}
}
