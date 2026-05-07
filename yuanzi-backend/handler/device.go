package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterDevice 注册推送设备
// @Summary 注册推送设备
// @Description 注册设备 token，用于推送通知
// @Tags 推送
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body RegisterDeviceRequest true "设备信息"
// @Success 200 {object} model.Response
// @Router /api/v1/device/register [post]
func RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var tags model.JSON
	if len(req.Tags) > 0 {
		if data, err := json.Marshal(req.Tags); err == nil {
			tags = model.JSON(data)
		}
	}

	var device model.PushDevice
	err := mysql.DB.Where("user_id = ? AND device_token = ?", userID, req.DeviceToken).First(&device).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询设备失败"})
		return
	}
	if err == nil {
		updates := map[string]interface{}{
			"platform":     req.Platform,
			"alias":        req.Alias,
			"tags":         tags,
			"is_active":    1,
			"last_used_at": time.Now(),
		}
		if err := mysql.DB.Model(&model.PushDevice{}).Where("id = ?", device.ID).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新设备失败"})
			return
		}
		c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "注册成功"})
		return
	}

	newDevice := model.PushDevice{
		UserID:      userID,
		Platform:    req.Platform,
		DeviceToken: req.DeviceToken,
		Alias:       req.Alias,
		Tags:        tags,
		IsActive:    1,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	if err := mysql.DB.Create(&newDevice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "注册失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "注册成功"})
}

type RegisterDeviceRequest struct {
	Platform    string   `json:"platform" binding:"required,oneof=ios android"`
	DeviceToken string   `json:"device_token" binding:"required"`
	Alias       string   `json:"alias"`
	Tags        []string `json:"tags"`
}
