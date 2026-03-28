package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
				io.WriteString(w, "routing to: "+chosen.URL+"\n")
				fmt.Println("routing to:", chosen.URL)
				return
			}
		}

		http.Error(w, "no backends available", http.StatusServiceUnavailable)
	}
}

func main() {
	s := []*backend.Server{
		{URL: "http://localhost:8001", Alive: true},
		{URL: "http://localhost:8002", Alive: true},
		{URL: "http://localhost:8003", Alive: true},
		{URL: "http://localhost:8004", Alive: true},
		{URL: "http://localhost:8005", Alive: true},
		{URL: "http://localhost:8006", Alive: true},
		{URL: "http://localhost:8007", Alive: true},
		{URL: "http://localhost:8008", Alive: true},
		{URL: "http://localhost:8009", Alive: true},
		{URL: "http://localhost:8000", Alive: true},
		{URL: "http://localhost:8010", Alive: true},
	}

	http.HandleFunc("/", makeHandles(s))

	ctx := context.Background()
	go healthcheck.StartHealthCheck(s, ctx)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
