package shared

import "net/http"

// RouteRegistrar 表示一个章节向统一 HTTP 服务注册自己接口的函数。
type RouteRegistrar func(mux *http.ServeMux)

// RouteInfo 描述一个章节对外提供的 HTTP 接口。
type RouteInfo struct {
	Method      string
	Path        string
	Example     string
	Description string
}
