-- AI token 使用记录
-- 用于管理端查询、趋势汇总和用户级明细分析

CREATE TABLE IF NOT EXISTS ai_usage_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    family_id VARCHAR(36) DEFAULT NULL,
    provider VARCHAR(64) NOT NULL COMMENT 'grokai|cloudflare_workers_ai|dashscope|cliproxyapi',
    model VARCHAR(128) NOT NULL,
    input_tokens INT NOT NULL DEFAULT 0,
    output_tokens INT NOT NULL DEFAULT 0,
    cached_tokens INT NOT NULL DEFAULT 0,
    total_tokens INT NOT NULL DEFAULT 0,
    request_type VARCHAR(32) NOT NULL COMMENT 'chat|speech',
    status VARCHAR(32) NOT NULL COMMENT 'success|error',
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ai_usage_user_created (user_id, created_at),
    INDEX idx_ai_usage_family_created (family_id, created_at),
    INDEX idx_ai_usage_provider_created (provider, created_at),
    INDEX idx_ai_usage_request_status (request_type, status),
    INDEX idx_ai_usage_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
