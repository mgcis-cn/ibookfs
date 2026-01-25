<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-backdrop" @click.self="handleBackdropClick">
        <div
          ref="modalRef"
          class="modal"
          :class="[`modal-${size}`, { 'modal-full-screen': fullScreen }]"
          role="dialog"
          aria-modal="true"
        >
          <div class="modal-header">
            <h2 class="modal-title">
              <slot name="title">{{ title }}</slot>
            </h2>
            <button class="modal-close" @click="close" aria-label="关闭">
              <X :size="20" />
            </button>
          </div>

          <div class="modal-body">
            <slot />
          </div>

          <div v-if="$slots.footer" class="modal-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { X } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    size?: 'sm' | 'md' | 'lg' | 'xl'
    fullScreen?: boolean
    closeOnBackdrop?: boolean
  }>(),
  {
    size: 'md',
    fullScreen: false,
    closeOnBackdrop: true,
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
  open: []
}>()

const modalRef = ref<HTMLElement | null>(null)

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop) {
    close()
  }
}

watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    emit('open')
    // Prevent body scroll when modal is open
    document.body.style.overflow = 'hidden'
    // Store scroll position for restoration
    const scrollY = window.scrollY
    document.body.style.position = 'fixed'
    document.body.style.top = `-${scrollY}px`
    document.body.style.width = '100%'
    nextTick(() => {
      modalRef.value?.focus()
    })
  } else {
    // Restore body scroll when modal is closed
    const scrollY = parseInt(document.body.style.top || '0', 10) * -1
    document.body.style.overflow = ''
    document.body.style.position = ''
    document.body.style.top = ''
    document.body.style.width = ''
    window.scrollTo(0, scrollY)
  }
})
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal-backdrop);
  padding: var(--space-4);
}

.modal {
  width: 100%;
  max-height: 90vh;
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-sm {
  max-width: 400px;
}

.modal-md {
  max-width: 560px;
}

.modal-lg {
  max-width: 720px;
}

.modal-xl {
  max-width: 900px;
}

.modal-full-screen {
  width: 100vw;
  height: 100vh;
  max-width: none;
  max-height: none;
  border-radius: 0;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}

.modal-title {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
  flex-shrink: 0;
}

.modal-close:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.modal-body {
  padding: var(--space-6);
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  display: flex;
  gap: var(--space-3);
  justify-content: flex-end;
  padding: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  flex-shrink: 0;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--duration-base) var(--ease-out);
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: all var(--duration-base) var(--ease-out);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

@media (max-width: 768px) {
  .modal {
    max-width: 100%;
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    margin: auto;
    margin-bottom: 0;
  }

  .modal-backdrop {
    align-items: flex-end;
    padding: 0;
  }
}
</style>
