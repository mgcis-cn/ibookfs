import { useState, useEffect, useRef, useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { useToastStore } from '@/shared/stores/toast'
import { BookOpen, ArrowLeft, Mail, Link as LinkIcon, Shield } from 'lucide-react'
import type { OAuthProvider } from '@/features/auth/types'
import './LoginPage.css'

const ALL_OAUTH_PROVIDERS: { id: OAuthProvider; name: string }[] = [
  { id: 'github', name: 'GitHub' },
  { id: 'gitee', name: 'Gitee' },
  { id: 'google', name: 'Google' },
  { id: 'wechat', name: 'WeChat' },
]

export default function LoginPage() {
  const navigate = useNavigate()
  const auth = useAuthStore()
  const toast = useToastStore(s => s.add)

  const [activeTab, setActiveTab] = useState<'email' | 'oauth'>('email')
  const [step, setStep] = useState(1)
  const [email, setEmail] = useState('')
  const [codeDigits, setCodeDigits] = useState(['', '', '', '', '', ''])
  const [isLoading, setIsLoading] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [showOAuthModal, setShowOAuthModal] = useState(false)
  const [oauthLoading, setOauthLoading] = useState<Record<string, boolean>>({})

  const codeInputs = useRef<(HTMLInputElement | null)[]>([])
  const oauthWindowRef = useRef<Window | null>(null)

  const oauthProviders = ALL_OAUTH_PROVIDERS.filter(p => auth.isProviderConfigured(p.id))

  useEffect(() => {
    const init = async () => {
      if (!auth.isAuthenticated()) await auth.initialize()
      if (auth.isAuthenticated()) { navigate('/library'); return }
      await auth.loadConfig()
    }
    init()

    const handleMsg = (event: MessageEvent) => {
      if (event.origin !== window.location.origin) return
      if (event.data.type === 'oauth_callback') {
        handleOAuthCallback(event.data.provider, event.data.code, event.data.state)
      }
    }
    window.addEventListener('message', handleMsg)
    return () => {
      window.removeEventListener('message', handleMsg)
      if (oauthWindowRef.current && !oauthWindowRef.current.closed) oauthWindowRef.current.close()
    }
  }, [])

  useEffect(() => {
    if (!auth.isEmailEnabled() && auth.isOAuthEnabled()) setActiveTab('oauth')
  }, [auth.isEmailEnabled, auth.isOAuthEnabled])

  useEffect(() => {
    if (countdown <= 0) return
    const t = setInterval(() => setCountdown(c => { if (c <= 1) { clearInterval(t); return 0 } return c - 1 }), 1000)
    return () => clearInterval(t)
  }, [countdown])

  const handleCodeInput = useCallback(async (index: number, e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    if (!/^\d*$/.test(value)) { e.target.value = ''; return }
    const newDigits = [...codeDigits]
    newDigits[index] = value.slice(-1)
    setCodeDigits(newDigits)
    if (value && index < 5) codeInputs.current[index + 1]?.focus()
    if (newDigits.every(d => d !== '')) {
      setErrors(prev => ({ ...prev, code: '' }))
      submitLogin(newDigits)
    }
  }, [codeDigits])

  const handleCodeKeydown = useCallback((index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace') {
      if (!(e.target as HTMLInputElement).value && index > 0) codeInputs.current[index - 1]?.focus()
      else { const nd = [...codeDigits]; nd[index] = ''; setCodeDigits(nd) }
    }
    if (e.key === 'ArrowLeft' && index > 0) { e.preventDefault(); codeInputs.current[index - 1]?.focus() }
    if (e.key === 'ArrowRight' && index < 5) { e.preventDefault(); codeInputs.current[index + 1]?.focus() }
  }, [codeDigits])

  const handleCodePaste = useCallback((e: React.ClipboardEvent) => {
    e.preventDefault()
    const digits = (e.clipboardData.getData('text') || '').replace(/\D/g, '').slice(0, 6)
    const newDigits = [...codeDigits]
    for (let i = 0; i < digits.length; i++) newDigits[i] = digits[i] ?? ''
    setCodeDigits(newDigits)
    codeInputs.current[Math.min(digits.length, 5)]?.focus()
    if (newDigits.every(d => d !== '')) submitLogin(newDigits)
  }, [codeDigits])

  const submitLogin = async (digits: string[]) => {
    const codeValue = digits.join('')
    if (codeValue.length !== 6 || !/^\d{6}$/.test(codeValue)) {
      setErrors(prev => ({ ...prev, code: '请输入完整的 6 位验证码' }))
      return
    }
    setIsSubmitting(true)
    const success = await auth.login({ email, code: codeValue })
    setIsSubmitting(false)
    if (success) {
      toast('success', '登录成功')
      setTimeout(() => navigate('/library'), 500)
    } else {
      setErrors(prev => ({ ...prev, code: auth.error || '登录失败' }))
      setCodeDigits(['', '', '', '', '', ''])
      setTimeout(() => codeInputs.current[0]?.focus(), 50)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    if (step === 1) {
      if (!email) { setErrors({ email: '请输入邮箱地址' }); return }
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) { setErrors({ email: '请输入有效的邮箱地址' }); return }
      setIsLoading(true)
      const ok = await auth.sendCode({ email, type: 'login' })
      setIsLoading(false)
      if (ok) {
        setStep(2)
        setCountdown(60)
        toast('success', '验证码已发送')
        setTimeout(() => codeInputs.current[0]?.focus(), 50)
      } else {
        setErrors({ email: auth.error || '发送失败' })
      }
    }
  }

  const handleResend = async () => {
    const ok = await auth.sendCode({ email, type: 'login' })
    if (ok) {
      toast('success', '验证码已重新发送')
      setCountdown(60)
      setCodeDigits(['', '', '', '', '', ''])
      setErrors(prev => ({ ...prev, code: '' }))
      setTimeout(() => codeInputs.current[0]?.focus(), 50)
    } else {
      toast('error', auth.error || '发送失败')
    }
  }

  const handleOAuthLogin = async (provider: OAuthProvider) => {
    setOauthLoading(prev => ({ ...prev, [provider]: true }))
    try {
      const redirectUri = `${window.location.origin}/auth/${provider}/callback`
      const result = await auth.getOAuthAuthorizeUrl(provider, redirectUri)
      if (!result?.url) { toast('error', auth.error || '获取授权链接失败'); return }
      localStorage.setItem('oauth_state', result.state)
      setShowOAuthModal(true)
      oauthWindowRef.current = window.open(result.url, 'oauth', 'width=600,height=700,scrollbars=yes')
      if (!oauthWindowRef.current) { toast('error', '弹窗被浏览器阻止，请允许弹窗后重试'); setShowOAuthModal(false); return }
      const checkPopup = setInterval(() => {
        if (oauthWindowRef.current?.closed) { clearInterval(checkPopup); setShowOAuthModal(false); setOauthLoading(prev => ({ ...prev, [provider]: false })) }
      }, 500)
      setTimeout(() => { if (oauthWindowRef.current && !oauthWindowRef.current.closed) { oauthWindowRef.current.close(); setShowOAuthModal(false); toast('error', '授权超时，请稍后再试') } }, 120000)
    } catch { setShowOAuthModal(false); toast('error', '获取授权链接失败') }
    finally { setOauthLoading(prev => ({ ...prev, [provider]: false })) }
  }

  const oauthProcessingRef = useRef(false)
  const handleOAuthCallback = async (provider: OAuthProvider, code: string, state: string) => {
    if (oauthProcessingRef.current) return
    oauthProcessingRef.current = true
    const savedState = localStorage.getItem('oauth_state')
    localStorage.removeItem('oauth_state')
    if (state !== savedState) { toast('error', '授权验证失败：状态不匹配'); oauthProcessingRef.current = false; return }
    const success = await auth.oauthCallback(provider, code, state)
    if (success) { toast('success', '登录成功'); setShowOAuthModal(false); setTimeout(() => navigate('/library'), 500) }
    else { toast('error', auth.error || 'OAuth 登录失败'); oauthProcessingRef.current = false }
    setShowOAuthModal(false)
  }

  return (
    <div className="auth-page">
      {/* Left Brand Section */}
      <div className="auth-brand">
        <div className="brand-content">
          <Link to="/" className="brand-logo">
            <BookOpen size={32} />
            <span className="logo-text">iBookFS</span>
          </Link>
          <h1 className="brand-title">欢迎回来</h1>
          <p className="brand-description">继续管理你的数字藏书，让知识触手可及</p>
          <div className="brand-quote">
            <p className="quote-text">"iBookFS 让我的书籍管理变得如此简单，支持多种登录方式，体验非常流畅。"</p>
            <div className="quote-author">
              <div className="author-avatar">张</div>
              <div>
                <p className="author-name">张明</p>
                <p className="author-title">独立研究者</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Right Form Section */}
      <div className="auth-form-container">
        <div className="auth-form-wrapper">
          <Link to="/" className="back-link">
            <ArrowLeft size={16} />
            返回首页
          </Link>

          <div className="form-header">
            <h2 className="form-title">登录账号</h2>
            <p className="form-subtitle">选择您喜欢的登录方式</p>
          </div>

          {/* Tab Switcher */}
          {auth.isEmailEnabled() && auth.isOAuthEnabled() && (
            <div className="auth-tabs">
              <button className={`auth-tab ${activeTab === 'email' ? 'active' : ''}`} onClick={() => setActiveTab('email')}>
                <Mail size={18} /> 邮箱登录
              </button>
              <button className={`auth-tab ${activeTab === 'oauth' ? 'active' : ''}`} onClick={() => setActiveTab('oauth')}>
                <LinkIcon size={18} /> 第三方登录
              </button>
            </div>
          )}

          {/* Email Login Form */}
          {activeTab === 'email' && (
            <div className="auth-form-content">
              <form className="auth-form" onSubmit={handleSubmit}>
                <div className="form-group">
                  <label className="form-label">邮箱地址</label>
                  <div className="input-wrapper">
                    <input
                      type="email"
                      className={`form-input ${errors.email ? 'error' : email ? 'success' : ''}`}
                      placeholder="your@email.com"
                      value={email}
                      onChange={e => setEmail(e.target.value)}
                      disabled={step === 2}
                      required
                      autoComplete="email"
                    />
                    <Mail className="input-icon" size={20} />
                  </div>
                  {errors.email && <span className="form-error">{errors.email}</span>}
                </div>

                {step === 2 && (
                  <div className="form-group">
                    <label className="form-label">验证码</label>
                    <div className="code-digits-container">
                      {[0, 1, 2, 3, 4, 5].map(index => (
                        <input
                          key={index}
                          ref={el => { codeInputs.current[index] = el }}
                          type="text"
                          inputMode="numeric"
                          className={`code-digit-input ${errors.code ? 'error' : codeDigits.join('').length === 6 ? 'success' : ''}`}
                          maxLength={1}
                          value={codeDigits[index]}
                          onChange={e => handleCodeInput(index, e)}
                          onKeyDown={e => handleCodeKeydown(index, e)}
                          onPaste={handleCodePaste}
                          autoComplete="one-time-code"
                        />
                      ))}
                    </div>
                    <span className="form-hint">验证码已发送至 <strong>{email}</strong></span>
                    {errors.code && <span className="form-error">{errors.code}</span>}
                  </div>
                )}

                {step === 1 && (
                  <button type="submit" className="button-primary button-lg button-full" disabled={isLoading}>
                    {!isLoading ? <span>发送验证码</span> : (
                      <span className="btn-loading">
                        <svg className="spinner" width="20" height="20" viewBox="0 0 20 20">
                          <circle cx="10" cy="10" r="8" stroke="currentColor" strokeWidth="2" fill="none" opacity="0.25" />
                          <path d="M10 2a8 8 0 018 8" stroke="currentColor" strokeWidth="2" fill="none" />
                        </svg>
                        处理中...
                      </span>
                    )}
                  </button>
                )}

                {step === 2 && !isSubmitting && (
                  <button type="button" className="resend-button button-lg button-full" disabled={countdown > 0 || isLoading} onClick={handleResend}>
                    {countdown > 0 ? `${countdown}s 后重新发送` : '重新发送验证码'}
                  </button>
                )}

                {step === 2 && isSubmitting && (
                  <div className="submitting-state">
                    <svg className="spinner" width="24" height="24" viewBox="0 0 24 24">
                      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" fill="none" opacity="0.25" />
                      <path d="M12 2a10 10 0 0110 10" stroke="currentColor" strokeWidth="3" fill="none" />
                    </svg>
                    <span>验证中...</span>
                  </div>
                )}
              </form>
            </div>
          )}

          {/* OAuth Login */}
          {activeTab === 'oauth' && (
            <div className="oauth-content">
              <p className="oauth-description">使用您的第三方账号快速登录</p>
              <div className="oauth-buttons">
                {oauthProviders.map(provider => (
                  <button key={provider.id} className="oauth-button" disabled={oauthLoading[provider.id]} onClick={() => handleOAuthLogin(provider.id)}>
                    使用 {provider.name} 登录
                  </button>
                ))}
              </div>
              <div className="oauth-hint">
                <Shield size={16} />
                <span>我们尊重您的隐私，不会获取不必要的权限</span>
              </div>
            </div>
          )}

          <div className="divider"><span>或</span></div>

          <div className="form-footer">
            <p className="footer-text">
              还没有账号？
              <Link to="/register" className="footer-link">立即注册</Link>
            </p>
          </div>

          <p className="privacy-notice">
            登录即表示您同意我们的
            <a href="/terms">服务条款</a>
            和
            <a href="/privacy">隐私政策</a>
          </p>
        </div>
      </div>

      {/* OAuth Callback Modal */}
      {showOAuthModal && (
        <div className="oauth-modal-overlay" onClick={e => { if (e.target === e.currentTarget) setShowOAuthModal(false) }}>
          <div className="oauth-modal">
            <div className="oauth-modal-content">
              <div className="oauth-modal-icon">
                <svg className="spinner" width="48" height="48" viewBox="0 0 48 48">
                  <circle cx="24" cy="24" r="20" stroke="currentColor" strokeWidth="3" fill="none" opacity="0.25" />
                  <path d="M24 4a20 20 0 0120 20" stroke="currentColor" strokeWidth="3" fill="none" />
                </svg>
              </div>
              <h3 className="oauth-modal-title">正在连接...</h3>
              <p className="oauth-modal-description">请在新窗口中完成授权</p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
