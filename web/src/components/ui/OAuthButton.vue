<template>
  <button
    :type="type"
    :disabled="disabled || isLoading"
    :class="oauthButtonClass"
    @click="handleClick"
  >
    <span v-if="isLoading" class="oauth-spinner">
      <svg class="spinner" width="18" height="18" viewBox="0 0 18 18">
        <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="2" fill="none" opacity="0.25" />
        <path d="M9 2a7 7 0 017 7" stroke="currentColor" stroke-width="2" fill="none" />
      </svg>
    </span>
    <span v-else class="oauth-icon">
      <!-- GitHub Icon -->
      <svg v-if="provider === 'github'" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
      </svg>

      <!-- Gitee Icon -->
      <svg v-else-if="provider === 'gitee'" viewBox="0 0 24 24" fill="currentColor">
        <path d="M11.984 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.016 0zm6.09 5.333c.328 0 .593.266.592.593v1.482a.594.594 0 0 1-.593.592H9.777c-.982 0-1.778.796-1.778 1.778v5.63c0 .327.266.592.593.592h5.63c.982 0 1.778-.796 1.778-1.778v-.296a.593.593 0 0 0-.592-.593h-4.037a.594.594 0 0 1-.592-.593v-1.482a.593.593 0 0 1 .593-.592h6.815c.327 0 .593.265.593.592v3.408a4 4 0 0 1-4 4H5.926a.593.593 0 0 1-.593-.593V9.778a4.444 4.444 0 0 1 4.445-4.444h8.296z"/>
      </svg>

      <!-- iCloud Icon -->
      <svg v-else-if="provider === 'icloud'" viewBox="0 0 24 24" fill="currentColor">
        <path d="M17.6 9.48c1.03-.68 1.73-1.83 1.73-3.18 0-2.12-1.74-3.85-3.88-3.85-.65 0-1.26.16-1.79.44-.68-2.15-2.71-3.72-5.1-3.72-2.93 0-5.31 2.35-5.31 5.25 0 .28.03.55.08.82C1.23 6.77 0 8.45 0 10.42c0 2.43 1.99 4.41 4.44 4.41h13.12c2.45 0 4.44-1.98 4.44-4.41 0-2.12-1.51-3.9-3.52-4.28-.28-.07-.57-.11-.88-.11zM12.5 15l-1.5-3-1.5 3h3zm-2.5 4h1v-2h-1v2zm3 0h1v-2h-1v2z"/>
      </svg>

      <!-- Google Icon -->
      <svg v-else-if="provider === 'google'" viewBox="0 0 24 24" fill="currentColor">
        <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
        <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
        <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
        <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
      </svg>

      <!-- WeChat Icon -->
      <svg v-else-if="provider === 'wechat'" viewBox="0 0 24 24" fill="currentColor">
        <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z"/>
      </svg>
    </span>

    <span class="oauth-text">
      <slot>{{ providerText }}</slot>
    </span>

    <span v-if="!isLoading && isLinked" class="oauth-linked">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M20 6L9 17l-5-5"/>
      </svg>
    </span>

    <span v-if="!isLoading && isPrimary" class="oauth-primary">
      主账号
    </span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { OAuthProvider } from '@/types'

interface Props {
  provider: OAuthProvider
  variant?: 'default' | 'outlined' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit'
  disabled?: boolean
  isLoading?: boolean
  isLinked?: boolean
  isPrimary?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'md',
  type: 'button',
  disabled: false,
  isLoading: false,
  isLinked: false,
  isPrimary: false,
})

const emit = defineEmits<{
  click: []
}>()

const handleClick = () => {
  if (!props.disabled && !props.isLoading) {
    emit('click')
  }
}

const providerText: Record<OAuthProvider, string> = {
  github: 'GitHub',
  gitee: 'Gitee',
  icloud: 'iCloud',
  google: 'Google',
  wechat: '微信',
}

const oauthButtonClass = computed(() => {
  const classes = ['oauth-button', `oauth-button--${props.variant}`, `oauth-button--${props.size}`]

  if (props.disabled) classes.push('oauth-button--disabled')
  if (props.isLoading) classes.push('oauth-button--loading')
  if (props.isLinked) classes.push('oauth-button--linked')
  if (props.isPrimary) classes.push('oauth-button--primary')

  return classes.join(' ')
})
</script>

<style scoped>
.oauth-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  font-family: var(--font-body);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  position: relative;
  overflow: hidden;
}

/* Sizes */
.oauth-button--sm {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  gap: var(--space-2);
}

.oauth-button--sm .oauth-icon {
  width: 16px;
  height: 16px;
}

.oauth-button--md {
  padding: var(--space-3) var(--space-5);
  font-size: var(--text-base);
}

.oauth-button--md .oauth-icon {
  width: 20px;
  height: 20px;
}

.oauth-button--lg {
  padding: var(--space-4) var(--space-6);
  font-size: var(--text-lg);
}

.oauth-button--lg .oauth-icon {
  width: 24px;
  height: 24px;
}

/* Variants */
.oauth-button--default {
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  color: var(--color-ink);
}

.oauth-button--default:hover:not(.oauth-button--disabled):not(.oauth-button--loading) {
  background: var(--bg-secondary);
  border-color: var(--border-strong);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.oauth-button--outlined {
  background: transparent;
  border: 1.5px solid var(--border-medium);
  color: var(--color-ink);
}

.oauth-button--outlined:hover:not(.oauth-button--disabled):not(.oauth-button--loading) {
  background: var(--bg-secondary);
  border-color: var(--color-accent);
}

.oauth-button--ghost {
  background: transparent;
  border: 1.5px solid transparent;
  color: var(--color-ink-light);
}

.oauth-button--ghost:hover:not(.oauth-button--disabled):not(.oauth-button--loading) {
  background: var(--bg-secondary);
  color: var(--color-ink);
}

/* States */
.oauth-button--disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

.oauth-button--loading {
  pointer-events: none;
}

.oauth-button--linked {
  border-color: var(--color-success);
  background: rgba(122, 155, 118, 0.1);
}

.oauth-button--linked:hover:not(.oauth-button--disabled):not(.oauth-button--loading) {
  background: rgba(122, 155, 118, 0.15);
}

.oauth-button--primary {
  border-color: var(--color-accent);
  background: rgba(193, 123, 92, 0.1);
}

.oauth-button--primary:hover:not(.oauth-button--disabled):not(.oauth-button--loading) {
  background: rgba(193, 123, 92, 0.15);
}

/* Icon */
.oauth-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-ink);
}

.oauth-button--loading .oauth-icon {
  display: none;
}

/* Spinner */
.oauth-spinner {
  display: flex;
  align-items: center;
  justify-content: center;
}

.oauth-button:not(.oauth-button--loading) .oauth-spinner {
  display: none;
}

.spinner {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Text */
.oauth-text {
  flex: 1;
}

/* Linked Badge */
.oauth-linked {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  background: var(--color-success);
  color: white;
  border-radius: var(--radius-full);
}

.oauth-linked svg {
  width: 12px;
  height: 12px;
  stroke-width: 2.5;
}

.oauth-button:not(.oauth-button--linked) .oauth-linked {
  display: none;
}

/* Primary Badge */
.oauth-primary {
  padding: var(--space-1) var(--space-2);
  font-size: var(--text-xs);
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-semibold);
}

.oauth-button:not(.oauth-button--primary) .oauth-primary {
  display: none;
}

/* Provider-specific colors when linked */
.oauth-button--linked.oauth-provider--github {
  border-color: #333;
}

.oauth-button--linked.oauth-provider--gitee {
  border-color: #c71d23;
}

.oauth-button--linked.oauth-provider--icloud {
  border-color: #007aff;
}

.oauth-button--linked.oauth-provider--google {
  border-color: #4285f4;
}

.oauth-button--linked.oauth-provider--wechat {
  border-color: #07c160;
}

/* Focus visible */
.oauth-button:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}

/* Active state */
.oauth-button:active:not(.oauth-button--disabled):not(.oauth-button--loading) {
  transform: translateY(0);
}
</style>
