-- iBookFS Database Schema
-- Generated from MySQL MCP tool
-- Database: ibookfs

-- =====================================================
-- 1. users - 用户表
-- =====================================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
  `first_name` VARCHAR(100) NOT NULL COMMENT '名',
  `last_name` VARCHAR(100) NOT NULL COMMENT '姓',
  `display_name` VARCHAR(255) DEFAULT NULL COMMENT '显示名称',
  `avatar_url` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `status` VARCHAR(20) DEFAULT 'active' COMMENT '状态 (active=活跃, suspended=暂停, deleted=已删除)',
  `role` VARCHAR(20) DEFAULT 'user' COMMENT '角色 (user=用户, admin=管理员)',
  `preferences` LONGTEXT DEFAULT NULL COMMENT '用户偏好设置 (JSON)',
  `language` VARCHAR(10) DEFAULT 'zh-CN' COMMENT '语言设置',
  `timezone` VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '时区',
  `last_login_at` DATETIME(3) DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_status` (`status`),
  KEY `idx_users_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 2. sessions - 会话表
-- =====================================================
CREATE TABLE IF NOT EXISTS `sessions` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `token_hash` VARCHAR(64) NOT NULL COMMENT '访问令牌哈希',
  `refresh_token_hash` VARCHAR(64) DEFAULT NULL COMMENT '刷新令牌哈希',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理（浏览器信息）',
  `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
  `device_type` VARCHAR(50) DEFAULT NULL COMMENT '设备类型 (mobile, tablet, desktop)',
  `device_name` VARCHAR(100) DEFAULT NULL COMMENT '设备名称',
  `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
  `last_activity_at` DATETIME(3) DEFAULT NULL COMMENT '最后活跃时间',
  `revoked` TINYINT(1) DEFAULT 0 COMMENT '是否已撤销 (0=正常, 1=已撤销)',
  `revoked_at` DATETIME(3) DEFAULT NULL COMMENT '撤销时间',
  `revoked_reason` VARCHAR(255) DEFAULT NULL COMMENT '撤销原因',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sessions_token_hash` (`token_hash`),
  KEY `idx_sessions_user_id` (`user_id`),
  KEY `idx_sessions_expires` (`expires_at`),
  KEY `idx_sessions_revoked` (`revoked`),
  CONSTRAINT `fk_users_sessions` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 3. oauth_identities - OAuth身份绑定表
-- =====================================================
CREATE TABLE IF NOT EXISTS `oauth_identities` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `provider` VARCHAR(20) NOT NULL COMMENT 'OAuth提供商 (github, gitee)',
  `provider_user_id` VARCHAR(255) NOT NULL COMMENT '第三方平台用户ID',
  `provider_username` VARCHAR(255) DEFAULT NULL COMMENT '第三方平台用户名',
  `provider_email` VARCHAR(255) DEFAULT NULL COMMENT '第三方平台邮箱',
  `access_token` TEXT DEFAULT NULL COMMENT 'OAuth访问令牌',
  `refresh_token` TEXT DEFAULT NULL COMMENT 'OAuth刷新令牌',
  `token_expires_at` DATETIME(3) DEFAULT NULL COMMENT '令牌过期时间',
  `scope` VARCHAR(500) DEFAULT NULL COMMENT 'OAuth授权范围',
  `profile_data` LONGTEXT DEFAULT NULL COMMENT '第三方平台用户资料 (JSON)',
  `is_primary` TINYINT(1) DEFAULT 0 COMMENT '是否为主账号 (0=否, 1=是)',
  `linked_at` DATETIME(3) DEFAULT NULL COMMENT '绑定时间',
  `last_used_at` DATETIME(3) DEFAULT NULL COMMENT '最后使用时间',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_oauth_provider_user` (`provider`, `provider_user_id`),
  KEY `idx_oauth_user_id` (`user_id`),
  KEY `idx_oauth_provider` (`provider`),
  KEY `idx_oauth_is_primary` (`is_primary`),
  CONSTRAINT `fk_users_o_auth_identities` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 4. login_history - 登录历史表
-- =====================================================
CREATE TABLE IF NOT EXISTS `login_history` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT(20) UNSIGNED DEFAULT NULL COMMENT '用户ID',
  `login_method` VARCHAR(20) NOT NULL COMMENT '登录方式 (email, oauth, password)',
  `success` TINYINT(1) DEFAULT NULL COMMENT '是否成功 (0=失败, 1=成功)',
  `failure_reason` VARCHAR(255) DEFAULT NULL COMMENT '失败原因',
  `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
  `device_fingerprint` VARCHAR(255) DEFAULT NULL COMMENT '设备指纹',
  `country` VARCHAR(100) DEFAULT NULL COMMENT '国家',
  `region` VARCHAR(100) DEFAULT NULL COMMENT '地区/省份',
  `city` VARCHAR(100) DEFAULT NULL COMMENT '城市',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_login_user_id` (`user_id`),
  KEY `idx_login_success` (`success`),
  KEY `idx_login_created_at` (`created_at`),
  CONSTRAINT `fk_login_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 5. email_verification_codes - 邮箱验证码表
-- =====================================================
CREATE TABLE IF NOT EXISTS `email_verification_codes` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `email` VARCHAR(255) NOT NULL COMMENT '邮箱地址',
  `code` VARCHAR(10) NOT NULL COMMENT '验证码（6位数字）',
  `type` VARCHAR(20) NOT NULL COMMENT '类型 (register=注册, login=登录, reset_password=重置密码)',
  `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
  `used` TINYINT(1) DEFAULT 0 COMMENT '是否已使用 (0=未使用, 1=已使用)',
  `attempts` INT(11) DEFAULT 0 COMMENT '验证尝试次数',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_verification_email_type` (`email`, `type`),
  KEY `idx_verification_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- 6. books - 图书表
-- =====================================================
CREATE TABLE IF NOT EXISTS `books` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `title` VARCHAR(255) NOT NULL COMMENT '书名',
  `author` VARCHAR(255) DEFAULT NULL COMMENT '作者',
  `isbn` VARCHAR(20) DEFAULT NULL COMMENT 'ISBN编号',
  `publisher` VARCHAR(255) DEFAULT NULL COMMENT '出版社',
  `year` BIGINT(20) DEFAULT NULL COMMENT '出版年份',
  `pages` BIGINT(20) DEFAULT NULL COMMENT '总页数',
  `cover` VARCHAR(500) DEFAULT NULL COMMENT '封面图片路径',
  `uploaded_pages` BIGINT(20) DEFAULT 0 COMMENT '已上传页数',
  `status` VARCHAR(20) DEFAULT 'draft' COMMENT '状态 (draft=草稿, processing=处理中, completed=完成, failed=失败)',
  `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_books_isbn` (`isbn`),
  KEY `idx_books_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- Indexes for better query performance
-- =====================================================

-- Users table indexes are already defined above

-- Additional useful indexes could be added based on query patterns
-- CREATE INDEX idx_users_last_login ON users(last_login_at DESC);
-- CREATE INDEX idx_sessions_last_activity ON sessions(last_activity_at DESC);
-- CREATE INDEX idx_oauth_last_used ON oauth_identities(last_used_at DESC);
