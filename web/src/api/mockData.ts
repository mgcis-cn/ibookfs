import type { Book, Photo } from '@/types'

// Mock Books Data
export const mockBooks: Book[] = [
  {
    id: '1',
    title: '深入理解计算机系统',
    author: 'Randal E. Bryant · David R. O\'Hallaron',
    isbn: '978-7-111-54493-7',
    cover: 'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=400&h=600&fit=crop',
    totalPages: 120,
    uploadedPages: 120,
    status: 'completed',
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-20T15:30:00Z',
  },
  {
    id: '2',
    title: '设计心理学',
    author: 'Don Norman',
    isbn: '978-7-508-64738-8',
    cover: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=400&h=600&fit=crop',
    totalPages: 86,
    uploadedPages: 86,
    status: 'completed',
    createdAt: '2024-01-10T08:30:00Z',
    updatedAt: '2024-01-18T12:00:00Z',
  },
  {
    id: '3',
    title: '算法导论',
    author: 'Thomas H. Cormen',
    isbn: '978-7-111-40701-0',
    cover: 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=400&h=600&fit=crop',
    totalPages: 200,
    uploadedPages: 156,
    status: 'uploading',
    createdAt: '2024-01-20T14:00:00Z',
    updatedAt: '2024-01-24T09:15:00Z',
  },
  {
    id: '4',
    title: 'JavaScript高级程序设计',
    author: 'Matt Frisbie',
    isbn: '978-7-115-54538-1',
    cover: 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=400&h=600&fit=crop',
    totalPages: 150,
    uploadedPages: 45,
    status: 'processing',
    createdAt: '2024-01-22T16:45:00Z',
    updatedAt: '2024-01-23T11:20:00Z',
  },
  {
    id: '5',
    title: '代码整洁之道',
    author: 'Robert C. Martin',
    cover: 'https://images.unsplash.com/photo-1589829085413-56de8ae18c73?w=400&h=600&fit=crop',
    totalPages: 80,
    uploadedPages: 0,
    status: 'draft',
    createdAt: '2024-01-24T10:00:00Z',
    updatedAt: '2024-01-24T10:00:00Z',
  },
]

// Mock Photos Data
export const mockPhotos: Photo[] = [
  // Photos for book 1
  {
    id: 'p1',
    bookId: '1',
    page: 1,
    url: 'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800&h=1200&fit=crop',
    thumbnailUrl: 'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=200&h=300&fit=crop',
    size: 2458624,
    width: 2400,
    height: 3600,
    createdAt: '2024-01-15T10:05:00Z',
  },
  {
    id: 'p2',
    bookId: '1',
    page: 2,
    url: 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=800&h=1200&fit=crop',
    thumbnailUrl: 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=200&h=300&fit=crop',
    size: 2154896,
    width: 2400,
    height: 3600,
    createdAt: '2024-01-15T10:06:00Z',
  },
  {
    id: 'p3',
    bookId: '1',
    page: 3,
    url: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800&h=1200&fit=crop',
    thumbnailUrl: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=200&h=300&fit=crop',
    size: 2356412,
    width: 2400,
    height: 3600,
    createdAt: '2024-01-15T10:07:00Z',
  },
  {
    id: 'p4',
    bookId: '1',
    page: 4,
    url: 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=800&h=1200&fit=crop',
    thumbnailUrl: 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=200&h=300&fit=crop',
    size: 2256789,
    width: 2400,
    height: 3600,
    createdAt: '2024-01-15T10:08:00Z',
  },
]

// Mock Settings
export const mockSettings = {
  theme: 'light' as const,
  uploadQuality: 'high' as const,
  autoOptimize: true,
  defaultStartPage: 1,
}

// Helper function to simulate API delay
export const delay = (ms: number = 500) =>
  new Promise((resolve) => setTimeout(resolve, ms))

// Helper function to get photos by book ID
export const getPhotosByBookId = (bookId: string): Photo[] => {
  return mockPhotos.filter((p) => p.bookId === bookId)
}
