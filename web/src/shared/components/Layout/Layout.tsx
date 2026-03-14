import { useState, useEffect } from 'react'
import { Link, useLocation, Outlet } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/store'
import { Library, Settings, User } from 'lucide-react'
import './Layout.css'

export function AppLayout() {
  const location = useLocation()
  const publicRoutes = ['/', '/login', '/register']
  const showNavigation = !publicRoutes.includes(location.pathname) && !location.pathname.startsWith('/auth/')

  const [isMobile, setIsMobile] = useState(window.innerWidth < 768)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768)
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  if (!showNavigation) {
    return <Outlet />
  }

  return (
    <div className={`app-layout${isMobile ? ' mobile' : ''}`}>
      {/* Desktop Sidebar */}
      {!isMobile && (
        <Sidebar />
      )}

      {/* Main Content */}
      <main className="main-content">
        <Outlet />
      </main>

      {/* Mobile Bottom Navigation */}
      {isMobile && (
        <MobileBottomNav />
      )}
    </div>
  )
}

function Sidebar() {
  const location = useLocation()
  const { userDisplayName } = useAuthStore()

  const navItems = [
    { path: '/library', label: '书架', icon: Library },
    { path: '/settings', label: '设置', icon: Settings },
  ]

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <Link to="/" className="sidebar-logo">
          <img src="/logo.svg" alt="iBookFS" style={{ height: 32 }} />
        </Link>
      </div>

      <nav className="sidebar-nav">
        {navItems.map(item => {
          const Icon = item.icon
          const active = location.pathname === item.path
          return (
            <Link key={item.path} to={item.path} className={`nav-item${active ? ' active' : ''}`}>
              <Icon size={20} />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </nav>

      <div className="sidebar-footer">
        <Link to="/profile" className={`user-profile${location.pathname === '/profile' ? ' active' : ''}`}>
          <User size={20} />
          <span>{userDisplayName() || '用户'}</span>
        </Link>
      </div>
    </aside>
  )
}

function MobileBottomNav() {
  const location = useLocation()

  const navItems = [
    { path: '/library', label: '书架', icon: Library },
    { path: '/settings', label: '设置', icon: Settings },
  ]

  return (
    <nav className="bottom-nav">
      {navItems.map(item => {
        const Icon = item.icon
        const active = location.pathname === item.path
        return (
          <Link key={item.path} to={item.path} className={`bottom-nav-item${active ? ' active' : ''}`}>
            <Icon size={24} />
            <span>{item.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
