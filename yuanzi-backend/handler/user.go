package handler

import (
	"net/http"
	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var user model.User
	if err := mysql.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: UserProfileResponse{
			ID:        user.ID,
			Phone:     user.Phone,
			Username:  user.Username,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
		},
	})
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if len(updates) > 0 {
		if err := mysql.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
			return
		}
	}

	var user model.User
	if err := mysql.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
		Data: UserProfileResponse{
			ID:        user.ID,
			Phone:     user.Phone,
			Username:  user.Username,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
		},
	})
}

type UserProfileResponse struct {
	ID        string `json:"id" example:"1"`
	Phone     string `json:"phone" example:"13800138000"`
	Username  string `json:"username,omitempty" example:"mom"`
	Nickname  string `json:"nickname" example:"小园子妈妈"`
	AvatarURL string `json:"avatar_url" example:"https://cdn.yuanzi.com/avatar.jpg"`
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname" binding:"omitempty,max=50" example:"新昵称"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url" example:"https://cdn.yuanzi.com/new_avatar.jpg"`
}
