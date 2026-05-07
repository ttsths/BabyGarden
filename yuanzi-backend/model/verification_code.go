package model

import (
	"time"

	"gorm.io/gorm"
)

type VerificationCode struct {
	ID        string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Phone     string     `gorm:"type:varchar(11);not null" json:"phone"`
	Code      string     `gorm:"type:varchar(6);not null" json:"code"`
	Type      string     `gorm:"type:varchar(20);default:'login'" json:"type"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	IPAddress string     `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time  `json:"created_at"`
}

func (vc *VerificationCode) BeforeCreate(_ *gorm.DB) error {
	if vc.ID == "" {
		vc.ID = NewID()
	}
	return nil
}
func (VerificationCode) TableName() string   { return "verification_codes" }
func (vc *VerificationCode) IsExpired() bool { return time.Now().After(vc.ExpiresAt) }
func (vc *VerificationCode) IsUsed() bool    { return vc.UsedAt != nil }
