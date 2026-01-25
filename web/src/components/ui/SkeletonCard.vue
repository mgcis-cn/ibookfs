<template>
  <div class="skeleton-book-card">
    <Skeleton class="skeleton-cover" variant="rect" />
    <div class="skeleton-content">
      <Skeleton class="skeleton-title" variant="rect" />
      <Skeleton class="skeleton-meta" variant="rect" />
      <Skeleton class="skeleton-progress" variant="rect" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Skeleton from './Skeleton.vue'
</script>

<style scoped>
.skeleton-book-card {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.skeleton-cover {
  aspect-ratio: 3 / 4;
}

.skeleton-content {
  padding: var(--space-4);
}

.skeleton-title {
  height: 20px;
  width: 80%;
  margin-bottom: var(--space-2);
}

.skeleton-meta {
  height: 14px;
  width: 60%;
  margin-bottom: var(--space-3);
}

.skeleton-progress {
  height: 4px;
  width: 100%;
}
</style>
