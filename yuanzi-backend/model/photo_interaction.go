package model

import (
	"time"

	"gorm.io/gorm"
)

type PhotoComment struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	PhotoID   string         `gorm:"type:varchar(36);not null;index" json:"photo_id"`
	FamilyID  string         `gorm:"type:varchar(36);not null;index" json:"family_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (c *PhotoComment) BeforeCreate(_ *gorm.DB) error {
	if c.ID == "" {
		c.ID = NewID()
	}
	return nil
}

func (PhotoComment) TableName() string { return "photo_comments" }

type PhotoLike struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	PhotoID   string    `gorm:"type:varchar(36);not null;uniqueIndex:uk_photo_user" json:"photo_id"`
	FamilyID  string    `gorm:"type:varchar(36);not null;index" json:"family_id"`
	UserID    string    `gorm:"type:varchar(36);not null;uniqueIndex:uk_photo_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *PhotoLike) BeforeCreate(_ *gorm.DB) error {
	if l.ID == "" {
		l.ID = NewID()
	}
	return nil
}

func (PhotoLike) TableName() string { return "photo_likes" }
