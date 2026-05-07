package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"yuanzi-backend/config"
)

// Client OSS客户端
type Client struct {
	region          string
	bucket          string
	accessKeyID     string
	accessKeySecret string
	cdnDomain       string
}

// NewClient 创建OSS客户端
func NewClient() *Client {
	cfg := config.GlobalConfig.OSS
	return &Client{
		region:          cfg.Region,
		bucket:          cfg.Bucket,
		accessKeyID:     cfg.AccessKeyID,
		accessKeySecret: cfg.AccessKeySecret,
		cdnDomain:       cfg.CdnDomain,
	}
}

// UploadPolicy 上传策略
type UploadPolicy struct {
	Expiration string          `json:"expiration"`
	Conditions [][]interface{} `json:"conditions"`
}

// PostSignature 直传签名结果
type PostSignature struct {
	OSSAccessKeyID string `json:"OSSAccessKeyId"`
	Policy         string `json:"policy"`
	Signature      string `json:"signature"`
}

// GetPostSignature 生成直传签名
func (c *Client) GetPostSignature(key string, maxSize int64, expireSeconds int) (*PostSignature, string, error) {
	if expireSeconds <= 0 {
		expireSeconds = 300 // 默认5分钟
	}

	// 构建上传策略
	expiration := time.Now().UTC().Add(time.Duration(expireSeconds) * time.Second).Format("2006-01-02T15:04:05Z")
	policy := UploadPolicy{
		Expiration: expiration,
		Conditions: [][]interface{}{
			{"key", key},
			{"bucket", c.bucket},
			{"content-length-range", 0, maxSize},
		},
	}

	policyJSON, _ := json.Marshal(policy)
	policyBase64 := base64.StdEncoding.EncodeToString(policyJSON)

	// 计算签名
	signature := c.sign(policyBase64)

	return &PostSignature{
		OSSAccessKeyID: c.accessKeyID,
		Policy:         policyBase64,
		Signature:      signature,
	}, c.getUploadURL(), nil
}

// sign 计算签名
func (c *Client) sign(data string) string {
	h := hmac.New(sha1.New, []byte(c.accessKeySecret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// getUploadURL 获取上传URL
func (c *Client) getUploadURL() string {
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", c.bucket, c.region)
}

// GetURL 获取文件访问URL
func (c *Client) GetURL(key string) string {
	if c.cdnDomain != "" {
		return fmt.Sprintf("https://%s/%s", c.cdnDomain, key)
	}
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", c.bucket, c.region, key)
}

// GetThumbnailURL 获取缩略图URL
func (c *Client) GetThumbnailURL(key string, width int) string {
	url := c.GetURL(key)
	if width > 0 {
		return url + fmt.Sprintf("?x-oss-process=image/resize,w_%d", width)
	}
	return url
}

// DeleteObject 删除 OSS 对象
func (c *Client) DeleteObject(key string) error {
	if key == "" {
		return fmt.Errorf("object key 不能为空")
	}
	date := time.Now().UTC().Format(http.TimeFormat)
	resource := fmt.Sprintf("/%s/%s", c.bucket, key)
	stringToSign := fmt.Sprintf("DELETE\n\n\n%s\n%s", date, resource)
	signature := c.sign(stringToSign)
	url := fmt.Sprintf("%s/%s", c.getUploadURL(), key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", fmt.Sprintf("OSS %s:%s", c.accessKeyID, signature))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("OSS 删除失败: status=%d", resp.StatusCode)
	}
	return nil
}
