package push

import (
	"sync"

	"yuanzi-backend/config"
)

var (
	apnsOnce   sync.Once
	apnsClient *APNsClient
	apnsErr    error
)

// IsAPNsEnabled 判断 APNs 是否配置完整
func IsAPNsEnabled(cfg config.APNsConfig) bool {
	return cfg.TeamID != "" && cfg.KeyID != "" && cfg.BundleID != "" && cfg.KeyPath != ""
}

// APNs 获取 APNs 客户端（单例）
func APNs() (*APNsClient, error) {
	apnsOnce.Do(func() {
		apnsClient, apnsErr = NewAPNsClient(config.GlobalConfig.Push.APNs)
	})
	return apnsClient, apnsErr
}
