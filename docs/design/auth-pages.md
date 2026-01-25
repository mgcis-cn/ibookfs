# iBookFS 认证页面设计

## 页面概述

认证系统包含登录和注册两个核心页面，采用邮箱 + 验证码的无密码登录方式，降低用户使用门槛，提升安全性。

### 设计原则
- **简洁优先**：最小化表单字段，减少认知负担
- **信任建立**：清晰的隐私说明和安全提示
- **流畅体验**：即时反馈，智能验证
- **品牌一致**：延续主站的视觉风格

---

## 整体布局架构

### 分屏布局（Desktop）

```
┌─────────────────────────────────────────┐
│              │                          │
│              │                          │
│   品牌展示    │       表单区域           │
│   视觉区域    │                          │
│              │                          │
│              │                          │
└─────────────────────────────────────────┘
    40%                 60%
```

### 居中布局（Mobile）

```
┌─────────────────┐
│                 │
│   Logo          │
│                 │
│   表单区域      │
│                 │
│                 │
└─────────────────┘
```

---

## 1. 登录页面 (Login)

### 功能流程

```
输入邮箱 → 点击"发送验证码" → 接收邮件 → 输入验证码 → 登录成功
```

### HTML 结构

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>登录 - iBookFS</title>
  <link rel="stylesheet" href="/css/design-tokens.css">
  <link rel="stylesheet" href="/css/auth.css">
</head>
<body>
  <div class="auth-page">
    <!-- 左侧品牌区 -->
    <div class="auth-brand">
      <div class="brand-content">
        <a href="/" class="brand-logo">
          <svg class="logo-icon">📚</svg>
          <span class="logo-text">iBookFS</span>
        </a>
        
        <h1 class="brand-title">欢迎回来</h1>
        <p class="brand-description">
          继续管理你的数字藏书，让知识触手可及
        </p>
        
        <div class="brand-visual">
          <img src="/images/auth-illustration.svg" alt="登录插图" />
        </div>
        
        <div class="brand-quote">
          <p class="quote-text">
            "iBookFS 让我的书籍管理变得如此简单，
            OCR 功能更是节省了大量时间。"
          </p>
          <div class="quote-author">
            <img src="/images/avatar-1.jpg" alt="用户头像" />
            <div>
              <p class="author-name">张明</p>
              <p class="author-title">独立研究者</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 右侧表单区 -->
    <div class="auth-form-container">
      <div class="auth-form-wrapper">
        <!-- 返回首页 -->
        <a href="/" class="back-link">
          <svg width="16" height="16" viewBox="0 0 16 16">
            <path d="M10 12L6 8l4-4" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
          返回首页
        </a>
        
        <!-- 表单头部 -->
        <div class="form-header">
          <h2 class="form-title">登录账号</h2>
          <p class="form-subtitle">
            使用邮箱验证码快速登录
          </p>
        </div>
        
        <!-- 登录表单 -->
        <form class="auth-form" id="loginForm">
          <!-- 邮箱输入 -->
          <div class="form-group">
            <label for="email" class="form-label">邮箱地址</label>
            <div class="input-wrapper">
              <input 
                type="email" 
                id="email" 
                name="email"
                class="form-input" 
                placeholder="your@email.com"
                required
                autocomplete="email"
              />
              <svg class="input-icon" width="20" height="20" viewBox="0 0 20 20">
                <path d="M2.5 6.5l7.5 5 7.5-5M2.5 6.5v7a1 1 0 001 1h13a1 1 0 001-1v-7m-15 0a1 1 0 011-1h13a1 1 0 011 1" stroke="currentColor" stroke-width="1.5" fill="none"/>
              </svg>
            </div>
            <span class="form-error" id="emailError"></span>
          </div>
          
          <!-- 验证码输入 -->
          <div class="form-group" id="codeGroup" style="display: none;">
            <label for="code" class="form-label">验证码</label>
            <div class="code-input-wrapper">
              <input 
                type="text" 
                id="code" 
                name="code"
                class="form-input" 
                placeholder="输入 6 位验证码"
                maxlength="6"
                pattern="[0-9]{6}"
                autocomplete="one-time-code"
              />
              <button type="button" class="resend-btn" id="resendBtn" disabled>
                <span class="countdown">60s</span>
              </button>
            </div>
            <span class="form-hint">
              验证码已发送至 <strong id="emailDisplay"></strong>
            </span>
            <span class="form-error" id="codeError"></span>
          </div>
          
          <!-- 提交按钮 -->
          <button type="submit" class="button-primary button-lg button-full" id="submitBtn">
            <span class="btn-text">发送验证码</span>
            <span class="btn-loading" style="display: none;">
              <svg class="spinner" width="20" height="20" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="8" stroke="currentColor" stroke-width="2" fill="none" opacity="0.25"/>
                <path d="M10 2a8 8 0 018 8" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
              处理中...
            </span>
          </button>
        </form>
        
        <!-- 分割线 -->
        <div class="divider">
          <span>或</span>
        </div>
        
        <!-- 注册引导 -->
        <div class="form-footer">
          <p class="footer-text">
            还没有账号？
            <a href="/register" class="footer-link">立即注册</a>
          </p>
        </div>
        
        <!-- 隐私说明 -->
        <p class="privacy-notice">
          登录即表示您同意我们的
          <a href="/terms">服务条款</a>
          和
          <a href="/privacy">隐私政策</a>
        </p>
      </div>
    </div>
  </div>
  
  <script src="/js/auth-login.js"></script>
</body>
</html>
```

---

## 2. 注册页面 (Register)

### 功能流程

```
输入姓名 → 输入邮箱 → 发送验证码 → 输入验证码 → 注册成功
```

### HTML 结构

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>注册 - iBookFS</title>
  <link rel="stylesheet" href="/css/design-tokens.css">
  <link rel="stylesheet" href="/css/auth.css">
</head>
<body>
  <div class="auth-page">
    <!-- 左侧品牌区 -->
    <div class="auth-brand">
      <div class="brand-content">
        <a href="/" class="brand-logo">
          <svg class="logo-icon">📚</svg>
          <span class="logo-text">iBookFS</span>
        </a>
        
        <h1 class="brand-title">开始数字化之旅</h1>
        <p class="brand-description">
          加入数千名用户，让纸质书籍焕发新生
        </p>
        
        <div class="brand-visual">
          <img src="/images/register-illustration.svg" alt="注册插图" />
        </div>
        
        <!-- 特性列表 -->
        <ul class="brand-features">
          <li class="feature-item">
            <svg class="feature-icon" width="20" height="20">
              <path d="M5 10l3 3 7-7" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
            <span>无限书籍存储</span>
          </li>
          <li class="feature-item">
            <svg class="feature-icon" width="20" height="20">
              <path d="M5 10l3 3 7-7" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
            <span>智能 OCR 识别</span>
          </li>
          <li class="feature-item">
            <svg class="feature-icon" width="20" height="20">
              <path d="M5 10l3 3 7-7" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
            <span>多设备同步</span>
          </li>
        </ul>
      </div>
    </div>
    
    <!-- 右侧表单区 -->
    <div class="auth-form-container">
      <div class="auth-form-wrapper">
        <!-- 返回首页 -->
        <a href="/" class="back-link">
          <svg width="16" height="16" viewBox="0 0 16 16">
            <path d="M10 12L6 8l4-4" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
          返回首页
        </a>
        
        <!-- 表单头部 -->
        <div class="form-header">
          <h2 class="form-title">创建账号</h2>
          <p class="form-subtitle">
            填写基本信息，开始使用 iBookFS
          </p>
        </div>
        
        <!-- 注册表单 -->
        <form class="auth-form" id="registerForm">
          <!-- 姓名输入 -->
          <div class="form-row">
            <div class="form-group">
              <label for="firstName" class="form-label">名字</label>
              <input 
                type="text" 
                id="firstName" 
                name="firstName"
                class="form-input" 
                placeholder="张"
                required
                autocomplete="given-name"
              />
              <span class="form-error" id="firstNameError"></span>
            </div>
            
            <div class="form-group">
              <label for="lastName" class="form-label">姓氏</label>
              <input 
                type="text" 
                id="lastName" 
                name="lastName"
                class="form-input" 
                placeholder="三"
                required
                autocomplete="family-name"
              />
              <span class="form-error" id="lastNameError"></span>
            </div>
          </div>
          
          <!-- 邮箱输入 -->
          <div class="form-group">
            <label for="email" class="form-label">邮箱地址</label>
            <div class="input-wrapper">
              <input 
                type="email" 
                id="email" 
                name="email"
                class="form-input" 
                placeholder="your@email.com"
                required
                autocomplete="email"
              />
              <svg class="input-icon" width="20" height="20" viewBox="0 0 20 20">
                <path d="M2.5 6.5l7.5 5 7.5-5M2.5 6.5v7a1 1 0 001 1h13a1 1 0 001-1v-7m-15 0a1 1 0 011-1h13a1 1 0 011 1" stroke="currentColor" stroke-width="1.5" fill="none"/>
              </svg>
            </div>
            <span class="form-hint">我们将向此邮箱发送验证码</span>
            <span class="form-error" id="emailError"></span>
          </div>
          
          <!-- 验证码输入 -->
          <div class="form-group" id="codeGroup" style="display: none;">
            <label for="code" class="form-label">验证码</label>
            <div class="code-input-wrapper">
              <input 
                type="text" 
                id="code" 
                name="code"
                class="form-input" 
                placeholder="输入 6 位验证码"
                maxlength="6"
                pattern="[0-9]{6}"
                autocomplete="one-time-code"
              />
              <button type="button" class="resend-btn" id="resendBtn" disabled>
                <span class="countdown">60s</span>
              </button>
            </div>
            <span class="form-hint">
              验证码已发送至 <strong id="emailDisplay"></strong>
            </span>
            <span class="form-error" id="codeError"></span>
          </div>
          
          <!-- 服务条款 -->
          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" name="terms" required />
              <span class="checkbox-text">
                我已阅读并同意
                <a href="/terms" target="_blank">服务条款</a>
                和
                <a href="/privacy" target="_blank">隐私政策</a>
              </span>
            </label>
          </div>
          
          <!-- 提交按钮 -->
          <button type="submit" class="button-primary button-lg button-full" id="submitBtn">
            <span class="btn-text">发送验证码</span>
            <span class="btn-loading" style="display: none;">
              <svg class="spinner" width="20" height="20" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="8" stroke="currentColor" stroke-width="2" fill="none" opacity="0.25"/>
                <path d="M10 2a8 8 0 018 8" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
              处理中...
            </span>
          </button>
        </form>
        
        <!-- 分割线 -->
        <div class="divider">
          <span>或</span>
        </div>
        
        <!-- 登录引导 -->
        <div class="form-footer">
          <p class="footer-text">
            已有账号？
            <a href="/login" class="footer-link">立即登录</a>
          </p>
        </div>
      </div>
    </div>
  </div>
  
  <script src="/js/auth-register.js"></script>
</body>
</html>
```

---

## 3. 样式定义 (auth.css)

```css
/* ========================================
   认证页面布局
   ======================================== */

.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 40% 60%;
  background: var(--bg-primary);
}

/* ========================================
   左侧品牌区
   ======================================== */

.auth-brand {
  position: relative;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-dark) 100%);
  color: white;
  padding: var(--space-8);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.auth-brand::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  border-radius: 50%;
  animation: float 20s ease-in-out infinite;
}

.auth-brand::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255,255,255,0.08) 0%, transparent 70%);
  border-radius: 50%;
  animation: float 15s ease-in-out infinite reverse;
}

.brand-content {
  position: relative;
  z-index: 1;
  max-width: 480px;
  width: 100%;
}

.brand-logo {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: var(--font-weight-bold);
  color: white;
  text-decoration: none;
  margin-bottom: var(--space-8);
  transition: opacity var(--duration-fast) var(--ease-out);
}

.brand-logo:hover {
  opacity: 0.9;
}

.brand-logo .logo-icon {
  width: 32px;
  height: 32px;
}

.brand-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--leading-tight);
  margin-bottom: var(--space-4);
  color: white;
}

.brand-description {
  font-size: var(--text-lg);
  line-height: var(--leading-relaxed);
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: var(--space-8);
}

.brand-visual {
  margin: var(--space-10) 0;
  border-radius: var(--radius-xl);
  overflow: hidden;
}

.brand-visual img {
  width: 100%;
  height: auto;
  display: block;
  opacity: 0.95;
}

/* 用户评价 */
.brand-quote {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.quote-text {
  font-size: var(--text-base);
  line-height: var(--leading-relaxed);
  font-style: italic;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: var(--space-4);
}

.quote-author {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.quote-author img {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.author-name {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-semibold);
  color: white;
}

.author-title {
  font-size: var(--text-xs);
  color: rgba(255, 255, 255, 0.8);
}

/* 特性列表 */
.brand-features {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin-top: var(--space-8);
}

.feature-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-size: var(--text-base);
  color: rgba(255, 255, 255, 0.95);
}

.feature-icon {
  flex-shrink: 0;
  color: rgba(255, 255, 255, 0.9);
}

/* ========================================
   右侧表单区
   ======================================== */

.auth-form-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-8);
  background: var(--bg-primary);
}

.auth-form-wrapper {
  width: 100%;
  max-width: 440px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-8);
  transition: color var(--duration-fast) var(--ease-out);
}

.back-link:hover {
  color: var(--color-accent);
}

.form-header {
  margin-bottom: var(--space-8);
}

.form-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-2);
  color: var(--text-primary);
}

.form-subtitle {
  font-size: var(--text-base);
  color: var(--text-secondary);
}

/* ========================================
   表单样式
   ======================================== */

.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-label {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.input-wrapper {
  position: relative;
}

.form-input {
  width: 100%;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  color: var(--text-primary);
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.input-wrapper .form-input {
  padding-right: var(--space-10);
}

.form-input::placeholder {
  color: var(--text-tertiary);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.form-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: var(--bg-secondary);
}

.form-input.error {
  border-color: var(--color-error);
}

.form-input.error:focus {
  box-shadow: 0 0 0 3px rgba(184, 92, 92, 0.2);
}

.input-icon {
  position: absolute;
  right: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-tertiary);
  pointer-events: none;
}

/* 验证码输入 */
.code-input-wrapper {
  position: relative;
  display: flex;
  gap: var(--space-2);
}

.code-input-wrapper .form-input {
  flex: 1;
  font-family: var(--font-mono);
  font-size: var(--text-lg);
  letter-spacing: 0.5em;
  text-align: center;
}

.resend-btn {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-accent);
  background: var(--bg-secondary);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  white-space: nowrap;
  transition: all var(--duration-fast) var(--ease-out);
}

.resend-btn:hover:not(:disabled) {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
}

.resend-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.countdown {
  font-family: var(--font-mono);
}

/* 提示和错误信息 */
.form-hint {
  font-size: var(--text-xs);
  color: var(--text-secondary);
}

.form-hint strong {
  color: var(--text-primary);
  font-weight: var(--font-weight-medium);
}

.form-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  display: none;
}

.form-error:not(:empty) {
  display: block;
}

/* 复选框 */
.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: var(--space-2);
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  margin-top: 2px;
  cursor: pointer;
  accent-color: var(--color-accent);
}

.checkbox-text {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  line-height: var(--leading-normal);
}

.checkbox-text a {
  color: var(--color-accent);
  text-decoration: underline;
}

/* 按钮 */
.button-full {
  width: 100%;
}

.button-primary {
  position: relative;
}

.btn-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
}

.spinner {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 分割线 */
.divider {
  position: relative;
  text-align: center;
  margin: var(--space-6) 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--border-subtle);
}

.divider span {
  position: relative;
  display: inline-block;
  padding: 0 var(--space-4);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  background: var(--bg-primary);
}

/* 页脚 */
.form-footer {
  text-align: center;
}

.footer-text {
  font-size: var(--text-base);
  color: var(--text-secondary);
}

.footer-link {
  color: var(--color-accent);
  font-weight: var(--font-weight-medium);
  transition: color var(--duration-fast) var(--ease-out);
}

.footer-link:hover {
  color: var(--color-accent-dark);
}

.privacy-notice {
  margin-top: var(--space-6);
  padding-top: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  font-size: var(--text-xs);
  color: var(--text-tertiary);
  text-align: center;
  line-height: var(--leading-relaxed);
}

.privacy-notice a {
  color: var(--text-secondary);
  text-decoration: underline;
}

/* ========================================
   响应式设计
   ======================================== */

@media (max-width: 1024px) {
  .auth-page {
    grid-template-columns: 1fr;
  }
  
  .auth-brand {
    display: none;
  }
  
  .auth-form-container {
    padding: var(--space-6) var(--space-4);
  }
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-title {
    font-size: var(--text-2xl);
  }
  
  .code-input-wrapper {
    flex-direction: column;
  }
  
  .resend-btn {
    width: 100%;
  }
}

/* ========================================
   动画
   ======================================== */

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(5deg);
  }
}

/* 表单进入动画 */
.auth-form-wrapper {
  animation: fadeInUp 0.6s var(--ease-out);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 成功状态 */
.form-input.success {
  border-color: var(--color-success);
}

.form-input.success:focus {
  box-shadow: 0 0 0 3px rgba(122, 155, 118, 0.2);
}

/* 加载状态 */
.button-primary.loading {
  pointer-events: none;
}

.button-primary.loading .btn-text {
  display: none;
}

.button-primary.loading .btn-loading {
  display: flex;
}
```

---

## 4. JavaScript 交互逻辑

### 登录页面 (auth-login.js)

```javascript
// 登录页面交互逻辑
class LoginForm {
  constructor() {
    this.form = document.getElementById('loginForm');
    this.emailInput = document.getElementById('email');
    this.codeInput = document.getElementById('code');
    this.codeGroup = document.getElementById('codeGroup');
    this.submitBtn = document.getElementById('submitBtn');
    this.resendBtn = document.getElementById('resendBtn');
    this.emailDisplay = document.getElementById('emailDisplay');
    
    this.step = 1; // 1: 输入邮箱, 2: 输入验证码
    this.countdown = 60;
    this.countdownTimer = null;
    
    this.init();
  }
  
  init() {
    this.form.addEventListener('submit', (e) => this.handleSubmit(e));
    this.resendBtn.addEventListener('click', () => this.resendCode());
    this.emailInput.addEventListener('input', () => this.validateEmail());
    this.codeInput.addEventListener('input', () => this.validateCode());
  }
  
  async handleSubmit(e) {
    e.preventDefault();
    
    if (this.step === 1) {
      await this.sendCode();
    } else {
      await this.verifyCode();
    }
  }
  
  validateEmail() {
    const email = this.emailInput.value.trim();
    const emailError = document.getElementById('emailError');
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (!email) {
      emailError.textContent = '请输入邮箱地址';
      this.emailInput.classList.add('error');
      return false;
    }
    
    if (!emailRegex.test(email)) {
      emailError.textContent = '请输入有效的邮箱地址';
      this.emailInput.classList.add('error');
      return false;
    }
    
    emailError.textContent = '';
    this.emailInput.classList.remove('error');
    this.emailInput.classList.add('success');
    return true;
  }
  
  validateCode() {
    const code = this.codeInput.value.trim();
    const codeError = document.getElementById('codeError');
    
    if (!code) {
      codeError.textContent = '请输入验证码';
      this.codeInput.classList.add('error');
      return false;
    }
    
    if (!/^\d{6}$/.test(code)) {
      codeError.textContent = '验证码必须是 6 位数字';
      this.codeInput.classList.add('error');
      return false;
    }
    
    codeError.textContent = '';
    this.codeInput.classList.remove('error');
    this.codeInput.classList.add('success');
    return true;
  }
  
  async sendCode() {
    if (!this.validateEmail()) return;
    
    this.setLoading(true);
    
    try {
      // 调用 API 发送验证码
      const response = await fetch('/api/auth/send-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: this.emailInput.value.trim() })
      });
      
      const data = await response.json();
      
      if (response.ok) {
        this.showCodeInput();
        this.startCountdown();
        this.showToast('验证码已发送', 'success');
      } else {
        throw new Error(data.message || '发送失败');
      }
    } catch (error) {
      this.showToast(error.message, 'error');
    } finally {
      this.setLoading(false);
    }
  }
  
  async verifyCode() {
    if (!this.validateCode()) return;
    
    this.setLoading(true);
    
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: this.emailInput.value.trim(),
          code: this.codeInput.value.trim()
        })
      });
      
      const data = await response.json();
      
      if (response.ok) {
        this.showToast('登录成功', 'success');
        setTimeout(() => {
          window.location.href = '/dashboard';
        }, 1000);
      } else {
        throw new Error(data.message || '验证失败');
      }
    } catch (error) {
      this.showToast(error.message, 'error');
      document.getElementById('codeError').textContent = error.message;
      this.codeInput.classList.add('error');
    } finally {
      this.setLoading(false);
    }
  }
  
  showCodeInput() {
    this.step = 2;
    this.codeGroup.style.display = 'block';
    this.emailDisplay.textContent = this.emailInput.value;
    this.emailInput.disabled = true;
    this.submitBtn.querySelector('.btn-text').textContent = '登录';
    this.codeInput.focus();
  }
  
  async resendCode() {
    await this.sendCode();
  }
  
  startCountdown() {
    this.countdown = 60;
    this.resendBtn.disabled = true;
    
    this.countdownTimer = setInterval(() => {
      this.countdown--;
      this.resendBtn.querySelector('.countdown').textContent = `${this.countdown}s`;
      
      if (this.countdown <= 0) {
        clearInterval(this.countdownTimer);
        this.resendBtn.disabled = false;
        this.resendBtn.querySelector('.countdown').textContent = '重新发送';
      }
    }, 1000);
  }
  
  setLoading(loading) {
    this.submitBtn.classList.toggle('loading', loading);
    this.submitBtn.disabled = loading;
  }
  
  showToast(message, type = 'info') {
    // 实现 Toast 通知
    console.log(`[${type}] ${message}`);
  }
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
  new LoginForm();
});
```

### 注册页面 (auth-register.js)

```javascript
// 注册页面交互逻辑
class RegisterForm {
  constructor() {
    this.form = document.getElementById('registerForm');
    this.firstNameInput = document.getElementById('firstName');
    this.lastNameInput = document.getElementById('lastName');
    this.emailInput = document.getElementById('email');
    this.codeInput = document.getElementById('code');
    this.codeGroup = document.getElementById('codeGroup');
    this.submitBtn = document.getElementById('submitBtn');
    this.resendBtn = document.getElementById('resendBtn');
    this.emailDisplay = document.getElementById('emailDisplay');
    
    this.step = 1; // 1: 填写信息, 2: 输入验证码
    this.countdown = 60;
    this.countdownTimer = null;
    
    this.init();
  }
  
  init() {
    this.form.addEventListener('submit', (e) => this.handleSubmit(e));
    this.resendBtn.addEventListener('click', () => this.resendCode());
    
    this.firstNameInput.addEventListener('input', () => this.validateName());
    this.lastNameInput.addEventListener('input', () => this.validateName());
    this.emailInput.addEventListener('input', () => this.validateEmail());
    this.codeInput.addEventListener('input', () => this.validateCode());
  }
  
  async handleSubmit(e) {
    e.preventDefault();
    
    if (this.step === 1) {
      await this.sendCode();
    } else {
      await this.register();
    }
  }
  
  validateName() {
    const firstName = this.firstNameInput.value.trim();
    const lastName = this.lastNameInput.value.trim();
    const firstNameError = document.getElementById('firstNameError');
    const lastNameError = document.getElementById('lastNameError');
    
    let isValid = true;
    
    if (!firstName) {
      firstNameError.textContent = '请输入名字';
      this.firstNameInput.classList.add('error');
      isValid = false;
    } else {
      firstNameError.textContent = '';
      this.firstNameInput.classList.remove('error');
      this.firstNameInput.classList.add('success');
    }
    
    if (!lastName) {
      lastNameError.textContent = '请输入姓氏';
      this.lastNameInput.classList.add('error');
      isValid = false;
    } else {
      lastNameError.textContent = '';
      this.lastNameInput.classList.remove('error');
      this.lastNameInput.classList.add('success');
    }
    
    return isValid;
  }
  
  validateEmail() {
    const email = this.emailInput.value.trim();
    const emailError = document.getElementById('emailError');
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (!email) {
      emailError.textContent = '请输入邮箱地址';
      this.emailInput.classList.add('error');
      return false;
    }
    
    if (!emailRegex.test(email)) {
      emailError.textContent = '请输入有效的邮箱地址';
      this.emailInput.classList.add('error');
      return false;
    }
    
    emailError.textContent = '';
    this.emailInput.classList.remove('error');
    this.emailInput.classList.add('success');
    return true;
  }
  
  validateCode() {
    const code = this.codeInput.value.trim();
    const codeError = document.getElementById('codeError');
    
    if (!code) {
      codeError.textContent = '请输入验证码';
      this.codeInput.classList.add('error');
      return false;
    }
    
    if (!/^\d{6}$/.test(code)) {
      codeError.textContent = '验证码必须是 6 位数字';
      this.codeInput.classList.add('error');
      return false;
    }
    
    codeError.textContent = '';
    this.codeInput.classList.remove('error');
    this.codeInput.classList.add('success');
    return true;
  }
  
  async sendCode() {
    if (!this.validateName() || !this.validateEmail()) return;
    
    this.setLoading(true);
    
    try {
      const response = await fetch('/api/auth/send-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          email: this.emailInput.value.trim(),
          type: 'register'
        })
      });
      
      const data = await response.json();
      
      if (response.ok) {
        this.showCodeInput();
        this.startCountdown();
        this.showToast('验证码已发送', 'success');
      } else {
        throw new Error(data.message || '发送失败');
      }
    } catch (error) {
      this.showToast(error.message, 'error');
    } finally {
      this.setLoading(false);
    }
  }
  
  async register() {
    if (!this.validateCode()) return;
    
    this.setLoading(true);
    
    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          firstName: this.firstNameInput.value.trim(),
          lastName: this.lastNameInput.value.trim(),
          email: this.emailInput.value.trim(),
          code: this.codeInput.value.trim()
        })
      });
      
      const data = await response.json();
      
      if (response.ok) {
        this.showToast('注册成功', 'success');
        setTimeout(() => {
          window.location.href = '/dashboard';
        }, 1000);
      } else {
        throw new Error(data.message || '注册失败');
      }
    } catch (error) {
      this.showToast(error.message, 'error');
      document.getElementById('codeError').textContent = error.message;
      this.codeInput.classList.add('error');
    } finally {
      this.setLoading(false);
    }
  }
  
  showCodeInput() {
    this.step = 2;
    this.codeGroup.style.display = 'block';
    this.emailDisplay.textContent = this.emailInput.value;
    this.firstNameInput.disabled = true;
    this.lastNameInput.disabled = true;
    this.emailInput.disabled = true;
    this.submitBtn.querySelector('.btn-text').textContent = '完成注册';
    this.codeInput.focus();
  }
  
  async resendCode() {
    await this.sendCode();
  }
  
  startCountdown() {
    this.countdown = 60;
    this.resendBtn.disabled = true;
    
    this.countdownTimer = setInterval(() => {
      this.countdown--;
      this.resendBtn.querySelector('.countdown').textContent = `${this.countdown}s`;
      
      if (this.countdown <= 0) {
        clearInterval(this.countdownTimer);
        this.resendBtn.disabled = false;
        this.resendBtn.querySelector('.countdown').textContent = '重新发送';
      }
    }, 1000);
  }
  
  setLoading(loading) {
    this.submitBtn.classList.toggle('loading', loading);
    this.submitBtn.disabled = loading;
  }
  
  showToast(message, type = 'info') {
    console.log(`[${type}] ${message}`);
  }
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
  new RegisterForm();
});
```

---

## 5. 交互细节

### 表单验证

**实时验证：**
- 邮箱格式验证（正则表达式）
- 姓名非空验证
- 验证码格式验证（6位数字）

**视觉反馈：**
- ✅ 成功：绿色边框
- ❌ 错误：红色边框 + 错误提示
- 🔄 输入中：默认边框
- 👁️ 聚焦：强调色边框 + 外发光

### 验证码流程

**发送验证码：**
1. 验证邮箱格式
2. 调用 API 发送验证码
3. 显示验证码输入框
4. 禁用邮箱输入
5. 开始 60 秒倒计时

**重新发送：**
- 倒计时结束后启用
- 点击后重新发送
- 重置倒计时

**验证码输入：**
- 等宽字体显示
- 自动聚焦
- 6 位数字限制
- 实时验证

### 加载状态

**按钮状态：**
- 默认：显示文本
- 加载中：显示 Spinner + "处理中..."
- 禁用：降低透明度，禁止点击

### 错误处理

**常见错误：**
- 邮箱已注册
- 验证码错误
- 验证码过期
- 网络错误

**错误展示：**
- 表单字段下方显示错误信息
- Toast 通知全局错误
- 红色边框标识错误字段

---

## 6. 可访问性

### 键盘导航
- Tab 键切换表单字段
- Enter 键提交表单
- Esc 键清除错误提示

### 屏幕阅读器
- 语义化 HTML 标签
- 适当的 label 关联
- ARIA 属性标注

### 表单自动填充
- `autocomplete` 属性
- 支持密码管理器
- 支持浏览器自动填充

---

## 7. 安全考虑

### 验证码安全
- 6 位随机数字
- 5 分钟有效期
- 单次使用
- 限制发送频率（60 秒）

### 邮箱验证
- 格式验证
- 域名验证
- 防止重复注册

### CSRF 防护
- Token 验证
- SameSite Cookie

### 速率限制
- 发送验证码：60 秒/次
- 登录尝试：5 次/小时
- 注册尝试：3 次/小时

---

## 8. 性能优化

### 首屏加载
- 关键 CSS 内联
- 延迟加载非关键资源
- 预连接 API 域名

### 表单优化
- 防抖输入验证
- 异步验证
- 乐观 UI 更新

### 图片优化
- SVG 图标
- WebP 格式
- 懒加载背景图

---

## 9. 响应式设计

### 移动端优化
- 单列布局
- 大号按钮（最小 44x44px）
- 自动聚焦输入框
- 数字键盘（验证码输入）

### 平板适配
- 保持分屏布局
- 调整间距和字号

### 桌面优化
- 分屏布局
- 悬停效果
- 键盘快捷键