<template>
  <div class="landing-page">
    <!-- Navigation -->
    <nav class="landing-nav">
      <div class="nav-container">
        <router-link to="/" class="nav-logo">
          <Logo />
        </router-link>

        <div class="nav-links">
          <a href="#features" class="nav-link">功能特性</a>
          <a href="#how-it-works" class="nav-link">使用方法</a>
          <a href="#pricing" class="nav-link">价格</a>
        </div>

        <div class="nav-actions">
          <router-link v-if="!authStore.isAuthenticated" to="/login" class="button-text">登录</router-link>
          <router-link v-else to="/settings" class="button-text">
            <User :size="18" />
          </router-link>
          <router-link v-if="!authStore.isAuthenticated" to="/register" class="button-primary">免费注册</router-link>
          <router-link v-else to="/library" class="button-primary">进入书架</router-link>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-container">
        <div class="hero-content">
          <h1 class="hero-title">
            将纸质书籍
            <span class="text-accent">优雅地</span>
            转化为数字资产
          </h1>
          <p class="hero-description">
            通过照片扫描、智能 OCR 识别和结构化整理，让您的纸质藏书成为可搜索、可分享的数字图书馆。
          </p>
          <div class="hero-actions">
            <router-link v-if="!authStore.isAuthenticated" to="/register" class="button-primary button-lg">
              <Plus :size="20" />
              开始免费使用
            </router-link>
            <router-link v-else to="/library" class="button-primary button-lg">
              <Library :size="20" />
              进入书架
            </router-link>
            <a href="#how-it-works" class="button-secondary button-lg">
              <PlayCircle :size="20" />
              了解工作原理
            </a>
          </div>
          <div class="hero-stats">
            <div class="stat-item">
              <div class="stat-value">10,000+</div>
              <div class="stat-label">已数字化书籍</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <div class="stat-value">99.8%</div>
              <div class="stat-label">OCR 准确率</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <div class="stat-value">5,000+</div>
              <div class="stat-label">活跃用户</div>
            </div>
          </div>
        </div>

        <div class="hero-visual">
          <div class="book-preview">
            <div class="preview-book">
              <div class="book-cover">
                <BookOpen :size="48" />
              </div>
              <div class="book-info">
                <div class="book-title">设计心理学</div>
                <div class="book-author">唐纳德·诺曼</div>
              </div>
            </div>
            <div class="preview-pages">
              <div v-for="i in 3" :key="i" class="page-card">
                <FileText :size="24" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features-section">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">强大的功能，简单的操作</h2>
          <p class="section-subtitle">
            专为书籍数字化设计，从拍照到导出，全程流畅体验
          </p>
        </div>

        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <Camera :size="32" />
            </div>
            <h3 class="feature-title">拍照扫描</h3>
            <p class="feature-description">
              支持批量拍照，智能排序，自动裁剪和优化，让扫描变得简单高效
            </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <Sparkles :size="32" />
            </div>
            <h3 class="feature-title">智能 OCR</h3>
            <p class="feature-description">
              先进的文字识别技术，支持中英文混合识别，准确率高达 99.8%
            </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <FileCode :size="32" />
            </div>
            <h3 class="feature-title">Markdown 导出</h3>
            <p class="feature-description">
              一键导出为 Markdown 格式，保留原文结构，方便编辑和分享
            </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <Library :size="32" />
            </div>
            <h3 class="feature-title">书架管理</h3>
            <p class="feature-description">
              直观的书架视图，支持搜索、筛选、分类，让藏书井井有条
            </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <Cloud :size="32" />
            </div>
            <h3 class="feature-title">云端同步</h3>
            <p class="feature-description">
              数据自动备份云端，多设备同步访问，随时随地阅读您的藏书
            </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <Shield :size="32" />
            </div>
            <h3 class="feature-title">隐私保护</h3>
            <p class="feature-description">
              端到端加密存储，您的藏书内容只有您能访问，安全无忧
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section id="how-it-works" class="how-it-works-section">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">三步完成书籍数字化</h2>
          <p class="section-subtitle">
            简单直观的流程，无需专业技能，轻松上手
          </p>
        </div>

        <div class="steps-container">
          <div class="step-item">
            <div class="step-number">1</div>
            <div class="step-content">
              <h3 class="step-title">创建书籍</h3>
              <p class="step-description">
                输入书名、作者等基本信息，创建一个新的数字书籍条目
              </p>
            </div>
            <div class="step-visual">
              <PlusCircle :size="64" />
            </div>
          </div>

          <div class="step-arrow">
            <ChevronRight :size="32" />
          </div>

          <div class="step-item">
            <div class="step-number">2</div>
            <div class="step-content">
              <h3 class="step-title">拍照上传</h3>
              <p class="step-description">
                使用手机或相机拍摄书页，批量上传，系统自动排序整理
              </p>
            </div>
            <div class="step-visual">
              <Camera :size="64" />
            </div>
          </div>

          <div class="step-arrow">
            <ChevronRight :size="32" />
          </div>

          <div class="step-item">
            <div class="step-number">3</div>
            <div class="step-content">
              <h3 class="step-title">导出使用</h3>
              <p class="step-description">
                OCR 识别文字，导出为 Markdown 或 PDF，完成数字化
              </p>
            </div>
            <div class="step-visual">
              <Download :size="64" />
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="cta-container">
        <div class="cta-content">
          <h2 class="cta-title">准备好开始数字化之旅了吗？</h2>
          <p class="cta-description">
            立即注册，获得 14 天免费试用，体验完整的书籍数字化功能
          </p>
          <div class="cta-actions">
            <router-link v-if="!authStore.isAuthenticated" to="/register" class="button-primary button-lg">
              免费开始使用
              <ArrowRight :size="20" />
            </router-link>
            <router-link v-else to="/library" class="button-primary button-lg">
              进入我的书架
              <ArrowRight :size="20" />
            </router-link>
          </div>
          <p class="cta-note">无需信用卡 · 随时取消 · 数据安全保证</p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="footer-container">
        <div class="footer-brand">
          <router-link to="/" class="footer-logo">
            <Logo />
          </router-link>
          <p class="footer-description">
            让纸质书籍焕发新生，将知识转化为数字资产
          </p>
        </div>

        <div class="footer-links">
          <div class="footer-column">
            <h4 class="footer-column-title">产品</h4>
            <a href="#features" class="footer-link">功能特性</a>
            <a href="#pricing" class="footer-link">价格</a>
            <a href="#" class="footer-link">更新日志</a>
          </div>

          <div class="footer-column">
            <h4 class="footer-column-title">资源</h4>
            <a href="#" class="footer-link">帮助文档</a>
            <a href="#" class="footer-link">API 文档</a>
            <a href="#" class="footer-link">社区</a>
          </div>

          <div class="footer-column">
            <h4 class="footer-column-title">公司</h4>
            <a href="#" class="footer-link">关于我们</a>
            <a href="#" class="footer-link">联系方式</a>
            <a href="#" class="footer-link">隐私政策</a>
          </div>
        </div>
      </div>

      <div class="footer-bottom">
        <p>&copy; 2026 iBookFS. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import {
  Plus,
  PlayCircle,
  Camera,
  Sparkles,
  FileCode,
  Library,
  Cloud,
  Shield,
  FileText,
  PlusCircle,
  ChevronRight,
  Download,
  ArrowRight,
  User,
} from 'lucide-vue-next'
import Logo from '@/components/ui/Logo.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Initialize auth store on page load to check if user is logged in
onMounted(async () => {
  if (!authStore.isAuthenticated) {
    await authStore.initialize()
  }
})
</script>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* Navigation */
.landing-nav {
  position: sticky;
  top: 0;
  z-index: var(--z-fixed);
  background: rgba(250, 248, 243, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-subtle);
}

.nav-container {
  max-width: var(--container-xl);
  margin: 0 auto;
  padding: 0 var(--space-6);
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  opacity: 1;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.nav-logo:hover {
  opacity: 0.85;
}

.nav-links {
  display: flex;
  gap: var(--space-6);
}

.nav-link {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: #2B2B2B;
  text-decoration: none;
  transition: color var(--duration-fast) var(--ease-out);
}

.nav-link:hover {
  color: var(--color-accent);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

/* Hero Section */
.hero-section {
  padding: var(--space-16) var(--space-6);
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
}

.hero-container {
  max-width: var(--container-xl);
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-12);
  align-items: center;
}

.hero-content {
  animation: fadeInUp 0.8s var(--ease-out);
}

.hero-title {
  font-size: var(--text-5xl);
  font-weight: var(--font-weight-bold);
  line-height: 1.2;
  margin-bottom: var(--space-6);
  color: var(--color-ink);
}

.text-accent {
  color: var(--color-accent);
}

.hero-description {
  font-size: var(--text-lg);
  line-height: 1.6;
  color: var(--color-ink-light);
  margin-bottom: var(--space-8);
  max-width: 540px;
}

.hero-actions {
  display: flex;
  gap: var(--space-4);
  margin-bottom: var(--space-10);
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: var(--space-6);
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
  color: var(--color-ink-light);
  margin-top: var(--space-1);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-subtle);
}

.hero-visual {
  animation: fadeIn 1s var(--ease-out);
}

.book-preview {
  position: relative;
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
  box-shadow: var(--shadow-lg);
}

.preview-book {
  display: flex;
  gap: var(--space-4);
  padding: var(--space-4);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-4);
}

.book-cover {
  width: 80px;
  height: 120px;
  background: var(--color-accent);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.book-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.book-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
  margin-bottom: var(--space-1);
}

.book-author {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
}

.preview-pages {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
}

.page-card {
  aspect-ratio: 3/4;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-ink-light);
}

/* Features Section */
.features-section {
  padding: var(--space-16) var(--space-6);
}

.section-container {
  max-width: var(--container-xl);
  margin: 0 auto;
}

.section-header {
  text-align: center;
  margin-bottom: var(--space-12);
}

.section-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-4);
  color: var(--color-ink);
}

.section-subtitle {
  font-size: var(--text-lg);
  color: var(--color-ink-light);
  max-width: 600px;
  margin: 0 auto;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-8);
}

.feature-card {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-8);
  transition: all var(--duration-base) var(--ease-out);
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.feature-icon-wrapper {
  width: 64px;
  height: 64px;
  background: var(--color-accent-light);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
  margin-bottom: var(--space-6);
}

.feature-title {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
  color: var(--color-ink);
}

.feature-description {
  font-size: var(--text-base);
  line-height: 1.6;
  color: var(--color-ink-light);
}

/* How It Works Section */
.how-it-works-section {
  padding: var(--space-16) var(--space-6);
  background: var(--bg-secondary);
}

.steps-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
}

.step-item {
  flex: 1;
  max-width: 320px;
  text-align: center;
}

.step-number {
  width: 48px;
  height: 48px;
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xl);
  font-weight: var(--font-weight-bold);
  margin: 0 auto var(--space-6);
}

.step-title {
  font-size: var(--text-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
  color: var(--color-ink);
}

.step-description {
  font-size: var(--text-sm);
  line-height: 1.6;
  color: var(--color-ink-light);
  margin-bottom: var(--space-6);
}

.step-visual {
  color: var(--color-accent);
  opacity: 0.5;
}

.step-arrow {
  color: var(--color-ink-light);
}

/* CTA Section */
.cta-section {
  padding: var(--space-16) var(--space-6);
}

.cta-container {
  max-width: var(--container-lg);
  margin: 0 auto;
}

.cta-content {
  background: linear-gradient(135deg, var(--color-accent) 0%, #a6624a 100%);
  border-radius: var(--radius-xl);
  padding: var(--space-12);
  text-align: center;
  color: white;
}

.cta-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-4);
}

.cta-description {
  font-size: var(--text-lg);
  opacity: 0.9;
  margin-bottom: var(--space-8);
  max-width: 540px;
  margin-left: auto;
  margin-right: auto;
}

.cta-actions {
  display: flex;
  justify-content: center;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.cta-note {
  font-size: var(--text-sm);
  opacity: 0.8;
}

/* Footer */
.landing-footer {
  background: var(--bg-elevated);
  border-top: 1px solid var(--border-subtle);
  padding: var(--space-12) var(--space-6) var(--space-6);
}

.footer-container {
  max-width: var(--container-xl);
  margin: 0 auto;
  display: grid;
  grid-template-columns: 2fr 3fr;
  gap: var(--space-12);
  margin-bottom: var(--space-8);
}

.footer-brand {
  max-width: 320px;
}

.footer-logo {
  display: inline-flex;
  align-items: center;
  text-decoration: none;
  margin-bottom: var(--space-4);
  opacity: 0.9;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.footer-logo:hover {
  opacity: 1;
}

.footer-description {
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  line-height: 1.6;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-8);
}

.footer-column-title {
  font-size: var(--text-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-ink);
  margin-bottom: var(--space-4);
}

.footer-link {
  display: block;
  font-size: var(--text-sm);
  color: var(--color-ink-light);
  text-decoration: none;
  margin-bottom: var(--space-3);
  transition: color var(--duration-fast) var(--ease-out);
}

.footer-link:hover {
  color: var(--color-accent);
}

.footer-bottom {
  max-width: var(--container-xl);
  margin: 0 auto;
  padding-top: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  text-align: center;
  font-size: var(--text-sm);
  color: var(--color-ink-light);
}

/* Buttons */
.button-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: white;
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
}

.button-primary:hover {
  background: #a6624a;
  transform: translateY(-1px);
}

.button-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-6);
  font-size: var(--text-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-ink);
  background: var(--bg-elevated);
  border: 1.5px solid var(--border-medium);
  border-radius: var(--radius-md);
  cursor: pointer;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
}

.button-secondary:hover {
  background: var(--bg-secondary);
  border-color: var(--color-accent);
}

.button-text {
  display: inline-flex;
  align-items: center;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-weight-medium);
  color: #2B2B2B;
  text-decoration: none;
  transition: color var(--duration-fast) var(--ease-out);
}

.button-text:hover {
  color: var(--color-accent);
}

.button-lg {
  padding: var(--space-4) var(--space-8);
  font-size: var(--text-lg);
}

/* Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Responsive */
@media (max-width: 1024px) {
  .hero-container {
    grid-template-columns: 1fr;
    text-align: center;
  }

  .hero-description {
    margin-left: auto;
    margin-right: auto;
  }

  .hero-actions {
    justify-content: center;
  }

  .hero-stats {
    justify-content: center;
  }

  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .steps-container {
    flex-direction: column;
  }

  .step-arrow {
    transform: rotate(90deg);
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .hero-title {
    font-size: var(--text-3xl);
  }

  .hero-actions {
    flex-direction: column;
  }

  .hero-stats {
    flex-direction: column;
    gap: var(--space-4);
  }

  .stat-divider {
    display: none;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .footer-container {
    grid-template-columns: 1fr;
  }

  .footer-links {
    grid-template-columns: 1fr;
  }
}
</style>
