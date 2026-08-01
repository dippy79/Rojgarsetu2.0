# ==========================================
# Stage 1: Builder Stage (Offline Build)
# ==========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy full source including local 'vendor' directory
COPY . .

# Build statically using local vendor directory (NO NETWORK CALLS!)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# ==========================================
# Stage 2: Production Runtime
# ==========================================
FROM alpine:3.20

# Install Chromium for headless browser crawling + other system deps
RUN apk add --no-cache ca-certificates tzdata wget chromium && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

# Set Chrome binary path for the detection logic in browser.go
ENV CHROME_BIN=/usr/bin/chromium-browser

WORKDIR /app

COPY --from=builder /app/bin/crawler .

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
