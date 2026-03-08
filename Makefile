.PHONY: build run test lint clean tidy help \
        docker-build docker-run docker-push docker-compose-up docker-compose-down \
        build-backend build-frontend deploy check-env install-env

# =============================================================================
# 变量定义
# =============================================================================
APP_NAME := ibookfs
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y%m%d-%H%M%S)
BIN_DIR := bin
DOCKER_REGISTRY ?= 
DOCKER_IMAGE := $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/,)$(APP_NAME)

# Go 编译参数
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
GOFLAGS := -ldflags="$(LDFLAGS)"

# =============================================================================
# 帮助信息
# =============================================================================
help:
	@echo "iBookFS Makefile"
	@echo ""
	@echo "使用: make [target]"
	@echo ""
	@echo "开发命令:"
	@echo "  build           编译后端 (本地平台)"
	@echo "  build-linux     编译后端 (Linux amd64)"
	@echo "  build-frontend  编译前端"
	@echo "  run             运行后端开发服务器"
	@echo "  test            运行测试"
	@echo "  lint            运行代码检查"
	@echo "  tidy            整理 Go 模块"
	@echo "  clean           清理构建产物"
	@echo ""
	@echo "Docker 命令:"
	@echo "  docker-build           构建完整镜像"
	@echo "  docker-build-apiserver 构建API Server镜像"
	@echo "  docker-build-web       构建Web前端镜像"
	@echo "  docker-run-apiserver   运行API Server容器"
	@echo "  docker-run-web         运行Web前端容器"
	@echo "  docker-stop-apiserver  停止API Server容器"
	@echo "  docker-stop-web        停止Web前端容器"
	@echo "  docker-compose-up      启动所有服务"
	@echo "  docker-compose-down    停止所有服务"
	@echo ""
	@echo "部署命令:"
	@echo "  check-env              检查环境依赖"
	@echo "  install-env            安装环境依赖"
	@echo "  deploy                 完整构建部署"
	@echo "  deploy-apiserver       API Server部署（systemd）"
	@echo "  deploy-web             Web前端部署（本地）"
	@echo "  deploy-apiserver-docker API Server Docker部署"
	@echo "  deploy-web-docker      Web前端Docker部署"
	@echo ""
	@echo "变量:"
	@echo "  DOCKER_REGISTRY 镜像仓库地址 (默认: 空)"
	@echo "  VERSION         版本号 (默认: git tag)"

# =============================================================================
# 开发命令
# =============================================================================
build:
	@mkdir -p $(BIN_DIR)
	go build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) ./cmd/apiserver

build-linux:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 ./cmd/apiserver

build-frontend:
	cd web && npm ci && npm run build

run:
	go run ./cmd/apiserver

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BIN_DIR)
	rm -rf web/dist
	rm -rf web/node_modules

tidy:
	go mod tidy

# =============================================================================
# Docker 命令
# =============================================================================
docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

docker-run:
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE):latest

docker-push:
	@if [ -z "$(DOCKER_REGISTRY)" ]; then echo "请设置 DOCKER_REGISTRY"; exit 1; fi
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down

docker-compose-logs:
	docker compose logs -f

# =============================================================================
# 独立构建命令（按服务目录组织）
# =============================================================================
docker-build-apiserver:
	docker build -f build/docker/apiserver/Dockerfile -t $(DOCKER_IMAGE)-apiserver:$(VERSION) -t $(DOCKER_IMAGE)-apiserver:latest .

docker-build-web:
	docker build -f build/docker/web/Dockerfile -t $(DOCKER_IMAGE)-web:$(VERSION) -t $(DOCKER_IMAGE)-web:latest .

docker-run-apiserver:
	docker compose -f build/docker/apiserver/docker-compose.yml up -d

docker-run-web:
	docker compose -f build/docker/web/docker-compose.yml up -d

docker-stop-apiserver:
	docker compose -f build/docker/apiserver/docker-compose.yml down

docker-stop-web:
	docker compose -f build/docker/web/docker-compose.yml down

# =============================================================================
# 部署命令
# =============================================================================
check-env:
	@bash scripts/check-env.sh

install-env:
	@sudo bash scripts/install-env.sh

deploy:
	@bash scripts/build.sh

deploy-apiserver:
	@bash scripts/deploy-backend.sh deploy

deploy-web:
	@bash scripts/deploy-frontend.sh deploy

deploy-apiserver-docker:
	@bash scripts/deploy-backend.sh docker-run

deploy-web-docker:
	@bash scripts/deploy-frontend.sh docker-run
