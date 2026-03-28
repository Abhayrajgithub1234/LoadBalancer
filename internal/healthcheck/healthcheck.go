package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Abhayrajgithub123/LoadBalancer/internal/backend"
)

func makeHttpReq(ctx context.Context, bes *backend.Server) {
	req, error := http.NewRequestWithContext(ctx, "GET", bes.URL+"/status", nil)
	if error != nil {
		log.Fatal(error)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		bes.SetAlive(false)
		fmt.Println(bes.URL, "is unreachable")
		return
	}

	if resp.StatusCode == http.StatusOK {
		bes.SetAlive(true)
	} else {
		bes.SetAlive(false)
		fmt.Println(bes.URL, "is unreachable")
		return
	}
	defer resp.Body.Close()
}

func StartHealthCheck(backends []*backend.Server, ctx context.Context) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			for _, b := range backends {
				reqCtx, reqCancel := context.WithTimeout(ctx, 500*time.Millisecond)
				makeHttpReq(reqCtx, b)
				reqCancel()
			}
		case <-ctx.Done():
			return
		}
	}
}
