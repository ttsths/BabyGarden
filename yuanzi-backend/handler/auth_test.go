package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

func TestPasswordLoginAcceptsPhoneOrUsername(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("pwdlogin"), "密码登录用户")
	user.Username = "pwduser"
	user.SetPassword("yuanzi123")
	if err := mysql.DB.Save(&user).Error; err != nil {
		t.Fatalf("更新测试用户失败: %v", err)
	}
	defer cleanupUsers(t, user.ID)

	for _, identifier := range []string{user.Phone, user.Username} {
		body := mustMarshal(t, PasswordLoginRequest{Identifier: identifier, Password: "yuanzi123"})
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/auth/password-login", bytes.NewReader(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		PasswordLogin(ctx)

		if recorder.Code != http.StatusOK {
			t.Fatalf("密码登录失败 identifier=%s status=%d body=%s", identifier, recorder.Code, recorder.Body.String())
		}
		var response struct {
			Data LoginResponse `json:"data"`
		}
		decodeResponse(t, recorder.Body.Bytes(), &response)
		if response.Data.AccessToken == "" || response.Data.RefreshToken == "" || response.Data.User == nil {
			t.Fatalf("密码登录响应缺少 token 或用户信息: %+v", response.Data)
		}
	}
}

func TestPasswordLoginRejectsWrongPassword(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("pwdfail"), "密码失败用户")
	user.SetPassword("yuanzi123")
	if err := mysql.DB.Save(&user).Error; err != nil {
		t.Fatalf("更新测试用户失败: %v", err)
	}
	defer cleanupUsers(t, user.ID)

	body := mustMarshal(t, PasswordLoginRequest{Identifier: user.Phone, Password: "wrong"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/auth/password-login", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	PasswordLogin(ctx)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("错误密码应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestLogoutAddsTokenToBlacklist(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("logout"), "退出用户")
	defer cleanupUsers(t, user.ID)

	_, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Phone)
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}
	claims, err := middleware.ParseToken(refreshToken)
	if err != nil {
		t.Fatalf("解析Token失败: %v", err)
	}

	body := mustMarshal(t, LogoutRequest{RefreshToken: refreshToken})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	Logout(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("退出登录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	blacklisted, err := middleware.IsBlacklisted(claims.JTI)
	if err != nil {
		t.Fatalf("校验黑名单失败: %v", err)
	}
	if !blacklisted {
		t.Fatalf("Refresh Token 未加入黑名单")
	}
}

func TestLogoutRejectsInvalidToken(t *testing.T) {
	setupFamilyTestDB(t)

	body := mustMarshal(t, LogoutRequest{RefreshToken: "invalid"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	Logout(ctx)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("无效Token应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp model.Response
	decodeResponse(t, recorder.Body.Bytes(), &resp)
	if resp.Code != model.ERROR_NOT_AUTH {
		t.Fatalf("返回码错误: %+v", resp)
	}
}
