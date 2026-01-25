<template>
  <div class="empty-state">
    <component :is="icon" class="empty-icon" :size="iconSize" />
    <h3 class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

withDefaults(
  defineProps<{
    icon?: Component
    title: string
    description?: string
    iconSize?: number
  }>(),
  {
    iconSize: 64,
  }
)
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12);
  text-align: center;
}

.empty-icon {
  margin-bottom: var(--space-6);
  color: var(--text-tertiary);
  opacity: 0.5;
}

.empty-title {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
  color: var(--text-primary);
}

.empty-description {
  font-size: var(--text-base);
  color: var(--text-secondary);
  margin-bottom: var(--space-6);
  max-width: 400px;
}

.empty-state :deep(.button) {
  margin-top: var(--space-4);
}
</style>
