<template>
  <div class="page-library">
    <!-- Page Header -->
    <header class="page-header">
      <div class="header-left">
        <h1 class="page-title">我的书架</h1>
        <span v-if="!loading" class="book-count">{{ total }} 本书籍</span>
      </div>
      <div class="header-right">
        <div class="search-box">
          <Search :size="18" class="search-icon" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索书籍..."
            @input="handleSearch"
          />
        </div>
        <div class="view-toggle">
          <button
            :class="['toggle-btn', { active: viewMode === 'grid' }]"
            @click="viewMode = 'grid'"
          >
            <Grid3x3 :size="16" />
          </button>
          <button
            :class="['toggle-btn', { active: viewMode === 'list' }]"
            @click="viewMode = 'list'"
          >
            <List :size="16" />
          </button>
        </div>
        <Button variant="primary" @click="showCreateModal = true">
          <Plus :size="18" />
          新建书籍
        </Button>
      </div>
    </header>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <button
        v-for="filter in filters"
        :key="filter.value"
        :class="['filter-chip', { active: currentFilter === filter.value }]"
        @click="handleFilterChange(filter.value)"
      >
        {{ filter.label }}
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="books-grid">
      <SkeletonCard v-for="i in 6" :key="i" />
    </div>

    <!-- Empty State -->
    <EmptyState
      v-else-if="books.length === 0"
      :icon="Library"
      title="还没有书籍"
      description="开始创建你的第一本数字书籍"
    >
      <Button variant="primary" @click="showCreateModal = true">
        <Plus :size="18" />
        创建书籍
      </Button>
    </EmptyState>

    <!-- Books Grid -->
    <div v-else :class="['books-grid', `books-grid-${viewMode}`]">
      <BookCard
        v-for="book in books"
        :key="book.id"
        :book="book"
        @click="handleBookClick"
      />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <Button
        variant="secondary"
        :disabled="currentPage === 1"
        @click="goToPage(currentPage - 1)"
      >
        上一页
      </Button>
      <span class="pagination-info">第 {{ currentPage }} / {{ totalPages }} 页</span>
      <Button
        variant="secondary"
        :disabled="currentPage === totalPages"
        @click="goToPage(currentPage + 1)"
      >
        下一页
      </Button>
    </div>

    <!-- Create Book Modal -->
    <Modal v-model="showCreateModal" title="创建新书籍" @close="handleModalClose">
      <form @submit.prevent="handleCreateBook">
        <Input
          v-model="newBook.title"
          label="书名"
          placeholder="输入书籍标题"
          required
        />
        <Input
          v-model="newBook.author"
          label="作者"
          placeholder="输入作者姓名"
          required
        />
        <Input
          v-model="newBook.isbn"
          label="ISBN (可选)"
          placeholder="输入ISBN编号"
        />
        <Input
          v-model.number="newBook.totalPages"
          type="number"
          label="预计总页数"
          placeholder="输入预计页数"
          required
        />
      </form>
      <template #footer>
        <Button variant="secondary" @click="showCreateModal = false">取消</Button>
        <Button variant="primary" :loading="creating" @click="handleCreateBook">
          创建书籍
        </Button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBooksStore } from '@/stores/books'
import { useToastStore } from '@/stores/toast'
import { Search, Grid3x3, List, Plus, Library } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import BookCard from '@/components/ui/BookCard.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import type { Book, BookStatus } from '@/types'

const router = useRouter()
const booksStore = useBooksStore()
const toastStore = useToastStore()

const viewMode = ref<'grid' | 'list'>('grid')
const searchQuery = ref('')
const currentFilter = ref<BookStatus | 'all'>('all')
const showCreateModal = ref(false)
const creating = ref(false)

const newBook = ref<Partial<Book>>({
  title: '',
  author: '',
  isbn: '',
  totalPages: 100,
  uploadedPages: 0,
  status: 'draft',
})

const books = computed(() => booksStore.books)
const loading = computed(() => booksStore.loading)
const total = computed(() => booksStore.total)
const currentPage = computed(() => booksStore.query.page || 1)
const pageSize = computed(() => booksStore.query.pageSize || 12)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const filters = [
  { label: '全部', value: 'all' as const },
  { label: '草稿', value: 'draft' as const },
  { label: '上传中', value: 'uploading' as const },
  { label: '处理中', value: 'processing' as const },
  { label: '已完成', value: 'completed' as const },
]

onMounted(() => {
  booksStore.fetchBooks()
})

const handleSearch = (event: Event) => {
  const target = event.target as HTMLInputElement
  booksStore.setSearch(target.value)
}

const handleFilterChange = (filter: BookStatus | 'all') => {
  currentFilter.value = filter
  booksStore.setFilter(filter)
}

const handleBookClick = (book: Book) => {
  router.push(`/book/${book.id}`)
}

const handleCreateBook = async () => {
  if (!newBook.value.title || !newBook.value.author || !newBook.value.totalPages) {
    toastStore.warning('请填写所有必填字段')
    return
  }

  creating.value = true

  const created = await booksStore.createBook({
    ...newBook.value,
  } as Omit<Book, 'id' | 'createdAt' | 'updatedAt'>)

  creating.value = false

  if (created) {
    toastStore.success('书籍创建成功')
    showCreateModal.value = false
    // Reset form
    newBook.value = {
      title: '',
      author: '',
      isbn: '',
      totalPages: 100,
      uploadedPages: 0,
      status: 'draft',
    }
    // Navigate to upload page
    router.push(`/upload/${created.id}`)
  } else {
    toastStore.error('创建失败，请重试')
  }
}

const handleModalClose = () => {
  newBook.value = {
    title: '',
    author: '',
    isbn: '',
    totalPages: 100,
    uploadedPages: 0,
    status: 'draft',
  }
}

const goToPage = (page: number) => {
  booksStore.setPage(page)
}
</script>

<style scoped>
.page-library {
  padding: var(--space-8);
  max-width: var(--content-max-width);
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
  gap: var(--space-4);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
}

.page-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
}

.book-count {
  font-size: var(--text-base);
  color: var(--text-secondary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.search-box {
  position: relative;
}

.search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
  pointer-events: none;
}

.search-box input {
  width: 280px;
  padding: var(--space-3) var(--space-4) var(--space-3) var(--space-10);
  font-size: var(--text-base);
  color: var(--text-primary);
  background: var(--bg-elevated);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.search-box input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.view-toggle {
  display: flex;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: var(--space-1);
}

.toggle-btn {
  padding: var(--space-2) var(--space-3);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-out);
  cursor: pointer;
}

.toggle-btn.active {
  background: var(--bg-elevated);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
}

.filter-bar {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.filter-chip {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border: none;
  border-radius: var(--radius-full);
  white-space: nowrap;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.filter-chip:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.filter-chip.active {
  background: var(--color-accent);
  color: white;
}

.books-grid {
  display: grid;
  gap: var(--space-5);
}

.books-grid-grid {
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
}

.books-grid-list {
  grid-template-columns: 1fr;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  margin-top: var(--space-8);
}

.pagination-info {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .page-library {
    padding: var(--space-4);
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .header-right {
    flex-direction: column;
  }

  .search-box input {
    width: 100%;
  }

  .books-grid-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: var(--space-3);
  }

  .filter-bar {
    padding-bottom: var(--space-2);
  }
}
</style>
