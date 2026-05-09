package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

// TestCreateRecordWithEmptyContent 验证 Bug #2 修复：空 content 时不报 500
func TestCreateRecordWithEmptyContent(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("rcc"), "记录测试管理员", 1)
	family := createTestFamily(t, admin.ID, "记录测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "记录宝宝",
		Gender:    1,
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	mysql.DB.Create(&baby)

	// 不传 content 字段，验证默认空 JSON 处理
	body := mustMarshal(t, map[string]interface{}{
		"type":       "feeding",
		"baby_id":    baby.ID,
		"family_id":  family.ID,
		"started_at": time.Now().Format("2006-01-02T15:04:05Z"),
		"note":       "不含 content 的测试记录",
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/records",
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	CreateRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("空 content 创建记录应成功: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Code != model.SUCCESS {
		t.Fatalf("业务返回失败: %+v", resp)
	}

	// 验证记录已存入 DB
	var record model.Record
	if err := mysql.DB.First(&record, "id = ?", resp.Data.ID).Error; err != nil {
		t.Fatalf("查询创建的记录失败: %v", err)
	}
	if record.Type != model.RecordTypeFeeding {
		t.Errorf("记录类型错误: %s", record.Type)
	}
	// Content 应为合法的 JSON（{}）
	if len(record.Content) == 0 {
		t.Error("Content 字段不应为空")
	}
}

// TestCreateRecordWithContent 验证带 content 的创建正常
func TestCreateRecordWithContent(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("rcw"), "有内容记录测试", 1)
	family := createTestFamily(t, admin.ID, "有内容记录家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "内容宝宝",
		Gender:    2,
		Birthday:  time.Now().AddDate(0, -6, 0),
	}
	mysql.DB.Create(&baby)

	body := mustMarshal(t, map[string]interface{}{
		"type":       "sleep",
		"baby_id":    baby.ID,
		"family_id":  family.ID,
		"started_at": time.Now().Format("2006-01-02T15:04:05Z"),
		"note":       "带 content 的测试记录",
		"content": map[string]interface{}{
			"duration": 120,
			"quality":  "good",
		},
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/records",
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	CreateRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("带 content 创建记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	var record model.Record
	mysql.DB.First(&record, "id = ?", resp.Data.ID)
	if string(record.Content) == "" {
		t.Error("Content 不应为空")
	}
}

// TestCreateRecordWithMultipleTimeFormats 验证多时间格式解析
func TestCreateRecordWithMultipleTimeFormats(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("rtf"), "时间格式测试", 1)
	family := createTestFamily(t, admin.ID, "时间格式家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "时间宝宝",
		Gender:    1,
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	mysql.DB.Create(&baby)

	tests := []struct {
		name      string
		timeStr   string
		expectOK  bool
	}{
		{"RFC3339", "2025-05-01T10:30:00Z", true},
		{"日期时间格式", "2025-05-01 10:30", true},
		{"纯日期", "2025-05-01", true},
		{"无效格式", "not-a-time", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := mustMarshal(t, map[string]interface{}{
				"type":       "feeding",
				"baby_id":    baby.ID,
				"family_id":  family.ID,
				"started_at": tc.timeStr,
			})

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/records",
				bytes.NewReader(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			CreateRecord(ctx)

			if tc.expectOK && recorder.Code != http.StatusOK {
				t.Errorf("格式 '%s' 应成功: status=%d body=%s", tc.timeStr, recorder.Code, recorder.Body.String())
			}
			if !tc.expectOK && recorder.Code == http.StatusOK {
				t.Errorf("格式 '%s' 应失败", tc.timeStr)
			}
		})
	}
}

// TestCreateRecordInvalidBaby 验证无效宝宝 ID 返回错误
func TestCreateRecordInvalidBaby(t *testing.T) {
	setupAdminTestDB(t)

	body := mustMarshal(t, map[string]interface{}{
		"type":       "feeding",
		"baby_id":    "non-existent-baby-id",
		"family_id":  "non-existent-family-id",
		"started_at": time.Now().Format("2006-01-02T15:04:05Z"),
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/records",
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	CreateRecord(ctx)

	if recorder.Code == http.StatusOK {
		t.Errorf("无效宝宝应返回错误: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
