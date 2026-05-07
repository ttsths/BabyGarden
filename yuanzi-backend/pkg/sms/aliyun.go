package sms

import "fmt"

// SMSConfig 短信服务配置
type SMSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// SMSClient 短信客户端
// 当前保持轻量封装，避免未接入 SDK 时阻断编译。
type SMSClient struct {
	config *SMSConfig
}

// NewSMSClient 创建短信客户端
func NewSMSClient(config *SMSConfig) *SMSClient {
	return &SMSClient{config: config}
}

// SendVerificationCode 发送验证码
func (c *SMSClient) SendVerificationCode(phone, code string) error {
	if c == nil || c.config == nil {
		return fmt.Errorf("sms config is nil")
	}
	if FormatPhone(phone) == "" {
		return fmt.Errorf("invalid phone")
	}
	if len(code) != 6 {
		return fmt.Errorf("invalid verification code")
	}
	if c.config.SignName == "" || c.config.TemplateCode == "" {
		return fmt.Errorf("sms template config is incomplete")
	}
	return nil
}

// FormatPhone 格式化手机号
func FormatPhone(phone string) string {
	result := ""
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			result += string(r)
		}
	}
	return result
}
