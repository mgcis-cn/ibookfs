# iBookFS 图片管理系统设计方案

## 文档元信息

- **类型**: 技术决策文档
- **版本**: v1.0
- **创建日期**: 2025-02-02
- **状态**: 草稿

---

## 1. 背景与目标

### 1.1 当前系统现状

iBookFS 是一个图书管理系统，主要用于管理图书的元数据（标题、作者、ISBN等）和上传进度追踪。目前系统中图片相关功能非常有限：

- **Book.Cover 字段**：仅是一个简单的字符串路径字段，存储封面图片的文件路径
- **缺乏图片元数据**：无法获取图片尺寸、格式、大小等信息
- **单一分辨率**：只支持原始图片，无法提供渐进式加载体验
- **与业务耦合**：图片路径直接嵌入业务表，缺乏灵活性

### 1.2 问题陈述

随着系统发展，当前图片处理方式暴露出以下问题：

1. **用户体验问题**
   - 大图片加载慢，影响页面性能
   - 无法提供BlurHash等渐进式加载效果
   - 移动端流量消耗大

2. **架构问题**
   - 图片与业务逻辑强耦合
   - 无法灵活切换存储方式（本地/云存储）
   - 缺乏统一的图片资源管理

3. **扩展性问题**
   - 难以支持多分辨率生成
   - 无法支持图片分组、标签等功能
   - 不利于未来功能扩展（如相册、水印等）

### 1.3 系统目标

构建一个**独立的图片资源管理子系统**，具备以下核心能力：

1. **渐进式加载**：支持 BlurHash + 多分辨率，提供流畅的加载体验
2. **存储抽象**：通过接口抽象存储层，支持本地存储和未来云存储迁移
3. **业务解耦**：图片作为独立资源，通过关联关系与业务对接
4. **可扩展性**：预留分组、标签等扩展接口，技术架构支持平滑演进

### 1.4 设计原则

- **技术栈统一**：使用 Go + GORM + Gin，与现有系统保持一致
- **可扩展性优先**：存储接口化，便于未来接入 OSS/S3 等云存储
- **简单实用**：暂不支持复杂图片处理（裁剪、滤镜等），聚焦核心功能
- **渐进式落地**：分阶段实施，优先核心功能，扩展功能预留接口

### 1.5 与现有系统集成

- **数据迁移**：现有 Book.Cover 字段需要迁移到新的图片系统
- **API 路由**：新增 `/api/v1/images` 路由组，与现有路由体系保持一致
- **认证复用**：使用现有的 JWT + AuthMiddleware 认证机制
- **数据库兼容**：支持 MySQL 和 SQLite，与现有配置保持一致

### 1.6 非目标（暂不实现）

- ❌ 复杂图片处理（裁剪、旋转、水印、滤镜等）
- ❌ 图片编辑功能
- ❌ 社交分享功能
- ❌ 图片版权管理

### 1.7 规模考量

- 单本书图片数量：100-200 张
- 单张图片大小：< 5MB
- 总体规模：初期以本地存储为主，暂不考虑 CDN 加速

### 1.8 交付物

本方案将为技术团队提供：

- 图片数据模型设计（GORM 结构体定义）
- 存储接口规范（便于扩展云存储）
- RESTful API 设计（上传、访问、删除）
- 技术选型建议（Go 图片处理库）
- 分阶段实施计划

---

## 2. 核心需求

### 2.1 功能需求

#### 2.1.1 图片上传

- **单文件上传**：使用 multipart/form-data 格式
- **上传进度**：可选的进度显示（WebSocket/SSE）
- **断点续传**：支持大文件上传中断后继续
- **文件校验**：MIME类型校验、文件扩展名校验
- **无大小限制**：暂不限制单文件大小（单张<5MB，但系统不强制限制）

#### 2.1.2 渐进式加载

- **BlurHash**：后端异步生成，存储于数据库
- **多分辨率**：自动生成三个变体
  - 原图（original）：保持原始尺寸
  - 中图（medium）：最大宽度 1200px
  - 小图（small）：最大宽度 400px
- **异步处理**：上传后立即返回，后台异步生成变体

#### 2.1.3 图片访问

- **私有化访问**：默认私有，通过签名URL访问
- **变体选择**：URL参数 `?variant=small/medium/original`
- **访问令牌**：基于密钥的签名机制，支持分享

#### 2.1.4 图片分组

- **基础分组**：按图书分组（book_id 关联）
- **扩展接口**：预留通用分组接口（未来支持相册、标签等）

#### 2.1.5 存储抽象

- **接口化设计**：定义 Storage 接口，支持多种存储实现
- **本地存储**：默认实现，文件系统存储
- **云存储扩展**：预留 OSS/S3 接口（阿里云OSS、腾讯云COS、AWS S3）
- **混合架构**：支持本地上云，新旧图片共存

### 2.2 非功能需求

#### 2.2.1 性能需求

- **异步处理**：BlurHash 和多分辨率生成不阻塞上传响应
- **队列支持**：图片处理任务进入队列，后台worker执行
- **并发控制**：限制并发处理数量，防止资源耗尽

#### 2.2.2 安全需求

- **访问控制**：基于 JWT 的认证 + 图片访问令牌
- **权限隔离**：用户只能访问自己的图片
- **路径安全**：随机文件名，防止路径遍历攻击

#### 2.2.3 可扩展性

- **存储可插拔**：通过接口抽象，支持新增存储类型
- **业务解耦**：图片独立于业务，通过关联关系对接
- **接口预留**：分组、标签等功能的接口预留

---

## 3. 数据模型设计

### 3.1 核心实体

#### 3.1.1 图片表 (images)

```go
// Image 图片实体
type Image struct {
    ID            uint            `json:"id" gorm:"primaryKey;comment:主键"`
    OwnerID       uint            `json:"owner_id" gorm:"not null;index:idx_images_owner_id;comment:所有者ID"`
    Filename      string          `json:"filename" gorm:"size:255;not null;comment:随机文件名"`
    OriginalName  string          `json:"original_name" gorm:"size:255;not null;comment:原始文件名"`
    MimeType      string          `json:"mime_type" gorm:"size:100;not null;comment:MIME类型"`
    Size          int64           `json:"size" gorm:"not null;comment:文件大小(字节)"`
    Width         int             `json:"width" gorm:"not null;comment:原始宽度"`
    Height        int             `json:"height" gorm:"not null;comment:原始高度"`
    Blurhash      string          `json:"blurhash" gorm:"size:40;comment:BlurHash字符串"`
    StorageType   string          `json:"storage_type" gorm:"size:20;default:'local';comment:存储类型"`
    StoragePath   string          `json:"storage_path" gorm:"size:500;comment:存储路径"`
    IsPublic      bool            `json:"is_public" gorm:"default:false;comment:是否公开"`
    AccessToken   string          `json:"access_token" gorm:"size:64;comment:访问令牌"`
    Status        string          `json:"status" gorm:"size:20;default:'processing';comment:状态"`
    CreatedAt     time.Time       `json:"created_at" gormorm:"comment:创建时间"`
    UpdatedAt     time.Time       `json:"updated_at" gorm:"comment:更新时间"`

    // Associations
    Variants      []ImageVariant  `json:"variants,omitempty" gorm:"foreignKey:ImageID"`
    Groups        []ImageGroup    `json:"groups,omitempty" gorm:"foreignKey:ImageID"`
}

// TableName specifies the table name for Image
func (Image) TableName() string {
    return "images"
}

// ImageStatus 图片状态
type ImageStatus string

const (
    ImageStatusProcessing ImageStatus = "processing" // 处理中
    ImageStatusReady      ImageStatus = "ready"       // 就绪
    ImageStatusFailed     ImageStatus = "failed"      // 失败
)
```

#### 3.1.2 图片变体表 (image_variants)

```go
// ImageVariant 图片变体（多分辨率）
type ImageVariant struct {
    ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
    ImageID   uint      `json:"image_id" gorm:"not null;index:idx_variants_image_id;comment:图片ID"`
    Variant  string    `json:"variant" gorm:"size:20;not null;comment:变体类型"`
    Width     int       `json:"width" gorm:"not null;comment:宽度"`
    Height    int       `json:"height" gorm:"not null;comment:高度"`
    FileSize  int64     `json:"file_size" gorm:"not null;comment:文件大小"`
    FilePath  string    `json:"file_path" gorm:"size:500;not null;comment:文件路径"`
    CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`

    // Associations
    Image *Image `json:"-" gorm:"foreignKey:ImageID"`
}

// TableName specifies the table name for ImageVariant
func (ImageVariant) TableName() string {
    return "image_variants"
}

// ImageVariantType 变体类型
type ImageVariantType string

const (
    VariantOriginal ImageVariantType = "original"
    VariantMedium   ImageVariantType = "medium"
    VariantSmall    ImageVariantType = "small"
)
```

#### 3.1.3 图片分组表 (image_groups)

```go
// ImageGroup 图片分组
type ImageGroup struct {
    ID          uint      `json:"id" gorm:"primaryKey;comment:主键"`
    OwnerID     uint      `json:"owner_id" gorm:"not null;index:idx_groups_owner_id;comment:所有者ID"`
    GroupType   string    `json:"group_type" gorm:"size:20;not null;index:idx_groups_type;comment:分组类型"`
    GroupName   string    `json:"group_name" gorm:"size:255;comment:分组名称"`
    RefID       *uint     `json:"ref_id" gorm:"index:idx_groups_ref_id;comment:关联ID(如book_id)"`
    RefType     *string   `json:"ref_type" gorm:"size:50;comment:关联类型(如book)"`
    CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"comment:更新时间"`

    // Associations
    Images      []Image `json:"images,omitempty" gorm:"many2many:image_group_members;"`
}

// TableName specifies the table name for ImageGroup
func (ImageGroup) TableName() string {
    return "image_groups"
}

// ImageGroupType 分组类型
type ImageGroupType string

const (
    GroupTypeBook    ImageGroupType = "book"     // 图书分组
    GroupTypeAlbum   ImageGroupType = "album"    // 相册（预留）
    GroupTypeCustom  ImageGroupType = "custom"   // 自定义（预留）
)
```

#### 3.1.4 分组成员关联表 (image_group_members)

```go
// ImageGroupMember 图片分组成员（多对多）
type ImageGroupMember struct {
    GroupID  uint `json:"group_id" gorm:"primaryKey;comment:分组ID"`
    ImageID  uint `json:"image_id" gorm:"primaryKey;comment:图片ID"`
    SortOrder int  `json:"sort_order" gorm:"default:0;comment:排序"`
    CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName specifies the table name for ImageGroupMember
func (ImageGroupMember) TableName() string {
    return "image_group_members"
}
```

### 3.2 ER图

```
┌─────────────┐       ┌────────────────┐       ┌─────────────┐
│   users     │       │ image_groups   │       │   images     │
└─────────────┘       └────────────────┘       └─────────────┘
       │                     │                        │
       │ owner_id            │ ref_id                 │ owner_id
       │                     │                        │
       ▼                     ▼                        ▼
┌────────────────────────────────────────────────────────────────┐
│                                                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ image_group_members (join table)                           │  │
│  │  - group_id + image_id (composite primary key)              │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ image_variants                                             │  │
│  │  - id (primary key)                                       │  │
│  │  - image_id (foreign key → images.id)                     │  │
│  │  - variant (original/medium/small)                         │  │
│  └──────────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────────┘
```

---

## 4. 存储架构设计

### 4.1 存储接口定义

```go
// Storage 存储接口
type Storage interface {
    // 上传文件
    Upload(ctx context.Context, reader io.Reader, path string) error

    // 下载文件
    Download(ctx context.Context, path string) (io.ReadCloser, error)

    // 删除文件
    Delete(ctx context.Context, path string) error

    // 检查文件是否存在
    Exists(ctx context.Context, path string) (bool, error)

    // 获取文件URL
    GetURL(path string) string

    // 获取存储类型
    GetType() string
}
```

### 4.2 本地存储实现

```go
// LocalStorage 本地文件系统存储
type LocalStorage struct {
    basePath string
    baseURL  string
}

func (s *LocalStorage) Upload(ctx context.Context, reader io.Reader, path string) error {
    fullPath := filepath.Join(s.basePath, path)
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return err
    }

    file, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, reader)
    return err
}

// ... 其他方法实现
```

### 4.3 存储配置

```yaml
storage:
  default: local

  local:
    base_path: "./data/images"
    base_url: "/images"

  # 阿里云OSS（预留）
  oss:
    endpoint: ""
    access_key_id: ""
    access_key_secret: ""
    bucket_name: ""

  # AWS S3（预留）
  s3:
    region: ""
    access_key_id: ""
    secret_access_key: ""
    bucket: ""
```

### 4.4 文件路径规范

```
/images/book/{book_id}/{variant_type}/{random_filename}.{ext}

示例：
/images/book/123/small/a3f5c8d2e9b1f4a6.jpg
/images/book/123/medium/a3f5c8d2e9b1f4a6.jpg
/images/book/123/original/a3f5c8d2e9b1f4a6.jpg
```

---

## 5. 渐进式加载方案

### 5.1 BlurHash 方案

**生成时机**：图片上传后异步生成

**技术选型**：`github.com/nfnt/blurhash`

**存储**：存储于 `images.blurhash` 字段（VARCHAR(40)）

**生成流程**：
1. 接收原图上传请求
2. 保存原图到本地存储
3. 创建 Image 记录，状态为 `processing`
4. 异步任务：
   - 读取原图
   - 缩放到 32x32
   - 生成 BlurHash 字符串
   - 更新 Image 状态为 `ready`

### 5.2 多分辨率方案

**变体规格**：

| 变体 | 最大宽度 | 最大高度 | 用途 |
|------|---------|---------|------|
| original | 原始 | 原始 | 高清查看 |
| medium | 1200px | 1200px | 常规浏览 |
| small | 400px | 400px | 缩略图、列表视图 |

**生成流程**：
1. 原图上传后保存到 `original/` 目录
2. 异步任务生成 medium 和 small 变体
3. 将变体信息写入 `image_variants` 表

**技术实现**：使用 `github.com/disintegration/imaging`

```go
// 生成多分辨率变体
func generateVariants(imagePath string) (mediumPath, smallPath string, mediumW, mediumH, smallW, smallH int, err error) {
    // 打开原图
    src, err := imaging.Open(imagePath)
    if err != nil {
        return "", "", 0, 0, 0, 0, err
    }
    defer src.Close()

    bounds := src.Bounds()
    originalW := bounds.Dx()
    originalH := bounds.Dy()

    // 生成中图
    mediumPath = strings.Replace(imagePath, "/original/", "/medium/", 1)
    mediumDst := imaging.Resize(src, 1200, 0, imaging.Lanczos)
    err = imaging.Save(mediumDst, mediumPath)
    if err != nil {
        return "", "", 0, 0, 0, 0, err
    }

    // 生成小图
    smallPath = strings.Replace(imagePath, "/original/", "/small/", 1)
    smallDst := imaging.Resize(src, 400, 0, imaging.Lanczos)
    err = imaging.Save(smallDst, smallPath)
    if err != nil {
        return "", "", 0, 0, 0, 0, err
    }

    return mediumPath, smallPath, 1200, 0, 400, 0, nil
}
```

### 5.3 前端加载流程

```
1. 显示 BlurHash 模糊占位符（< 1KB）
   ↓
2. 并行加载 small 变体（~10-50KB）
   ↓
3. 用户交互或进入视口时加载 medium 变体（~100-500KB）
   ↓
4. 需要时加载原图（可能数MB）
```

---

## 6. API 设计

### 6.1 图片上传

**请求**：
```
POST /api/v1/images
Content-Type: multipart/form-data

FormData:
- file: (binary) 图片文件
- group_id: (optional) 分组ID
```

**响应**：
```json
{
  "success": true,
  "data": {
    "id": 123,
    "filename": "a3f5c8d2e9b1f4a6.jpg",
    "original_name": "book_cover.jpg",
    "mime_type": "image/jpeg",
    "size": 2048576,
    "width": 1920,
    "height": 1080,
    "blurhash": "LlL#0H}3RkWB9FbbWB9jFfFfFf",
    "status": "processing",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 6.2 图片访问

**获取图片信息**：
```
GET /api/v1/images/:id
Authorization: Bearer {jwt_token}
```

**访问图片文件**：
```
GET /api/v1/images/:id/file?variant=small&token={access_token}
```

**响应**：直接返回图片二进制流

### 6.3 图片列表

**请求**：
```
GET /api/v1/images?page=1&page_size=20&group_id={group_id}
Authorization: Bearer {jwt_token}
```

**响应**：
```json
{
  "success": true,
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 6.4 图片删除

**请求**：
```
DELETE /api/v1/images/:id
Authorization: Bearer {jwt_token}
```

### 6.5 分组管理

**创建分组**：
```
POST /api/v1/images/groups
Authorization: Bearer {jwt_token}

{
  "group_type": "book",
  "group_name": "我的图书封面",
  "ref_id": 123,
  "ref_type": "book"
}
```

**添加图片到分组**：
```
POST /api/v1/images/groups/:id/images
Authorization: Bearer {jwt_token}

{
  "image_ids": [1, 2, 3]
}
```

---

## 7. 技术选型

### 7.1 图片处理库

| 库 | Star数 | 用途 | 推荐理由 |
|---|--------|------|----------|
| [disintegration/imaging](https://github.com/disintegration/imaging) | 9k+ | 图片处理 | 纯Go实现，活跃维护，性能好 |
| [nfnt/blurhash](https://github.com/nfnt/blurhash) | 3k+ | BlurHash | 官方库，稳定可靠 |

### 7.2 文件上传

| 功能 | 技术选型 |
|------|---------|
| 单文件上传 | multipart/form-data |
| 断点续传 | 自定义 chunk 上传 |
| 进度推送 | WebSocket/SSE |
| 并发控制 | 信号量 + Worker Pool |

### 7.3 任务队列

| 方案 | 说明 |
|------|------|
| 第一阶段 | 内存队列 + Goroutine Worker |
| 第二阶段（可选） | Redis List + 自定义 Worker |
| 第三阶段（可选） | RabbitMQ / Redis Stream |

---

## 8. 实施计划

### 8.1 第一阶段：核心功能（MVP）

**目标**：实现基本的图片上传、存储、访问功能

**功能清单**：
1. 数据模型创建（images、image_variants 表）
2. 本地存储实现
3. 图片上传接口（单文件，无进度）
4. 图片访问接口（签名URL）
5. 基础分组功能（按图书）

**预计工作量**：2-3 周

### 8.2 第二阶段：渐进式加载

**目标**：实现 BlurHash + 多分辨率

**功能清单**：
1. BlurHash 生成（异步）
2. 多分辨率生成（异步）
3. 图片处理队列实现
4. Worker Pool 实现

**预计工作量**：2 周

### 8.3 第三阶段：扩展功能

**目标**：存储扩展、高级分组

**功能清单**：
1. 云存储接口（OSS/S3）
2. 批量上传
3. 高级分组（相册、标签）
4. 图片分享功能

**预计工作量**：2-3 周

### 8.4 数据迁移计划

**现有 Book.Cover 字段迁移**：

1. 创建 `image_migrations` 表记录迁移状态
2. 扫描所有有 Cover 值的 Book
3. 上传封面到图片系统
4. 生成 BlurHash 和多分辨率
5. 更新 Book.Cover 为新的 Image ID
6. 记录迁移状态

---

## 9. 风险与挑战

### 9.1 技术风险

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 大文件处理超时 | 上传失败 | 断点续传、超时重试 |
| 异步处理失败 | 变体未生成 | 失败重试、监控告警 |
| 存储迁移 | 数据丢失 | 备份机制、灰度迁移 |

### 9.2 业务风险

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 用户习惯改变 | 使用体验下降 | 渐进式灰度发布 |
| 存储成本增加 | 运营成本 | 定期清理无用图片 |
| 性能问题 | 系统响应变慢 | 缓存优化、CDN加速 |

---

## 附录

### A. 术语表

| 术语 | 说明 |
|------|------|
| BlurHash | 一种图片哈希算法，生成短字符串，用于生成模糊占位图 |
| 多分辨率 | 同一图片生成不同尺寸的版本 |
| 变体 | 图片的不同分辨率版本 |
| 签名URL | 带有访问令牌的URL，用于私有资源访问 |

### B. 参考文档

- [imaging 库文档](https://github.com/disintegration/imaging)
- [BlurHash 算法说明](https://blurha.sh/)
- [GORM 中文文档](https://gorm.cn/)

### C. 数据库迁移 SQL

```sql
-- 图片表
CREATE TABLE `images` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
  `filename` VARCHAR(255) NOT NULL COMMENT '随机文件名',
  `original_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
  `mime_type` VARCHAR(100) NOT NULL COMMENT 'MIME类型',
  `size` BIGINT NOT NULL COMMENT '文件大小(字节)',
  `width` INT NOT NULL COMMENT '原始宽度',
  `height` INT NOT NULL COMMENT '原始高度',
  `blurhash` VARCHAR(40) DEFAULT NULL COMMENT 'BlurHash字符串',
  `storage_type` VARCHAR(20) DEFAULT 'local' COMMENT '存储类型',
  `storage_path` VARCHAR(500) DEFAULT NULL COMMENT '存储路径',
  `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开',
  `access_token` VARCHAR(64) NOT NULL COMMENT '访问令牌',
  `status` VARCHAR(20) DEFAULT 'processing' COMMENT '状态',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_images_owner_id` (`owner_id`),
  KEY `idx_images_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片表';

-- 图片变体表
CREATE TABLE `image_variants` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `image_id` BIGINT UNSIGNED NOT NULL COMMENT '图片ID',
  `variant` VARCHAR(20) NOT NULL COMMENT '变体类型',
  `width` INT NOT NULL COMMENT '宽度',
  `height` INT NOT NULL COMMENT '高度',
  `file_size` BIGINT NOT NULL COMMENT '文件大小',
  `file_path` VARCHAR(500) NOT NULL COMMENT '文件路径',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_variants_image_id` (`image_id`),
  KEY `idx_variants_variant` (`variant`),
  CONSTRAINT `fk_variants_image` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片变体表';

-- 图片分组表
CREATE TABLE `image_groups` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
  `group_type` VARCHAR(20) NOT NULL COMMENT '分组类型',
  `group_name` VARCHAR(255) DEFAULT NULL COMMENT '分组名称',
  `ref_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联ID',
  `ref_type` VARCHAR(50) DEFAULT NULL COMMENT '关联类型',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_groups_owner_id` (`owner_id`),
  KEY `idx_groups_type` (`group_type`),
  KEY `idx_groups_ref_id` (`ref_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片分组表';

-- 分组成员关联表
CREATE TABLE `image_group_members` (
  `group_id` BIGINT UNSIGNED NOT NULL COMMENT '分组ID',
  `image_id` BIGINT UNSIGNED NOT NULL COMMENT '图片ID',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`group_id`, `image_id`),
  KEY `idx_members_image_id` (`image_id`),
  CONSTRAINT `fk_members_group` FOREIGN KEY (`group_id`) REFERENCES `image_groups` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_members_image` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分组成员关联表';
```

### D. 目录结构

```
internal/
├── model/
│   ├── image.go           # 图片模型定义
│   ├── image_variant.go  # 图片变体模型
│   └── image_group.go     # 图片分组模型
├── repository/
│   └── image_repository.go
├── service/
│   ├── image_service.go       # 图片业务逻辑
│   ├── storage/
│   │   ├── storage.go        # 存储接口
│   │   ├── local.go          # 本地存储实现
│   │   └── oss.go            # 云存储（预留）
│   └── processor/
│       ├── blurhash.go       # BlurHash生成
│       └── variant.go        # 多分辨率生成
├── handler/
│   └── image_handler.go       # HTTP 处理器
└── worker/
    └── image_worker.go        # 异步处理Worker

data/
└── images/
    └── book/
        └── {book_id}/
            ├── original/
            ├── medium/
            └── small/
```

### E. 配置示例

```yaml
# 图片服务配置
image:
  # 存储配置
  storage:
    default: local
    local:
      base_path: "./data/images"
      base_url: "/images"

  # 处理配置
  processing:
    blurhash_enabled: true
    variants:
      - name: medium
        max_width: 1200
        max_height: 1200
      - name: small
        max_width: 400
        max_height: 400

  # Worker配置
  worker:
    concurrent: 3
    queue_size: 1000
    retry_times: 3
    retry_interval: 60s

  # 上传配置
  upload:
    chunk_size: 1048576  # 1MB
    max_file_size: 52428800  # 50MB
    allowed_types:
      - image/jpeg
      - image/png
      - image/gif
      - image/webp
```
