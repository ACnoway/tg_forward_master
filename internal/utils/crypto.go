package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"strings"
)

// Encrypt 加密字符串
func Encrypt(plaintext, key string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("创建加密器失败: %w", err)
	}

	// 生成随机IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("生成IV失败: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func Decrypt(ciphertext, key string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("解码失败: %w", err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("创建解密器失败: %w", err)
	}

	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("密文太短")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return string(data), nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// FormatTelegramUser 格式化Telegram用户信息
func FormatTelegramUser(firstName, lastName, username string, telegramID int64) string {
	parts := []string{}

	if firstName != "" {
		name := firstName
		if lastName != "" {
			name += " " + lastName
		}
		parts = append(parts, name)
	}

	if username != "" {
		parts = append(parts, "@"+username)
	}

	parts = append(parts, fmt.Sprintf("ID:%d", telegramID))

	return strings.Join(parts, " | ")
}

// EscapeMarkdown 转义Markdown特殊字符
func EscapeMarkdown(text string) string {
	special := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range special {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
