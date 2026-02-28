/**
 * Pinia Auth Store
 *
 * Manages authentication state including OAuth provider support.
 * Handles login, registration, OAuth linking, and session management.
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
    User,
    OAuthProvider,
    AuthResponse,
    LoginRequest,
    RegisterRequest,
    SendCodeRequest,
    LinkOAuthRequest,
    UnlinkOAuthRequest,
    BindEmailRequest,
    SetPrimaryAccountRequest,
    UpdateProfileRequest,
    AuthConfig, ApiResponse,
} from '@/types'
import { authApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshTokenValue = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const authConfig = ref<AuthConfig | null>(null)

  // Computed
  const isAuthenticated = computed(() => !!user.value && !!token.value)
  const userEmail = computed(() => user.value?.email ?? '')
  const userDisplayName = computed(() => user.value?.display_name || `${user.value?.first_name} ${user.value?.last_name}` || '')
  const userAvatar = computed(() => user.value?.avatarUrl)
  const linkedAccounts = computed(() => user.value?.linkedAccounts ?? [])
  const primaryAccount = computed(() => linkedAccounts.value.find((a) => a.isPrimary))
  const hasEmail = computed(() => linkedAccounts.value.some((a) => a.provider === 'email'))
  const hasOAuth = computed(() =>
    linkedAccounts.value.some((a) => a.provider !== 'email')
  )
  const linkedProviders = computed(() =>
    linkedAccounts.value.filter((a) => a.provider !== 'email').map((a) => a.provider)
  )

  // Check if an OAuth provider is configured
  const isProviderConfigured = (provider: OAuthProvider): boolean => {
    return authConfig.value?.oauth[provider]?.configured ?? false
  }

  // Get supported OAuth providers
  const supportedProviders = computed(() => authConfig.value?.oauthProviders ?? [])

  // Check if email login is enabled
  const isEmailEnabled = computed(() => authConfig.value?.emailEnabled ?? true)

  // Check if OAuth login is enabled (any provider configured)
  const isOAuthEnabled = computed(() => authConfig.value?.oauthEnabled ?? false)

  // Actions
  const setError = (message: string | null) => {
    error.value = message
  }

  const clearError = () => {
    error.value = null
  }

  const setLoading = (loading: boolean) => {
    isLoading.value = loading
  }

  const setAuth = (authData: AuthResponse) => {
    user.value = authData.user
    token.value = authData.token
    refreshTokenValue.value = authData.refreshToken || null

    // Persist to localStorage
    localStorage.setItem('authToken', authData.token)
    localStorage.setItem('currentUser', JSON.stringify(authData.user))
    if (authData.refreshToken) {
      localStorage.setItem('refreshToken', authData.refreshToken)
    }
  }

  const clearAuth = () => {
    user.value = null
    token.value = null
    refreshTokenValue.value = null

    // Clear localStorage
    localStorage.removeItem('authToken')
    localStorage.removeItem('currentUser')
    localStorage.removeItem('refreshToken')
  }

  // Load auth configuration from server
  const loadConfig = async () => {
    try {
      const response = await authApi.getAuthConfig()
      if (response.success && response.data) {
        // Map snake_case API response to camelCase
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const data = response.data as any
        authConfig.value = {
          oauthProviders: data.oauth_providers || data.oauthProviders || [],
          emailEnabled: data.email_enabled ?? data.emailEnabled ?? true,
          oauthEnabled: data.oauth_enabled ?? data.oauthEnabled ?? false,
          oauth: data.oauth,
        }
      }
    } catch {
      // Use default config if API fails
      authConfig.value = {
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
      }
    }
  }

  // Initialize from localStorage
  const initialize = async () => {
    const savedToken = localStorage.getItem('authToken')
    const savedUser = localStorage.getItem('currentUser')
    const savedRefreshToken = localStorage.getItem('refreshToken')

    if (savedToken && savedUser) {
      try {
        token.value = savedToken
        user.value = JSON.parse(savedUser)
        refreshTokenValue.value = savedRefreshToken

        // Verify token is still valid by fetching current user
        const response = await authApi.getCurrentUser() as ApiResponse<{user: User}>
        if (response.success && response.data?.user) {
          user.value = response.data.user
        } else {
          clearAuth()
        }
      } catch {
        clearAuth()
      }
    }
  }

  // Send verification code
  const sendCode = async (data: SendCodeRequest) => {
    clearError()
    setLoading(true)

    try {
      console.log('[Auth] Sending code request:', data)
      const response = await authApi.sendCode(data)
      console.log('[Auth] Send code response:', response)

      if (!response.success) {
        setError(response.message || '发送验证码失败')
        return false
      }

      return true
    } catch (err) {
      console.error('[Auth] Send code error:', err)
      setError(err instanceof Error ? err.message : '发送验证码失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Login with email and code
  const login = async (data: LoginRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.login(data)

      if (!response.success) {
        setError(response.message || '登录失败')
        return false
      }

      setAuth(response.data)
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Register new account
  const register = async (data: RegisterRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.register(data)

      if (!response.success) {
        setError(response.message || '注册失败')
        return false
      }

      setAuth(response.data)
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '注册失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // OAuth - Get authorization URL
  const getOAuthAuthorizeUrl = async (
    provider: OAuthProvider,
    redirectUri: string
  ): Promise<{ url: string; state: string } | null> => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.oauthAuthorize(provider, redirectUri)

      if (!response.success) {
        setError(response.message || '获取授权链接失败')
        return null
      }

      // Handle both snake_case and camelCase field names
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.data as any
      return {
        url: data.authorize_url || data.authorizeUrl,
        state: data.state,
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取授权链接失败')
      return null
    } finally {
      setLoading(false)
    }
  }

  // OAuth - Handle callback
  const oauthCallback = async (
    provider: OAuthProvider,
    code: string,
    state: string
  ) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.oauthCallback(provider, code, state)

      if (!response.success) {
        setError(response.message || 'OAuth 登录失败')
        return false
      }

      // Handle nested Data structure from backend
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = (response.data as any)?.Data || response.data
      setAuth(data)
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : 'OAuth 登录失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Link OAuth to current account
  const linkOAuth = async (data: LinkOAuthRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.linkOAuth(data)

      if (!response.success) {
        setError(response.message || '绑定失败')
        return false
      }

      // Refresh user data
      await refreshUser()

      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '绑定失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Unlink OAuth from account
  const unlinkOAuth = async (data: UnlinkOAuthRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.unlinkOAuth(data)

      if (!response.success) {
        setError(response.message || '解绑失败')
        return false
      }

      // Refresh user data
      await refreshUser()

      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '解绑失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Bind email to account
  const bindEmail = async (data: BindEmailRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.bindEmail(data)

      if (!response.success) {
        setError(response.message || '绑定邮箱失败')
        return false
      }

      // Refresh user data
      await refreshUser()

      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '绑定邮箱失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Set primary account
  const setPrimaryAccount = async (data: SetPrimaryAccountRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.setPrimaryAccount(data)

      if (!response.success) {
        setError(response.message || '设置主账号失败')
        return false
      }

      // Refresh user data
      await refreshUser()

      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '设置主账号失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Refresh token
  const refreshAuthToken = async () => {
    if (!refreshTokenValue.value) return false

    try {
      const response = await authApi.refreshToken({ refreshToken: refreshTokenValue.value })

      if (response.success) {
        token.value = response.data.token
        if (response.data.refreshToken) {
          refreshTokenValue.value = response.data.refreshToken
        }
        localStorage.setItem('authToken', response.data.token)
        return true
      }

      clearAuth()
      return false
    } catch {
      clearAuth()
      return false
    }
  }

  // Refresh user data
  const refreshUser = async () => {
    if (!isAuthenticated.value) return false

    try {
      const response = await authApi.getCurrentUser() as ApiResponse<{user: User}>

      if (response.success && response.data?.user) {
        user.value = response.data.user
        localStorage.setItem('currentUser', JSON.stringify(response.data.user))
        return true
      }

      return false
    } catch {
      return false
    }
  }

  // Update profile
  const updateProfile = async (data: UpdateProfileRequest) => {
    clearError()
    setLoading(true)

    try {
      const response = await authApi.updateProfile(data)

      if (!response.success) {
        setError(response.message || '更新资料失败')
        return false
      }

      await refreshUser()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新资料失败')
      return false
    } finally {
      setLoading(false)
    }
  }

  // Get linked accounts
  const getLinkedAccountsList = async () => {
    if (!isAuthenticated.value) return []

    try {
      const response = await authApi.getLinkedAccounts()

      if (response.success && response.data) {
        return response.data.accounts
      }

      return []
    } catch {
      return []
    }
  }

  // Logout
  const logout = async () => {
    setLoading(true)

    try {
      await authApi.logout()
    } finally {
      clearAuth()
      setLoading(false)
    }
  }

  return {
    // State
    user,
    token,
    refreshToken: refreshTokenValue,
    isLoading,
    error,
    authConfig,

    // Computed
    isAuthenticated,
    userEmail,
    userDisplayName,
    userAvatar,
    linkedAccounts,
    primaryAccount,
    hasEmail,
    hasOAuth,
    linkedProviders,
    isProviderConfigured,
    supportedProviders,
    isEmailEnabled,
    isOAuthEnabled,

    // Actions
    loadConfig,
    initialize,
    sendCode,
    login,
    register,
    getOAuthAuthorizeUrl,
    oauthCallback,
    linkOAuth,
    unlinkOAuth,
    bindEmail,
    setPrimaryAccount,
    refreshAuthToken,
    refreshUser,
    updateProfile,
    getLinkedAccountsList,
    logout,
    clearError,
  }
})
