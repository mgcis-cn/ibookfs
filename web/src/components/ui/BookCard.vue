<template>
  <div class="card-book" @click="handleClick">
    <div class="card-book-cover">
      <img v-if="book.cover" :src="book.cover" :alt="book.title" />
      <div v-else class="card-book-placeholder">
        <Book :size="32" />
      </div>
      <div class="card-book-badge">{{ book.uploadedPages }} / {{ book.totalPages }} 页</div>
    </div>
    <div class="card-book-content">
      <h3 class="card-book-title">{{ book.title }}</h3>
      <p class="card-book-meta">{{ book.author }}</p>
      <div class="card-book-progress">
        <div class="progress-bar" :style="{ width: progressPercent }" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Book } from 'lucide-vue-next'
import type { Book as BookType } from '@/types'

const props = defineProps<{
  book: BookType
}>()

const emit = defineEmits<{
  click: [book: BookType]
}>()

const progressPercent = computed(() => {
  if (props.book.totalPages === 0) return '0%'
  return `${(props.book.uploadedPages / props.book.totalPages) * 100}%`
})

const handleClick = () => {
  emit('click', props.book)
}
</script>

<style scoped>
.card-book {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-book);
  transition: all var(--duration-base) var(--ease-out);
  cursor: pointer;
}

.card-book:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-book-hover);
}

.card-book-cover {
  position: relative;
  aspect-ratio: 3 / 4;
  overflow: hidden;
  background: var(--bg-secondary);
}

.card-book-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--duration-slow) var(--ease-out);
}

.card-book:hover .card-book-cover img {
  transform: scale(1.05);
}

.card-book-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
}

.card-book-badge {
  position: absolute;
  bottom: var(--space-2);
  right: var(--space-2);
  padding: var(--space-1) var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-weight-medium);
  color: white;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-sm);
}

.card-book-content {
  padding: var(--space-4);
}

.card-book-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-primary);
}

.card-book-meta {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-3);
}

.card-book-progress {
  height: 4px;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--color-accent);
  border-radius: var(--radius-full);
  transition: width var(--duration-slow) var(--ease-out);
}
</style>
