<template>
  <div class="app" :class="{ 'mobile': isMobile, 'no-nav': !showNavigation }">
    <!-- Desktop Sidebar -->
    <aside v-if="!isMobile && showNavigation" class="sidebar">
      <div class="sidebar-header">
        <router-link to="/library" class="sidebar-logo">
          <BookOpen :size="24" />
          <span>iBookFS</span>
        </router-link>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/library" class="nav-item" :class="{ active: currentRoute === '/library' }">
          <Library :size="20" />
          <span>书架</span>
        </router-link>
        <router-link to="/settings" class="nav-item" :class="{ active: currentRoute === '/settings' }">
          <Settings :size="20" />
          <span>设置</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-profile">
          <User :size="20" />
          <span>{{ userDisplayName }}</span>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content" :class="{ 'full-width': !showNavigation }">
      <router-view />
    </main>

    <!-- Mobile Bottom Navigation -->
    <nav v-if="isMobile && showNavigation" class="bottom-nav">
      <router-link to="/library" class="bottom-nav-item" :class="{ active: currentRoute === '/library' }">
        <Library :size="24" />
        <span>书架</span>
      </router-link>
      <router-link to="/settings" class="bottom-nav-item" :class="{ active: currentRoute === '/settings' }">
        <Settings :size="24" />
        <span>设置</span>
      </router-link>
    </nav>

    <!-- Toast Container -->
    <Teleport to="body">
      <div class="toast-container">
        <Toast
          v-for="toast in toasts"
          :key="toast.id"
          :type="toast.type"
          :message="toast.message"
          @close="removeToast(toast.id)"
        />
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  BookOpen,
  Library,
  Settings,
  User
} from 'lucide-vue-next'
import Toast from '@/components/ui/Toast.vue'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const toastStore = useToastStore()

const currentRoute = computed(() => route.path)
const toasts = computed(() => toastStore.toasts)

// Public routes that don't show navigation
const publicRoutes = ['/', '/login', '/register']
const showNavigation = computed(() => !publicRoutes.includes(currentRoute.value))

// Get user display name from localStorage
const userDisplayName = computed(() => {
  const userStr = localStorage.getItem('currentUser')
  if (userStr) {
    const user = JSON.parse(userStr)
    return user.firstName || '用户'
  }
  return '用户'
})

const isMobile = ref(false)
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

const removeToast = (id: number) => {
  toastStore.removeToast(id)
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.app {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
  color: var(--color-ink);
}

.sidebar {
  width: 240px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-elevated);
  border-right: 1px solid var(--border-subtle);
  position: fixed;
  left: 0;
  top: 0;
}

.sidebar-header {
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-accent);
  text-decoration: none;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.sidebar-logo:hover {
  opacity: 0.9;
}

.sidebar-nav {
  flex: 1;
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-ink-light);
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
}

.nav-item:hover {
  background: var(--bg-secondary);
  color: var(--color-ink);
}

.nav-item.active {
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.sidebar-footer {
  padding: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

.user-profile {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  color: var(--color-ink-light);
}

.main-content {
  flex: 1;
  margin-left: 240px;
  min-height: 100vh;
}

.main-content.full-width {
  margin-left: 0;
}

.app.mobile .main-content {
  margin-left: 0;
  padding-bottom: 64px;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  background: var(--bg-elevated);
  border-top: 1px solid var(--border-subtle);
  padding: var(--space-2);
  z-index: var(--z-fixed);
}

.bottom-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2);
  color: var(--color-ink-light);
  text-decoration: none;
  transition: color var(--duration-fast) var(--ease-out);
}

.bottom-nav-item span {
  font-size: var(--text-xs);
}

.bottom-nav-item.active {
  color: var(--color-accent);
}

.toast-container {
  position: fixed;
  top: var(--space-6);
  right: var(--space-6);
  z-index: var(--z-toast);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

@media (max-width: 768px) {
  .toast-container {
    left: var(--space-4);
    right: var(--space-4);
    top: var(--space-4);
  }
}
</style>
