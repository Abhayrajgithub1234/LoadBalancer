FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o load-balancer ./cmd/server/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/load-balancer .
ENTRYPOINT ["./load-balancer"]
CMD ["-backends", "http://host.docker.internal:8001,http://host.docker.internal:8002"]