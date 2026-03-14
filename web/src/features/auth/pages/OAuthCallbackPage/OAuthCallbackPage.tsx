import { useEffect, useState } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { Loader2, AlertCircle } from 'lucide-react'
import type { OAuthProvider } from '@/features/auth/types'

export default function OAuthCallbackPage() {
  const { provider } = useParams<{ provider: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { oauthCallback, error } = useAuthStore()
  const [processing, setProcessing] = useState(true)

  useEffect(() => {
    const code = searchParams.get('code')
    const state = searchParams.get('state')

    if (!provider || !code || !state) {
      setProcessing(false)
      return
    }

    // If opened as a popup by LoginPage/RegisterPage, send data back to parent
    if (window.opener) {
      window.opener.postMessage(
        { type: 'oauth_callback', provider, code, state },
        window.location.origin,
      )
      window.close()
      return
    }

    // Direct navigation (not a popup) — handle callback directly
    const handle = async () => {
      const ok = await oauthCallback(provider as OAuthProvider, code, state)
      if (ok) {
        navigate('/library', { replace: true })
      } else {
        setProcessing(false)
      }
    }
    handle()
  }, [])

  if (processing && !error) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center gap-4">
        <Loader2 className="w-10 h-10 text-[var(--color-accent)] animate-spin" />
        <p className="text-[var(--color-ink-light)]">正在处理 OAuth 登录...</p>
      </div>
    )
  }

  return (
    <div className="min-h-[60vh] flex flex-col items-center justify-center gap-4">
      <AlertCircle className="w-10 h-10 text-[var(--color-error)]" />
      <p className="text-[var(--color-error)] font-medium">{error || 'OAuth 登录失败，请重试'}</p>
      <button
        onClick={() => navigate('/login')}
        className="px-6 py-2.5 rounded-lg bg-[var(--color-accent)] text-white font-medium hover:bg-[var(--color-accent-dark)] transition"
      >
        返回登录
      </button>
    </div>
  )
}
