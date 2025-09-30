package utils

import (
	"fmt"
	"time"
)

const (
	// 常用时间格式
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	RFC3339Format  = time.RFC3339
)

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(format, timeStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串
func GetCurrentTimeString(format string) string {
	return FormatTime(GetCurrentTime(), format)
}

// GetBeginningOfDay 获取一天的开始时间
func GetBeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取一天的结束时间
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetBeginningOfWeek 获取一周的开始时间（周一）
func GetBeginningOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日调整为7
	}
	return GetBeginningOfDay(t.AddDate(0, 0, -(weekday-1)))
}

// GetEndOfWeek 获取一周的结束时间（周日）
func GetEndOfWeek(t time.Time) time.Time {
	return GetEndOfDay(GetBeginningOfWeek(t).AddDate(0, 0, 6))
}

// GetBeginningOfMonth 获取一个月的开始时间
func GetBeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth 获取一个月的结束时间
func GetEndOfMonth(t time.Time) time.Time {
	return GetEndOfDay(GetBeginningOfMonth(t).AddDate(0, 1, -1))
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 添加小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟数
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// DiffDays 计算两个时间相差的天数
func DiffDays(t1, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	return int(t2.Sub(t1).Hours() / 24)
}

// DiffHours 计算两个时间相差的小时数
func DiffHours(t1, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	return int(t2.Sub(t1).Hours())
}

// DiffMinutes 计算两个时间相差的分钟数
func DiffMinutes(t1, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	return int(t2.Sub(t1).Minutes())
}

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 判断是否是昨天
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// IsThisWeek 判断是否是本周
func IsThisWeek(t time.Time) bool {
	now := time.Now()
	beginningOfWeek := GetBeginningOfWeek(now)
	endOfWeek := GetEndOfWeek(now)
	return t.After(beginningOfWeek) && t.Before(endOfWeek) || t.Equal(beginningOfWeek) || t.Equal(endOfWeek)
}

// IsThisMonth 判断是否是本月
func IsThisMonth(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month()
}

// FormatDuration 格式化持续时间
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1f小时", d.Hours())
	} else {
		return fmt.Sprintf("%.1f天", d.Hours()/24)
	}
}

// GetTimestamp 获取时间戳（秒）
func GetTimestamp(t time.Time) int64 {
	return t.Unix()
}

// GetTimestampMilli 获取时间戳（毫秒）
func GetTimestampMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// FromTimestamp 从时间戳创建时间
func FromTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// FromTimestampMilli 从毫秒时间戳创建时间
func FromTimestampMilli(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}