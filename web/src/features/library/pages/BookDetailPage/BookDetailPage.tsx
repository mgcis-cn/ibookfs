import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { booksApi } from '@/features/library/api'
import { useToastStore } from '@/shared/stores/toast'
import type { Book, BookStatus } from '@/features/library/types'
import { BookOpen, Upload, Edit, Trash2 } from 'lucide-react'
import './BookDetailPage.css'

const STATUS_LABELS: Record<BookStatus, string> = {
  draft: '草稿',
  uploading: '上传中',
  processing: '处理中',
  completed: '已完成',
}

export default function BookDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToastStore(s => s.add)

  const [book, setBook] = useState<Book | null>(null)
  const [loading, setLoading] = useState(true)
  const [showEditModal, setShowEditModal] = useState(false)
  const [updating, setUpdating] = useState(false)
  const [editForm, setEditForm] = useState({ title: '', author: '', isbn: '', totalPages: 0 })

  useEffect(() => {
    if (!id) return
    setLoading(true)
    booksApi.get(id)
      .then(b => {
        setBook(b)
        setEditForm({ title: b.title, author: b.author, isbn: b.isbn || '', totalPages: b.totalPages })
      })
      .catch(() => toast('error', '获取书籍详情失败'))
      .finally(() => setLoading(false))
  }, [id])

  const goToUpload = () => {
    if (book) navigate(`/upload/${book.id}`)
  }

  const handleUpdateBook = async () => {
    if (!id || !book) return
    setUpdating(true)
    try {
      const updated = await booksApi.update(id, { title: editForm.title, author: editForm.author, isbn: editForm.isbn, pages: editForm.totalPages })
      setBook(updated)
      setShowEditModal(false)
      toast('success', '书籍信息已更新')
    } catch {
      toast('error', '更新失败，请重试')
    } finally {
      setUpdating(false)
    }
  }

  const handleDelete = async () => {
    if (!book) return
    if (!confirm(`确定要删除《${book.title}》吗？此操作不可撤销。`)) return
    try {
      await booksApi.delete(book.id)
      toast('success', '书籍已删除')
      navigate('/library')
    } catch {
      toast('error', '删除失败，请重试')
    }
  }

  if (loading) {
    return (
      <div className="page-book-detail">
        <div className="book-info-skeleton">
          <div className="skeleton-cover" />
          <div className="skeleton-meta">
            <div className="skeleton-title" />
            <div className="skeleton-author" />
          </div>
        </div>
      </div>
    )
  }

  if (!book) {
    return (
      <div className="page-book-detail">
        <div className="empty-state">
          <BookOpen size={64} />
          <h2>书籍未找到</h2>
          <Link to="/library">返回书架</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="page-book-detail">
      {/* Breadcrumb */}
      <nav className="breadcrumb">
        <Link to="/library">书架</Link>
        <span className="breadcrumb-separator">/</span>
        <span>{book.title}</span>
      </nav>

      {/* Book Info */}
      <section className="book-info">
        <div className="book-cover-large">
          {book.cover ? (
            <img src={book.cover} alt={book.title} />
          ) : (
            <div className="book-cover-placeholder">
              <BookOpen size={48} />
            </div>
          )}
        </div>
        <div className="book-meta">
          <h1 className="book-title">{book.title}</h1>
          <p className="book-author">{book.author}</p>
          {book.isbn && <div className="book-isbn">ISBN: {book.isbn}</div>}

          <div className="book-stats">
            <div className="stat-item">
              <span className="stat-label">总页数</span>
              <span className="stat-value">{book.totalPages}</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">已上传</span>
              <span className="stat-value">{book.uploadedPages}</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">状态</span>
              <span className={`stat-value stat-status status-${book.status}`}>
                {STATUS_LABELS[book.status]}
              </span>
            </div>
          </div>

          <div className="book-actions">
            <button className="btn-primary" onClick={goToUpload}>
              <Upload size={18} />
              <span>上传照片</span>
            </button>
            <button className="btn-secondary" onClick={() => setShowEditModal(true)}>
              <Edit size={18} />
              <span>编辑信息</span>
            </button>
            <button className="btn-secondary btn-danger" onClick={handleDelete}>
              <Trash2 size={18} />
              <span>删除</span>
            </button>
          </div>
        </div>
      </section>

      {/* Edit Modal */}
      {showEditModal && (
        <div className="modal-overlay" onClick={e => { if (e.target === e.currentTarget) setShowEditModal(false) }}>
          <div className="modal">
            <h2 className="modal-title">编辑书籍信息</h2>
            <form onSubmit={e => { e.preventDefault(); handleUpdateBook() }} className="modal-form">
              <div className="modal-field">
                <label>书名</label>
                <input type="text" placeholder="输入书籍标题" value={editForm.title} onChange={e => setEditForm(f => ({ ...f, title: e.target.value }))} required />
              </div>
              <div className="modal-field">
                <label>作者</label>
                <input type="text" placeholder="输入作者姓名" value={editForm.author} onChange={e => setEditForm(f => ({ ...f, author: e.target.value }))} required />
              </div>
              <div className="modal-field">
                <label>ISBN (可选)</label>
                <input type="text" placeholder="输入ISBN编号" value={editForm.isbn} onChange={e => setEditForm(f => ({ ...f, isbn: e.target.value }))} />
              </div>
              <div className="modal-field">
                <label>预计总页数</label>
                <input type="number" placeholder="输入预计页数" value={editForm.totalPages} onChange={e => setEditForm(f => ({ ...f, totalPages: Number(e.target.value) }))} required />
              </div>
            </form>
            <div className="modal-footer">
              <button className="btn-secondary" onClick={() => setShowEditModal(false)}>取消</button>
              <button className="btn-primary" disabled={updating} onClick={handleUpdateBook}>
                {updating ? '保存中...' : '保存更改'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
