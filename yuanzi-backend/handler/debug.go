package handler

import (
	"os"

	"yuanzi-backend/config"

	"github.com/gin-gonic/gin"
)

// DebugEnv 返回容器内的环境变量和 Viper 配置状态（仅用于诊断 R2）
func DebugEnv(c *gin.Context) {
	r2Keys := []string{
		"R2_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_ACCESS_KEY_SECRET",
		"R2_BUCKET", "R2_PUBLIC_URL",
	}

	envVars := make(map[string]string)
	for _, key := range r2Keys {
		val := os.Getenv(key)
		if val != "" {
			masked := val
			if len(val) > 8 {
				masked = val[:4] + "***" + val[len(val)-4:]
			}
			envVars[key] = masked
		} else {
			envVars[key] = "(empty)"
		}
	}

	viperVars := map[string]string{
		"r2.account_id":        config.GlobalConfig.R2.AccountID,
		"r2.access_key_id":     config.GlobalConfig.R2.AccessKeyID,
		"r2.access_key_secret": config.GlobalConfig.R2.AccessKeySecret,
		"r2.bucket":            config.GlobalConfig.R2.Bucket,
		"r2.public_url":        config.GlobalConfig.R2.PublicURL,
	}
	maskedViper := make(map[string]string)
	for k, v := range viperVars {
		if v != "" {
			masked := v
			if len(v) > 8 {
				masked = v[:4] + "***" + v[len(v)-4:]
			}
			maskedViper[k] = masked
		} else {
			maskedViper[k] = "(empty)"
		}
	}

	c.JSON(200, gin.H{
		"os_getenv":   envVars,
		"viper_config": maskedViper,
		"storage_provider": config.GlobalConfig.Storage.Provider,
	})
}
