package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

func TestCreateRecordFeedingHoursSinceLastFeed(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("feed"), "喂养用户")
	family := createTestFamily(t, admin.ID, "喂养家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "记录宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	startedAt := time.Now().UTC().Truncate(time.Second)
	previousAt := startedAt.Add(-2 * time.Hour)
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: previousAt,
		Content:   mustJSON(t, map[string]interface{}{"type": "breast", "side": "left", "duration": 10}),
		CreatedBy: admin.ID,
	})

	body := mustMarshal(t, CreateRecordRequest{
		BabyID:    baby.ID,
		Type:      "feeding",
		StartedAt: startedAt.Format(time.RFC3339),
		Content:   map[string]interface{}{"type": "breast", "side": "right", "duration": 12},
	})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/record", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", admin.ID)

	CreateRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建喂养记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int            `json:"code"`
		Data RecordResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Code != model.SUCCESS {
		t.Fatalf("业务返回失败: %+v", response)
	}
	if response.Data.HoursSinceLastFeed == nil || *response.Data.HoursSinceLastFeed != 2 {
		t.Fatalf("喂养间隔计算错误: %+v", response.Data)
	}
}

func TestCreateRecordTemperature(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("temp"), "测温用户")
	family := createTestFamily(t, admin.ID, "测温家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "测温宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	body := mustMarshal(t, CreateRecordRequest{
		BabyID:    baby.ID,
		Type:      "temperature",
		StartedAt: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
		Content:   map[string]interface{}{"value": 36.6, "unit": "C", "position": "armpit"},
	})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/record", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", admin.ID)

	CreateRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建测温记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data RecordResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.Type != "temperature" || response.Data.Content["value"] == nil {
		t.Fatalf("测温记录响应错误: %+v", response.Data)
	}
}

func TestUpdateRecordRejectNonCreator(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("owner"), "记录管理员")
	memberUser := createTestUser(t, uniquePhone("member"), "普通成员")
	family := createTestFamily(t, admin.ID, "更新记录家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	baby := createTestBaby(t, family.ID, "更新记录宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: time.Now().Add(-1 * time.Hour).UTC(),
		EndedAt:   sqlNullTime(time.Now().UTC()),
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		CreatedBy: admin.ID,
	})

	body := mustMarshal(t, UpdateRecordRequest{Note: stringPtr("修改备注")})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/record/"+record.ID, bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", memberUser.ID)

	UpdateRecord(ctx)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("非创建者修改应拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestDeleteRecordByAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("admin2"), "记录管理员")
	memberUser := createTestUser(t, uniquePhone("member2"), "普通成员")
	family := createTestFamily(t, admin.ID, "删除记录家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	baby := createTestBaby(t, family.ID, "删除记录宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeDiaper,
		StartedAt: time.Now().UTC(),
		Content:   mustJSON(t, map[string]interface{}{"type": "wet", "color": "yellow"}),
		CreatedBy: memberUser.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/record/"+record.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", admin.ID)

	DeleteRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("管理员删除失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var deleted model.Record
	if err := mysql.DB.Unscoped().First(&deleted, "id = ?", record.ID).Error; err != nil {
		t.Fatalf("查询删除记录失败: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatalf("记录未被软删除: %+v", deleted)
	}
}

func TestListRecordsFilterByDate(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("list"), "列表用户")
	family := createTestFamily(t, admin.ID, "记录列表家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "列表宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	day := time.Date(2024, 3, 9, 10, 0, 0, 0, time.Local)
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: day,
		Content:   mustJSON(t, map[string]interface{}{"type": "breast", "duration": 8}),
		CreatedBy: admin.ID,
	})
	createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: day.Add(24 * time.Hour),
		EndedAt:   sqlNullTime(day.Add(26 * time.Hour)),
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		CreatedBy: admin.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/record?baby_id="+baby.ID+"&type=feeding&date=2024-03-09", nil)
	ctx.Set("userId", admin.ID)

	ListRecords(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("查询记录列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			List []RecordResponse `json:"list"`
		} `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if len(response.Data.List) != 1 {
		t.Fatalf("日期过滤结果错误: %+v", response.Data.List)
	}
}

func createTestBaby(t *testing.T, familyID, name string) model.Baby {
	t.Helper()
	baby := model.Baby{
		FamilyID: familyID,
		Name:     name,
		Birthday: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		Gender:   1,
	}
	if err := mysql.DB.Create(&baby).Error; err != nil {
		t.Fatalf("创建测试宝宝失败: %v", err)
	}
	return baby
}

func createTestRecord(t *testing.T, record model.Record) model.Record {
	t.Helper()
	if err := mysql.DB.Create(&record).Error; err != nil {
		t.Fatalf("创建测试记录失败: %v", err)
	}
	return record
}

func mustJSON(t *testing.T, value interface{}) model.JSON {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("JSON 编码失败: %v", err)
	}
	return model.JSON(data)
}

func stringPtr(value string) *string { return &value }

func sqlNullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: true}
}

func TestGetRecordSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("getrecord"), "详情用户")
	family := createTestFamily(t, admin.ID, "详情家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "详情宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	startedAt := time.Now().Add(-2 * time.Hour).UTC()
	endedAt := startedAt.Add(30 * time.Minute)
	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: startedAt,
		EndedAt:   sqlNullTime(endedAt),
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		Note:      "测试备注",
		CreatedBy: admin.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/record/"+record.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", admin.ID)

	GetRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取记录详情失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int                  `json:"code"`
		Data RecordDetailResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.ID != record.ID || response.Data.BabyID != baby.ID {
		t.Fatalf("记录详情返回错误: %+v", response.Data)
	}
	if response.Data.CreatedBy != admin.ID {
		t.Fatalf("创建人返回错误: %+v", response.Data)
	}
	if response.Data.EndedAt == nil {
		t.Fatalf("结束时间缺失: %+v", response.Data)
	}
}

func TestUpdateRecordByCreatorSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("uprecord"), "更新用户")
	family := createTestFamily(t, admin.ID, "更新家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "更新宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeDiaper,
		StartedAt: time.Now().Add(-30 * time.Minute).UTC(),
		Content:   mustJSON(t, map[string]interface{}{"type": "wet", "color": "yellow"}),
		Note:      "原始备注",
		CreatedBy: admin.ID,
	})

	body := mustMarshal(t, UpdateRecordRequest{Note: stringPtr("更新备注"), Content: map[string]interface{}{"type": "dirty", "color": "brown"}})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/record/"+record.ID, bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", admin.ID)

	UpdateRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("更新记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var updated model.Record
	if err := mysql.DB.First(&updated, "id = ?", record.ID).Error; err != nil {
		t.Fatalf("查询更新记录失败: %v", err)
	}
	if updated.Note != "更新备注" {
		t.Fatalf("备注未更新: %+v", updated)
	}
}

func TestUpdateRecordRejectInvalidTime(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("uptime"), "时间用户")
	family := createTestFamily(t, admin.ID, "时间家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "时间宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	startedAt := time.Now().Add(-1 * time.Hour).UTC()
	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeSleep,
		StartedAt: startedAt,
		Content:   mustJSON(t, map[string]interface{}{"quality": "good", "location": "crib"}),
		CreatedBy: admin.ID,
	})

	invalidEnd := startedAt.Add(-10 * time.Minute).Format(time.RFC3339)
	body := mustMarshal(t, UpdateRecordRequest{EndedAt: &invalidEnd})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/record/"+record.ID, bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", admin.ID)

	UpdateRecord(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("结束时间早于开始时间应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestDeleteRecordByCreator(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("delrecord"), "删除用户")
	family := createTestFamily(t, admin.ID, "删除家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "删除宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: time.Now().UTC(),
		Content:   mustJSON(t, map[string]interface{}{"type": "formula", "amount": 90}),
		CreatedBy: admin.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/record/"+record.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: record.ID}}
	ctx.Set("userId", admin.ID)

	DeleteRecord(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建者删除记录失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var deleted model.Record
	if err := mysql.DB.Unscoped().First(&deleted, "id = ?", record.ID).Error; err != nil {
		t.Fatalf("查询删除记录失败: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatalf("记录未被软删除")
	}
}
