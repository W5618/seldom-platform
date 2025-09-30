package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 日志记录器
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var logger *Logger

// InitLogger 初始化日志记录器
func InitLogger() error {
	// 创建logs目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 创建日志文件
	today := time.Now().Format("2006-01-02")
	infoFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("info_%s.log", today)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("error_%s.log", today)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	debugFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("debug_%s.log", today)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger = &Logger{
		infoLogger:  log.New(infoFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(errorFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(debugFile, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return nil
}

// GetLogger 获取日志记录器实例
func GetLogger() *Logger {
	return logger
}

// LogInfo 记录信息日志
func LogInfo(message string, args ...interface{}) {
	if logger != nil {
		logger.infoLogger.Printf(message, args...)
	}
	// 同时输出到控制台
	log.Printf("[INFO] "+message, args...)
}

// LogError 记录错误日志
func LogError(message string, args ...interface{}) {
	if logger != nil {
		logger.errorLogger.Printf(message, args...)
	}
	// 同时输出到控制台
	log.Printf("[ERROR] "+message, args...)
}

// LogDebug 记录调试日志
func LogDebug(message string, args ...interface{}) {
	if logger != nil {
		logger.debugLogger.Printf(message, args...)
	}
	// 调试模式下输出到控制台
	if os.Getenv("DEBUG") == "true" {
		log.Printf("[DEBUG] "+message, args...)
	}
}

// LogRequest 记录HTTP请求日志
func LogRequest(method, path, ip string, statusCode int, duration time.Duration) {
	message := fmt.Sprintf("%s %s from %s - Status: %d, Duration: %v", method, path, ip, statusCode, duration)
	LogInfo(message)
}

// LogDatabaseOperation 记录数据库操作日志
func LogDatabaseOperation(operation, table string, affected int64, duration time.Duration) {
	message := fmt.Sprintf("DB %s on %s - Affected: %d, Duration: %v", operation, table, affected, duration)
	LogDebug(message)
}

// LogAuth 记录认证相关日志
func LogAuth(username, action, ip string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	message := fmt.Sprintf("Auth %s for user %s from %s - %s", action, username, ip, status)
	LogInfo(message)
}

// Logger结构体的方法
func (l *Logger) LogInfo(category, message string, data map[string]interface{}) {
	if l != nil && l.infoLogger != nil {
		logMsg := fmt.Sprintf("[%s] %s", category, message)
		if data != nil {
			logMsg += fmt.Sprintf(" - Data: %+v", data)
		}
		l.infoLogger.Println(logMsg)
	}
}

func (l *Logger) LogError(category, message string, data map[string]interface{}) {
	if l != nil && l.errorLogger != nil {
		logMsg := fmt.Sprintf("[%s] %s", category, message)
		if data != nil {
			logMsg += fmt.Sprintf(" - Data: %+v", data)
		}
		l.errorLogger.Println(logMsg)
	}
}

func (l *Logger) LogDebug(category, message string, data map[string]interface{}) {
	if l != nil && l.debugLogger != nil {
		logMsg := fmt.Sprintf("[%s] %s", category, message)
		if data != nil {
			logMsg += fmt.Sprintf(" - Data: %+v", data)
		}
		l.debugLogger.Println(logMsg)
	}
}