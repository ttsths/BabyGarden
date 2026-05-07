package push

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"yuanzi-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	apnsTokenTTL  = 50 * time.Minute
	apnsHostProd  = "https://api.push.apple.com"
	apnsHostDev   = "https://api.sandbox.push.apple.com"
	apnsPushType  = "alert"
	apnsTopicHint = "bundle_id is required"
)

// APNsClient APNs 推送客户端
// 注意：此实现基于 Token 认证（.p8 Key），不包含证书模式。
type APNsClient struct {
	bundleID string
	teamID   string
	keyID    string
	host     string
	client   *http.Client
	key      *ecdsa.PrivateKey

	mu        sync.Mutex
	cachedTok string
	issuedAt  time.Time
}

// NewAPNsClient 创建 APNs 客户端
func NewAPNsClient(cfg config.APNsConfig) (*APNsClient, error) {
	if cfg.BundleID == "" {
		return nil, errors.New(apnsTopicHint)
	}
	if cfg.TeamID == "" || cfg.KeyID == "" || cfg.KeyPath == "" {
		return nil, errors.New("apns team_id/key_id/key_path 不能为空")
	}
	keyBytes, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取 APNs Key 失败: %w", err)
	}
	key, err := jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("解析 APNs Key 失败: %w", err)
	}

	host := apnsHostProd
	if cfg.UseSandbox {
		host = apnsHostDev
	}
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &APNsClient{
		bundleID: cfg.BundleID,
		teamID:   cfg.TeamID,
		keyID:    cfg.KeyID,
		host:     host,
		client: &http.Client{
			Timeout: timeout,
		},
		key: key,
	}, nil
}

// Send 推送通知到 APNs
func (c *APNsClient) Send(ctx context.Context, deviceToken string, payload map[string]interface{}) error {
	if deviceToken == "" {
		return errors.New("device token 不能为空")
	}
	if payload == nil {
		payload = map[string]interface{}{}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化 payload 失败: %w", err)
	}

	token, err := c.getToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/3/device/%s", c.host, deviceToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("authorization", "bearer "+token)
	req.Header.Set("apns-topic", c.bundleID)
	req.Header.Set("apns-push-type", apnsPushType)
	req.Header.Set("content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送 APNs 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("APNs 返回错误: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *APNsClient) getToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cachedTok != "" && time.Since(c.issuedAt) < apnsTokenTTL {
		return c.cachedTok, nil
	}

	claims := jwt.MapClaims{
		"iss": c.teamID,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.keyID

	signed, err := token.SignedString(c.key)
	if err != nil {
		return "", fmt.Errorf("生成 APNs Token 失败: %w", err)
	}
	c.cachedTok = signed
	c.issuedAt = time.Now()
	return signed, nil
}
