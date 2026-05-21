package model

import (
	"time"

	"gorm.io/gorm"
)

type PhotoStatus string

const (
	PhotoStatusPending PhotoStatus = "pending"
	PhotoStatusActive  PhotoStatus = "active"
	PhotoStatusDeleted PhotoStatus = "deleted"
)

type Photo struct {
	ID           string      `gorm:"type:varchar(36);primaryKey" json:"id"`
	BabyID       string      `gorm:"type:varchar(36);not null;index" json:"baby_id"`
	FamilyID     string      `gorm:"type:varchar(36);not null;index" json:"family_id"`
	OSSKey       string      `gorm:"type:varchar(500);not null" json:"oss_key"`
	ThumbnailKey string      `gorm:"type:varchar(500)" json:"thumbnail_key,omitempty"`
	Width        *int        `json:"width,omitempty"`
	Height       *int        `json:"height,omitempty"`
	Size         int64       `json:"size"`
	ContentType  string      `gorm:"type:varchar(50);default:'image/jpeg'" json:"content_type"`
	TakenAt      *time.Time  `json:"taken_at,omitempty"`
	Description  string      `gorm:"type:text" json:"description"`
	UploadedBy   string      `gorm:"type:varchar(36);not null" json:"uploaded_by"`
	UploadedAt   time.Time   `json:"uploaded_at"`
	Status       PhotoStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
}

func (p *Photo) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = NewID()
	}
	return nil
}

func (Photo) TableName() string                 { return "photos" }
func (p *Photo) GetURL(cdnDomain string) string { return "https://" + cdnDomain + "/" + p.OSSKey }
func (p *Photo) GetThumbnailURL(cdnDomain string) string {
	if p.ThumbnailKey == "" {
		return p.GetURL(cdnDomain)
	}
	return "https://" + cdnDomain + "/" + p.ThumbnailKey
}

type PhotoComment struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	PhotoID   string         `gorm:"type:varchar(36);not null;index" json:"photo_id"`
	FamilyID  string         `gorm:"type:varchar(36);not null;index" json:"family_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Content   string         `gorm:"type:varchar(500);not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"-" json:"user,omitempty"`
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
	User      User      `gorm:"-" json:"user,omitempty"`
}

func (l *PhotoLike) BeforeCreate(_ *gorm.DB) error {
	if l.ID == "" {
		l.ID = NewID()
	}
	return nil
}

func (PhotoLike) TableName() string { return "photo_likes" }
