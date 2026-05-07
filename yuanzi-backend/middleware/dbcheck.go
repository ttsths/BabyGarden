package middleware

import (
	"net/http"

	"yuanzi-backend/mysql"
	"yuanzi-backend/model"

	"github.com/gin-gonic/gin"
)

// RequireDB 检查数据库是否可用，不可用则返回 503
// 允许服务在数据库未连接时优雅降级而不是 panic
func RequireDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !mysql.IsConnected() {
			c.JSON(http.StatusServiceUnavailable, model.Response{
				Code: model.ERROR,
				Msg:  "数据库服务暂不可用，请稍后重试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
