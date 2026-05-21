package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"
	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"

	"github.com/gin-gonic/gin"
)

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Type  string `json:"type" binding:"omitempty,oneof=login reset bind" example:"login"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	ExpiresIn int `json:"expires_in" example:"300"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Code  string `json:"code" binding:"required,len=6"`
}

// PasswordLoginRequest 用户名/手机号密码登录请求。
type PasswordLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User         *model.UserInfo `json:"user"`
	AccessToken  string          `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string          `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	ExpiresIn    int             `json:"expires_in" example:"7200"`
	IsNewUser    bool            `json:"is_new_user"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 退出登录请求
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// SendVerificationCode 发送短信验证码
func SendVerificationCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	minuteKey := fmt.Sprintf("sms:limit:%s:minute", req.Phone)
	exists, err := gredis.Exists(minuteKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "限流检查失败"})
		return
	}
	if exists {
		c.JSON(http.StatusTooManyRequests, model.Response{Code: model.ERROR_RATE_LIMITED, Msg: "1分钟内只能发送一次"})
		return
	}

	hourKey := fmt.Sprintf("sms:limit:%s:hour", req.Phone)
	count, err := gredis.Incr(hourKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "限流检查失败"})
		return
	}
	if count == 1 {
		_ = gredis.Expire(hourKey, 3600)
	}
	if count > 5 {
		c.JSON(http.StatusTooManyRequests, model.Response{Code: model.ERROR_RATE_LIMITED, Msg: "1小时内最多发送5次"})
		return
	}

	code, err := generateVerificationCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成验证码失败"})
		return
	}

	codeKey := fmt.Sprintf("sms:code:%s:%s", req.Phone, "login")
	if err := gredis.SetEx(codeKey, code, 300); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "保存验证码失败"})
		return
	}
	_ = gredis.SetEx(minuteKey, "1", 60)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "验证码已发送",
		Data: SendCodeResponse{ExpiresIn: 300},
	})
}

// Login 验证码登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	usedKey := fmt.Sprintf("sms:used:%s:%s", req.Phone, req.Code)
	used, err := gredis.Exists(usedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "验证码校验失败"})
		return
	}
	if used {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "验证码已使用"})
		return
	}

	codeKey := fmt.Sprintf("sms:code:%s:%s", req.Phone, "login")
	storedCode, err := gredis.Get(codeKey)
	if err != nil || storedCode == "" || storedCode != req.Code {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "验证码错误或已过期"})
		return
	}

	_ = gredis.SetEx(usedKey, "1", 300)
	_ = gredis.Del(codeKey)

	user, isNewUser, err := findOrCreateUser(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "登录失败"})
		return
	}

	if err := mysql.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
		"last_login_ip": c.ClientIP(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新登录信息失败"})
		return
	}

	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成Token失败"})
		return
	}

	userInfo := user.ToUserInfo()
	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "登录成功",
		Data: LoginResponse{
			User:         &userInfo,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
			IsNewUser:    isNewUser,
		},
	})
}

// PasswordLogin 用户名或手机号 + 密码登录。
func PasswordLogin(c *gin.Context) {
	var req PasswordLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var user model.User
	if err := mysql.DB.Where("phone = ? OR username = ?", req.Identifier, req.Identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "账号或密码错误"})
		return
	}
	if user.Status != 1 || !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "账号或密码错误"})
		return
	}

	if err := mysql.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
		"last_login_ip": c.ClientIP(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新登录信息失败"})
		return
	}

	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成Token失败"})
		return
	}

	userInfo := user.ToUserInfo()
	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "登录成功",
		Data: LoginResponse{
			User:         &userInfo,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
			IsNewUser:    false,
		},
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	claims, err := middleware.ParseToken(req.RefreshToken)
	if err != nil || claims == nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "无效的Refresh Token"})
		return
	}
	if claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "Token类型错误"})
		return
	}

	ttl := 7 * 24 * time.Hour
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
		if ttl < time.Minute {
			ttl = time.Minute
		}
	}
	if err := middleware.AddToBlacklist(claims.JTI, ttl); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "退出失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "退出成功"})
}

// RefreshToken 刷新 Access Token
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	claims, err := middleware.ParseToken(req.RefreshToken)
	if err != nil || claims == nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "无效的Refresh Token"})
		return
	}
	if claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "Token类型错误"})
		return
	}

	blacklisted, err := middleware.IsBlacklisted(claims.JTI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "Token状态校验失败"})
		return
	}
	if blacklisted {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "Token已被撤销"})
		return
	}

	var user model.User
	if err := mysql.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "用户不存在"})
		return
	}

	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成Token失败"})
		return
	}

	if claims.ExpiresAt != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl < time.Minute {
			ttl = time.Minute
		}
		_ = middleware.AddToBlacklist(claims.JTI, ttl)
	}

	userInfo := user.ToUserInfo()
	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "刷新成功",
		Data: LoginResponse{
			User:         &userInfo,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
			IsNewUser:    false,
		},
	})
}

func generateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func findOrCreateUser(phone string) (*model.User, bool, error) {
	var user model.User
	if err := mysql.DB.Where("phone = ?", phone).First(&user).Error; err == nil {
		return &user, false, nil
	}

	user = model.User{
		Phone:    phone,
		Username: phone,
		Nickname: "用户" + phone[len(phone)-4:],
		Status:   1,
	}
	if err := mysql.DB.Create(&user).Error; err != nil {
		return nil, false, err
	}
	return &user, true, nil
}
