-- 照片互动表：评论与点赞

CREATE TABLE IF NOT EXISTS photo_comments (
    id VARCHAR(36) PRIMARY KEY,
    photo_id VARCHAR(36) NOT NULL,
    family_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    content VARCHAR(500) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3),
    INDEX idx_photo_comments_photo (photo_id),
    INDEX idx_photo_comments_family (family_id),
    INDEX idx_photo_comments_user (user_id),
    INDEX idx_photo_comments_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片评论表';

CREATE TABLE IF NOT EXISTS photo_likes (
    id VARCHAR(36) PRIMARY KEY,
    photo_id VARCHAR(36) NOT NULL,
    family_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_photo_user (photo_id, user_id),
    INDEX idx_photo_likes_family (family_id),
    INDEX idx_photo_likes_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片点赞表';
