package handler

import (
	"os"

	"yuanzi-backend/config"
	"yuanzi-backend/pkg/storage"

	"github.com/gin-gonic/gin"
)

// DebugEnv 返回容器内的环境变量和 Viper 配置状态（仅用于诊断 R2）
func DebugEnv(c *gin.Context) {
	r2Keys := []string{
		"R2_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_ACCESS_KEY_SECRET",
		"R2_BUCKET", "R2_PUBLIC_URL",
	}

	mask := func(v string) string {
		if v == "" {
			return "(empty)"
		}
		if len(v) > 8 {
			return v[:4] + "***" + v[len(v)-4:]
		}
		return v
	}

	envVars := make(map[string]string)
	workerVars := make(map[string]string)
	for _, key := range r2Keys {
		envVars[key] = mask(os.Getenv(key))
		workerVars[key] = mask(storage.GetWorkerVar(key))
	}

	viperVars := map[string]string{
		"r2.account_id":        mask(config.GlobalConfig.R2.AccountID),
		"r2.access_key_id":     mask(config.GlobalConfig.R2.AccessKeyID),
		"r2.access_key_secret": mask(config.GlobalConfig.R2.AccessKeySecret),
		"r2.bucket":            mask(config.GlobalConfig.R2.Bucket),
		"r2.public_url":        mask(config.GlobalConfig.R2.PublicURL),
	}

	c.JSON(200, gin.H{
		"os_getenv":        envVars,
		"worker_vars":      workerVars,
		"viper_config":     viperVars,
		"storage_provider": config.GlobalConfig.Storage.Provider,
		"headers": gin.H{
			"x_worker_vars_raw": len(c.GetHeader("X-Worker-Vars")) > 0,
			"x_test_ping":       c.GetHeader("X-Test-Ping"),
		},
	})
}
