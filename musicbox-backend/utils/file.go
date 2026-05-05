package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ZipDir 将指定目录 src 压缩为 dst 文件
func ZipDir(src, dst string) error {
	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径作为 zip 内部的文件名
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// NormalizeUploadPath 将文件系统路径或 URL 风格路径归一化为 uploads/...。
func NormalizeUploadPath(filePath string) string {
	normalized := strings.TrimSpace(filePath)
	if normalized == "" {
		return ""
	}

	normalized = strings.ReplaceAll(normalized, "\\", "/")
	normalized = path.Clean(normalized)
	normalized = strings.TrimPrefix(normalized, "./")
	normalized = strings.TrimPrefix(normalized, "/")

	parts := strings.Split(normalized, "/")
	for i, part := range parts {
		if part == "uploads" {
			return strings.Join(parts[i:], "/")
		}
	}

	return normalized
}

// ResolveUploadFilePath 将 uploads/... 解析为当前工作目录下的真实文件路径。
// 如果路径不在 uploads 根目录内，返回错误，避免误删任意文件。
func ResolveUploadFilePath(filePath string) (string, error) {
	publicPath := NormalizeUploadPath(filePath)
	if publicPath == "" {
		return "", nil
	}
	if publicPath == "uploads" || !strings.HasPrefix(publicPath, "uploads/") {
		return "", fmt.Errorf("非法上传资源路径: %s", filePath)
	}

	root, err := filepath.Abs("uploads")
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(filepath.FromSlash(publicPath))
	if err != nil {
		return "", err
	}

	rootWithSep := root + string(os.PathSeparator)
	if target != root && !strings.HasPrefix(target, rootWithSep) {
		return "", fmt.Errorf("上传资源路径越界: %s", filePath)
	}

	return target, nil
}

// DeleteFile 尝试删除 uploads 下的文件。
// 如果文件不存在，它不会返回错误；如果路径越界，它会返回错误。
func DeleteFile(filePath string) error {
	resolvedPath, err := ResolveUploadFilePath(filePath)
	if err != nil {
		return err
	}
	if resolvedPath == "" {
		return nil
	}

	info, err := os.Stat(resolvedPath)
	if os.IsNotExist(err) {
		// 文件不存在，无需操作
		return nil
	}
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("拒绝删除目录: %s", filePath)
	}

	if err := os.Remove(resolvedPath); err != nil {
		return fmt.Errorf("无法删除文件 %s: %w", filePath, err)
	}
	return nil
}

// CopyFile 将文件从 src 复制到 dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dirPath string) error {
	dirPath = filepath.Clean(dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}
	return nil
}
