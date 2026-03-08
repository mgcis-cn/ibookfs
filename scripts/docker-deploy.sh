#!/bin/bash

################################################################################
# iBookFS Docker 部署脚本
#
# 功能：
# 1. 构建 Docker 镜像
# 2. 推送到镜像仓库（可选）
# 3. 在目标服务器部署
################################################################################

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置变量
APP_NAME="ibookfs"
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
DOCKER_REGISTRY="${DOCKER_REGISTRY:-}"
DOCKER_IMAGE="${APP_NAME}"

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 Docker 环境
check_docker() {
    log_info "检查 Docker 环境..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        echo "安装: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行"
        echo "启动: sudo systemctl start docker"
        exit 1
    fi
    
    log_success "Docker 环境正常"
}

# 构建镜像
build_image() {
    log_info "构建 Docker 镜像: ${DOCKER_IMAGE}:${VERSION}..."
    
    local image_tag="${DOCKER_IMAGE}:${VERSION}"
    local latest_tag="${DOCKER_IMAGE}:latest"
    
    if [ -n "$DOCKER_REGISTRY" ]; then
        image_tag="${DOCKER_REGISTRY}/${image_tag}"
        latest_tag="${DOCKER_REGISTRY}/${latest_tag}"
    fi
    
    docker build \
        --build-arg VERSION="${VERSION}" \
        --build-arg BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -t "${image_tag}" \
        -t "${latest_tag}" \
        .
    
    log_success "镜像构建完成: ${image_tag}"
}

# 推送镜像
push_image() {
    if [ -z "$DOCKER_REGISTRY" ]; then
        log_warn "未配置 DOCKER_REGISTRY，跳过推送"
        return
    fi
    
    log_info "推送镜像到 ${DOCKER_REGISTRY}..."
    
    local image_tag="${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${VERSION}"
    local latest_tag="${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest"
    
    docker push "${image_tag}"
    docker push "${latest_tag}"
    
    log_success "镜像推送完成"
}

# 本地部署
deploy_local() {
    log_info "本地部署..."
    
    # 停止旧容器
    if docker ps -q -f name="${APP_NAME}" | grep -q .; then
        log_info "停止旧容器..."
        docker stop "${APP_NAME}"
        docker rm "${APP_NAME}"
    fi
    
    # 启动新容器
    docker compose up -d
    
    # 等待服务启动
    sleep 5
    
    # 检查健康状态
    if docker ps -f name="${APP_NAME}-app" --format '{{.Status}}' | grep -q "healthy"; then
        log_success "服务启动成功"
    else
        log_warn "服务启动中，请稍后检查状态"
    fi
    
    # 显示状态
    docker compose ps
}

# 清理旧镜像
cleanup_images() {
    log_info "清理旧镜像..."
    
    # 删除悬空镜像
    docker image prune -f
    
    # 保留最近3个版本，删除更旧的镜像
    docker images "${DOCKER_IMAGE}" --format '{{.Tag}}' | \
        grep -v latest | \
        sort -rV | \
        tail -n +4 | \
        xargs -I {} docker rmi "${DOCKER_IMAGE}:{}" 2>/dev/null || true
    
    log_success "清理完成"
}

# 显示帮助
show_help() {
    cat << EOF
用法: $0 [命令] [选项]

命令:
    build       构建 Docker 镜像
    push        推送镜像到仓库
    deploy      本地部署（docker-compose）
    all         构建 + 推送 + 部署
    cleanup     清理旧镜像

选项:
    -v, --version VERSION   指定版本号
    -r, --registry URL      指定镜像仓库地址
    -h, --help              显示帮助

环境变量:
    VERSION         版本号（默认: git tag）
    DOCKER_REGISTRY 镜像仓库地址

示例:
    $0 build                    # 构建镜像
    $0 deploy                   # 本地部署
    $0 all -r registry.cn-hangzhou.aliyuncs.com/myns  # 完整流程
EOF
}

# 解析参数
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        build|push|deploy|all|cleanup)
            COMMAND=$1
            shift
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -r|--registry)
            DOCKER_REGISTRY="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 主流程
main() {
    echo ""
    echo "========================================"
    echo "  iBookFS Docker 部署"
    echo "  版本: ${VERSION}"
    echo "========================================"
    echo ""
    
    check_docker
    
    case $COMMAND in
        build)
            build_image
            ;;
        push)
            push_image
            ;;
        deploy)
            deploy_local
            ;;
        all)
            build_image
            push_image
            deploy_local
            ;;
        cleanup)
            cleanup_images
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
    
    echo ""
    log_success "完成！"
}

main
