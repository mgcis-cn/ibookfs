#!/bin/bash

################################################################################
# iBookFS Build and Deploy Script
#
# 功能：
# 1. 从 GitHub 获取最新代码
# 2. 构建 Go 后端
# 3. 构建前端（npm）
# 4. 部署到 /opt/ibookfs
# 5. 配置 systemd 服务
################################################################################

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置变量
REPO_URL="https://github.com/mgcis-cn/ibookfs.git"
BUILD_DIR="/tmp/ibookfs-build"
DEPLOY_DIR="/opt/ibookfs"
SERVICE_NAME="ibookfs"
BRANCH="master"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查构建依赖..."

    if ! command -v git &> /dev/null; then
        log_error "git 未安装，请先安装 git"
        exit 1
    fi

    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go (>= 1.21)"
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装，请先安装 Node.js (>= 18)"
        exit 1
    fi

    log_info "依赖检查通过 ✓"
}

# 获取代码
fetch_code() {
    # 检查是否在项目根目录
    if [ -f "$(dirname "$0")/../go.mod" ] && [ -f "$(dirname "$0")/../cmd/server/main.go" ]; then
        log_info "使用本地源码构建..."
        BUILD_DIR="$(cd "$(dirname "$0")/.." && pwd)"
        log_info "源码目录: $BUILD_DIR"
        log_info "代码获取完成 ✓ (本地模式)"
        return
    fi

    log_info "从 GitHub 获取代码 (branch: $BRANCH)..."

    # 清理旧的构建目录
    if [ -d "$BUILD_DIR" ]; then
        log_warn "清理旧的构建目录: $BUILD_DIR"
        rm -rf "$BUILD_DIR"
    fi

    # 克隆代码
    git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$BUILD_DIR"

    # 检查克隆是否成功
    if [ ! -f "$BUILD_DIR/go.mod" ]; then
        log_error "GitHub 仓库似乎为空或未正确初始化"
        echo ""
        echo "解决方案："
        echo "  1. 将代码推送到 GitHub: https://github.com/mgcis-cn/ibookfs"
        echo "  2. 或在项目根目录运行此脚本（自动使用本地源码）"
        echo ""
        exit 1
    fi

    log_info "代码获取完成 ✓"
}

# 构建后端
build_backend() {
    log_info "构建 Go 后端..."

    cd "$BUILD_DIR"

    # 创建输出目录
    mkdir -p "$BUILD_DIR/bin"

    # 设置 Go 代理（加速国内下载）
    export GOPROXY=https://goproxy.cn,direct
    export GOSUMDB=sum.golang.google.cn
    export GOTOOLCHAIN=auto

    log_info "整理 Go 模块依赖..."
    go mod tidy

    log_info "下载 Go 模块依赖..."
    go mod download -x

    log_info "验证 Go 模块..."
    go mod verify || log_warn "模块验证失败，继续构建..."

    log_info "开始编译..."
    log_warn "注意: github.com/ugorji/go/codec 包编译较慢，可能需要 5-10 分钟，请耐心等待..."
    # 构建 Go 二进制 (使用多核并行编译)
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -p=$(nproc) \
        -ldflags="-s -w" \
        -o "$BUILD_DIR/bin/ibookfs" \
        ./cmd/server

    log_info "编译命令执行完成"

    # 检查构建结果
    if [ ! -f "$BUILD_DIR/bin/ibookfs" ]; then
        log_error "后端构建失败"
        exit 1
    fi

    log_info "后端构建完成 ✓"
}

# 构建前端
build_frontend() {
    log_info "构建前端..."

    cd "$BUILD_DIR/web"

    # 安装依赖
    log_info "安装 npm 依赖..."
    npm ci --silent

    # 构建
    log_info "执行 npm build..."
    npm run build

    # 检查构建结果
    if [ ! -d "$BUILD_DIR/web/dist" ]; then
        log_error "前端构建失败"
        exit 1
    fi

    log_info "前端构建完成 ✓"
}

# 部署到生产目录
deploy() {
    log_info "部署到 $DEPLOY_DIR..."

    # 创建部署目录（如果不存在）
    sudo mkdir -p "$DEPLOY_DIR/web"

    # 停止现有服务（如果运行）
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "停止现有服务..."
        sudo systemctl stop "$SERVICE_NAME"
    fi

    # 复制后端二进制
    log_info "复制后端二进制..."
    sudo cp "$BUILD_DIR/bin/ibookfs" "$DEPLOY_DIR/"
    sudo chmod +x "$DEPLOY_DIR/ibookfs"

    # 复制前端静态文件到 web 目录
    log_info "复制前端静态文件..."
    sudo rm -rf "$DEPLOY_DIR/web/"*
    sudo cp -r "$BUILD_DIR/web/dist/"* "$DEPLOY_DIR/web/"

    # 复制 assets 静态资源到 web 目录
    if [ -d "$BUILD_DIR/web/assets" ]; then
        log_info "复制 assets 静态资源..."
        sudo rm -rf "$DEPLOY_DIR/web/assets"
        sudo cp -r "$BUILD_DIR/web/assets" "$DEPLOY_DIR/web/"
    fi

    # 复制配置文件（如果不存在）
    if [ ! -f "$DEPLOY_DIR/config.yaml" ]; then
        if [ -f "$BUILD_DIR/config.yaml" ]; then
            log_info "复制配置文件..."
            sudo cp "$BUILD_DIR/config.yaml" "$DEPLOY_DIR/"
        else
            log_warn "config.yaml 不存在，请手动配置"
        fi
    fi

    log_info "部署完成 ✓"
    log_info "  - 后端: $DEPLOY_DIR/ibookfs"
    log_info "  - 前端: $DEPLOY_DIR/web/"
}

# 配置 systemd 服务
setup_systemd() {
    log_info "配置 systemd 服务..."

    cat <<EOF | sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null
[Unit]
Description=iBookFS Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DEPLOY_DIR
ExecStart=$DEPLOY_DIR/ibookfs
Restart=on-failure
RestartSec=5s

# 环境变量
Environment=GIN_MODE=release

# 安全加固
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载 systemd 配置
    sudo systemctl daemon-reload

    # 启用服务（开机自启）
    sudo systemctl enable "$SERVICE_NAME"

    log_info "systemd 服务配置完成 ✓"
}

# 启动服务
start_service() {
    log_info "启动 $SERVICE_NAME 服务..."

    sudo systemctl start "$SERVICE_NAME"

    # 等待服务启动
    sleep 2

    # 检查服务状态
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "服务启动成功 ✓"

        # 显示服务状态
        echo ""
        systemctl status "$SERVICE_NAME" --no-pager
        echo ""

        # 显示日志（最近 10 行）
        log_info "最近日志："
        sudo journalctl -u "$SERVICE_NAME" -n 10 --no-pager
    else
        log_error "服务启动失败"
        sudo journalctl -u "$SERVICE_NAME" -n 20 --no-pager
        exit 1
    fi
}

# 清理构建目录
cleanup() {
    log_info "清理构建目录..."
    rm -rf "$BUILD_DIR"
    log_info "清理完成 ✓"
}

# 显示帮助信息
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    -s, --skip-build    跳过构建，仅部署已有文件
    -b, --backend-only  仅构建后端
    -f, --frontend-only 仅构建前端
    -n, --no-restart    部署后不重启服务
    --branch BRANCH    指定分支（默认: master）

示例:
    $0                  # 完整构建和部署
    $0 -s               # 跳过构建，仅部署
    $0 -b               # 仅构建后端
    $0 -f               # 仅构建前端
    $0 --branch dev    # 使用 dev 分支
EOF
}

# 解析命令行参数
SKIP_BUILD=false
BACKEND_ONLY=false
FRONTEND_ONLY=false
NO_RESTART=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -s|--skip-build)
            SKIP_BUILD=true
            shift
            ;;
        -b|--backend-only)
            BACKEND_ONLY=true
            shift
            ;;
        -f|--frontend-only)
            FRONTEND_ONLY=true
            shift
            ;;
        -n|--no-restart)
            NO_RESTART=true
            shift
            ;;
        --branch)
            BRANCH="$2"
            shift 2
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 检查环境（调用 check-env.sh）
check_environment() {
    log_info "检查环境依赖..."

    local check_script="$(dirname "$0")/check-env.sh"

    if [ -f "$check_script" ]; then
        if bash "$check_script"; then
            log_success "环境检查通过"
        else
            log_error "环境检查失败"
            echo ""
            echo "解决方案："
            echo "  1. 手动安装缺失的依赖"
            echo "  2. 运行: sudo ./scripts/install-env.sh"
            echo ""
            exit 1
        fi
    else
        log_warn "未找到环境检查脚本，跳过..."
    fi
}

# 主流程
main() {
    echo ""
    echo "=========================================="
    echo "  iBookFS 自动构建部署脚本"
    echo "=========================================="
    echo ""

    # 检查环境
    check_environment

    # 检查依赖
    check_dependencies

    # 获取代码并构建
    if [ "$SKIP_BUILD" = false ]; then
        fetch_code

        if [ "$FRONTEND_ONLY" = false ]; then
            build_backend
        fi

        if [ "$BACKEND_ONLY" = false ]; then
            build_frontend
        fi
    else
        log_warn "跳过构建，使用本地已有文件"
        BUILD_DIR="/opt/cn.mgcis/ibookfs"
    fi

    # 部署
    deploy

    # 配置 systemd 服务
    setup_systemd

    # 启动服务
    if [ "$NO_RESTART" = false ]; then
        start_service
    else
        log_info "跳过服务重启"
    fi

    # 清理构建目录
    if [ "$SKIP_BUILD" = false ]; then
        cleanup
    fi

    echo ""
    echo "=========================================="
    log_info "构建部署完成！"
    echo ""
    echo "服务管理命令："
    echo "  启动服务: sudo systemctl start $SERVICE_NAME"
    echo "  停止服务: sudo systemctl stop $SERVICE_NAME"
    echo "  重启服务: sudo systemctl restart $SERVICE_NAME"
    echo "  查看状态: sudo systemctl status $SERVICE_NAME"
    echo "  查看日志: sudo journalctl -u $SERVICE_NAME -f"
    echo "=========================================="
    echo ""
}

# 运行主流程
main
