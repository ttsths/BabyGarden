-- ============================================
-- Yuanzi 数据库迁移脚本 v1.0
-- 版本: 001
-- 日期: 2026-03-09
-- 描述: 用户系统表结构 (users, babies, family_members)
-- ============================================

-- 1. 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
  `avatar_url` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 宝宝表
CREATE TABLE IF NOT EXISTS `babies` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '宝宝ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '主要照顾者用户ID',
  `nickname` VARCHAR(50) NOT NULL COMMENT '宝宝昵称',
  `avatar_url` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `birth_date` DATE NOT NULL COMMENT '出生日期',
  `gender` TINYINT NOT NULL DEFAULT 0 COMMENT '性别: 0-未知, 1-男, 2-女',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_birth_date` (`birth_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='宝宝表';

-- 3. 家庭表
CREATE TABLE IF NOT EXISTS `families` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '家庭ID',
  `name` VARCHAR(100) NOT NULL COMMENT '家庭名称',
  `invite_code` VARCHAR(20) NOT NULL COMMENT '邀请码',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_invite_code` (`invite_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭表';

-- 4. 家庭成员表
CREATE TABLE IF NOT EXISTS `family_members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '成员ID',
  `family_id` BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `role` VARCHAR(20) NOT NULL DEFAULT 'member' COMMENT '角色: admin(管理员), member(成员)',
  `baby_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联宝宝ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_user` (`family_id`, `user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_family_id` (`family_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭成员表';

-- 初始化数据：创建默认家庭（可选）
-- INSERT INTO `families` (`name`, `invite_code`) VALUES ('默认家庭', 'DEFAULT001');
