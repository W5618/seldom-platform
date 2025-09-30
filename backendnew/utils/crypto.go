package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// GenerateMD5 生成MD5哈希
func GenerateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GenerateSHA256 生成SHA256哈希
func GenerateSHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// EncryptAES AES加密
func EncryptAES(plaintext, key string) (string, error) {
	// 创建cipher.Block接口
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES AES解密
func DecryptAES(ciphertext, key string) (string, error) {
	// 解码base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建cipher.Block接口
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 获取nonce大小
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 分离nonce和密文
	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	// 转换为可打印字符
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	
	return string(bytes), nil
}

// GenerateAPIKey 生成API密钥
func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

// HashPassword 哈希密码（用于存储）
func HashPassword(password, salt string) string {
	return GenerateSHA256(password + salt)
}

// VerifyPassword 验证密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	return HashPassword(password, salt) == hashedPassword
}