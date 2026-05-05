package service

import (
	"errors"
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
	"musicbox-backend/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserDisabled = errors.New("user disabled")

// Register 用户注册业务逻辑
func Register(user *model.User) error {
	// 检查用户名是否已存在
	var count int64
	if err := global.DB.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err // 返回真正的数据库错误
	}
	if count > 0 {
		return ErrUsernameExists
	}

	// 使用 bcrypt 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	result := global.DB.Create(user)
	if result.Error == nil {
		go GetNotificationService().PushDashboardRefresh()
	}
	return result.Error
}

// Login 登录业务逻辑
// 返回：用户信息, Token字符串, 过期时间, 错误信息
func Login(username, password string) (*model.User, string, time.Time, error) {
	var user model.User

	// 根据用户名查询用户
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", time.Time{}, ErrUserNotFound
		}
		return nil, "", time.Time{}, err
	}

	// 校验密码 (使用 bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", time.Time{}, ErrPasswordIncorrect
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, "", time.Time{}, errors.New("账号已被禁用")
	}

	// 生成 JWT Token
	token, expiration, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", time.Time{}, errors.New("Token生成失败")
	}

	return &user, token, expiration, nil
}

// ChangePassword 修改用户密码
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User

	// 查找用户
	if err := global.DB.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrPasswordIncorrect
	}

	// 哈希新密码
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("新密码哈希失败")
	}

	// 更新数据库中的密码
	if err := global.DB.Model(&user).Update("password_hash", string(hashedNewPassword)).Error; err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}

// UpdateAvatar 更新用户头像
func UpdateAvatar(userID uint, newAvatarURL string) (string, error) {
	var user model.User

	// 查找用户以获取旧头像路径
	if err := global.DB.First(&user, userID).Error; err != nil {
		return "", ErrUserNotFound
	}
	oldAvatarURL := user.AvatarURL

	// 更新数据库中的头像URL
	if err := global.DB.Model(&user).Update("avatar_url", newAvatarURL).Error; err != nil {
		return "", errors.New("更新头像失败")
	}

	// 返回旧头像的路径，以便在 handler 中删除
	return oldAvatarURL, nil
}

// 根据 ID 获取用户信息
func GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetActiveUserForClaims 通过 JWT 中的用户 ID 复查数据库里的最新账号状态。
func GetActiveUserForClaims(claims *utils.MyClaims) (*model.User, error) {
	if claims == nil {
		return nil, ErrNotFound
	}
	user, err := GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status == 0 {
		return nil, ErrUserDisabled
	}
	return user, nil
}
