package service

import (
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
)

type SystemService struct{}

var systemService = &SystemService{}

func GetSystemService() *SystemService {
	return systemService
}

// GetSetting 获取设置
func (s *SystemService) GetSetting(key string) (string, error) {
	var setting model.SystemSetting
	err := global.DB.Where("setting_key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.SettingValue, nil
}

// UpdateSetting 更新设置
func (s *SystemService) UpdateSetting(key, value string) error {
	return global.DB.Model(&model.SystemSetting{}).
		Where("setting_key = ?", key).
		Update("setting_value", value).Error
}

// IsRegisterAllowed 检查是否允许注册
func (s *SystemService) IsRegisterAllowed() bool {
	val, err := s.GetSetting("allow_register")
	if err != nil {
		return true // 默认允许
	}
	return val == "true"
}
