# iBookFS 部署脚本说明

本目录包含 iBookFS 项目的自动部署相关脚本。

## 脚本列表

### 1. check-env.sh - 环境检查脚本
检查项目部署所需的环境依赖是否已安装。

**检查项：**
- 系统工具（curl, wget, tar, gzip）
- Git
- Go (>= 1.21)
- Node.js (>= 18) 和 npm
- 端口占用（8080, 3000）
- 磁盘空间
- 目录权限
- 网络连接

**使用方法：**
```bash
# 基本检查
./scripts/check-env.sh

# 详细输出
./scripts/check-env.sh -v

# 跳过端口检查
./scripts/check-env.sh --skip-ports
```

### 2. install-env.sh - 环境安装脚本
自动安装项目所需的环境依赖。

**支持系统：**
- Ubuntu / Debian
- CentOS / RHEL

**安装内容：**
- 基础工具（curl, wget, git, tar, gzip, gcc）
- Go 1.21.6
- Node.js 20.x LTS
- 配置 Go 代理（goproxy.cn）

**使用方法：**
```bash
# 安装所有依赖
sudo ./scripts/install-env.sh

# 跳过 Go 安装
sudo ./scripts/install-env.sh --skip-go

# 跳过 Node.js 安装
sudo ./scripts/install-env.sh --skip-nodejs
```

### 3. build.sh - 构建部署脚本
自动从 GitHub 获取代码、构建并部署到生产环境。

**功能：**
1. 从 GitHub 克隆最新代码
2. 构建 Go 后端
3. 构建前端（npm）
4. 部署到 `/opt/ibookfs`
5. 配置 systemd 服务
6. 启动服务

**使用方法：**
```bash
# 完整构建和部署
./scripts/build.sh

# 仅部署（跳过构建）
./scripts/build.sh -s

# 仅构建后端
./scripts/build.sh -b

# 仅构建前端
./scripts/build.sh -f

# 使用指定分支
./scripts/build.sh --branch dev

# 部署后不重启服务
./scripts/build.sh -n
```

## 部署流程

### 首次部署

1. **检查环境**
   ```bash
   ./scripts/check-env.sh
   ```

2. **安装依赖**（如环境检查失败）
   ```bash
   sudo ./scripts/install-env.sh
   ```

3. **构建部署**
   ```bash
   ./scripts/build.sh
   ```

### 更新部署

```bash
# 完整更新（从 GitHub 获取最新代码并构建）
./scripts/build.sh

# 仅重新部署已有代码
./scripts/build.sh -s
```

## 服务管理

部署后，使用 systemd 管理服务：

```bash
# 启动服务
sudo systemctl start ibookfs

# 停止服务
sudo systemctl stop ibookfs

# 重启服务
sudo systemctl restart ibookfs

# 查看服务状态
sudo systemctl status ibookfs

# 查看日志
sudo journalctl -u ibookfs -f

# 开机自启
sudo systemctl enable ibookfs

# 禁用自启
sudo systemctl disable ibookfs
```

## 目录结构

```
/opt/ibookfs/           # 生产部署目录
├── ibookfs             # Go 后端二进制
├── config.yaml        # 配置文件
├── index.html          # 前端入口
├── assets/            # 前端静态资源
└── db/                # 数据库文件（如使用 SQLite）
```

## 常见问题

### 1. 端口被占用
```bash
# 查看端口占用
sudo lsof -i :8080

# 停止占用端口的进程
sudo kill <PID>
```

### 2. 服务启动失败
```bash
# 查看服务状态
sudo systemctl status ibookfs

# 查看详细日志
sudo journalctl -u ibookfs -n 50 --no-pager

# 查看实时日志
sudo journalctl -u ibookfs -f
```

### 3. 权限问题
```bash
# 修复部署目录权限
sudo chown -R root:root /opt/ibookfs
sudo chmod -R 755 /opt/ibookfs
```

### 4. 网络连接问题
```bash
# 测试网络连接
curl -I https://github.com
curl -I https://goproxy.cn

# 配置 Go 代理
export GOPROXY=https://goproxy.cn,direct
```
