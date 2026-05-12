package admin

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

// TestGetPhotosReturnsFilenameAndOriginalURL 验证 Issue #54：GetPhotos 返回 filename 和 original_url
func TestGetPhotosReturnsFilenameAndOriginalURL(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("photofn"), "照片测试管理员", 1)
	family := createTestFamily(t, admin.ID, "照片下载测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	// 插入测试照片
	now := time.Now()
	photo := model.Photo{
		BabyID:      "test-baby-photo-id",
		FamilyID:    family.ID,
		OSSKey:      family.ID + "/test-baby-photo-id/cat.jpg",
		Size:        12345,
		ContentType: "image/jpeg",
		Status:      model.PhotoStatusActive,
		UploadedBy:  admin.ID,
		UploadedAt:  now,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建测试照片失败: %v", err)
	}
	defer mysql.DB.Where("id = ?", photo.ID).Delete(&model.Photo{})

	// 设置 R2_PUBLIC_URL 环境变量
	os.Setenv("R2_PUBLIC_URL", "https://pub-test.r2.dev")
	defer os.Unsetenv("R2_PUBLIC_URL")

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/photos?page=1&page_size=20", nil)

	GetPhotos(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取照片列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
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

	// 验证列表非空
	if len(resp.Data.List) == 0 {
		t.Fatal("照片列表为空，无法验证字段")
	}

	// 查找自己插入的测试照片（列表可能包含其他测试的残留数据）
	expectedURL := "https://pub-test.r2.dev/" + photo.OSSKey
	found := false
	for _, item := range resp.Data.List {
		filename, _ := item["filename"].(string)
		originalURL, _ := item["original_url"].(string)
		if originalURL == expectedURL {
			found = true
			if filename != "cat.jpg" {
				t.Errorf("filename 期望 cat.jpg, 实际 %s", filename)
			}
			break
		}
	}
	if !found {
		t.Errorf("未找到测试照片 (期望 original_url=%s)，响应列表: %+v", expectedURL, resp.Data.List)
	}
}

// TestGetPhotosFilenameWithoutSlash 验证边界：OSSKey 不含 "/" 时 filename = whole key
func TestGetPhotosFilenameWithoutSlash(t *testing.T) {
	setupAdminTestDB(t)

	admin := createTestUser(t, uniquePhone("phonos"), "无分隔线照片测试", 1)
	family := createTestFamily(t, admin.ID, "无分隔线测试家庭")
	defer cleanupBabyRecords(t, family.ID)
	defer cleanupUsers(t, admin.ID)

	now := time.Now()
	photo := model.Photo{
		BabyID:      "test-baby-noslash",
		FamilyID:    family.ID,
		OSSKey:      "singlefile.jpg",
		Size:        999,
		ContentType: "image/jpeg",
		Status:      model.PhotoStatusActive,
		UploadedBy:  admin.ID,
		UploadedAt:  now,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建测试照片失败: %v", err)
	}
	defer mysql.DB.Where("id = ?", photo.ID).Delete(&model.Photo{})

	os.Setenv("R2_PUBLIC_URL", "https://pub-test.r2.dev")
	defer os.Unsetenv("R2_PUBLIC_URL")

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/photos?page=1&page_size=20", nil)

	GetPhotos(ctx)

	var resp struct {
		Code int `json:"code"`
		Data struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if len(resp.Data.List) == 0 {
		t.Fatal("照片列表为空")
	}

	item := resp.Data.List[0]

	// OSSKey 不含 "/" → filename = 整个 key
	if item["filename"] != "singlefile.jpg" {
		t.Errorf("filename 期望 singlefile.jpg, 实际 %v", item["filename"])
	}
}
