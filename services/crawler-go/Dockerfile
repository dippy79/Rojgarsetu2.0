# --- Stage 1: Builder Stage ---
FROM golang:1.26-bookworm AS builder
WORKDIR /app
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install the playwright driver and browsers
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
RUN go run github.com/mxschmitt/playwright-go/cmd/playwright install --with-deps chromium

# Build the crawler binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# --- Stage 2: Production Runtime ---
FROM mcr.microsoft.com/playwright:v1.49.1-jammy
WORKDIR /app

# Security: Create non-root user
RUN groupadd -r appgroup && useradd -r -g appgroup -G audio,video appuser && \
    mkdir -p /home/appuser/Downloads && \
    chown -R appuser:appgroup /home/appuser && \
    chown -R appuser:appgroup /app

ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

# Copy browsers and driver from builder
COPY --from=builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go
COPY --from=builder /ms-playwright /ms-playwright
# Copy the compiled binary
COPY --from=builder /app/bin/crawler .

# Ensure permissions
RUN chown -R appuser:appgroup /app && \
    chown -R appuser:appgroup /ms-playwright && \
    chown -R appuser:appgroup /root/.cache/ms-playwright-go

EXPOSE 8080
USER appuser

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
