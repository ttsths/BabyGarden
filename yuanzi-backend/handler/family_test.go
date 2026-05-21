package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"
)

var familyTestSetupOnce sync.Once

func setupFamilyTestDB(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_EXTERNAL_INTEGRATION_TESTS") != "1" {
		t.Skip("跳过依赖 MySQL 的集成测试：设置 RUN_EXTERNAL_INTEGRATION_TESTS=1 后执行")
	}
	familyTestSetupOnce.Do(func() {
		projectRoot := getFamilyProjectRoot(t)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(filepath.Join(projectRoot, "config"))
		config.Setup()
		logger.Setup()
		mysql.Setup()
		gredis.Setup()
		if mysql.DB != nil {
			_ = mysql.DB.AutoMigrate(&model.PhotoComment{}, &model.PhotoLike{})
		}
		gin.SetMode(gin.TestMode)
	})
	if mysql.DB == nil || !mysql.IsConnected() {
		t.Skip("跳过依赖 MySQL 的集成测试：数据库未连接")
	}
}

func TestCreateFamilySuccess(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("create"), "创建家庭用户")
	defer cleanupUsers(t, user.ID)

	body := mustMarshal(t, CreateFamilyRequest{Name: "测试家庭-创建"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/family", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", user.ID)

	CreateFamily(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("创建家庭失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int            `json:"code"`
		Data FamilyResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Code != model.SUCCESS {
		t.Fatalf("业务返回失败: %+v", response)
	}
	if response.Data.ID == "" || response.Data.InviteCode == "" {
		t.Fatalf("返回的家庭信息不完整: %+v", response.Data)
	}
	defer cleanupFamilies(t, response.Data.ID)

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", response.Data.ID, user.ID).First(&member).Error; err != nil {
		t.Fatalf("未创建管理员成员关系: %v", err)
	}
	if member.Role != model.FamilyRoleAdmin {
		t.Fatalf("成员角色错误: %s", member.Role)
	}
}

func TestGetFamilyMembersSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("admin"), "家庭管理员")
	memberUser := createTestUser(t, uniquePhone("member"), "家庭成员")
	family := createTestFamily(t, admin.ID, "成员列表家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/family/%s/members", family.ID), nil)
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}}
	ctx.Set("userId", admin.ID)

	GetFamilyMembers(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取家庭成员失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int                    `json:"code"`
		Data []FamilyMemberResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if len(response.Data) != 2 {
		t.Fatalf("成员数量错误: %+v", response.Data)
	}
	roles := map[string]bool{}
	for _, item := range response.Data {
		roles[item.Role] = true
	}
	if !roles[string(model.FamilyRoleAdmin)] || !roles[string(model.FamilyRoleMember)] {
		t.Fatalf("成员角色集合错误: %+v", response.Data)
	}
}

func TestInviteFamilyMemberRejectNonAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("owner"), "家庭所有者")
	normalUser := createTestUser(t, uniquePhone("normal"), "普通成员")
	target := createTestUser(t, uniquePhone("target"), "目标用户")
	family := createTestFamily(t, admin.ID, "邀请权限家庭")
	ownerMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	normalMember := createTestFamilyMember(t, family.ID, normalUser.ID, model.FamilyRoleMember)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, normalUser.ID, target.ID)
	defer cleanupMembers(t, ownerMember.ID, normalMember.ID)

	body := mustMarshal(t, InviteMemberRequest{Phone: target.Phone, Role: "elder"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/family/%s/invite", family.ID), bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}}
	ctx.Set("userId", normalUser.ID)

	InviteFamilyMember(ctx)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("非管理员邀请应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestInviteFamilyMemberRejectUserAlreadyInOtherFamily(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("owner2"), "家庭所有者2")
	target := createTestUser(t, uniquePhone("other2"), "已有家庭用户")
	currentFamily := createTestFamily(t, admin.ID, "当前家庭")
	ownerMember := createTestFamilyMember(t, currentFamily.ID, admin.ID, model.FamilyRoleAdmin)
	otherFamily := createTestFamily(t, target.ID, "目标用户原家庭")
	otherOwnerMember := createTestFamilyMember(t, otherFamily.ID, target.ID, model.FamilyRoleAdmin)
	defer cleanupFamilies(t, currentFamily.ID, otherFamily.ID)
	defer cleanupUsers(t, admin.ID, target.ID)
	defer cleanupMembers(t, ownerMember.ID, otherOwnerMember.ID)

	body := mustMarshal(t, InviteMemberRequest{Phone: target.Phone, Role: "member"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/family/%s/invite", currentFamily.ID), bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: currentFamily.ID}}
	ctx.Set("userId", admin.ID)

	InviteFamilyMember(ctx)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("已在其他家庭的用户应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestJoinAndLeaveFamily(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("joinown"), "邀请家庭所有者")
	joining := createTestUser(t, uniquePhone("joiner"), "加入用户")
	family := createTestFamily(t, admin.ID, "可加入家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, joining.ID)
	defer cleanupMembers(t, adminMember.ID)

	joinBody := mustMarshal(t, JoinFamilyRequest{InviteCode: family.InviteCode, Role: "elder"})
	joinRecorder := httptest.NewRecorder()
	joinCtx, _ := gin.CreateTestContext(joinRecorder)
	joinCtx.Request = httptest.NewRequest(http.MethodPost, "/family/join", bytes.NewReader(joinBody))
	joinCtx.Request.Header.Set("Content-Type", "application/json")
	joinCtx.Set("userId", joining.ID)

	JoinFamilyByInviteCode(joinCtx)

	if joinRecorder.Code != http.StatusOK {
		t.Fatalf("加入家庭失败: status=%d body=%s", joinRecorder.Code, joinRecorder.Body.String())
	}
	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", family.ID, joining.ID).First(&member).Error; err != nil {
		t.Fatalf("未创建加入成员: %v", err)
	}
	if member.Role != model.FamilyRoleElder || member.ElderMode != 1 {
		t.Fatalf("加入角色错误: %+v", member)
	}

	leaveRecorder := httptest.NewRecorder()
	leaveCtx, _ := gin.CreateTestContext(leaveRecorder)
	leaveCtx.Request = httptest.NewRequest(http.MethodDelete, "/family/"+family.ID+"/leave", nil)
	leaveCtx.Params = gin.Params{{Key: "id", Value: family.ID}}
	leaveCtx.Set("userId", joining.ID)

	LeaveFamily(leaveCtx)

	if leaveRecorder.Code != http.StatusOK {
		t.Fatalf("离开家庭失败: status=%d body=%s", leaveRecorder.Code, leaveRecorder.Body.String())
	}
	var count int64
	mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ? AND user_id = ?", family.ID, joining.ID).Count(&count)
	if count != 0 {
		t.Fatalf("成员未离开家庭")
	}
}

func uniquePhone(prefix string) string {
	seed := time.Now().UnixNano() % 1000000000
	return fmt.Sprintf("13%09d", (seed+int64(len(prefix))*97)%1000000000)
}

func createTestUser(t *testing.T, phone, nickname string) model.User {
	t.Helper()
	user := model.User{Phone: phone, Nickname: nickname, Status: 1}
	if err := mysql.DB.Create(&user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	return user
}

func createTestFamily(t *testing.T, createdBy string, name string) model.Family {
	t.Helper()
	code, err := generateInviteCode()
	if err != nil {
		t.Fatalf("生成邀请码失败: %v", err)
	}
	family := model.Family{Name: name, InviteCode: code, CreatedBy: createdBy, StorageLimit: defaultStorageLimit}
	if err := mysql.DB.Create(&family).Error; err != nil {
		t.Fatalf("创建测试家庭失败: %v", err)
	}
	return family
}

func createTestFamilyMember(t *testing.T, familyID, userID string, role model.FamilyRole) model.FamilyMember {
	t.Helper()
	member := model.FamilyMember{FamilyID: familyID, UserID: userID, Role: role, ElderMode: elderModeFromRole(role), Notifications: model.JSON([]byte(`{"feed":true}`)), JoinedAt: time.Now()}
	if err := mysql.DB.Create(&member).Error; err != nil {
		t.Fatalf("创建测试成员失败: %v", err)
	}
	return member
}

func cleanupFamilies(t *testing.T, familyIDs ...string) {
	t.Helper()
	if len(familyIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.PhotoComment{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.PhotoLike{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.Photo{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.PhotoComment{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.PhotoLike{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.Baby{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.Record{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.FamilyMember{}).Error
	_ = mysql.DB.Where("id IN ?", familyIDs).Delete(&model.Family{}).Error
}

func cleanupMembers(t *testing.T, memberIDs ...string) {
	t.Helper()
	if len(memberIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("id IN ?", memberIDs).Delete(&model.FamilyMember{}).Error
}

func cleanupUsers(t *testing.T, userIDs ...string) {
	t.Helper()
	if len(userIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("id IN ?", userIDs).Delete(&model.User{}).Error
}

func mustMarshal(t *testing.T, value interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("JSON 编码失败: %v", err)
	}
	return data
}

func decodeResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("JSON 解码失败: %v body=%s", err, string(body))
	}
}

func getFamilyProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}
	for i := 0; i < 5; i++ {
		configPath := filepath.Join(dir, "config", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatalf("未找到项目根目录")
	return ""
}

func TestRemoveFamilyMemberSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("rmadmin"), "移除管理员")
	memberUser := createTestUser(t, uniquePhone("rmmember"), "待移除成员")
	family := createTestFamily(t, admin.ID, "移除成员家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/family/%s/members/%s", family.ID, memberUser.ID), nil)
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}, {Key: "userId", Value: memberUser.ID}}
	ctx.Set("userId", admin.ID)

	RemoveFamilyMember(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("移除成员失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var check model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", family.ID, memberUser.ID).First(&check).Error; err == nil {
		t.Fatalf("成员未被移除")
	}
}

func TestRemoveFamilyMemberRejectSelf(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("rmself"), "移除自己")
	family := createTestFamily(t, admin.ID, "移除自己家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, adminMember.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/family/%s/members/%s", family.ID, admin.ID), nil)
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}, {Key: "userId", Value: admin.ID}}
	ctx.Set("userId", admin.ID)

	RemoveFamilyMember(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("移除自己应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRemoveFamilyMemberRejectNonAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("rmown"), "家庭管理员")
	memberUser := createTestUser(t, uniquePhone("rmuser"), "普通成员")
	family := createTestFamily(t, admin.ID, "权限家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	member := createTestFamilyMember(t, family.ID, memberUser.ID, model.FamilyRoleMember)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, memberUser.ID)
	defer cleanupMembers(t, adminMember.ID, member.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/family/%s/members/%s", family.ID, admin.ID), nil)
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}, {Key: "userId", Value: admin.ID}}
	ctx.Set("userId", memberUser.ID)

	RemoveFamilyMember(ctx)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("非管理员移除应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestJoinFamilyByInviteCode(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("joinadmin"), "邀请管理员")
	joiner := createTestUser(t, uniquePhone("joiner"), "加入成员")
	family := createTestFamily(t, admin.ID, "邀请码家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID, joiner.ID)
	defer cleanupMembers(t, adminMember.ID)

	body := mustMarshal(t, JoinFamilyRequest{InviteCode: family.InviteCode})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/family/join", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", joiner.ID)

	JoinFamily(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("加入家庭失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", family.ID, joiner.ID).First(&member).Error; err != nil {
		t.Fatalf("未创建加入成员: %v", err)
	}
	if member.Role != model.FamilyRoleMember {
		t.Fatalf("加入成员角色错误: %s", member.Role)
	}
}

func TestLeaveFamilyRejectLastAdmin(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("leaveadmin"), "最后管理员")
	family := createTestFamily(t, admin.ID, "退出家庭")
	adminMember := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, adminMember.ID)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/family/%s/leave", family.ID), nil)
	ctx.Params = gin.Params{{Key: "id", Value: family.ID}}
	ctx.Set("userId", admin.ID)

	LeaveFamily(ctx)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("最后管理员退出应被拒绝: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
