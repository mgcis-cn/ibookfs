################################################################################
# iBookFS API Server Dockerfile
################################################################################

# =============================================================================
# 阶段1: 构建
# =============================================================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o /app/bin/apiserver \
    ./cmd/apiserver

# =============================================================================
# 阶段2: 运行
# =============================================================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai

RUN addgroup -S ibookfs && adduser -S ibookfs -G ibookfs

WORKDIR /app

COPY --from=builder /app/bin/apiserver /app/
COPY configs/ /app/configs/

RUN mkdir -p /app/data /app/logs && \
    chown -R ibookfs:ibookfs /app

USER ibookfs

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/apiserver"]
