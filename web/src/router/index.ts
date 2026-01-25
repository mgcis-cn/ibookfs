import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  // Public pages
  {
    path: '/',
    name: 'Landing',
    component: () => import('@/pages/LandingPage.vue'),
    meta: { title: 'iBookFS - 将纸质书籍转化为数字资产' },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/LoginPage.vue'),
    meta: { title: '登录 - iBookFS' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/RegisterPage.vue'),
    meta: { title: '注册 - iBookFS' },
  },

  // App pages (require authentication)
  {
    path: '/library',
    name: 'Library',
    component: () => import('@/pages/LibraryPage.vue'),
    meta: { title: '我的书架' },
  },
  {
    path: '/book/:id',
    name: 'BookDetail',
    component: () => import('@/pages/BookDetailPage.vue'),
    meta: { title: '书籍详情' },
  },
  {
    path: '/upload/:bookId',
    name: 'UploadToBook',
    component: () => import('@/pages/UploadPage.vue'),
    meta: { title: '上传照片' },
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/pages/SettingsPage.vue'),
    meta: { title: '设置' },
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
router.beforeEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - iBookFS` : 'iBookFS'
})

export default router
