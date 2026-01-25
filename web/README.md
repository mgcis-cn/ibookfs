# iBookFS Web 前端

基于 Vue 3 + TypeScript 的 iBookFS 前端项目。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 下一代前端构建工具
- **Vue Router** - Vue.js 官方路由
- **Pinia** - Vue 状态管理库
- **Lucide Vue Next** - 图标库

## 项目结构

```
web/
├── src/
│   ├── api/           # API 接口层（当前为 Mock 数据）
│   │   ├── index.ts   # API 客户端
│   │   └── mockData.ts # Mock 数据
│   ├── assets/        # 静态资源
│   ├── components/    # Vue 组件
│   │   └── ui/        # UI 基础组件
│   ├── pages/         # 页面组件
│   ├── router/        # 路由配置
│   ├── stores/        # Pinia 状态管理
│   ├── styles/        # 全局样式
│   ├── types/         # TypeScript 类型定义
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── .env               # 环境变量
├── index.html         # HTML 模板
├── package.json       # 项目依赖
├── tsconfig.json      # TypeScript 配置
└── vite.config.ts     # Vite 配置
```

## 安装依赖

```bash
npm install
```

## 开发

```bash
npm run dev
```

访问 http://localhost:3000

## 构建

```bash
npm run build
```

## 预览生产构建

```bash
npm run preview
```

## 与后端对接

当前项目使用 Mock 数据，需要对接后端 API 时：

1. 修改 `.env` 文件：
   ```env
   VITE_API_BASE_URL=http://your-api-url
   VITE_USE_MOCK=false
   ```

2. 在 `src/api/index.ts` 中实现真实的 API 调用（已预留 TODO 标记位置）

## 设计系统

项目遵循 iBookFS 设计规范，详见 `/docs/design` 目录。

### 核心特性

- **响应式设计** - 支持桌面端和移动端
- **深色模式** - 支持浅色/深色主题切换
- **类型安全** - 完整的 TypeScript 类型定义
- **组件化** - 可复用的 UI 组件库
