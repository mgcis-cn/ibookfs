import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import {
  Plus, PlayCircle, Camera, Sparkles, FileCode,
  Library, Cloud, Shield, FileText, PlusCircle,
  ChevronRight, Download, ArrowRight, User, BookOpen,
} from 'lucide-react'
import { useAuthStore } from '@/features/auth/store'
import './HomePage.css'

export default function HomePage() {
  const { isAuthenticated, initialize } = useAuthStore()

  useEffect(() => {
    if (!isAuthenticated()) initialize()
  }, [])

  return (
    <div className="landing-page">
      {/* Navigation */}
      <nav className="landing-nav">
        <div className="nav-container">
          <Link to="/" className="nav-logo">
            <img src="/logo.svg" alt="iBookFS" style={{ height: 32 }} />
          </Link>

          <div className="nav-links">
            <a href="#features" className="nav-link">功能特性</a>
            <a href="#how-it-works" className="nav-link">使用方法</a>
            <a href="#pricing" className="nav-link">价格</a>
          </div>

          <div className="nav-actions">
            {!isAuthenticated() ? (
              <>
                <Link to="/login" className="button-text">登录</Link>
                <Link to="/register" className="button-primary">免费注册</Link>
              </>
            ) : (
              <>
                <Link to="/settings" className="button-text"><User size={18} /></Link>
                <Link to="/library" className="button-primary">进入书架</Link>
              </>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="hero-section">
        <div className="hero-container">
          <div className="hero-content">
            <h1 className="hero-title">
              将纸质书籍
              <span className="text-accent"> 优雅地 </span>
              转化为数字资产
            </h1>
            <p className="hero-description">
              通过照片扫描、智能 OCR 识别和结构化整理，让您的纸质藏书成为可搜索、可分享的数字图书馆。
            </p>
            <div className="hero-actions">
              {!isAuthenticated() ? (
                <Link to="/register" className="button-primary button-lg">
                  <Plus size={20} />
                  开始免费使用
                </Link>
              ) : (
                <Link to="/library" className="button-primary button-lg">
                  <Library size={20} />
                  进入书架
                </Link>
              )}
              <a href="#how-it-works" className="button-secondary button-lg">
                <PlayCircle size={20} />
                了解工作原理
              </a>
            </div>
            <div className="hero-stats">
              <div className="stat-item">
                <div className="stat-value">10,000+</div>
                <div className="stat-label">已数字化书籍</div>
              </div>
              <div className="stat-divider" />
              <div className="stat-item">
                <div className="stat-value">99.8%</div>
                <div className="stat-label">OCR 准确率</div>
              </div>
              <div className="stat-divider" />
              <div className="stat-item">
                <div className="stat-value">5,000+</div>
                <div className="stat-label">活跃用户</div>
              </div>
            </div>
          </div>

          <div className="hero-visual">
            <div className="book-preview">
              <div className="preview-book">
                <div className="book-cover">
                  <BookOpen size={48} />
                </div>
                <div className="book-info">
                  <div className="book-title">设计心理学</div>
                  <div className="book-author">唐纳德·诺曼</div>
                </div>
              </div>
              <div className="preview-pages">
                {[1, 2, 3].map(i => (
                  <div key={i} className="page-card">
                    <FileText size={24} />
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="features-section">
        <div className="section-container">
          <div className="section-header">
            <h2 className="section-title">强大的功能，简单的操作</h2>
            <p className="section-subtitle">
              专为书籍数字化设计，从拍照到导出，全程流畅体验
            </p>
          </div>

          <div className="features-grid">
            {[
              { Icon: Camera, title: '拍照扫描', desc: '支持批量拍照，智能排序，自动裁剪和优化，让扫描变得简单高效' },
              { Icon: Sparkles, title: '智能 OCR', desc: '先进的文字识别技术，支持中英文混合识别，准确率高达 99.8%' },
              { Icon: FileCode, title: 'Markdown 导出', desc: '一键导出为 Markdown 格式，保留原文结构，方便编辑和分享' },
              { Icon: Library, title: '书架管理', desc: '直观的书架视图，支持搜索、筛选、分类，让藏书井井有条' },
              { Icon: Cloud, title: '云端同步', desc: '数据自动备份云端，多设备同步访问，随时随地阅读您的藏书' },
              { Icon: Shield, title: '隐私保护', desc: '端到端加密存储，您的藏书内容只有您能访问，安全无忧' },
            ].map(({ Icon, title, desc }) => (
              <div key={title} className="feature-card">
                <div className="feature-icon-wrapper">
                  <Icon size={32} />
                </div>
                <h3 className="feature-title">{title}</h3>
                <p className="feature-description">{desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section id="how-it-works" className="how-it-works-section">
        <div className="section-container">
          <div className="section-header">
            <h2 className="section-title">三步完成书籍数字化</h2>
            <p className="section-subtitle">
              简单直观的流程，无需专业技能，轻松上手
            </p>
          </div>

          <div className="steps-container">
            <div className="step-item">
              <div className="step-number">1</div>
              <div className="step-content">
                <h3 className="step-title">创建书籍</h3>
                <p className="step-description">
                  输入书名、作者等基本信息，创建一个新的数字书籍条目
                </p>
              </div>
              <div className="step-visual">
                <PlusCircle size={64} />
              </div>
            </div>

            <div className="step-arrow">
              <ChevronRight size={32} />
            </div>

            <div className="step-item">
              <div className="step-number">2</div>
              <div className="step-content">
                <h3 className="step-title">拍照上传</h3>
                <p className="step-description">
                  使用手机或相机拍摄书页，批量上传，系统自动排序整理
                </p>
              </div>
              <div className="step-visual">
                <Camera size={64} />
              </div>
            </div>

            <div className="step-arrow">
              <ChevronRight size={32} />
            </div>

            <div className="step-item">
              <div className="step-number">3</div>
              <div className="step-content">
                <h3 className="step-title">导出使用</h3>
                <p className="step-description">
                  OCR 识别文字，导出为 Markdown 或 PDF，完成数字化
                </p>
              </div>
              <div className="step-visual">
                <Download size={64} />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="cta-section">
        <div className="cta-container">
          <div className="cta-content">
            <h2 className="cta-title">准备好开始数字化之旅了吗？</h2>
            <p className="cta-description">
              立即注册，获得 14 天免费试用，体验完整的书籍数字化功能
            </p>
            <div className="cta-actions">
              {!isAuthenticated() ? (
                <Link to="/register" className="button-primary button-lg">
                  免费开始使用
                  <ArrowRight size={20} />
                </Link>
              ) : (
                <Link to="/library" className="button-primary button-lg">
                  进入我的书架
                  <ArrowRight size={20} />
                </Link>
              )}
            </div>
            <p className="cta-note">无需信用卡 · 随时取消 · 数据安全保证</p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="landing-footer">
        <div className="footer-container">
          <div className="footer-brand">
            <Link to="/" className="footer-logo">
              <img src="/logo.svg" alt="iBookFS" style={{ height: 32 }} />
            </Link>
            <p className="footer-description">
              让纸质书籍焕发新生，将知识转化为数字资产
            </p>
          </div>

          <div className="footer-links">
            <div className="footer-column">
              <h4 className="footer-column-title">产品</h4>
              <a href="#features" className="footer-link">功能特性</a>
              <a href="#pricing" className="footer-link">价格</a>
              <a href="#" className="footer-link">更新日志</a>
            </div>
            <div className="footer-column">
              <h4 className="footer-column-title">资源</h4>
              <a href="#" className="footer-link">帮助文档</a>
              <a href="#" className="footer-link">API 文档</a>
              <a href="#" className="footer-link">社区</a>
            </div>
            <div className="footer-column">
              <h4 className="footer-column-title">公司</h4>
              <a href="#" className="footer-link">关于我们</a>
              <a href="#" className="footer-link">联系方式</a>
              <a href="#" className="footer-link">隐私政策</a>
            </div>
          </div>
        </div>

        <div className="footer-bottom">
          <p>© 2026 iBookFS. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}
