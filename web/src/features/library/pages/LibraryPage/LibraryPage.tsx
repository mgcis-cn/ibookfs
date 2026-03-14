import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { booksApi } from '@/features/library/api'
import { useToastStore } from '@/shared/stores/toast'
import type { Book, BookStatus } from '@/features/library/types'
import { Search, Grid3x3, List, Plus, BookOpen, Library as LibraryIcon } from 'lucide-react'
import './LibraryPage.css'

const STATUS_LABELS: Record<BookStatus, string> = {
  draft: '草稿',
  uploading: '上传中',
  processing: '处理中',
  completed: '已完成',
}

const FILTERS: { label: string; value: BookStatus | 'all' }[] = [
  { label: '全部', value: 'all' },
  { label: '草稿', value: 'draft' },
  { label: '上传中', value: 'uploading' },
  { label: '处理中', value: 'processing' },
  { label: '已完成', value: 'completed' },
]

export default function LibraryPage() {
  const navigate = useNavigate()
  const toast = useToastStore(s => s.add)

  const [books, setBooks] = useState<Book[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [currentFilter, setCurrentFilter] = useState<BookStatus | 'all'>('all')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [currentPage, setCurrentPage] = useState(1)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [creating, setCreating] = useState(false)
  const [newBook, setNewBook] = useState({ title: '', author: '', isbn: '', totalPages: 100 })
  const pageSize = 12

  const totalPages = Math.ceil(total / pageSize)

  useEffect(() => { fetchBooks() }, [currentPage, currentFilter, searchQuery])

  const fetchBooks = async () => {
    setLoading(true)
    try {
      const res = await booksApi.list({ page: currentPage, pageSize, filter: currentFilter === 'all' ? undefined : currentFilter, search: searchQuery || undefined })
      setBooks(res.items || [])
      setTotal(res.total || 0)
    } catch {
      setBooks([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value)
    setCurrentPage(1)
  }

  const handleFilterChange = (filter: BookStatus | 'all') => {
    setCurrentFilter(filter)
    setCurrentPage(1)
  }

  const handleCreateBook = async () => {
    if (!newBook.title || !newBook.author || !newBook.totalPages) {
      toast('warning', '请填写所有必填字段')
      return
    }
    setCreating(true)
    try {
      const created = await booksApi.create({ title: newBook.title, author: newBook.author, isbn: newBook.isbn, pages: newBook.totalPages })
      if (created) {
        toast('success', '书籍创建成功')
        setShowCreateModal(false)
        setNewBook({ title: '', author: '', isbn: '', totalPages: 100 })
        navigate(`/upload/${created.id}`)
      }
    } catch {
      toast('error', '创建失败，请重试')
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="page-library">
      {/* Page Header */}
      <header className="page-header">
        <div className="header-left">
          <h1 className="page-title">我的书架</h1>
          {!loading && <span className="book-count">{total} 本书籍</span>}
        </div>
        <div className="header-right">
          <div className="search-box">
            <Search size={18} className="search-icon" />
            <input type="text" placeholder="搜索书籍..." value={searchQuery} onChange={handleSearch} />
          </div>
          <div className="view-toggle">
            <button className={`toggle-btn ${viewMode === 'grid' ? 'active' : ''}`} onClick={() => setViewMode('grid')}>
              <Grid3x3 size={16} />
            </button>
            <button className={`toggle-btn ${viewMode === 'list' ? 'active' : ''}`} onClick={() => setViewMode('list')}>
              <List size={16} />
            </button>
          </div>
          <button className="btn-primary" onClick={() => setShowCreateModal(true)}>
            <Plus size={18} />
            <span>新建书籍</span>
          </button>
        </div>
      </header>

      {/* Filter Bar */}
      <div className="filter-bar">
        {FILTERS.map(f => (
          <button key={f.value} className={`filter-chip ${currentFilter === f.value ? 'active' : ''}`} onClick={() => handleFilterChange(f.value)}>
            {f.label}
          </button>
        ))}
      </div>

      {/* Loading Skeleton */}
      {loading && (
        <div className={`books-grid books-grid-${viewMode}`}>
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="skeleton-card">
              <div className="skeleton-cover" />
              <div className="skeleton-title" />
              <div className="skeleton-author" />
            </div>
          ))}
        </div>
      )}

      {/* Empty State */}
      {!loading && books.length === 0 && (
        <div className="empty-state">
          <LibraryIcon size={64} className="empty-icon" />
          <h2 className="empty-title">还没有书籍</h2>
          <p className="empty-description">开始创建你的第一本数字书籍</p>
          <button className="btn-primary" onClick={() => setShowCreateModal(true)}>
            <Plus size={18} />
            <span>创建书籍</span>
          </button>
        </div>
      )}

      {/* Books Grid */}
      {!loading && books.length > 0 && (
        <div className={`books-grid books-grid-${viewMode}`}>
          {books.map(book => (
            <div key={book.id} className="book-card" onClick={() => navigate(`/book/${book.id}`)}>
              <div className="book-card-cover">
                {book.cover ? (
                  <img src={book.cover} alt={book.title} />
                ) : (
                  <div className="book-cover-placeholder">
                    <BookOpen size={32} />
                  </div>
                )}
                <div className="book-card-overlay">
                  <span className={`status-badge status-${book.status}`}>{STATUS_LABELS[book.status]}</span>
                </div>
              </div>
              <div className="book-card-info">
                <h3 className="book-card-title">{book.title}</h3>
                <p className="book-card-author">{book.author}</p>
                <div className="book-card-progress">
                  <div className="progress-bar-bg">
                    <div className="progress-bar-fill" style={{ width: `${book.totalPages ? (book.uploadedPages / book.totalPages) * 100 : 0}%` }} />
                  </div>
                  <span className="progress-text">{book.uploadedPages}/{book.totalPages}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="pagination">
          <button className="btn-secondary" disabled={currentPage === 1} onClick={() => setCurrentPage(p => p - 1)}>上一页</button>
          <span className="pagination-info">第 {currentPage} / {totalPages} 页</span>
          <button className="btn-secondary" disabled={currentPage === totalPages} onClick={() => setCurrentPage(p => p + 1)}>下一页</button>
        </div>
      )}

      {/* Create Book Modal */}
      {showCreateModal && (
        <div className="modal-overlay" onClick={e => { if (e.target === e.currentTarget) setShowCreateModal(false) }}>
          <div className="modal">
            <h2 className="modal-title">创建新书籍</h2>
            <form onSubmit={e => { e.preventDefault(); handleCreateBook() }} className="modal-form">
              <div className="modal-field">
                <label>书名</label>
                <input type="text" placeholder="输入书籍标题" value={newBook.title} onChange={e => setNewBook(b => ({ ...b, title: e.target.value }))} required />
              </div>
              <div className="modal-field">
                <label>作者</label>
                <input type="text" placeholder="输入作者姓名" value={newBook.author} onChange={e => setNewBook(b => ({ ...b, author: e.target.value }))} required />
              </div>
              <div className="modal-field">
                <label>ISBN (可选)</label>
                <input type="text" placeholder="输入ISBN编号" value={newBook.isbn} onChange={e => setNewBook(b => ({ ...b, isbn: e.target.value }))} />
              </div>
              <div className="modal-field">
                <label>预计总页数</label>
                <input type="number" placeholder="输入预计页数" value={newBook.totalPages} onChange={e => setNewBook(b => ({ ...b, totalPages: Number(e.target.value) }))} required />
              </div>
            </form>
            <div className="modal-footer">
              <button className="btn-secondary" onClick={() => setShowCreateModal(false)}>取消</button>
              <button className="btn-primary" disabled={creating} onClick={handleCreateBook}>
                {creating ? '创建中...' : '创建书籍'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
