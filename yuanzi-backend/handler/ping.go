package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping 返回 pong，供 Cloudflare Containers 健康检查使用
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
