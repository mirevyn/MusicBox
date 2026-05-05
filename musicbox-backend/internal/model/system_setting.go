package model

import (
	"time"
)

type SystemSetting struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SettingKey   string    `gorm:"uniqueIndex;not null" json:"settingKey"`
	SettingValue string    `gorm:"type:text;not null" json:"settingValue"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName 指定数据库表名为 system_settings
func (SystemSetting) TableName() string {
	return "system_settings"
}
