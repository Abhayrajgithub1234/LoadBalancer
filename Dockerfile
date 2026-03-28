FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o load-balancer ./cmd/server/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/load-balancer .
CMD ["./load-balancer"]