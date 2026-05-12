-- BabyGarden AI 请求审计日志
-- 用于排查 fallback、观测 Provider 稳定性、估算 Cloudflare 免费额度消耗

CREATE TABLE IF NOT EXISTS ai_request_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    family_id BIGINT DEFAULT NULL,
    provider VARCHAR(64) NOT NULL COMMENT 'grokai|cloudflare_workers_ai|dashscope|cliproxyapi',
    model VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL COMMENT 'success|error|fallback|rate_limited',
    latency_ms INT NOT NULL DEFAULT 0,
    input_tokens INT NOT NULL DEFAULT 0,
    output_tokens INT NOT NULL DEFAULT 0,
    neurons_est INT NOT NULL DEFAULT 0 COMMENT 'Cloudflare Workers AI 估算 Neurons',
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_created (user_id, created_at),
    INDEX idx_provider_created (provider, created_at),
    INDEX idx_request_id (request_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
