#!/bin/bash

################################################################################
# iBookFS 前端独立部署脚本
#
# 功能：仅构建和部署前端服务
# 支持 Linux 和 macOS 系统
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
    DEPLOY_DIR="${DEPLOY_DIR:-$HOME/.ibookfs/web}"
    # Homebrew Nginx 配置目录（支持 Intel 和 Apple Silicon）
    if [ -d "/opt/homebrew/etc/nginx/servers" ]; then
        NGINX_CONF_DIR="/opt/homebrew/etc/nginx/servers"
    else
        NGINX_CONF_DIR="/usr/local/etc/nginx/servers"
    fi
else
    OS_TYPE="linux"
    DEPLOY_DIR="${DEPLOY_DIR:-/var/www/ibookfs}"
    NGINX_CONF_DIR="/etc/nginx/sites-available"
fi

# 配置变量
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
WEB_DIR="$PROJECT_DIR/web"
NGINX_DIR="$DEPLOY_DIR"

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查并安装 Nginx
check_nginx() {
    log_info "检查 Nginx 环境..."
    
    if command -v nginx &> /dev/null; then
        log_success "Nginx 已安装 ($(nginx -v 2>&1 | head -1))"
        return 0
    fi
    
    log_warn "Nginx 未安装，正在安装..."
    
    if [ "$OS_TYPE" = "macos" ]; then
        if ! command -v brew &> /dev/null; then
            log_error "Homebrew 未安装，请先安装 Homebrew"
            exit 1
        fi
        brew install nginx
    else
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y nginx
        elif command -v yum &> /dev/null; then
            sudo yum install -y nginx
        else
            log_error "无法自动安装 Nginx，请手动安装"
            exit 1
        fi
    fi
    
    if command -v nginx &> /dev/null; then
        log_success "Nginx 安装完成"
    else
        log_error "Nginx 安装失败"
        exit 1
    fi
}

# 检查 Node.js 环境
check_node() {
    log_info "检查 Node.js 环境..."
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi
    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装"
        exit 1
    fi
    log_success "Node.js $(node -v), npm $(npm -v)"
}

# 构建前端
build_frontend() {
    log_info "构建前端..."
    
    cd "$WEB_DIR"
    
    log_info "安装依赖..."
    npm ci --silent
    
    log_info "执行构建..."
    npm run build
    
    if [ ! -d "$WEB_DIR/dist" ]; then
        log_error "构建失败"
        exit 1
    fi
    
    log_success "前端构建完成: $WEB_DIR/dist"
}

# 部署到本地目录
deploy_local() {
    log_info "部署前端到 $DEPLOY_DIR..."
    
    mkdir -p "$DEPLOY_DIR"
    rm -rf "$DEPLOY_DIR"/*
    cp -r "$WEB_DIR/dist/"* "$DEPLOY_DIR/"
    
    # 复制 assets（如果有）
    if [ -d "$WEB_DIR/assets" ]; then
        cp -r "$WEB_DIR/assets" "$DEPLOY_DIR/"
    fi
    
    log_success "前端部署完成"
}

# 部署到 Nginx
deploy_nginx() {
    log_info "部署前端到 Nginx ($NGINX_DIR)..."
    
    mkdir -p "$NGINX_DIR"
    rm -rf "$NGINX_DIR"/*
    cp -r "$WEB_DIR/dist/"* "$NGINX_DIR/"
    
    if [ -d "$WEB_DIR/assets" ]; then
        cp -r "$WEB_DIR/assets" "$NGINX_DIR/"
    fi
    
    # 配置 Nginx
    setup_nginx_config
    
    # 重载 Nginx
    reload_nginx
    
    log_success "前端部署到 Nginx 完成"
}

# 配置 Nginx（根据操作系统）
setup_nginx_config() {
    local nginx_conf="$NGINX_CONF_DIR/ibookfs.conf"
    
    if [ -f "$nginx_conf" ]; then
        log_info "Nginx 配置已存在: $nginx_conf"
        return
    fi
    
    log_info "创建 Nginx 配置..."
    mkdir -p "$NGINX_CONF_DIR"
    
    # macOS 使用非特权端口 3000，Linux 使用 80
    local listen_port="80"
    if [ "$OS_TYPE" = "macos" ]; then
        listen_port="3000"
    fi
    
    cat > "$nginx_conf" <<EOF
server {
    listen $listen_port;
    server_name _;
    root $NGINX_DIR;
    index index.html;

    location /health {
        access_log off;
        return 200 'OK';
        add_header Content-Type text/plain;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)\$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location / {
        try_files \$uri \$uri/ /index.html;
    }
}
EOF
    
    # Linux 需要创建软链接
    if [ "$OS_TYPE" = "linux" ]; then
        sudo ln -sf "$nginx_conf" /etc/nginx/sites-enabled/ibookfs.conf
    fi
    
    log_success "Nginx 配置创建完成: $nginx_conf"
}

# 重载 Nginx
reload_nginx() {
    if ! command -v nginx &> /dev/null; then
        log_warn "Nginx 未安装，跳过重载"
        return
    fi
    
    log_info "测试 Nginx 配置..."
    if ! nginx -t 2>/dev/null; then
        log_error "Nginx 配置测试失败"
        nginx -t
        exit 1
    fi
    
    log_info "重载 Nginx..."
    if [ "$OS_TYPE" = "macos" ]; then
        brew services restart nginx 2>/dev/null || nginx -s reload
    else
        sudo systemctl reload nginx
    fi
    
    log_success "Nginx 重载完成"
}

# Docker 构建
docker_build() {
    log_info "构建 Web Docker 镜像..."
    
    cd "$PROJECT_DIR"
    VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
    
    docker build \
        -f build/docker/web/Dockerfile \
        -t ibookfs-web:${VERSION} \
        -t ibookfs-web:latest \
        .
    
    log_success "镜像构建完成: ibookfs-web:${VERSION}"
}

# Docker 部署
docker_deploy() {
    log_info "Docker 部署 Web..."
    
    cd "$PROJECT_DIR/build/docker/web"
    docker compose up -d
    
    log_success "Web 容器启动完成"
    docker compose ps
}

# 显示帮助
show_help() {
    cat << EOF
用法: $0 [命令]

命令:
    build       仅构建前端
    deploy      构建并部署到本地目录
    nginx       构建并部署到 Nginx
    docker      构建 Docker 镜像
    docker-run  Docker 部署运行

选项:
    -h, --help  显示帮助

环境变量:
    DEPLOY_DIR  部署目录（默认: /opt/ibookfs/web）
    NGINX_DIR   Nginx 目录（默认: /var/www/ibookfs）
    VERSION     版本号（默认: git tag）

示例:
    $0 build        # 仅构建
    $0 deploy       # 部署到本地目录
    $0 nginx        # 部署到 Nginx
    $0 docker       # 构建 Docker 镜像
    $0 docker-run   # Docker 部署
EOF
}

# 主流程
main() {
    local cmd="${1:-deploy}"
    
    echo ""
    echo "========================================"
    echo "  iBookFS 前端部署"
    echo "========================================"
    echo ""
    
    case $cmd in
        build)
            check_node
            build_frontend
            ;;
        deploy)
            check_node
            build_frontend
            deploy_local
            ;;
        nginx)
            check_nginx
            check_node
            build_frontend
            deploy_nginx
            ;;
        docker)
            docker_build
            ;;
        docker-run)
            docker_build
            docker_deploy
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
