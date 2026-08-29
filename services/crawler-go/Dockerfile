# ==========================================
# Stage 1: Builder Stage
# ==========================================
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y git curl && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the crawler binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# Install playwright dependencies for browser installation
RUN apt-get update && apt-get install -y \
    libglib2.0-0 libnss3 libnspr4 libdbus-1-3 \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 \
    libdrm2 libxkbcommon0 libxcomposite1 libxdamage1 \
    libxfixes3 libxrandr2 libgbm1 libasound2 \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

# Pre-download playwright browsers at build time
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
# Use exact package from go.mod to avoid mismatches
RUN go run -mod=mod github.com/mxschmitt/playwright-go/cmd/playwright install --with-deps chromium

# ==========================================
# Stage 2: Production Runtime
# ==========================================
FROM mcr.microsoft.com/playwright:v1.45.0-jammy

WORKDIR /app

# Set browsers path to reuse from builder stage
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
# Enforce NO runtime downloads
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

# Copy pre-downloaded browsers (Chromium only)
COPY --from=builder /ms-playwright /ms-playwright
# Copy playwright driver cache to avoid runtime download of the Go driver
COPY --from=builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go

# Copy the compiled binary
COPY --from=builder /app/bin/crawler .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
