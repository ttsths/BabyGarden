package handler

import (
	"net/http"
	"time"

	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns the service health status.
func HealthCheck(c *gin.Context) {
	dbOK := mysql.IsConnected()
	redisOK := gredis.IsConnected()
	status := http.StatusOK
	if !dbOK || !redisOK {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":    map[bool]string{true: "ok", false: "degraded"}[dbOK && redisOK],
		"database":  dbOK,
		"redis":     redisOK,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
