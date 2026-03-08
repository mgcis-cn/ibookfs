#!/bin/bash

################################################################################
# iBookFS 后端独立部署脚本
#
# 功能：仅构建和部署后端服务
# 支持 Linux (systemd) 和 macOS (launchd)
################################################################################

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检测操作系统
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS_TYPE="macos"
    DEPLOY_DIR="${DEPLOY_DIR:-$HOME/.ibookfs}"
else
    OS_TYPE="linux"
    DEPLOY_DIR="${DEPLOY_DIR:-/opt/ibookfs}"
fi

# 配置变量
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SERVICE_NAME="ibookfs"
PLIST_NAME="cn.mgcis.ibookfs"

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 Go 环境
check_go() {
    log_info "检查 Go 环境..."
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装"
        exit 1
    fi
    log_success "Go $(go version | awk '{print $3}')"
}

# 构建后端
build_backend() {
    log_info "构建后端..."
    
    cd "$PROJECT_DIR"
    mkdir -p "$PROJECT_DIR/bin"
    
    export GOPROXY=https://goproxy.cn,direct
    export GOSUMDB=sum.golang.google.cn
    
    go mod tidy
    go mod download
    
    VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
    BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    
    # 根据操作系统设置编译参数
    if [ "$OS_TYPE" = "macos" ]; then
        go build \
            -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
            -o "$PROJECT_DIR/bin/ibookfs" \
            ./cmd/apiserver
    else
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
            -o "$PROJECT_DIR/bin/ibookfs" \
            ./cmd/apiserver
    fi
    
    if [ ! -f "$PROJECT_DIR/bin/ibookfs" ]; then
        log_error "构建失败"
        exit 1
    fi
    
    log_success "后端构建完成: $PROJECT_DIR/bin/ibookfs"
}

# 部署后端
deploy_backend() {
    log_info "部署后端到 $DEPLOY_DIR..."
    
    mkdir -p "$DEPLOY_DIR"
    mkdir -p "$DEPLOY_DIR/configs"
    mkdir -p "$DEPLOY_DIR/logs"
    mkdir -p "$DEPLOY_DIR/data"
    
    # 停止服务
    stop_service
    
    # 复制文件（源码目录和部署目录相同时跳过）
    if [ "$PROJECT_DIR" != "$DEPLOY_DIR" ]; then
        cp "$PROJECT_DIR/bin/ibookfs" "$DEPLOY_DIR/"
        cp "$PROJECT_DIR/configs/application.yaml" "$DEPLOY_DIR/configs/"
        cp "$PROJECT_DIR/configs/application-test.yaml" "$DEPLOY_DIR/configs/" 2>/dev/null || true
        cp "$PROJECT_DIR/configs/application-prod.yaml" "$DEPLOY_DIR/configs/" 2>/dev/null || true
    else
        cp "$PROJECT_DIR/bin/ibookfs" "$DEPLOY_DIR/ibookfs"
    fi
    chmod +x "$DEPLOY_DIR/ibookfs"
    
    log_success "后端部署完成"
}

# 停止服务
stop_service() {
    if [ "$OS_TYPE" = "macos" ]; then
        if launchctl list | grep -q "$PLIST_NAME" 2>/dev/null; then
            log_info "停止现有服务..."
            launchctl unload ~/Library/LaunchAgents/${PLIST_NAME}.plist 2>/dev/null || true
        fi
    else
        if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
            log_info "停止现有服务..."
            sudo systemctl stop "$SERVICE_NAME"
        fi
    fi
}

# 配置服务（根据操作系统选择 systemd 或 launchd）
setup_service() {
    if [ "$OS_TYPE" = "macos" ]; then
        setup_launchd
    else
        setup_systemd
    fi
}

# 配置 macOS launchd
setup_launchd() {
    log_info "配置 launchd 服务..."
    
    mkdir -p ~/Library/LaunchAgents
    
    cat <<EOF > ~/Library/LaunchAgents/${PLIST_NAME}.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>${PLIST_NAME}</string>
    <key>ProgramArguments</key>
    <array>
        <string>${DEPLOY_DIR}/ibookfs</string>
        <string>--conf</string>
        <string>${DEPLOY_DIR}/configs</string>
        <string>--env</string>
        <string>test</string>
    </array>
    <key>WorkingDirectory</key>
    <string>${DEPLOY_DIR}</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${DEPLOY_DIR}/logs/ibookfs.log</string>
    <key>StandardErrorPath</key>
    <string>${DEPLOY_DIR}/logs/ibookfs.error.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>GIN_MODE</key>
        <string>release</string>
    </dict>
</dict>
</plist>
EOF
    
    log_success "launchd 服务配置完成"
}

# 配置 Linux systemd
setup_systemd() {
    log_info "配置 systemd 服务..."
    
    cat <<EOF | sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null
[Unit]
Description=iBookFS Backend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DEPLOY_DIR
ExecStart=$DEPLOY_DIR/ibookfs --conf $DEPLOY_DIR/configs --env prod
Restart=on-failure
RestartSec=5s
Environment=GIN_MODE=release
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE_NAME"
    
    log_success "systemd 服务配置完成"
}

# 启动服务
start_service() {
    log_info "启动服务..."
    
    if [ "$OS_TYPE" = "macos" ]; then
        launchctl load ~/Library/LaunchAgents/${PLIST_NAME}.plist
        sleep 2
        
        if launchctl list | grep -q "$PLIST_NAME"; then
            log_success "服务启动成功"
            log_info "查看日志: tail -f $DEPLOY_DIR/logs/ibookfs.log"
        else
            log_error "服务启动失败"
            cat "$DEPLOY_DIR/logs/ibookfs.error.log" 2>/dev/null || true
            exit 1
        fi
    else
        sudo systemctl start "$SERVICE_NAME"
        sleep 2
        
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            log_success "服务启动成功"
            systemctl status "$SERVICE_NAME" --no-pager
        else
            log_error "服务启动失败"
            sudo journalctl -u "$SERVICE_NAME" -n 20 --no-pager
            exit 1
        fi
    fi
}

# Docker 构建
docker_build() {
    log_info "构建 API Server Docker 镜像..."
    
    cd "$PROJECT_DIR"
    VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
    
    docker build \
        -f build/docker/apiserver/Dockerfile \
        -t ibookfs-apiserver:${VERSION} \
        -t ibookfs-apiserver:latest \
        .
    
    log_success "镜像构建完成: ibookfs-apiserver:${VERSION}"
}

# Docker 部署
docker_deploy() {
    log_info "Docker 部署 API Server..."
    
    cd "$PROJECT_DIR/build/docker/apiserver"
    docker compose up -d
    
    log_success "API Server 容器启动完成"
    docker compose ps
}

# 显示帮助
show_help() {
    cat << EOF
用法: $0 [命令]

命令:
    build       仅构建后端
    deploy      构建并部署到系统（systemd）
    docker      构建 Docker 镜像
    docker-run  Docker 部署运行
    all         完整流程（构建+部署+启动）

选项:
    -h, --help  显示帮助

环境变量:
    DEPLOY_DIR  部署目录（默认: /opt/ibookfs）
    VERSION     版本号（默认: git tag）

示例:
    $0 build        # 仅构建
    $0 deploy       # 构建并部署
    $0 docker       # 构建 Docker 镜像
    $0 docker-run   # Docker 部署
EOF
}

# 主流程
main() {
    local cmd="${1:-all}"
    
    echo ""
    echo "========================================"
    echo "  iBookFS 后端部署"
    echo "========================================"
    echo ""
    
    case $cmd in
        build)
            check_go
            build_backend
            ;;
        deploy)
            check_go
            build_backend
            deploy_backend
            setup_service
            start_service
            ;;
        docker)
            docker_build
            ;;
        docker-run)
            docker_build
            docker_deploy
            ;;
        all)
            check_go
            build_backend
            deploy_backend
            setup_service
            start_service
            ;;
        -h|--help)
            show_help
            ;;
        *)
            log_error "未知命令: $cmd"
            show_help
            exit 1
            ;;
    esac
    
    echo ""
    log_success "完成！"
}

main "$@"
