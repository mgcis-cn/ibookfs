<template>
  <div :class="['skeleton', `skeleton-${variant}`]" />
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'rect' | 'circle' | 'text'
  }>(),
  {
    variant: 'rect',
  }
)
</script>

<style scoped>
.skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 0%,
    var(--bg-tertiary) 50%,
    var(--bg-secondary) 100%
  );
  background-size: 200% 100%;
  animation: skeletonLoading 1.5s ease-in-out infinite;
  border-radius: var(--radius-md);
}

.skeleton-circle {
  border-radius: var(--radius-full);
}

.skeleton-text {
  height: 1em;
  border-radius: var(--radius-sm);
}

@keyframes skeletonLoading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
