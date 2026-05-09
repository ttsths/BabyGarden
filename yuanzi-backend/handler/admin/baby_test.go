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

// TestGetBabiesReturnsGenderString 验证 Bug #3 修复：性别返回字符串而非 int8
func TestGetBabiesReturnsGenderString(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bgen"), "性别测试管理员", 1)
	family := createTestFamily(t, admin.ID, "性别测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	// 创建不同性别的宝宝
	baby1 := model.Baby{
		FamilyID:  family.ID,
		Name:      "男宝宝",
		Gender:    1, // male
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	baby2 := model.Baby{
		FamilyID:  family.ID,
		Name:      "女宝宝",
		Gender:    2, // female
		Birthday:  time.Now().AddDate(0, -6, 0),
	}
	baby3 := model.Baby{
		FamilyID:  family.ID,
		Name:      "未知性别宝宝",
		Gender:    0, // unknown
		Birthday:  time.Now(),
	}
	mysql.DB.Create(&baby1)
	mysql.DB.Create(&baby2)
	mysql.DB.Create(&baby3)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/babies?page=1&page_size=20", nil)

	GetBabies(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取宝宝列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Code int `json:"code"`
		Data struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Code != model.SUCCESS {
		t.Fatalf("业务返回失败: %+v", resp)
	}

	genderMap := make(map[string]string)
	for _, item := range resp.Data.List {
		name := item["name"].(string)
		gender := item["gender"].(string)
		genderMap[name] = gender
	}

	// 验证性别是字符串 not 数字
	if genderMap["男宝宝"] != "male" {
		t.Errorf("男宝宝性别期望 male, 实际 %s", genderMap["男宝宝"])
	}
	if genderMap["女宝宝"] != "female" {
		t.Errorf("女宝宝性别期望 female, 实际 %s", genderMap["女宝宝"])
	}
	if genderMap["未知性别宝宝"] != "unknown" {
		t.Errorf("未知性别宝宝期望 unknown, 实际 %s", genderMap["未知性别宝宝"])
	}
}

// TestGetBabiesReturnsFamilyName 验证 Bug #4 修复：列表返回 family_name
func TestGetBabiesReturnsFamilyName(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bfn"), "家庭名测试管理员", 1)
	family := createTestFamily(t, admin.ID, "小园子测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "家庭名测试宝宝",
		Gender:    2,
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	mysql.DB.Create(&baby)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/babies", nil)

	GetBabies(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取宝宝列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Data struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if len(resp.Data.List) == 0 {
		t.Fatal("宝宝列表为空")
	}

	// 在列表中查找自己的宝宝（不假设它是第一条）
	var item map[string]interface{}
	found := false
	for _, bm := range resp.Data.List {
		if bm["id"] == baby.ID {
			found = true
			item = bm
			break
		}
	}
	if !found {
		t.Fatalf("未在列表中找到宝宝 ID=%s, 列表中共%d条", baby.ID, len(resp.Data.List))
	}

	// 验证返回了 family_name
	fn, ok := item["family_name"].(string)
	if !ok || fn == "" {
		t.Errorf("family_name 缺失或为空: %v", item)
	}
	if fn != "小园子测试家庭" {
		t.Errorf("family_name 期望 '小园子测试家庭', 实际 '%s'", fn)
	}

	// 验证也返回了 family_id
	if _, ok := item["family_id"].(string); !ok {
		t.Errorf("family_id 缺失: %v", item)
	}
}

// TestGetBabyReturnsStringGenderAndFamilyName 验证 Bug #3+#4：详情返回字符串性别+家庭名
func TestGetBabyReturnsStringGenderAndFamilyName(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bdet"), "详情测试管理员", 1)
	family := createTestFamily(t, admin.ID, "详情测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "详情宝宝",
		Gender:    1, // male
		Birthday:  time.Now().AddDate(-2, 0, 0),
	}
	mysql.DB.Create(&baby)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/babies/"+baby.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: baby.ID}}

	GetBaby(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取宝宝详情失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	// 验证 gender 是字符串
	gender, ok := resp.Data["gender"].(string)
	if !ok {
		t.Fatalf("gender 应为字符串: type=%T value=%v", resp.Data["gender"], resp.Data["gender"])
	}
	if gender != "male" {
		t.Errorf("gender 期望 male, 实际 %s", gender)
	}

	// 验证 family_name 存在
	fn, ok := resp.Data["family_name"].(string)
	if !ok || fn == "" {
		t.Errorf("family_name 缺失或为空: %v", resp.Data)
	}
	if fn != "详情测试家庭" {
		t.Errorf("family_name 期望 '详情测试家庭', 实际 '%s'", fn)
	}
}

// TestUpdateBabyGenderWithString 验证 Bug #3 修复：用字符串更新性别
func TestUpdateBabyGenderWithString(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bupd"), "更新性别管理员", 1)
	family := createTestFamily(t, admin.ID, "更新性别家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "待更新宝宝",
		Gender:    0, // 初始 unknown
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	mysql.DB.Create(&baby)

	// 用字符串 "female" 更新
	body := mustMarshal(t, map[string]interface{}{
		"gender":   "female",
		"name":     "待更新宝宝",
		"birthday": time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/babies/"+baby.ID,
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: baby.ID}}

	UpdateBaby(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("更新宝宝性别失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	// 验证 DB 中正确存储
	var updated model.Baby
	mysql.DB.First(&updated, "id = ?", baby.ID)
	if updated.Gender != 2 {
		t.Errorf("DB 中 gender 应为 2 (female), 实际 %d", updated.Gender)
	}
}

// TestGenderToString 验证性别转换函数
func TestGenderToString(t *testing.T) {
	tests := []struct {
		input    int8
		expected string
	}{
		{0, "unknown"},
		{1, "male"},
		{2, "female"},
		{99, "unknown"}, // 未知值
		{-1, "unknown"},
	}

	for _, tc := range tests {
		got := genderToString(tc.input)
		if got != tc.expected {
			t.Errorf("genderToString(%d) = %s, want %s", tc.input, got, tc.expected)
		}
	}
}

// TestGetBabiesListBabyNotFound 验证不存在的宝宝
func TestGetBabiesListBabyNotFound(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bnf"), "不存在测试", 1)
	family := createTestFamily(t, admin.ID, "空家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	// 不创建宝宝，直接查列表
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/babies?keyword=不存在", nil)

	GetBabies(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("空列表查询失败: status=%d", recorder.Code)
	}

	var resp struct {
		Data struct {
			List  []interface{} `json:"list"`
			Total int64         `json:"total"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Data.Total > 0 {
		t.Errorf("搜索 '不存在' 应返回 0 条: got %d", resp.Data.Total)
	}
}

// TestUpdateBabyWithoutGender 验证不带 gender 的更新不报错
func TestUpdateBabyWithoutGender(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bng"), "无性别更新", 1)
	family := createTestFamily(t, admin.ID, "无性别更新家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	baby := model.Baby{
		FamilyID:  family.ID,
		Name:      "原名",
		Gender:    1,
		Birthday:  time.Now().AddDate(-1, 0, 0),
	}
	mysql.DB.Create(&baby)

	// 只更新名称，不带 gender
	body := mustMarshal(t, map[string]interface{}{
		"name":     "新名字",
		"birthday": time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/babies/"+baby.ID,
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: baby.ID}}

	UpdateBaby(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("不带 gender 更新失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var updated model.Baby
	mysql.DB.First(&updated, "id = ?", baby.ID)
	if updated.Name != "新名字" {
		t.Errorf("名称未更新: %s", updated.Name)
	}
	// gender 应保持不变
	if updated.Gender != 1 {
		t.Errorf("gender 被意外修改: %d", updated.Gender)
	}
}

// TestAdminCreateBabyWithStringGender 验证创建宝宝时 gender 字符串映射
func TestAdminCreateBabyWithStringGender(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("bcr"), "创建宝宝测试", 1)
	family := createTestFamily(t, admin.ID, "创建宝宝家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	body := mustMarshal(t, map[string]interface{}{
		"family_id": family.ID,
		"name":      "新创建宝宝",
		"gender":    "male",
		"birthday":  time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/babies",
		bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	CreateBaby(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建宝宝失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Data struct {
			ID     string `json:"id"`
			Gender string `json:"gender"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Data.Gender != "male" {
		t.Errorf("返回 gender 应为 male, 实际 %s", resp.Data.Gender)
	}

	// 验证 DB 存储
	var saved model.Baby
	mysql.DB.First(&saved, "id = ?", resp.Data.ID)
	if saved.Gender != 1 {
		t.Errorf("DB gender 应为 1 (male), 实际 %d", saved.Gender)
	}
}
