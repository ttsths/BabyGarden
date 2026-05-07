package admin

import (
	"net/http"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetFamilies returns paginated family list.
// GET /api/v1/admin/families
func GetFamilies(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var families []model.Family

	query := mysql.DB.Model(&model.Family{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&families).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	type familyItem struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		InviteCode string `json:"invite_code"`
		MemberCount int64 `json:"member_count"`
		CreatedAt  string `json:"created_at"`
	}

	items := make([]familyItem, len(families))
	for i, f := range families {
		var count int64
		mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ?", f.ID).Count(&count)
		items[i] = familyItem{
			ID:         f.ID,
			Name:       f.Name,
			InviteCode: f.InviteCode,
			MemberCount: count,
			CreatedAt:  f.CreatedAt.Format("2006-01-02 15:04:05"),
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

// GetFamily returns family detail with members.
// GET /api/v1/admin/families/:id
func GetFamily(c *gin.Context) {
	id := c.Param("id")
	var family model.Family
	if err := mysql.DB.Where("id = ?", id).First(&family).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "家庭不存在"})
		return
	}

	var members []model.FamilyMember
	mysql.DB.Preload("User").Where("family_id = ?", family.ID).Order("joined_at ASC").Find(&members)

	type memberItem struct {
		UserID    string `json:"user_id"`
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
		Role      string `json:"role"`
		ElderMode bool   `json:"elder_mode"`
	}
	memberItems := make([]memberItem, len(members))
	for i, m := range members {
		memberItems[i] = memberItem{
			UserID:    m.UserID,
			Nickname:  m.User.Nickname,
			AvatarURL: m.User.AvatarURL,
			Role:      string(m.Role),
			ElderMode: m.ElderMode == 1,
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"id":          family.ID,
			"name":        family.Name,
			"invite_code": family.InviteCode,
			"members":     memberItems,
			"created_at":  family.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteFamily disables a family.
// DELETE /api/v1/admin/families/:id
func DeleteFamily(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Model(&model.Family{}).Where("id = ?", id).Update("is_paid", -1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已禁用"})
}
