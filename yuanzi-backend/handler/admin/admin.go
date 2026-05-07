package admin

import (
	"net/http"
	"yuanzi-backend/config"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AdminLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	User      struct {
		ID        string `json:"id"`
		Phone     string `json:"phone"`
		Nickname  string `json:"nickname"`
		IsAdmin   int8   `json:"is_admin"`
	} `json:"user"`
}

// AdminLogin handles admin authentication via phone + password.
// POST /api/v1/admin/login
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}

	var user model.User
	if err := mysql.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{
			Code: model.ERROR_AUTH,
			Msg:  "手机号、密码错误或非管理员",
		})
		return
	}

	if user.IsAdmin != 1 {
		c.JSON(http.StatusForbidden, model.Response{
			Code: model.ERROR_AUTH,
			Msg:  "非管理员账户",
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusForbidden, model.Response{
			Code: model.ERROR_AUTH,
			Msg:  "手机号、密码错误或非管理员",
		})
		return
	}

	// Generate JWT token with admin flag
	claims := middlewareClaims{
		UserID:  user.ID,
		Phone:   user.Phone,
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "yuanzi-admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code: model.ERROR,
			Msg:  "生成Token失败",
		})
		return
	}

	var resp AdminLoginResponse
	resp.Token = tokenString
	resp.ExpiresIn = 86400
	resp.User.ID = user.ID
	resp.User.Phone = user.Phone
	resp.User.Nickname = user.Nickname
	resp.User.IsAdmin = user.IsAdmin

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "登录成功",
		Data: resp,
	})
}

type middlewareClaims struct {
	UserID  string `json:"user_id"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}
