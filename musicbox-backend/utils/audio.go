package utils

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mewkiz/flac"
	"github.com/tcolgate/mp3"
)

// GetAudioDuration 读取一个音频文件并返回其时长（秒
// 支持格式: .mp3, .flac
func GetAudioDuration(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		return getMP3Duration(file)
	case ".flac":
		return getFLACDuration(file)
	default:
		return 0, fmt.Errorf("不支持的音频格式: %s", ext)
	}
}

func getMP3Duration(r io.Reader) (int, error) {
	decoder := mp3.NewDecoder(r)
	var totalDuration time.Duration
	var frame mp3.Frame
	skipped := 0

	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			// 如果已经解码了一些帧，遇到错误（如文件末尾的元数据）则忽略错误返回当前时长
			if totalDuration > 0 {
				fmt.Printf("MP3解码警告: %v\n", err)
				break
			}
			return 0, fmt.Errorf("解码MP3帧时出错: %w", err)
		}
		totalDuration += frame.Duration()
	}
	return int(math.Round(totalDuration.Seconds())), nil
}

func getFLACDuration(r io.Reader) (int, error) {
	stream, err := flac.Parse(r)
	if err != nil {
		return 0, fmt.Errorf("解析FLAC失败: %w", err)
	}
	defer stream.Close()

	info := stream.Info
	if info.SampleRate == 0 {
		return 0, fmt.Errorf("无效的采样率")
	}

	seconds := float64(info.NSamples) / float64(info.SampleRate)
	return int(math.Round(seconds)), nil
}
