package admin

import (
	"net/http"
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
