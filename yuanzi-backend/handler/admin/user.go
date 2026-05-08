package admin

import (
	"net/http"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetUsers returns paginated user list.
// GET /api/v1/admin/users
func GetUsers(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var users []model.User

	query := mysql.DB.Model(&model.User{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("phone LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	infos := make([]model.UserInfo, len(users))
	for i, u := range users {
		infos[i] = u.ToUserInfo()
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: infos,
			Pagination: model.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	})
}

// GetUser returns a single user by ID.
// GET /api/v1/admin/users/:id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := mysql.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: user.ToUserInfo(),
	})
}

// CreateUser creates a new user.
// POST /api/v1/admin/users
func CreateUser(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required,len=11"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname" binding:"required"`
		Status   int8   `json:"status" binding:"oneof=0 1"`
		IsAdmin  int8   `json:"is_admin" binding:"oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var exists model.User
	if err := mysql.DB.Where("phone = ?", req.Phone).First(&exists).Error; err == nil {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_INVALID, Msg: "手机号已存在"})
		return
	}

	user := model.User{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Status:   req.Status,
		IsAdmin:  req.IsAdmin,
	}
	user.SetPassword(req.Password)

	if err := mysql.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "创建成功",
		Data: user.ToUserInfo(),
	})
}

// UpdateUser updates user info.
// PUT /api/v1/admin/users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := mysql.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "用户不存在"})
		return
	}

	var req struct {
		Nickname string `json:"nickname" binding:"omitempty"`
		Status   *int8  `json:"status" binding:"omitempty,oneof=0 1"`
		IsAdmin  *int8  `json:"is_admin" binding:"omitempty,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}

	if len(updates) > 0 {
		if err := mysql.DB.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
		Data: user.ToUserInfo(),
	})
}

// DeleteUser soft-deletes a user (sets status=0).
// DELETE /api/v1/admin/users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Model(&model.User{}).Where("id = ?", id).Update("status", 0).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已禁用"})
}

// UpdateUserStatus enables or disables a user.
// PUT /api/v1/admin/users/:id/status
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int8 `json:"status" binding:"oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	if err := mysql.DB.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}

	msg := "已启用"
	if req.Status == 0 {
		msg = "已禁用"
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: msg})
}

func paginate(c *gin.Context) (int, int) {
	page := 1
	pageSize := 20
	if p, ok := c.GetQuery("page"); ok {
		if v := atoi(p); v > 0 {
			page = v
		}
	}
	if ps, ok := c.GetQuery("page_size"); ok {
		if v := atoi(ps); v > 0 && v <= 100 {
			pageSize = v
		}
	}
	return page, pageSize
}

func atoi(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
