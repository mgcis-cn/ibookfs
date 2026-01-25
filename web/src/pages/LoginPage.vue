<template>
  <div class="auth-page">
    <!-- Left Brand Section -->
    <div class="auth-brand">
      <div class="brand-content">
        <router-link to="/" class="brand-logo">
          <BookOpen :size="32" />
          <span class="logo-text">iBookFS</span>
        </router-link>

        <h1 class="brand-title">欢迎回来</h1>
        <p class="brand-description">继续管理你的数字藏书，让知识触手可及</p>

        <div class="brand-quote">
          <p class="quote-text">
            "iBookFS 让我的书籍管理变得如此简单，OCR 功能更是节省了大量时间。"
          </p>
          <div class="quote-author">
            <div class="author-avatar">张</div>
            <div>
              <p class="author-name">张明</p>
              <p class="author-title">独立研究者</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Form Section -->
    <div class="auth-form-container">
      <div class="auth-form-wrapper">
        <!-- Back to Home -->
        <router-link to="/" class="back-link">
          <ArrowLeft :size="16" />
          返回首页
        </router-link>

        <!-- Form Header -->
        <div class="form-header">
          <h2 class="form-title">登录账号</h2>
          <p class="form-subtitle">使用邮箱验证码快速登录</p>
        </div>

        <!-- Login Form -->
        <form class="auth-form" @submit.prevent="handleSubmit">
          <!-- Email Input -->
          <div class="form-group">
            <label for="email" class="form-label">邮箱地址</label>
            <div class="input-wrapper">
              <input
                id="email"
                v-model="email"
                type="email"
                class="form-input"
                :class="{ error: errors.email, success: !errors.email && email }"
                placeholder="your@email.com"
                :disabled="step === 2"
                required
                autocomplete="email"
              />
              <Mail class="input-icon" :size="20" />
            </div>
            <span v-if="errors.email" class="form-error">{{ errors.email }}</span>
          </div>

          <!-- Verification Code Input -->
          <div v-show="step === 2" class="form-group">
            <label for="code" class="form-label">验证码</label>
            <div class="code-input-wrapper">
              <input
                id="code"
                v-model="code"
                type="text"
                class="form-input code-input"
                :class="{ error: errors.code, success: !errors.code && code }"
                placeholder="输入 6 位验证码"
                maxlength="6"
                required
                autocomplete="one-time-code"
              />
              <button
                type="button"
                class="resend-btn"
                :disabled="countdown > 0"
                @click="handleResend"
              >
                {{ countdown > 0 ? `${countdown}s` : '重新发送' }}
              </button>
            </div>
            <span class="form-hint">
              验证码已发送至 <strong>{{ email }}</strong>
            </span>
            <span v-if="errors.code" class="form-error">{{ errors.code }}</span>
          </div>

          <!-- Demo Hint -->
          <div v-if="step === 2" class="demo-hint">
            <Info :size="16" />
            <span>演示模式：验证码为 <strong>123456</strong></span>
          </div>

          <!-- Submit Button -->
          <button type="submit" class="button-primary button-lg button-full" :disabled="isLoading">
            <span v-if="!isLoading" class="btn-text">{{ step === 1 ? '发送验证码' : '登录' }}</span>
            <span v-else class="btn-loading">
              <svg class="spinner" width="20" height="20" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="8" stroke="currentColor" stroke-width="2" fill="none" opacity="0.25" />
                <path d="M10 2a8 8 0 018 8" stroke="currentColor" stroke-width="2" fill="none" />
              </svg>
              处理中...
            </span>
          </button>
        </form>

        <!-- Divider -->
        <div class="divider">
          <span>或</span>
        </div>

        <!-- Register Link -->
        <div class="form-footer">
          <p class="footer-text">
            还没有账号？
            <router-link to="/register" class="footer-link">立即注册</router-link>
          </p>
        </div>

        <!-- Privacy Notice -->
        <p class="privacy-notice">
          登录即表示您同意我们的
          <a href="/terms">服务条款</a>
          和
          <a href="/privacy">隐私政策</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { BookOpen, ArrowLeft, Mail, Info } from 'lucide-vue-next'
import { sendCode, login } from '@/api'
import { useToastStore } from '@/stores/toast'

const router = useRouter()
const toastStore = useToastStore()

const email = ref('')
const code = ref('')
const step = ref(1)
const isLoading = ref(false)
const countdown = ref(0)
const errors = ref<Record<string, string>>({})

let countdownTimer: ReturnType<typeof setInterval> | null = null

const validateEmail = (): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!email.value) {
    errors.value.email = '请输入邮箱地址'
    return false
  }
  if (!emailRegex.test(email.value)) {
    errors.value.email = '请输入有效的邮箱地址'
    return false
  }
  errors.value.email = ''
  return true
}

const validateCode = (): boolean => {
  if (!code.value) {
    errors.value.code = '请输入验证码'
    return false
  }
  if (!/^\d{6}$/.test(code.value)) {
    errors.value.code = '验证码必须是 6 位数字'
    return false
  }
  errors.value.code = ''
  return true
}

const startCountdown = () => {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (countdownTimer) clearInterval(countdownTimer)
    }
  }, 1000)
}

const handleResend = async () => {
  const response = await sendCode({ email: email.value, type: 'login' })
  if (response.success) {
    toastStore.success('验证码已重新发送')
    startCountdown()
  } else {
    toastStore.error(response.message || '发送失败')
  }
}

const handleSubmit = async () => {
  errors.value = {}

  if (step.value === 1) {
    if (!validateEmail()) return

    isLoading.value = true
    const response = await sendCode({ email: email.value, type: 'login' })
    isLoading.value = false

    if (response.success) {
      step.value = 2
      startCountdown()
      toastStore.success('验证码已发送')
    } else {
      errors.value.email = response.message || '发送失败'
    }
  } else {
    if (!validateCode()) return

    isLoading.value = true
    const response = await login({ email: email.value, code: code.value })
    isLoading.value = false

    if (response.success) {
      toastStore.success('登录成功')
      setTimeout(() => {
        router.push('/library')
      }, 500)
    } else {
      errors.value.code = response.message || '登录失败'
    }
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 40% 60%;
  background: var(--bg-primary);
}

/* Left Brand Section */
.auth-brand {
  position: relative;
  background: linear-gradient(135deg, var(--color-accent) 0%, #a6624a 100%);
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
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
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
  background: radial-gradient(circle, rgba(255, 255, 255, 0.08) 0%, transparent 70%);
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

.logo-text {
  margin-left: var(--space-2);
}

.brand-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  line-height: 1.2;
  margin-bottom: var(--space-4);
  color: white;
}

.brand-description {
  font-size: var(--text-lg);
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: var(--space-8);
}

.brand-quote {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.quote-text {
  font-size: var(--text-base);
  line-height: 1.6;
  font-style: italic;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: var(--space-4);
}

.quote-author {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
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

/* Right Form Section */
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
  animation: fadeInUp 0.6s var(--ease-out);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  margin-bottom: var(--space-8);
  text-decoration: none;
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
  color: var(--color-ink);
}

.form-subtitle {
  font-size: var(--text-base);
  color: var(--color-ink-light);
}

/* Form Styles */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-label {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-ink);
}

.input-wrapper {
  position: relative;
}

.form-input {
  width: 100%;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  color: var(--color-ink);
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.input-wrapper .form-input {
  padding-right: var(--space-10);
}

.form-input::placeholder {
  color: var(--color-ink-light);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(193, 123, 92, 0.2);
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

.form-input.success {
  border-color: var(--color-success);
}

.input-icon {
  position: absolute;
  right: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-ink-light);
  pointer-events: none;
}

/* Code Input */
.code-input-wrapper {
  position: relative;
  display: flex;
  gap: var(--space-2);
}

.code-input {
  font-family: var(--font-mono);
  font-size: var(--text-lg);
  letter-spacing: 0.5em;
  text-align: center;
  flex: 1;
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
  cursor: pointer;
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

/* Hint and Error */
.form-hint {
  font-size: var(--text-xs);
  color: var(--color-ink-light);
}

.form-hint strong {
  color: var(--color-ink);
  font-weight: var(--font-weight-medium);
}

.form-error {
  font-size: var(--text-xs);
  color: var(--color-error);
}

.demo-hint {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-ink-light);
}

.demo-hint strong {
  color: var(--color-accent);
}

/* Button */
.button-full {
  width: 100%;
}

.button-primary {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: white;
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.button-primary:hover:not(:disabled) {
  background: #a6624a;
  transform: translateY(-1px);
}

.button-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.button-lg {
  padding: var(--space-4) var(--space-6);
}

.btn-loading {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.spinner {
  animation: spin 0.8s linear infinite;
}

/* Divider */
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
  color: var(--color-ink-light);
  background: var(--bg-primary);
}

/* Footer */
.form-footer {
  text-align: center;
}

.footer-text {
  font-size: var(--text-base);
  color: var(--color-ink-light);
}

.footer-link {
  color: var(--color-accent);
  font-weight: var(--font-weight-medium);
  text-decoration: none;
  transition: color var(--duration-fast) var(--ease-out);
}

.footer-link:hover {
  color: #a6624a;
}

.privacy-notice {
  margin-top: var(--space-6);
  padding-top: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  font-size: var(--text-xs);
  color: var(--color-ink-light);
  text-align: center;
  line-height: 1.6;
}

.privacy-notice a {
  color: var(--color-ink-light);
  text-decoration: underline;
}

/* Animations */
@keyframes float {
  0%,
  100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(5deg);
  }
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

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Responsive */
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
</style>
