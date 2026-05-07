package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type RecordType string

const (
	RecordTypeFeeding RecordType = "feeding"
	RecordTypeSleep   RecordType = "sleep"
	RecordTypeDiaper  RecordType = "diaper"
	RecordTypeGrowth  RecordType = "growth"
)

type Record struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	BabyID    string         `gorm:"type:varchar(36);not null;index" json:"baby_id"`
	FamilyID  string         `gorm:"type:varchar(36);not null;index" json:"family_id"`
	Type      RecordType     `gorm:"type:varchar(20);not null;index" json:"type"`
	StartedAt time.Time      `gorm:"not null;index" json:"started_at"`
	EndedAt   sql.NullTime   `json:"ended_at,omitempty"`
	Content   JSON           `gorm:"type:json;not null" json:"content"`
	Note      string         `gorm:"type:text" json:"note"`
	CreatedBy string         `gorm:"type:varchar(36);not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Record) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = NewID()
	}
	return nil
}
func (Record) TableName() string { return "records" }
func (r *Record) Duration() int {
	if !r.EndedAt.Valid {
		return 0
	}
	return int(r.EndedAt.Time.Sub(r.StartedAt).Minutes())
}

type JSON json.RawMessage

func (j JSON) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
	case string:
		*j = JSON(v)
	default:
		*j = nil
	}
	return nil
}
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}
func (j *JSON) UnmarshalJSON(data []byte) error { *j = JSON(data); return nil }

type FeedingContent struct {
	Type     string `json:"type"`
	Side     string `json:"side,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Amount   int    `json:"amount,omitempty"`
	Unit     string `json:"unit,omitempty"`
}
type SleepContent struct {
	Quality  string `json:"quality"`
	Location string `json:"location"`
}
type DiaperContent struct {
	Type        string `json:"type"`
	Color       string `json:"color,omitempty"`
	Consistency string `json:"consistency,omitempty"`
}
type GrowthContent struct {
	Weight            float64 `json:"weight"`
	Height            float64 `json:"height"`
	HeadCircumference float64 `json:"head_circumference,omitempty"`
}
