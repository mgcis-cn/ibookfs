/**
 * API Client Module
 *
 * This module provides a centralized API client connecting to the Go backend.
 */

import type {
  ApiResponse,
  AuthResponse,
  Book,
  BookQuery,
  LoginRequest,
  PaginatedResponse,
  Photo,
  RegisterRequest,
  SendCodeRequest,
  Settings,
  User,
} from '@/types'
import { mockPhotos, mockSettings, delay, getPhotosByBookId } from './mockData'

// API Configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const USE_MOCK = import.meta.env.VITE_USE_MOCK !== 'false'

/**
 * Generic API handler for real backend calls
 */
async function apiCall<T>(
  endpoint: string,
  method: string = 'GET',
  data?: unknown
): Promise<T> {
  // Get token from localStorage
  const token = localStorage.getItem('authToken')

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }

  // Add Authorization header if token exists
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || error.message || 'Request failed')
  }

  return response.json()
}

/**
 * Transform backend book to frontend format
 */
function transformBook(backendBook: {
  id: number
  title: string
  author: string
  isbn: string
  publisher: string
  year: number
  pages: number
  uploaded_pages: number
  status: string
  cover: string
  created_at: string
  updated_at: string
}): Book {
  return {
    id: String(backendBook.id),
    title: backendBook.title,
    author: backendBook.author || '',
    isbn: backendBook.isbn,
    cover: backendBook.cover,
    totalPages: backendBook.pages || 0,
    uploadedPages: backendBook.uploaded_pages || 0,
    status: (backendBook.status || 'draft') as Book['status'],
    createdAt: backendBook.created_at,
    updatedAt: backendBook.updated_at,
  }
}

// ====================
// Book APIs
// ====================

/**
 * Get all books with optional filtering and pagination
 */
export async function getBooks(query: BookQuery = {}): Promise<ApiResponse<PaginatedResponse<Book>>> {
  if (USE_MOCK) {
    await delay(200)
    return {
      success: true,
      data: { items: [], total: 0, page: 1, pageSize: 12 },
    }
  }

  try {
    const params = new URLSearchParams()
    if (query.page) params.set('page', String(query.page))
    if (query.pageSize) params.set('page_size', String(query.pageSize))
    if (query.filter && query.filter !== 'all') params.set('status', query.filter)
    if (query.search) params.set('search', query.search)

    const response = await apiCall<{ data: unknown[]; total: number; message: string }>(
      `/books?${params.toString()}`
    )

    const books = (response.data || []).map((item) => transformBook(item as Parameters<typeof transformBook>[0]))

    return {
      success: true,
      data: {
        items: books,
        total: response.total || books.length,
        page: query.page || 1,
        pageSize: query.pageSize || 12,
      },
    }
  } catch (err) {
    return {
      success: false,
      data: { items: [], total: 0, page: 1, pageSize: 12 },
      message: err instanceof Error ? err.message : 'Failed to fetch books',
    }
  }
}

/**
 * Get a single book by ID
 */
export async function getBook(id: string): Promise<ApiResponse<Book>> {
  if (USE_MOCK) {
    await delay(150)
    return { success: false, data: null as unknown as Book, message: 'Book not found' }
  }

  try {
    const response = await apiCall<{ data: unknown; message: string }>(`/books/${id}`)
    return {
      success: true,
      data: transformBook(response.data as Parameters<typeof transformBook>[0]),
    }
  } catch (err) {
    return {
      success: false,
      data: null as unknown as Book,
      message: err instanceof Error ? err.message : 'Book not found',
    }
  }
}

/**
 * Create a new book
 */
export async function createBook(
  data: Omit<Book, 'id' | 'createdAt' | 'updatedAt'>
): Promise<ApiResponse<Book>> {
  if (USE_MOCK) {
    await delay(300)
    return { success: false, data: null as unknown as Book, message: 'Mock mode' }
  }

  try {
    const payload = {
      title: data.title,
      author: data.author,
      isbn: data.isbn || '',
      pages: data.totalPages || 0,
      uploaded_pages: data.uploadedPages || 0,
      status: data.status || 'draft',
      cover: data.cover || '',
    }

    const response = await apiCall<{ data: unknown; message: string }>('/books', 'POST', payload)
    return {
      success: true,
      data: transformBook(response.data as Parameters<typeof transformBook>[0]),
    }
  } catch (err) {
    return {
      success: false,
      data: null as unknown as Book,
      message: err instanceof Error ? err.message : 'Failed to create book',
    }
  }
}

/**
 * Update a book
 */
export async function updateBook(
  id: string,
  data: Partial<Book>
): Promise<ApiResponse<Book>> {
  if (USE_MOCK) {
    await delay(250)
    return { success: false, data: null as unknown as Book, message: 'Mock mode' }
  }

  try {
    const payload = {
      title: data.title,
      author: data.author,
      isbn: data.isbn,
      pages: data.totalPages,
      cover: data.cover,
    }

    const response = await apiCall<{ data: unknown; message: string }>(`/books/${id}`, 'PUT', payload)
    return {
      success: true,
      data: transformBook(response.data as Parameters<typeof transformBook>[0]),
    }
  } catch (err) {
    return {
      success: false,
      data: null as unknown as Book,
      message: err instanceof Error ? err.message : 'Failed to update book',
    }
  }
}

/**
 * Delete a book
 */
export async function deleteBook(id: string): Promise<ApiResponse<void>> {
  if (USE_MOCK) {
    await delay(200)
    return { success: false, data: undefined, message: 'Mock mode' }
  }

  try {
    await apiCall<{ message: string }>(`/books/${id}`, 'DELETE')
    return { success: true, data: undefined }
  } catch (err) {
    return {
      success: false,
      data: undefined,
      message: err instanceof Error ? err.message : 'Failed to delete book',
    }
  }
}

// ====================
// Photo APIs
// ====================

/**
 * Get all photos for a book
 *
 * TODO: Replace with GET /api/books/:id/photos
 */
export async function getPhotos(bookId: string): Promise<ApiResponse<Photo[]>> {
  await delay(150)

  const photos = getPhotosByBookId(bookId)

  return {
    success: true,
    data: photos,
  }
}

/**
 * Upload photos for a book
 *
 * TODO: Replace with POST /api/books/:id/photos
 */
export async function uploadPhotos(
  bookId: string,
  files: File[]
): Promise<ApiResponse<Photo[]>> {
  // Simulate upload progress
  // @ts-ignore - Intentionally unused for future progress tracking
  const _totalSize = files.reduce((acc, file) => acc + file.size, 0)
  let uploaded = 0

  for (const file of files) {
    await delay(100 + Math.random() * 200)
    uploaded += file.size
    // TODO: Emit progress event
  }

  const newPhotos: Photo[] = files.map((file, index) => ({
    id: `p${Date.now()}-${index}`,
    bookId,
    page: 0, // Will be set during organize step
    url: URL.createObjectURL(file),
    thumbnailUrl: URL.createObjectURL(file),
    size: file.size,
    width: 1200,
    height: 1800,
    createdAt: new Date().toISOString(),
  }))

  return {
    success: true,
    data: newPhotos,
  }
}

/**
 * Delete a photo
 *
 * TODO: Replace with DELETE /api/photos/:id
 */
export async function deletePhoto(photoId: string): Promise<ApiResponse<void>> {
  await delay(150)

  const index = mockPhotos.findIndex((p) => p.id === photoId)

  if (index === -1) {
    return {
      success: false,
      data: null as unknown as void,
      message: 'Photo not found',
    }
  }

  mockPhotos.splice(index, 1)

  return {
    success: true,
    data: undefined,
  }
}

/**
 * Reorder photos in a book
 *
 * TODO: Replace with PUT /api/books/:id/photos/reorder
 */
export async function reorderPhotos(
  // @ts-ignore - Intentionally unused for future API implementation
  _bookId: string,
  photoIds: string[]
): Promise<ApiResponse<void>> {
  await delay(200)

  // Update page numbers based on new order
  photoIds.forEach((id, index) => {
    const photo = mockPhotos.find((p) => p.id === id)
    if (photo) {
      photo.page = index + 1
    }
  })

  return {
    success: true,
    data: undefined,
  }
}

// ====================
// Settings APIs
// ====================

/**
 * Get user settings
 *
 * TODO: Replace with GET /api/settings
 */
export async function getSettings(): Promise<ApiResponse<Settings>> {
  await delay(100)

  return {
    success: true,
    data: mockSettings,
  }
}

/**
 * Update user settings
 *
 * TODO: Replace with PUT /api/settings
 */
export async function updateSettings(
  data: Partial<Settings>
): Promise<ApiResponse<Settings>> {
  await delay(200)

  Object.assign(mockSettings, data)

  return {
    success: true,
    data: mockSettings,
  }
}

// ====================
// Auth APIs
// ====================

/**
 * Send verification code
 *
 * TODO: Replace with POST /api/auth/send-code
 */
export async function sendCode(
  data: SendCodeRequest
): Promise<ApiResponse<{ message: string }>> {
  await delay(500)

  // Mock: simulate rate limiting (60 seconds)
  // In real implementation, this would be handled by the backend
  const lastSent = localStorage.getItem(`lastCodeSent_${data.email}`)
  if (lastSent) {
    const elapsed = Date.now() - parseInt(lastSent)
    if (elapsed < 60000) {
      const remaining = Math.ceil((60000 - elapsed) / 1000)
      return {
        success: false,
        data: { message: '' },
        message: `请等待 ${remaining} 秒后重新发送`,
      }
    }
  }

  localStorage.setItem(`lastCodeSent_${data.email}`, Date.now().toString())

  // Mock: Store a verification code (in real app, this goes to email)
  // For demo purposes: use '123456' as the valid code
  localStorage.setItem(`verificationCode_${data.email}`, '123456')

  return {
    success: true,
    data: { message: '验证码已发送' },
  }
}

/**
 * Login with email and verification code
 *
 * TODO: Replace with POST /api/auth/login
 */
export async function login(
  data: LoginRequest
): Promise<ApiResponse<AuthResponse>> {
  await delay(600)

  // Mock: verify the code
  const storedCode = localStorage.getItem(`verificationCode_${data.email}`)
  if (!storedCode || storedCode !== data.code) {
    return {
      success: false,
      data: {} as AuthResponse,
      message: '验证码错误或已过期',
    }
  }

  // Mock: create user session
  const authResponse: AuthResponse = {
    user: {
      id: 'user_1',
      email: data.email,
      emailVerified: true,
      first_name: '',
      last_name: '小李',
      status: 'active',
      role: 'user',
      language: 'zh-CN',
      timezone: 'Asia/Shanghai',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    token: 'mock_token_' + Date.now(),
  }

  // Store auth state
  localStorage.setItem('authToken', authResponse.token)
  localStorage.setItem('currentUser', JSON.stringify(authResponse.user))

  return {
    success: true,
    data: authResponse,
  }
}

/**
 * Register with name, email and verification code
 *
 * TODO: Replace with POST /api/auth/register
 */
export async function register(
  data: RegisterRequest
): Promise<ApiResponse<AuthResponse>> {
  await delay(700)

  // Mock: verify the code
  const storedCode = localStorage.getItem(`verificationCode_${data.email}`)
  if (!storedCode || storedCode !== data.code) {
    return {
      success: false,
      data: {} as AuthResponse,
      message: '验证码错误或已过期',
    }
  }

  // Mock: check if email already exists
  const existingUsers = JSON.parse(localStorage.getItem('users') || '[]')
  if (existingUsers.some((u: { email: string }) => u.email === data.email)) {
    return {
      success: false,
      data: {} as AuthResponse,
      message: '该邮箱已被注册',
    }
  }

  // Mock: create new user
  const newUser = {
    id: 'user_' + Date.now(),
    email: data.email,
    emailVerified: true,
    first_name: data.first_name,
    last_name: data.last_name,
    status: 'active' as const,
    role: 'user' as const,
    language: 'zh-CN',
    timezone: 'Asia/Shanghai',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  existingUsers.push(newUser)
  localStorage.setItem('users', JSON.stringify(existingUsers))

  const authResponse: AuthResponse = {
    user: newUser,
    token: 'mock_token_' + Date.now(),
  }

  // Store auth state
  localStorage.setItem('authToken', authResponse.token)
  localStorage.setItem('currentUser', JSON.stringify(authResponse.user))

  return {
    success: true,
    data: authResponse,
  }
}

/**
 * Logout current user
 *
 * TODO: Replace with POST /api/auth/logout
 */
export async function logout(): Promise<ApiResponse<void>> {
  await delay(200)

  localStorage.removeItem('authToken')
  localStorage.removeItem('currentUser')

  return {
    success: true,
    data: undefined,
  }
}

/**
 * Get current user
 *
 * TODO: Replace with GET /api/auth/me
 */
export async function getCurrentUser(): Promise<ApiResponse<User | null>> {
  await delay(150)

  const userStr = localStorage.getItem('currentUser')
  if (!userStr) {
    return {
      success: true,
      data: null,
    }
  }

  return {
    success: true,
    data: JSON.parse(userStr),
  }
}

// Export API object for convenient importing
export const api = {
  // Books
  getBooks,
  getBook,
  createBook,
  updateBook,
  deleteBook,

  // Photos
  getPhotos,
  uploadPhotos,
  deletePhoto,
  reorderPhotos,

  // Settings
  getSettings,
  updateSettings,

  // Auth
  sendCode,
  login,
  register,
  logout,
  getCurrentUser,
}
