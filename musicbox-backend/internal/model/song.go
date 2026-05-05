package model

import "time"

// Song 对应数据库中的歌曲记录。
type Song struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title    string    `gorm:"type:varchar(200);not null" json:"title"`
	Artist   string    `gorm:"type:varchar(100);not null" json:"artist"`
	Album    string    `gorm:"type:varchar(100)" json:"album"`
	Duration int       `json:"duration"`
	FileUrl  string    `gorm:"type:varchar(300);not null;column:file_url" json:"fileUrl"`
	CoverUrl string    `gorm:"type:varchar(300);column:cover_url" json:"coverUrl"`
	LyricUrl string    `gorm:"type:varchar(300);column:lyric_url" json:"lyricUrl"`
	UploadAt time.Time `gorm:"column:upload_at;default:CURRENT_TIMESTAMP" json:"uploadAt"`
}

// TableName 指定数据库表名为 songs
func (Song) TableName() string {
	return "songs"
}
