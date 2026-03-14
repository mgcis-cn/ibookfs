import { useEffect, useState } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { Loader2 } from 'lucide-react'

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const { isAuthenticated, initialize } = useAuthStore()
  // Skip async check if already authenticated (avoids flash on route changes)
  const [checking, setChecking] = useState(!isAuthenticated())

  useEffect(() => {
    if (!checking) return
    let cancelled = false
    const check = async () => {
      await initialize()
      if (!cancelled) setChecking(false)
    }
    check()
    return () => { cancelled = true }
  }, [checking])

  if (checking) {
    return (
      <div style={{ minHeight: '60vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Loader2 size={32} style={{ color: 'var(--color-accent)', animation: 'spin 1s linear infinite' }} />
      </div>
    )
  }

  if (!isAuthenticated()) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />
  }

  return <>{children}</>
}
