################################################################################
# iBookFS Multi-stage Dockerfile
# 
# 构建阶段：编译Go后端 + 构建前端
# 运行阶段：最小化镜像运行服务
################################################################################

# =============================================================================
# 阶段1: 前端构建
# =============================================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装依赖
RUN npm ci --silent

# 复制前端源码
COPY web/ ./

# 构建前端
RUN npm run build

# =============================================================================
# 阶段2: 后端构建
# =============================================================================
FROM golang:1.25-alpine AS backend-builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 设置Go代理（加速国内下载）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# 复制Go依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 编译后端
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=$(date +%Y%m%d)" \
    -o /app/bin/ibookfs \
    ./cmd/apiserver

# =============================================================================
# 阶段3: 运行镜像
# =============================================================================
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -S ibookfs && adduser -S ibookfs -G ibookfs

WORKDIR /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/bin/ibookfs /app/
COPY --from=frontend-builder /app/web/dist /app/web/

# 复制配置文件模板
COPY configs/ /app/configs/

# 创建必要目录
RUN mkdir -p /app/data /app/logs && \
    chown -R ibookfs:ibookfs /app

# 切换到非root用户
USER ibookfs

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["/app/ibookfs"]
CMD ["--config", "/app/configs/config.yaml"]
