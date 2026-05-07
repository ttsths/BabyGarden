package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"
	"yuanzi-backend/router"
)

var e2eSetupOnce sync.Once

// setupE2E initializes the full test environment once per test run.
func setupE2E(t *testing.T) {
	t.Helper()
	e2eSetupOnce.Do(func() {
		projectRoot := findProjectRoot(t)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(filepath.Join(projectRoot, "config"))
		config.Setup()
		logger.Setup()
		mysql.Setup()
		gredis.Setup()
		gin.SetMode(gin.TestMode)
	})
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(dir, "config", "config.yaml")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("未找到项目根目录")
	return ""
}

// e2eRouter returns a fully initialized Gin engine for E2E testing.
func e2eRouter(t *testing.T) *gin.Engine {
	t.Helper()
	setupE2E(t)
	return router.SetupRouter()
}

// doRequest sends an HTTP request through the full Gin middleware chain.
func doRequest(t *testing.T, r *gin.Engine, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// getJSON sends a GET request with optional auth header.
func getJSON(t *testing.T, r *gin.Engine, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return doRequest(t, r, http.MethodGet, path, "", headers)
}

// postJSON sends a POST request with JSON body and optional auth header.
func postJSON(t *testing.T, r *gin.Engine, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	return doRequest(t, r, http.MethodPost, path, string(b), headers)
}

// parseResponse unmarshals the response body into a model.Response.
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) model.Response {
	t.Helper()
	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, w.Body.String())
	}
	return resp
}

// assertStatus checks HTTP status code.
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int, context string) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("%s: 期望状态码 %d, 实际 %d, body=%s", context, expected, w.Code, w.Body.String())
	}
}

// assertCode checks business code in model.Response.
func assertCode(t *testing.T, resp model.Response, expected int, context string) {
	t.Helper()
	if resp.Code != expected {
		t.Errorf("%s: 期望业务码 %d, 实际 %d, msg=%s", context, expected, resp.Code, resp.Msg)
	}
}

// seedAdminUser creates an admin user for E2E testing and returns cleanup function.
func seedAdminUser(t *testing.T) (model.User, func()) {
	t.Helper()
	setupE2E(t)

	admin := model.User{
		Phone:    "13899990001",
		Nickname: "E2E管理员",
		IsAdmin:  1,
		Password: "e2eadmin",
		Status:   1,
	}
	if err := mysql.DB.Create(&admin).Error; err != nil {
		t.Fatalf("创建E2E管理员失败: %v", err)
	}
	cleanup := func() {
		mysql.DB.Where("id = ?", admin.ID).Delete(&model.User{})
	}
	return admin, cleanup
}

// seedRegularUser creates a regular (non-admin) user for E2E testing.
func seedRegularUser(t *testing.T) (model.User, func()) {
	t.Helper()
	setupE2E(t)

	user := model.User{
		Phone:    "13899990002",
		Nickname: "E2E普通用户",
		IsAdmin:  0,
		Status:   1,
	}
	if err := mysql.DB.Create(&user).Error; err != nil {
		t.Fatalf("创建E2E普通用户失败: %v", err)
	}
	cleanup := func() {
		mysql.DB.Where("id = ?", user.ID).Delete(&model.User{})
	}
	return user, cleanup
}
