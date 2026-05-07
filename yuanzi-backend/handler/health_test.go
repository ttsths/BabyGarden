package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	HealthCheck(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("健康检查状态码错误: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析健康检查响应失败: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("健康检查状态错误: %q", response.Status)
	}

	if response.Timestamp == "" {
		t.Fatalf("健康检查缺少 timestamp 字段")
	}
	if _, err := time.Parse(time.RFC3339, response.Timestamp); err != nil {
		t.Fatalf("timestamp 不是有效的 RFC3339 格式: %q, err=%v", response.Timestamp, err)
	}
}
