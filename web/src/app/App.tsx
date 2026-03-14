import { Routes, Route } from 'react-router-dom'
import { AppLayout } from '@/shared/components/Layout/Layout'
import { AuthGuard } from '@/features/auth/components/AuthGuard'
import { ToastContainer } from '@/shared/components/Toast/Toast'

import HomePage from '@/features/home/pages/HomePage/HomePage'
import LoginPage from '@/features/auth/pages/LoginPage/LoginPage'
import RegisterPage from '@/features/auth/pages/RegisterPage/RegisterPage'
import OAuthCallbackPage from '@/features/auth/pages/OAuthCallbackPage/OAuthCallbackPage'
import LibraryPage from '@/features/library/pages/LibraryPage/LibraryPage'
import BookDetailPage from '@/features/library/pages/BookDetailPage/BookDetailPage'
import UploadPage from '@/features/upload/pages/UploadPage/UploadPage'
import SettingsPage from '@/features/settings/pages/SettingsPage/SettingsPage'
import ProfilePage from '@/features/profile/pages/ProfilePage/ProfilePage'

export default function App() {
  return (
    <>
      <ToastContainer />
      <Routes>
        <Route element={<AppLayout />}>
          {/* Public routes */}
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/auth/:provider/callback" element={<OAuthCallbackPage />} />

          {/* Protected routes */}
          <Route path="/library" element={<AuthGuard><LibraryPage /></AuthGuard>} />
          <Route path="/book/:id" element={<AuthGuard><BookDetailPage /></AuthGuard>} />
          <Route path="/upload/:bookId" element={<AuthGuard><UploadPage /></AuthGuard>} />
          <Route path="/settings" element={<AuthGuard><SettingsPage /></AuthGuard>} />
          <Route path="/profile" element={<AuthGuard><ProfilePage /></AuthGuard>} />
        </Route>
      </Routes>
    </>
  )
}
