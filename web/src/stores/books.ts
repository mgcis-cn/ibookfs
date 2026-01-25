import { defineStore } from 'pinia'
import { api } from '@/api'
import type { Book, BookQuery } from '@/types'

interface BooksState {
  books: Book[]
  currentBook: Book | null
  query: BookQuery
  loading: boolean
  error: string | null
  total: number
}

export const useBooksStore = defineStore('books', {
  state: (): BooksState => ({
    books: [],
    currentBook: null,
    query: {
      filter: 'all',
      sort: 'createdAt',
      page: 1,
      pageSize: 12,
    },
    loading: false,
    error: null,
    total: 0,
  }),

  getters: {
    filteredBooks: (state) => state.books,
    bookById: (state) => (id: string) => state.books.find((b) => b.id === id),
    hasBooks: (state) => state.books.length > 0,
  },

  actions: {
    async fetchBooks(query?: BookQuery) {
      this.loading = true
      this.error = null

      try {
        if (query) {
          this.query = { ...this.query, ...query }
        }

        const response = await api.getBooks(this.query)

        if (response.success) {
          this.books = response.data.items
          this.total = response.data.total
        } else {
          this.error = response.message || 'Failed to fetch books'
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
      } finally {
        this.loading = false
      }
    },

    async fetchBook(id: string) {
      this.loading = true
      this.error = null

      try {
        const response = await api.getBook(id)

        if (response.success) {
          this.currentBook = response.data
        } else {
          this.error = response.message || 'Failed to fetch book'
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
      } finally {
        this.loading = false
      }
    },

    async createBook(data: Omit<Book, 'id' | 'createdAt' | 'updatedAt'>) {
      this.loading = true
      this.error = null

      try {
        const response = await api.createBook(data)

        if (response.success) {
          this.books.unshift(response.data)
          this.total++
          return response.data
        } else {
          this.error = response.message || 'Failed to create book'
          return null
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
        return null
      } finally {
        this.loading = false
      }
    },

    async updateBook(id: string, data: Partial<Book>) {
      this.loading = true
      this.error = null

      try {
        const response = await api.updateBook(id, data)

        if (response.success) {
          const index = this.books.findIndex((b) => b.id === id)
          if (index !== -1) {
            this.books[index] = response.data
          }
          if (this.currentBook?.id === id) {
            this.currentBook = response.data
          }
          return response.data
        } else {
          this.error = response.message || 'Failed to update book'
          return null
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
        return null
      } finally {
        this.loading = false
      }
    },

    async deleteBook(id: string) {
      this.loading = true
      this.error = null

      try {
        const response = await api.deleteBook(id)

        if (response.success) {
          const index = this.books.findIndex((b) => b.id === id)
          if (index !== -1) {
            this.books.splice(index, 1)
          }
          if (this.currentBook?.id === id) {
            this.currentBook = null
          }
          this.total--
          return true
        } else {
          this.error = response.message || 'Failed to delete book'
          return false
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
        return false
      } finally {
        this.loading = false
      }
    },

    setFilter(filter: Book['status'] | 'all') {
      this.query.filter = filter
      this.query.page = 1
      this.fetchBooks()
    },

    setSort(sort: BookQuery['sort']) {
      this.query.sort = sort
      this.fetchBooks()
    },

    setSearch(search: string) {
      this.query.search = search
      this.query.page = 1
      this.fetchBooks()
    },

    setPage(page: number) {
      this.query.page = page
      this.fetchBooks()
    },
  },
})
