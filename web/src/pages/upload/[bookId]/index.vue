<template>
  <div class="page-upload">
    <!-- Steps Indicator -->
    <div class="steps-indicator">
      <div :class="['step', { active: currentStep === 1, completed: currentStep > 1 }]">
        <div class="step-number">
          <Check v-if="currentStep > 1" :size="20" />
          <span v-else>1</span>
        </div>
        <span class="step-label">选择照片</span>
      </div>
      <div class="step-divider" :class="{ completed: currentStep > 1 }"></div>
      <div :class="['step', { active: currentStep === 2, completed: currentStep > 2 }]">
        <div class="step-number">
          <Check v-if="currentStep > 2" :size="20" />
          <span v-else>2</span>
        </div>
        <span class="step-label">整理排序</span>
      </div>
      <div class="step-divider" :class="{ completed: currentStep > 2 }"></div>
      <div :class="['step', { active: currentStep === 3 }]">
        <div class="step-number">3</div>
        <span class="step-label">确认上传</span>
      </div>
    </div>

    <!-- Step 1: Select Photos -->
    <div v-show="currentStep === 1" class="upload-step">
      <div
        class="upload-area"
        :class="{ dragging: isDragging }"
        @drop.prevent="handleDrop"
        @dragover.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        @click="handleSelectFiles"
      >
        <Upload class="upload-icon" :size="48" />
        <span class="upload-text">拖拽照片到此处，或点击选择</span>
        <span class="upload-hint">支持 JPG、PNG 格式，最大 10MB</span>
        <input
          ref="fileInputRef"
          type="file"
          multiple
          accept="image/*"
          style="display: none"
          @change="handleFileSelect"
        />
      </div>
    </div>

    <!-- Step 2: Organize Photos -->
    <div v-show="currentStep === 2" class="upload-step">
      <div class="preview-toolbar">
        <span class="photo-count">已选择 {{ uploadPhotos.length }} 张照片</span>
        <div class="toolbar-actions">
          <Button variant="text" @click="handleAutoSort">
            <ArrowUpDown :size="16" />
            自动排序
          </Button>
          <Button variant="text" @click="handleClearAll">
            <X :size="16" />
            清空
          </Button>
        </div>
      </div>
      <div class="preview-grid">
        <div
          v-for="(photo, index) in uploadPhotos"
          :key="photo.id"
          :class="['preview-item', { selected: selectedPhotos.has(photo.id) }]"
          @click="toggleSelectPhoto(photo.id)"
        >
          <img :src="photo.preview" alt="" />
          <div class="preview-overlay">
            <span class="page-number">第 {{ index + 1 }} 页</span>
            <button class="remove-btn" @click.stop="handleRemovePhoto(photo.id)">
              <X :size="16" />
            </button>
          </div>
          <div v-if="selectedPhotos.has(photo.id)" class="selected-badge">
            <Check :size="16" />
          </div>
        </div>
      </div>
      <div class="upload-actions">
        <Button variant="secondary" @click="currentStep = 1">上一步</Button>
        <Button
          variant="primary"
          :disabled="uploadPhotos.length === 0"
          @click="currentStep = 3"
        >
          下一步
        </Button>
      </div>
    </div>

    <!-- Step 3: Confirm -->
    <div v-show="currentStep === 3" class="upload-step">
      <div class="upload-summary">
        <h2>确认上传</h2>
        <div class="summary-info">
          <p v-if="targetBook">
            书籍：<strong>{{ targetBook.title }}</strong>
          </p>
          <p>照片数量：<strong>{{ uploadPhotos.length }} 张</strong></p>
          <p>起始页码：<strong>第 {{ startPage }} 页</strong></p>
        </div>
      </div>

      <div v-if="uploading" class="upload-progress">
        <div class="progress-bar-container">
          <div class="progress-bar" :style="{ width: `${uploadProgress}%` }"></div>
        </div>
        <p class="progress-text">
          {{ uploadStatus }}
        </p>
      </div>

      <div class="upload-actions">
        <Button variant="secondary" :disabled="uploading" @click="currentStep = 2">
          上一步
        </Button>
        <Button
          variant="primary"
          :loading="uploading"
          :disabled="uploadPhotos.length === 0"
          @click="handleUpload"
        >
          开始上传
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBooksStore } from '@/stores/books'
import { useToastStore } from '@/stores/toast'
import { Upload, Check, X, ArrowUpDown } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import type { UploadPhoto, Book } from '@/types'

const route = useRoute()
const router = useRouter()
const booksStore = useBooksStore()
const toastStore = useToastStore()

const currentStep = ref(1)
const isDragging = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref('')
const fileInputRef = ref<HTMLInputElement | null>(null)

const uploadPhotos = ref<UploadPhoto[]>([])
const selectedPhotos = ref<Set<string>>(new Set())
const startPage = ref(1)
const targetBookId = ref<string | null>(null)
const targetBook = ref<Book | null>(null)

onMounted(async () => {
  const bookId = route.params.bookId as string
  if (bookId) {
    targetBookId.value = bookId
    await booksStore.fetchBook(bookId)
    targetBook.value = booksStore.currentBook
  }
})

const handleSelectFiles = () => {
  fileInputRef.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files) {
    addFiles(Array.from(files))
  }
  target.value = ''
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files) {
    addFiles(Array.from(files))
  }
}

const addFiles = (files: File[]) => {
  const imageFiles = files.filter((file) => file.type.startsWith('image/'))

  imageFiles.forEach((file) => {
    const preview = URL.createObjectURL(file)
    uploadPhotos.value.push({
      id: `upload-${Date.now()}-${Math.random()}`,
      file,
      preview,
      page: uploadPhotos.value.length + 1,
      status: 'pending',
      progress: 0,
    })
  })

  if (imageFiles.length < files.length) {
    toastStore.warning(`${files.length - imageFiles.length} 个非图片文件已被忽略`)
  }

  toastStore.success(`已添加 ${imageFiles.length} 张照片`)
}

const handleRemovePhoto = (id: string) => {
  const index = uploadPhotos.value.findIndex((p) => p.id === id)
  if (index !== -1) {
    uploadPhotos.value.splice(index, 1)
    selectedPhotos.value.delete(id)
  }
}

const handleClearAll = () => {
  if (confirm('确定要清空所有照片吗？')) {
    uploadPhotos.value = []
    selectedPhotos.value.clear()
  }
}

const handleAutoSort = () => {
  // Sort by file name (which usually contains page numbers)
  uploadPhotos.value.sort((a, b) =>
    a.file.name.localeCompare(b.file.name, undefined, { numeric: true })
  )
  toastStore.success('已按文件名自动排序')
}

const toggleSelectPhoto = (id: string) => {
  if (selectedPhotos.value.has(id)) {
    selectedPhotos.value.delete(id)
  } else {
    selectedPhotos.value.add(id)
  }
}

const handleUpload = async () => {
  if (!targetBookId.value) {
    toastStore.error('未指定目标书籍')
    return
  }

  uploading.value = true
  uploadProgress.value = 0
  uploadStatus.value = '准备上传...'

  try {
    const total = uploadPhotos.value.length
    // const files = uploadPhotos.value.map((p) => p.file)

    // Simulate upload progress
    const interval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 10
        uploadStatus.value = `上传中... ${Math.round(uploadProgress.value)}%`
      }
    }, 200)

    // TODO: Replace with actual API call
    // const response = await api.uploadPhotos(targetBookId.value, files)

    // Simulate completion
    await new Promise((resolve) => setTimeout(resolve, 2000))
    clearInterval(interval)
    uploadProgress.value = 100
    uploadStatus.value = '上传完成！'

    toastStore.success(`成功上传 ${total} 张照片`)

    // Navigate to book detail
    setTimeout(() => {
      router.push(`/book/${targetBookId.value}`)
    }, 1000)
  } catch (error) {
    toastStore.error('上传失败，请重试')
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.page-upload {
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-8);
}

.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--space-10);
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
}

.step-number {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  transition: all var(--duration-base) var(--ease-out);
}

.step.active .step-number {
  background: var(--color-accent);
  color: white;
}

.step.completed .step-number {
  background: var(--color-success);
  color: white;
}

.step-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.step.active .step-label {
  color: var(--text-primary);
  font-weight: var(--font-weight-medium);
}

.step-divider {
  width: 80px;
  height: 2px;
  background: var(--border-medium);
  margin: 0 var(--space-4);
}

.step-divider.completed {
  background: var(--color-success);
}

.upload-step {
  animation: fadeIn var(--duration-base) var(--ease-out);
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: var(--space-12);
  border: 2px dashed var(--border-medium);
  border-radius: var(--radius-xl);
  background: var(--bg-elevated);
  cursor: pointer;
  transition: all var(--duration-base) var(--ease-out);
}

.upload-area:hover,
.upload-area.dragging {
  border-color: var(--color-accent);
  background: var(--bg-secondary);
}

.upload-icon {
  margin-bottom: var(--space-6);
  color: var(--color-accent);
}

.upload-text {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.upload-hint {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.preview-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.photo-count {
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.toolbar-actions {
  display: flex;
  gap: var(--space-2);
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.preview-item {
  position: relative;
  aspect-ratio: 3 / 4;
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: pointer;
  transition: all var(--duration-base) var(--ease-out);
}

.preview-item:hover {
  transform: scale(1.05);
  box-shadow: var(--shadow-lg);
  z-index: 1;
}

.preview-item.selected {
  ring: 2px solid var(--color-accent);
}

.preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: var(--space-2);
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.preview-item:hover .preview-overlay {
  opacity: 1;
}

.page-number {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-semibold);
  color: white;
}

.remove-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  background: rgba(184, 92, 92, 0.9);
  border-radius: var(--radius-full);
  transition: all var(--duration-fast) var(--ease-out);
}

.remove-btn:hover {
  background: rgba(184, 92, 92, 1);
}

.selected-badge {
  position: absolute;
  top: var(--space-2);
  right: var(--space-2);
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-full);
}

.upload-actions {
  display: flex;
  gap: var(--space-3);
  justify-content: center;
}

.upload-summary {
  text-align: center;
  padding: var(--space-10);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-8);
}

.upload-summary h2 {
  font-size: var(--text-3xl);
  margin-bottom: var(--space-6);
  color: var(--text-primary);
}

.summary-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  font-size: var(--text-lg);
  color: var(--text-secondary);
}

.summary-info strong {
  color: var(--text-primary);
}

.upload-progress {
  text-align: center;
  padding: var(--space-8);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-8);
}

.progress-bar-container {
  height: 8px;
  background: var(--bg-secondary);
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: var(--space-4);
}

.progress-bar {
  height: 100%;
  background: var(--color-accent);
  border-radius: var(--radius-full);
  transition: width var(--duration-base) var(--ease-out);
}

.progress-text {
  font-size: var(--text-base);
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .page-upload {
    padding: var(--space-4);
  }

  .steps-indicator {
    flex-wrap: wrap;
    gap: var(--space-2);
  }

  .step-divider {
    display: none;
  }

  .upload-area {
    min-height: 300px;
    padding: var(--space-6);
  }

  .preview-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: var(--space-2);
  }
}
</style>
