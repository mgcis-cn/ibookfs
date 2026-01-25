<template>
  <div :class="['toast', `toast-${type}`]">
    <component :is="iconComponent" class="toast-icon" :size="20" />
    <span class="toast-message">{{ message }}</span>
    <button class="toast-close" @click="$emit('close')">×</button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle, AlertCircle, XCircle, Info } from 'lucide-vue-next'

const props = defineProps<{
  type: 'success' | 'warning' | 'error' | 'info'
  message: string
}>()

defineEmits<{
  close: []
}>()

const iconComponent = computed(() => {
  switch (props.type) {
    case 'success':
      return CheckCircle
    case 'warning':
      return AlertCircle
    case 'error':
      return XCircle
    case 'info':
      return Info
  }
})
</script>

<style scoped>
.toast {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 320px;
  max-width: 400px;
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

.toast-info {
  border-left: 4px solid var(--color-accent);
}

.toast-icon {
  flex-shrink: 0;
}

.toast-success .toast-icon {
  color: var(--color-success);
}

.toast-warning .toast-icon {
  color: var(--color-warning);
}

.toast-error .toast-icon {
  color: var(--color-error);
}

.toast-info .toast-icon {
  color: var(--color-accent);
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
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-out);
}

.toast-close:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
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
</style>
