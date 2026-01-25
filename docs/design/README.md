# iBookFS 前端设计文档

## 项目概述

iBookFS 是一个通过照片形式管理书籍的数字化工具，帮助用户将纸质书籍转变为数字产品。

### 核心功能

**第一阶段：数字化基础**
- 建书：创建书籍条目
- 图片上传：批量上传书籍照片
- 图片成册：组织照片形成数字书籍

**第二阶段：智能优化**
- 图片优化：自动增强照片质量
- 图片OCR：文字识别提取
- 图片转Markdown：结构化内容输出

---

## 设计理念

### 核心定位
**"从物理到数字的优雅转换"**

iBookFS 不仅是工具，更是连接纸质书与数字世界的桥梁。设计应体现：
- **尊重原物**：保留书籍的温度与质感
- **高效流畅**：简化复杂的数字化流程
- **专注内容**：让书籍本身成为主角

### 美学方向

**Editorial Minimalism（编辑式极简主义）**

灵感来源于高品质出版物与图书馆系统，强调：
- 清晰的信息层级
- 充足的留白空间
- 精致的排版细节
- 克制而有力的视觉元素

**关键特征：**
- 以内容为中心的布局
- 书籍封面/照片作为视觉焦点
- 类似书架/图书馆的组织方式
- 纸质质感的微妙暗示

---

## 视觉系统

### 色彩方案

**主色调：温暖中性色**
```css
:root {
  /* 主色 - 纸张与墨水 */
  --color-paper: #FAF8F3;        /* 温暖的纸白色 */
  --color-ink: #2B2B2B;          /* 深墨色 */
  --color-ink-light: #5A5A5A;    /* 浅墨色 */
  
  /* 强调色 - 书签与标记 */
  --color-accent: #C17B5C;       /* 古铜色 */
  --color-accent-light: #E8C4B3; /* 浅古铜 */
  
  /* 功能色 */
  --color-success: #7A9B76;      /* 柔和绿 */
  --color-warning: #D4A574;      /* 琥珀色 */
  --color-error: #B85C5C;        /* 砖红色 */
  
  /* 背景层次 */
  --bg-primary: #FAF8F3;
  --bg-secondary: #F2EFE8;
  --bg-elevated: #FFFFFF;
  
  /* 边框与分割 */
  --border-subtle: rgba(43, 43, 43, 0.08);
  --border-medium: rgba(43, 43, 43, 0.15);
}
```

**暗色模式：夜间阅读**
```css
[data-theme="dark"] {
  --color-paper: #1A1A1A;
  --color-ink: #E8E6E1;
  --color-ink-light: #A8A6A1;
  --color-accent: #D4956F;
  
  --bg-primary: #1A1A1A;
  --bg-secondary: #242424;
  --bg-elevated: #2E2E2E;
  
  --border-subtle: rgba(232, 230, 225, 0.08);
  --border-medium: rgba(232, 230, 225, 0.15);
}
```

### 字体系统

**原则：可读性优先，细节精致**

```css
:root {
  /* 标题字体 - 优雅衬线 */
  --font-display: 'Crimson Pro', 'Source Serif Pro', 'Noto Serif SC', serif;
  
  /* 正文字体 - 清晰无衬线 */
  --font-body: 'Inter', 'SF Pro Text', 'PingFang SC', sans-serif;
  
  /* 等宽字体 - 技术信息 */
  --font-mono: 'JetBrains Mono', 'SF Mono', 'Menlo', monospace;
  
  /* 字重 */
  --font-weight-light: 300;
  --font-weight-regular: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
  
  /* 字号比例 (1.25 - Major Third) */
  --text-xs: 0.64rem;    /* 10.24px */
  --text-sm: 0.8rem;     /* 12.8px */
  --text-base: 1rem;     /* 16px */
  --text-lg: 1.25rem;    /* 20px */
  --text-xl: 1.563rem;   /* 25px */
  --text-2xl: 1.953rem;  /* 31.25px */
  --text-3xl: 2.441rem;  /* 39px */
  --text-4xl: 3.052rem;  /* 48.83px */
}
```

### 间距系统

**8px 基准网格**

```css
:root {
  --space-1: 0.25rem;   /* 4px */
  --space-2: 0.5rem;    /* 8px */
  --space-3: 0.75rem;   /* 12px */
  --space-4: 1rem;      /* 16px */
  --space-5: 1.5rem;    /* 24px */
  --space-6: 2rem;      /* 32px */
  --space-8: 3rem;      /* 48px */
  --space-10: 4rem;     /* 64px */
  --space-12: 6rem;     /* 96px */
  --space-16: 8rem;     /* 128px */
  
  /* 容器宽度 */
  --container-sm: 640px;
  --container-md: 768px;
  --container-lg: 1024px;
  --container-xl: 1280px;
}
```

### 圆角与阴影

```css
:root {
  /* 圆角 - 柔和但不过度 */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;
  
  /* 阴影 - 微妙的深度 */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.04);
  --shadow-md: 0 2px 8px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 4px 16px rgba(0, 0, 0, 0.12);
  --shadow-xl: 0 8px 32px rgba(0, 0, 0, 0.16);
  
  /* 书籍卡片特殊阴影 */
  --shadow-book: 
    0 1px 3px rgba(0, 0, 0, 0.08),
    0 4px 12px rgba(0, 0, 0, 0.06),
    0 0 0 1px rgba(0, 0, 0, 0.04);
}
```

---

## 动效系统

### 过渡时长

```css
:root {
  --duration-fast: 150ms;
  --duration-base: 250ms;
  --duration-slow: 400ms;
  --duration-slower: 600ms;
  
  --ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-out: cubic-bezier(0, 0, 0.2, 1);
  --ease-in: cubic-bezier(0.4, 0, 1, 1);
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);
}
```

### 动效原则

1. **微妙而有意义**：每个动效都应有明确目的
2. **性能优先**：优先使用 transform 和 opacity
3. **渐进增强**：支持 `prefers-reduced-motion`
4. **一致性**：相似操作使用相似动效

**关键动效场景：**
- 页面加载：渐入 + 轻微上移
- 卡片悬停：轻微上浮 + 阴影增强
- 图片上传：进度指示 + 成功反馈
- 页面切换：淡入淡出 + 位移
- 模态框：背景模糊 + 缩放渐入

---

## 组件设计规范

### 按钮系统

**主要按钮（Primary）**
- 用于主要操作：上传照片、创建书籍、开始处理
- 背景：`--color-accent`
- 悬停：加深 10%，轻微上浮

**次要按钮（Secondary）**
- 用于次要操作：取消、返回
- 边框：`--border-medium`
- 悬停：背景 `--bg-secondary`

**文本按钮（Text）**
- 用于低优先级操作
- 颜色：`--color-ink-light`
- 悬停：颜色 `--color-ink`

### 卡片系统

**书籍卡片**
```
┌─────────────────┐
│                 │
│   [封面图片]    │  ← 3:4 比例
│                 │
├─────────────────┤
│ 书名            │  ← 标题字体
│ 作者 · 页数     │  ← 元信息
│ [进度条]        │  ← 数字化进度
└─────────────────┘
```

**照片卡片**
```
┌─────────────────┐
│                 │
│  [照片缩略图]   │  ← 保持原始比例
│                 │
│  页码标记       │  ← 右下角
└─────────────────┘
```

### 表单元素

**输入框**
- 边框：`--border-medium`
- 聚焦：边框色变为 `--color-accent`，添加微妙外发光
- 占位符：`--color-ink-light`，斜体

**文件上传区域**
- 虚线边框，拖拽时高亮
- 大图标 + 清晰提示文字
- 支持拖拽、点击、粘贴多种方式

---

## 布局模式

### 主导航

**侧边栏导航（Desktop）**
```
┌────────┬──────────────────┐
│        │                  │
│  LOGO  │                  │
│        │                  │
│  书架  │   主内容区域     │
│  上传  │                  │
│  设置  │                  │
│        │                  │
│ [用户] │                  │
└────────┴──────────────────┘
```

**底部导航（Mobile）**
```
┌──────────────────────────┐
│                          │
│     主内容区域           │
│                          │
├──────────────────────────┤
│  书架  上传  我的  设置  │
└──────────────────────────┘
```

### 页面布局

**书架视图**
- 网格布局：响应式列数（2-6列）
- 卡片间距：`--space-5`
- 支持列表/网格切换

**书籍详情**
- 左侧：封面 + 元信息
- 右侧：照片网格
- 顶部：操作栏（编辑、导出、删除）

**上传流程**
- 居中布局
- 步骤指示器
- 大面积拖拽区域

---

## 响应式设计

### 断点

```css
/* Mobile First */
--breakpoint-sm: 640px;   /* 手机横屏 */
--breakpoint-md: 768px;   /* 平板竖屏 */
--breakpoint-lg: 1024px;  /* 平板横屏/小笔记本 */
--breakpoint-xl: 1280px;  /* 桌面 */
--breakpoint-2xl: 1536px; /* 大屏 */
```

### 适配策略

**移动端（< 768px）**
- 单列布局
- 底部导航
- 全屏模态框
- 手势操作优化

**平板（768px - 1024px）**
- 2-3列网格
- 侧边栏可收起
- 混合导航模式

**桌面（> 1024px）**
- 多列网格（3-6列）
- 固定侧边栏
- 悬停交互增强
- 键盘快捷键

---

## 交互模式

### 第一阶段核心交互

**1. 创建书籍**
```
点击"新建书籍" 
  → 模态框弹出
  → 填写基本信息（书名、作者、ISBN可选）
  → 确认创建
  → 跳转到上传页面
```

**2. 上传照片**
```
拖拽/选择文件
  → 显示预览网格
  → 自动排序（可手动调整）
  → 批量上传（进度显示）
  → 完成后跳转到书籍详情
```

**3. 照片管理**
```
网格视图
  → 悬停显示操作按钮
  → 点击放大查看
  → 拖拽重新排序
  → 批量选择操作（删除、导出）
```

### 第二阶段交互

**图片优化**
- 一键批量优化
- 实时预览对比
- 可调整参数

**OCR识别**
- 选择页面范围
- 后台处理 + 进度提示
- 结果预览与编辑

**导出Markdown**
- 选择导出范围
- 模板选择
- 下载/复制

---

## 状态与反馈

### 加载状态

**骨架屏**
- 用于初始加载
- 模拟真实内容结构
- 微妙的脉动动画

**进度条**
- 用于上传、处理等长时操作
- 显示百分比和预计时间
- 可取消操作

**加载指示器**
- 用于短时等待（< 3秒）
- 简洁的旋转动画

### 反馈提示

**Toast通知**
- 右上角滑入
- 3秒自动消失
- 成功/警告/错误不同颜色

**空状态**
- 友好的插图
- 清晰的引导文案
- 明确的操作按钮

**错误处理**
- 明确的错误信息
- 建议的解决方案
- 重试选项

---

## 可访问性

### 基本要求

- **色彩对比度**：至少 4.5:1（正文），3:1（大文本）
- **键盘导航**：所有功能可通过键盘操作
- **焦点指示**：清晰的焦点样式
- **语义化HTML**：正确使用标签和ARIA属性
- **屏幕阅读器**：提供有意义的alt文本和标签

### 快捷键

```
全局：
  Cmd/Ctrl + K    → 快速搜索
  Cmd/Ctrl + U    → 上传照片
  Cmd/Ctrl + N    → 新建书籍
  Esc             → 关闭模态框

书架：
  ←/→             → 切换书籍
  Enter           → 打开书籍
  Delete          → 删除选中项

照片视图：
  Space           → 全屏查看
  ←/→             → 上一张/下一张
  Cmd/Ctrl + A    → 全选
```

---

## 技术实现建议

### 技术栈推荐

**前端框架**
- React 18+ / Vue 3+
- TypeScript

**样式方案**
- CSS Modules / Tailwind CSS
- CSS Variables 用于主题

**动画库**
- Framer Motion (React)
- GSAP (复杂动画)

**图片处理**
- react-dropzone / vue-dropzone
- react-image-crop
- 虚拟滚动（大量图片）

**状态管理**
- Zustand / Pinia
- React Query / VueUse

### 性能优化

**图片优化**
- 懒加载（Intersection Observer）
- 响应式图片（srcset）
- WebP格式 + 降级
- 缩略图 + 原图分离

**代码分割**
- 路由级别代码分割
- 组件懒加载
- 动态导入

**缓存策略**
- Service Worker
- IndexedDB 存储元数据
- 图片本地缓存

---

## 设计交付物

### 必需文件

1. **设计系统文档**（本文档）
2. **组件库 Storybook**
3. **原型图**（Figma/Sketch）
4. **图标资源**
5. **示例页面**

### 开发协作

**设计稿标注**
- 使用 8px 网格对齐
- 标注关键间距和尺寸
- 提供多种状态示例

**组件命名规范**
```
Button
  - ButtonPrimary
  - ButtonSecondary
  - ButtonText

Card
  - CardBook
  - CardPhoto

Input
  - InputText
  - InputFile
```

---

## 未来扩展

### 第二阶段设计考虑

**智能功能UI**
- 优化前后对比滑块
- OCR结果编辑器（类似文档标注）
- Markdown预览与编辑

**高级功能**
- 书籍分类与标签系统
- 搜索与过滤
- 批量操作工具栏
- 导出模板编辑器

**社交功能（可选）**
- 分享书籍
- 协作编辑
- 评论与标注

---

## 参考资源

### 设计灵感

- **Readwise Reader**：优秀的阅读应用设计
- **Notion**：清晰的信息架构
- **Linear**：精致的交互细节
- **Dropbox Paper**：文档编辑体验

### 字体资源

- [Google Fonts](https://fonts.google.com/)
- [Adobe Fonts](https://fonts.adobe.com/)
- [Font Squirrel](https://www.fontsquirrel.com/)

### 图标库

- [Lucide Icons](https://lucide.dev/)
- [Heroicons](https://heroicons.com/)
- [Phosphor Icons](https://phosphoricons.com/)

---

## 版本历史

- **v1.0** (2026-01-24)：初始版本，覆盖第一阶段核心功能设计