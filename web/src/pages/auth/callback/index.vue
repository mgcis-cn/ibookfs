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
  const provider = route.query.provider as OAuthProvider
  const code = route.query.code as string
  const state = route.query.state as string
  const errorParam = route.query.error as string
  const errorDescription = route.query.error_description as string

  // Handle OAuth errors
  if (errorParam) {
    error.value = errorDescription || errorParam
    sendMessage({ type: 'oauth_error', provider, error: error.value })
    setTimeout(() => window.close(), 3000)
    return
  }

  // Validate required parameters
  if (!provider || !code || !state) {
    error.value = 'Missing required OAuth parameters'
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

  // Wait a bit before closing to ensure message is sent
  setTimeout(() => {
    window.close()
  }, 1000)
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
