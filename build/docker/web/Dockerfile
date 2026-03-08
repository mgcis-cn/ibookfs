################################################################################
# iBookFS Web Frontend Dockerfile
################################################################################

# =============================================================================
# 阶段1: 构建
# =============================================================================
FROM node:20-alpine AS builder

WORKDIR /app

COPY web/package*.json ./

RUN npm ci --silent

COPY web/ ./

ARG API_BASE_URL=/api/v1
ENV VITE_API_BASE_URL=${API_BASE_URL}
RUN npm run build

# =============================================================================
# 阶段2: Nginx 运行
# =============================================================================
FROM nginx:alpine

COPY build/docker/web/nginx.conf /etc/nginx/nginx.conf

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

CMD ["nginx", "-g", "daemon off;"]
