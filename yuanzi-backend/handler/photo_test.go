package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/config"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

func TestGetPhotoUploadURLCreatesPending(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("photo"), "照片用户")
	family := createTestFamily(t, admin.ID, "照片家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "照片宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	body := mustMarshal(t, PhotoUploadURLRequest{
		BabyID:      baby.ID,
		Filename:    "test.jpg",
		ContentType: "image/jpeg",
		Size:        2048,
	})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photo/upload-url", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", admin.ID)

	GetPhotoUploadURL(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取上传地址失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int                    `json:"code"`
		Data PhotoUploadURLResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.PhotoID == "" || response.Data.UploadURL == "" {
		t.Fatalf("返回数据不完整: %+v", response.Data)
	}

	var photo model.Photo
	if err := mysql.DB.Where("id = ?", response.Data.PhotoID).First(&photo).Error; err != nil {
		t.Fatalf("照片记录未创建: %v", err)
	}
	if photo.Status != model.PhotoStatusPending {
		t.Fatalf("照片状态错误: %s", photo.Status)
	}
}

func TestPhotoUploadCallbackActivatesPhoto(t *testing.T) {
	setupFamilyTestDB(t)

	prevCallbackSecret := config.GlobalConfig.OSS.CallbackSecret
	config.GlobalConfig.OSS.CallbackSecret = "test-secret"
	defer func() { config.GlobalConfig.OSS.CallbackSecret = prevCallbackSecret }()

	admin := createTestUser(t, uniquePhone("photocb"), "回调用户")
	family := createTestFamily(t, admin.ID, "回调家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "回调宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    family.ID,
		OSSKey:      "families/test/photo.jpg",
		Size:        1,
		ContentType: "image/jpeg",
		UploadedBy:  admin.ID,
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusPending,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}

	body, _ := json.Marshal(PhotoCallbackRequest{PhotoID: photo.ID, Size: 4096, ETag: "etag"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photo/callback", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("X-OSS-Callback-Token", "test-secret")

	PhotoUploadCallback(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("回调处理失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var updated model.Photo
	if err := mysql.DB.First(&updated, "id = ?", photo.ID).Error; err != nil {
		t.Fatalf("查询照片失败: %v", err)
	}
	if updated.Status != model.PhotoStatusActive || updated.Size != 4096 {
		t.Fatalf("照片状态未更新: %+v", updated)
	}
}

func TestListPhotos(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("plist"), "列表用户")
	family := createTestFamily(t, admin.ID, "列表家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "列表宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    family.ID,
		OSSKey:      "families/test/list.jpg",
		Size:        100,
		ContentType: "image/jpeg",
		UploadedBy:  admin.ID,
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusActive,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/photo?baby_id="+baby.ID, nil)
	ctx.Set("userId", admin.ID)

	ListPhotos(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取照片列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestListPhotosDateRange(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("pdate"), "日期用户")
	family := createTestFamily(t, admin.ID, "日期家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "日期宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	today := time.Now().In(time.Local)
	oldDay := today.AddDate(0, 0, -2)
	photoToday := model.Photo{BabyID: baby.ID, FamilyID: family.ID, OSSKey: "families/test/today.jpg", Size: 100, ContentType: "image/jpeg", UploadedBy: admin.ID, UploadedAt: today, Status: model.PhotoStatusActive}
	photoOld := model.Photo{BabyID: baby.ID, FamilyID: family.ID, OSSKey: "families/test/old.jpg", Size: 100, ContentType: "image/jpeg", UploadedBy: admin.ID, UploadedAt: oldDay, Status: model.PhotoStatusActive}
	if err := mysql.DB.Create(&photoToday).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}
	if err := mysql.DB.Create(&photoOld).Error; err != nil {
		t.Fatalf("创建旧照片失败: %v", err)
	}

	dateStr := today.Format("2006-01-02")
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/photo?baby_id="+baby.ID+"&date_from="+dateStr+"&date_to="+dateStr, nil)
	ctx.Set("userId", admin.ID)

	ListPhotos(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取照片列表失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			List []PhotoResponse `json:"list"`
		} `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if len(response.Data.List) != 1 {
		t.Fatalf("日期筛选数量错误: %+v", response.Data.List)
	}
}

func TestPhotoCommentsAndLikes(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("pinter"), "互动用户")
	family := createTestFamily(t, admin.ID, "互动家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "互动宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	photo := model.Photo{BabyID: baby.ID, FamilyID: family.ID, OSSKey: "families/test/interact.jpg", Size: 100, ContentType: "image/jpeg", UploadedBy: admin.ID, UploadedAt: time.Now(), Status: model.PhotoStatusActive}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}

	commentBody := mustMarshal(t, CreatePhotoCommentRequest{Content: "真可爱"})
	commentRecorder := httptest.NewRecorder()
	commentCtx, _ := gin.CreateTestContext(commentRecorder)
	commentCtx.Request = httptest.NewRequest(http.MethodPost, "/photo/"+photo.ID+"/comments", bytes.NewReader(commentBody))
	commentCtx.Request.Header.Set("Content-Type", "application/json")
	commentCtx.Params = gin.Params{{Key: "id", Value: photo.ID}}
	commentCtx.Set("userId", admin.ID)
	CreatePhotoComment(commentCtx)
	if commentRecorder.Code != http.StatusOK {
		t.Fatalf("评论失败: status=%d body=%s", commentRecorder.Code, commentRecorder.Body.String())
	}

	likeRecorder := httptest.NewRecorder()
	likeCtx, _ := gin.CreateTestContext(likeRecorder)
	likeCtx.Request = httptest.NewRequest(http.MethodPost, "/photo/"+photo.ID+"/like", nil)
	likeCtx.Params = gin.Params{{Key: "id", Value: photo.ID}}
	likeCtx.Set("userId", admin.ID)
	LikePhoto(likeCtx)
	if likeRecorder.Code != http.StatusOK {
		t.Fatalf("点赞失败: status=%d body=%s", likeRecorder.Code, likeRecorder.Body.String())
	}

	listRecorder := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listRecorder)
	listCtx.Request = httptest.NewRequest(http.MethodGet, "/photo?baby_id="+baby.ID, nil)
	listCtx.Set("userId", admin.ID)
	ListPhotos(listCtx)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("获取照片列表失败: status=%d body=%s", listRecorder.Code, listRecorder.Body.String())
	}
	var response struct {
		Data struct {
			List []PhotoResponse `json:"list"`
		} `json:"data"`
	}
	decodeResponse(t, listRecorder.Body.Bytes(), &response)
	if len(response.Data.List) != 1 || response.Data.List[0].LikeCount != 1 || response.Data.List[0].CommentCount != 1 || !response.Data.List[0].LikedByMe {
		t.Fatalf("互动统计错误: %+v", response.Data.List)
	}
}

func TestDeletePhotoByUploader(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("pdel"), "删除用户")
	family := createTestFamily(t, admin.ID, "删除家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "删除宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    family.ID,
		OSSKey:      "families/test/delete.jpg",
		Size:        100,
		ContentType: "image/jpeg",
		UploadedBy:  admin.ID,
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusActive,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/photo/"+photo.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: photo.ID}}
	ctx.Set("userId", admin.ID)

	DeletePhoto(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("删除照片失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var updated model.Photo
	if err := mysql.DB.First(&updated, "id = ?", photo.ID).Error; err != nil {
		t.Fatalf("查询照片失败: %v", err)
	}
	if updated.Status != model.PhotoStatusDeleted {
		t.Fatalf("照片状态未更新: %+v", updated)
	}
}

func TestPhotoUploadCallbackRejectsMissingToken(t *testing.T) {
	setupFamilyTestDB(t)

	prevCallbackSecret := config.GlobalConfig.OSS.CallbackSecret
	config.GlobalConfig.OSS.CallbackSecret = "test-secret"
	defer func() { config.GlobalConfig.OSS.CallbackSecret = prevCallbackSecret }()

	admin := createTestUser(t, uniquePhone("cbtok"), "回调鉴权")
	family := createTestFamily(t, admin.ID, "回调鉴权家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "回调鉴权宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)

	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    family.ID,
		OSSKey:      "families/test/callback.jpg",
		Size:        1,
		ContentType: "image/jpeg",
		UploadedBy:  admin.ID,
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusPending,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		t.Fatalf("创建照片失败: %v", err)
	}

	body, _ := json.Marshal(PhotoCallbackRequest{PhotoID: photo.ID, Size: 4096, ETag: "etag"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photo/callback", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	PhotoUploadCallback(ctx)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("缺少回调鉴权应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
