# --- Stage 1: Builder Stage ---
FROM golang:1.23-bookworm AS builder
WORKDIR /app
ENV GOPROXY=https://proxy.golang.org,direct
RUN apt-get update && apt-get install -y git curl && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# --- Stage 2: Production Runtime ---
FROM mcr.microsoft.com/playwright:v1.45.0-jammy
WORKDIR /app

# Security: Create non-root user for Playwright
RUN groupadd -r appgroup && useradd -r -g appgroup -G audio,video appuser && \
    mkdir -p /home/appuser/Downloads && \
    chown -R appuser:appgroup /home/appuser && \
    chown -R appuser:appgroup /app

ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

COPY --from=builder /app/bin/crawler .
# Ensure permissions for browsers
RUN chown -R appuser:appgroup /ms-playwright

EXPOSE 8080
USER appuser

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
