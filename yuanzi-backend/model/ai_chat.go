package model

import (
	"time"

	"gorm.io/gorm"
)

type AIChatRecord struct {
	ID         string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID     string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	BabyID     *string    `gorm:"type:varchar(36);index" json:"baby_id,omitempty"`
	Question   string     `gorm:"type:text;not null" json:"question"`
	Answer     string     `gorm:"type:text;not null" json:"answer"`
	VoiceURL   string     `gorm:"type:varchar(500)" json:"voice_url,omitempty"`
	TokensUsed int        `json:"tokens_used"`
	Model      string     `gorm:"type:varchar(50)" json:"model"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (a *AIChatRecord) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = NewID()
	}
	return nil
}

func (AIChatRecord) TableName() string { return "ai_chat_records" }
