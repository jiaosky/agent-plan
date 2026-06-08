package chapter02

import (
	"net/http"

	"agent-plan/internal/shared"
)

// RegisterRoutes 注册第 02 章的 HTTP 接口。
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/chapter02/time-sse", TimeSSEHandler())
}

// Routes 返回第 02 章对外提供的 HTTP 接口信息。
func Routes() []shared.RouteInfo {
	return []shared.RouteInfo{
		{
			Method:      http.MethodGet,
			Path:        "/chapter02/time-sse",
			Example:     "/chapter02/time-sse?client_id=demo",
			Description: "第 02 章 SSE 时间流示例",
		},
	}
}
