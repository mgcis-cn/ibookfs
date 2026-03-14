import { create } from 'zustand'
import type {
  User,
  AuthConfig,
  OAuthProvider,
  LoginRequest,
  RegisterRequest,
  SendCodeRequest,
  LinkOAuthRequest,
  UnlinkOAuthRequest,
  BindEmailRequest,
  SetPrimaryAccountRequest,
  UpdateProfileRequest,
} from './types'
import { authApi } from './api'

interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isLoading: boolean
  error: string | null
  authConfig: AuthConfig | null

  // Computed-like getters
  isAuthenticated: () => boolean
  userDisplayName: () => string
  linkedProviders: () => string[]
  isEmailEnabled: () => boolean
  isOAuthEnabled: () => boolean
  isProviderConfigured: (provider: OAuthProvider) => boolean

  // Actions
  initialize: () => Promise<void>
  loadConfig: () => Promise<void>
  sendCode: (data: SendCodeRequest) => Promise<boolean>
  login: (data: LoginRequest) => Promise<boolean>
  register: (data: RegisterRequest) => Promise<boolean>
  getOAuthAuthorizeUrl: (provider: OAuthProvider, redirectUri: string) => Promise<{ url: string; state: string } | null>
  oauthCallback: (provider: OAuthProvider, code: string, state: string) => Promise<boolean>
  linkOAuth: (data: LinkOAuthRequest) => Promise<boolean>
  unlinkOAuth: (data: UnlinkOAuthRequest) => Promise<boolean>
  bindEmail: (data: BindEmailRequest) => Promise<boolean>
  setPrimaryAccount: (data: SetPrimaryAccountRequest) => Promise<boolean>
  updateProfile: (data: UpdateProfileRequest) => Promise<boolean>
  refreshUser: () => Promise<boolean>
  logout: () => Promise<void>
  clearError: () => void
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: null,
  refreshToken: null,
  isLoading: false,
  error: null,
  authConfig: null,

  isAuthenticated: () => !!get().user && !!get().token,
  userDisplayName: () => {
    const u = get().user
    if (!u) return ''
    return u.display_name || [u.first_name, u.last_name].filter(Boolean).join(' ') || u.email
  },
  linkedProviders: () =>
    (get().user?.linkedAccounts ?? [])
      .filter(a => a.provider !== 'email')
      .map(a => a.provider),

  isEmailEnabled: () => get().authConfig?.emailEnabled ?? true,
  isOAuthEnabled: () => get().authConfig?.oauthEnabled ?? false,
  isProviderConfigured: (provider: OAuthProvider) => {
    const config = get().authConfig
    if (!config) return false
    const oauthConfig = config.oauth?.[provider]
    return oauthConfig?.configured ?? config.oauthProviders.includes(provider)
  },

  clearError: () => set({ error: null }),

  initialize: async () => {
    const savedToken = localStorage.getItem('authToken')
    const savedUser = localStorage.getItem('currentUser')
    const savedRefresh = localStorage.getItem('refreshToken')

    if (!savedToken || !savedUser) return

    try {
      set({ token: savedToken, user: JSON.parse(savedUser), refreshToken: savedRefresh })
      const response = await authApi.getCurrentUser()
      if (response?.user) {
        set({ user: response.user })
        localStorage.setItem('currentUser', JSON.stringify(response.user))
      } else {
        set({ user: null, token: null, refreshToken: null })
        localStorage.removeItem('authToken')
        localStorage.removeItem('currentUser')
        localStorage.removeItem('refreshToken')
      }
    } catch {
      set({ user: null, token: null, refreshToken: null })
      localStorage.removeItem('authToken')
      localStorage.removeItem('currentUser')
      localStorage.removeItem('refreshToken')
    }
  },

  loadConfig: async () => {
    try {
      const data = await authApi.getAuthConfig()
      // Handle snake_case from backend
      const raw = data as unknown as Record<string, unknown>
      set({
        authConfig: {
          oauthProviders: (raw.oauth_providers ?? raw.oauthProviders ?? []) as string[],
          emailEnabled: (raw.email_enabled ?? raw.emailEnabled ?? true) as boolean,
          oauthEnabled: (raw.oauth_enabled ?? raw.oauthEnabled ?? false) as boolean,
          oauth: raw.oauth as AuthConfig['oauth'],
        },
      })
    } catch {
      set({
        authConfig: {
          oauthProviders: [],
          emailEnabled: true,
          oauthEnabled: false,
          oauth: {
            github: { configured: false },
            gitee: { configured: false },
            icloud: { configured: false },
            google: { configured: false },
            wechat: { configured: false },
          },
        },
      })
    }
  },

  sendCode: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.sendCode(data)
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '发送验证码失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  login: async (data) => {
    set({ error: null, isLoading: true })
    try {
      const res = await authApi.login(data)
      set({ user: res.user, token: res.token, refreshToken: res.refreshToken ?? null })
      localStorage.setItem('authToken', res.token)
      localStorage.setItem('currentUser', JSON.stringify(res.user))
      if (res.refreshToken) localStorage.setItem('refreshToken', res.refreshToken)
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '登录失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  register: async (data) => {
    set({ error: null, isLoading: true })
    try {
      const res = await authApi.register(data)
      set({ user: res.user, token: res.token, refreshToken: res.refreshToken ?? null })
      localStorage.setItem('authToken', res.token)
      localStorage.setItem('currentUser', JSON.stringify(res.user))
      if (res.refreshToken) localStorage.setItem('refreshToken', res.refreshToken)
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '注册失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  getOAuthAuthorizeUrl: async (provider, redirectUri) => {
    set({ error: null, isLoading: true })
    try {
      const data = await authApi.oauthAuthorize(provider, redirectUri)
      const raw = data as unknown as Record<string, string>
      return { url: raw['authorize_url'] || raw['authorizeUrl'] || '', state: raw['state'] || '' }
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '获取授权链接失败' })
      return null
    } finally {
      set({ isLoading: false })
    }
  },

  oauthCallback: async (provider, code, state) => {
    set({ error: null, isLoading: true })
    try {
      const res = await authApi.oauthCallback(provider, code, state)
      // Handle nested Data from backend
      const authData = (res as unknown as Record<string, unknown>)['Data'] as typeof res ?? res
      set({ user: authData.user, token: authData.token, refreshToken: authData.refreshToken ?? null })
      localStorage.setItem('authToken', authData.token)
      localStorage.setItem('currentUser', JSON.stringify(authData.user))
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : 'OAuth 登录失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  linkOAuth: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.linkOAuth(data)
      await get().refreshUser()
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '绑定失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  unlinkOAuth: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.unlinkOAuth(data)
      await get().refreshUser()
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '解绑失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  bindEmail: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.bindEmail(data)
      await get().refreshUser()
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '绑定邮箱失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  setPrimaryAccount: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.setPrimaryAccount(data)
      await get().refreshUser()
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '设置主账号失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  updateProfile: async (data) => {
    set({ error: null, isLoading: true })
    try {
      await authApi.updateProfile(data)
      await get().refreshUser()
      return true
    } catch (err) {
      set({ error: err instanceof Error ? err.message : '更新资料失败' })
      return false
    } finally {
      set({ isLoading: false })
    }
  },

  refreshUser: async () => {
    if (!get().isAuthenticated()) return false
    try {
      const res = await authApi.getCurrentUser()
      if (res?.user) {
        set({ user: res.user })
        localStorage.setItem('currentUser', JSON.stringify(res.user))
        return true
      }
      return false
    } catch {
      return false
    }
  },

  logout: async () => {
    set({ isLoading: true })
    try {
      await authApi.logout()
    } finally {
      set({ user: null, token: null, refreshToken: null, isLoading: false })
      localStorage.removeItem('authToken')
      localStorage.removeItem('currentUser')
      localStorage.removeItem('refreshToken')
    }
  },
}))
