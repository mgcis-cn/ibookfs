import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
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
