/**
 * API Client Module
 *
 * This module provides a centralized API client that currently uses mock data.
 * In the future, this can be easily replaced with real API calls.
 *
 * To integrate with a real backend:
 * 1. Replace the mock implementations with actual fetch/axios calls
 * 2. Update the base URL configuration
 * 3. Add authentication headers if needed
 * 4. Implement proper error handling for network failures
 */

import type {
  ApiResponse,
  Book,
  BookQuery,
  PaginatedResponse,
  Photo,
  Settings,
} from '@/types'
import { mockBooks, mockPhotos, mockSettings, delay, getPhotosByBookId } from './mockData'

// API Configuration
// @ts-ignore - Intentionally unused for future API implementation
const _API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const USE_MOCK = import.meta.env.VITE_USE_MOCK !== 'false'

/**
 * Generic API handler that can switch between mock and real API
 * This is a template for future real API implementation
 */
// @ts-ignore - Intentionally unused template function
async function apiCall<T>(
  // @ts-ignore - Intentionally unused for future API implementation
  _endpoint: string,
  // @ts-ignore - Intentionally unused for future API implementation
  _method: string = 'GET',
  // @ts-ignore - Intentionally unused for future API implementation
  _data?: unknown
): Promise<ApiResponse<T>> {
  if (USE_MOCK) {
    // Simulate network delay
    await delay()
  }

  // TODO: Implement real API call here when USE_MOCK is false
  // Example:
  // if (!USE_MOCK) {
  //   const response = await fetch(`${_API_BASE_URL}${_endpoint}`, {
  //     method: _method,
  //     headers: {
  //       'Content-Type': 'application/json',
  //       // Add auth headers if needed
  //     },
  //     body: _data ? JSON.stringify(_data) : undefined,
  //   })
  //   if (!response.ok) {
  //     throw new Error(await response.text())
  //   }
  //   return response.json()
  // }

  throw new Error('API endpoint not implemented in mock mode')
}

// ====================
// Book APIs
// ====================

/**
 * Get all books with optional filtering and pagination
 *
 * TODO: Replace with GET /api/books
 */
export async function getBooks(query: BookQuery = {}): Promise<ApiResponse<PaginatedResponse<Book>>> {
  await delay(200)

  let filtered = [...mockBooks]

  // Apply filter
  if (query.filter && query.filter !== 'all') {
    filtered = filtered.filter((book) => book.status === query.filter)
  }

  // Apply search
  if (query.search) {
    const searchLower = query.search.toLowerCase()
    filtered = filtered.filter(
      (book) =>
        book.title.toLowerCase().includes(searchLower) ||
        book.author.toLowerCase().includes(searchLower)
    )
  }

  // Apply sorting
  if (query.sort) {
    filtered.sort((a, b) => {
      switch (query.sort) {
        case 'title':
          return a.title.localeCompare(b.title)
        case 'author':
          return a.author.localeCompare(b.author)
        case 'createdAt':
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        case 'status':
          return a.status.localeCompare(b.status)
        default:
          return 0
      }
    })
  }

  // Apply pagination
  const page = query.page || 1
  const pageSize = query.pageSize || 12
  const start = (page - 1) * pageSize
  const paginated = filtered.slice(start, start + pageSize)

  return {
    success: true,
    data: {
      items: paginated,
      total: filtered.length,
      page,
      pageSize,
    },
  }
}

/**
 * Get a single book by ID
 *
 * TODO: Replace with GET /api/books/:id
 */
export async function getBook(id: string): Promise<ApiResponse<Book>> {
  await delay(150)

  const book = mockBooks.find((b) => b.id === id)

  if (!book) {
    return {
      success: false,
      data: null as unknown as Book,
      message: 'Book not found',
    }
  }

  return {
    success: true,
    data: book,
  }
}

/**
 * Create a new book
 *
 * TODO: Replace with POST /api/books
 */
export async function createBook(
  data: Omit<Book, 'id' | 'createdAt' | 'updatedAt'>
): Promise<ApiResponse<Book>> {
  await delay(300)

  const newBook: Book = {
    ...data,
    id: Date.now().toString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  mockBooks.push(newBook)

  return {
    success: true,
    data: newBook,
  }
}

/**
 * Update a book
 *
 * TODO: Replace with PUT /api/books/:id
 */
export async function updateBook(
  id: string,
  data: Partial<Book>
): Promise<ApiResponse<Book>> {
  await delay(250)

  const index = mockBooks.findIndex((b) => b.id === id)

  if (index === -1) {
    return {
      success: false,
      data: null as unknown as Book,
      message: 'Book not found',
    }
  }

  mockBooks[index] = {
    ...mockBooks[index],
    ...data,
    updatedAt: new Date().toISOString(),
  }

  return {
    success: true,
    data: mockBooks[index],
  }
}

/**
 * Delete a book
 *
 * TODO: Replace with DELETE /api/books/:id
 */
export async function deleteBook(id: string): Promise<ApiResponse<void>> {
  await delay(200)

  const index = mockBooks.findIndex((b) => b.id === id)

  if (index === -1) {
    return {
      success: false,
      data: null as unknown as void,
      message: 'Book not found',
    }
  }

  mockBooks.splice(index, 1)

  return {
    success: true,
    data: undefined,
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
}
