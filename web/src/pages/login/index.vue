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
            "iBookFS 让我的书籍管理变得如此简单，支持多种登录方式，体验非常流畅。"
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
          <p class="form-subtitle">选择您喜欢的登录方式</p>
        </div>

        <!-- Tab Switcher (only show if both email and OAuth are enabled) -->
        <div v-if="authStore.isEmailEnabled && authStore.isOAuthEnabled" class="auth-tabs">
          <button
            :class="['auth-tab', { active: activeTab === 'email' }]"
            @click="activeTab = 'email'"
          >
            <Mail :size="18" />
            邮箱登录
          </button>
          <button
            :class="['auth-tab', { active: activeTab === 'oauth' }]"
            @click="activeTab = 'oauth'"
          >
            <Link :size="18" />
            第三方登录
          </button>
        </div>

        <!-- Email Login Form -->
        <div v-show="activeTab === 'email'" class="auth-form-content">
          <form class="auth-form" @submit.prevent="handleSubmit">
            <!-- Email Input -->
            <div class="form-group">
              <label for="email" class="form-label">邮箱地址</label>
              <div class="input-wrapper">
                <input
                  id="email"
                  ref="emailInput"
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

            <!-- Verification Code Input (6-digit split) -->
            <div v-show="step === 2" class="form-group">
              <label class="form-label">验证码</label>
              <div class="code-digits-container">
                <input
                  v-for="(_digit, index) in 6"
                  :key="index"
                  :ref="el => codeInputs[index] = el as HTMLInputElement"
                  v-model="codeDigits[index]"
                  type="text"
                  inputmode="numeric"
                  class="code-digit-input"
                  :class="{ error: errors.code, success: !errors.code && codeDigits.join('').length === 6 }"
                  maxlength="1"
                  @input="handleCodeInput(index, $event)"
                  @keydown="handleCodeKeydown(index, $event)"
                  @paste="handleCodePaste"
                  autocomplete="one-time-code"
                />
              </div>
              <span class="form-hint">
                验证码已发送至 <strong>{{ email }}</strong>
              </span>
              <span v-if="errors.code" class="form-error">{{ errors.code }}</span>
            </div>

            <!-- Demo Hint -->
            <div v-if="step === 2 && !authStore.isEmailEnabled" class="demo-hint">
              <Info :size="16" />
              <span>演示模式：验证码为 <strong>123456</strong></span>
            </div>

            <!-- Submit Button / Resend Button -->
            <button
              v-if="step === 1"
              type="submit"
              class="button-primary button-lg button-full"
              :disabled="isLoading"
            >
              <span v-if="!isLoading">发送验证码</span>
              <span v-else class="btn-loading">
                <svg class="spinner" width="20" height="20" viewBox="0 0 20 20">
                  <circle cx="10" cy="10" r="8" stroke="currentColor" stroke-width="2" fill="none" opacity="0.25" />
                  <path d="M10 2a8 8 0 018 8" stroke="currentColor" stroke-width="2" fill="none" />
                </svg>
                处理中...
              </span>
            </button>

            <!-- Step 2: Resend button replaces submit button -->
            <button
              v-if="step === 2 && !isSubmitting"
              type="button"
              class="resend-button button-lg button-full"
              :disabled="countdown > 0 || isLoading"
              @click="handleResend"
            >
              {{ countdown > 0 ? `${countdown}s 后重新发送` : '重新发送验证码' }}
            </button>

            <!-- Loading state during login -->
            <div v-if="step === 2 && isSubmitting" class="submitting-state">
              <svg class="spinner" width="24" height="24" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" fill="none" opacity="0.25" />
                <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="3" fill="none" />
              </svg>
              <span>验证中...</span>
            </div>
          </form>
        </div>

        <!-- OAuth Login -->
        <div v-show="activeTab === 'oauth'" class="oauth-content">
          <p class="oauth-description">使用您的第三方账号快速登录</p>

          <div class="oauth-buttons">
            <OAuthButton
              v-for="provider in oauthProviders"
              :key="provider.id"
              :provider="provider.id"
              :is-loading="oauthLoading[provider.id]"
              :is-linked="linkedProviders.includes(provider.id)"
              size="lg"
              variant="outlined"
              @click="handleOAuthLogin(provider.id)"
            >
              使用 {{ provider.name }} 登录
            </OAuthButton>
          </div>

          <div class="oauth-hint">
            <Shield :size="16" />
            <span>我们尊重您的隐私，不会获取不必要的权限</span>
          </div>
        </div>

        <!-- Divider -->
        <div class="divider">
          <span>或</span>
        </div>

        <!-- Mode Switch -->
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

    <!-- OAuth Callback Modal -->
    <div v-if="showOAuthModal" class="oauth-modal-overlay" @click.self="closeOAuthModal">
      <div class="oauth-modal">
        <div class="oauth-modal-content">
          <div class="oauth-modal-icon">
            <svg class="spinner" width="48" height="48" viewBox="0 0 48 48">
              <circle cx="24" cy="24" r="20" stroke="currentColor" stroke-width="3" fill="none" opacity="0.25" />
              <path d="M24 4a20 20 0 0120 20" stroke="currentColor" stroke-width="3" fill="none" />
            </svg>
          </div>
          <h3 class="oauth-modal-title">正在连接...</h3>
          <p class="oauth-modal-description">请在新窗口中完成授权</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { BookOpen, ArrowLeft, Mail, Link, Info, Shield } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import OAuthButton from '@/components/ui/OAuthButton.vue'
import type { OAuthProvider } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

// Form state
const activeTab = ref<'email' | 'oauth'>('email')
const step = ref(1)

// Computed - determine initial active tab based on config
const initialTab = computed<'email' | 'oauth'>(() => {
  if (!authStore.isEmailEnabled && authStore.isOAuthEnabled) {
    return 'oauth'
  }
  return 'email'
})

// Auto-select initial tab when config loads
watch(initialTab, (newTab) => {
  if (activeTab.value !== newTab) {
    activeTab.value = newTab
  }
})

const isLoading = ref(false)
const isSubmitting = ref(false)
const countdown = ref(0)
const errors = ref<Record<string, string>>({})

// Form fields
const email = ref('')
const codeDigits = ref(['', '', '', '', '', ''])
const codeInputs = ref<(HTMLInputElement | null)[]>([])

// Refs
const emailInput = ref<HTMLInputElement>()

// OAuth state
const allOAuthProviders = [
  { id: 'github' as OAuthProvider, name: 'GitHub' },
  { id: 'gitee' as OAuthProvider, name: 'Gitee' },
  { id: 'icloud' as OAuthProvider, name: 'iCloud' },
  { id: 'google' as OAuthProvider, name: 'Google' },
  { id: 'wechat' as OAuthProvider, name: 'WeChat' },
]
const oauthProviders = computed(() =>
  allOAuthProviders.filter((p) => authStore.isProviderConfigured(p.id))
)
const oauthLoading = ref<Record<OAuthProvider, boolean>>({
  github: false,
  gitee: false,
  icloud: false,
  google: false,
  wechat: false,
})
const showOAuthModal = ref(false)
const linkedProviders = computed(() => authStore.linkedProviders)

let countdownTimer: ReturnType<typeof setInterval> | null = null
let oauthWindow: Window | null = null

// Methods
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
  // Use codeDigits directly instead of code.value (watch is async)
  const codeValue = codeDigits.value.join('')
  if (codeValue.length !== 6) {
    errors.value.code = '请输入完整的 6 位验证码'
    return false
  }
  if (!/^\d{6}$/.test(codeValue)) {
    errors.value.code = '验证码必须是 6 位数字'
    return false
  }
  errors.value.code = ''
  return true
}

const startCountdown = () => {
  countdown.value = 60
  if (countdownTimer) clearInterval(countdownTimer)
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (countdownTimer) clearInterval(countdownTimer)
    }
  }, 1000)
}

// Handle code input in each digit box
const handleCodeInput = async (index: number, event: Event) => {
  const input = event.target as HTMLInputElement
  const value = input.value

  console.log('[LoginPage] handleCodeInput:', { index, value, codeDigits: codeDigits.value })

  // Only allow digits
  if (!/^\d*$/.test(value)) {
    input.value = ''
    return
  }

  // Take only the last character if multiple entered
  codeDigits.value[index] = value.slice(-1)

  // Auto-focus next input
  if (value && index < 5) {
    const nextInput = codeInputs.value[index + 1]
    if (nextInput) {
      nextInput.focus()
    }
  }

  // Wait for Vue to process reactivity
  await nextTick()

  // Auto-submit when all 6 digits are entered
  const allDigitsEntered = codeDigits.value.every(d => d !== '')
  console.log('[LoginPage] allDigitsEntered:', allDigitsEntered, 'codeDigits:', codeDigits.value)
  if (allDigitsEntered) {
    errors.value.code = ''
    submitLogin()
  }
}

// Handle keydown in code input boxes
const handleCodeKeydown = (index: number, event: KeyboardEvent) => {
  const input = event.target as HTMLInputElement

  // Handle Backspace: delete current digit and focus previous
  if (event.key === 'Backspace') {
    if (!input.value && index > 0) {
      const prevInput = codeInputs.value[index - 1]
      if (prevInput) {
        prevInput.focus()
      }
    } else {
      codeDigits.value[index] = ''
    }
  }

  // Handle Arrow keys for navigation
  if (event.key === 'ArrowLeft' && index > 0) {
    event.preventDefault()
    const prevInput = codeInputs.value[index - 1]
    if (prevInput) prevInput.focus()
  }
  if (event.key === 'ArrowRight' && index < 5) {
    event.preventDefault()
    const nextInput = codeInputs.value[index + 1]
    if (nextInput) nextInput.focus()
  }
}

// Handle paste event
const handleCodePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text') || ''

  // Extract only digits from pasted content
  const digits = pastedData.replace(/\D/g, '').slice(0, 6)

  // Fill the code digits
  for (let i = 0; i < digits.length; i++) {
    codeDigits.value[i] = digits[i]
  }

  // Focus the appropriate input
  const focusIndex = Math.min(digits.length, 5)
  const targetInput = codeInputs.value[focusIndex]
  if (targetInput) {
    targetInput.focus()
  }

  // Auto-submit if we have all 6 digits
  const allDigitsEntered = codeDigits.value.every(d => d !== '')
  if (allDigitsEntered) {
    submitLogin()
  }
}

const submitLogin = async () => {
  if (!validateCode()) return

  // Use codeDigits directly instead of code.value (watch is async)
  const codeValue = codeDigits.value.join('')

  isSubmitting.value = true
  const success = await authStore.login({ email: email.value, code: codeValue })
  isSubmitting.value = false

  if (success) {
    toastStore.success('登录成功')
    setTimeout(() => {
      router.push('/library')
    }, 500)
  } else {
    errors.value.code = authStore.error || '登录失败'
    // Clear code on error
    codeDigits.value = ['', '', '', '', '', '']
    // Focus first input
    nextTick(() => {
      const firstInput = codeInputs.value[0]
      if (firstInput) firstInput.focus()
    })
  }
}

const handleResend = async () => {
  const response = await authStore.sendCode({ email: email.value, type: 'login' })
  if (response) {
    toastStore.success('验证码已重新发送')
    startCountdown()
    // Clear code inputs and focus first one
    codeDigits.value = ['', '', '', '', '', '']
    errors.value.code = ''
    nextTick(() => {
      const firstInput = codeInputs.value[0]
      if (firstInput) firstInput.focus()
    })
  } else {
    toastStore.error(authStore.error || '发送失败')
  }
}

const handleSubmit = async () => {
  errors.value = {}

  if (step.value === 1) {
    if (!validateEmail()) return

    isLoading.value = true
    const success = await authStore.sendCode({ email: email.value, type: 'login' })
    isLoading.value = false

    if (success) {
      step.value = 2
      startCountdown()
      toastStore.success('验证码已发送')
      // Focus first code input
      nextTick(() => {
        const firstInput = codeInputs.value[0]
        if (firstInput) firstInput.focus()
      })
    } else {
      errors.value.email = authStore.error || '发送失败'
    }
  }
  // Step 2 is handled automatically by code input
}

const handleOAuthLogin = async (provider: OAuthProvider) => {
  oauthLoading.value[provider] = true

  try {
    const redirectUri = `${window.location.origin}/auth/callback`
    const result = await authStore.getOAuthAuthorizeUrl(provider, redirectUri)

    if (result) {
      sessionStorage.setItem('oauth_state', result.state)
      sessionStorage.setItem('oauth_provider', provider)

      showOAuthModal.value = true
      oauthWindow = window.open(result.url, 'oauth', 'width=600,height=700,scrollbars=yes')

      const checkPopup = setInterval(() => {
        if (oauthWindow?.closed) {
          clearInterval(checkPopup)
          closeOAuthModal()
          oauthLoading.value[provider] = false
        }
      }, 500)
    } else {
      toastStore.error(authStore.error || '获取授权链接失败')
    }
  } finally {
    if (!showOAuthModal.value) {
      oauthLoading.value[provider] = false
    }
  }
}

const closeOAuthModal = () => {
  showOAuthModal.value = false
  oauthWindow = null
}

// Handle OAuth callback from popup
const handleOAuthCallback = async (provider: OAuthProvider, code: string, state: string) => {
  const savedState = sessionStorage.getItem('oauth_state')
  const savedProvider = sessionStorage.getItem('oauth_provider')

  if (state !== savedState || provider !== savedProvider) {
    toastStore.error('授权验证失败')
    return
  }

  oauthLoading.value[provider] = true

  try {
    const success = await authStore.oauthCallback(provider, code, state)

    if (success) {
      toastStore.success('登录成功')
      closeOAuthModal()
      setTimeout(() => {
        router.push('/library')
      }, 500)
    } else {
      toastStore.error(authStore.error || 'OAuth 登录失败')
    }
  } finally {
    oauthLoading.value[provider] = false
    closeOAuthModal()
  }

  sessionStorage.removeItem('oauth_state')
  sessionStorage.removeItem('oauth_provider')
}

// Listen for OAuth callback message
onMounted(async () => {
  // Initialize auth store to check if user is already logged in
  if (!authStore.isAuthenticated) {
    await authStore.initialize()
  }
  // If already logged in, redirect to library
  if (authStore.isAuthenticated) {
    router.push('/library')
    return
  }

  await authStore.loadConfig()

  window.addEventListener('message', (event) => {
    if (event.origin !== window.location.origin) return

    if (event.data.type === 'oauth_callback') {
      handleOAuthCallback(event.data.provider, event.data.code, event.data.state)
    }
  })
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
  if (oauthWindow && !oauthWindow.closed) {
    oauthWindow.close()
  }
})
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
  margin-bottom: var(--space-6);
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

/* Tabs */
.auth-tabs {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
  background: var(--bg-secondary);
  padding: var(--space-1);
  border-radius: var(--radius-md);
}

.auth-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-3);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-ink-light);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.auth-tab:hover {
  color: var(--color-ink);
}

.auth-tab.active {
  background: var(--bg-elevated);
  color: var(--color-ink);
  box-shadow: var(--shadow-xs);
}

/* Form Styles */
.auth-form-content {
  animation: fadeIn 0.3s var(--ease-out);
}

.auth-form {
  display: flex;
  flex-direction: column;
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

/* Code Digits Input - 6 separate boxes */
.code-digits-container {
  display: flex;
  gap: var(--space-3);
  justify-content: space-between;
}

.code-digit-input {
  width: 56px;
  height: 56px;
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  line-height: 56px;
  text-align: center;
  color: var(--color-ink);
  background: var(--bg-elevated);
  border: 2px solid var(--border-medium);
  border-radius: var(--radius-lg);
  transition: all var(--duration-fast) var(--ease-out);
  font-family: var(--font-mono);
  letter-spacing: 0;
}

.code-digit-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(193, 123, 92, 0.2);
  transform: translateY(-2px);
}

.code-digit-input.error {
  border-color: var(--color-error);
}

.code-digit-input.error:focus {
  box-shadow: 0 0 0 3px rgba(184, 92, 92, 0.2);
}

.code-digit-input.success {
  border-color: var(--color-success);
}

/* Resend Button (replaces submit button in step 2) */
.resend-button {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-accent);
  background: var(--bg-secondary);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.resend-button:hover:not(:disabled) {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
  color: white;
  transform: translateY(-1px);
}

.resend-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Submitting State */
.submitting-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-accent);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
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

/* OAuth Content */
.oauth-content {
  animation: fadeIn 0.3s var(--ease-out);
}

.oauth-description {
  text-align: center;
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  margin-bottom: var(--space-6);
}

.oauth-buttons {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.oauth-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-3);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  color: var(--color-ink-light);
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

/* OAuth Modal */
.oauth-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  animation: fadeIn 0.2s var(--ease-out);
}

.oauth-modal {
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
  max-width: 360px;
  width: 90%;
  box-shadow: var(--shadow-2xl);
  animation: scaleIn 0.3s var(--ease-out);
}

.oauth-modal-content {
  text-align: center;
}

.oauth-modal-icon {
  display: flex;
  justify-content: center;
  margin-bottom: var(--space-5);
  color: var(--color-accent);
}

.oauth-modal-title {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
  margin-bottom: var(--space-2);
}

.oauth-modal-description {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
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

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
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

  .code-digits-container {
    gap: var(--space-2);
  }

  .code-digit-input {
    width: 52px;
    height: 52px;
    font-size: var(--text-xl);
    line-height: 52px;
  }
}
</style>
