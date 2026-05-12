package model

import (
	"time"

	"gorm.io/gorm"
)

type AIUsageLog struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID       string    `gorm:"type:varchar(36);not null;index:idx_ai_usage_user_created,priority:1" json:"user_id"`
	FamilyID     string    `gorm:"type:varchar(36);index:idx_ai_usage_family_created,priority:1" json:"family_id,omitempty"`
	Provider     string    `gorm:"type:varchar(64);not null;index:idx_ai_usage_provider_created,priority:1" json:"provider"`
	Model        string    `gorm:"type:varchar(128);not null" json:"model"`
	InputTokens  int       `gorm:"not null;default:0" json:"input_tokens"`
	OutputTokens int       `gorm:"not null;default:0" json:"output_tokens"`
	CachedTokens int       `gorm:"not null;default:0" json:"cached_tokens"`
	TotalTokens  int       `gorm:"not null;default:0" json:"total_tokens"`
	RequestType  string    `gorm:"type:varchar(32);not null;index" json:"request_type"`
	Status       string    `gorm:"type:varchar(32);not null;index" json:"status"`
	ErrorMessage string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time `gorm:"index:idx_ai_usage_user_created,priority:2;index:idx_ai_usage_family_created,priority:2;index:idx_ai_usage_provider_created,priority:2" json:"created_at"`
}

func (a *AIUsageLog) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = NewID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil
}

func (AIUsageLog) TableName() string { return "ai_usage_logs" }
