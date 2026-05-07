package middleware

import (
	"net/http"
	"yuanzi-backend/model"

	"github.com/gin-gonic/gin"
)

// AdminAuth checks if the authenticated user has admin privileges.
// Must be used AFTER the JWT middleware which sets claims in context.
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": model.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  "Authentication required",
				"data": nil,
			})
			c.Abort()
			return
		}

		claims, ok := claimsVal.(*Claims)
		if !ok || !claims.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code": model.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  "Admin access required",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
