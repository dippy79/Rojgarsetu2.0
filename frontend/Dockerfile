# --- Stage 1: Build Stage ---
FROM node:22-alpine3.20 AS builder
WORKDIR /app

# Add production dependencies only first for caching
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Build args for Next.js public envs
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_WS_HOST
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_WS_HOST=$NEXT_PUBLIC_WS_HOST

# Copy source and build
COPY . .
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# --- Stage 2: Production Stage ---
FROM nginx:1.27-alpine3.20
WORKDIR /usr/share/nginx/html

# Security: Create a non-root appuser
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    chown -R appuser:appgroup /var/cache/nginx && \
    touch /var/run/nginx.pid && \
    chown -R appuser:appgroup /var/run/nginx.pid

# Copy built static files
COPY --from=builder /app/out .

EXPOSE 80
USER appuser
CMD ["nginx", "-g", "daemon off;"]
