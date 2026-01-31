#!/bin/bash

################################################################################
# iBookFS 环境检查脚本
#
# 功能：检查项目部署所需的环境依赖
################################################################################

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 计数器
WARN_COUNT=0
ERROR_COUNT=0

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[⚠]${NC} $1"
    ((WARN_COUNT++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((ERROR_COUNT++))
}

log_section() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

# 检查命令是否存在
check_command() {
    local cmd=$1
    local min_version=${2:-""}
    local package=$3

    if command -v "$cmd" &> /dev/null; then
        local version=$($cmd --version 2>&1 | head -n1 || echo "unknown")

        if [ -n "$min_version" ]; then
            # 简单的版本比较（需要根据具体命令调整）
            log_success "$cmd 已安装 (版本: $version)"
        else
            log_success "$cmd 已安装"
        fi
        return 0
    else
        log_error "$cmd 未安装"
        if [ -n "$package" ]; then
            echo "     安装: $package"
        fi
        return 1
    fi
}

# 检查 Go 环境
check_go() {
    log_section "检查 Go 环境"

    local go_version="1.21"

    if command -v go &> /dev/null; then
        local installed_version=$(go version | awk '{print $3}' | sed 's/go//')
        local major=$(echo $installed_version | cut -d. -f1)
        local minor=$(echo $installed_version | cut -d. -f2)

        if [ "$major" -gt 1 ] || ([ "$major" -eq 1 ] && [ "$minor" -ge 21 ]); then
            log_success "Go 已安装 (版本: go$installed_version)"

            # 检查 GOPATH 和 GOROOT
            if [ -n "$GOPATH" ]; then
                log_info "GOPATH: $GOPATH"
            else
                log_warn "GOPATH 未设置"
            fi

            # 检查 Go modules 代理
            if [ -z "$GOPROXY" ] && [ -z "$GO_NOPROXY" ]; then
                log_warn "建议设置 GOPROXY (如: https://goproxy.cn,direct)"
            fi
        else
            log_error "Go 版本过低 (当前: go$installed_version, 需要: >= go$go_version)"
            echo "     升级: 访问 https://go.dev/dl/"
        fi
    else
        log_error "Go 未安装 (需要: >= $go_version)"
        echo "     安装: https://go.dev/dl/"
        echo "     或运行: ./install-env.sh go"
    fi
}

# 检查 Node.js 环境
check_nodejs() {
    log_section "检查 Node.js 环境"

    local node_version="18"
    local npm_version="9"

    if command -v node &> /dev/null; then
        local node_ver=$(node -v | sed 's/v//' | sed 's/\..*//')
        if [ "$node_ver" -ge "$node_version" ]; then
            log_success "Node.js 已安装 (版本: $(node -v))"
        else
            log_error "Node.js 版本过低 (当前: $(node -v), 需要: >= v$node_version)"
            echo "     升级: 访问 https://nodejs.org/"
        fi
    else
        log_error "Node.js 未安装 (需要: >= v$node_version)"
        echo "     安装: https://nodejs.org/"
        echo "     或运行: ./install-env.sh node"
    fi

    if command -v npm &> /dev/null; then
        local npm_ver=$(npm -v | sed 's/\..*//')
        if [ "$npm_ver" -ge "$npm_version" ]; then
            log_success "npm 已安装 (版本: $(npm -v))"
        else
            log_warn "npm 版本过低 (当前: $(npm -v), 建议: >= v$npm_version)"
        fi
    else
        log_error "npm 未安装"
    fi

    # 检查 yarn/pnpm（可选）
    if command -v pnpm &> /dev/null; then
        log_info "pnpm 已安装 ($(pnpm -v))"
    elif command -v yarn &> /dev/null; then
        log_info "yarn 已安装 ($(yarn -v))"
    fi
}

# 检查 Git
check_git() {
    log_section "检查 Git"

    if command -v git &> /dev/null; then
        log_success "Git 已安装 (版本: $(git --version))"
    else
        log_error "Git 未安装"
        echo "     安装: sudo apt-get install git"
        echo "     或运行: ./install-env.sh git"
    fi
}

# 检查系统工具
check_system_tools() {
    log_section "检查系统工具"

    local tools=(
        "curl:curl"
        "wget:wget"
        "tar:tar"
        "gzip:gzip"
    )

    for tool_info in "${tools[@]}"; do
        IFS=':' read -r cmd pkg <<< "$tool_info"
        check_command "$cmd" "" "$pkg"
    done
}

# 检查端口占用
check_ports() {
    log_section "检查端口占用"

    # 检查 8080 端口
    if command -v ss &> /dev/null; then
        if ss -tuln 2>/dev/null | grep -q ":8080 "; then
            # 端口被占用，检查是否是 ibookfs 服务
            local pid=$(ss -tulnp 2>/dev/null | grep ":8080 " | grep -oP 'pid=\K[0-9]+' | head -1)
            if [ -n "$pid" ]; then
                local process_name=$(ps -p "$pid" -o comm= 2>/dev/null)
                if [ "$process_name" = "ibookfs" ]; then
                    log_success "端口 8080 已被 ibookfs 服务使用 (部署时会自动重启)"
                else
                    log_warn "端口 8080 被其他程序占用 ($process_name, PID: $pid)"
                    echo "     停止: sudo kill $pid"
                fi
            else
                log_warn "端口 8080 已被占用，无法确定进程"
            fi
        else
            log_success "端口 8080 可用 (后端服务)"
        fi
    elif command -v netstat &> /dev/null; then
        if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
            local pid=$(lsof -ti:8080 2>/dev/null | head -1)
            if [ -n "$pid" ]; then
                local process_name=$(ps -p "$pid" -o comm= 2>/dev/null)
                if [ "$process_name" = "ibookfs" ]; then
                    log_success "端口 8080 已被 ibookfs 服务使用 (部署时会自动重启)"
                else
                    log_warn "端口 8080 被其他程序占用 ($process_name, PID: $pid)"
                    echo "     停止: sudo kill $pid"
                fi
            else
                log_warn "端口 8080 已被占用，无法确定进程"
            fi
        else
            log_success "端口 8080 可用 (后端服务)"
        fi
    fi

    # 检查 3000 端口 (仅警告，不阻塞)
    if command -v ss &> /dev/null; then
        if ss -tuln 2>/dev/null | grep -q ":3000 "; then
            log_info "端口 3000 已被占用 (前端开发服务器)"
        else
            log_success "端口 3000 可用 (前端开发服务器)"
        fi
    fi
}

# 检查磁盘空间
check_disk_space() {
    log_section "检查磁盘空间"

    local min_space_gb=5

    if command -v df &> /dev/null; then
        local available_gb=$(df / | awk 'NR==2 {printf "%.0f", $4/1024/1024}')

        if [ "$available_gb" -ge "$min_space_gb" ]; then
            log_success "磁盘空间充足 (可用: ${available_gb}GB)"
        else
            log_warn "磁盘空间不足 (可用: ${available_gb}GB, 建议: >= ${min_space_gb}GB)"
        fi
    fi
}

# 检查目录权限
check_permissions() {
    log_section "检查目录权限"

    local dirs=(
        "/opt/ibookfs:部署目录"
        "/var/log:日志目录"
    )

    for dir_info in "${dirs[@]}"; do
        IFS=':' read -r dir desc <<< "$dir_info"

        if [ -d "$dir" ]; then
            if [ -w "$dir" ]; then
                log_success "$dir 可写 ($desc)"
            else
                log_warn "$dir 不可写 ($desc)"
                echo "     修复: sudo chmod 755 $dir"
            fi
        else
            log_info "$dir 不存在 ($desc)"
            echo "     创建: sudo mkdir -p $dir"
        fi
    done
}

# 检查网络连接
check_network() {
    log_section "检查网络连接"

    local urls=(
        "github.com:GitHub"
        "goproxy.cn:Go代理"
    )

    for url_info in "${urls[@]}"; do
        IFS=':' read -r url desc <<< "$url_info"

        # 使用函数封装，避免 set -e 导致脚本退出
        local check_result=0
        (curl -s --connect-timeout 5 --max-time 10 "https://$url" > /dev/null 2>&1 || \
         curl -s --connect-timeout 5 --max-time 10 "http://$url" > /dev/null 2>&1) || check_result=$?

        if [ $check_result -eq 0 ]; then
            log_success "$desc 可访问"
        else
            log_warn "$desc 连接失败 (超时或不可达)"
        fi
    done
}

# 显示帮助信息
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    -v, --verbose       详细输出
    -f, --fix          尝试自动修复问题
    --skip-ports       跳过端口检查
    --skip-disk        跳过磁盘检查

示例:
    $0                  # 执行完整检查
    $0 -f               # 检查并尝试修复问题
    $0 --skip-ports     # 跳过端口检查
EOF
}

# 解析命令行参数
VERBOSE=false
FIX_MODE=false
SKIP_PORTS=false
SKIP_DISK=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -f|--fix)
            FIX_MODE=true
            shift
            ;;
        --skip-ports)
            SKIP_PORTS=true
            shift
            ;;
        --skip-disk)
            SKIP_DISK=true
            shift
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 主流程
main() {
    echo ""
    echo "========================================"
    echo "  iBookFS 环境检查工具"
    echo "========================================"
    echo ""

    # 执行检查
    check_system_tools
    check_git
    check_go
    check_nodejs
    check_network

    if [ "$SKIP_PORTS" = false ]; then
        check_ports
    fi

    if [ "$SKIP_DISK" = false ]; then
        check_disk_space
    fi

    check_permissions

    # 输出结果摘要
    echo ""
    echo "========================================"
    echo "  检查摘要"
    echo "========================================"
    echo ""

    if [ $ERROR_COUNT -eq 0 ] && [ $WARN_COUNT -eq 0 ]; then
        log_success "所有检查通过！环境已准备就绪。"
        echo ""
        echo "下一步: 运行 ./scripts/build.sh 开始构建部署"
        exit 0
    elif [ $ERROR_COUNT -eq 0 ]; then
        log_warn "发现 $WARN_COUNT 个警告，但可以继续部署。"
        echo ""
        echo "建议: 运行 ./scripts/install-env.sh 安装缺失的依赖"
        echo "或运行: $0 -f 尝试自动修复"
        exit 0
    else
        log_error "发现 $ERROR_COUNT 个错误，需要解决后才能部署。"
        echo ""
        echo "解决方案:"
        echo "  1. 手动安装缺失的依赖"
        echo "  2. 运行 ./scripts/install-env.sh 自动安装"
        echo "  3. 运行 $0 -f 尝试自动修复"
        exit 1
    fi
}

# 运行主流程
main
