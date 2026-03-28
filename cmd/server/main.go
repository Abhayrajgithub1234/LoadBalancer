package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync/atomic"

	"github.com/Abhayrajgithub123/LoadBalancer/internal/backend"
	"github.com/Abhayrajgithub123/LoadBalancer/internal/healthcheck"
)

func makeHandles(backends []*backend.Server) http.HandlerFunc {

	var counter int64

	return func(w http.ResponseWriter, req *http.Request) {
		current := atomic.AddInt64(&counter, 1)

		for i := 0; i < len(backends); i++ {
			index := (current - 1 + int64(i)) % int64(len(backends))
			chosen := backends[index]
			if chosen.IsAlive() {
				proxy := httputil.NewSingleHostReverseProxy(chosen.ParsedUrl())
				proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
					slog.Warn("backend failed", "url", chosen.URL)
					chosen.SetAlive(false)
				}
				proxy.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, "no backends available", http.StatusServiceUnavailable)
	}
}

func status(s []*backend.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result := "["
		for i, b := range s {
			if i > 0 {
				result += ","
			}
			status := "dead"
			if b.IsAlive() {
				status = "alive"
			}
			result += fmt.Sprintf(`{"url":"%s","status":"%s"}`, b.URL, status)
		}
		result += "]"
		io.WriteString(w, result)
	}
}

func main() {

	addr := flag.String("addr", ":8080", "load balancer address")
	backendsFlag := flag.String("backends", "http://localhost:8001,http://localhost:8002", "comma separated backend URLs")
	flag.Parse()

	urls := strings.Split(*backendsFlag, ",")
	s := make([]*backend.Server, len(urls))
	for i, u := range urls {
		s[i] = &backend.Server{URL: strings.TrimSpace(u), Alive: true}
	}

	http.HandleFunc("/", makeHandles(s))
	http.HandleFunc("/status", status(s))

	ctx := context.Background()
	go healthcheck.StartHealthCheck(s, ctx)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		slog.Error("server failed", "error", err)
	}

}
