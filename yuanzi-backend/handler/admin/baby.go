package admin

import (
	"net/http"
	"time"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetBabies returns paginated baby list.
// GET /api/v1/admin/babies
func GetBabies(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var babies []model.Baby

	query := mysql.DB.Model(&model.Baby{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&babies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	type babyItem struct {
		ID        string `json:"id"`
		FamilyID  string `json:"family_id"`
		Name      string `json:"name"`
		Birthday  string `json:"birthday"`
		Gender    int8   `json:"gender"`
		AvatarURL string `json:"avatar_url"`
		AgeMonths int    `json:"age_months"`
	}

	items := make([]babyItem, len(babies))
	for i, b := range babies {
		items[i] = babyItem{
			ID:        b.ID,
			FamilyID:  b.FamilyID,
			Name:      b.Name,
			Birthday:  b.Birthday.Format("2006-01-02"),
			Gender:    b.Gender,
			AvatarURL: b.AvatarURL,
			AgeMonths: b.AgeInMonths(),
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: items,
			Pagination: model.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	})
}

// GetBaby returns baby detail.
// GET /api/v1/admin/babies/:id
func GetBaby(c *gin.Context) {
	id := c.Param("id")
	var baby model.Baby
	if err := mysql.DB.Where("id = ?", id).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"id":            baby.ID,
			"family_id":     baby.FamilyID,
			"name":          baby.Name,
			"birthday":      baby.Birthday.Format("2006-01-02"),
			"gender":        baby.Gender,
			"birth_weight":  baby.BirthWeight,
			"birth_height":  baby.BirthHeight,
			"avatar_url":    baby.AvatarURL,
			"note":          baby.Note,
			"is_premature":  baby.IsPremature == 1,
			"age_months":    baby.AgeInMonths(),
			"created_at":    baby.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// CreateBaby creates a new baby.
// POST /api/v1/admin/babies
func CreateBaby(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Gender    int8   `json:"gender" binding:"required,oneof=1 2"`
		Birthday  string `json:"birthday" binding:"required"`
		FamilyID  string `json:"family_id" binding:"required"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "生日格式错误，应为 YYYY-MM-DD"})
		return
	}

	baby := model.Baby{
		FamilyID:   req.FamilyID,
		Name:       req.Name,
		Birthday:   birthday,
		Gender:     req.Gender,
		AvatarURL:  req.AvatarURL,
	}

	if err := mysql.DB.Create(&baby).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "创建成功",
		Data: gin.H{
			"id":         baby.ID,
			"family_id":  baby.FamilyID,
			"name":       baby.Name,
			"birthday":   baby.Birthday.Format("2006-01-02"),
			"gender":     baby.Gender,
			"avatar_url": baby.AvatarURL,
		},
	})
}

// UpdateBaby admin 更新宝宝信息
// PUT /api/v1/admin/babies/:id
func UpdateBaby(c *gin.Context) {
	id := c.Param("id")
	var baby model.Baby
	if err := mysql.DB.Where("id = ?", id).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return
	}

	var req struct {
		Name      string `json:"name" binding:"omitempty"`
		Gender    string `json:"gender" binding:"omitempty,oneof=male female unknown"`
		Birthday  string `json:"birthday" binding:"omitempty"`
		AvatarURL string `json:"avatar_url" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	if req.Name != "" {
		baby.Name = req.Name
	}
	if req.Gender != "" {
		// 字符串映射到 int8
		genderMap := map[string]int8{"male": 1, "female": 2, "unknown": 0}
		if g, ok := genderMap[req.Gender]; ok {
			baby.Gender = g
		}
	}
	if req.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "生日格式错误，应为 YYYY-MM-DD"})
			return
		}
		baby.Birthday = birthday
	}
	if req.AvatarURL != "" {
		baby.AvatarURL = req.AvatarURL
	}

	if err := mysql.DB.Save(&baby).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
		Data: gin.H{
			"id":         baby.ID,
			"family_id":  baby.FamilyID,
			"name":       baby.Name,
			"birthday":   baby.Birthday.Format("2006-01-02"),
			"gender":     baby.Gender,
			"avatar_url": baby.AvatarURL,
		},
	})
}

// DeleteBaby soft-deletes a baby.
// DELETE /api/v1/admin/babies/:id
func DeleteBaby(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Where("id = ?", id).Delete(&model.Baby{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已删除"})
}
