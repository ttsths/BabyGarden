package model

import (
	"time"

	"gorm.io/gorm"
)

type Family struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	InviteCode   string    `gorm:"type:varchar(8);uniqueIndex;not null" json:"invite_code"`
	CreatedBy    string    `gorm:"type:varchar(36);not null" json:"created_by"`
	IsPaid       int8      `gorm:"type:tinyint;default:0" json:"is_paid"`
	StorageLimit int64     `gorm:"default:1073741824" json:"storage_limit"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (f *Family) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = NewID()
	}
	return nil
}

func (Family) TableName() string { return "families" }

type FamilyRole string

const (
	FamilyRoleAdmin  FamilyRole = "admin"
	FamilyRoleMember FamilyRole = "member"
	FamilyRoleElder  FamilyRole = "elder"
)

type FamilyMember struct {
	ID            string     `gorm:"type:varchar(255);primaryKey" json:"id"`
	FamilyID      string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_family_user" json:"family_id"`
	UserID        string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_family_user" json:"user_id"`
	Role          FamilyRole `gorm:"type:varchar(20);default:'member'" json:"role"`
	ElderMode     int8       `gorm:"type:tinyint;default:0" json:"elder_mode"`
	Notifications JSON       `json:"notifications,omitempty"`
	JoinedAt      time.Time  `json:"joined_at"`
	User          User       `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (fm *FamilyMember) BeforeCreate(_ *gorm.DB) error {
	if fm.ID == "" {
		fm.ID = NewID()
	}
	return nil
}

func (FamilyMember) TableName() string          { return "family_members" }
func (fm *FamilyMember) IsAdmin() bool          { return fm.Role == FamilyRoleAdmin }
func (fm *FamilyMember) CanManageMembers() bool { return fm.Role == FamilyRoleAdmin }
