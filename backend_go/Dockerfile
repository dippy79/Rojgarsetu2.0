# --- Stage 1: Build Stage ---
FROM golang:1.23-alpine3.20 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
RUN apk add --no-cache git build-base
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -mod=mod -a -installsuffix cgo -o rojgar_api ./cmd/server

# --- Stage 2: Production Stage ---
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Security: Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/rojgar_api .
COPY --from=builder --chown=appuser:appgroup /app/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/docs ./docs

EXPOSE 8083
USER appuser

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8083/health || exit 1

CMD ["./rojgar_api"]
