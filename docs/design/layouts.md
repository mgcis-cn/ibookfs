# iBookFS 页面布局规范

本文档定义项目中主要页面的布局结构和交互流程。

---

## 整体布局架构

### Desktop 布局

```
┌─────────────────────────────────────────────┐
│  Sidebar  │         Main Content            │
│           │                                  │
│   Logo    │  ┌─────────────────────────┐   │
│           │  │     Page Header         │   │
│  ┌─────┐  │  └─────────────────────────┘   │
│  │书架 │  │                                  │
│  └─────┘  │  ┌─────────────────────────┐   │
│  ┌─────┐  │  │                         │   │
│  │上传 │  │  │                         │   │
│  └─────┘  │  │    Page Content         │   │
│  ┌─────┐  │  │                         │   │
│  │设置 │  │  │                         │   │
│  └─────┘  │  └─────────────────────────┘   │
│           │                                  │
│  [User]   │                                  │
└─────────────────────────────────────────────┘
```

### Mobile 布局

```
┌─────────────────────┐
│   Mobile Header     │
├─────────────────────┤
│                     │
│                     │
│   Page Content      │
│                     │
│                     │
├─────────────────────┤
│  Bottom Navigation  │
└─────────────────────┘
```

---

## 页面详细设计

## 1. 书架页面 (Library)

### 功能概述
- 展示所有已创建的书籍
- 支持网格/列表视图切换
- 搜索和筛选功能
- 快速创建新书籍

### 布局结构

```html
<div class="page-library">
  <!-- 页面头部 -->
  <header class="page-header">
    <div class="header-left">
      <h1 class="page-title">我的书架</h1>
      <span class="book-count">12 本书籍</span>
    </div>
    <div class="header-right">
      <div class="search-box">
        <input type="text" placeholder="搜索书籍..." />
      </div>
      <div class="view-toggle">
        <button class="toggle-btn active">网格</button>
        <button class="toggle-btn">列表</button>
      </div>
      <button class="button-primary">
        <svg>+</svg>
        新建书籍
      </button>
    </div>
  </header>

  <!-- 筛选栏 -->
  <div class="filter-bar">
    <button class="filter-chip active">全部</button>
    <button class="filter-chip">进行中</button>
    <button class="filter-chip">已完成</button>
    <button class="filter-chip">最近添加</button>
  </div>

  <!-- 书籍网格 -->
  <div class="books-grid">
    <!-- 书籍卡片 -->
    <div class="card-book">...</div>
    <div class="card-book">...</div>
    <div class="card-book">...</div>
  </div>
</div>
```

### 样式定义

```css
.page-library {
  padding: var(--space-8);
  max-width: var(--content-max-width);
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
  gap: var(--space-4);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
}

.page-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
}

.book-count {
  font-size: var(--text-base);
  color: var(--text-secondary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.search-box {
  position: relative;
}

.search-box input {
  width: 280px;
  padding: var(--space-3) var(--space-4) var(--space-3) var(--space-10);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  background: var(--bg-elevated);
}

.search-box::before {
  content: '🔍';
  position: absolute;
  left: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
}

.view-toggle {
  display: flex;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: var(--space-1);
}

.toggle-btn {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-out);
}

.toggle-btn.active {
  background: var(--bg-elevated);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
}

.filter-bar {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
  overflow-x: auto;
}

.filter-chip {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  white-space: nowrap;
  transition: all var(--duration-fast) var(--ease-out);
}

.filter-chip:hover {
  background: var(--bg-tertiary);
}

.filter-chip.active {
  background: var(--color-accent);
  color: white;
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: var(--space-5);
}

/* 响应式 */
@media (max-width: 768px) {
  .page-library {
    padding: var(--space-4);
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .header-left {
    flex-direction: column;
    gap: var(--space-2);
  }

  .header-right {
    flex-direction: column;
  }

  .search-box input {
    width: 100%;
  }

  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: var(--space-3);
  }
}
```

### 交互状态

**空状态：**
```html
<div class="empty-state">
  <svg class="empty-icon">📚</svg>
  <h3>还没有书籍</h3>
  <p>开始创建你的第一本数字书籍</p>
  <button class="button-primary">创建书籍</button>
</div>
```

**加载状态：**
```html
<div class="books-grid">
  <div class="skeleton-book-card"></div>
  <div class="skeleton-book-card"></div>
  <div class="skeleton-book-card"></div>
</div>
```

---

## 2. 书籍详情页 (Book Detail)

### 功能概述
- 查看书籍基本信息
- 浏览所有照片页面
- 管理照片（排序、删除）
- 导出和分享

### 布局结构

```html
<div class="page-book-detail">
  <!-- 返回导航 -->
  <nav class="breadcrumb">
    <a href="#">书架</a>
    <span>/</span>
    <span>深入理解计算机系统</span>
  </nav>

  <!-- 书籍信息区 -->
  <section class="book-info">
    <div class="book-cover-large">
      <img src="cover.jpg" alt="书籍封面" />
    </div>
    <div class="book-meta">
      <h1 class="book-title">深入理解计算机系统</h1>
      <p class="book-author">Randal E. Bryant · David R. O'Hallaron</p>
      <div class="book-stats">
        <div class="stat-item">
          <span class="stat-label">总页数</span>
          <span class="stat-value">120</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">已上传</span>
          <span class="stat-value">120</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">创建时间</span>
          <span class="stat-value">2024-01-15</span>
        </div>
      </div>
      <div class="book-actions">
        <button class="button-primary">上传照片</button>
        <button class="button-secondary">编辑信息</button>
        <button class="button-secondary">导出</button>
      </div>
    </div>
  </section>

  <!-- 照片网格区 -->
  <section class="photos-section">
    <div class="section-header">
      <h2>书籍页面</h2>
      <div class="section-actions">
        <button class="button-text">全选</button>
        <button class="button-text">排序</button>
      </div>
    </div>
    <div class="photos-grid">
      <div class="card-photo">...</div>
      <div class="card-photo">...</div>
      <div class="card-photo">...</div>
    </div>
  </section>
</div>
```

### 样式定义

```css
.page-book-detail {
  padding: var(--space-8);
  max-width: var(--content-max-width);
  margin: 0 auto;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.breadcrumb a {
  color: var(--text-secondary);
  transition: color var(--duration-fast) var(--ease-out);
}

.breadcrumb a:hover {
  color: var(--color-accent);
}

.book-info {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: var(--space-8);
  margin-bottom: var(--space-10);
  padding: var(--space-8);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
}

.book-cover-large {
  aspect-ratio: 3 / 4;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-book);
}

.book-cover-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.book-meta {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.book-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--leading-tight);
}

.book-author {
  font-size: var(--text-lg);
  color: var(--text-secondary);
}

.book-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  padding: var(--space-5) 0;
  border-top: 1px solid var(--border-subtle);
  border-bottom: 1px solid var(--border-subtle);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.stat-value {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.book-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: auto;
}

.photos-section {
  margin-top: var(--space-8);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.section-header h2 {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
}

.section-actions {
  display: flex;
  gap: var(--space-2);
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-4);
}

/* 响应式 */
@media (max-width: 768px) {
  .book-info {
    grid-template-columns: 1fr;
    gap: var(--space-5);
  }

  .book-stats {
    grid-template-columns: repeat(3, 1fr);
  }

  .book-actions {
    flex-direction: column;
  }

  .photos-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: var(--space-3);
  }
}
```

---

## 3. 上传页面 (Upload)

### 功能概述
- 批量上传照片
- 预览和排序
- 设置页码
- 提交处理

### 布局结构

```html
<div class="page-upload">
  <!-- 步骤指示器 -->
  <div class="steps-indicator">
    <div class="step active">
      <div class="step-number">1</div>
      <span class="step-label">选择照片</span>
    </div>
    <div class="step-divider"></div>
    <div class="step">
      <div class="step-number">2</div>
      <span class="step-label">整理排序</span>
    </div>
    <div class="step-divider"></div>
    <div class="step">
      <div class="step-number">3</div>
      <span class="step-label">确认上传</span>
    </div>
  </div>

  <!-- 步骤 1: 选择照片 -->
  <div class="upload-step step-1">
    <div class="upload-area">
      <input type="file" id="file-input" multiple accept="image/*" />
      <label for="file-input" class="upload-label">
        <svg class="upload-icon">📸</svg>
        <span class="upload-text">拖拽照片到此处，或点击选择</span>
        <span class="upload-hint">支持 JPG、PNG 格式，最大 10MB</span>
      </label>
    </div>
  </div>

  <!-- 步骤 2: 整理排序 -->
  <div class="upload-step step-2" style="display: none;">
    <div class="preview-toolbar">
      <span class="photo-count">已选择 12 张照片</span>
      <div class="toolbar-actions">
        <button class="button-text">自动排序</button>
        <button class="button-text">清空</button>
      </div>
    </div>
    <div class="preview-grid sortable">
      <div class="preview-item" draggable="true">
        <img src="photo1.jpg" alt="" />
        <div class="preview-overlay">
          <span class="page-number">1</span>
          <button class="remove-btn">×</button>
        </div>
      </div>
      <!-- 更多预览项 -->
    </div>
    <div class="upload-actions">
      <button class="button-secondary">上一步</button>
      <button class="button-primary">下一步</button>
    </div>
  </div>

  <!-- 步骤 3: 确认上传 -->
  <div class="upload-step step-3" style="display: none;">
    <div class="upload-summary">
      <h2>确认上传</h2>
      <div class="summary-info">
        <p>书籍：<strong>深入理解计算机系统</strong></p>
        <p>照片数量：<strong>12 张</strong></p>
        <p>起始页码：<strong>第 1 页</strong></p>
      </div>
    </div>
    <div class="upload-progress" style="display: none;">
      <div class="progress-bar-container">
        <div class="progress-bar" style="width: 45%"></div>
      </div>
      <p class="progress-text">上传中... 45%</p>
    </div>
    <div class="upload-actions">
      <button class="button-secondary">上一步</button>
      <button class="button-primary">开始上传</button>
    </div>
  </div>
</div>
```

### 样式定义

```css
.page-upload {
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-8);
}

.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--space-10);
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
}

.step-number {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  transition: all var(--duration-base) var(--ease-out);
}

.step.active .step-number {
  background: var(--color-accent);
  color: white;
}

.step.completed .step-number {
  background: var(--color-success);
  color: white;
}

.step-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.step.active .step-label {
  color: var(--text-primary);
  font-weight: var(--font-weight-medium);
}

.step-divider {
  width: 80px;
  height: 2px;
  background: var(--border-medium);
  margin: 0 var(--space-4);
}

.upload-step {
  animation: fadeIn var(--duration-base) var(--ease-out);
}

.preview-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.photo-count {
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
}

.toolbar-actions {
  display: flex;
  gap: var(--space-2);
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.preview-item {
  position: relative;
  aspect-ratio: 3 / 4;
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: move;
  transition: all var(--duration-base) var(--ease-out);
}

.preview-item:hover {
  transform: scale(1.05);
  box-shadow: var(--shadow-lg);
  z-index: 1;
}

.preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.6), transparent);
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: var(--space-2);
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.preview-item:hover .preview-overlay {
  opacity: 1;
}

.page-number {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-semibold);
  color: white;
}

.remove-btn {
  width: 24px;
  height: 24px;
  font-size: var(--text-lg);
  color: white;
  background: rgba(184, 92, 92, 0.8);
  border-radius: var(--radius-full);
}

.upload-actions {
  display: flex;
  gap: var(--space-3);
  justify-content: center;
  margin-top: var(--space-8);
}

.upload-summary {
  text-align: center;
  padding: var(--space-8);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-6);
}

.upload-summary h2 {
  font-size: var(--text-3xl);
  margin-bottom: var(--space-5);
}

.summary-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  font-size: var(--text-lg);
}

.upload-progress {
  text-align: center;
  padding: var(--space-6);
}

.progress-bar-container {
  height: 8px;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: var(--space-3);
}

.progress-bar {
  height: 100%;
  background: var(--color-accent);
  border-radius: var(--radius-full);
  transition: width var(--duration-base) var(--ease-out);
}

.progress-text {
  font-size: var(--text-base);
  color: var(--text-secondary);
}
```

---

## 4. 设置页面 (Settings)

### 布局结构

```html
<div class="page-settings">
  <h1 class="page-title">设置</h1>

  <div class="settings-sections">
    <!-- 外观设置 -->
    <section class="settings-section">
      <h2 class="section-title">外观</h2>
      <div class="setting-item">
        <div class="setting-info">
          <label>主题模式</label>
          <p class="setting-description">选择浅色或深色主题</p>
        </div>
        <select class="setting-control">
          <option>跟随系统</option>
          <option>浅色</option>
          <option>深色</option>
        </select>
      </div>
    </section>

    <!-- 上传设置 -->
    <section class="settings-section">
      <h2 class="section-title">上传</h2>
      <div class="setting-item">
        <div class="setting-info">
          <label>图片质量</label>
          <p class="setting-description">上传时的压缩质量</p>
        </div>
        <select class="setting-control">
          <option>原图</option>
          <option>高质量</option>
          <option>标准</option>
        </select>
      </div>
    </section>

    <!-- 关于 -->
    <section class="settings-section">
      <h2 class="section-title">关于</h2>
      <div class="setting-item">
        <div class="setting-info">
          <label>版本</label>
          <p class="setting-description">iBookFS v1.0.0</p>
        </div>
      </div>
    </section>
  </div>
</div>
```

```css
.page-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-8);
}

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.settings-section {
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
}

.section-title {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--border-subtle);
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) 0;
}

.setting-item:not(:last-child) {
  border-bottom: 1px solid var(--border-subtle);
}

.setting-info label {
  display: block;
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--space-1);
}

.setting-description {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.setting-control {
  padding: var(--space-2) var(--space-4);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
}
```

---

## 响应式布局策略

### 断点使用

```css
/* 移动端优先 */
.container {
  padding: var(--space-4);
}

/* 平板 */
@media (min-width: 768px) {
  .container {
    padding: var(--space-6);
  }
}

/* 桌面 */
@media (min-width: 1024px) {
  .container {
    padding: var(--space-8);
  }
}
```

### 网格自适应

```css
/* 自适应列数 */
.responsive-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-4);
}

/* 固定断点 */
@media (max-width: 640px) {
  .responsive-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 641px) and (max-width: 1024px) {
  .responsive-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1025px) {
  .responsive-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
```