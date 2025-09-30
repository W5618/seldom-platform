package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	if email == "" {
		return true // 允许空邮箱
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidUsername 验证用户名格式
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 150 {
		return false
	}
	
	// 用户名只能包含字母、数字、下划线和连字符
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword 验证密码强度
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	
	// 至少包含一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	return hasLetter && hasNumber
}

// SanitizeString 清理字符串，移除危险字符
func SanitizeString(input string) string {
	// 移除前后空格
	input = strings.TrimSpace(input)
	
	// 移除HTML标签
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlRegex.ReplaceAllString(input, "")
	
	// 移除SQL注入相关字符
	sqlRegex := regexp.MustCompile(`[';\"\\]`)
	input = sqlRegex.ReplaceAllString(input, "")
	
	return input
}

// IsValidPort 验证端口号
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsValidCronExpression 验证Cron表达式格式
func IsValidCronExpression(cron string) bool {
	if cron == "" {
		return true // 允许空的cron表达式
	}
	
	// 支持6个字段的cron表达式（包含秒）
	parts := strings.Fields(cron)
	return len(parts) == 6 || len(parts) == 5
}

// ValidateProjectName 验证项目名称
func ValidateProjectName(name string) bool {
	if len(name) < 1 || len(name) > 100 {
		return false
	}
	
	// 项目名称不能包含特殊字符
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\u4e00-\u9fa5_-\s]+$`)
	return nameRegex.MatchString(name)
}

// ValidateURL 验证URL格式
func ValidateURL(url string) bool {
	if url == "" {
		return true // 允许空URL
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}