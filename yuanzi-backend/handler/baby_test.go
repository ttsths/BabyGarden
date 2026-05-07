package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
)

func TestDeleteBabyByAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("badmin"), "宝宝管理员")
	memberUser := createTestUser(t, uniquePhone("bmember"), "普通成员")
	family := createTestFamily(t, admin.ID, "宝宝删除家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	baby := createTestBaby(t, family.ID, "删除宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	record := createTestRecord(t, model.Record{
		BabyID:    baby.ID,
		FamilyID:  family.ID,
		Type:      model.RecordTypeFeeding,
		StartedAt: time.Now().UTC(),
		Content:   mustJSON(t, map[string]interface{}{"type": "breast", "duration": 6}),
		CreatedBy: admin.ID,
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/baby/"+baby.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: baby.ID}}
	ctx.Set("userId", admin.ID)

	DeleteBaby(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("管理员删除宝宝失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var deletedBaby model.Baby
	if err := mysql.DB.Unscoped().First(&deletedBaby, "id = ?", baby.ID).Error; err == nil {
		t.Fatalf("宝宝未被删除")
	}
	var deletedRecord model.Record
	if err := mysql.DB.Unscoped().First(&deletedRecord, "id = ?", record.ID).Error; err != nil {
		t.Fatalf("查询删除记录失败: %v", err)
	}
	if !deletedRecord.DeletedAt.Valid {
		t.Fatalf("记录未被软删除")
	}
}

func TestDeleteBabyRejectNonAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("bowner"), "宝宝管理员2")
	memberUser := createTestUser(t, uniquePhone("buser"), "普通成员2")
	family := createTestFamily(t, admin.ID, "宝宝权限家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	baby := createTestBaby(t, family.ID, "权限宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/baby/"+baby.ID, nil)
	ctx.Params = gin.Params{{Key: "id", Value: baby.ID}}
	ctx.Set("userId", memberUser.ID)

	DeleteBaby(ctx)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("非管理员删除宝宝应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
