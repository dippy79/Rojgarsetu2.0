FROM golang:alpine AS builder

WORKDIR /app
# Install build essentials
RUN apk add --no-cache git build-base

# Set Go proxy for better reliability
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build the Go binary (Static, No C bindings)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rojgar_api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /app
# Copy ONLY the compiled binary and required assets
COPY --from=builder /app/rojgar_api .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

EXPOSE 8083

# Healthcheck to verify service is ready
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8083/health || exit 1

CMD ["./rojgar_api"]
