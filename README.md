# iBookFS

iBookFS 是一个书籍文件管理服务，采用 Go + Vue.js 全栈开发。

## 特性

- 📚 书籍元数据管理（标题、作者、ISBN、出版社等）
- 📄 书籍页面上传与状态追踪
- 🔍 按状态筛选和搜索
- 🎨 现代化 Vue.js 前端界面
- 🗄️ 多数据库支持（SQLite / MySQL）
- ⚙️ 灵活配置（YAML 文件 + 环境变量）

## 项目结构

```
ibookfs/
├── cmd/server/              # 应用入口
├── internal/
│   ├── config/              # 配置管理（支持 YAML + 环境变量）
│   ├── database/            # 数据库连接（atomic.Pointer 模式）
│   │   └── sources/         # 数据源抽象
│   │       ├── sources.go   # Source 接口 + Registry
│   │       ├── sqlite/      # SQLite 实现
│   │       └── mysql/       # MySQL 实现
│   ├── handler/             # HTTP 处理器
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由定义
│   └── service/             # 业务逻辑
├── web/                     # Vue.js 前端
│   ├── src/
│   │   ├── api/             # API 客户端
│   │   ├── components/      # 组件
│   │   ├── pages/           # 页面
│   │   ├── stores/          # Pinia 状态管理
│   │   └── types/           # TypeScript 类型
│   └── ...
├── config.yaml.example      # 配置示例
├── go.mod
├── Makefile
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+ (前端开发)
- MySQL 5.7+ (可选，默认使用 SQLite)

### 后端启动

```bash
# 安装依赖
go mod tidy

# 运行服务（默认使用 SQLite）
go run ./cmd/server

# 或使用 Makefile
make run
```

### 前端启动

```bash
cd web
npm install
npm run dev
```

### 构建

```bash
# 后端
make build

# 前端
cd web && npm run build
```

## API 接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /api/v1/books | 获取书籍列表（支持分页、状态筛选、搜索） |
| GET | /api/v1/books/:id | 获取书籍详情 |
| POST | /api/v1/books | 创建书籍 |
| PUT | /api/v1/books/:id | 更新书籍 |
| DELETE | /api/v1/books/:id | 删除书籍 |

### 查询参数

```
GET /api/v1/books?page=1&pageSize=10&status=completed&search=关键词
```

## 配置

### 方式一：YAML 文件

复制 `config.yaml.example` 为 `config.yaml`：

```yaml
server:
  address: ":8080"
  allowedOrigins:
    - "http://localhost:3000"
    - "http://localhost:5173"

database:
  kind: sqlite          # sqlite 或 mysql
  path: "ibookfs.db"    # SQLite 文件路径
  
  # MySQL 配置
  # kind: mysql
  # host: "localhost"
  # port: "3306"
  # user: "root"
  # password: "password"
  # name: "ibookfs"
```

### 方式二：环境变量

环境变量会覆盖 YAML 配置：

| 变量 | 默认值 | 描述 |
|------|--------|------|
| SERVER_ADDR | :8080 | 服务地址 |
| ALLOWED_ORIGINS | localhost:3000,localhost:5173 | CORS 允许的源 |
| DB_KIND | sqlite | 数据库类型 |
| DB_PATH | ibookfs.db | SQLite 文件路径 |
| DB_HOST | localhost | MySQL 主机 |
| DB_PORT | 3306 | MySQL 端口 |
| DB_USER | - | MySQL 用户 |
| DB_PASSWORD | - | MySQL 密码 |
| DB_NAME | ibookfs | MySQL 数据库名 |

## 架构设计

### 数据库抽象

采用 Registry 模式实现多数据库支持：

```go
// 使用 SQLite
source := sqlite.New("ibookfs.db")
database.OpenSource(source)

// 使用 MySQL
source := mysql.New(host, port, user, password, dbname)
database.OpenSource(source)

// 获取默认数据库
db := database.Default()
```

### 配置优先级

1. 环境变量（最高）
2. config.yaml 文件
3. 默认值（最低）

## License

MIT
