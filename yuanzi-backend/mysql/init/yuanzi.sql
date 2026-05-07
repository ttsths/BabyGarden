-- Yuanzi 初始化数据库脚本
-- 执行此脚本创建数据库和基础表

-- 创建数据库
CREATE DATABASE IF NOT EXISTS yuanzi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE yuanzi;

-- 1. 用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY COMMENT '用户ID',
    phone VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号',
    nickname VARCHAR(50) COMMENT '昵称',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
    is_admin TINYINT DEFAULT 0 COMMENT '是否管理员: 0-否 1-是',
    password VARCHAR(255) COMMENT '管理员密码',
    last_login_at DATETIME(3) COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_phone (phone),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 家庭表
CREATE TABLE IF NOT EXISTS families (
    id VARCHAR(36) PRIMARY KEY COMMENT '家庭ID',
    name VARCHAR(100) NOT NULL COMMENT '家庭名称',
    invite_code VARCHAR(8) UNIQUE NOT NULL COMMENT '8位邀请码',
    created_by VARCHAR(36) NOT NULL COMMENT '创建者ID',
    is_paid TINYINT DEFAULT 0 COMMENT '是否付费',
    storage_limit BIGINT DEFAULT 1073741824 COMMENT '存储配额(字节)',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_invite_code (invite_code),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭表';

-- 3. 家庭成员表
CREATE TABLE IF NOT EXISTS family_members (
    id VARCHAR(36) PRIMARY KEY COMMENT '家庭成员ID',
    family_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role VARCHAR(20) DEFAULT 'member' COMMENT '角色: admin/member/elder',
    elder_mode TINYINT DEFAULT 0 COMMENT '祖辈模式开关',
    notifications JSON DEFAULT '{"feed": true, "sleep": true}' COMMENT '推送通知配置',
    joined_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_family_user (family_id, user_id),
    INDEX idx_family_id (family_id),
    INDEX idx_user_id (user_id),
    INDEX idx_family_role (family_id, role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='家庭成员表';

-- 4. 宝宝表
CREATE TABLE IF NOT EXISTS babies (
    id VARCHAR(36) PRIMARY KEY COMMENT '宝宝ID',
    family_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL COMMENT '宝宝姓名/昵称',
    birthday DATE NOT NULL COMMENT '出生日期',
    gender TINYINT NOT NULL COMMENT '性别: 1-男 2-女',
    birth_weight DECIMAL(5,2) COMMENT '出生体重(kg)',
    birth_height DECIMAL(4,1) COMMENT '出生身高(cm)',
    avatar_url VARCHAR(500) COMMENT '宝宝头像',
    note TEXT COMMENT '备注',
    is_premature TINYINT DEFAULT 0 COMMENT '是否早产',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_family_id (family_id),
    INDEX idx_birthday (birthday)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='宝宝表';

-- 5. 记录表
CREATE TABLE IF NOT EXISTS records (
    id VARCHAR(36) PRIMARY KEY COMMENT '记录ID',
    baby_id VARCHAR(36) NOT NULL,
    family_id VARCHAR(36) NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT '类型: feeding/sleep/diaper/growth',
    started_at DATETIME(3) NOT NULL COMMENT '开始时间',
    ended_at DATETIME(3) COMMENT '结束时间',
    content JSON NOT NULL COMMENT '类型特定内容',
    note TEXT COMMENT '备注',
    created_by VARCHAR(36) NOT NULL COMMENT '创建者ID',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) COMMENT '软删除时间',
    INDEX idx_baby_id (baby_id),
    INDEX idx_family_id (family_id),
    INDEX idx_type (type),
    INDEX idx_started_at (started_at),
    INDEX idx_baby_started (baby_id, started_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='记录表';

-- 6. 照片表
CREATE TABLE IF NOT EXISTS photos (
    id VARCHAR(36) PRIMARY KEY COMMENT '照片ID',
    baby_id VARCHAR(36) NOT NULL,
    family_id VARCHAR(36) NOT NULL,
    oss_key VARCHAR(500) NOT NULL COMMENT 'OSS对象key',
    thumbnail_key VARCHAR(500) COMMENT '缩略图key',
    width INT COMMENT '宽度',
    height INT COMMENT '高度',
    size BIGINT NOT NULL COMMENT '文件大小(bytes)',
    content_type VARCHAR(50) DEFAULT 'image/jpeg' COMMENT '文件类型',
    taken_at DATETIME(3) COMMENT '拍摄时间',
    description TEXT COMMENT '描述',
    uploaded_by VARCHAR(36) NOT NULL COMMENT '上传者ID',
    uploaded_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态: pending/active/deleted',
    INDEX idx_baby_id (baby_id),
    INDEX idx_family_id (family_id),
    INDEX idx_taken_at (taken_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片表';

-- 7. AI问答记录表
CREATE TABLE IF NOT EXISTS ai_chat_records (
    id VARCHAR(36) PRIMARY KEY COMMENT '问答记录ID',
    user_id VARCHAR(36) NOT NULL,
    baby_id VARCHAR(36),
    question TEXT NOT NULL COMMENT '问题',
    answer TEXT NOT NULL COMMENT '回答',
    voice_url VARCHAR(500) COMMENT '语音输入URL',
    tokens_used INT COMMENT '消耗token数',
    model VARCHAR(50) DEFAULT 'qwen-turbo',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_user_id (user_id),
    INDEX idx_baby_id (baby_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI问答记录表';

-- 8. 验证码表
CREATE TABLE IF NOT EXISTS verification_codes (
    id VARCHAR(36) PRIMARY KEY COMMENT '验证码ID',
    phone VARCHAR(11) NOT NULL,
    code VARCHAR(6) NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT '类型: login/reset/bind',
    expires_at DATETIME(3) NOT NULL COMMENT '过期时间',
    used_at DATETIME(3) COMMENT '使用时间',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_phone (phone, type),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证码表';

-- 9. 推送设备表
CREATE TABLE IF NOT EXISTS push_devices (
    id VARCHAR(36) PRIMARY KEY COMMENT '设备ID',
    user_id VARCHAR(36) NOT NULL,
    platform VARCHAR(20) NOT NULL COMMENT '平台: ios/android',
    device_token VARCHAR(255) NOT NULL COMMENT '设备token',
    alias VARCHAR(100) COMMENT '别名',
    tags JSON COMMENT '标签JSON数组',
    is_active TINYINT DEFAULT 1 COMMENT '是否激活',
    last_used_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_user_device (user_id, device_token),
    INDEX idx_user_id (user_id),
    INDEX idx_device_token (device_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='推送设备表';

-- 创建示例用户和家庭
INSERT INTO users (id, phone, nickname, is_admin, password) VALUES 
('100e8400-e29b-41d4-a716-446655440000', '13800138000', '妈妈', 1, 'admin123'),
('100e8400-e29b-41d4-a716-446655440010', '139001390010', '奶奶');

-- 创建示例家庭
INSERT INTO families (id, name, invite_code, created_by, is_paid) VALUES 
('200e8400-e29b-41d4-a716-446655440000', '小园子的家', 'ABC12345', '100e8400-e29b-41d4-a716-446655440000', 0);

-- 创建家庭成员
INSERT INTO family_members (id, family_id, user_id, role) VALUES 
('300e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440000', 'admin'),
('300e8400-e29b-41d4-a716-446655440010', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440010', 'member'];

-- 创建示例宝宝
INSERT INTO babies (id, family_id, name, birthday, gender, birth_weight, birth_height) VALUES 
('400e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', '小园子', '2024-01-01', 2, 3.20, 50.0);

-- 初始化完成
SELECT 'Database yuanzi initialized successfully!' AS status;
