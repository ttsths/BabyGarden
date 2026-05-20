-- BabyGarden family live data additions.

ALTER TABLE `records`
  MODIFY COLUMN `type` VARCHAR(20) NOT NULL COMMENT '类型: feeding/sleep/diaper/excretion/temperature/growth';

CREATE TABLE IF NOT EXISTS `photo_comments` (
  `id` VARCHAR(36) NOT NULL COMMENT '评论ID',
  `photo_id` VARCHAR(36) NOT NULL COMMENT '照片ID',
  `family_id` VARCHAR(36) NOT NULL COMMENT '家庭ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '评论用户ID',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_photo_id` (`photo_id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片评论表';

CREATE TABLE IF NOT EXISTS `photo_likes` (
  `id` VARCHAR(36) NOT NULL COMMENT '点赞ID',
  `photo_id` VARCHAR(36) NOT NULL COMMENT '照片ID',
  `family_id` VARCHAR(36) NOT NULL COMMENT '家庭ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '点赞用户ID',
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_photo_user` (`photo_id`, `user_id`),
  KEY `idx_family_id` (`family_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片点赞表';
