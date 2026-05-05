package service

import (
	"errors"

	"gorm.io/gorm"
)

// 定义通用的业务错误变量
// 使用 var 方式定义，配合 errors.Is 使用
var (
	// ErrForbiddenAccess 当用户无权执行某个操作时返回此错误
	ErrForbiddenAccess = errors.New("无权进行此操作")

	// ErrDuplicateEntry 当数据重复时返回此错误 (通用)
	ErrDuplicateEntry = errors.New("数据已存在")

	// ErrUsernameExists 当注册或修改用户信息时，用户名已存在
	// 对应 admin.go 和 user.go 中的 "用户名已存在" 判断
	ErrUsernameExists = errors.New("用户名已存在")

	// ErrUserNotFound 当操作针对的用户不存在时
	// 对应 admin.go 和 song_likes.go 中的 "用户不存在" / "user not found"
	ErrUserNotFound = errors.New("用户不存在")

	// ErrProtectedAdmin 表示默认管理员账号不可通过后台用户管理修改或删除
	ErrProtectedAdmin = errors.New("默认管理员不可修改")

	// ErrSelfOperation 表示管理员不可在后台用户管理中操作自己
	ErrSelfOperation = errors.New("不能在后台用户管理中操作自己")

	// ErrSongNotFound 当操作针对的歌曲不存在时
	// 对应 song.go 和 song_likes.go 中的 "歌曲未找到" / "song not found"
	ErrSongNotFound = errors.New("歌曲不存在")

	// ErrPasswordIncorrect 登录或修改密码时旧密码错误
	ErrPasswordIncorrect = errors.New("密码错误")

	// ErrNotFound 是 gorm.ErrRecordNotFound 的别名
	// 用于通用的"记录未找到"情况
	ErrNotFound = gorm.ErrRecordNotFound
)
