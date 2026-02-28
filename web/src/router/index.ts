import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  // Public pages
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/home/index.vue'),
    meta: { title: 'iBookFS - 将纸质书籍转化为数字资产', requiresAuth: false },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { title: '登录 - iBookFS', requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/register/index.vue'),
    meta: { title: '注册 - iBookFS', requiresAuth: false },
  },
  // OAuth callback page with provider in path
  {
    path: '/auth/:provider/callback',
    name: 'OAuthCallback',
    component: () => import('@/pages/auth/callback/index.vue'),
    meta: { title: 'OAuth 登录 - iBookFS', requiresAuth: false },
  },

  // App pages (require authentication)
  {
    path: '/library',
    name: 'Library',
    component: () => import('@/pages/library/index.vue'),
    meta: { title: '我的书架', requiresAuth: true },
  },
  {
    path: '/book/:id',
    name: 'BookDetail',
    component: () => import('@/pages/book/[id]/index.vue'),
    meta: { title: '书籍详情', requiresAuth: true },
  },
  {
    path: '/upload/:bookId',
    name: 'UploadToBook',
    component: () => import('@/pages/upload/[bookId]/index.vue'),
    meta: { title: '上传照片', requiresAuth: true },
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/pages/settings/index.vue'),
    meta: { title: '设置', requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/pages/profile/index.vue'),
    meta: { title: '个人中心', requiresAuth: true },
  },

  // Catch all - redirect to landing
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
})

// Update page title
router.beforeEach(async (to, _from, next) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - iBookFS` : 'iBookFS'

  // Check authentication
  const requiresAuth = to.meta.requiresAuth as boolean | undefined
  if (requiresAuth) {
    const authStore = useAuthStore()

    // Initialize auth store if not already initialized
    if (!authStore.isAuthenticated) {
      await authStore.initialize()
    }

    // Redirect to login if not authenticated
    if (!authStore.isAuthenticated) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }
  }

  next()
})

export default router
