import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { booksApi } from '@/features/library/api'
import { useToastStore } from '@/shared/stores/toast'
import type { Book } from '@/features/library/types'
import { Upload as UploadIcon, Check, X, ArrowUpDown, CheckCircle } from 'lucide-react'
import './UploadPage.css'

interface UploadPhoto {
  id: string
  file: File
  preview: string
  page: number
  status: 'pending' | 'uploading' | 'done' | 'error'
  progress: number
}

export default function UploadPage() {
  const { bookId } = useParams<{ bookId: string }>()
  const navigate = useNavigate()
  const toast = useToastStore(s => s.add)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [currentStep, setCurrentStep] = useState(1)
  const [isDragging, setIsDragging] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [uploadStatus, setUploadStatus] = useState('')
  const [uploadPhotos, setUploadPhotos] = useState<UploadPhoto[]>([])
  const [selectedPhotos, setSelectedPhotos] = useState<Set<string>>(new Set())
  const [targetBook, setTargetBook] = useState<Book | null>(null)
  const [uploadDone, setUploadDone] = useState(false)
  const [uploadedCount, setUploadedCount] = useState(0)
  const startPage = 1

  useEffect(() => {
    if (bookId) {
      booksApi.get(bookId).then(b => setTargetBook(b)).catch(() => {})
    }
  }, [bookId])

  const handleSelectFiles = () => fileInputRef.current?.click()

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) addFiles(Array.from(e.target.files))
    e.target.value = ''
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    if (e.dataTransfer?.files) addFiles(Array.from(e.dataTransfer.files))
  }

  const addFiles = (files: File[]) => {
    const imageFiles = files.filter(f => f.type.startsWith('image/'))
    
    const newPhotos: UploadPhoto[] = imageFiles.map((file, index) => ({
      id: `upload-${Date.now()}-${Math.random()}`,
      file,
      preview: URL.createObjectURL(file),
      page: uploadPhotos.length + index + 1,
      status: 'pending',
      progress: 0,
    }))
    
    setUploadPhotos(prev => [...prev, ...newPhotos])
    
    if (imageFiles.length < files.length) toast('warning', `${files.length - imageFiles.length} 个非图片文件已被忽略`)
    if (imageFiles.length > 0) {
      toast('success', `已添加 ${imageFiles.length} 张照片`)
      if (currentStep === 1) setCurrentStep(2)
    }
  }

  const handleRemovePhoto = (id: string) => {
    setUploadPhotos(prev => prev.filter(p => p.id !== id))
    setSelectedPhotos(prev => { const s = new Set(prev); s.delete(id); return s })
  }

  const handleClearAll = () => {
    if (confirm('确定要清空所有照片吗？')) {
      uploadPhotos.forEach(p => URL.revokeObjectURL(p.preview))
      setUploadPhotos([])
      setSelectedPhotos(new Set())
    }
  }

  const handleAutoSort = () => {
    setUploadPhotos(prev => [...prev].sort((a, b) => a.file.name.localeCompare(b.file.name, undefined, { numeric: true })))
    toast('success', '已按文件名自动排序')
  }

  const toggleSelectPhoto = (id: string) => {
    setSelectedPhotos(prev => {
      const s = new Set(prev)
      if (s.has(id)) s.delete(id); else s.add(id)
      return s
    })
  }

  const handleUpload = async () => {
    if (!bookId) { toast('error', '未指定目标书籍'); return }
    setUploading(true)
    setUploadProgress(0)
    setUploadStatus('准备上传...')
    try {
      const total = uploadPhotos.length
      let successCount = 0
      const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
      const token = localStorage.getItem('authToken')
      const headers: HeadersInit = {}
      if (token) headers['Authorization'] = `Bearer ${token}`

      for (let i = 0; i < total; i++) {
        const photo = uploadPhotos[i]
        if (!photo) continue
        setUploadStatus(`正在上传第 ${i + 1}/${total} 张照片...`)

        const formData = new FormData()
        formData.append('file', photo.file)

        const res = await fetch(`${API_BASE_URL}/images?book_id=${bookId}`, {
          method: 'POST',
          headers,
          body: formData
        })

        if (!res.ok) {
          throw new Error(`Upload failed for ${photo.file.name}`)
        }

        successCount++
        setUploadProgress(Math.round((successCount / total) * 100))
        setUploadPhotos(prev => prev.map(p => p.id === photo.id ? { ...p, status: 'done', progress: 100 } : p))
      }

      setUploadStatus('上传完成！')
      setUploadDone(true)
      setUploadedCount(successCount)
      toast('success', `成功上传 ${successCount} 张照片`)
    } catch (error) {
      console.error('Upload Error:', error)
      toast('error', '上传失败，请重试')
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="page-upload">
      {/* Steps Indicator */}
      <div className="steps-indicator">
        <div className={`step ${currentStep === 1 ? 'active' : ''} ${currentStep > 1 ? 'completed' : ''}`}>
          <div className="step-number">{currentStep > 1 ? <Check size={20} /> : <span>1</span>}</div>
          <span className="step-label">选择照片</span>
        </div>
        <div className={`step-divider ${currentStep > 1 ? 'completed' : ''}`} />
        <div className={`step ${currentStep === 2 ? 'active' : ''} ${currentStep > 2 ? 'completed' : ''}`}>
          <div className="step-number">{currentStep > 2 ? <Check size={20} /> : <span>2</span>}</div>
          <span className="step-label">整理排序</span>
        </div>
        <div className={`step-divider ${currentStep > 2 ? 'completed' : ''}`} />
        <div className={`step ${currentStep === 3 ? 'active' : ''}`}>
          <div className="step-number"><span>3</span></div>
          <span className="step-label">确认上传</span>
        </div>
      </div>

      {/* Step 1: Select Photos */}
      {currentStep === 1 && (
        <div className="upload-step">
          <div
            className={`upload-area ${isDragging ? 'dragging' : ''}`}
            onDrop={handleDrop}
            onDragOver={e => { e.preventDefault(); setIsDragging(true) }}
            onDragLeave={e => { e.preventDefault(); setIsDragging(false) }}
            onClick={handleSelectFiles}
          >
            <UploadIcon className="upload-icon" size={48} />
            <span className="upload-text">拖拽照片到此处，或点击选择</span>
            <span className="upload-hint">支持 JPG、PNG 格式，最大 10MB</span>
            <input ref={fileInputRef} type="file" multiple accept="image/*" style={{ display: 'none' }} onChange={handleFileSelect} />
          </div>
        </div>
      )}

      {/* Step 2: Organize Photos */}
      {currentStep === 2 && (
        <div className="upload-step">
          <div className="preview-toolbar">
            <span className="photo-count">已选择 {uploadPhotos.length} 张照片</span>
            <div className="toolbar-actions">
              <button className="btn-text" onClick={handleAutoSort}><ArrowUpDown size={16} /> 自动排序</button>
              <button className="btn-text" onClick={handleClearAll}><X size={16} /> 清空</button>
            </div>
          </div>
          <div className="preview-grid">
            {uploadPhotos.map((photo, index) => (
              <div key={photo.id} className={`preview-item ${selectedPhotos.has(photo.id) ? 'selected' : ''}`} onClick={() => toggleSelectPhoto(photo.id)}>
                <img src={photo.preview} alt="" />
                <div className="preview-overlay">
                  <span className="page-number">第 {index + 1} 页</span>
                  <button className="remove-btn" onClick={e => { e.stopPropagation(); handleRemovePhoto(photo.id) }}><X size={16} /></button>
                </div>
                {selectedPhotos.has(photo.id) && <div className="selected-badge"><Check size={16} /></div>}
              </div>
            ))}
          </div>
          <div className="upload-actions">
            <button className="btn-secondary" onClick={() => setCurrentStep(1)}>上一步</button>
            <button className="btn-primary" disabled={uploadPhotos.length === 0} onClick={() => setCurrentStep(3)}>下一步</button>
          </div>
        </div>
      )}

      {/* Step 3: Confirm / Uploading / Done */}
      {currentStep === 3 && (
        <div className="upload-step">
          {uploadDone ? (
            <>
              <div className="upload-success">
                <CheckCircle size={64} className="success-icon" />
                <h2>上传成功</h2>
                <p>已成功上传 <strong>{uploadedCount}</strong> 张照片到《{targetBook?.title}》</p>
              </div>
              <div className="upload-actions">
                <button className="btn-secondary" onClick={() => navigate(`/book/${bookId}`)}>查看书籍详情</button>
                <button className="btn-primary" onClick={() => {
                  uploadPhotos.forEach(p => URL.revokeObjectURL(p.preview))
                  setUploadPhotos([])
                  setSelectedPhotos(new Set())
                  setUploadDone(false)
                  setUploadedCount(0)
                  setUploadProgress(0)
                  setUploadStatus('')
                  setCurrentStep(1)
                }}>
                  继续上传
                </button>
              </div>
            </>
          ) : (
            <>
              <div className="upload-summary">
                <h2>确认上传</h2>
                <div className="summary-info">
                  {targetBook && <p>书籍：<strong>{targetBook.title}</strong></p>}
                  <p>照片数量：<strong>{uploadPhotos.length} 张</strong></p>
                  <p>起始页码：<strong>第 {startPage} 页</strong></p>
                </div>
              </div>

              {uploading && (
                <div className="upload-progress">
                  <div className="progress-bar-container">
                    <div className="progress-bar" style={{ width: `${uploadProgress}%` }} />
                  </div>
                  <p className="progress-text">{uploadStatus}</p>
                </div>
              )}

              <div className="upload-actions">
                <button className="btn-secondary" disabled={uploading} onClick={() => setCurrentStep(2)}>上一步</button>
                <button className="btn-primary" disabled={uploadPhotos.length === 0 || uploading} onClick={handleUpload}>
                  {uploading ? '上传中...' : '开始上传'}
                </button>
              </div>
            </>
          )}
        </div>
      )}
    </div>
  )
}
