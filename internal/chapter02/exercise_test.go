package chapter02

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestChapterMetadata(t *testing.T) {
	if !Chapter().IsValid() {
		t.Fatalf("chapter metadata should be populated")
	}
	if Exercises() == nil {
		t.Log("exercises are ready to be filled")
	}
}

func TestTimeSSEHandlerRequiresClientID(t *testing.T) {
	handler := TimeSSEHandler()
	req := httptest.NewRequest(http.MethodGet, "/chapter02/time-sse", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTimeSSEHandlerStreamsCurrentTime(t *testing.T) {
	fixedTime := time.Date(2026, 6, 8, 12, 30, 0, 0, time.UTC)
	handler := newTimeSSEHandler(25*time.Millisecond, 5*time.Millisecond, func() time.Time {
		return fixedTime
	})
	req := httptest.NewRequest(http.MethodGet, "/chapter02/time-sse?client_id=demo", nil)
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
		"event: time\n",
		`"client_id":"demo"`,
		`"time":"2026-06-08T12:30:00Z"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body %q should contain %q", body, want)
		}
	}
}

func TestTimeSSEHandlerStopsWhenRequestContextIsCanceled(t *testing.T) {
	handler := newTimeSSEHandler(time.Minute, 5*time.Millisecond, time.Now)
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/chapter02/time-sse?client_id=demo", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		handler.ServeHTTP(rec, req)
	}()

	time.Sleep(15 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("handler did not stop after request context was canceled")
	}
}
