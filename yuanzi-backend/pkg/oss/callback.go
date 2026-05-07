package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"strings"
)

// VerifyCallbackToken 校验回调鉴权 Token（与 OSS 配置的回调密钥一致）。
func VerifyCallbackToken(secret, provided string) error {
	if secret == "" {
		return errors.New("callback secret is empty")
	}
	if provided == "" {
		return errors.New("callback token is empty")
	}
	if subtleTokenMatch(secret, provided) {
		return nil
	}
	return errors.New("invalid callback token")
}

// VerifyCallbackSignature 校验 OSS 回调签名（Authorization 头）。
// 说明：使用回调 body 与 accessKeySecret 计算 HMAC-SHA1，base64 后与 authorization 进行对比。
func VerifyCallbackSignature(accessKeySecret string, body []byte, authorization string) error {
	if accessKeySecret == "" {
		return errors.New("accessKeySecret is empty")
	}
	if authorization == "" {
		return errors.New("authorization is empty")
	}

	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !subtleTokenMatch(expected, authorization) {
		return errors.New("invalid authorization signature")
	}
	return nil
}

func subtleTokenMatch(expected, provided string) bool {
	return strings.TrimSpace(expected) == strings.TrimSpace(provided)
}
