package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo 文件信息结构
type FileInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	MimeType string `json:"mime_type"`
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, uploadDir string) (*FileInfo, error) {
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), GenerateRandomFileName(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	return &FileInfo{
		Name:     filename,
		Size:     file.Size,
		Path:     filePath,
		MimeType: file.Header.Get("Content-Type"),
	}, nil
}

// GenerateRandomFileName 生成随机文件名
func GenerateRandomFileName() string {
	randomStr, _ := GenerateRandomString(8)
	return randomStr
}

// IsAllowedFileType 检查文件类型是否允许
func IsAllowedFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedTypes {
		if ext == strings.ToLower(allowedType) {
			return true
		}
	}
	return false
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileExists 检查文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	if !FileExists(filePath) {
		return nil // 文件不存在，认为删除成功
	}
	return os.Remove(filePath)
}

// CreateDirectory 创建目录
func CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// CopyFile 复制文件
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
	return err
}

// GetMimeType 根据文件扩展名获取MIME类型
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".csv":  "text/csv",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(file *multipart.FileHeader, maxSize int64) error {
	if file.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", file.Size, maxSize)
	}
	return nil
}