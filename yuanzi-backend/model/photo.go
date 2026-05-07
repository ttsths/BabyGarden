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
