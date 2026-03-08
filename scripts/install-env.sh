#!/bin/bash

################################################################################
# iBookFS 环境安装脚本
#
# 功能：自动安装项目所需的环境依赖
# 支持 Ubuntu/Debian、CentOS/RHEL 和 macOS 系统
################################################################################

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS_ID="macos"
        OS_VERSION=$(sw_vers -productVersion)
    elif [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_ID=$ID
        OS_VERSION=$VERSION_ID
    else
        log_error "无法检测操作系统类型"
        exit 1
    fi

    log_info "检测到操作系统: $OS_ID $OS_VERSION"
}

# 检查是否为 root 用户（macOS 不强制要求）
check_root() {
    if [ "$OS_ID" = "macos" ]; then
        log_info "macOS 系统，使用 Homebrew 安装（无需 root）"
        return
    fi
    if [ "$EUID" -ne 0 ]; then
        log_error "此脚本需要 root 权限运行"
        echo "请使用: sudo $0"
        exit 1
    fi
}

# 更新包管理器索引
update_package_manager() {
    log_info "更新包管理器索引..."

    case $OS_ID in
        macos)
            if ! command -v brew &> /dev/null; then
                log_info "安装 Homebrew..."
                /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
            fi
            brew update
            ;;
        ubuntu|debian)
            export DEBIAN_FRONTEND=noninteractive
            apt-get update -qq
            ;;
        centos|rhel|fedora)
            yum makecache fast -q 2>/dev/null || dnf makecache fast -q
            ;;
    esac

    log_success "包管理器索引已更新"
}

# 安装基础工具
install_base_tools() {
    log_info "安装基础工具..."

    case $OS_ID in
        macos)
            brew install curl wget git
            ;;
        ubuntu|debian)
            apt-get install -y -qq curl wget git tar gzip build-essential
            ;;
        centos|rhel|fedora)
            yum install -y -q curl wget git tar gzip gcc make
            ;;
    esac

    log_success "基础工具安装完成"
}

# 安装 Go
install_go() {
    if command -v go &> /dev/null; then
        log_success "Go 已安装，跳过"
        return
    fi

    log_info "安装 Go..."

    case $OS_ID in
        macos)
            brew install go
            ;;
        *)
            local go_version="1.25.6"
            local go_arch="amd64"

            # 下载 Go
            cd /tmp
            curl -sLO "https://go.dev/dl/go${go_version}.linux-${go_arch}.tar.gz"

            # 解压并安装
            tar -C /usr/local -xzf "go${go_version}.linux-${go_arch}.tar.gz"
            rm -f "go${go_version}.linux-${go_arch}.tar.gz"

            # 设置环境变量
            cat <<'EOF' > /etc/profile.d/go.sh
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
EOF

            # 立即生效
            export PATH=$PATH:/usr/local/go/bin
            export GOPATH=$HOME/go
            ;;
    esac

    # 验证安装
    if command -v go &> /dev/null; then
        log_success "Go 安装完成 ($(go version))"
    else
        log_error "Go 安装失败"
        exit 1
    fi
}

# 安装 Node.js
install_nodejs() {
    if command -v node &> /dev/null && command -v npm &> /dev/null; then
        log_success "Node.js 已安装，跳过"
        return
    fi

    log_info "安装 Node.js..."

    case $OS_ID in
        macos)
            brew install node
            ;;
        ubuntu|debian)
            # 安装必要的工具
            apt-get install -y -qq ca-certificates curl gnupg

            # 添加 NodeSource GPG 密钥
            mkdir -p /etc/apt/keyrings
            curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource-repo.gpg

            # 添加 NodeSource 仓库
            echo "deb [signed-by=/etc/apt/keyrings/nodesource-repo.gpg] https://deb.nodesource.com/node_20.x nodistro main" \
                > /etc/apt/sources.list.d/nodesource.list

            # 更新并安装
            export DEBIAN_FRONTEND=noninteractive
            apt-get update -qq
            apt-get install -y -qq nodejs
            ;;

        centos|rhel|fedora)
            # 添加 NodeSource 仓库
            cat <<'EOF' > /etc/yum.repos.d/nodesource-el.repo
[nodesource]
name=Node.js Packages for Enterprise Linux
baseurl=https://repo.nodesource.com/20.x/el/$releasever/
enabled=1
gpgcheck=1
gpgkey=https://rpm.nodesource.com/pub/el/NODESOURCE-GPG-SIGNING-KEY
EOF

            # 安装
            yum install -y nodejs npm
            ;;
    esac

    # 验证安装
    if command -v node &> /dev/null && command -v npm &> /dev/null; then
        log_success "Node.js 安装完成 ($(node -v), npm $(npm -v))"
    else
        log_error "Node.js 安装失败"
        exit 1
    fi
}

# 配置 Git
configure_git() {
    if ! command -v git &> /dev/null; then
        return
    fi

    log_info "配置 Git..."

    # 设置全局配置
    git config --global core.autocrlf true
    git config --global init.defaultBranch main

    log_success "Git 配置完成"
}

# 安装 Nginx
install_nginx() {
    if command -v nginx &> /dev/null; then
        log_success "Nginx 已安装，跳过"
        return
    fi

    log_info "安装 Nginx..."

    case $OS_ID in
        macos)
            brew install nginx
            ;;
        ubuntu|debian)
            apt-get install -y -qq nginx
            ;;
        centos|rhel|fedora)
            yum install -y nginx
            ;;
    esac

    # 验证安装
    if command -v nginx &> /dev/null; then
        log_success "Nginx 安装完成 ($(nginx -v 2>&1))"
    else
        log_error "Nginx 安装失败"
        exit 1
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."

    case $OS_ID in
        macos)
            mkdir -p /usr/local/var/ibookfs
            mkdir -p /usr/local/var/log/ibookfs
            mkdir -p /usr/local/var/www/ibookfs
            ;;
        *)
            mkdir -p /opt/ibookfs
            mkdir -p /var/log/ibookfs
            mkdir -p /var/www/ibookfs
            ;;
    esac

    log_success "目录创建完成"
}

# 设置 Go 代理（国内用户）
setup_go_proxy() {
    log_info "配置 Go 代理..."

    cat <<'EOF' > /etc/profile.d/go.sh
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
EOF

    log_success "Go 代理配置完成"
}

# 显示帮助信息
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    --skip-go          跳过 Go 安装
    --skip-nodejs       跳过 Node.js 安装
    --skip-git          跳过 Git 安装
    --skip-tools        跳过基础工具安装
    --skip-dirs         跳过创建目录

示例:
    $0                  # 安装所有依赖
    $0 --skip-go        # 跳过 Go 安装
    $0 --skip-nodejs     # 跳过 Node.js 安装
EOF
}

# 解析命令行参数
SKIP_GO=false
SKIP_NODEJS=false
SKIP_NGINX=false
SKIP_GIT=false
SKIP_TOOLS=false
SKIP_DIRS=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --skip-go)
            SKIP_GO=true
            shift
            ;;
        --skip-nodejs)
            SKIP_NODEJS=true
            shift
            ;;
        --skip-git)
            SKIP_GIT=true
            shift
            ;;
        --skip-tools)
            SKIP_TOOLS=true
            shift
            ;;
        --skip-dirs)
            SKIP_DIRS=true
            shift
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 主流程
main() {
    echo ""
    echo "========================================"
    echo "  iBookFS 环境安装工具"
    echo "========================================"
    echo ""

    # 检查 root 权限
    check_root

    # 检测操作系统
    detect_os

    # 更新包管理器
    update_package_manager

    # 安装基础工具
    if [ "$SKIP_TOOLS" = false ]; then
        install_base_tools
    fi

    # 安装 Go
    if [ "$SKIP_GO" = false ]; then
        install_go
        setup_go_proxy
    fi

    # 安装 Node.js
    if [ "$SKIP_NODEJS" = false ]; then
        install_nodejs
    fi

    # 安装 Nginx
    if [ "$SKIP_NGINX" = false ]; then
        install_nginx
    fi

    # 配置 Git
    if [ "$SKIP_GIT" = false ]; then
        configure_git
    fi

    # 创建目录
    if [ "$SKIP_DIRS" = false ]; then
        create_directories
    fi

    # 显示摘要
    echo ""
    echo "========================================"
    log_success "环境安装完成！"
    echo ""
    echo "已安装的工具:"
    echo "  - Go:       $(go version 2>/dev/null | head -n1 || echo '未安装')"
    echo "  - Node.js:   $(node -v 2>/dev/null || echo '未安装')"
    echo "  - npm:      $(npm -v 2>/dev/null || echo '未安装')"
    echo "  - Git:      $(git --version 2>/dev/null | head -n1 || echo '未安装')"
    echo ""
    echo "下一步:"
    echo "  1. 运行 ./scripts/check-env.sh 检查环境"
    echo "  2. 运行 ./scripts/build.sh 开始构建部署"
    echo "========================================"
    echo ""
}

# 运行主流程
main
