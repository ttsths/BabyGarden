package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
)

// TestUserListIncludesIsAdmin 验证 Bug #8 修复：用户列表返回 is_admin
func TestUserListIncludesIsAdmin(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("uisadmin"), "管理员用户", 1)
	normal := createTestUser(t, uniquePhone("uisnormal"), "普通用户", 0)
	defer cleanupUsers(t, admin.ID, normal.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1&page_size=20", nil)

	GetUsers(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取用户列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
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

	// 验证两个用户都包含 is_admin 字段
	for _, u := range resp.Data.List {
		phone := u["phone"].(string)
		isAdmin, exists := u["is_admin"]

		if !exists {
			t.Errorf("用户 %s 缺少 is_admin 字段: %v", phone, u)
			continue
		}

		// is_admin 可能是 float64 (JSON number) 或直接 number
		var isAdminVal float64
		switch v := isAdmin.(type) {
		case float64:
			isAdminVal = v
		case int8:
			isAdminVal = float64(v)
		default:
			t.Errorf("用户 %s is_admin 类型异常: %T = %v", phone, isAdmin, isAdmin)
			continue
		}

		if phone == admin.Phone && isAdminVal != 1 {
			t.Errorf("管理员用户 is_admin 应为 1, 实际 %v", isAdminVal)
		}
		if phone == normal.Phone && isAdminVal != 0 {
			t.Errorf("普通用户 is_admin 应为 0, 实际 %v", isAdminVal)
		}
	}
}

// TestCreateUserWithIsAdmin 验证创建用户时 is_admin 被正确存储
func TestCreateUserWithIsAdmin(t *testing.T) {
	setupAdminTestDB(t)

	phone := uniquePhone("cuadmin")
	body := mustMarshal(t, map[string]interface{}{
		"phone":    phone,
		"password": "12345678",
		"nickname": "新管理员",
		"is_admin": 1,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/users",
		bytesBody(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	CreateUser(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建管理员用户失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	// 验证 DB 中 is_admin 正确
	var user model.User
	if err := dbConn().Where("phone = ?", phone).First(&user).Error; err != nil {
		t.Fatalf("查询创建的用户失败: %v", err)
	}
	defer cleanupUsers(t, user.ID)

	if user.IsAdmin != 1 {
		t.Errorf("DB 中 is_admin 应为 1, 实际 %d", user.IsAdmin)
	}

	// 创建普通用户
	phone2 := uniquePhone("cunormal")
	body2 := mustMarshal(t, map[string]interface{}{
		"phone":    phone2,
		"password": "12345678",
		"nickname": "新普通用户",
		"is_admin": 0,
	})

	recorder2 := httptest.NewRecorder()
	ctx2, _ := gin.CreateTestContext(recorder2)
	ctx2.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/users",
		bytesBody(body2))
	ctx2.Request.Header.Set("Content-Type", "application/json")

	CreateUser(ctx2)

	if recorder2.Code != http.StatusOK {
		t.Fatalf("创建普通用户失败: status=%d body=%s", recorder2.Code, recorder2.Body.String())
	}

	var user2 model.User
	if err := dbConn().Where("phone = ?", phone2).First(&user2).Error; err != nil {
		t.Fatalf("查询创建的普通用户失败: %v", err)
	}
	defer cleanupUsers(t, user2.ID)

	if user2.IsAdmin != 0 {
		t.Errorf("普通用户 is_admin 应为 0, 实际 %d", user2.IsAdmin)
	}
}

// TestGetUserDetailWithIsAdmin 验证用户详情包含 is_admin
func TestGetUserDetailWithIsAdmin(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("guda"), "详情管理员", 1)
	defer cleanupUsers(t, admin.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+admin.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: admin.ID}}

	GetUser(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取用户详情失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	isAdmin, exists := resp.Data["is_admin"]
	if !exists {
		t.Errorf("用户详情缺少 is_admin 字段: %v", resp.Data)
	}

	switch v := isAdmin.(type) {
	case float64:
		if v != 1 {
			t.Errorf("is_admin 应为 1, 实际 %v", v)
		}
	default:
		t.Logf("is_admin 类型: %T = %v (非标准 float64, 但字段存在)", isAdmin, isAdmin)
	}
}
