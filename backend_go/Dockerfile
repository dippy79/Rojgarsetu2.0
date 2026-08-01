FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
COPY go.mod go.sum* ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8083
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8083/health || exit 1
CMD ["./server"]
