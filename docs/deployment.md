# iBookFS 部署指南

## 目录

- [快速开始](#快速开始)
- [Shell 脚本部署](#shell-脚本部署)
- [Docker 部署](#docker-部署)
- [生产环境配置](#生产环境配置)

---

## 快速开始

### 环境要求

| 依赖 | 最低版本 | 说明 |
|------|---------|------|
| Go | 1.21+ | 后端编译 |
| Node.js | 18+ | 前端编译 |
| Docker | 20+ | 容器化部署（可选） |
| PostgreSQL | 14+ | 数据库 |

### 一键检查环境

```bash
# 检查环境依赖
make check-env

# 或直接运行脚本
./scripts/check-env.sh
```

### 一键安装依赖

```bash
# 安装缺失的依赖
make install-env

# 或直接运行脚本
sudo ./scripts/install-env.sh
```

---

## Shell 脚本部署

### 完整构建部署

```bash
# 完整流程：拉取代码 → 编译后端 → 编译前端 → 部署 → 启动服务
./scripts/build.sh
```

### 分步部署

```bash
# 仅构建后端
./scripts/build.sh -b

# 仅构建前端
./scripts/build.sh -f

# 跳过构建，仅部署
./scripts/build.sh -s

# 部署后不重启服务
./scripts/build.sh -n
```

### 服务管理

```bash
# 启动服务
sudo systemctl start ibookfs

# 停止服务
sudo systemctl stop ibookfs

# 重启服务
sudo systemctl restart ibookfs

# 查看状态
sudo systemctl status ibookfs

# 查看日志
sudo journalctl -u ibookfs -f
```

---

## 前后端独立部署

目录结构（按服务名组织）：
```
build/docker/
├── apiserver/          # API Server 后端
│   ├── Dockerfile
│   └── docker-compose.yml
└── web/                # Web 前端
    ├── Dockerfile
    ├── nginx.conf
    └── docker-compose.yml
```

### API Server 独立部署

#### Shell 部署（systemd）

```bash
# 仅构建
./scripts/deploy-backend.sh build

# 构建并部署（systemd）
./scripts/deploy-backend.sh deploy
# 或
make deploy-apiserver
```

#### Docker 部署

```bash
# 构建镜像
make docker-build-apiserver

# 运行容器（含数据库）
make docker-run-apiserver

# 停止
make docker-stop-apiserver
```

### Web 前端独立部署

#### Shell 部署（本地/Nginx）

```bash
# 仅构建
./scripts/deploy-frontend.sh build

# 部署到本地目录
./scripts/deploy-frontend.sh deploy
# 或
make deploy-web

# 部署到 Nginx
./scripts/deploy-frontend.sh nginx
```

#### Docker 部署

```bash
# 构建镜像
make docker-build-web

# 运行容器
make docker-run-web

# 停止
make docker-stop-web
```

---

## Docker 部署

### 方式一：Docker Compose（推荐）

```bash
# 1. 复制环境变量配置
cp .env.example .env
vim .env

# 2. 启动所有服务
make docker-compose-up
# 或
docker compose up -d

# 3. 查看日志
make docker-compose-logs
# 或
docker compose logs -f

# 4. 停止服务
make docker-compose-down
```

### 方式二：单独构建镜像

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run

# 推送到仓库（需配置 DOCKER_REGISTRY）
DOCKER_REGISTRY=registry.example.com/myns make docker-push
```

### 方式三：使用部署脚本

```bash
# 构建镜像
./scripts/docker-deploy.sh build

# 本地部署
./scripts/docker-deploy.sh deploy

# 完整流程（构建 + 推送 + 部署）
./scripts/docker-deploy.sh all -r registry.example.com/myns
```

---

## 生产环境配置

### 配置文件

```bash
# 复制配置模板
cp configs/config.yaml.example configs/config.yaml

# 编辑配置
vim configs/config.yaml
```

### 关键配置项

```yaml
server:
  http:
    addr: ":8080"
    timeout: 30s

data:
  database:
    - name: default
      driver: postgres
      host: ${DB_HOST}
      port: ${DB_PORT}
      user: ${DB_USER}
      password: ${DB_PASSWORD}
      database: ${DB_NAME}

auth:
  jwt:
    secret: "your-secret-key"  # 生产环境必须修改
    expired: 24h
```

### 反向代理（Nginx）

```nginx
upstream ibookfs {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://ibookfs;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSL 配置（Let's Encrypt）

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 申请证书
sudo certbot --nginx -d your-domain.com
```

---

## Make 命令速查

| 命令 | 说明 |
|------|------|
| `make help` | 显示帮助信息 |
| `make build` | 编译后端 |
| `make build-frontend` | 编译前端 |
| `make run` | 运行开发服务器 |
| `make test` | 运行测试 |
| `make lint` | 代码检查 |
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-compose-up` | 启动所有服务 |
| `make deploy` | 完整构建部署 |
| `make check-env` | 检查环境 |

---

## 故障排查

### 服务无法启动

```bash
# 查看详细日志
sudo journalctl -u ibookfs -n 50 --no-pager

# 检查端口占用
ss -tuln | grep 8080

# 检查配置文件
cat /opt/ibookfs/config.yaml
```

### Docker 容器异常

```bash
# 查看容器日志
docker logs ibookfs-app

# 进入容器调试
docker exec -it ibookfs-app sh

# 检查健康状态
docker inspect ibookfs-app --format='{{.State.Health.Status}}'
```
