/**
 * Authentication API Client
 *
 * Connects to the Go backend for authentication operations.
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

/**
 * Generic API handler with JWT token support
 */
async function apiCall<T>(
  endpoint: string,
  method: string = 'GET',
  data?: unknown,
  requiresAuth: boolean = false
): Promise<T> {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }

  // Add auth token if required
  if (requiresAuth) {
    const token = localStorage.getItem('authToken')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }))
    throw new Error(error.message || 'Request failed')
  }

  const json = await response.json()
  // Wrap raw backend response into ApiResponse format
  return { success: true, data: json } as T
}

// Import types from global types
import type {
  SendCodeRequest,
  LoginRequest,
  RegisterRequest,
  LinkOAuthRequest,
  UnlinkOAuthRequest,
  BindEmailRequest,
  SetPrimaryAccountRequest,
  RefreshTokenRequest,
  UpdateProfileRequest,
  User,
  LinkedAccount,
  AuthResponse,
  ApiResponse,
} from '@/types'

// ====================
// Auth API Functions
// ====================

// Auth config types
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

/**
 * Get auth configuration
 * GET /api/v1/auth/config
 */
export async function getAuthConfig(): Promise<ApiResponse<AuthConfig>> {
  return apiCall<ApiResponse<AuthConfig>>('/auth/config', 'GET', undefined, false)
}

/**
 * Send verification code
 * POST /api/v1/auth/send-code
 */
export async function sendCode(data: SendCodeRequest): Promise<ApiResponse<{ message: string }>> {
  return apiCall<ApiResponse<{ message: string }>>('/auth/send-code', 'POST', data)
}

/**
 * Login with email and verification code
 * POST /api/v1/auth/login
 */
export async function login(data: LoginRequest): Promise<ApiResponse<AuthResponse>> {
  return apiCall<ApiResponse<AuthResponse>>('/auth/login', 'POST', data)
}

/**
 * Register with email and verification code
 * POST /api/v1/auth/register
 */
export async function register(data: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
  return apiCall<ApiResponse<AuthResponse>>('/auth/register', 'POST', data)
}

/**
 * OAuth authorize - get authorization URL
 * GET /api/v1/auth/oauth/authorize
 */
export async function oauthAuthorize(provider: string, redirectUri: string): Promise<
  ApiResponse<{ authorizeUrl: string; state: string }>
> {
  return apiCall<ApiResponse<{ authorizeUrl: string; state: string }>>(
    `/auth/oauth/authorize?provider=${provider}&redirect_uri=${encodeURIComponent(redirectUri)}`,
    'GET'
  )
}

/**
 * OAuth callback
 * POST /api/v1/auth/oauth/callback
 */
export async function oauthCallback(
  provider: string,
  code: string,
  state: string
): Promise<ApiResponse<AuthResponse>> {
  return apiCall<ApiResponse<AuthResponse>>('/auth/oauth/callback', 'POST', {
    provider,
    code,
    state,
  })
}

/**
 * Refresh token
 * POST /api/v1/auth/refresh
 */
export async function refreshToken(data: RefreshTokenRequest): Promise<
  ApiResponse<{ token: string; refreshToken?: string; expiresIn: number }>
> {
  return apiCall<ApiResponse<{ token: string; refreshToken?: string; expiresIn: number }>>(
    '/auth/refresh',
    'POST',
    data
  )
}

/**
 * Logout
 * POST /api/v1/auth/logout
 */
export async function logout(): Promise<ApiResponse<{ message: string }>> {
  return apiCall<ApiResponse<{ message: string }>>('/auth/logout', 'POST', {}, true)
}

/**
 * Get current user
 * GET /api/v1/auth/me
 */
export async function getCurrentUser(): Promise<ApiResponse<{user: User}>> {
  return apiCall<ApiResponse<{user: User}>>('/auth/me', 'GET', undefined, true)
}

/**
 * Get linked accounts
 * GET /api/v1/auth/accounts
 */
export async function getLinkedAccounts(): Promise<ApiResponse<{ accounts: LinkedAccount[] }>> {
  return apiCall<ApiResponse<{ accounts: LinkedAccount[] }>>('/auth/accounts', 'GET', undefined, true)
}

/**
 * Link OAuth to account
 * POST /api/v1/auth/oauth/link
 */
export async function linkOAuth(data: LinkOAuthRequest): Promise<ApiResponse<{
  linked_account: LinkedAccount
}>> {
  return apiCall<ApiResponse<{ linked_account: LinkedAccount }>>('/auth/oauth/link', 'POST', data, true)
}

/**
 * Unlink OAuth from account
 * DELETE /api/v1/auth/oauth/unlink
 */
export async function unlinkOAuth(data: UnlinkOAuthRequest): Promise<ApiResponse<{ message: string }>> {
  return apiCall<ApiResponse<{ message: string }>>('/auth/oauth/unlink', 'DELETE', data, true)
}

/**
 * Bind email to account
 * POST /api/v1/auth/bind-email
 */
export async function bindEmail(data: BindEmailRequest): Promise<ApiResponse<{ user: User }>> {
  return apiCall<ApiResponse<{ user: User }>>('/auth/bind-email', 'POST', data, true)
}

/**
 * Set primary account
 * PUT /api/v1/auth/primary-account
 */
export async function setPrimaryAccount(data: SetPrimaryAccountRequest): Promise<
  ApiResponse<{ message: string }>
> {
  return apiCall<ApiResponse<{ message: string }>>('/auth/primary-account', 'PUT', data, true)
}

/**
 * Update profile
 * PUT /api/v1/auth/profile
 */
export async function updateProfile(data: UpdateProfileRequest): Promise<ApiResponse<{ user: User }>> {
  return apiCall<ApiResponse<{ user: User }>>('/auth/profile', 'PUT', data, true)
}

// Export API object
export const authApi = {
  getAuthConfig,
  sendCode,
  login,
  register,
  oauthAuthorize,
  oauthCallback,
  refreshToken,
  logout,
  getCurrentUser,
  getLinkedAccounts,
  linkOAuth,
  unlinkOAuth,
  bindEmail,
  setPrimaryAccount,
  updateProfile,
}
