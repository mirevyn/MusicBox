package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
	"musicbox-backend/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DashboardStats struct {
	TotalSongs         int64   `json:"totalSongs"`
	TotalUsers         int64   `json:"totalUsers"`
	TodayPlays         int64   `json:"todayPlays"`
	NewSongs           int64   `json:"newSongs"`
	NewUsers           int64   `json:"newUsers"`
	SongTrendPercent   float64 `json:"songTrendPercent"`
	UserTrendPercent   float64 `json:"userTrendPercent"`
	PlayTrendPercent   float64 `json:"playTrendPercent"`
	PlayTrendDirection string  `json:"playTrendDirection"`
}

type DashboardAnalytics struct {
	TrendSeries            DashboardTrendSeries `json:"trendSeries"`
	TopArtists             []DashboardValueItem `json:"topArtists"`
	PlaySourceDistribution []DashboardValueItem `json:"playSourceDistribution"`
	PlaylistStatus         []DashboardValueItem `json:"playlistStatus"`
}

type DashboardTrendSeries struct {
	Labels []string `json:"labels"`
	Plays  []int64  `json:"plays"`
	Songs  []int64  `json:"songs"`
	Users  []int64  `json:"users"`
}

type DashboardValueItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// GetUsers 根据过滤条件检索用户分页列表。
func GetUsers(keyword, role string, pageIndex, pageSize int) ([]model.User, int64, error) {
	var users []model.User

	// 创建查询构建器
	query := global.DB.Model(&model.User{})

	// 应用过滤器
	if keyword != "" {
		// 在用户名中搜索
		query = query.Where("username LIKE ?", "%"+keyword+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query = query.Order("created_at DESC")

	// 使用通用的分页函数
	total, err := utils.Paginate(query, pageIndex, pageSize, &users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetDashboardStats 获取管理端仪表盘概览数据。
func GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	if err := global.DB.Model(&model.Song{}).Count(&stats.TotalSongs).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrowStart := todayStart.Add(24 * time.Hour)
	weekStart := todayStart.AddDate(0, 0, -6)
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	prevWeekEnd := weekStart

	if err := global.DB.Model(&model.PlayHistory{}).
		Where("played_at >= ? AND played_at < ?", todayStart, tomorrowStart).
		Count(&stats.TodayPlays).Error; err != nil {
		return nil, err
	}

	var currentWeekSongs, previousWeekSongs int64
	if err := global.DB.Model(&model.Song{}).
		Where("upload_at >= ? AND upload_at < ?", weekStart, tomorrowStart).
		Count(&currentWeekSongs).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Model(&model.Song{}).
		Where("upload_at >= ? AND upload_at < ?", prevWeekStart, prevWeekEnd).
		Count(&previousWeekSongs).Error; err != nil {
		return nil, err
	}

	var currentWeekUsers, previousWeekUsers int64
	if err := global.DB.Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", weekStart, tomorrowStart).
		Count(&currentWeekUsers).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", prevWeekStart, prevWeekEnd).
		Count(&previousWeekUsers).Error; err != nil {
		return nil, err
	}

	yesterdayStart := todayStart.Add(-24 * time.Hour)
	var yesterdayPlays int64
	if err := global.DB.Model(&model.PlayHistory{}).
		Where("played_at >= ? AND played_at < ?", yesterdayStart, todayStart).
		Count(&yesterdayPlays).Error; err != nil {
		return nil, err
	}

	stats.NewSongs = currentWeekSongs
	stats.NewUsers = currentWeekUsers
	stats.SongTrendPercent = calculateTrendPercent(currentWeekSongs, previousWeekSongs)
	stats.UserTrendPercent = calculateTrendPercent(currentWeekUsers, previousWeekUsers)
	stats.PlayTrendPercent = calculateTrendPercent(stats.TodayPlays, yesterdayPlays)
	stats.PlayTrendDirection = trendDirection(stats.TodayPlays, yesterdayPlays)

	return stats, nil
}

// GetDashboardAnalytics 获取仪表盘图表数据。
func GetDashboardAnalytics() (*DashboardAnalytics, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := todayStart.AddDate(0, 0, -6)
	end := todayStart.Add(24 * time.Hour)

	playMap, err := getDailyCountMap("play_histories", "played_at", start, end)
	if err != nil {
		return nil, err
	}
	songMap, err := getDailyCountMap("songs", "upload_at", start, end)
	if err != nil {
		return nil, err
	}
	userMap, err := getDailyCountMap("users", "created_at", start, end)
	if err != nil {
		return nil, err
	}

	labels := make([]string, 0, 7)
	plays := make([]int64, 0, 7)
	songs := make([]int64, 0, 7)
	users := make([]int64, 0, 7)
	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i)
		key := day.Format("2006-01-02")
		labels = append(labels, day.Format("01-02"))
		plays = append(plays, playMap[key])
		songs = append(songs, songMap[key])
		users = append(users, userMap[key])
	}

	topArtists, err := getTopArtistsDistribution()
	if err != nil {
		return nil, err
	}

	// 播放来源直接反映首页、歌单、推荐等入口的使用情况。
	playSourceDistribution, err := getPlaySourceDistribution()
	if err != nil {
		return nil, err
	}

	playlistStatus, err := getPlaylistStatusDistribution()
	if err != nil {
		return nil, err
	}

	return &DashboardAnalytics{
		TrendSeries: DashboardTrendSeries{
			Labels: labels,
			Plays:  plays,
			Songs:  songs,
			Users:  users,
		},
		TopArtists:             topArtists,
		PlaySourceDistribution: playSourceDistribution,
		PlaylistStatus:         playlistStatus,
	}, nil
}

// BuildDashboardReportCSV 生成管理端仪表盘当前快照的 CSV 报告，供后台直接下载。
func BuildDashboardReportCSV() ([]byte, error) {
	stats, err := GetDashboardStats()
	if err != nil {
		return nil, err
	}

	analytics, err := GetDashboardAnalytics()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	// Excel 在 Windows 上更稳定地识别 UTF-8 编码。
	buffer.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(&buffer)

	rows := [][]string{
		{"MusicBox 管理后台数据报告"},
		{"导出时间", time.Now().Format("2006-01-02 15:04:05")},
		{},
		{"概览指标"},
		{"指标", "数值", "趋势"},
		{"总歌曲数", fmt.Sprintf("%d", stats.TotalSongs), formatTrendValue(stats.SongTrendPercent, "近7天 vs 上一个7天")},
		{"注册用户数", fmt.Sprintf("%d", stats.TotalUsers), formatTrendValue(stats.UserTrendPercent, "近7天 vs 上一个7天")},
		{"今日播放", fmt.Sprintf("%d", stats.TodayPlays), formatTrendValue(stats.PlayTrendPercent, "今日 vs 昨日")},
		{},
		{"近7天趋势"},
		{"日期", "播放次数", "新增歌曲", "新增用户"},
	}

	for i, label := range analytics.TrendSeries.Labels {
		rows = append(rows, []string{
			label,
			fmt.Sprintf("%d", valueAt(analytics.TrendSeries.Plays, i)),
			fmt.Sprintf("%d", valueAt(analytics.TrendSeries.Songs, i)),
			fmt.Sprintf("%d", valueAt(analytics.TrendSeries.Users, i)),
		})
	}

	rows = append(rows, []string{}, []string{"播放来源分布"}, []string{"来源", "播放次数"})
	for _, item := range analytics.PlaySourceDistribution {
		rows = append(rows, []string{item.Name, fmt.Sprintf("%d", item.Value)})
	}

	rows = append(rows, []string{}, []string{"曲库歌手分布"}, []string{"歌手", "歌曲数"})
	for _, item := range analytics.TopArtists {
		rows = append(rows, []string{item.Name, fmt.Sprintf("%d", item.Value)})
	}

	rows = append(rows, []string{}, []string{"歌单审核状态"}, []string{"状态", "数量"})
	for _, item := range analytics.PlaylistStatus {
		rows = append(rows, []string{item.Name, fmt.Sprintf("%d", item.Value)})
	}

	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}

	return buffer.Bytes(), writer.Error()
}

func calculateTrendPercent(current, previous int64) float64 {
	if previous <= 0 {
		if current <= 0 {
			return 0
		}
		return 100
	}
	return (float64(current-previous) / float64(previous)) * 100
}

func formatTrendValue(value float64, suffix string) string {
	if value == 0 {
		return "0% " + suffix
	}
	return fmt.Sprintf("%.1f%% %s", value, suffix)
}

func valueAt(values []int64, index int) int64 {
	if index < 0 || index >= len(values) {
		return 0
	}
	return values[index]
}

func trendDirection(current, previous int64) string {
	switch {
	case current > previous:
		return "up"
	case current < previous:
		return "down"
	default:
		return "flat"
	}
}

func getDailyCountMap(tableName, timeColumn string, start, end time.Time) (map[string]int64, error) {
	type result struct {
		Day   string `gorm:"column:day"`
		Total int64  `gorm:"column:total"`
	}

	var rows []result
	query := "SELECT DATE_FORMAT(" + timeColumn + ", '%Y-%m-%d') AS day, COUNT(*) AS total FROM " + tableName + " WHERE " + timeColumn + " >= ? AND " + timeColumn + " < ? GROUP BY DATE_FORMAT(" + timeColumn + ", '%Y-%m-%d')"
	if err := global.DB.Raw(query, start, end).Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make(map[string]int64, len(rows))
	for _, row := range rows {
		data[row.Day] = row.Total
	}
	return data, nil
}

func getTopArtistsDistribution() ([]DashboardValueItem, error) {
	type result struct {
		Name  string `gorm:"column:name"`
		Value int64  `gorm:"column:value"`
	}

	var rows []result
	query := `
		SELECT COALESCE(NULLIF(TRIM(artist), ''), '未知歌手') AS name, COUNT(*) AS value
		FROM songs
		GROUP BY COALESCE(NULLIF(TRIM(artist), ''), '未知歌手')
		ORDER BY value DESC
		LIMIT 6
	`
	if err := global.DB.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]DashboardValueItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, DashboardValueItem{Name: row.Name, Value: row.Value})
	}
	return items, nil
}

func getPlaySourceDistribution() ([]DashboardValueItem, error) {
	type result struct {
		Name  string `gorm:"column:name"`
		Value int64  `gorm:"column:value"`
	}

	var rows []result
	query := `
		SELECT COALESCE(NULLIF(TRIM(source), ''), '其他来源') AS name, COUNT(*) AS value
		FROM play_histories
		GROUP BY COALESCE(NULLIF(TRIM(source), ''), '其他来源')
		ORDER BY value DESC
		LIMIT 5
	`
	if err := global.DB.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]DashboardValueItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, DashboardValueItem{Name: row.Name, Value: row.Value})
	}
	return items, nil
}

func getPlaylistStatusDistribution() ([]DashboardValueItem, error) {
	type result struct {
		Status int   `gorm:"column:status"`
		Value  int64 `gorm:"column:value"`
	}

	var rows []result
	if err := global.DB.Table("playlists").
		Select("status, COUNT(*) AS value").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	labelMap := map[int]string{
		0: "待审核",
		1: "已通过",
		2: "已驳回",
	}

	items := []DashboardValueItem{
		{Name: "待审核", Value: 0},
		{Name: "已通过", Value: 0},
		{Name: "已驳回", Value: 0},
	}
	for _, row := range rows {
		name, ok := labelMap[row.Status]
		if !ok {
			continue
		}
		for i := range items {
			if items[i].Name == name {
				items[i].Value = row.Value
			}
		}
	}

	return items, nil
}

// UpdateUserByAdmin 管理员更新用户信息
func UpdateUserByAdmin(actorUserID, targetUserID uint, updates map[string]interface{}) (*model.User, error) {
	var user model.User
	// 查找用户
	if err := global.DB.First(&user, targetUserID).Error; err != nil {
		return nil, ErrUserNotFound
	}
	if err := validateAdminUserMutation(actorUserID, user); err != nil {
		return nil, err
	}

	// 如果更新了密码，对新密码进行哈希处理
	if newPassword, ok := updates["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("密码哈希失败")
		}
		updates["password_hash"] = string(hashedPassword)
		delete(updates, "password") // 从 map 中删除原始密码
	}

	// 用户状态更新 (前端传递 int)
	if newStatus, ok := updates["status"].(int); ok {
		updates["status"] = newStatus
	}

	// 执行更新
	if err := global.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, errors.New("更新用户信息失败")
	}

	go GetNotificationService().PushDashboardRefresh()

	return &user, nil
}

// DeleteUser 管理员删除用户及其关联资源
func DeleteUser(actorUserID, targetUserID uint) error {
	var user model.User
	// 查找用户
	if err := global.DB.First(&user, targetUserID).Error; err != nil {
		return ErrUserNotFound
	}
	if err := validateAdminUserMutation(actorUserID, user); err != nil {
		return err
	}

	// 从数据库删除用户记录
	if err := global.DB.Delete(&user).Error; err != nil {
		return errors.New("删除用户失败")
	}
	// 删除用户头像文件 (如果存在且不是默认头像)
	if user.AvatarURL != "" && user.AvatarURL != "avatar/default.jpg" {
		_ = utils.DeleteFile(user.AvatarURL)
	}

	go GetNotificationService().PushDashboardRefresh()

	return nil
}

func validateAdminUserMutation(actorUserID uint, target model.User) error {
	if isDefaultAdminUser(target) {
		return ErrProtectedAdmin
	}
	if actorUserID != 0 && actorUserID == target.ID {
		return ErrSelfOperation
	}
	return nil
}

func isDefaultAdminUser(user model.User) bool {
	return user.ID == 1 || strings.EqualFold(strings.TrimSpace(user.Username), "admin")
}

// GetAdminPlaylists 管理员获取歌单列表
func GetAdminPlaylists(keyword string, status int, pageIndex, pageSize int) ([]model.Playlist, int64, error) {
	var playlists []model.Playlist
	query := global.DB.Model(&model.Playlist{}).Preload("User")

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// status: -1 表示全部
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	total, err := utils.Paginate(query, pageIndex, pageSize, &playlists)
	if err != nil {
		return nil, 0, err
	}

	// 填充每个歌单的 SongCount
	for i := range playlists {
		var count int64
		global.DB.Model(&model.PlaylistSong{}).Where("playlist_id = ?", playlists[i].ID).Count(&count)
		playlists[i].SongCount = int(count)
	}

	return playlists, total, nil
}

// UpdatePlaylistStatus 管理员更新歌单状态 (审核)
func UpdatePlaylistStatus(id uint, status int, rejectReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == 2 { // 驳回
		updates["reject_reason"] = rejectReason
	} else {
		updates["reject_reason"] = "" // 通过则清空驳回理由
	}

	result := global.DB.Model(&model.Playlist{}).Where("id = ?", id).Updates(updates)
	if result.Error == nil {
		go GetNotificationService().PushDashboardRefresh()
	}
	return result.Error
}

// DeletePlaylistByAdmin 管理员删除歌单
func DeletePlaylistByAdmin(id uint) error {
	// 查找歌单 (确认存在)
	var playlist model.Playlist
	if err := global.DB.First(&playlist, id).Error; err != nil {
		return err
	}

	// 数据库删除 (级联删除 playlist_songs 由数据库外键处理，或者 Gorm 钩子)
	if err := global.DB.Delete(&playlist).Error; err != nil {
		return err
	}
	// 删除封面 (如果是本地上传的)
	if playlist.CoverUrl != "" {
		_ = utils.DeleteFile(playlist.CoverUrl)
	}

	go GetNotificationService().PushDashboardRefresh()

	return nil
}

// GetAdminPlaylistDetails 管理员获取歌单详情 (包含歌曲列表)
func GetAdminPlaylistDetails(id uint) (*model.Playlist, []model.Song, error) {
	var playlist model.Playlist

	// 获取歌单基本信息
	if err := global.DB.Preload("User").First(&playlist, id).Error; err != nil {
		return nil, nil, err
	}

	// 获取歌单歌曲列表 (按排序和添加时间)
	var songs []model.Song
	err := global.DB.
		Table("songs").
		Joins("JOIN playlist_songs ps ON songs.id = ps.song_id").
		Where("ps.playlist_id = ?", id).
		Order("ps.sort_order ASC, ps.added_at DESC").
		Find(&songs).Error

	if err != nil {
		return nil, nil, err
	}

	playlist.SongCount = len(songs)

	return &playlist, songs, nil
}
