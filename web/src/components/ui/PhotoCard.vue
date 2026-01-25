<template>
  <div class="card-photo" @click="handleClick">
    <img :src="photo.thumbnailUrl || photo.url" :alt="`第 ${photo.page} 页`" />
    <div class="card-photo-overlay">
      <span class="card-photo-page">第 {{ photo.page }} 页</span>
      <button v-if="showDelete" class="card-photo-delete" @click.stop="handleDelete">
        <Trash2 :size="16" />
      </button>
    </div>
    <div v-if="photo.ocrText" class="card-photo-ocr-indicator">
      <FileText :size="14" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { FileText, Trash2 } from 'lucide-vue-next'
import type { Photo } from '@/types'

const props = defineProps<{
  photo: Photo
  showDelete?: boolean
}>()

const emit = defineEmits<{
  click: [photo: Photo]
  delete: [photo: Photo]
}>()

const handleClick = () => {
  emit('click', props.photo)
}

const handleDelete = () => {
  emit('delete', props.photo)
}
</script>

<style scoped>
.card-photo {
  position: relative;
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--duration-base) var(--ease-out);
  cursor: pointer;
  aspect-ratio: 3 / 4;
}

.card-photo:hover {
  box-shadow: var(--shadow-md);
  transform: scale(1.02);
}

.card-photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.card-photo-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-out);
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: var(--space-3);
}

.card-photo:hover .card-photo-overlay {
  opacity: 1;
}

.card-photo-page {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: white;
}

.card-photo-delete {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  background: rgba(184, 92, 92, 0.8);
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-out);
}

.card-photo-delete:hover {
  background: rgba(184, 92, 92, 1);
}

.card-photo-ocr-indicator {
  position: absolute;
  top: var(--space-2);
  left: var(--space-2);
  padding: var(--space-1);
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-sm);
}
</style>
