# --- Stage 1: Builder Stage ---
FROM golang:1.26-bookworm AS builder
WORKDIR /app
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
# Clear vendor if it exists in the build context
RUN rm -rf vendor
RUN go mod download

COPY . .
RUN rm -rf vendor

# Install the driver and browsers
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
RUN go run github.com/mxschmitt/playwright-go/cmd/playwright install --with-deps chromium

# Build the crawler binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# --- Stage 2: Production Runtime ---
FROM mcr.microsoft.com/playwright:v1.49.1-jammy
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/bin/crawler .

# Copy driver and browsers
COPY --from=builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go
COPY --from=builder /ms-playwright /ms-playwright

ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
