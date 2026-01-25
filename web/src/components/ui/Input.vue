<template>
  <div class="input-group">
    <label v-if="label" :for="inputId" class="input-label">
      {{ label }}
      <span v-if="required" class="input-required">*</span>
    </label>
    <div class="input-wrapper">
      <component
        :is="textarea ? 'textarea' : 'input'"
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :rows="rows"
        :class="['input-field', { 'input-error': error, 'input-has-icon': $slots.prefix || icon }]"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />
      <span v-if="$slots.prefix" class="input-prefix">
        <slot name="prefix" />
      </span>
      <component v-if="icon" :is="icon" class="input-icon" :size="18" />
    </div>
    <span v-if="hint && !error" class="input-hint">{{ hint }}</span>
    <span v-if="error" class="input-error-text">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { v4 as uuidv4 } from 'uuid'

withDefaults(
  defineProps<{
    modelValue?: string | number
    type?: string
    placeholder?: string
    label?: string
    hint?: string
    error?: string
    icon?: Component
    disabled?: boolean
    required?: boolean
    textarea?: boolean
    rows?: number
  }>(),
  {
    type: 'text',
    disabled: false,
    required: false,
    textarea: false,
    rows: 3,
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  blur: [event: FocusEvent]
  focus: [event: FocusEvent]
}>()

const inputId = computed(() => `input-${uuidv4()}`)

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  emit('update:modelValue', target.value)
}

const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}

const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
}
</script>

<style scoped>
.input-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input-label {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

.input-required {
  color: var(--color-error);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-field {
  width: 100%;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  color: var(--text-primary);
  background: var(--bg-elevated);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.input-field::placeholder {
  color: var(--text-tertiary);
  font-style: italic;
}

.input-field:hover:not(:disabled) {
  border-color: var(--border-strong);
}

.input-field:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.input-field:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--bg-secondary);
}

.input-field.input-error {
  border-color: var(--color-error);
}

.input-field.input-error:focus {
  box-shadow: 0 0 0 3px rgba(184, 92, 92, 0.2);
}

.input-has-icon > .input-field {
  padding-right: var(--space-10);
}

.input-prefix {
  position: absolute;
  left: var(--space-4);
  color: var(--text-secondary);
  pointer-events: none;
}

.input-has-icon > .input-field {
  padding-left: var(--space-10);
}

.input-icon {
  position: absolute;
  right: var(--space-3);
  color: var(--text-secondary);
  pointer-events: none;
}

.input-hint {
  font-size: var(--text-xs);
  color: var(--text-secondary);
}

.input-error-text {
  font-size: var(--text-xs);
  color: var(--color-error);
}

textarea.input-field {
  resize: vertical;
  min-height: 80px;
}
</style>
