package model

import "time"

// User
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(256);not null" json:"-"`
	AvatarURL    string    `gorm:"column:avatar_url" json:"avatarUrl"`
	Role         string    `gorm:"default:'User'" json:"role"`
	Status       int       `gorm:"default:1;comment:'0-Disabled, 1-Active'" json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

// TableName 指定数据库表名为 users
func (User) TableName() string {
	return "users"
}
