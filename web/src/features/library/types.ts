export interface Book {
  id: string
  title: string
  author: string
  isbn?: string
  cover?: string
  totalPages: number
  uploadedPages: number
  status: BookStatus
  createdAt: string
  updatedAt: string
}

export type BookStatus = 'draft' | 'uploading' | 'processing' | 'completed'

export interface Photo {
  id: string
  bookId: string
  page: number
  url: string
  thumbnailUrl: string
  size: number
  width: number
  height: number
  ocrText?: string
  createdAt: string
}

export interface UploadPhoto {
  id: string
  file: File
  preview: string
  page: number
  status: 'pending' | 'uploading' | 'processing' | 'completed' | 'failed'
  progress: number
}

export type BookFilter = 'all' | 'draft' | 'uploading' | 'processing' | 'completed'
export type BookSort = 'createdAt' | 'title' | 'author' | 'status'

export interface BookQuery {
  filter?: BookFilter
  sort?: BookSort
  search?: string
  page?: number
  pageSize?: number
}
