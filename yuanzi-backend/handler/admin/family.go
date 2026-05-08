package admin

import (
	"net/http"
	"time"
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
		ID          string `json:"id"`
		Name        string `json:"name"`
		InviteCode  string `json:"invite_code"`
		MemberCount int64  `json:"member_count"`
		CreatedAt   string `json:"created_at"`
	}

	items := make([]familyItem, len(families))
	for i, f := range families {
		var count int64
		mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ?", f.ID).Count(&count)
		items[i] = familyItem{
			ID:          f.ID,
			Name:        f.Name,
			InviteCode:  f.InviteCode,
			MemberCount: count,
			CreatedAt:   f.CreatedAt.Format("2006-01-02 15:04:05"),
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

// CreateFamily creates a new family.
// POST /api/v1/admin/families
func CreateFamily(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required,max=100"`
		InviteCode string `json:"invite_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	family := model.Family{
		Name:       req.Name,
		InviteCode: req.InviteCode,
	}

	if err := mysql.DB.Create(&family).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "创建成功",
		Data: gin.H{
			"id":          family.ID,
			"name":        family.Name,
			"invite_code": family.InviteCode,
		},
	})
}

// UpdateFamily updates family info.
// PUT /api/v1/admin/families/:id
func UpdateFamily(c *gin.Context) {
	id := c.Param("id")
	var family model.Family
	if err := mysql.DB.Where("id = ?", id).First(&family).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "家庭不存在"})
		return
	}

	var req struct {
		Name       string `json:"name" binding:"omitempty"`
		InviteCode string `json:"invite_code" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.InviteCode != "" {
		updates["invite_code"] = req.InviteCode
	}

	if len(updates) > 0 {
		if err := mysql.DB.Model(&family).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
		Data: gin.H{
			"id":          family.ID,
			"name":        family.Name,
			"invite_code": family.InviteCode,
		},
	})
}

// AddFamilyMember adds a member to a family.
// POST /api/v1/admin/families/:id/members
func AddFamilyMember(c *gin.Context) {
	id := c.Param("id")
	var family model.Family
	if err := mysql.DB.Where("id = ?", id).First(&family).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "家庭不存在"})
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required,oneof=admin member elder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var user model.User
	if err := mysql.DB.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "用户不存在"})
		return
	}

	var exists int64
	if err := mysql.DB.Model(&model.FamilyMember{}).Where("family_id = ? AND user_id = ?", family.ID, req.UserID).Count(&exists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "用户已在家庭中"})
		return
	}

	role := model.FamilyRoleMember
	if req.Role == "admin" {
		role = model.FamilyRoleAdmin
	} else if req.Role == "elder" {
		role = model.FamilyRoleElder
	}

	member := model.FamilyMember{
		FamilyID: family.ID,
		UserID:   req.UserID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := mysql.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "添加成员失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "添加成功",
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
