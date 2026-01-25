# iBookFS 组件规范

本文档定义 iBookFS 项目中所有核心组件的设计规范和实现指南。

---

## 按钮组件 (Button)

### 变体

#### Primary Button - 主要按钮
用于最重要的操作，每个页面/区域应只有一个主要按钮。

```css
.button-primary {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: white;
  background-color: var(--color-accent);
  border-radius: var(--radius-md);
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  box-shadow: var(--shadow-sm);
}

.button-primary:hover {
  background-color: var(--color-accent-dark);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.button-primary:active {
  transform: translateY(0);
  box-shadow: var(--shadow-xs);
}

.button-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}
```

**使用场景：** 上传照片、创建书籍、保存更改

#### Secondary Button - 次要按钮

```css
.button-secondary {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  background-color: transparent;
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.button-secondary:hover {
  background-color: var(--bg-secondary);
  border-color: var(--border-strong);
}
```

**使用场景：** 取消、返回、次要操作

#### Text Button - 文本按钮

```css
.button-text {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-base);
  font-weight: var(--font-weight-regular);
  color: var(--text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  transition: color var(--duration-fast) var(--ease-out);
}

.button-text:hover {
  color: var(--text-primary);
}
```

**使用场景：** 低优先级操作、内联链接

### 尺寸变体

```css
/* Small */
.button-sm {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
}

/* Medium (默认) */
.button-md {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
}

/* Large */
.button-lg {
  padding: var(--space-4) var(--space-8);
  font-size: var(--text-lg);
}
```

### 图标按钮

```css
.button-icon {
  width: 40px;
  height: 40px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
}
```

---

## 卡片组件 (Card)

### Book Card - 书籍卡片

书架中展示书籍的核心组件。

```html
<div class="card-book">
  <div class="card-book-cover">
    <img src="cover.jpg" alt="书籍封面" />
    <div class="card-book-badge">120页</div>
  </div>
  <div class="card-book-content">
    <h3 class="card-book-title">书籍标题</h3>
    <p class="card-book-meta">作者名 · 2024</p>
    <div class="card-book-progress">
      <div class="progress-bar" style="width: 60%"></div>
    </div>
  </div>
</div>
```

```css
.card-book {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-book);
  transition: all var(--duration-base) var(--ease-out);
  cursor: pointer;
}

.card-book:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-book-hover);
}

.card-book-cover {
  position: relative;
  aspect-ratio: 3 / 4;
  overflow: hidden;
  background: var(--bg-secondary);
}

.card-book-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--duration-slow) var(--ease-out);
}

.card-book:hover .card-book-cover img {
  transform: scale(1.05);
}

.card-book-badge {
  position: absolute;
  bottom: var(--space-2);
  right: var(--space-2);
  padding: var(--space-1) var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-weight-medium);
  color: white;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-sm);
}

.card-book-content {
  padding: var(--space-4);
}

.card-book-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-book-meta {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-3);
}

.card-book-progress {
  height: 4px;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--color-accent);
  border-radius: var(--radius-full);
  transition: width var(--duration-slow) var(--ease-out);
}
```

### Photo Card - 照片卡片

```html
<div class="card-photo">
  <img src="photo.jpg" alt="书籍页面" />
  <div class="card-photo-overlay">
    <span class="card-photo-page">第 5 页</span>
  </div>
</div>
```

```css
.card-photo {
  position: relative;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--duration-base) var(--ease-out);
  cursor: pointer;
}

.card-photo:hover {
  box-shadow: var(--shadow-md);
}

.card-photo img {
  width: 100%;
  height: auto;
  display: block;
}

.card-photo-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: var(--space-3);
  background: linear-gradient(to top, rgba(0,0,0,0.6), transparent);
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.card-photo:hover .card-photo-overlay {
  opacity: 1;
}

.card-photo-page {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: white;
}
```

---

## 输入组件 (Input)

### Text Input - 文本输入框

```html
<div class="input-group">
  <label class="input-label" for="book-title">书籍标题</label>
  <input 
    type="text" 
    id="book-title" 
    class="input-text" 
    placeholder="输入书籍标题"
  />
  <span class="input-hint">必填项</span>
</div>
```

```css
.input-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input-label {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.input-text {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  color: var(--text-primary);
  background: var(--bg-elevated);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.input-text::placeholder {
  color: var(--text-tertiary);
  font-style: italic;
}

.input-text:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.input-text:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-hint {
  font-size: var(--text-xs);
  color: var(--text-secondary);
}

.input-error {
  border-color: var(--color-error);
}

.input-error:focus {
  box-shadow: 0 0 0 3px rgba(184, 92, 92, 0.2);
}
```

### File Upload - 文件上传

```html
<div class="upload-area">
  <input type="file" id="file-input" class="upload-input" multiple accept="image/*" />
  <label for="file-input" class="upload-label">
    <svg class="upload-icon"><!-- 上传图标 --></svg>
    <span class="upload-text">拖拽照片到此处，或点击选择</span>
    <span class="upload-hint">支持 JPG、PNG 格式，最大 10MB</span>
  </label>
</div>
```

```css
.upload-area {
  position: relative;
}

.upload-input {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
  overflow: hidden;
}

.upload-label {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 240px;
  padding: var(--space-8);
  border: 2px dashed var(--border-medium);
  border-radius: var(--radius-lg);
  background: var(--bg-elevated);
  cursor: pointer;
  transition: all var(--duration-base) var(--ease-out);
}

.upload-label:hover {
  border-color: var(--color-accent);
  background: var(--bg-secondary);
}

.upload-area.dragging .upload-label {
  border-color: var(--color-accent);
  background: var(--color-accent-light);
}

.upload-icon {
  width: 48px;
  height: 48px;
  margin-bottom: var(--space-4);
  color: var(--text-secondary);
}

.upload-text {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.upload-hint {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}
```

---

## 模态框 (Modal)

```html
<div class="modal-backdrop">
  <div class="modal">
    <div class="modal-header">
      <h2 class="modal-title">创建新书籍</h2>
      <button class="modal-close">×</button>
    </div>
    <div class="modal-body">
      <!-- 内容 -->
    </div>
    <div class="modal-footer">
      <button class="button-secondary">取消</button>
      <button class="button-primary">确认</button>
    </div>
  </div>
</div>
```

```css
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: var(--bg-overlay);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal-backdrop);
  animation: fadeIn var(--duration-base) var(--ease-out);
}

.modal {
  width: 90%;
  max-width: 560px;
  max-height: 90vh;
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  overflow: hidden;
  animation: modalSlideIn var(--duration-base) var(--ease-out);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
}

.modal-title {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
}

.modal-close {
  width: 32px;
  height: 32px;
  font-size: var(--text-2xl);
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.modal-close:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.modal-body {
  padding: var(--space-6);
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  gap: var(--space-3);
  justify-content: flex-end;
  padding: var(--space-6);
  border-top: 1px solid var(--border-subtle);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
```

---

## Toast 通知

```html
<div class="toast toast-success">
  <svg class="toast-icon"><!-- 图标 --></svg>
  <span class="toast-message">书籍创建成功</span>
  <button class="toast-close">×</button>
</div>
```

```css
.toast {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 320px;
  padding: var(--space-4) var(--space-5);
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  animation: toastSlideIn var(--duration-base) var(--ease-spring);
}

.toast-success {
  border-left: 4px solid var(--color-success);
}

.toast-warning {
  border-left: 4px solid var(--color-warning);
}

.toast-error {
  border-left: 4px solid var(--color-error);
}

.toast-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.toast-message {
  flex: 1;
  font-size: var(--text-base);
  color: var(--text-primary);
}

.toast-close {
  width: 24px;
  height: 24px;
  font-size: var(--text-xl);
  color: var(--text-secondary);
  flex-shrink: 0;
}

@keyframes toastSlideIn {
  from {
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
```

---

## 加载状态

### Skeleton - 骨架屏

```html
<div class="skeleton-book-card">
  <div class="skeleton skeleton-cover"></div>
  <div class="skeleton skeleton-title"></div>
  <div class="skeleton skeleton-meta"></div>
</div>
```

```css
.skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 0%,
    var(--bg-tertiary) 50%,
    var(--bg-secondary) 100%
  );
  background-size: 200% 100%;
  animation: skeletonLoading 1.5s ease-in-out infinite;
  border-radius: var(--radius-md);
}

.skeleton-cover {
  aspect-ratio: 3 / 4;
  margin-bottom: var(--space-4);
}

.skeleton-title {
  height: 24px;
  width: 80%;
  margin-bottom: var(--space-2);
}

.skeleton-meta {
  height: 16px;
  width: 60%;
}

@keyframes skeletonLoading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

### Spinner - 加载指示器

```html
<div class="spinner"></div>
```

```css
.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-subtle);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
```

---

## 空状态

```html
<div class="empty-state">
  <svg class="empty-icon"><!-- 插图 --></svg>
  <h3 class="empty-title">还没有书籍</h3>
  <p class="empty-description">开始创建你的第一本数字书籍</p>
  <button class="button-primary">创建书籍</button>
</div>
```

```css
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12);
  text-align: center;
}

.empty-icon {
  width: 120px;
  height: 120px;
  margin-bottom: var(--space-6);
  color: var(--text-tertiary);
  opacity: 0.5;
}

.empty-title {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
}

.empty-description {
  font-size: var(--text-base);
  color: var(--text-secondary);
  margin-bottom: var(--space-6);
  max-width: 400px;
}
```

---

## 导航组件

### Sidebar - 侧边栏 (Desktop)

```html
<nav class="sidebar">
  <div class="sidebar-header">
    <h1 class="sidebar-logo">iBookFS</h1>
  </div>
  
  <ul class="sidebar-nav">
    <li><a href="#" class="nav-item active">书架</a></li>
    <li><a href="#" class="nav-item">上传</a></li>
    <li><a href="#" class="nav-item">设置</a></li>
  </ul>
  
  <div class="sidebar-footer">
    <div class="user-profile">
      <img src="avatar.jpg" alt="用户头像" />
      <span>用户名</span>
    </div>
  </div>
</nav>
```

```css
.sidebar {
  width: 240px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-elevated);
  border-right: 1px solid var(--border-subtle);
}

.sidebar-header {
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
}

.sidebar-logo {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-accent);
}

.sidebar-nav {
  flex: 1;
  padding: var(--space-4);
  list-style: none;
}

.nav-item {
  display: block;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.nav-item:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
}

.sidebar-footer {
  padding: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}
```

### Bottom Navigation - 底部导航 (Mobile)

```html
<nav class="bottom-nav">
  <a href="#" class="bottom-nav-item active">
    <svg><!-- 图标 --></svg>
    <span>书架</span>
  </a>
  <a href="#" class="bottom-nav-item">
    <svg><!-- 图标 --></svg>
    <span>上传</span>
  </a>
  <a href="#" class="bottom-nav-item">
    <svg><!-- 图标 --></svg>
    <span>我的</span>
  </a>
</nav>
```

```css
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background: var(--bg-elevated);
  border-top: 1px solid var(--border-subtle);
  padding: var(--space-2);
  z-index: var(--z-fixed);
}

.bottom-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2);
  color: var(--text-secondary);
  transition: color var(--duration-fast) var(--ease-out);
}

.bottom-nav-item svg {
  width: 24px;
  height: 24px;
}

.bottom-nav-item span {
  font-size: var(--text-xs);
}

.bottom-nav-item.active {
  color: var(--color-accent);
}
```

---

## 使用指南

### 组件组合示例

**创建书籍流程：**
1. 点击按钮 → 打开模态框
2. 填写表单 → 输入组件
3. 提交成功 → Toast 通知
4. 跳转到上传 → 文件上传组件

**书架展示：**
1. 加载中 → 骨架屏
2. 无数据 → 空状态
3. 有数据 → 书籍卡片网格

### 响应式适配

```css
/* 移动端 */
@media (max-width: 768px) {
  .card-book-title {
    font-size: var(--text-base);
  }
  
  .modal {
    max-width: 100%;
    border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  }
}

/* 平板 */
@media (min-width: 768px) and (max-width: 1024px) {
  .sidebar {
    width: 200px;
  }
}
```