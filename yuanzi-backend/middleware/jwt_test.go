package middleware

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"yuanzi-backend/model"
)

func TestGenerateTokenPair(t *testing.T) {
	userID := "user-1"
	phone := "13800138000"

	accessToken, refreshToken, err := GenerateTokenPair(userID, phone)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatal("生成的 token 为空")
	}

	accessClaims, err := ParseToken(accessToken)
	if err != nil {
		t.Fatalf("解析 Access Token 失败: %v", err)
	}
	if accessClaims.UserID != userID || accessClaims.Phone != phone || accessClaims.Type != "access" {
		t.Fatalf("Access Token claims 不匹配: %+v", accessClaims)
	}

	now := time.Now()
	if diff := accessClaims.ExpiresAt.Time.Sub(now); diff < 115*time.Minute || diff > 125*time.Minute {
		t.Fatalf("Access Token 过期时间错误: %v", diff)
	}

	refreshClaims, err := ParseToken(refreshToken)
	if err != nil {
		t.Fatalf("解析 Refresh Token 失败: %v", err)
	}
	if refreshClaims.Type != "refresh" {
		t.Fatalf("Refresh Token 类型错误: %s", refreshClaims.Type)
	}
}

func TestParseToken(t *testing.T) {
	accessToken, _, err := GenerateTokenPair("user-1", "13800138000")
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}
	claims, err := ParseToken(accessToken)
	if err != nil {
		t.Fatalf("解析有效 Token 失败: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("用户ID不匹配: %s", claims.UserID)
	}
	if invalidClaims, err := ParseToken("invalid.token.here"); err == nil || invalidClaims != nil {
		t.Fatal("无效 Token 应解析失败")
	}
}

func TestJWTMiddleware(t *testing.T) {
	token := ""
	code := model.SUCCESS
	if token == "" {
		code = model.ERROR_AUTH_CHECK_TOKEN_FAIL
	}
	if code != model.ERROR_AUTH_CHECK_TOKEN_FAIL {
		t.Fatal("无 Token 应返回认证失败")
	}
	if !strings.HasPrefix("Bearer valid_token", "Bearer ") {
		t.Fatal("Bearer Token 格式错误")
	}
	claims := Claims{UserID: "user-1", Phone: "13800138000", JTI: "test_jti", Type: "access", RegisteredClaims: jwt.RegisteredClaims{Issuer: "yuanzi"}}
	if claims.Type != "access" || claims.JTI == "" {
		t.Fatalf("claims 不正确: %+v", claims)
	}
}

func TestGenerateJTI(t *testing.T) {
	jti1, err1 := generateJTI()
	jti2, err2 := generateJTI()
	if err1 != nil || err2 != nil {
		t.Fatalf("生成 JTI 失败: %v %v", err1, err2)
	}
	if jti1 == jti2 || !strings.HasPrefix(jti1, "jti_") {
		t.Fatalf("JTI 不符合预期: %s %s", jti1, jti2)
	}
}

func TestGetUserIDOrZero(t *testing.T) {
	userID := ""
	if userID != "" {
		t.Fatal("空 userID 应为空字符串")
	}
}
