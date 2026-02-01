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

// Auth Configuration Type
export interface AuthConfig {
  oauthProviders: string[]
  emailEnabled: boolean
  oauthEnabled: boolean
  oauth: {
    github: { configured: boolean }
    gitee: { configured: boolean }
    icloud: { configured: boolean }
    google: { configured: boolean }
    wechat: { configured: boolean }
  }
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

// ================================================================
// AUTH TYPES - Enhanced with OAuth Support
// ================================================================

export type OAuthProvider = 'github' | 'gitee' | 'icloud' | 'google' | 'wechat'

export interface User {
  id: string
  email: string
  emailVerified: boolean
  first_name?: string
  last_name?: string
  display_name?: string
  avatarUrl?: string
  status: UserStatus
  role: UserRole
  preferences?: UserPreferences
  language: string
  timezone: string
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
  linkedAccounts?: LinkedAccount[]
}

export type UserStatus = 'active' | 'suspended' | 'deleted'

export type UserRole = 'user' | 'admin' | 'super_admin'

export interface UserPreferences {
  theme?: 'light' | 'dark' | 'system'
  notifications?: {
    email?: boolean
    push?: boolean
  }
}

export interface LinkedAccount {
  provider: 'email' | OAuthProvider
  providerEmail?: string
  providerUserId?: string
  providerUsername?: string
  isPrimary: boolean
  linkedAt?: string
  lastUsedAt?: string
}

export interface AuthResponse {
  user: User
  token: string
  refreshToken?: string
  expiresIn?: number
  isNewUser?: boolean
}

export interface SendCodeRequest {
  email: string
  type: 'login' | 'register' | 'bind_email' | 'reset_password'
}

export interface LoginRequest {
  email: string
  code: string
}

export interface RegisterRequest {
  first_name: string
  last_name: string
  email: string
  code: string
}

export interface OAuthAuthorizeRequest {
  provider: OAuthProvider
  redirectUri: string
  state?: string
}

export interface OAuthAuthorizeResponse {
  authorizeUrl: string
  state: string
}

export interface OAuthCallbackRequest {
  provider: OAuthProvider
  code: string
  state: string
}

export interface LinkOAuthRequest {
  provider: OAuthProvider
  code: string
  state: string
}

export interface UnlinkOAuthRequest {
  provider: OAuthProvider
  providerUserId?: string
}

export interface BindEmailRequest {
  email: string
  code: string
}

export interface SetPrimaryAccountRequest {
  provider: OAuthProvider | 'email'
  providerUserId?: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface UpdateProfileRequest {
  first_name?: string
  last_name?: string
  displayName?: string
  language?: string
  timezone?: string
}

// ================================================================
// STORE TYPES
// ================================================================

export interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

export type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: { user: User; token: string; refreshToken?: string } }
  | { type: 'AUTH_FAILURE'; payload: string }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'AUTH_UPDATE_USER'; payload: Partial<User> }
  | { type: 'AUTH_CLEAR_ERROR' }
