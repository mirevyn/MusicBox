package model

import "time"

// PlayHistory 记录用户播放历史。
type PlayHistory struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"column:user_id;not null" json:"userId"`
	SongID   uint      `gorm:"column:song_id;not null" json:"songId"`
	PlayedAt time.Time `gorm:"column:played_at;autoCreateTime" json:"playedAt"`
	Duration int       `gorm:"column:duration;default:0" json:"duration"`
	Source   string    `gorm:"column:source;type:varchar(100);default:''" json:"source"`
}

func (PlayHistory) TableName() string {
	return "play_histories"
}
