<template>
  <div class="page-book-detail">
    <!-- Breadcrumb -->
    <nav class="breadcrumb">
      <router-link to="/">书架</router-link>
      <span class="breadcrumb-separator">/</span>
      <span v-if="loading">加载中...</span>
      <span v-else>{{ currentBook?.title }}</span>
    </nav>

    <!-- Loading State -->
    <div v-if="loading" class="book-info-skeleton">
      <Skeleton class="skeleton-cover" variant="rect" />
      <div class="skeleton-meta">
        <Skeleton class="skeleton-title" variant="rect" />
        <Skeleton class="skeleton-author" variant="rect" />
      </div>
    </div>

    <!-- Book Info -->
    <section v-else-if="currentBook" class="book-info">
      <div class="book-cover-large">
        <img v-if="currentBook.cover" :src="currentBook.cover" :alt="currentBook.title" />
        <div v-else class="book-cover-placeholder">
          <BookOpen :size="48" />
        </div>
      </div>
      <div class="book-meta">
        <h1 class="book-title">{{ currentBook.title }}</h1>
        <p class="book-author">{{ currentBook.author }}</p>
        <div v-if="currentBook.isbn" class="book-isbn">ISBN: {{ currentBook.isbn }}</div>

        <div class="book-stats">
          <div class="stat-item">
            <span class="stat-label">总页数</span>
            <span class="stat-value">{{ currentBook.totalPages }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">已上传</span>
            <span class="stat-value">{{ currentBook.uploadedPages }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">状态</span>
            <span class="stat-value stat-status" :class="`status-${currentBook.status}`">
              {{ statusLabels[currentBook.status] }}
            </span>
          </div>
        </div>

        <div class="book-actions">
          <Button variant="primary" @click="goToUpload">
            <Upload :size="18" />
            上传照片
          </Button>
          <Button variant="secondary" @click="showEditModal = true">
            <Edit :size="18" />
            编辑信息
          </Button>
          <Button variant="secondary" @click="handleDelete">
            <Trash2 :size="18" />
            删除
          </Button>
        </div>
      </div>
    </section>

    <!-- Photos Section -->
    <section class="photos-section">
      <div class="section-header">
        <h2>书籍页面</h2>
        <div class="section-actions">
          <Button v-if="photos.length > 0" variant="text" @click="handleSelectAll">
            {{ allSelected ? '取消全选' : '全选' }}
          </Button>
          <Button v-if="selectedCount > 0" variant="text" @click="handleDeleteSelected">
            删除 ({{ selectedCount }})
          </Button>
        </div>
      </div>

      <!-- Photos Grid -->
      <div v-if="loading" class="photos-grid">
        <Skeleton v-for="i in 8" :key="i" class="photo-skeleton" variant="rect" />
      </div>

      <EmptyState
        v-else-if="photos.length === 0"
        :icon="Upload"
        title="还没有照片"
        description="上传第一张照片开始数字化"
      >
        <Button variant="primary" @click="goToUpload">
          <Upload :size="18" />
          上传照片
        </Button>
      </EmptyState>

      <div v-else class="photos-grid">
        <PhotoCard
          v-for="photo in photos"
          :key="photo.id"
          :photo="photo"
          :show-delete="selectedPhotos.has(photo.id)"
          @click="handlePhotoClick"
          @delete="handlePhotoDelete"
        />
      </div>
    </section>

    <!-- Edit Modal -->
    <Modal v-model="showEditModal" title="编辑书籍信息" @close="closeEditModal">
      <form @submit.prevent="handleUpdateBook">
        <Input
          v-model="editForm.title"
          label="书名"
          placeholder="输入书籍标题"
          required
        />
        <Input
          v-model="editForm.author"
          label="作者"
          placeholder="输入作者姓名"
          required
        />
        <Input
          v-model="editForm.isbn"
          label="ISBN (可选)"
          placeholder="输入ISBN编号"
        />
        <Input
          v-model.number="editForm.totalPages"
          type="number"
          label="预计总页数"
          placeholder="输入预计页数"
          required
        />
      </form>
      <template #footer>
        <Button variant="secondary" @click="closeEditModal">取消</Button>
        <Button variant="primary" :loading="updating" @click="handleUpdateBook">
          保存更改
        </Button>
      </template>
    </Modal>

    <!-- Photo Viewer Modal -->
    <Modal
      v-model="showPhotoViewer"
      :full-screen="true"
      :close-on-backdrop="true"
      @close="closePhotoViewer"
    >
      <div v-if="currentPhoto" class="photo-viewer">
        <img :src="currentPhoto.url" :alt="`第 ${currentPhoto.page} 页`" />
        <div class="photo-viewer-info">
          <span>第 {{ currentPhoto.page }} 页</span>
          <span v-if="currentPhoto.ocrText" class="ocr-badge">已识别文字</span>
        </div>
        <button class="photo-viewer-close" @click="closePhotoViewer">
          <X :size="24" />
        </button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBooksStore } from '@/stores/books'
import { useToastStore } from '@/stores/toast'
import { api } from '@/api'
import { Upload, Edit, Trash2, BookOpen, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import PhotoCard from '@/components/ui/PhotoCard.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import type { Book, Photo } from '@/types'

const route = useRoute()
const router = useRouter()
const booksStore = useBooksStore()
const toastStore = useToastStore()

const currentBook = computed(() => booksStore.currentBook)
const loading = computed(() => booksStore.loading)

const photos = ref<Photo[]>([])
const selectedPhotos = ref<Set<string>>(new Set())
const showEditModal = ref(false)
const showPhotoViewer = ref(false)
const currentPhoto = ref<Photo | null>(null)
const updating = ref(false)

const editForm = ref<Partial<Book>>({
  title: '',
  author: '',
  isbn: '',
  totalPages: 0,
})

const statusLabels: Record<Book['status'], string> = {
  draft: '草稿',
  uploading: '上传中',
  processing: '处理中',
  completed: '已完成',
}

const selectedCount = computed(() => selectedPhotos.value.size)
const allSelected = computed(() => photos.value.length > 0 && selectedPhotos.value.size === photos.value.length)

onMounted(async () => {
  const bookId = route.params.id as string
  await booksStore.fetchBook(bookId)
  await loadPhotos(bookId)
})

const loadPhotos = async (bookId: string) => {
  const response = await api.getPhotos(bookId)
  if (response.success) {
    photos.value = response.data
  }
}

const goToUpload = () => {
  if (currentBook.value) {
    router.push(`/upload/${currentBook.value.id}`)
  }
}

const handlePhotoClick = (photo: Photo) => {
  currentPhoto.value = photo
  showPhotoViewer.value = true
}

const handlePhotoDelete = async (photo: Photo) => {
  const confirmed = confirm(`确定要删除第 ${photo.page} 页吗？`)
  if (!confirmed) return

  const response = await api.deletePhoto(photo.id)
  if (response.success) {
    photos.value = photos.value.filter((p) => p.id !== photo.id)
    toastStore.success('照片已删除')
    await booksStore.fetchBook(currentBook.value!.id)
  } else {
    toastStore.error(response.message || '删除失败')
  }
}

const handleSelectAll = () => {
  if (allSelected.value) {
    selectedPhotos.value.clear()
  } else {
    photos.value.forEach((photo) => selectedPhotos.value.add(photo.id))
  }
}

const handleDeleteSelected = async () => {
  const count = selectedCount.value
  const confirmed = confirm(`确定要删除选中的 ${count} 张照片吗？`)
  if (!confirmed) return

  for (const photoId of selectedPhotos.value) {
    await api.deletePhoto(photoId)
  }

  photos.value = photos.value.filter((p) => !selectedPhotos.value.has(p.id))
  selectedPhotos.value.clear()
  toastStore.success(`已删除 ${count} 张照片`)
  await booksStore.fetchBook(currentBook.value!.id)
}

const closeEditModal = () => {
  showEditModal.value = false
}

const handleUpdateBook = async () => {
  if (!currentBook.value) return

  updating.value = true
  const updated = await booksStore.updateBook(currentBook.value.id, editForm.value)
  updating.value = false

  if (updated) {
    toastStore.success('书籍信息已更新')
    closeEditModal()
  } else {
    toastStore.error('更新失败，请重试')
  }
}

const handleDelete = async () => {
  if (!currentBook.value) return

  const confirmed = confirm(`确定要删除《${currentBook.value.title}》吗？此操作不可撤销。`)
  if (!confirmed) return

  const deleted = await booksStore.deleteBook(currentBook.value.id)
  if (deleted) {
    toastStore.success('书籍已删除')
    router.push('/')
  } else {
    toastStore.error('删除失败，请重试')
  }
}

const closePhotoViewer = () => {
  showPhotoViewer.value = false
  currentPhoto.value = null
}
</script>

<style scoped>
.page-book-detail {
  padding: var(--space-8);
  max-width: var(--content-max-width);
  margin: 0 auto;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.breadcrumb a {
  color: var(--text-secondary);
  transition: color var(--duration-fast) var(--ease-out);
}

.breadcrumb a:hover {
  color: var(--color-accent);
}

.breadcrumb-separator {
  color: var(--text-tertiary);
}

.book-info {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: var(--space-8);
  margin-bottom: var(--space-10);
  padding: var(--space-8);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
}

.book-cover-large {
  aspect-ratio: 3 / 4;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-book);
  background: var(--bg-secondary);
}

.book-cover-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.book-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
}

.book-meta {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.book-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--leading-tight);
  color: var(--text-primary);
}

.book-author {
  font-size: var(--text-lg);
  color: var(--text-secondary);
}

.book-isbn {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}

.book-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  padding: var(--space-5) 0;
  border-top: 1px solid var(--border-subtle);
  border-bottom: 1px solid var(--border-subtle);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.stat-value {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.stat-status {
  font-size: var(--text-base);
}

.status-draft { color: var(--text-tertiary); }
.status-uploading { color: var(--color-warning); }
.status-processing { color: var(--color-accent); }
.status-completed { color: var(--color-success); }

.book-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: auto;
}

.photos-section {
  margin-top: var(--space-8);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.section-header h2 {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.section-actions {
  display: flex;
  gap: var(--space-2);
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-4);
}

.photo-skeleton {
  aspect-ratio: 3 / 4;
}

/* Photo Viewer */
.photo-viewer {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: var(--space-8);
}

.photo-viewer img {
  max-width: 100%;
  max-height: calc(100vh - 200px);
  object-fit: contain;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

.photo-viewer-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-4);
  font-size: var(--text-lg);
  color: var(--text-primary);
}

.cr-badge {
  padding: var(--space-1) var(--space-2);
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--text-sm);
}

.photo-viewer-close {
  position: absolute;
  top: var(--space-6);
  right: var(--space-6);
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-elevated);
  border-radius: var(--radius-full);
  box-shadow: var(--shadow-lg);
}

/* Skeleton Loading */
.book-info-skeleton {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: var(--space-8);
  margin-bottom: var(--space-10);
  padding: var(--space-8);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
}

.skeleton-cover {
  aspect-ratio: 3 / 4;
}

.skeleton-meta {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.skeleton-title {
  height: 40px;
  width: 80%;
}

.skeleton-author {
  height: 24px;
  width: 60%;
}

@media (max-width: 768px) {
  .page-book-detail {
    padding: var(--space-4);
  }

  .book-info,
  .book-info-skeleton {
    grid-template-columns: 1fr;
    gap: var(--space-5);
    padding: var(--space-5);
  }

  .book-stats {
    grid-template-columns: repeat(3, 1fr);
    gap: var(--space-3);
  }

  .book-actions {
    flex-direction: column;
  }

  .photos-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: var(--space-3);
  }
}
</style>
