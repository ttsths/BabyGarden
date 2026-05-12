package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yuanzi-backend/config"
	"yuanzi-backend/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.GlobalConfig = config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
		},
	}
}

// TestAdminMiddlewareChain tests that admin routes have correct auth middleware
func TestAdminRoutesRequireJWT(t *testing.T) {
	r := SetupRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		// Login 不需要认证 → 缺少 body 返回 400 (不是 401)
		{"login no body", "POST", "/api/v1/admin/login", http.StatusBadRequest},
		// 所有 admin 受保护接口 → 无 token 应返回 401
		{"stats overview no auth", "GET", "/api/v1/admin/stats/overview", http.StatusUnauthorized},
		{"stats daily no auth", "GET", "/api/v1/admin/stats/daily", http.StatusUnauthorized},
		{"users list no auth", "GET", "/api/v1/admin/users", http.StatusUnauthorized},
		{"families list no auth", "GET", "/api/v1/admin/families", http.StatusUnauthorized},
		{"babies list no auth", "GET", "/api/v1/admin/babies", http.StatusUnauthorized},
		{"photos list no auth", "GET", "/api/v1/admin/photos", http.StatusUnauthorized},
		{"records list no auth", "GET", "/api/v1/admin/records", http.StatusUnauthorized},
		{"ai usage list no auth", "GET", "/api/v1/admin/ai/usage", http.StatusUnauthorized},
		{"ai usage summary no auth", "GET", "/api/v1/admin/ai/usage/summary", http.StatusUnauthorized},
		{"ai usage user no auth", "GET", "/api/v1/admin/ai/usage/user-1", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				t.Errorf("%s %s: expected %d, got %d (msg: %v)",
					tt.method, tt.path, tt.wantStatus, w.Code, resp["msg"])
			}
		})
	}
}

// TestAdminRoutesRejectNonAdminToken tests that a valid but non-admin token is rejected
func TestAdminRoutesRejectNonAdminToken(t *testing.T) {
	r := SetupRouter()

	// Generate a regular user token (is_admin = false)
	accessToken, _, err := middleware.GenerateTokenPair("user-1", "13800138000")
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/admin/stats/overview", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Regular user should get 403 Forbidden (not admin)
	if w.Code != http.StatusForbidden {
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		t.Errorf("Non-admin token: expected 403, got %d (msg: %v)", w.Code, resp["msg"])
	}
}

// TestLoginIsPublic ensures admin login is accessible without auth
func TestAdminLoginIsPublic(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("POST", "/api/v1/admin/login", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Without body, should be 400 (not 401 — means JWT middleware is NOT applied to login)
	if w.Code == http.StatusUnauthorized {
		t.Fatal("admin/login should NOT require JWT auth, but got 401")
	}
}
