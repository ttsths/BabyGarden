package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

func TestRegisterDeviceCreate(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("device"), "推送用户")
	defer cleanupUsers(t, user.ID)
	defer cleanupPushDevices(t, user.ID)

	body := mustMarshal(t, RegisterDeviceRequest{
		Platform:    "ios",
		DeviceToken: "token-123",
		Alias:       "user-01",
		Tags:        []string{"feed", "sleep"},
	})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/device/register", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", user.ID)

	RegisterDevice(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("注册设备失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Code != model.SUCCESS {
		t.Fatalf("业务返回失败: %+v", response)
	}

	var device model.PushDevice
	if err := mysql.DB.Where("user_id = ? AND device_token = ?", user.ID, "token-123").First(&device).Error; err != nil {
		t.Fatalf("设备未写入: %v", err)
	}

	if device.Platform != "ios" || device.Alias != "user-01" || device.IsActive != 1 {
		t.Fatalf("设备字段不符合预期: %+v", device)
	}

	var tags []string
	if err := json.Unmarshal(device.Tags, &tags); err != nil {
		t.Fatalf("解析 tags 失败: %v", err)
	}
	if len(tags) != 2 || tags[0] != "feed" || tags[1] != "sleep" {
		t.Fatalf("tags 不符合预期: %+v", tags)
	}
}

func TestRegisterDeviceUpdate(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("device2"), "推送用户2")
	defer cleanupUsers(t, user.ID)
	defer cleanupPushDevices(t, user.ID)

	seed := model.PushDevice{
		UserID:      user.ID,
		Platform:    "android",
		DeviceToken: "token-456",
		Alias:       "old",
		Tags:        mustJSON(t, []string{"old"}),
		IsActive:    0,
		LastUsedAt:  time.Now().Add(-2 * time.Hour),
		CreatedAt:   time.Now().Add(-24 * time.Hour),
	}
	if err := mysql.DB.Create(&seed).Error; err != nil {
		t.Fatalf("创建种子设备失败: %v", err)
	}

	body := mustMarshal(t, RegisterDeviceRequest{
		Platform:    "android",
		DeviceToken: "token-456",
		Alias:       "new",
		Tags:        []string{"sleep"},
	})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/device/register", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", user.ID)

	RegisterDevice(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("更新设备失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var device model.PushDevice
	if err := mysql.DB.Where("id = ?", seed.ID).First(&device).Error; err != nil {
		t.Fatalf("设备更新后查询失败: %v", err)
	}

	if device.Alias != "new" || device.IsActive != 1 {
		t.Fatalf("设备更新未生效: %+v", device)
	}

	if time.Since(device.LastUsedAt) > time.Minute {
		t.Fatalf("last_used_at 未更新: %+v", device.LastUsedAt)
	}
}

func cleanupPushDevices(t *testing.T, userIDs ...string) {
	t.Helper()
	if len(userIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("user_id IN ?", userIDs).Delete(&model.PushDevice{}).Error
}
