export type OAuthProvider = 'github' | 'gitee' | 'icloud' | 'google' | 'wechat'

export interface User {
  id: string
  email: string
  emailVerified: boolean
  first_name?: string
  last_name?: string
  display_name?: string
  avatarUrl?: string
  status: 'active' | 'suspended' | 'deleted'
  role: 'user' | 'admin' | 'super_admin'
  preferences?: {
    theme?: 'light' | 'dark' | 'system'
    notifications?: { email?: boolean; push?: boolean }
  }
  language: string
  timezone: string
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
  linkedAccounts?: LinkedAccount[]
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
