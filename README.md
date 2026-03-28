# HTTP Load Balancer

A lightweight HTTP load balancer built in Go that distributes incoming
requests across multiple backend servers.

## What it is

When an application is hit with many requests, a single server can get
overwhelmed. This load balancer sits in front of multiple backend servers
and distributes traffic across them using round-robin selection. If a
backend goes down, the load balancer detects it automatically and stops
sending traffic to it.

## How it works

- Incoming requests hit the load balancer on `:8080`
- A round-robin algorithm picks the next available backend
- A background goroutine pings each backend every 10 seconds
- Each health check has a 500ms timeout
- Dead backends are marked and skipped automatically
- If all backends are down, the client receives a `503 Service Unavailable`

## Technical decisions

**Why `sync/atomic` for the counter?**
The round-robin counter is a single integer that only needs to be
incremented. Atomic operations handle this without the overhead of
a full mutex lock.

**Why `sync.RWMutex` for backend health?**
Multiple goroutines read the health status on every request, but only
the health checker writes it. RWMutex allows concurrent reads while
ensuring exclusive writes.

**Why a background goroutine for health checks?**
Checking health on every request would add latency. A background goroutine
checks independently every 10 seconds without affecting request routing.

## Failure handling

- Unreachable backend → marked dead, skipped in routing
- Backend fails mid-request → marked dead, client gets 503
- All backends down → client receives `503 Service Unavailable`
- Backend recovers → health checker marks it alive within 10 seconds

## What I'd add next

- Automatic backend discovery instead of hardcoded URLs
- Weighted round-robin for servers with different capacities
- Circuit breaker pattern to stop hammering failing backends
- Metrics endpoint showing request counts per backend

## Architecture Diagram

```text
Clients
    Client 1 ----\
    Client 2 -----+--> Load Balancer :8080 --> Backend A :8001 (alive)
    Client 3 ----/                              \-> Backend B :8002 (alive)
                                                                                             \-> Backend C :8003 (dead, skipped)

Health Checker (every 10s, timeout 500ms)
    checks /status on A, B, and C
```

## Prerequisites

- Backends must expose a `GET /status` endpoint that returns `200 OK`
- Example backend addresses used in this README: `http://localhost:8001` and `http://localhost:8002`

## How to run

**Without Docker:**
```bash
go build -o load-balancer ./cmd/server/
./load-balancer -backends "http://localhost:8001,http://localhost:8002"
```

**With Docker:**
```bash
docker build -t loader .
docker run -p 8080:8080 loader
```

The Docker image default uses:
- `http://host.docker.internal:8001`
- `http://host.docker.internal:8002`

Make sure those backend services are running on your host machine and expose `/status`.

**Check backend health:**
```bash
curl http://localhost:8080/status
```