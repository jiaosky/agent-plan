package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"agent-plan/internal/chapter02"
)

func TestChapterRoutesIncludeChapter02SSEExample(t *testing.T) {
	routes := allChapterRoutes()
	if len(routes) != 1 {
		t.Fatalf("expected 1 chapter route, got %d", len(routes))
	}

	route := routes[0]
	if route.Method != http.MethodGet {
		t.Fatalf("expected method %q, got %q", http.MethodGet, route.Method)
	}
	if route.Path != "/chapter02/time-sse" {
		t.Fatalf("expected path %q, got %q", "/chapter02/time-sse", route.Path)
	}
	if route.Example != "/chapter02/time-sse?client_id=demo" {
		t.Fatalf("expected example %q, got %q", "/chapter02/time-sse?client_id=demo", route.Example)
	}
}

func TestRootHandlerPrintsChapterRouteExamples(t *testing.T) {
	mux := newHTTPMux(nil, chapter02.Routes())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "GET /chapter02/time-sse?client_id=demo") {
		t.Fatalf("expected root output to include chapter02 SSE route, got:\n%s", body)
	}
}
