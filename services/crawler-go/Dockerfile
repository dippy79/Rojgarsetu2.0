# ==========================================
# Stage 1: Builder Stage
# ==========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for dependency resolution if needed
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/crawler ./cmd/scheduler

# ==========================================
# Stage 2: Production Runtime
# ==========================================
FROM mcr.microsoft.com/playwright:v1.45.0-jammy

WORKDIR /app

# Copy the binary
COPY --from=builder /app/bin/crawler .

# Ensure enough resources for Playwright
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
RUN npx playwright install chromium

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

CMD ["./crawler"]
