package model

import (
	"time"

	"gorm.io/gorm"
)

type Baby struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	FamilyID    string    `gorm:"type:varchar(36);not null;index" json:"family_id"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Birthday    time.Time `gorm:"type:date;not null" json:"birthday"`
	Gender      int8      `gorm:"type:tinyint;not null" json:"gender"`
	BirthWeight *float64  `gorm:"type:decimal(5,2)" json:"birth_weight,omitempty"`
	BirthHeight *float64  `gorm:"type:decimal(4,1)" json:"birth_height,omitempty"`
	AvatarURL   string    `gorm:"type:varchar(500)" json:"avatar_url"`
	Note        string    `gorm:"type:text" json:"note"`
	IsPremature int8      `gorm:"type:tinyint;default:0" json:"is_premature"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Baby) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = NewID()
	}
	return nil
}

func (Baby) TableName() string { return "babies" }
func (b *Baby) AgeInMonths() int {
	now := time.Now()
	months := (now.Year()-b.Birthday.Year())*12 + int(now.Month()-b.Birthday.Month())
	if now.Day() < b.Birthday.Day() {
		months--
	}
	if months < 0 {
		return 0
	}
	return months
}
