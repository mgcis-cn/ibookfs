<template>
  <div class="binding-page">
    <!-- Header -->
    <div class="binding-header">
      <div class="header-container">
        <router-link to="/settings" class="back-link">
          <ArrowLeft :size="20" />
          返回设置
        </router-link>

        <div class="header-content">
          <div class="header-icon">
            <Link :size="32" />
          </div>
          <div>
            <h1 class="header-title">账号绑定</h1>
            <p class="header-description">管理您的登录方式，支持多种方式绑定</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="binding-content">
      <div class="content-container">
        <!-- Primary Account Section -->
        <section class="binding-section">
          <h2 class="section-title">
            <Shield :size="20" />
            主账号
          </h2>
          <p class="section-description">设置您的主要登录方式</p>

          <div class="accounts-list">
            <!-- Email Account -->
            <div
              v-if="hasEmail"
              :class="['account-card', { primary: isPrimary('email') }]"
            >
              <div class="account-info">
                <div class="account-icon email-icon">
                  <Mail :size="24" />
                </div>
                <div class="account-details">
                  <div class="account-name">邮箱登录</div>
                  <div class="account-email">{{ user?.email }}</div>
                  <div v-if="isPrimary('email')" class="account-badge primary">主账号</div>
                </div>
              </div>
              <div class="account-actions">
                <button
                  v-if="!isPrimary('email')"
                  class="action-btn set-primary-btn"
                  @click="handleSetPrimary('email')"
                >
                  设为主账号
                </button>
              </div>
            </div>

            <!-- OAuth Accounts -->
            <div
              v-for="account in oauthAccounts"
              :key="account.provider"
              :class="['account-card', { primary: account.isPrimary }]"
            >
              <div class="account-info">
                <div :class="['account-icon', `${account.provider}-icon`]">
                  <component :is="getProviderIcon(account.provider)" :size="24" />
                </div>
                <div class="account-details">
                  <div class="account-name">{{ getProviderName(account.provider) }}</div>
                  <div class="account-username">{{ account.providerUsername }}</div>
                  <div v-if="account.isPrimary" class="account-badge primary">主账号</div>
                  <div class="account-meta">
                    <span>绑定于 {{ formatDate(account.linkedAt) }}</span>
                    <span v-if="account.lastUsedAt">
                      · 最后使用于 {{ formatDate(account.lastUsedAt) }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="account-actions">
                <button
                  v-if="!account.isPrimary"
                  class="action-btn set-primary-btn"
                  @click="handleSetPrimary(account.provider, account.providerUserId)"
                >
                  设为主账号
                </button>
                <button
                  class="action-btn unlink-btn"
                  @click="handleUnlink(account.provider, account.providerUserId)"
                >
                  解绑
                </button>
              </div>
            </div>
          </div>
        </section>

        <!-- Link New Account Section -->
        <section class="binding-section">
          <h2 class="section-title">
            <PlusCircle :size="20" />
            绑定新账号
          </h2>
          <p class="section-description">添加更多登录方式，让账号更安全</p>

          <!-- Available Providers -->
          <div class="providers-grid">
            <div
              v-for="provider in availableProviders"
              :key="provider.id"
              :class="['provider-card', { linked: isProviderLinked(provider.id) }]"
            >
              <div class="provider-icon">
                <component :is="getProviderIcon(provider.id)" :size="32" />
              </div>
              <div class="provider-info">
                <div class="provider-name">{{ provider.name }}</div>
                <div v-if="isProviderLinked(provider.id)" class="provider-status linked">
                  <CheckCircle :size="16" />
                  已绑定
                </div>
                <div v-else class="provider-status">
                  点击绑定
                </div>
              </div>
              <div class="provider-action">
                <OAuthButton
                  v-if="!isProviderLinked(provider.id)"
                  :provider="provider.id"
                  :is-loading="linkLoading[provider.id]"
                  variant="default"
                  @click="handleLinkProvider(provider.id)"
                >
                  绑定
                </OAuthButton>
                <div v-else class="linked-badge">
                  <CheckCircle :size="18" />
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Bind Email Section (for OAuth-only users) -->
        <section v-if="!hasEmail" class="binding-section bind-email-section">
          <h2 class="section-title">
            <Mail :size="20" />
            绑定邮箱
          </h2>
          <p class="section-description">
            添加邮箱作为登录方式，方便您随时访问账号
          </p>

          <form class="bind-email-form" @submit.prevent="handleBindEmail">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">邮箱地址</label>
                <input
                  v-model="bindEmailForm.email"
                  type="email"
                  class="form-input"
                  placeholder="your@email.com"
                  :disabled="bindEmailForm.step === 2"
                />
              </div>

              <div v-if="bindEmailForm.step === 2" class="form-group">
                <label class="form-label">验证码</label>
                <div class="code-input-wrapper">
                  <input
                    v-model="bindEmailForm.code"
                    type="text"
                    class="form-input code-input"
                    placeholder="6 位验证码"
                    maxlength="6"
                  />
                  <button
                    type="button"
                    class="resend-btn"
                    :disabled="bindEmailForm.countdown > 0"
                    @click="handleSendBindCode"
                  >
                    {{ bindEmailForm.countdown > 0 ? `${bindEmailForm.countdown}s` : '发送验证码' }}
                  </button>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button
                v-if="bindEmailForm.step === 1"
                type="button"
                class="button-secondary"
                @click="handleSendBindCode"
              >
                发送验证码
              </button>
              <button type="submit" class="button-primary" :disabled="bindEmailForm.isLoading">
                {{ bindEmailForm.isLoading ? '处理中...' : '绑定邮箱' }}
              </button>
            </div>
          </form>
        </section>
      </div>
    </div>

    <!-- OAuth Modal -->
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
          <button class="button-secondary" @click="closeOAuthModal">取消</button>
        </div>
      </div>
    </div>

    <!-- Confirm Modal -->
    <div v-if="showConfirmModal" class="confirm-modal-overlay" @click.self="closeConfirmModal">
      <div class="confirm-modal">
        <div class="confirm-modal-header">
          <h3 class="confirm-modal-title">确认解绑</h3>
          <button class="close-btn" @click="closeConfirmModal">
            <X :size="20" />
          </button>
        </div>
        <div class="confirm-modal-body">
          <p class="confirm-modal-description">
            确定要解绑 <strong>{{ getProviderName(confirmModal.provider) }}</strong> 账号吗？
          </p>
          <p v-if="confirmModal.isPrimary" class="confirm-modal-warning">
            这是您的主账号，解绑后需要设置其他账号为主账号。
          </p>
        </div>
        <div class="confirm-modal-footer">
          <button class="button-secondary" @click="closeConfirmModal">取消</button>
          <button class="button-primary danger" @click="confirmUnlink">
            确认解绑
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import {
  ArrowLeft,
  Link,
  Shield,
  PlusCircle,
  Mail,
  CheckCircle,
  X,
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import OAuthButton from '@/components/ui/OAuthButton.vue'
import type { OAuthProvider, LinkedAccount } from '@/types'

// Provider icons (inline SVGs)
const GithubIcon = {
  template: h('svg', {
    viewBox: '0 0 24 24',
    fill: 'currentColor'
  }, [
    h('path', {
      d: 'M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z'
    })
  ])
}

const GiteeIcon = {
  template: h('svg', {
    viewBox: '0 0 24 24',
    fill: 'currentColor'
  }, [
    h('path', {
      d: 'M11.984 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.016 0zm6.09 5.333c.328 0 .593.266.592.593v1.482a.594.594 0 0 1-.593.592H9.777c-.982 0-1.778.796-1.778 1.778v5.63c0 .327.266.592.593.592h5.63c.982 0 1.778-.796 1.778-1.778v-.296a.593.593 0 0 0-.592-.593h-4.037a.594.594 0 0 1-.592-.593v-1.482a.593.593 0 0 1 .593-.592h6.815c.327 0 .593.265.593.592v3.408a4 4 0 0 1-4 4H5.926a.593.593 0 0 1-.593-.593V9.778a4.444 4.444 0 0 1 4.445-4.444h8.296z'
    })
  ])
}

const IcloudIcon = {
  template: h('svg', {
    viewBox: '0 0 24 24',
    fill: 'currentColor'
  }, [
    h('path', {
      d: 'M17.6 9.48c1.03-.68 1.73-1.83 1.73-3.18 0-2.12-1.74-3.85-3.88-3.85-.65 0-1.26.16-1.79.44-.68-2.15-2.71-3.72-5.1-3.72-2.93 0-5.31 2.35-5.31 5.25 0 .28.03.55.08.82C1.23 6.77 0 8.45 0 10.42c0 2.43 1.99 4.41 4.44 4.41h13.12c2.45 0 4.44-1.98 4.44-4.41 0-2.12-1.51-3.9-3.52-4.28-.28-.07-.57-.11-.88-.11zM12.5 15l-1.5-3-1.5 3h3zm-2.5 4h1v-2h-1v2zm3 0h1v-2h-1v2z'
    })
  ])
}

const GoogleIcon = {
  template: h('svg', {
    viewBox: '0 0 24 24',
    fill: 'currentColor'
  }, [
    h('path', {
      d: 'M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z'
    }),
    h('path', {
      d: 'M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z'
    }),
    h('path', {
      d: 'M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z'
    }),
    h('path', {
      d: 'M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z'
    })
  ])
}

const WechatIcon = {
  template: h('svg', {
    viewBox: '0 0 24 24',
    fill: 'currentColor'
  }, [
    h('path', {
      d: 'M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z'
    })
  ])
}

const authStore = useAuthStore()
const toastStore = useToastStore()

const user = computed(() => authStore.user)
const linkedAccounts = computed(() => authStore.linkedAccounts)
const hasEmail = computed(() => authStore.hasEmail)
const oauthAccounts = computed(() =>
  linkedAccounts.value.filter((a) => a.provider !== 'email') as LinkedAccount[]
)

// Available providers
const availableProviders = [
  { id: 'github' as OAuthProvider, name: 'GitHub' },
  { id: 'gitee' as OAuthProvider, name: 'Gitee' },
  { id: 'icloud' as OAuthProvider, name: 'iCloud' },
  { id: 'google' as OAuthProvider, name: 'Google' },
  { id: 'wechat' as OAuthProvider, name: 'WeChat' },
]

// Link loading state
const linkLoading = ref<Record<OAuthProvider, boolean>>({
  github: false,
  gitee: false,
  icloud: false,
  google: false,
  wechat: false,
})

// OAuth modal
const showOAuthModal = ref(false)
const linkingProvider = ref<OAuthProvider | null>(null)
let oauthWindow: Window | null = null

// Confirm modal
const showConfirmModal = ref(false)
const confirmModal = ref<{
  provider: OAuthProvider
  providerUserId?: string
  isPrimary: boolean
}>({
  provider: 'github',
  isPrimary: false,
})

// Bind email form
const bindEmailForm = ref({
  email: '',
  code: '',
  step: 1,
  isLoading: false,
  countdown: 0,
})

let bindEmailTimer: ReturnType<typeof setInterval> | null = null

// Methods
const getProviderIcon = (provider: OAuthProvider | 'email') => {
  if (provider === 'email') return Mail
  const icons: Record<OAuthProvider, unknown> = {
    github: GithubIcon,
    gitee: GiteeIcon,
    icloud: IcloudIcon,
    google: GoogleIcon,
    wechat: WechatIcon,
  }
  return (icons[provider] as unknown) ?? GithubIcon
}

const getProviderName = (provider: OAuthProvider | 'email'): string => {
  if (provider === 'email') return '邮箱登录'
  const names: Record<OAuthProvider, string> = {
    github: 'GitHub',
    gitee: 'Gitee',
    icloud: 'iCloud',
    google: 'Google',
    wechat: '微信',
  }
  return names[provider] || provider
}

const isPrimary = (provider: OAuthProvider | 'email'): boolean => {
  return linkedAccounts.value.some((a) => a.provider === provider && a.isPrimary)
}

const isProviderLinked = (provider: OAuthProvider): boolean => {
  return linkedAccounts.value.some((a) => a.provider === provider)
}

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  return date.toLocaleDateString('zh-CN')
}

const handleSetPrimary = async (provider: OAuthProvider | 'email', providerUserId?: string) => {
  const success = await authStore.setPrimaryAccount({
    provider,
    providerUserId,
  })

  if (success) {
    toastStore.success('主账号设置成功')
  } else {
    toastStore.error(authStore.error || '设置失败')
  }
}

const handleUnlink = (provider: OAuthProvider | 'email', providerUserId?: string) => {
  if (provider === 'email') {
    toastStore.error('无法解绑邮箱账号')
    return
  }
  const account = linkedAccounts.value.find((a) => a.provider === provider)
  confirmModal.value = {
    provider,
    providerUserId,
    isPrimary: account?.isPrimary || false,
  }
  showConfirmModal.value = true
}

const confirmUnlink = async () => {
  const { provider, providerUserId } = confirmModal.value

  const success = await authStore.unlinkOAuth({
    provider,
    providerUserId,
  })

  if (success) {
    toastStore.success('解绑成功')
    closeConfirmModal()
  } else {
    toastStore.error(authStore.error || '解绑失败')
  }
}

const closeConfirmModal = () => {
  showConfirmModal.value = false
}

const handleLinkProvider = async (provider: OAuthProvider) => {
  linkLoading.value[provider] = true
  linkingProvider.value = provider

  try {
    const redirectUri = `${window.location.origin}/settings/accounts/link`
    const result = await authStore.getOAuthAuthorizeUrl(provider, redirectUri)

    if (result) {
      sessionStorage.setItem('oauth_action', 'link')
      sessionStorage.setItem('oauth_state', result.state)
      sessionStorage.setItem('oauth_provider', provider)

      showOAuthModal.value = true
      oauthWindow = window.open(result.url, 'oauth_link', 'width=600,height=700,scrollbars=yes')

      const checkPopup = setInterval(() => {
        if (oauthWindow?.closed) {
          clearInterval(checkPopup)
          closeOAuthModal()
          linkLoading.value[provider] = false
        }
      }, 500)
    } else {
      toastStore.error(authStore.error || '获取授权链接失败')
      linkLoading.value[provider] = false
    }
  } catch {
    toastStore.error('操作失败')
    linkLoading.value[provider] = false
  }
}

const closeOAuthModal = () => {
  showOAuthModal.value = false
  if (linkingProvider.value) {
    linkLoading.value[linkingProvider.value] = false
    linkingProvider.value = null
  }
  if (oauthWindow && !oauthWindow.closed) {
    oauthWindow.close()
  }
  oauthWindow = null
}

const handleSendBindCode = async () => {
  if (!bindEmailForm.value.email) {
    toastStore.error('请输入邮箱地址')
    return
  }

  const success = await authStore.sendCode({
    email: bindEmailForm.value.email,
    type: 'bind_email',
  })

  if (success) {
    toastStore.success('验证码已发送')
    bindEmailForm.value.step = 2
    bindEmailForm.value.countdown = 60
    bindEmailTimer = setInterval(() => {
      bindEmailForm.value.countdown--
      if (bindEmailForm.value.countdown <= 0) {
        if (bindEmailTimer) clearInterval(bindEmailTimer)
      }
    }, 1000)
  } else {
    toastStore.error(authStore.error || '发送失败')
  }
}

const handleBindEmail = async () => {
  if (!bindEmailForm.value.code) {
    toastStore.error('请输入验证码')
    return
  }

  bindEmailForm.value.isLoading = true

  const success = await authStore.bindEmail({
    email: bindEmailForm.value.email,
    code: bindEmailForm.value.code,
  })

  bindEmailForm.value.isLoading = false

  if (success) {
    toastStore.success('邮箱绑定成功')
    bindEmailForm.value = {
      email: '',
      code: '',
      step: 1,
      isLoading: false,
      countdown: 0,
    }
  } else {
    toastStore.error(authStore.error || '绑定失败')
  }
}

// Handle OAuth callback from popup
const handleOAuthCallback = async (provider: OAuthProvider, code: string, state: string) => {
  const savedAction = sessionStorage.getItem('oauth_action')
  const savedState = sessionStorage.getItem('oauth_state')
  const savedProvider = sessionStorage.getItem('oauth_provider')

  if (state !== savedState || provider !== savedProvider || savedAction !== 'link') {
    toastStore.error('授权验证失败')
    return
  }

  if (linkingProvider.value) {
    linkLoading.value[linkingProvider.value] = true
  }

  try {
    const success = await authStore.linkOAuth({
      provider,
      code,
      state,
    })

    if (success) {
      toastStore.success('账号绑定成功')
    } else {
      toastStore.error(authStore.error || '绑定失败')
    }
  } finally {
    if (linkingProvider.value) {
      linkLoading.value[linkingProvider.value] = false
    }
    closeOAuthModal()
  }

  sessionStorage.removeItem('oauth_action')
  sessionStorage.removeItem('oauth_state')
  sessionStorage.removeItem('oauth_provider')
}

// Listen for OAuth callback message
onMounted(() => {
  window.addEventListener('message', (event) => {
    if (event.origin !== window.location.origin) return

    if (event.data.type === 'oauth_callback') {
      handleOAuthCallback(event.data.provider, event.data.code, event.data.state)
    }
  })
})

onUnmounted(() => {
  if (bindEmailTimer) clearInterval(bindEmailTimer)
  if (oauthWindow && !oauthWindow.closed) {
    oauthWindow.close()
  }
})
</script>

<style scoped>
.binding-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* Header */
.binding-header {
  background: var(--bg-elevated);
  border-bottom: 1px solid var(--border-subtle);
  padding: var(--space-6) 0;
}

.header-container {
  max-width: var(--container-lg);
  margin: 0 auto;
  padding: 0 var(--space-6);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  text-decoration: none;
  margin-bottom: var(--space-6);
  transition: color var(--duration-fast) var(--ease-out);
}

.back-link:hover {
  color: var(--color-accent);
}

.header-content {
  display: flex;
  align-items: center;
  gap: var(--space-5);
}

.header-icon {
  width: 56px;
  height: 56px;
  background: var(--color-accent-light);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.header-title {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-ink);
  margin-bottom: var(--space-1);
}

.header-description {
  font-size: var(--text-base);
  color: var(--color-ink-light);
}

/* Main Content */
.binding-content {
  padding: var(--space-8) 0;
}

.content-container {
  max-width: var(--container-lg);
  margin: 0 auto;
  padding: 0 var(--space-6);
}

/* Section */
.binding-section {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  border: 1px solid var(--border-subtle);
}

.section-title {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
  margin-bottom: var(--space-2);
}

.section-description {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  margin-bottom: var(--space-6);
}

/* Account Cards */
.accounts-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.account-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border: 1.5px solid var(--border-subtle);
  transition: all var(--duration-fast) var(--ease-out);
}

.account-card:hover {
  border-color: var(--border-medium);
}

.account-card.primary {
  border-color: var(--color-accent);
  background: rgba(193, 123, 92, 0.05);
}

.account-info {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.account-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.account-icon.email-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.account-icon.github-icon {
  background: #333;
}

.account-icon.gitee-icon {
  background: #c71d23;
}

.account-icon.icloud-icon {
  background: #007aff;
}

.account-icon.google-icon {
  background: #4285f4;
}

.account-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.account-name {
  font-size: var(--text-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
}

.account-email,
.account-username {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
}

.account-meta {
  font-size: var(--text-xs);
  color: var(--color-ink-lighter);
}

.account-badge {
  display: inline-flex;
  align-items: center;
  padding: var(--space-1) var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-sm);
}

.account-badge.primary {
  background: var(--color-accent);
  color: white;
}

.account-actions {
  display: flex;
  gap: var(--space-2);
}

.action-btn {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.set-primary-btn {
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  color: var(--color-ink);
}

.set-primary-btn:hover {
  background: var(--bg-secondary);
  border-color: var(--color-accent);
}

.unlink-btn {
  background: transparent;
  border: 1.5px solid var(--border-medium);
  color: var(--color-error);
}

.unlink-btn:hover {
  background: rgba(184, 92, 92, 0.1);
  border-color: var(--color-error);
}

/* Providers Grid */
.providers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--space-4);
}

.provider-card {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border: 1.5px solid var(--border-subtle);
  transition: all var(--duration-fast) var(--ease-out);
}

.provider-card:hover {
  border-color: var(--border-medium);
}

.provider-card.linked {
  border-color: var(--color-success);
  background: rgba(122, 155, 118, 0.05);
}

.provider-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-ink-light);
}

.provider-info {
  flex: 1;
}

.provider-name {
  font-size: var(--text-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
}

.provider-status {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  margin-top: var(--space-1);
}

.provider-status.linked {
  color: var(--color-success);
}

.provider-action {
  display: flex;
  align-items: center;
}

.linked-badge {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-success);
}

/* Bind Email Section */
.bind-email-section {
  border-color: var(--color-accent-light);
  background: linear-gradient(135deg, rgba(193, 123, 92, 0.03) 0%, var(--bg-elevated) 100%);
}

.bind-email-form {
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
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  color: var(--color-ink);
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(193, 123, 92, 0.2);
}

.code-input-wrapper {
  display: flex;
  gap: var(--space-2);
}

.code-input {
  font-family: var(--font-mono);
  text-align: center;
  flex: 1;
}

.resend-btn {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-accent);
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  white-space: nowrap;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.resend-btn:hover:not(:disabled) {
  background: var(--bg-secondary);
  border-color: var(--color-accent);
}

.resend-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: var(--space-3);
  justify-content: flex-end;
}

/* Buttons */
.button-primary,
.button-secondary {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.button-primary {
  background: var(--color-accent);
  color: white;
  border: none;
}

.button-primary:hover:not(:disabled) {
  background: #a6624a;
}

.button-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.button-primary.danger {
  background: var(--color-error);
}

.button-primary.danger:hover:not(:disabled) {
  background: #a34f4f;
}

.button-secondary {
  background: var(--bg-elevated);
  color: var(--color-ink);
  border: 1.5px solid var(--border-medium);
}

.button-secondary:hover {
  background: var(--bg-secondary);
  border-color: var(--border-strong);
}

/* Modals */
.oauth-modal-overlay,
.confirm-modal-overlay {
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
  margin-bottom: var(--space-6);
}

.confirm-modal {
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  max-width: 440px;
  width: 90%;
  box-shadow: var(--shadow-2xl);
  animation: scaleIn 0.3s var(--ease-out);
}

.confirm-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
}

.confirm-modal-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
}

.close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-ink-light);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.close-btn:hover {
  background: var(--bg-secondary);
  color: var(--color-ink);
}

.confirm-modal-body {
  padding: var(--space-6);
}

.confirm-modal-description {
  font-size: var(--text-base);
  color: var(--color-ink);
  margin-bottom: var(--space-4);
}

.confirm-modal-warning {
  font-size: var(--text-sm);
  color: var(--color-warning);
  padding: var(--space-3);
  background: rgba(212, 165, 116, 0.1);
  border-radius: var(--radius-md);
}

.confirm-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-6);
  border-top: 1px solid var(--border-subtle);
}

/* Animations */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
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

.spinner {
  animation: spin 0.8s linear infinite;
}

/* Responsive */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .providers-grid {
    grid-template-columns: 1fr;
  }

  .account-card {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-4);
  }

  .account-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .form-actions {
    flex-direction: column;
  }

  .button-primary,
  .button-secondary {
    width: 100%;
  }
}
</style>
