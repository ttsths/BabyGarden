package e2e

import (
	"net/http"
	"testing"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
)

// ============================================================
// E2E Smoke Tests — 5 条关键路径冒烟用例
// 耗时 < 3 分钟，覆盖 admin 模块最可能坏的路径
// ============================================================

// Smoke 1: 管理员登录 + 获取仪表盘数据
func TestSmoke_AdminLoginAndDashboard(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanup := seedAdminUser(t)
	defer cleanup()

	// Step 1: 登录
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)

	assertStatus(t, loginW, http.StatusOK, "管理员登录")
	loginResp := parseResponse(t, loginW)
	assertCode(t, loginResp, model.SUCCESS, "管理员登录")

	// 提取 token
	loginData, ok := loginResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("登录响应 data 格式错误")
	}
	token, ok := loginData["token"].(string)
	if !ok || token == "" {
		t.Fatalf("登录未返回token, data=%+v", loginData)
	}

	// Step 2: 用 token 获取仪表盘
	statsW := getJSON(t, r, "/api/v1/admin/stats/overview", map[string]string{
		"Authorization": "Bearer " + token,
	})

	assertStatus(t, statsW, http.StatusOK, "仪表盘")
	statsResp := parseResponse(t, statsW)
	assertCode(t, statsResp, model.SUCCESS, "仪表盘")

	// Step 3: 验证 data 有核心字段
	statsData, ok := statsResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("仪表盘 data 格式错误")
	}
	if _, exists := statsData["users"]; !exists {
		t.Error("仪表盘缺少 users 字段")
	}
	if _, exists := statsData["families"]; !exists {
		t.Error("仪表盘缺少 families 字段")
	}
	if _, exists := statsData["babies"]; !exists {
		t.Error("仪表盘缺少 babies 字段")
	}
	t.Logf("仪表盘数据: users=%v, families=%v, babies=%v, photos=%v, records=%v",
		statsData["users"], statsData["families"], statsData["babies"],
		statsData["photos"], statsData["records"])
}

// Smoke 2: 管理员登录 + 查看用户列表（含分页）
func TestSmoke_AdminLoginAndUsersList(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanup := seedAdminUser(t)
	defer cleanup()

	// Step 1: 登录
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginResp := parseResponse(t, loginW)
	loginData := loginResp.Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Step 2: 获取用户列表（分页）
	usersW := getJSON(t, r, "/api/v1/admin/users?page=1&size=10", map[string]string{
		"Authorization": "Bearer " + token,
	})

	assertStatus(t, usersW, http.StatusOK, "用户列表")
	usersResp := parseResponse(t, usersW)
	assertCode(t, usersResp, model.SUCCESS, "用户列表")

	// Step 3: 验证分页结构（total 在 pagination 对象内）
	usersData, ok := usersResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("用户列表 data 格式错误")
	}
	if list, ok := usersData["list"]; !ok || list == nil {
		t.Error("用户列表缺少 list 字段")
	}
	pagination, ok := usersData["pagination"].(map[string]interface{})
	if !ok {
		t.Error("用户列表缺少 pagination 字段")
	} else {
		if _, ok := pagination["total"]; !ok {
			t.Error("用户列表 pagination 缺少 total 字段")
		}
		if _, ok := pagination["page"]; !ok {
			t.Error("用户列表 pagination 缺少 page 字段")
		}
	}
	listLen := 0
	if l, ok := usersData["list"].([]interface{}); ok {
		listLen = len(l)
	}
	t.Logf("用户列表: pagination=%+v, list_len=%d", pagination, listLen)
}

// Smoke 3: 无 token 访问 admin 接口 → 401
func TestSmoke_NoTokenReturns401(t *testing.T) {
	r := e2eRouter(t)

	endpoints := []string{
		"/api/v1/admin/stats/overview",
		"/api/v1/admin/stats/daily?days=7",
		"/api/v1/admin/users",
		"/api/v1/admin/families",
		"/api/v1/admin/babies",
		"/api/v1/admin/photos",
		"/api/v1/admin/records",
	}

	for _, path := range endpoints {
		t.Run(path, func(t *testing.T) {
			w := getJSON(t, r, path, nil)
			assertStatus(t, w, http.StatusUnauthorized, "无token: "+path)

			resp := parseResponse(t, w)
			if resp.Code == model.SUCCESS {
				t.Errorf("%s: 无token应返回错误, 实际 code=%d", path, resp.Code)
			}
		})
	}

	// 验证 login 不需要 token（不放 JWT 中间件）
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    "13800000000",
		"password": "wrong",
	}, nil)
	if loginW.Code == http.StatusUnauthorized {
		t.Fatal("admin/login 不应该要求 JWT token")
	}
}

// Smoke 4: 普通用户 token 访问 admin 接口 → 403
func TestSmoke_NonAdminTokenReturns403(t *testing.T) {
	r := e2eRouter(t)
	user, cleanupUser := seedRegularUser(t)
	defer cleanupUser()

	// 生成普通用户的 JWT
	accessToken, _, err := middleware.GenerateTokenPair(user.ID, user.Phone)
	if err != nil {
		t.Fatalf("生成普通用户token失败: %v", err)
	}

	adminEndpoints := []string{
		"/api/v1/admin/stats/overview",
		"/api/v1/admin/stats/daily?days=7",
		"/api/v1/admin/users",
		"/api/v1/admin/families",
	}

	for _, path := range adminEndpoints {
		t.Run(path, func(t *testing.T) {
			w := getJSON(t, r, path, map[string]string{
				"Authorization": "Bearer " + accessToken,
			})
			assertStatus(t, w, http.StatusForbidden, "普通用户token: "+path)

			resp := parseResponse(t, w)
			if resp.Code == model.SUCCESS {
				t.Errorf("%s: 普通用户应被拒绝, 实际 code=%d", path, resp.Code)
			}
		})
	}
}

// Smoke 5: 管理员登录 + 创建家庭 + 获取家庭详情（自给自足，不依赖预存数据）
func TestSmoke_AdminLoginAndFamilyDetail(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Step 1: 登录
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginData := parseResponse(t, loginW).Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Step 2: 创建测试家庭（确保数据自给自足）
	family, cleanupFamily := seedFamily(t, admin.ID)
	defer cleanupFamily()

	// Step 3: 获取家庭列表，验证新创建的家庭存在
	familiesW := getJSON(t, r, "/api/v1/admin/families?page=1&size=10", map[string]string{
		"Authorization": "Bearer " + token,
	})

	assertStatus(t, familiesW, http.StatusOK, "家庭列表")
	familiesResp := parseResponse(t, familiesW)
	assertCode(t, familiesResp, model.SUCCESS, "家庭列表")

	// Step 4: 获取新创建家庭的详情（用我们自己的 family ID）
	detailW := getJSON(t, r, "/api/v1/admin/families/"+family.ID, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, detailW, http.StatusOK, "家庭详情")
	detailResp := parseResponse(t, detailW)
	assertCode(t, detailResp, model.SUCCESS, "家庭详情")

	t.Logf("家庭详情: id=%s, name=%s", family.ID, family.Name)
}

// Smoke 6: Record CREATE → LIST → DELETE full lifecycle
func TestSmoke_RecordCRUD(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Login
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginData := parseResponse(t, loginW).Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Create family + baby for record context
	family, cleanupFamily := seedFamily(t, admin.ID)
	defer cleanupFamily()
	baby, cleanupBaby := seedBaby(t, family.ID)
	defer cleanupBaby()

	// Step 1: Create record
	createBody := map[string]interface{}{
		"type":       "feeding",
		"baby_id":    baby.ID,
		"family_id":  family.ID,
		"started_at": "2026-05-09 12:00",
		"note":       "e2e record test",
	}
	createW := postJSON(t, r, "/api/v1/admin/records", createBody, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, createW, http.StatusOK, "创建记录")
	createResp := parseResponse(t, createW)
	assertCode(t, createResp, model.SUCCESS, "创建记录")

	recordData, ok := createResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("创建记录响应 data 格式错误")
	}
	recordID, ok := recordData["id"].(string)
	if !ok || recordID == "" {
		t.Fatalf("创建记录未返回 id: %+v", recordData)
	}

	// Verify type is preserved
	if typ, ok := recordData["type"].(string); !ok || typ != "feeding" {
		t.Errorf("记录 type 期望 feeding, 实际: %v", recordData["type"])
	}

	// Step 2: List records and verify our record exists
	listW := getJSON(t, r, "/api/v1/admin/records?baby_id="+baby.ID+"&page=1&size=10", map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, listW, http.StatusOK, "记录列表")
	listResp := parseResponse(t, listW)
	assertCode(t, listResp, model.SUCCESS, "记录列表")

	listData, ok := listResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("记录列表 data 格式错误")
	}
	list, ok := listData["list"].([]interface{})
	if !ok || len(list) == 0 {
		t.Fatal("记录列表为空")
	}
	t.Logf("记录列表长度: %d", len(list))

	// Step 3: Delete the record
	deleteW := doRequest(t, r, http.MethodDelete, "/api/v1/admin/records/"+recordID, "", map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, deleteW, http.StatusOK, "删除记录")
	deleteResp := parseResponse(t, deleteW)
	assertCode(t, deleteResp, model.SUCCESS, "删除记录")

	t.Logf("记录 CRUD 完整流程通过: id=%s", recordID)
}

// Smoke 7: Baby list returns gender string and family_name (Bug #3/#4 fix verification)
func TestSmoke_BabyListGenderAndFamilyName(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Login
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginData := parseResponse(t, loginW).Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Create family + baby
	family, cleanupFamily := seedFamily(t, admin.ID)
	defer cleanupFamily()
	baby, cleanupBaby := seedBabyWithGender(t, family.ID, "male")
	defer cleanupBaby()

	// Get baby list
	listW := getJSON(t, r, "/api/v1/admin/babies?page=1&size=10", map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, listW, http.StatusOK, "宝宝列表")
	listResp := parseResponse(t, listW)
	assertCode(t, listResp, model.SUCCESS, "宝宝列表")

	listData := listResp.Data.(map[string]interface{})
	babies, ok := listData["list"].([]interface{})
	if !ok || len(babies) == 0 {
		t.Fatal("宝宝列表为空")
	}

	// Find our baby and verify fields
	found := false
	for _, b := range babies {
		bm, ok := b.(map[string]interface{})
		if !ok {
			continue
		}
		if bm["id"] == baby.ID {
			found = true
			// Verify gender is string not int
			gender, ok := bm["gender"].(string)
			if !ok {
				t.Errorf("gender 应为字符串, 实际类型=%T 值=%v", bm["gender"], bm["gender"])
			} else if gender != "male" {
				t.Errorf("gender 期望 male, 实际=%s", gender)
			}
			// Verify family_name is present
			familyName, ok := bm["family_name"].(string)
			if !ok || familyName == "" {
				t.Errorf("family_name 缺失或为空: %v", bm["family_name"])
			}
			t.Logf("宝宝验证: gender=%s family_name=%s", gender, familyName)
			break
		}
	}
	if !found {
		t.Error("未在列表中找到创建的宝宝")
	}
}

// Smoke 8: User list returns is_admin field (Bug #8 fix verification)
func TestSmoke_UserListHasIsAdmin(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Login
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginData := parseResponse(t, loginW).Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Get user list
	listW := getJSON(t, r, "/api/v1/admin/users?page=1&size=10", map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, listW, http.StatusOK, "用户列表")
	listResp := parseResponse(t, listW)
	assertCode(t, listResp, model.SUCCESS, "用户列表")

	listData := listResp.Data.(map[string]interface{})
	users, ok := listData["list"].([]interface{})
	if !ok {
		t.Fatal("用户列表 list 字段缺失")
	}

	// Verify each user has is_admin field
	for _, u := range users {
		um, ok := u.(map[string]interface{})
		if !ok {
			continue
		}
		isAdmin, exists := um["is_admin"]
		if !exists {
			t.Error("用户列表缺少 is_admin 字段")
			break
		}
		// is_admin can be float64 (JSON number) or int
		switch v := isAdmin.(float64); {
		case v != 1 && v != 0:
			// Check if it's actually a different type
			if ia, ok := isAdmin.(int); !ok || (ia != 0 && ia != 1) {
				t.Logf("is_admin 值: %v (type=%T)", isAdmin, isAdmin)
			}
		}
	}
	t.Logf("用户列表验证通过: 共 %d 个用户", len(users))
}

// Smoke 9: Wrong password returns 403
func TestSmoke_WrongPasswordReturns403(t *testing.T) {
	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Attempt login with wrong password
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "wrongpassword",
	}, nil)
	assertStatus(t, loginW, http.StatusForbidden, "错误密码")
	loginResp := parseResponse(t, loginW)
	if loginResp.Code == model.SUCCESS {
		t.Error("错误密码不应登录成功")
	}
}

// Smoke 10: Photo upload URL (known R2 dependency — skip until R2 configured)
func TestSmoke_PhotoUploadURL(t *testing.T) {
	t.Skip("跳过：需要 R2 环境变量配置（Bug #1）")

	r := e2eRouter(t)
	admin, cleanupAdmin := seedAdminUser(t)
	defer cleanupAdmin()

	// Login
	loginW := postJSON(t, r, "/api/v1/admin/login", map[string]string{
		"phone":    admin.Phone,
		"password": "e2eadmin",
	}, nil)
	loginData := parseResponse(t, loginW).Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Create family + baby
	family, cleanupFamily := seedFamily(t, admin.ID)
	defer cleanupFamily()
	baby, cleanupBaby := seedBaby(t, family.ID)
	defer cleanupBaby()

	// Request upload URL
	uploadW := postJSON(t, r, "/api/v1/admin/photos/upload-url", map[string]interface{}{
		"baby_id":      baby.ID,
		"filename":     "test.jpg",
		"content_type": "image/jpeg",
		"size":         12345,
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assertStatus(t, uploadW, http.StatusOK, "照片上传URL")
	uploadResp := parseResponse(t, uploadW)
	assertCode(t, uploadResp, model.SUCCESS, "照片上传URL")

	// Verify upload URL fields
	data, ok := uploadResp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("照片上传URL data 格式错误")
	}
	if url, ok := data["upload_url"].(string); !ok || url == "" {
		t.Error("upload_url 缺失或为空")
	}
	if url, ok := data["access_url"].(string); !ok || url == "" {
		t.Error("access_url 缺失或为空")
	}
	t.Logf("照片上传URL: upload_url=%s", data["upload_url"])
}
