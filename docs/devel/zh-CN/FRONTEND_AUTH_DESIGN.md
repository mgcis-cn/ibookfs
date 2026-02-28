# iBookFS 多Provider认证系统设计文档

## 概述

本文档描述了 iBookFS 项目的多Provider认证系统设计，支持邮箱验证码、GitHub、Gitee、iCloud 等多种登录方式，允许一个用户绑定多种登录方式。

---

## 设计理念

### 美学方向
- **风格**: 现代、简约、温暖
- **色彩**: 纸质质感色调 (#FAF8F3) + 暖色调强调色 (#C17B5C)
- **排版**: Crimson Pro (衬线) 用于标题，Inter 用于正文
- **视觉层次**: 清晰的信息架构，柔和的阴影和圆角

### 设计原则
1. **一致性**: 所有组件遵循统一的设计语言
2. **可访问性**: 支持键盘导航、屏幕阅读器
3. **响应式**: 适配桌面、平板、手机
4. **渐进增强**: 从基础功能开始，逐步添加高级特性

---

## 文件结构

```
ibookfs/
├── internal/
│   └── migrations/
│       └── 001_auth_schema.sql          # 数据库Schema
├── docs/
│   ├── AUTH_API_DESIGN.md                # API设计文档
│   └── FRONTEND_AUTH_DESIGN.md           # 本文档
└── web/
    └── src/
        ├── types/
        │   └── index.ts                  # TypeScript类型定义
        ├── stores/
        │   └── auth.ts                   # Pinia认证Store
        ├── components/
        │   └── ui/
        │       └── OAuthButton.vue       # OAuth按钮组件
        └── pages/
            ├── AuthPage.vue              # 统一认证页面
            └── AccountBindingPage.vue    # 账号绑定管理页面
```

---

## 数据库设计

### 核心表结构

| 表名 | 说明 |
|------|------|
| `users` | 用户核心信息 |
| `oauth_identities` | OAuth身份绑定 |
| `email_verification_codes` | 邮箱验证码 |
| `sessions` | 会话管理 |
| `login_history` | 登录历史审计 |

### 关键关系

```
users (1) ----< (*) oauth_identities
users (1) ----< (*) sessions
users (1) ----< (*) login_history
```

**特点**:
- 一个用户可以绑定多个OAuth身份
- OAuth身份包含token信息（加密存储）
- 支持设置主登录方式

---

## API设计

### 认证端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/auth/send-code` | POST | 发送邮箱验证码 |
| `/api/v1/auth/login` | POST | 邮箱验证码登录 |
| `/api/v1/auth/register` | POST | 邮箱验证码注册 |
| `/api/v1/auth/oauth/authorize` | GET | 获取OAuth授权URL |
| `/api/v1/auth/oauth/callback` | POST | OAuth回调处理 |
| `/api/v1/auth/oauth/link` | POST | 绑定OAuth账号 |
| `/api/v1/auth/oauth/unlink` | DELETE | 解绑OAuth账号 |
| `/api/v1/auth/bind-email` | POST | 绑定邮箱 |
| `/api/v1/auth/accounts` | GET | 获取已绑定账号列表 |
| `/api/v1/auth/me` | GET | 获取当前用户信息 |
| `/api/v1/auth/logout` | POST | 登出 |

详细API文档见 `docs/AUTH_API_DESIGN.md`

---

## 前端组件设计

### 1. OAuthButton 组件

**文件**: `web/src/components/ui/OAuthButton.vue`

可复用的OAuth登录/绑定按钮组件。

**Props**:
```typescript
interface Props {
  provider: OAuthProvider      // github | gitee | icloud | google | wechat
  variant?: 'default' | 'outlined' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  isLoading?: boolean
  isLinked?: boolean           // 是否已绑定
  isPrimary?: boolean          // 是否为主账号
}
```

**特点**:
- 内置各Provider品牌图标
- 支持加载状态
- 支持已绑定/主账号标识
- 完整的hover/active/focus状态

### 2. AuthPage 组件

**文件**: `web/src/pages/AuthPage.vue`

统一的认证页面，整合了登录和注册功能。

**功能**:
- Tab切换: 邮箱登录 / OAuth登录
- 两步验证: 先邮箱 → 后验证码
- 支持GitHub/Gitee/iCloud/Google登录
- OAuth弹窗授权处理
- 响应式布局（左侧品牌区 + 右侧表单区）

**状态管理**:
```typescript
const state = {
  mode: 'login' | 'register'       // 登录/注册模式
  activeTab: 'email' | 'oauth'     // 当前Tab
  step: 1 | 2                       // 邮箱步骤: 1输入邮箱 2输入验证码
}
```

### 3. AccountBindingPage 组件

**文件**: `web/src/pages/AccountBindingPage.vue`

账号绑定管理页面。

**功能模块**:

1. **主账号区域**: 显示已绑定的所有登录方式
   - 邮箱账号（如果已绑定）
   - OAuth账号列表
   - 设置主账号
   - 解绑账号

2. **绑定新账号区域**:
   - 显示可用的OAuth Provider
   - 已绑定状态标识
   - 一键绑定

3. **绑定邮箱区域** (仅OAuth用户):
   - 发送验证码
   - 输入验证码完成绑定

---

## 状态管理 (Pinia Store)

### Auth Store

**文件**: `web/src/stores/auth.ts`

**State**:
```typescript
interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}
```

**Computed**:
```typescript
{
  userEmail: string
  userDisplayName: string
  userAvatar: string | null
  linkedAccounts: LinkedAccount[]
  primaryAccount: LinkedAccount | undefined
  hasEmail: boolean
  hasOAuth: boolean
  linkedProviders: OAuthProvider[]
}
```

**Actions**:
```typescript
{
  // 认证相关
  sendCode(data: SendCodeRequest): Promise<boolean>
  login(data: LoginRequest): Promise<boolean>
  register(data: RegisterRequest): Promise<boolean>

  // OAuth相关
  getOAuthAuthorizeUrl(provider, redirectUri): Promise<{url, state} | null>
  oauthCallback(provider, code, state): Promise<boolean>
  linkOAuth(data: LinkOAuthRequest): Promise<boolean>
  unlinkOAuth(data: UnlinkOAuthRequest): Promise<boolean>

  // 账号管理
  bindEmail(data: BindEmailRequest): Promise<boolean>
  setPrimaryAccount(data: SetPrimaryAccountRequest): Promise<boolean>
  getLinkedAccountsList(): Promise<LinkedAccount[]>

  // 用户相关
  refreshAuthToken(): Promise<boolean>
  refreshUser(): Promise<boolean>
  updateProfile(data: UpdateProfileRequest): Promise<boolean>
  getCurrentUser(): Promise<User | null>
  logout(): Promise<void>

  // 初始化
  initialize(): Promise<void>
}
```

---

## 类型定义

**文件**: `web/src/types/index.ts`

新增/更新的类型:

```typescript
// OAuth Provider类型
type OAuthProvider = 'github' | 'gitee' | 'icloud' | 'google' | 'wechat'

// 用户信息（增强版）
interface User {
  id: string
  email: string
  emailVerified: boolean
  firstName: string
  lastName: string
  displayName?: string
  avatarUrl?: string
  status: UserStatus
  role: UserRole
  preferences?: UserPreferences
  language: string
  timezone: string
  createdAt: string
  updatedAt: string
  linkedAccounts?: LinkedAccount[]  // 新增
}

// 绑定的账号信息
interface LinkedAccount {
  provider: 'email' | OAuthProvider
  providerEmail?: string
  providerUserId?: string
  providerUsername?: string
  isPrimary: boolean
  linkedAt?: string
  lastUsedAt?: string
}

// 认证响应（增强版）
interface AuthResponse {
  user: User
  token: string
  refreshToken?: string
  expiresIn: number
  isNewUser?: boolean
}
```

---

## 用户流程

### 注册流程

```
┌──────────────────────────────────────────────────────────────────┐
│                         注册流程                                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. 访问 /register                                                │
│     ↓                                                            │
│  2. 选择注册方式:                                                 │
│     ├─ 邮箱注册                                                   │
│     │  ├─ 输入姓名 + 邮箱                                        │
│     │  ├─ 发送验证码                                             │
│     │  └─ 输入验证码 → 创建账号                                   │
│     │                                                            │
│     └─ OAuth注册                                                  │
│        ├─ 选择Provider (GitHub/Gitee/iCloud/Google)               │
│        ├─ 弹窗授权                                                │
│        └─ 自动创建账号 → 登录                                      │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 登录流程

```
┌──────────────────────────────────────────────────────────────────┐
│                         登录流程                                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. 访问 /login                                                   │
│     ↓                                                            │
│  2. 选择登录方式:                                                 │
│     ├─ 邮箱登录                                                   │
│     │  ├─ 输入邮箱                                                │
│     │  ├─ 发送验证码                                              │
│     │  └─ 输入验证码 → 登录成功                                   │
│     │                                                            │
│     └─ OAuth登录                                                  │
│        ├─ 选择Provider                                            │
│        ├─ 弹窗授权                                                │
│        └─ 登录成功                                                │
│                                                                  │
│  3. 已有账号，OAuth邮箱不同的情况:                                │
│     ├─ 显示提示: "检测到已有账号"                                  │
│     └─ 用户选择: 绑定到已有账号 / 创建新账号                        │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 账号绑定流程

```
┌──────────────────────────────────────────────────────────────────┐
│                       账号绑定流程                                 │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. 访问 /settings/accounts                                       │
│     ↓                                                            │
│  2. 查看已绑定账号列表                                            │
│     ↓                                                            │
│  3. 绑定新账号:                                                   │
│     ├─ OAuth绑定                                                 │
│     │  └─ 点击Provider → 授权 → 完成                              │
│     │                                                            │
│     └─ 邮箱绑定 (仅OAuth用户)                                    │
│        ├─ 输入邮箱                                                │
│        ├─ 发送验证码                                              │
│        └─ 输入验证码 → 完成                                      │
│                                                                  │
│  4. 设置主账号:                                                   │
│     └─ 点击"设为主账号" → 确认                                    │
│                                                                  │
│  5. 解绑账号:                                                     │
│     ├─ 点击"解绑"                                                 │
│     ├─ 确认对话框                                                 │
│     └─ 完成解绑                                                  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 路由配置

```typescript
// router/index.ts
const routes = [
  // 替换原有的 /login 和 /route 为统一入口
  {
    path: '/login',
    component: () => import('@/pages/AuthPage.vue'),
    props: { mode: 'login' }
  },
  {
    path: '/register',
    component: () => import('@/pages/AuthPage.vue'),
    props: { mode: 'register' }
  },
  // OAuth回调
  {
    path: '/auth/callback',
    component: () => import('@/pages/OAuthCallbackPage.vue')
  },
  // 账号绑定管理
  {
    path: '/settings/accounts',
    component: () => import('@/pages/AccountBindingPage.vue'),
    meta: { requiresAuth: true }
  }
]
```

---

## 环境变量配置

```bash
# .env.production
VITE_API_BASE_URL=https://api.ibookfs.com/api
VITE_USE_MOCK=false

# OAuth配置（后端）
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
GITEE_CLIENT_ID=xxx
GITEE_CLIENT_SECRET=xxx
ICLOUD_CLIENT_ID=xxx
ICLOUD_CLIENT_SECRET=xxx
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx

# OAuth重定向URI
OAUTH_REDIRECT_URI=https://ibookfs.com/auth/callback
```

---

## 安全考虑

### 前端安全
1. **Token存储**: 使用HttpOnly、Secure、SameSite Cookie
2. **XSS防护**: 输入转义、CSP策略
3. **CSRF防护**: State参数验证
4. **Token刷新**: 自动刷新过期token

### 后端安全
1. **密码哈希**: bcrypt/scrypt
2. **Token加密**: OAuth access_token/refresh_token加密存储
3. **Rate Limiting**: 防止暴力破解
4. **会话管理**: 设备指纹、IP追踪
5. **审计日志**: 登录历史、安全事件

---

## 浏览器兼容性

| 浏览器 | 最低版本 |
|--------|----------|
| Chrome | 90+ |
| Safari | 14+ |
| Firefox | 88+ |
| Edge | 90+ |

移动端:
- iOS Safari 14+
- Chrome Android 90+

---

## 后续开发步骤

### Phase 1: 后端实现
1. 创建Go模型 (`internal/model/user.go`, `internal/model/oauth.go`)
2. 实现认证handler (`internal/handler/auth.go`)
3. OAuth集成 (`internal/service/oauth.go`)
4. JWT中间件 (`internal/middleware/auth.go`)

### Phase 2: 前端集成
1. 更新路由配置
2. 替换mock API为真实调用
3. 添加路由守卫
4. 实现token自动刷新

### Phase 3: 测试与优化
1. 单元测试
2. E2E测试
3. 性能优化
4. 安全审计

---

## 参考资料

- [RFC 6749: OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7519: JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
- [Sign in with Apple](https://developer.apple.com/sign-in-with-apple/)
- [GitHub OAuth Apps](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [Gitee OAuth文档](https://gitee.com/api/v5/oauth_doc)
