package model

import "time"

// SongLike 表示用户点赞歌曲的记录。
type SongLike struct {
	SongID  uint      `gorm:"primaryKey;autoIncrement:false" json:"songId"` // 歌曲ID
	UserID  uint      `gorm:"primaryKey;autoIncrement:false" json:"userId"` // 用户ID
	LikedAt time.Time `gorm:"autoCreateTime" json:"likedAt"`                // 点赞时间，自动创建

	// Associations 关联模型
	User User `gorm:"foreignKey:UserID"` // 关联的用户
	Song Song `gorm:"foreignKey:SongID"` // 关联的歌曲
}

// TableName 指定数据库表名为 song_likes
func (SongLike) TableName() string {
	return "song_likes"
}
