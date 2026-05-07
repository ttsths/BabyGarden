package model

import (
	"time"

	"gorm.io/gorm"
)

type PushDevice struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Platform    string    `gorm:"type:varchar(20);not null" json:"platform"`
	DeviceToken string    `gorm:"type:varchar(255);not null" json:"device_token"`
	Alias       string    `gorm:"type:varchar(100)" json:"alias,omitempty"`
	Tags        JSON      `json:"tags,omitempty"`
	IsActive    int8      `gorm:"type:tinyint;default:1" json:"is_active"`
	LastUsedAt  time.Time `json:"last_used_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *PushDevice) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = NewID()
	}
	return nil
}

func (PushDevice) TableName() string { return "push_devices" }
