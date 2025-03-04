package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ChatLog 定义聊天日志结构
// @Description 聊天日志记录
type ChatLog struct {
	// 唯一标识符
	ID string `json:"id"`
	// 聊天文本内容
	Text string `json:"text" binding:"required" example:"这是一条聊天记录"`
	// 验证结果
	ValidationResult string `json:"validation_result" example:"不涉密"`
	// 输入时间
	InputTime string `json:"input_time" example:"2023-01-01 12:00:00.000"`
}

var (
	// 存储目录
	storageDir = "./storage"
	// 文件锁，防止并发写入问题
	fileLock sync.Mutex
	// 日志文件大小限制 (100MB)
	maxFileSize int64 = 100 * 1024 * 1024
)

// InitStorage 初始化存储系统
func InitStorage() error {
	// 确保存储目录存在
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("创建存储目录失败: %w", err)
	}
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	return nil
}

// 自定义时间格式
const (
	// 日期格式 (用于文件名)
	dateFormat = "2006-01-02"
	// 完整时间格式 (用于日志记录)
	timeFormat = "2006-01-02 15:04:05.000"
	// ID时间格式
	idTimeFormat = "20060102150405"
	// 随机字符集
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 随机部分长度
	randomPartLength = 6
)

// 生成指定长度的随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// SaveLog 保存日志到文件
func SaveLog(log *ChatLog) error {
	// 生成唯一ID (yyyyMMddHHmmss+6位随机数)
	if log.ID == "" {
		timepart := time.Now().Format(idTimeFormat)
		randompart := generateRandomString(randomPartLength)
		log.ID = timepart + randompart
	}

	// 设置输入时间
	log.InputTime = time.Now().Format(timeFormat)

	// 按日期组织文件
	currentTime := time.Now()
	dateStr := currentTime.Format(dateFormat)
	baseFilePath := filepath.Join(storageDir, dateStr)
	filePath := baseFilePath + ".json"

	// 加锁防止并发写入
	fileLock.Lock()
	defer fileLock.Unlock()

	// 检查文件是否存在及其大小
	fileIndex := 1
	if fileInfo, err := os.Stat(filePath); err == nil {
		// 如果文件大小接近或超过限制，创建新文件
		if fileInfo.Size() >= maxFileSize {
			// 查找可用的文件名
			for {
				newFilePath := fmt.Sprintf("%s_%d.json", baseFilePath, fileIndex)
				if fileInfo, err := os.Stat(newFilePath); err != nil || fileInfo.Size() < maxFileSize {
					filePath = newFilePath
					break
				}
				fileIndex++
			}
		}
	}

	// 读取现有日志
	var logs []ChatLog
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取日志文件失败: %w", err)
		}
		if err := json.Unmarshal(data, &logs); err != nil {
			// 如果解析失败，创建新的日志集合
			logs = []ChatLog{}
		}
	}

	// 添加新日志
	logs = append(logs, *log)

	// 写入文件
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化日志失败: %w", err)
	}

	// 使用追加模式写入文件
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("写入日志文件失败: %w", err)
	}

	return nil
}

// GetLogs 获取指定日期范围内的日志
func GetLogs(startDate, endDate time.Time) ([]ChatLog, error) {
	var allLogs []ChatLog

	// 如果结束日期为零值，设置为当前日期
	if endDate.IsZero() {
		endDate = time.Now()
	}

	// 如果开始日期为零值，设置为结束日期的前7天
	if startDate.IsZero() {
		startDate = endDate.AddDate(0, 0, -7)
	}

	// 遍历日期范围
	current := startDate
	for !current.After(endDate) {
		dateStr := current.Format(dateFormat)
		baseFilePath := filepath.Join(storageDir, dateStr)

		// 查找所有与该日期相关的日志文件
		files, err := filepath.Glob(baseFilePath + "*.json")
		if err != nil {
			return nil, fmt.Errorf("查找日志文件失败: %w", err)
		}

		// 处理找到的所有文件
		for _, filePath := range files {
			// 读取日志文件
			data, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("读取日志文件失败: %w", err)
			}

			// 解析日志
			var logs []ChatLog
			if err := json.Unmarshal(data, &logs); err != nil {
				return nil, fmt.Errorf("解析日志文件失败: %w", err)
			}

			// 添加到结果集
			allLogs = append(allLogs, logs...)
		}

		// 前进到下一天
		current = current.AddDate(0, 0, 1)
	}

	return allLogs, nil
}
