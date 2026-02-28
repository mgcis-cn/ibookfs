<template>
  <div class="oauth-callback-page">
    <div class="callback-container">
      <div class="callback-content">
        <div class="loading-spinner">
          <svg class="spinner" width="48" height="48" viewBox="0 0 48 48">
            <circle cx="24" cy="24" r="20" stroke="currentColor" stroke-width="3" fill="none" opacity="0.25" />
            <path d="M24 4a20 20 0 0120 20" stroke="currentColor" stroke-width="3" fill="none" />
          </svg>
        </div>
        <h2 class="callback-title">正在登录...</h2>
        <p class="callback-description">请稍候，我们正在完成您的登录</p>
        <p v-if="error" class="callback-error">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { OAuthProvider } from '@/types'

const route = useRoute()
const error = ref<string | null>(null)

onMounted(async () => {
  const code = route.query.code as string
  const state = route.query.state as string
  const errorParam = route.query.error as string
  const errorDescription = route.query.error_description as string

  // Get provider from URL path (e.g., /auth/callback/github)
  const provider = route.params.provider as OAuthProvider

  // Handle OAuth errors from provider (e.g., rate limit, access denied)
  if (errorParam) {
    // Map common error codes to user-friendly messages
    const errorMessages: Record<string, string> = {
      'access_denied': '授权被拒绝',
      'rate_limit': '请求太频繁，请10分钟后再试',
      'temporarily_unavailable': '服务暂时不可用，请稍后再试',
      'server_error': '服务器错误，请稍后再试',
    }
    error.value = errorMessages[errorParam] || errorDescription || errorParam
    
    // Check for rate limit indicators
    if (errorParam.includes('rate') || errorParam.includes('limit') || errorParam === '429') {
      error.value = '请求太频繁，请等待10-30分钟后再试'
    }
    
    if (provider) {
      sendMessage({ type: 'oauth_error', provider, error: error.value })
    }
    setTimeout(() => window.close(), 5000)
    return
  }

  // Check provider first - if missing, likely a stale/invalid callback
  if (!provider) {
    error.value = '授权已过期，请重新点击登录按钮'
    setTimeout(() => window.close(), 3000)
    return
  }

  // Validate required parameters from OAuth provider
  if (!code || !state) {
    error.value = '授权失败：缺少必要参数'
    sendMessage({ type: 'oauth_error', provider, error: error.value })
    setTimeout(() => window.close(), 3000)
    return
  }

  // Send callback data to parent window
  sendMessage({
    type: 'oauth_callback',
    provider,
    code,
    state,
  })

  // Close popup after sending message
  setTimeout(() => window.close(), 500)
})

/**
 * Send message to parent window
 */
function sendMessage(data: {
  type: 'oauth_callback' | 'oauth_error'
  provider: OAuthProvider
  code?: string
  state?: string
  error?: string
}) {
  if (window.opener) {
    window.opener.postMessage(data, window.location.origin)
  } else {
    // If no opener (opened in same tab), redirect to home with data
    const params = new URLSearchParams()
    if (data.type === 'oauth_callback' && data.code && data.state) {
      params.set('type', 'oauth_callback')
      params.set('provider', data.provider)
      params.set('code', data.code)
      params.set('state', data.state)
    } else if (data.type === 'oauth_error') {
      params.set('type', 'oauth_error')
      params.set('provider', data.provider)
      params.set('error', data.error || 'Unknown error')
    }
    window.location.href = `/?${params.toString()}`
  }
}
</script>

<style scoped>
.oauth-callback-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
}

.callback-container {
  text-align: center;
  padding: var(--space-8);
}

.callback-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
}

.loading-spinner {
  color: var(--color-accent);
}

.spinner {
  animation: spin 0.8s linear infinite;
}

.callback-title {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
}

.callback-description {
  font-size: var(--text-base);
  color: var(--color-ink-light);
}

.callback-error {
  font-size: var(--text-sm);
  color: var(--color-error);
  background: var(--color-error-light);
  padding: var(--space-3);
  border-radius: var(--radius-md);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
