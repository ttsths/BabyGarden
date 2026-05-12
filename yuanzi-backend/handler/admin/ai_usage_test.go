package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

func TestGetAIUsage(t *testing.T) {
	setupAdminTestDB(t)
	ensureAIUsageLogTable(t)

	user := createTestUser(t, uniquePhone("aiusage"), "AI使用用户", 0)
	defer cleanupUsers(t, user.ID)
	defer cleanupAIUsageLogs(t, user.ID)

	createAIUsageLog(t, user.ID, "grokai", 100, time.Now())
	createAIUsageLog(t, user.ID, "dashscope", 50, time.Now())

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/admin/ai/usage?page=1&page_size=10&provider=grokai", nil)

	GetAIUsage(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("查询AI使用记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			List       []model.AIUsageLog `json:"list"`
			Pagination model.Pagination   `json:"pagination"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &response)
	if response.Data.Pagination.Total != 1 || len(response.Data.List) != 1 {
		t.Fatalf("筛选结果错误: %+v", response.Data)
	}
	if response.Data.List[0].Provider != "grokai" {
		t.Fatalf("provider筛选失效: %+v", response.Data.List[0])
	}
}

func TestGetAIUsageSummary(t *testing.T) {
	setupAdminTestDB(t)
	ensureAIUsageLogTable(t)

	user := createTestUser(t, uniquePhone("aisummary"), "AI汇总用户", 0)
	defer cleanupUsers(t, user.ID)
	defer cleanupAIUsageLogs(t, user.ID)

	createAIUsageLog(t, user.ID, "grokai", 100, time.Now())
	createAIUsageLog(t, user.ID, "grokai", 40, time.Now())

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/admin/ai/usage/summary?period=day&days=7", nil)

	GetAIUsageSummary(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("查询AI使用汇总失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			Items []AIUsageSummaryItem `json:"items"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &response)
	if len(response.Data.Items) == 0 {
		t.Fatalf("汇总结果为空: %+v", response.Data)
	}
}

func TestGetUserAIUsage(t *testing.T) {
	setupAdminTestDB(t)
	ensureAIUsageLogTable(t)

	user := createTestUser(t, uniquePhone("aiuser"), "AI明细用户", 0)
	defer cleanupUsers(t, user.ID)
	defer cleanupAIUsageLogs(t, user.ID)

	createAIUsageLog(t, user.ID, "cloudflare_workers_ai", 88, time.Now())

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/admin/ai/usage/"+user.ID, nil)
	ctx.Params = gin.Params{{Key: "userId", Value: user.ID}}

	GetUserAIUsage(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("查询用户AI使用详情失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			Overview AIUsageOverview `json:"overview"`
			Logs     struct {
				List []model.AIUsageLog `json:"list"`
			} `json:"logs"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &response)
	if response.Data.Overview.TodayTotalTokens < 88 || len(response.Data.Logs.List) != 1 {
		t.Fatalf("用户AI详情错误: %+v", response.Data)
	}
}

func ensureAIUsageLogTable(t *testing.T) {
	t.Helper()
	if mysql.DB == nil {
		t.Skip("数据库不可用，跳过AI使用记录管理端测试")
	}
	if err := mysql.DB.AutoMigrate(&model.AIUsageLog{}); err != nil {
		t.Fatalf("创建AI使用记录表失败: %v", err)
	}
}

func createAIUsageLog(t *testing.T, userID, provider string, totalTokens int, createdAt time.Time) {
	t.Helper()
	log := model.AIUsageLog{
		UserID:       userID,
		Provider:     provider,
		Model:        "test-model",
		InputTokens:  totalTokens / 2,
		OutputTokens: totalTokens - totalTokens/2,
		TotalTokens:  totalTokens,
		RequestType:  "chat",
		Status:       "success",
		CreatedAt:    createdAt,
	}
	if err := mysql.DB.Create(&log).Error; err != nil {
		t.Fatalf("创建AI使用记录失败: %v", err)
	}
}

func cleanupAIUsageLogs(t *testing.T, userIDs ...string) {
	t.Helper()
	if len(userIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("user_id IN ?", userIDs).Delete(&model.AIUsageLog{}).Error
}
