# iBookFS 设计资源索引

本目录包含 iBookFS 项目的完整前端设计文档。

---

## 📚 文档导航

### 核心设计系统

| 文档 | 描述 | 关键内容 |
|------|------|----------|
| [README.md](./README.md) | 设计系统总览 | 设计理念、视觉系统、色彩字体规范 |
| [design-tokens.css](./design-tokens.css) | CSS 设计变量 | 可直接使用的 CSS 变量，支持主题切换 |
| [components.md](./components.md) | 组件库文档 | 按钮、卡片、表单等基础组件实现 |
| [layouts.md](./layouts.md) | 页面布局规范 | 书架、详情、上传、设置页面布局 |

### 营销和认证

| 文档 | 描述 | 关键内容 |
|------|------|----------|
| [landing-page.md](./landing-page.md) | 网站主页设计 | 英雄区、功能展示、CTA 转化设计 |
| [auth-pages.md](./auth-pages.md) | 登录注册页面 | 邮箱验证码认证、表单交互逻辑 |

---

## 🎨 设计理念

**Editorial Minimalism（编辑式极简主义）**

灵感来源于高品质出版物与图书馆系统，强调：
- 清晰的信息层级
- 充足的留白空间
- 精致的排版细节
- 克制而有力的视觉元素

---

## 🚀 快速开始

### 1. 导入设计变量

```html
<link rel="stylesheet" href="/docs/design/design-tokens.css">
```

### 2. 使用 CSS 变量

```css
.my-component {
  color: var(--text-primary);
  background: var(--bg-elevated);
  padding: var(--space-4);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}
```

### 3. 参考组件实现

查看 [components.md](./components.md) 获取完整的组件 HTML/CSS 代码。

---

## 📐 核心规范速查

### 色彩系统

```css
/* 主色 */
--color-paper: #FAF8F3;      /* 纸白色 */
--color-ink: #2B2B2B;        /* 墨色 */
--color-accent: #C17B5C;     /* 古铜色 */

/* 功能色 */
--color-success: #7A9B76;
--color-warning: #D4A574;
--color-error: #B85C5C;
```

### 字体系统

```css
/* 字体家族 */
--font-display: 'Crimson Pro', serif;  /* 标题 */
--font-body: 'Inter', sans-serif;      /* 正文 */
--font-mono: 'JetBrains Mono', monospace;

/* 字号（Major Third Scale 1.25） */
--text-xs: 0.64rem;    /* 10.24px */
--text-sm: 0.8rem;     /* 12.8px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.25rem;    /* 20px */
--text-xl: 1.563rem;   /* 25px */
--text-2xl: 1.953rem;  /* 31.25px */
--text-3xl: 2.441rem;  /* 39px */
--text-4xl: 3.052rem;  /* 48.83px */
```

### 间距系统（8px 网格）

```css
--space-1: 0.25rem;   /* 4px */
--space-2: 0.5rem;    /* 8px */
--space-3: 0.75rem;   /* 12px */
--space-4: 1rem;      /* 16px */
--space-5: 1.5rem;    /* 24px */
--space-6: 2rem;      /* 32px */
--space-8: 3rem;      /* 48px */
--space-10: 4rem;     /* 64px */
```

### 圆角和阴影

```css
/* 圆角 */
--radius-sm: 4px;
--radius-md: 8px;
--radius-lg: 12px;
--radius-xl: 16px;

/* 阴影 */
--shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.08);
--shadow-md: 0 2px 8px rgba(0, 0, 0, 0.08);
--shadow-lg: 0 4px 16px rgba(0, 0, 0, 0.12);
--shadow-book: /* 书籍卡片专用阴影 */
  0 1px 3px rgba(0, 0, 0, 0.08),
  0 4px 12px rgba(0, 0, 0, 0.06),
  0 0 0 1px rgba(0, 0, 0, 0.04);
```

---

## 🎯 页面类型

### 应用内页面

1. **书架页面** - 展示所有书籍，支持搜索筛选
2. **书籍详情** - 查看书籍信息和照片
3. **上传页面** - 三步上传流程
4. **设置页面** - 用户偏好设置

详见 [layouts.md](./layouts.md)

### 营销页面

5. **网站主页** - 产品介绍和功能展示
   - 英雄区（Hero Section）
   - 功能展示（Features）
   - 工作流程（How It Works）
   - CTA 转化区域

详见 [landing-page.md](./landing-page.md)

### 认证页面

6. **登录页面** - 邮箱 + 验证码登录
7. **注册页面** - 姓名 + 邮箱 + 验证码注册

详见 [auth-pages.md](./auth-pages.md)

---

## 🔧 组件清单

### 基础组件

- ✅ Button（Primary / Secondary / Text）
- ✅ Card（Book / Photo）
- ✅ Input（Text / Email / File Upload）
- ✅ Modal
- ✅ Toast
- ✅ Skeleton
- ✅ Spinner
- ✅ Empty State

### 导航组件

- ✅ Sidebar（Desktop）
- ✅ Bottom Navigation（Mobile）
- ✅ Breadcrumb

### 表单组件

- ✅ Text Input
- ✅ Email Input
- ✅ Code Input（验证码）
- ✅ Checkbox
- ✅ File Upload Area

详见 [components.md](./components.md)

---

## 📱 响应式断点

```css
--breakpoint-sm: 640px;   /* 手机横屏 */
--breakpoint-md: 768px;   /* 平板竖屏 */
--breakpoint-lg: 1024px;  /* 平板横屏/小笔记本 */
--breakpoint-xl: 1280px;  /* 桌面 */
--breakpoint-2xl: 1536px; /* 大屏 */
```

### 适配策略

- **移动端（< 768px）**：单列布局、底部导航、全屏模态框
- **平板（768px - 1024px）**：2-3列网格、可收起侧边栏
- **桌面（> 1024px）**：多列网格、固定侧边栏、悬停交互

---

## ♿ 可访问性

### 基本要求

- ✅ 色彩对比度：至少 4.5:1（正文），3:1（大文本）
- ✅ 键盘导航：所有功能可通过键盘操作
- ✅ 焦点指示：清晰的焦点样式
- ✅ 语义化 HTML：正确使用标签和 ARIA 属性
- ✅ 屏幕阅读器：有意义的 alt 文本和标签

### 全局快捷键

```
Cmd/Ctrl + K    → 快速搜索
Cmd/Ctrl + U    → 上传照片
Cmd/Ctrl + N    → 新建书籍
Esc             → 关闭模态框
```

---

## 🎬 动效原则

### 时长

```css
--duration-fast: 150ms;    /* 快速交互 */
--duration-base: 250ms;    /* 标准过渡 */
--duration-slow: 400ms;    /* 复杂动画 */
```

### 缓动函数

```css
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
--ease-out: cubic-bezier(0, 0, 0.2, 1);
--ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);
```

### 关键动效场景

- 页面加载：渐入 + 轻微上移
- 卡片悬停：轻微上浮 + 阴影增强
- 模态框：背景模糊 + 缩放渐入
- Toast：右侧滑入

---

## 🔐 安全考虑（认证页面）

### 验证码机制

- 6 位随机数字
- 5 分钟有效期
- 单次使用
- 60 秒发送间隔

### 速率限制

- 发送验证码：60 秒/次
- 登录尝试：5 次/小时
- 注册尝试：3 次/小时

详见 [auth-pages.md](./auth-pages.md)

---

## 📊 性能优化

### 图片优化

- WebP 格式 + PNG/JPG 降级
- 懒加载（Intersection Observer）
- 响应式图片（srcset）
- 缩略图 + 原图分离

### 代码优化

- 路由级别代码分割
- 组件懒加载
- CSS 变量减少重复代码

### 缓存策略

- Service Worker
- IndexedDB 存储元数据
- 图片本地缓存

---

## 🎨 设计资源

### 推荐字体

- **Crimson Pro** - Google Fonts（标题）
- **Inter** - Google Fonts（正文）
- **JetBrains Mono** - 等宽字体

### 推荐图标库

- [Lucide Icons](https://lucide.dev/)
- [Heroicons](https://heroicons.com/)
- [Phosphor Icons](https://phosphoricons.com/)

### 设计工具

- Figma - 原型设计
- Storybook - 组件开发
- Chrome DevTools - 调试

---

## 📝 版本历史

- **v1.1** (2026-01-26)
  - ✨ 新增网站主页设计（landing-page.md）
  - ✨ 新增登录注册页面设计（auth-pages.md）
  - 📝 更新文档索引和导航

- **v1.0** (2026-01-24)
  - 🎉 初始版本发布
  - 📚 核心设计系统文档
  - 🎨 组件库和布局规范

---

## 🤝 贡献指南

### 修改设计规范

1. 在对应的 `.md` 文件中修改
2. 更新 `design-tokens.css` 中的变量
3. 更新本索引文件的版本历史

### 添加新组件

1. 在 `components.md` 中添加组件文档
2. 提供完整的 HTML/CSS 实现
3. 包含响应式适配方案
4. 添加可访问性说明

### 添加新页面

1. 在 `layouts.md` 中添加页面布局
2. 或创建独立的页面设计文档
3. 更新本索引文件的导航

---

## 📧 联系方式

如有设计相关问题，请联系设计团队。

---

**最后更新：** 2026-01-26  
**维护者：** iBookFS 设计团队