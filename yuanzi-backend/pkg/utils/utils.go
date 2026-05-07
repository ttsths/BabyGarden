package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
)

// GenerateVerificationCode 生成6位验证码
func GenerateVerificationCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// GenerateInviteCode 生成8位邀请码
func GenerateInviteCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsValidDate 验证日期格式
func IsValidDate(date string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	matched, _ := regexp.MatchString(pattern, date)
	return matched
}
