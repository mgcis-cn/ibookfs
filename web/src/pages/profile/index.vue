<template>
  <div class="page-profile">
    <!-- Back Button -->
    <button class="back-button" @click="router.back()">
      <ArrowLeft :size="20" />
      <span>返回</span>
    </button>

    <!-- Profile Header -->
    <div class="profile-header">
      <div class="avatar-large">
        {{ authStore.userDisplayName?.charAt(0) || 'U' }}
      </div>
      <h1 class="profile-name">{{ authStore.userDisplayName || '用户' }}</h1>
      <p class="profile-email">{{ authStore.userEmail }}</p>
    </div>

    <!-- Profile Stats -->
    <div class="profile-stats">
      <div class="stat-item">
        <div class="stat-value">{{ stats.totalBooks }}</div>
        <div class="stat-label">藏书</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-value">{{ stats.totalPages }}</div>
        <div class="stat-label">总页数</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-value">{{ stats.linkedAccounts }}</div>
        <div class="stat-label">绑定账号</div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="action-section">
      <h2 class="section-title">快捷操作</h2>
      <div class="action-grid">
        <router-link to="/settings" class="action-card">
          <div class="action-icon" style="background: var(--color-accent-light); color: var(--color-accent);">
            <Settings :size="24" />
          </div>
          <span class="action-label">设置</span>
        </router-link>
        <router-link to="/library" class="action-card">
          <div class="action-icon" style="background: var(--color-success); color: white;">
            <Library :size="24" />
          </div>
          <span class="action-label">书架</span>
        </router-link>
      </div>
    </div>

    <!-- Account Actions -->
    <div class="account-section">
      <h2 class="section-title">账户</h2>
      <div class="action-list">
        <button class="action-item">
          <div class="action-item-icon">
            <User :size="20" />
          </div>
          <div class="action-item-content">
            <span class="action-item-title">编辑资料</span>
            <span class="action-item-desc">修改个人信息</span>
          </div>
          <ChevronRight :size="18" class="action-item-arrow" />
        </button>
        <button class="action-item">
          <div class="action-item-icon">
            <Shield :size="20" />
          </div>
          <div class="action-item-content">
            <span class="action-item-title">账号安全</span>
            <span class="action-item-desc">密码与绑定管理</span>
          </div>
          <ChevronRight :size="18" class="action-item-arrow" />
        </button>
      </div>
    </div>

    <!-- Logout Button -->
    <div class="logout-section">
      <button @click="handleLogout" class="logout-button" :disabled="isLoggingOut">
        <svg v-if="!isLoggingOut" class="logout-icon" width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2h2" />
          <polyline points="16 17 5 12 13 8 17" />
        </svg>
        <svg v-else class="spinner" width="20" height="20" viewBox="0 0 20 20">
          <circle cx="10" cy="10" r="8" stroke="currentColor" stroke-width="2" fill="none" opacity="0.3" />
          <path d="M10 2a8 8 0 018 8" stroke="currentColor" stroke-width="2" fill="none" />
        </svg>
        <span>{{ isLoggingOut ? '退出中...' : '退出登录' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import {
  ArrowLeft,
  Settings,
  Library,
  User,
  Shield,
  ChevronRight
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

const isLoggingOut = ref(false)
const stats = ref({
  totalBooks: 0,
  totalPages: 0,
  linkedAccounts: 0
})

onMounted(async () => {
  // Fetch user stats
  await loadStats()
})

const loadStats = async () => {
  // Get linked accounts count
  stats.value.linkedAccounts = authStore.linkedAccounts.length

  // TODO: Fetch books and pages stats from API
  // For now, use placeholder values
  stats.value.totalBooks = 0
  stats.value.totalPages = 0
}

const handleLogout = async () => {
  isLoggingOut.value = true
  try {
    await authStore.logout()
    toastStore.success('已退出登录')
    router.push('/')
  } catch {
    toastStore.error('退出失败，请重试')
  } finally {
    isLoggingOut.value = false
  }
}
</script>

<style scoped>
.page-profile {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--space-6);
}

/* Back Button */
.back-button {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  margin-bottom: var(--space-4);
}

.back-button:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

/* Profile Header */
.profile-header {
  text-align: center;
  padding: var(--space-8) 0;
}

.avatar-large {
  width: 96px;
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto var(--space-4);
  background: var(--color-accent);
  color: white;
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  border-radius: var(--radius-full);
  box-shadow: var(--shadow-md);
}

.profile-name {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--space-2);
}

.profile-email {
  font-size: var(--text-base);
  color: var(--text-secondary);
  margin: 0;
}

/* Profile Stats */
.profile-stats {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: var(--space-6);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
  margin-bottom: var(--space-6);
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: var(--text-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-accent);
}

.stat-label {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-top: var(--space-1);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-subtle);
}

/* Action Section */
.action-section,
.account-section,
.logout-section {
  margin-bottom: var(--space-6);
}

.section-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--space-4);
  padding-left: var(--space-2);
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-4);
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-5);
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
}

.action-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.action-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
}

.action-label {
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

/* Action List */
.action-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.action-item {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4);
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.action-item:hover {
  border-color: var(--border-medium);
  box-shadow: var(--shadow-sm);
}

.action-item-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.action-item-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
}

.action-item-title {
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.action-item-desc {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.action-item-arrow {
  color: var(--text-tertiary);
  flex-shrink: 0;
}

/* Logout Section */
.logout-button {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-4);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-error);
  background: var(--bg-elevated);
  border: 1px solid var(--color-error);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.logout-button:hover:not(:disabled) {
  background: var(--color-error);
  color: white;
}

.logout-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.logout-icon,
.spinner {
  flex-shrink: 0;
}

.spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .page-profile {
    padding: var(--space-4);
  }

  .avatar-large {
    width: 80px;
    height: 80px;
    font-size: var(--text-3xl);
  }

  .stat-value {
    font-size: var(--text-xl);
  }

  .stat-divider {
    height: 32px;
  }
}
</style>
