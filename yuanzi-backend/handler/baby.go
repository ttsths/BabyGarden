package handler

import (
	"net/http"
	"strconv"
	"time"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateBaby 创建宝宝
func CreateBaby(c *gin.Context) {
	var req CreateBabyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	member, err := currentFamilyMember(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "请先加入家庭"})
		return
	}

	birthday, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "生日格式错误，应为 YYYY-MM-DD"})
		return
	}

	baby := model.Baby{
		FamilyID:    member.FamilyID,
		Name:        req.Name,
		Birthday:    birthday,
		Gender:      int8(req.Gender),
		AvatarURL:   req.AvatarURL,
		Note:        req.Note,
		IsPremature: int8(req.IsPremature),
	}
	if req.BirthWeight > 0 {
		baby.BirthWeight = &req.BirthWeight
	}
	if req.BirthHeight > 0 {
		baby.BirthHeight = &req.BirthHeight
	}

	if err := mysql.DB.Create(&baby).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "创建成功", Data: baby})
}

// GetBaby 获取宝宝信息
func GetBaby(c *gin.Context) {
	baby, err := loadBabyWithAccess(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: baby})
}

// UpdateBaby 更新宝宝信息
func UpdateBaby(c *gin.Context) {
	baby, err := loadBabyWithAccess(c)
	if err != nil {
		return
	}

	var req UpdateBabyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	if req.Name != "" {
		baby.Name = req.Name
	}
	if req.BirthDate != "" {
		birthday, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "生日格式错误，应为 YYYY-MM-DD"})
			return
		}
		baby.Birthday = birthday
	}
	if req.Gender != 0 {
		baby.Gender = int8(req.Gender)
	}
	if req.AvatarURL != "" {
		baby.AvatarURL = req.AvatarURL
	}
	if req.Note != "" {
		baby.Note = req.Note
	}
	if req.IsPremature != 0 {
		baby.IsPremature = int8(req.IsPremature)
	}

	if err := mysql.DB.Save(&baby).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "更新成功", Data: baby})
}

// DeleteBaby 删除宝宝（软删除）
func DeleteBaby(c *gin.Context) {
	baby, member, err := loadBabyWithMemberAccess(c)
	if err != nil {
		return
	}
	if !member.IsAdmin() {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "仅管理员可删除宝宝"})
		return
	}

	if err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("baby_id = ?", baby.ID).Delete(&model.Record{}).Error; err != nil {
			return err
		}
		if err := tx.Where("baby_id = ?", baby.ID).Delete(&model.Photo{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&baby).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "删除失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "删除成功"})
}

// ListBabies 获取宝宝列表
func ListBabies(c *gin.Context) {
	member, err := currentFamilyMember(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "请先加入家庭"})
		return
	}

	var babies []model.Baby
	if err := mysql.DB.Where("family_id = ?", member.FamilyID).Order("id desc").Find(&babies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: babies})
}

func currentFamilyMember(c *gin.Context) (*model.FamilyMember, error) {
	userID := middleware.GetUserIDOrZero(c)
	var member model.FamilyMember
	if err := mysql.DB.Where("user_id = ?", userID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func loadBabyWithAccess(c *gin.Context) (*model.Baby, error) {
	id, err := func() (string, error) {
		id := c.Param("id")
		if id == "" {
			return "", strconv.ErrSyntax
		}
		return id, nil
	}()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "参数错误"})
		return nil, err
	}

	var baby model.Baby
	if err := mysql.DB.Where("id = ?", id).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return nil, err
	}

	userID := middleware.GetUserIDOrZero(c)
	var member model.FamilyMember
	if err := mysql.DB.Where("user_id = ? AND family_id = ?", userID, baby.FamilyID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权访问"})
		return nil, err
	}

	return &baby, nil
}

func loadBabyWithMemberAccess(c *gin.Context) (*model.Baby, *model.FamilyMember, error) {
	id, err := func() (string, error) {
		id := c.Param("id")
		if id == "" {
			return "", strconv.ErrSyntax
		}
		return id, nil
	}()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "参数错误"})
		return nil, nil, err
	}

	var baby model.Baby
	if err := mysql.DB.Where("id = ?", id).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return nil, nil, err
	}

	userID := middleware.GetUserIDOrZero(c)
	var member model.FamilyMember
	if err := mysql.DB.Where("user_id = ? AND family_id = ?", userID, baby.FamilyID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权访问"})
		return nil, nil, err
	}

	return &baby, &member, nil
}

type CreateBabyRequest struct {
	Name        string  `json:"name" binding:"required"`
	BirthDate   string  `json:"birth_date" binding:"required"`
	Gender      int     `json:"gender" binding:"required,oneof=1 2"`
	AvatarURL   string  `json:"avatar_url"`
	Note        string  `json:"note"`
	BirthWeight float64 `json:"birth_weight"`
	BirthHeight float64 `json:"birth_height"`
	IsPremature int     `json:"is_premature"`
}

type UpdateBabyRequest struct {
	Name        string `json:"name"`
	BirthDate   string `json:"birth_date"`
	Gender      int    `json:"gender" binding:"omitempty,oneof=1 2"`
	AvatarURL   string `json:"avatar_url"`
	Note        string `json:"note"`
	IsPremature int    `json:"is_premature"`
}
