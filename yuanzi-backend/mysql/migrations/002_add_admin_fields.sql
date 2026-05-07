-- ============================================
-- Yuanzi 数据库迁移脚本
-- 版本: 002
-- 日期: 2026-05-07
-- 描述: 为 users 表添加管理员字段 (is_admin, password)
-- 关联: PR #24 - 管理后台 API
-- ============================================

ALTER TABLE `users`
  ADD COLUMN IF NOT EXISTS `is_admin` TINYINT DEFAULT 0 COMMENT '是否管理员: 0-否 1-是' AFTER `status`,
  ADD COLUMN IF NOT EXISTS `password` VARCHAR(255) COMMENT '管理员密码(明文存储,生产请用bcrypt)' AFTER `is_admin`;

-- 为示例管理员设置密码 (密码: admin123)
UPDATE `users` 
SET `is_admin` = 1, `password` = 'admin123' 
WHERE `phone` = '13800138000';
