package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型，对齐当前开发库 `users` 表结构。
type User struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Phone       string     `gorm:"type:varchar(32);uniqueIndex;not null" json:"phone"`
	Nickname    string     `gorm:"type:varchar(50)" json:"nickname"`
	AvatarURL   string     `gorm:"type:varchar(500)" json:"avatar_url"`
	Status      int8       `gorm:"type:tinyint;default:1" json:"status"`
	IsAdmin     int8       `gorm:"type:tinyint;default:0" json:"is_admin"`
	Password    string     `gorm:"type:varchar(255)" json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"type:varchar(45)" json:"last_login_ip"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = NewID()
	}
	return nil
}

func (User) TableName() string { return "users" }

type UserInfo struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Status    int8   `json:"status"`
}

func (u *User) ToUserInfo() UserInfo {
	return UserInfo{ID: u.ID, Phone: u.Phone, Nickname: u.Nickname, AvatarURL: u.AvatarURL, Status: u.Status}
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

// CheckPassword compares a plaintext password with the stored password.
// In production, use bcrypt.CompareHashAndPassword.
func (u *User) CheckPassword(password string) bool {
	if u.Password == "" {
		return false
	}
	return u.Password == password
}
