<template>
  <div class="auth-page">
    <!-- Left Brand Section -->
    <div class="auth-brand">
      <div class="brand-content">
        <router-link to="/" class="brand-logo">
          <BookOpen :size="32" />
          <span class="logo-text">iBookFS</span>
        </router-link>

        <h1 class="brand-title">开始数字化之旅</h1>
        <p class="brand-description">加入数千名用户，让纸质书籍焕发新生</p>

        <!-- Features List -->
        <ul class="brand-features">
          <li class="feature-item">
            <Check class="feature-icon" :size="20" />
            <span>无限书籍存储</span>
          </li>
          <li class="feature-item">
            <Check class="feature-icon" :size="20" />
            <span>智能 OCR 识别</span>
          </li>
          <li class="feature-item">
            <Check class="feature-icon" :size="20" />
            <span>多设备同步</span>
          </li>
        </ul>
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
          <h2 class="form-title">创建账号</h2>
          <p class="form-subtitle">填写基本信息，开始使用 iBookFS</p>
        </div>

        <!-- Register Form -->
        <form class="auth-form" @submit.prevent="handleSubmit">
          <!-- Name Fields -->
          <div class="form-row">
            <div class="form-group">
              <label for="firstName" class="form-label">名字</label>
              <input
                id="firstName"
                v-model="firstName"
                type="text"
                class="form-input"
                :class="{ error: errors.firstName, success: !errors.firstName && firstName }"
                placeholder="张"
                :disabled="step === 2"
                required
                autocomplete="given-name"
              />
              <span v-if="errors.firstName" class="form-error">{{ errors.firstName }}</span>
            </div>

            <div class="form-group">
              <label for="lastName" class="form-label">姓氏</label>
              <input
                id="lastName"
                v-model="lastName"
                type="text"
                class="form-input"
                :class="{ error: errors.lastName, success: !errors.lastName && lastName }"
                placeholder="三"
                :disabled="step === 2"
                required
                autocomplete="family-name"
              />
              <span v-if="errors.lastName" class="form-error">{{ errors.lastName }}</span>
            </div>
          </div>

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
            <span class="form-hint">我们将向此邮箱发送验证码</span>
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

          <!-- Terms Checkbox -->
          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="agreeTerms" type="checkbox" required />
              <span class="checkbox-text">
                我已阅读并同意
                <a href="/terms" target="_blank">服务条款</a>
                和
                <a href="/privacy" target="_blank">隐私政策</a>
              </span>
            </label>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="button-primary button-lg button-full"
            :disabled="isLoading || !agreeTerms"
          >
            <span v-if="!isLoading" class="btn-text">{{ step === 1 ? '发送验证码' : '完成注册' }}</span>
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

        <!-- Login Link -->
        <div class="form-footer">
          <p class="footer-text">
            已有账号？
            <router-link to="/login" class="footer-link">立即登录</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { BookOpen, ArrowLeft, Mail, Check, Info } from 'lucide-vue-next'
import { sendCode, register } from '@/api'
import { useToastStore } from '@/stores/toast'

const router = useRouter()
const toastStore = useToastStore()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const code = ref('')
const step = ref(1)
const isLoading = ref(false)
const countdown = ref(0)
const agreeTerms = ref(false)
const errors = ref<Record<string, string>>({})

let countdownTimer: ReturnType<typeof setInterval> | null = null

const validateName = (): boolean => {
  let isValid = true

  if (!firstName.value) {
    errors.value.firstName = '请输入名字'
    isValid = false
  } else {
    errors.value.firstName = ''
  }

  if (!lastName.value) {
    errors.value.lastName = '请输入姓氏'
    isValid = false
  } else {
    errors.value.lastName = ''
  }

  return isValid
}

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
  const response = await sendCode({ email: email.value, type: 'register' })
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
    if (!validateName() || !validateEmail()) return

    isLoading.value = true
    const response = await sendCode({ email: email.value, type: 'register' })
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
    const response = await register({
      firstName: firstName.value,
      lastName: lastName.value,
      email: email.value,
      code: code.value,
    })
    isLoading.value = false

    if (response.success) {
      toastStore.success('注册成功')
      setTimeout(() => {
        router.push('/library')
      }, 500)
    } else {
      errors.value.code = response.message || '注册失败'
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
  color: var(--color-ink);
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

.input-wrapper {
  position: relative;
}

.input-wrapper .form-input {
  padding-right: var(--space-10);
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

/* Checkbox */
.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: var(--space-2);
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type='checkbox'] {
  width: 18px;
  height: 18px;
  margin-top: 2px;
  cursor: pointer;
  accent-color: var(--color-accent);
}

.checkbox-text {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  line-height: 1.5;
}

.checkbox-text a {
  color: var(--color-accent);
  text-decoration: underline;
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
</style>
