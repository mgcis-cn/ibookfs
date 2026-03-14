import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { useToastStore } from '@/shared/stores/toast'
import { ArrowLeft, Settings, Library, User, Shield, ChevronRight } from 'lucide-react'
import './ProfilePage.css'

export default function ProfilePage() {
  const navigate = useNavigate()
  const authStore = useAuthStore()
  const toast = useToastStore(s => s.add)

  const [isLoggingOut, setIsLoggingOut] = useState(false)

  const displayName = authStore.userDisplayName()
  const linkedCount = authStore.linkedProviders()?.length || 0
  const stats = { totalBooks: 0, totalPages: 0, linkedAccounts: linkedCount }

  const handleLogout = async () => {
    setIsLoggingOut(true)
    try {
      await authStore.logout()
      toast('success', '已退出登录')
      navigate('/')
    } catch {
      toast('error', '退出失败，请重试')
    } finally {
      setIsLoggingOut(false)
    }
  }

  return (
    <div className="page-profile">
      {/* Back Button */}
      <button className="back-button" onClick={() => navigate(-1)}>
        <ArrowLeft size={20} />
        <span>返回</span>
      </button>

      {/* Profile Header */}
      <div className="profile-header">
        <div className="avatar-large">
          {displayName?.charAt(0) || 'U'}
        </div>
        <h1 className="profile-name">{displayName || '用户'}</h1>
        <p className="profile-email">{authStore.user?.email}</p>
      </div>

      {/* Profile Stats */}
      <div className="profile-stats">
        <div className="stat-item">
          <div className="stat-value">{stats.totalBooks}</div>
          <div className="stat-label">藏书</div>
        </div>
        <div className="stat-divider" />
        <div className="stat-item">
          <div className="stat-value">{stats.totalPages}</div>
          <div className="stat-label">总页数</div>
        </div>
        <div className="stat-divider" />
        <div className="stat-item">
          <div className="stat-value">{stats.linkedAccounts}</div>
          <div className="stat-label">绑定账号</div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="action-section">
        <h2 className="section-title">快捷操作</h2>
        <div className="action-grid">
          <Link to="/settings" className="action-card">
            <div className="action-icon" style={{ background: 'rgba(193, 123, 92, 0.15)', color: 'var(--color-accent)' }}>
              <Settings size={24} />
            </div>
            <span className="action-label">设置</span>
          </Link>
          <Link to="/library" className="action-card">
            <div className="action-icon" style={{ background: 'var(--color-success)', color: 'white' }}>
              <Library size={24} />
            </div>
            <span className="action-label">书架</span>
          </Link>
        </div>
      </div>

      {/* Account Actions */}
      <div className="account-section">
        <h2 className="section-title">账户</h2>
        <div className="action-list">
          <button className="action-item">
            <div className="action-item-icon"><User size={20} /></div>
            <div className="action-item-content">
              <span className="action-item-title">编辑资料</span>
              <span className="action-item-desc">修改个人信息</span>
            </div>
            <ChevronRight size={18} className="action-item-arrow" />
          </button>
          <button className="action-item">
            <div className="action-item-icon"><Shield size={20} /></div>
            <div className="action-item-content">
              <span className="action-item-title">账号安全</span>
              <span className="action-item-desc">密码与绑定管理</span>
            </div>
            <ChevronRight size={18} className="action-item-arrow" />
          </button>
        </div>
      </div>

      {/* Logout */}
      <div className="logout-section">
        <button className="logout-button" onClick={handleLogout} disabled={isLoggingOut}>
          {isLoggingOut ? (
            <svg className="spinner" width="20" height="20" viewBox="0 0 20 20">
              <circle cx="10" cy="10" r="8" stroke="currentColor" strokeWidth="2" fill="none" opacity="0.3" />
              <path d="M10 2a8 8 0 018 8" stroke="currentColor" strokeWidth="2" fill="none" />
            </svg>
          ) : (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
              <polyline points="16 17 21 12 16 7" />
              <line x1="21" y1="12" x2="9" y2="12" />
            </svg>
          )}
          <span>{isLoggingOut ? '退出中...' : '退出登录'}</span>
        </button>
      </div>
    </div>
  )
}
