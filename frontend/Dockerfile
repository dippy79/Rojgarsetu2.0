FROM node:22-alpine AS builder
WORKDIR /app

# Copy package files
COPY package*.json ./

# Use legacy-peer-deps to avoid resolution crashes
RUN npm install --legacy-peer-deps

# Declare build arguments for Next.js public variables
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_WS_HOST

# Set environment variables for the build process
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_WS_HOST=$NEXT_PUBLIC_WS_HOST

# Copy source code
COPY . .

# Environment setup to prevent memory leaks & strict warning crashes
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV CI=false
ENV NEXT_TELEMETRY_DISABLED=1

# Run build with static export
RUN npm run build

# Nginx Production Stage
FROM nginx:alpine
COPY --from=builder /app/out /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
