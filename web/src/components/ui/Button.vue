<template>
  <component
    :is="tag"
    :type="tag === 'button' ? nativeType : undefined"
    :to="tag === 'router-link' ? to : undefined"
    :class="classes"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <Loader2 v-if="loading" class="button-spinner" :size="iconSize" />
    <component v-else-if="icon" :is="icon" class="button-icon" :size="iconSize" />
    <span v-if="$slots.default" class="button-content">
      <slot />
    </span>
  </component>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { Loader2 } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'text'
    size?: 'sm' | 'md' | 'lg'
    icon?: Component
    iconOnly?: boolean
    disabled?: boolean
    loading?: boolean
    to?: string
    nativeType?: 'button' | 'submit' | 'reset'
    fullWidth?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
    disabled: false,
    loading: false,
    nativeType: 'button',
  }
)

const emit = defineEmits<{
  click: [event: Event]
}>()

const tag = computed(() => {
  if (props.to) return 'router-link'
  return 'button'
})

const classes = computed(() => [
  'button',
  `button-${props.variant}`,
  `button-${props.size}`,
  {
    'button-icon-only': props.iconOnly,
    'button-loading': props.loading,
    'button-full-width': props.fullWidth,
    'button-has-icon': props.icon,
  },
])

const iconSize = computed(() => {
  switch (props.size) {
    case 'sm':
      return 16
    case 'lg':
      return 20
    default:
      return 18
  }
})

const handleClick = (event: Event) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
  white-space: nowrap;
}

/* Variants */
.button-primary {
  background: var(--color-accent);
  color: white;
  border: none;
  box-shadow: var(--shadow-sm);
}

.button-primary:hover:not(:disabled) {
  background: var(--color-accent-dark);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.button-primary:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: var(--shadow-xs);
}

.button-secondary {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border-medium);
}

.button-secondary:hover:not(:disabled) {
  background: var(--bg-secondary);
  border-color: var(--border-strong);
}

.button-text {
  background: transparent;
  color: var(--text-secondary);
  border: none;
  padding: var(--space-2) var(--space-4);
}

.button-text:hover:not(:disabled) {
  color: var(--text-primary);
  background: var(--bg-secondary);
}

/* Sizes */
.button-sm {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
}

.button-md {
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
}

.button-lg {
  padding: var(--space-4) var(--space-8);
  font-size: var(--text-lg);
}

/* States */
.button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.button-icon-only {
  padding: var(--space-3);
  aspect-ratio: 1;
}

.button-full-width {
  width: 100%;
}

.button-loading {
  pointer-events: none;
}

.button-spinner {
  animation: spin 1s linear infinite;
}

.button-has-icon .button-content {
  display: inline-flex;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
