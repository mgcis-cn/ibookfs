// Book Types
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

// Photo Types
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

// Upload Types
export interface UploadSession {
  id: string
  bookId: string
  step: UploadStep
  photos: UploadPhoto[]
  startPage: number
  status: UploadStatus
  createdAt: string
}

export type UploadStep = 'select' | 'organize' | 'confirm'

export type UploadStatus = 'pending' | 'uploading' | 'processing' | 'completed' | 'failed'

export interface UploadPhoto {
  id: string
  file: File
  preview: string
  page: number
  status: UploadStatus
  progress: number
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

export interface ApiError {
  success: false
  error: string
  code?: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

// Filter & Sort Types
export type BookFilter = 'all' | 'draft' | 'uploading' | 'processing' | 'completed'

export type BookSort = 'createdAt' | 'title' | 'author' | 'status'

export interface BookQuery {
  filter?: BookFilter
  sort?: BookSort
  search?: string
  page?: number
  pageSize?: number
}

// Settings Types
export interface Settings {
  theme: 'light' | 'dark' | 'system'
  uploadQuality: 'original' | 'high' | 'standard'
  autoOptimize: boolean
  defaultStartPage: number
}

// Toast Types
export interface Toast {
  id: number
  type: 'success' | 'warning' | 'error' | 'info'
  message: string
  duration?: number
}

// Auth Types
export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  createdAt: string
}

export interface AuthResponse {
  user: User
  token: string
}

export interface SendCodeRequest {
  email: string
  type?: 'login' | 'register'
}

export interface LoginRequest {
  email: string
  code: string
}

export interface RegisterRequest {
  firstName: string
  lastName: string
  email: string
  code: string
}
