-- Add request_id column to ai_usage_logs for cross-provider request tracing
-- Issue: #66

ALTER TABLE ai_usage_logs
    ADD COLUMN request_id VARCHAR(36) DEFAULT NULL AFTER id,
    ADD INDEX idx_ai_usage_request_id (request_id);