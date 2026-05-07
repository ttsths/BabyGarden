package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
)

func TestGetDailyStats(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("stats"), "统计用户")
	family := createTestFamily(t, admin.ID, "统计家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "统计宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	base := time.Now().In(time.Local)
	startedAt := time.Date(base.Year(), base.Month(), base.Day(), 10, 0, 0, 0, time.Local)
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: startedAt,
		Content:   mustJSON(t, map[string]interface{}{"type": "formula", "amount": 120, "unit": "ml"}),
		CreatedBy: admin.ID,
	})
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: startedAt.Add(-2 * time.Hour),
		EndedAt:   sqlNullTime(startedAt),
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		CreatedBy: admin.ID,
	})
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeDiaper,
		StartedAt: startedAt.Add(-1 * time.Hour),
		Content:   mustJSON(t, map[string]interface{}{"type": "wet"}),
		CreatedBy: admin.ID,
	})

	dateStr := startedAt.Format("2006-01-02")
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/stats/daily?baby_id="+baby.ID+"&date="+dateStr, nil)
	ctx.Set("userId", admin.ID)

	GetDailyStats(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取日统计失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int                `json:"code"`
		Data DailyStatsResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.Feeding.Count != 1 || response.Data.Feeding.TotalAmount != 120 {
		t.Fatalf("喂养统计错误: %+v", response.Data.Feeding)
	}
	if response.Data.Sleep.Count != 1 || response.Data.Sleep.TotalHours < 1.9 {
		t.Fatalf("睡眠统计错误: %+v", response.Data.Sleep)
	}
	if response.Data.Diaper.Count != 1 {
		t.Fatalf("排泄统计错误: %+v", response.Data.Diaper)
	}
}

func TestGetWeeklyStats(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("weekly"), "周统计用户")
	family := createTestFamily(t, admin.ID, "周统计家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "周统计宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	base := time.Now().In(time.Local)
	day0 := time.Date(base.Year(), base.Month(), base.Day(), 9, 0, 0, 0, time.Local)
	day1 := day0.AddDate(0, 0, -1)

	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: day0,
		Content:   mustJSON(t, map[string]interface{}{"type": "formula", "amount": 100}),
		CreatedBy: admin.ID,
	})
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: day1.Add(-2 * time.Hour),
		EndedAt:   sqlNullTime(day1),
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		CreatedBy: admin.ID,
	})
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeDiaper,
		StartedAt: day1,
		Content:   mustJSON(t, map[string]interface{}{"type": "wet"}),
		CreatedBy: admin.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/stats/weekly?baby_id="+baby.ID, nil)
	ctx.Set("userId", admin.ID)

	GetWeeklyStats(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取周统计失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int                 `json:"code"`
		Data WeeklyStatsResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)

	if len(response.Data.Dates) != 7 {
		t.Fatalf("日期数量错误: %+v", response.Data.Dates)
	}

	idxToday := 6
	idxYesterday := 5
	if response.Data.Feeding[idxToday] != 1 {
		t.Fatalf("今日喂养统计错误: %+v", response.Data.Feeding)
	}
	if response.Data.Diaper[idxYesterday] != 1 {
		t.Fatalf("昨日排泄统计错误: %+v", response.Data.Diaper)
	}
	if response.Data.Sleep[idxYesterday] < 1.9 {
		t.Fatalf("昨日日睡眠统计错误: %+v", response.Data.Sleep)
	}
}
