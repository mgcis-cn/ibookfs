import { apiCall } from '@/shared/api/client'
import type {
  AuthResponse,
  AuthConfig,
  User,
  LinkedAccount,
  SendCodeRequest,
  LoginRequest,
  RegisterRequest,
  LinkOAuthRequest,
  UnlinkOAuthRequest,
  BindEmailRequest,
  SetPrimaryAccountRequest,
  RefreshTokenRequest,
  UpdateProfileRequest,
} from './types'

export const authApi = {
  getAuthConfig: () =>
    apiCall<AuthConfig>('/auth/config', 'GET', undefined, false),

  sendCode: (data: SendCodeRequest) =>
    apiCall<{ message: string }>('/auth/send-code', 'POST', data, false),

  login: (data: LoginRequest) =>
    apiCall<AuthResponse>('/auth/login', 'POST', data, false),

  register: (data: RegisterRequest) =>
    apiCall<AuthResponse>('/auth/register', 'POST', data, false),

  oauthAuthorize: (provider: string, redirectUri: string) =>
    apiCall<{ authorizeUrl: string; state: string }>(
      `/auth/oauth/authorize?provider=${provider}&redirect_uri=${encodeURIComponent(redirectUri)}`,
      'GET',
      undefined,
      false,
    ),

  oauthCallback: (provider: string, code: string, state: string) =>
    apiCall<AuthResponse>('/auth/oauth/callback', 'POST', { provider, code, state }, false),

  refreshToken: (data: RefreshTokenRequest) =>
    apiCall<{ token: string; refreshToken?: string; expiresIn: number }>(
      '/auth/refresh', 'POST', data, false,
    ),

  logout: () =>
    apiCall<{ message: string }>('/auth/logout', 'POST', {}),

  getCurrentUser: () =>
    apiCall<{ user: User }>('/auth/me'),

  getLinkedAccounts: () =>
    apiCall<{ accounts: LinkedAccount[] }>('/auth/accounts'),

  linkOAuth: (data: LinkOAuthRequest) =>
    apiCall<{ linked_account: LinkedAccount }>('/auth/oauth/link', 'POST', data),

  unlinkOAuth: (data: UnlinkOAuthRequest) =>
    apiCall<{ message: string }>('/auth/oauth/unlink', 'DELETE', data),

  bindEmail: (data: BindEmailRequest) =>
    apiCall<{ user: User }>('/auth/bind-email', 'POST', data),

  setPrimaryAccount: (data: SetPrimaryAccountRequest) =>
    apiCall<{ message: string }>('/auth/accounts/primary', 'PUT', data),

  updateProfile: (data: UpdateProfileRequest) =>
    apiCall<{ user: User }>('/auth/me', 'PUT', data),
}
