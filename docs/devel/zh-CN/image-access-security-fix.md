# 图片访问安全修复方案

> **背景**：当前系统存在严重安全漏洞，图片存储路径直接通过 `r.Static()` 暴露，绕过了所有访问控制。
> **目标**：参考云厂商 OSS 的鉴权设计，实现安全的图片访问机制。

---

## 一、问题现状

### 1.1 数据库设计（正确）

```go
// internal/model/image.go
type Image struct {
    ID          uint   `json:"id"`
    OwnerID     uint   `json:"owner_id"`
    IsPublic    bool   `json:"is_public" gorm:"default:false"`      // 公开标识
    AccessToken string `json:"access_token" gorm:"size:64"`         // 访问令牌
    // ...
}
```

### 1.2 API 实现（正确）

`internal/handler/image.go` 的 `ServeImage` 已实现双认证：

```go
// 支持两种访问方式：
// 1. JWT 认证: Authorization: Bearer {jwt_token}
// 2. Token 认证: ?token={access_token}
func ServeImage(c *gin.Context) {
    token := c.Query("token")

    if token != "" {
        // 通过 access_token 访问（共享场景）
        image, err = imageService.GetByAccessToken(ctx, token)
    } else {
        // 通过 JWT 访问（登录用户）
        userID, exists := middleware.GetUserID(c)
        image, err = imageService.GetByID(ctx, uint(id), userID)
    }
}
```

### 1.3 路由配置（安全漏洞）

`internal/router/router.go:125`

```go
// ❌ 问题代码：直接暴露整个存储目录
r.Static(cfg.Storage.Local.BaseURL, cfg.Storage.Local.BasePath)
// 实际: r.Static("/static", "./storages")

// 结果：任何人都可以直接访问
// http://localhost:8080/static/images/book/1/original/xxx.jpg
// 完全绕过了 is_public 和 access_token 的控制
```

---

## 二、云厂商 OSS 鉴权设计参考

### 2.1 OSS 的两种核心共享方式

| 模式 | URL 形式 | 鉴权方式 | 适用场景 |
|------|----------|----------|----------|
| **公共读** | `https://bucket.oss.com/path/file.jpg` | 资源权限替代鉴权 | 完全公开的资源 |
| **私有读 + 签名URL** | `https://bucket.oss.com/path.jpg?OSSAccessKeyId=xxx&Expires=xxx&Signature=xxx` | URL 参数伪装鉴权 | 临时分享、隐私保护 |

### 2.2 OSS 签名 URL 参数解析

```
https://bucket.oss-cn-hangzhou.aliyuncs.com/images/book/1/original/xxx.jpg
?OSSAccessKeyId=TMP.3xxxxx        # 临时 AK
&Expires=1738360800                # Unix 时间戳过期时间
&Signature=vMPkxxxxx               # HMAC-SHA1 签名
```

**核心原理**：
1. 服务端用 `SecretKey` + `请求参数` 计算签名
2. 客户端 URL 中的 `Signature` 必须与服务端计算结果一致
3. `Expires` 控制签名有效期

### 2.3 OSS 签名算法（简化版）

```
StringToSign = HTTP-Verb + "\n" +
               Content-MD5 + "\n" +
               Content-Type + "\n" +
               Expires + "\n" +
               CanonicalizedResource

Signature = Base64( HMAC-SHA1( SecretKey, StringToSign ) )
```

---

## 三、iBookFS 修复方案

### 3.1 设计决策

采用**简化的双轨制**，参考 OSS 但不过度复杂：

```
┌─────────────────────────────────────────────────────────────┐
│  图片访问策略                                                │
├─────────────────────────────────────────────────────────────┤
│  公开图片 (is_public=true)                                   │
│    → /api/v1/images/:id/file?variant=small                  │
│    → 无需任何认证                                            │
│                                                             │
│  私有图片 (is_public=false)                                  │
│    → /api/v1/images/:id/file?variant=small&token=xxx       │
│    → 需要 access_token                                      │
│    → 或 JWT 认证 (登录用户)                                  │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 修复步骤

#### Step 1: 移除静态文件直接暴露

**文件**: `internal/router/router.go`

```go
// ❌ 删除这行
r.Static(cfg.Storage.Local.BaseURL, cfg.Storage.Local.BasePath)
```

#### Step 2: 更新 ServeImage 支持 is_public

**文件**: `internal/handler/image.go`

```go
func ServeImage(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    variant := c.DefaultQuery("variant", "original")
    token := c.Query("token")

    var image *model.Image
    var err error

    if token != "" {
        // 方式1: 通过 access_token 访问
        image, err = imageService.GetByAccessToken(ctx, token)
        if err != nil || image.ID != uint(id) {
            c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
            return
        }
    } else {
        // 方式2: 先尝试获取图片，检查是否公开
        image, err = imageService.GetByIDForAccess(ctx, uint(id))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
            return
        }

        // 如果是私有图片，需要 JWT 认证
        if !image.IsPublic {
            userID, exists := middleware.GetUserID(c)
            if !exists || image.OwnerID != userID {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
                return
            }
        }
        // 如果是公开图片，直接允许访问
    }

    // 下载并返回文件
    reader, mimeType, _ := imageService.DownloadFile(ctx, image, variant)
    c.Header("Content-Type", mimeType)
    c.Header("Cache-Control", "public, max-age=31536000")
    io.Copy(c.Writer, reader)
}
```

#### Step 3: Service 层添加新方法

**文件**: `internal/service/image_service.go`

```go
// GetByIDForAccess 根据ID获取图片（用于访问控制检查，不验证所有权）
func (s *ImageService) GetByIDForAccess(ctx context.Context, id uint) (*model.Image, error) {
    var image model.Image
    err := s.db.WithContext(ctx).Where("id = ?", id).First(&image).Error
    return &image, err
}
```

### 3.3 访问 URL 对比

| 场景 | 修复前 | 修复后 |
|------|--------|--------|
| 用户访问自己的私有图 | `http://localhost:8080/static/xxx.jpg` ❌ | `Authorization: Bearer {jwt}` → `/api/v1/images/:id/file` |
| 分享私有图片给他人 | 无法安全分享 ❌ | `/api/v1/images/:id/file?token={access_token}` |
| 公开图片访问 | 同上 ❌ | `/api/v1/images/:id/file` (无需认证) |

---

## 四、后续扩展（可选）

### 4.1 带过期时间的签名 URL

参考 OSS 的 `Expires` 参数，可扩展 `access_token` 机制：

```go
type Image struct {
    AccessToken    string    `json:"access_token" gorm:"size:64"`
    AccessExpires  *time.Time `json:"access_expires,omitempty"` // 新增
}
```

生成签名 URL：
```go
// /api/v1/images/:id/file?signature=xxx&expires=1738360800
```

### 4.2 响应头控制

```go
// 下载而非预览
c.Header("Content-Disposition", "attachment; filename=xxx.jpg")

// 跨域共享
c.Header("Access-Control-Allow-Origin", "*")
```

---

## 五、检查清单

- [ ] 移除 `router.go` 中的 `r.Static()` 调用
- [ ] 更新 `ServeImage` handler 支持 `is_public` 检查
- [ ] Service 层添加 `GetByIDForAccess` 方法
- [ ] 测试：JWT 用户访问自己的私有图片
- [ ] 测试：使用 access_token 共享私有图片
- [ ] 测试：无需认证访问公开图片
- [ ] 测试：确保直接访问 `/static/*` 返回 404

---

## 六、相关文件

| 文件 | 修改类型 |
|------|----------|
| `internal/router/router.go` | 删除静态路由 |
| `internal/handler/image.go` | 更新 ServeImage 逻辑 |
| `internal/service/image_service.go` | 添加 GetByIDForAccess |
| `internal/model/image.go` | 无需修改（字段已存在） |
