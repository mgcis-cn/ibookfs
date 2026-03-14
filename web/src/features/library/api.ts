import { apiCall } from '@/shared/api/client'
import type { Book, BookQuery } from './types'

interface BackendBook {
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
}

function transformBook(b: BackendBook): Book {
  return {
    id: String(b.id),
    title: b.title,
    author: b.author || '',
    isbn: b.isbn,
    cover: b.cover,
    totalPages: b.pages || 0,
    uploadedPages: b.uploaded_pages || 0,
    status: (b.status || 'draft') as Book['status'],
    createdAt: b.created_at,
    updatedAt: b.updated_at,
  }
}

export const booksApi = {
  list: async (query: BookQuery = {}) => {
    const params = new URLSearchParams()
    if (query.page) params.set('page', String(query.page))
    if (query.pageSize) params.set('page_size', String(query.pageSize))
    if (query.filter && query.filter !== 'all') params.set('status', query.filter)
    if (query.search) params.set('search', query.search)

    const qs = params.toString()
    const response = await apiCall<{ data: BackendBook[]; total: number }>(
      `/books${qs ? `?${qs}` : ''}`,
    )
    return {
      items: (response.data || []).map(transformBook),
      total: response.total || 0,
      page: query.page || 1,
      pageSize: query.pageSize || 12,
    }
  },

  get: async (id: string) => {
    const response = await apiCall<{ data: BackendBook }>(`/books/${id}`)
    return transformBook(response.data)
  },

  create: async (data: {
    title: string
    author?: string
    isbn?: string
    pages?: number
  }) => {
    const response = await apiCall<{ data: BackendBook }>('/books', 'POST', data)
    return transformBook(response.data)
  },

  update: async (id: string, data: Partial<{
    title: string
    author: string
    isbn: string
    pages: number
    cover: string
  }>) => {
    const response = await apiCall<{ data: BackendBook }>(`/books/${id}`, 'PUT', data)
    return transformBook(response.data)
  },

  delete: async (id: string) => {
    await apiCall<{ message: string }>(`/books/${id}`, 'DELETE')
  },
}
