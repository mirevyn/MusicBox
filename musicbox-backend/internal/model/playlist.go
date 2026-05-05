package model

import "time"

// Playlist 歌单模型
type Playlist struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;column:user_id" json:"userId"`
	Title string `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`
	CoverUrl    string    `gorm:"column:cover_url;type:varchar(300)" json:"coverUrl"`
	IsPublic    bool      `gorm:"column:is_public" json:"isPublic"`
	PlayCount   int       `gorm:"column:play_count;default:0" json:"playCount"`
	Status      int       `gorm:"column:status;default:1;comment:'0-Pending, 1-Approved, 2-Rejected'" json:"status"`
	RejectReason string   `gorm:"column:reject_reason;type:varchar(255)" json:"rejectReason"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// 关联关系
	User      *User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联创建者
	SongCount int   `gorm:"-" json:"songCount"`                      // 统计字段
}

// PlaylistSong 歌单与歌曲的关联表 (中间表)
type PlaylistSong struct {
	ID         uint      `gorm:"primaryKey"`
	PlaylistID uint      `gorm:"column:playlist_id;not null"`
	SongID     uint      `gorm:"column:song_id;not null"`
	SortOrder  int       `gorm:"column:sort_order;default:0"`
	AddedAt    time.Time `gorm:"column:added_at;autoCreateTime"`
}

// TableName 指定数据库表名为 playlists
func (Playlist) TableName() string {
	return "playlists"
}

// TableName 指定数据库表名为 playlist_songs
func (PlaylistSong) TableName() string {
	return "playlist_songs"
}
