package utils

import (
	"errors"
	"musicbox-backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type MyClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID uint, username, role string) (string, time.Time, error) {
	// 从配置中获取 JWT 密钥和过期时间
	jwtKey := []byte(config.Conf.JWT.Secret)
	expiresDuration := time.Duration(config.Conf.JWT.ExpiresTime) * time.Hour

	expirationTime := time.Now().Add(expiresDuration)

	claims := MyClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 使用配置的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "musicbox-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (*MyClaims, error) {
	jwtKey := []byte(config.Conf.JWT.Secret)
	
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

