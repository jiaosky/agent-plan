package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"agent-plan/internal/chapter01"
	"agent-plan/internal/chapter02"
	"agent-plan/internal/chapter03"
	"agent-plan/internal/chapter04"
	"agent-plan/internal/chapter05"
	"agent-plan/internal/chapter06"
	"agent-plan/internal/chapter07"
	"agent-plan/internal/chapter08"
	"agent-plan/internal/chapter09"
	"agent-plan/internal/chapter10"
	"agent-plan/internal/chapter11"
	"agent-plan/internal/chapter12"
	"agent-plan/internal/chapter13"
	"agent-plan/internal/chapter14"
	"agent-plan/internal/shared"
)

func main() {
	chapterFlag := flag.String("chapter", "", "chapter number, for example 01 or 13")
	serveFlag := flag.Bool("serve", false, "start chapter HTTP server")
	addrFlag := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	chapters := allChapters()

	if *serveFlag {
		routes := allChapterRoutes()
		mux := newHTTPMux(chapters, routes)
		baseURL := localBaseURL(*addrFlag)
		fmt.Printf("chapter HTTP server listening on %s\n", *addrFlag)
		fmt.Println("health:   " + baseURL + "/healthz")
		fmt.Println("chapters: " + baseURL + "/chapters")
		printChapterRoutes(baseURL, routes)
		server := &http.Server{
			Addr:              *addrFlag,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "start HTTP server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *chapterFlag == "" {
		for _, chapter := range chapters {
			fmt.Printf("%s  %s\n", chapter.Number, chapter.Title)
		}
		return
	}

	want := strings.TrimSpace(*chapterFlag)
	for _, chapter := range chapters {
		if chapter.Number == want {
			fmt.Printf("Chapter %s: %s\nSource: %s\nPackage: internal/chapter%s\n", chapter.Number, chapter.Title, chapter.Source, chapter.Number)
			if chapter.Number == "01" {
				if err := chapter01.RunMockDemo(context.Background(), os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "run chapter 01 demo: %v\n", err)
					os.Exit(1)
				}
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "chapter %q not found\n", want)
	os.Exit(1)
}

func allChapters() []shared.Chapter {
	return []shared.Chapter{
		chapter01.Chapter(),
		chapter02.Chapter(),
		chapter03.Chapter(),
		chapter04.Chapter(),
		chapter05.Chapter(),
		chapter06.Chapter(),
		chapter07.Chapter(),
		chapter08.Chapter(),
		chapter09.Chapter(),
		chapter10.Chapter(),
		chapter11.Chapter(),
		chapter12.Chapter(),
		chapter13.Chapter(),
		chapter14.Chapter(),
	}
}

func allChapterRoutes() []shared.RouteInfo {
	var routes []shared.RouteInfo
	for _, chapterRoutes := range [][]shared.RouteInfo{
		chapter02.Routes(),
	} {
		routes = append(routes, chapterRoutes...)
	}
	return routes
}

func printChapterRoutes(baseURL string, routes []shared.RouteInfo) {
	for _, route := range routes {
		path := route.Path
		if route.Example != "" {
			path = route.Example
		}
		fmt.Printf("%s: %s%s\n", route.Description, baseURL, path)
	}
}

func newHTTPMux(chapters []shared.Chapter, routes []shared.RouteInfo) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "agent-plan chapter HTTP server")
		fmt.Fprintln(w, "GET /healthz")
		fmt.Fprintln(w, "GET /chapters")
		for _, route := range routes {
			path := route.Path
			if route.Example != "" {
				path = route.Example
			}
			fmt.Fprintf(w, "%s %s\n", route.Method, path)
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/chapters", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(chapters); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	for _, register := range []shared.RouteRegistrar{
		chapter01.RegisterRoutes,
		chapter02.RegisterRoutes,
		chapter03.RegisterRoutes,
		chapter04.RegisterRoutes,
		chapter05.RegisterRoutes,
		chapter06.RegisterRoutes,
		chapter07.RegisterRoutes,
		chapter08.RegisterRoutes,
		chapter09.RegisterRoutes,
		chapter10.RegisterRoutes,
		chapter11.RegisterRoutes,
		chapter12.RegisterRoutes,
		chapter13.RegisterRoutes,
		chapter14.RegisterRoutes,
	} {
		register(mux)
	}

	return mux
}

func localBaseURL(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}
	return "http://" + addr
}
