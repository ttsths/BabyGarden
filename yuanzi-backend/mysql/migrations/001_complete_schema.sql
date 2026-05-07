-- ============================================
-- Yuanzi 数据库迁移脚本 v1.1
-- 版本: 001
-- 日期: 2026-03-09
-- 描述: 完整表结构（基于 PRD-V1.1）
-- ============================================

-- 1. 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `phone` VARCHAR(11) NOT NULL COMMENT '手机号',
  `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
  `avatar_url` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
  `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 验证码表
CREATE TABLE IF NOT EXISTS `verification_codes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `phone` VARCHAR(11) NOT NULL COMMENT '手机号',
  `code` VARCHAR(6) NOT NULL COMMENT '验证码',
  `type` VARCHAR(20) NOT NULL DEFAULT 'login' COMMENT '类型: login/reset/bind',
  `expires_at` DATETIME NOT NULL COMMENT '过期时间',
  `used_at` DATETIME DEFAULT NULL COMMENT '使用时间',
  `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_phone_type` (`phone`, `type`),
  KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证码表';

-- 3. 家庭表
CREATE TABLE IF NOT EXISTS `families` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '家庭ID',
  `name` VARCHAR(100) NOT NULL COMMENT '家庭名称',
  `invite_code` VARCHAR(8) NOT NULL COMMENT '8位邀请码',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
  `is_paid` TINYINT NOT NULL DEFAULT 0 COMMENT '是否付费: 0-否 1-是',
  `storage_limit` BIGINT NOT NULL DEFAULT 1073741824 COMMENT '存储限制(字节), 默认1GB',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_invite_code` (`invite_code`),
  KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭表';

-- 4. 家庭成员表
CREATE TABLE IF NOT EXISTS `family_members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '成员ID',
  `family_id` BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `role` VARCHAR(20) NOT NULL DEFAULT 'member' COMMENT '角色: admin/member/elder',
  `elder_mode` TINYINT NOT NULL DEFAULT 0 COMMENT '祖辈模式: 0-否 1-是',
  `notifications` JSON DEFAULT NULL COMMENT '通知设置',
  `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_user` (`family_id`, `user_id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role` (`family_id`, `role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭成员表';

-- 5. 宝宝表
CREATE TABLE IF NOT EXISTS `babies` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '宝宝ID',
  `family_id` BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
  `name` VARCHAR(50) NOT NULL COMMENT '宝宝姓名',
  `birthday` DATE NOT NULL COMMENT '出生日期',
  `gender` TINYINT NOT NULL COMMENT '性别: 1-男 2-女',
  `birth_weight` DECIMAL(5,2) DEFAULT NULL COMMENT '出生体重(kg)',
  `birth_height` DECIMAL(4,1) DEFAULT NULL COMMENT '出生身高(cm)',
  `avatar_url` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `note` TEXT DEFAULT NULL COMMENT '备注',
  `is_premature` TINYINT NOT NULL DEFAULT 0 COMMENT '是否早产: 0-否 1-是',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_birthday` (`birthday`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='宝宝表';

-- 6. 记录表（核心表 - 喂养/睡眠/排泄/成长）
CREATE TABLE IF NOT EXISTS `records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `baby_id` BIGINT UNSIGNED NOT NULL COMMENT '宝宝ID',
  `family_id` BIGINT UNSIGNED NOT NULL COMMENT '家庭ID(冗余)',
  `type` VARCHAR(20) NOT NULL COMMENT '类型: feeding/sleep/diaper/growth',
  `started_at` DATETIME NOT NULL COMMENT '开始时间',
  `ended_at` DATETIME DEFAULT NULL COMMENT '结束时间',
  `content` JSON NOT NULL COMMENT '类型特定内容',
  `note` TEXT DEFAULT NULL COMMENT '备注',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_baby_id` (`baby_id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_type` (`type`),
  KEY `idx_started_at` (`started_at`),
  KEY `idx_baby_started` (`baby_id`, `started_at`),
  KEY `idx_not_deleted` (`baby_id`, `started_at`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='记录表';

-- 7. 照片表
CREATE TABLE IF NOT EXISTS `photos` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '照片ID',
  `baby_id` BIGINT UNSIGNED NOT NULL COMMENT '宝宝ID',
  `family_id` BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
  `oss_key` VARCHAR(500) NOT NULL COMMENT 'OSS对象key',
  `thumbnail_key` VARCHAR(500) DEFAULT NULL COMMENT '缩略图key',
  `width` INT DEFAULT NULL COMMENT '宽度',
  `height` INT DEFAULT NULL COMMENT '高度',
  `size` BIGINT NOT NULL COMMENT '文件大小(字节)',
  `content_type` VARCHAR(50) NOT NULL DEFAULT 'image/jpeg' COMMENT '内容类型',
  `taken_at` DATETIME DEFAULT NULL COMMENT '照片拍摄时间',
  `description` TEXT DEFAULT NULL COMMENT '描述',
  `uploaded_by` BIGINT UNSIGNED NOT NULL COMMENT '上传者用户ID',
  `uploaded_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: pending/active/deleted',
  PRIMARY KEY (`id`),
  KEY `idx_baby_id` (`baby_id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_taken_at` (`taken_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片表';

-- 8. AI问答记录表
CREATE TABLE IF NOT EXISTS `ai_chat_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `baby_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '宝宝ID',
  `question` TEXT NOT NULL COMMENT '问题',
  `answer` TEXT NOT NULL COMMENT '回答',
  `voice_url` VARCHAR(500) DEFAULT NULL COMMENT '语音输入URL',
  `tokens_used` INT DEFAULT NULL COMMENT '消耗token数',
  `model` VARCHAR(50) NOT NULL DEFAULT 'qwen-turbo' COMMENT '模型',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_baby_id` (`baby_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI问答记录表';

-- 9. 配额使用表
CREATE TABLE IF NOT EXISTS `usage_quotas` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `date` DATE NOT NULL COMMENT '日期',
  `speech_count` INT NOT NULL DEFAULT 0 COMMENT '语音识别次数',
  `ai_chat_count` INT NOT NULL DEFAULT 0 COMMENT 'AI问答次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_date` (`user_id`, `date`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配额使用表';

-- 10. 推送设备表
CREATE TABLE IF NOT EXISTS `push_devices` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `platform` VARCHAR(20) NOT NULL COMMENT '平台: ios/android',
  `device_token` VARCHAR(255) NOT NULL COMMENT '设备token',
  `alias` VARCHAR(100) DEFAULT NULL COMMENT '别名',
  `tags` JSON DEFAULT NULL COMMENT '标签数组',
  `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否激活: 0-否 1-是',
  `last_used_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后使用时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_device` (`user_id`, `device_token`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_token` (`device_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='推送设备表';
