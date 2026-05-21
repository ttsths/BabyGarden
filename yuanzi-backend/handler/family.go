package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultInviteTTLSeconds = 86400
	defaultStorageLimit     = int64(1073741824)
)

// CreateFamily 创建家庭。
// 当前实现约束为：一个用户同一时间只能属于一个家庭。
func CreateFamily(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var req CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var exists int64
	if err := mysql.DB.Model(&model.FamilyMember{}).Where("user_id = ?", userID).Count(&exists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "校验家庭状态失败"})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "用户已加入家庭"})
		return
	}

	inviteCode, err := generateUniqueInviteCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成邀请码失败"})
		return
	}

	family := model.Family{
		Name:         req.Name,
		InviteCode:   inviteCode,
		CreatedBy:    userID,
		StorageLimit: defaultStorageLimit,
	}
	member := model.FamilyMember{
		UserID:        userID,
		Role:          model.FamilyRoleAdmin,
		ElderMode:     0,
		Notifications: model.JSON([]byte(`{"feed":true,"sleep":true}`)),
		JoinedAt:      time.Now(),
	}

	if err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&family).Error; err != nil {
			return err
		}
		member.FamilyID = family.ID
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建家庭失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "创建成功", Data: FamilyResponse{ID: family.ID, Name: family.Name, InviteCode: family.InviteCode}})
}

// GetFamily 获取家庭详情，并校验当前用户是否属于该家庭。
func GetFamily(c *gin.Context) {
	family, _, err := loadFamilyWithMemberAccess(c)
	if err != nil {
		return
	}

	storageUsed, err := calculateFamilyStorageUsed(family.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "统计家庭存储失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: FamilyDetailResponse{ID: family.ID, Name: family.Name, InviteCode: family.InviteCode, IsPaid: family.IsPaid == 1, StorageLimit: family.StorageLimit, StorageUsed: storageUsed},
	})
}

// InviteFamilyMember 邀请家庭成员。
// 若目标手机号已注册，则直接加入家庭；否则返回邀请码供前端后续分享。
func InviteFamilyMember(c *gin.Context) {
	family, member, err := loadFamilyWithMemberAccess(c)
	if err != nil {
		return
	}
	if !member.CanManageMembers() {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "仅管理员可邀请成员"})
		return
	}

	var req InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	role := normalizeInviteRole(req.Role)
	var targetUser model.User
	if err := mysql.DB.Where("phone = ?", req.Phone).First(&targetUser).Error; err == nil {
		var count int64
		if err := mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ? AND user_id = ?", family.ID, targetUser.ID).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "校验成员状态失败"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "用户已在家庭中"})
			return
		}
		if err := ensureUserNotInAnotherFamily(targetUser.ID); err != nil {
			c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "目标用户已属于其他家庭"})
			return
		}

		newMember := model.FamilyMember{FamilyID: family.ID, UserID: targetUser.ID, Role: role, ElderMode: elderModeFromRole(role), Notifications: model.JSON([]byte(`{"feed":true,"sleep":true}`)), JoinedAt: time.Now()}
		if err := mysql.DB.Create(&newMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "添加成员失败"})
			return
		}
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "邀请已发送", Data: InviteMemberResponse{InviteCode: family.InviteCode, ExpiresIn: defaultInviteTTLSeconds}})
}

// JoinFamily 通过邀请码加入家庭。
func JoinFamily(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var req JoinFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}
	if err := ensureUserNotInAnotherFamily(userID); err != nil {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "用户已加入家庭"})
		return
	}

	var family model.Family
	if err := mysql.DB.Where("invite_code = ?", req.InviteCode).First(&family).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "邀请码无效"})
		return
	}

	role := normalizeInviteRole(req.Role)
	member := model.FamilyMember{
		FamilyID:      family.ID,
		UserID:        userID,
		Role:          role,
		ElderMode:     elderModeFromRole(role),
		Notifications: model.JSON([]byte(`{"feed":true,"sleep":true}`)),
		JoinedAt:      time.Now(),
	}
	if err := mysql.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "加入家庭失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "加入成功", Data: FamilyResponse{ID: family.ID, Name: family.Name, InviteCode: family.InviteCode}})
}

// LeaveFamily 当前用户离开家庭。
func LeaveFamily(c *gin.Context) {
	family, member, err := loadFamilyWithMemberAccess(c)
	if err != nil {
		return
	}
	if member.Role == model.FamilyRoleAdmin {
		var adminCount int64
		if err := mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ? AND role = ?", family.ID, model.FamilyRoleAdmin).Count(&adminCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "校验管理员失败"})
			return
		}
		if adminCount <= 1 {
			c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "最后一个管理员不能离开家庭"})
			return
		}
	}

	if err := mysql.DB.Delete(member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "离开家庭失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已离开家庭"})
}

// GetFamilyMembers 获取家庭成员列表。
func GetFamilyMembers(c *gin.Context) {
	family, _, err := loadFamilyWithMemberAccess(c)
	if err != nil {
		return
	}

	var members []model.FamilyMember
	if err := mysql.DB.Preload("User").Where("family_id = ?", family.ID).Order("id asc").Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询成员列表失败"})
		return
	}

	response := make([]FamilyMemberResponse, 0, len(members))
	for _, item := range members {
		response = append(response, FamilyMemberResponse{UserID: item.UserID, Nickname: item.User.Nickname, AvatarURL: item.User.AvatarURL, Role: string(item.Role), ElderMode: item.ElderMode == 1})
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: response})
}

// RemoveFamilyMember 移除家庭成员。
func RemoveFamilyMember(c *gin.Context) {
	family, member, err := loadFamilyWithMemberAccess(c)
	if err != nil {
		return
	}
	if !member.CanManageMembers() {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "仅管理员可移除成员"})
		return
	}

	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "成员ID错误"})
		return
	}
	if targetUserID == member.UserID {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "不能移除自己"})
		return
	}

	var target model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", family.ID, targetUserID).First(&target).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "成员不存在"})
		return
	}

	if target.Role == model.FamilyRoleAdmin {
		var adminCount int64
		if err := mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ? AND role = ?", family.ID, model.FamilyRoleAdmin).Count(&adminCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "校验管理员失败"})
			return
		}
		if adminCount <= 1 {
			c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "不能移除最后一个管理员"})
			return
		}
	}

	if err := mysql.DB.Delete(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "移除失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "移除成功"})
}

type CreateFamilyRequest struct {
	Name string `json:"name" binding:"required,max=100" example:"小园子的家"`
}

type FamilyResponse struct {
	ID         string `json:"id" example:"2c3d4e5f-1234-5678-9012-abcdef123456"`
	Name       string `json:"name" example:"小园子的家"`
	InviteCode string `json:"invite_code" example:"ABC12345"`
}

type FamilyDetailResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	InviteCode   string `json:"invite_code"`
	IsPaid       bool   `json:"is_paid"`
	StorageLimit int64  `json:"storage_limit"`
	StorageUsed  int64  `json:"storage_used"`
}

type InviteMemberRequest struct {
	Phone string `json:"phone" binding:"required,len=11" example:"13900139000"`
	Role  string `json:"role" binding:"omitempty,oneof=member elder" example:"member"`
}

type JoinFamilyRequest struct {
	InviteCode string `json:"invite_code" binding:"required,len=8" example:"ABC12345"`
	Role       string `json:"role" binding:"omitempty,oneof=member elder" example:"member"`
}

type InviteMemberResponse struct {
	InviteCode string `json:"invite_code" example:"ABC12345"`
	ExpiresIn  int    `json:"expires_in" example:"86400"`
}

type FamilyMemberResponse struct {
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Role      string `json:"role"`
	ElderMode bool   `json:"elder_mode"`
}

func loadFamilyWithMemberAccess(c *gin.Context) (*model.Family, *model.FamilyMember, error) {
	familyID := c.Param("id")
	if familyID == "" {
		err := fmt.Errorf("invalid family id")
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "家庭ID错误"})
		return nil, nil, err
	}

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return nil, nil, fmt.Errorf("unauthorized")
	}

	var family model.Family
	if err := mysql.DB.Where("id = ?", familyID).First(&family).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "家庭不存在"})
		return nil, nil, err
	}

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "非家庭成员无法访问"})
		return nil, nil, err
	}

	return &family, &member, nil
}

func calculateFamilyStorageUsed(familyID string) (int64, error) {
	var total int64
	if err := mysql.DB.Model(&model.Photo{}).Where("family_id = ? AND status <> ?", familyID, model.PhotoStatusDeleted).Select("COALESCE(SUM(size), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func normalizeInviteRole(role string) model.FamilyRole {
	if role == string(model.FamilyRoleElder) {
		return model.FamilyRoleElder
	}
	return model.FamilyRoleMember
}

func elderModeFromRole(role model.FamilyRole) int8 {
	if role == model.FamilyRoleElder {
		return 1
	}
	return 0
}

func generateUniqueInviteCode() (string, error) {
	for i := 0; i < 5; i++ {
		code, err := generateInviteCode()
		if err != nil {
			return "", err
		}
		var count int64
		if err := mysql.DB.Model(&model.Family{}).Where("invite_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique invite code")
}

func generateInviteCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, 8)
	for i := range buf {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		buf[i] = alphabet[n.Int64()]
	}
	return string(buf), nil
}

func ensureUserNotInAnotherFamily(userID string) error {
	var count int64
	if err := mysql.DB.Model(&model.FamilyMember{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user already belongs to another family")
	}
	return nil
}
