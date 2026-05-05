package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrInvalidUpload = errors.New("invalid upload")

type invalidUploadError string

func (e invalidUploadError) Error() string {
	return string(e)
}

func (e invalidUploadError) Is(target error) bool {
	return target == ErrInvalidUpload
}

type uploadPolicy struct {
	name       string
	maxBytes   int64
	extensions map[string]struct{}
	acceptMIME func(string) bool
}

var (
	imageUploadPolicy = uploadPolicy{
		name:     "图片",
		maxBytes: 5 * 1024 * 1024,
		extensions: map[string]struct{}{
			".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {}, ".gif": {},
		},
		acceptMIME: func(contentType string) bool {
			return strings.HasPrefix(contentType, "image/")
		},
	}
	audioUploadPolicy = uploadPolicy{
		name:     "音频",
		maxBytes: 200 * 1024 * 1024,
		extensions: map[string]struct{}{
			".mp3": {}, ".flac": {}, ".wav": {}, ".m4a": {}, ".aac": {}, ".ogg": {},
		},
		acceptMIME: func(contentType string) bool {
			return strings.HasPrefix(contentType, "audio/") ||
				contentType == "video/mp4" ||
				contentType == "application/ogg" ||
				contentType == "application/octet-stream"
		},
	}
	lyricUploadPolicy = uploadPolicy{
		name:     "歌词",
		maxBytes: 1 * 1024 * 1024,
		extensions: map[string]struct{}{
			".lrc": {}, ".txt": {},
		},
		acceptMIME: func(contentType string) bool {
			return strings.HasPrefix(contentType, "text/") ||
				contentType == "application/octet-stream"
		},
	}
)

// SaveUploadedFile 将上传的文件以唯一名称保存到目标目录。
// 它返回文件保存的路径或错误。
func SaveUploadedFile(c *gin.Context, header *multipart.FileHeader, dstDir string) (string, error) {
	policy := uploadPolicyForDir(dstDir)
	if err := validateUploadedFile(header, policy); err != nil {
		return "", err
	}

	// 确保目标目录存在
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", fmt.Errorf("未能创建目录: %w", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	uniqueFilename := uuid.New().String() + ext

	// 创建目标文件路径
	dstPath := filepath.Join(dstDir, uniqueFilename)

	// 打开上传的文件
	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("无法打开上传的文件: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("复制文件内容失败: %w", err)
	}

	return publicUploadPath(dstPath), nil
}

func publicUploadPath(filePath string) string {
	return utils.NormalizeUploadPath(filePath)
}

func uploadPolicyForDir(dstDir string) uploadPolicy {
	publicDir := publicUploadPath(dstDir)
	switch {
	case strings.Contains(publicDir, "uploads/music"):
		return audioUploadPolicy
	case strings.Contains(publicDir, "uploads/lyrics"):
		return lyricUploadPolicy
	default:
		return imageUploadPolicy
	}
}

func validateUploadedFile(header *multipart.FileHeader, policy uploadPolicy) error {
	if header == nil {
		return invalidUploadError("未提供上传文件")
	}
	if header.Size <= 0 {
		return invalidUploadError(fmt.Sprintf("%s文件不能为空", policy.name))
	}
	if header.Size > policy.maxBytes {
		return invalidUploadError(fmt.Sprintf("%s文件不能超过 %s", policy.name, formatMaxUploadSize(policy.maxBytes)))
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if _, ok := policy.extensions[ext]; !ok {
		return invalidUploadError(fmt.Sprintf("不支持的%s文件类型: %s", policy.name, ext))
	}

	src, err := header.Open()
	if err != nil {
		return fmt.Errorf("无法打开上传的文件: %w", err)
	}
	defer src.Close()

	sample := make([]byte, 512)
	n, readErr := src.Read(sample)
	if readErr != nil && readErr != io.EOF {
		return fmt.Errorf("无法读取上传文件内容: %w", readErr)
	}
	if seeker, ok := src.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("无法重置上传文件读取位置: %w", err)
		}
	}

	detected := http.DetectContentType(sample[:n])
	declared := strings.ToLower(strings.TrimSpace(header.Header.Get("Content-Type")))
	if policy.acceptMIME(detected) {
		return nil
	}
	if detected == "application/octet-stream" && declared != "" && policy.acceptMIME(declared) {
		return nil
	}

	return invalidUploadError(fmt.Sprintf("%s文件内容类型不合法: %s", policy.name, detected))
}

func formatMaxUploadSize(size int64) string {
	const mb = 1024 * 1024
	if size%mb == 0 {
		return fmt.Sprintf("%dMB", size/mb)
	}
	return fmt.Sprintf("%d bytes", size)
}
