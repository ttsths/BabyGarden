package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetRecords returns paginated growth record list.
// GET /api/v1/admin/records
func GetRecords(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var records []model.Record

	query := mysql.DB.Model(&model.Record{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("type = ?", keyword)
	}
	query.Count(&total)

	if err := query.Order("started_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	type recordItem struct {
		ID        string `json:"id"`
		BabyID    string `json:"baby_id"`
		FamilyID  string `json:"family_id"`
		Type      string `json:"type"`
		StartedAt string `json:"started_at"`
		Duration  int    `json:"duration_min"`
		CreatedBy string `json:"created_by"`
		Note      string `json:"note,omitempty"`
	}

	items := make([]recordItem, len(records))
	for i, r := range records {
		items[i] = recordItem{
			ID:        r.ID,
			BabyID:    r.BabyID,
			FamilyID:  r.FamilyID,
			Type:      string(r.Type),
			StartedAt: r.StartedAt.Format("2006-01-02 15:04:05"),
			Duration:  r.Duration(),
			CreatedBy: r.CreatedBy,
			Note:      r.Note,
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

// GetRecord returns record detail.
// GET /api/v1/admin/records/:id
func GetRecord(c *gin.Context) {
	id := c.Param("id")
	var record model.Record
	if err := mysql.DB.Where("id = ?", id).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "记录不存在"})
		return
	}

	contentMap := recordContentToMap(record.Content)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"id":         record.ID,
			"baby_id":    record.BabyID,
			"family_id":  record.FamilyID,
			"type":       string(record.Type),
			"started_at": record.StartedAt.Format("2006-01-02 15:04:05"),
			"duration":   record.Duration(),
			"note":       record.Note,
			"created_by": record.CreatedBy,
			"created_at": record.CreatedAt.Format("2006-01-02 15:04:05"),
			"data":       contentMap,
		},
	})
}

// CreateRecord admin 创建记录
// POST /api/v1/admin/records
func CreateRecord(c *gin.Context) {
	var req struct {
		Type      string                 `json:"type" binding:"required,oneof=feeding sleep diaper growth"`
		BabyID    string                 `json:"baby_id" binding:"required"`
		FamilyID  string                 `json:"family_id" binding:"required"`
		StartedAt string                 `json:"started_at" binding:"required"`
		EndedAt   string                 `json:"ended_at"`
		Note      string                 `json:"note"`
		Content   map[string]interface{} `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	// 支持多种时间格式
	startedAt, err := parseTime(req.StartedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "时间格式错误"})
		return
	}

	var endedAt *time.Time
	if req.EndedAt != "" {
		et, err := parseTime(req.EndedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束时间格式错误"})
			return
		}
		endedAt = &et
	}

	var baby model.Baby
	if err := mysql.DB.Where("id = ?", req.BabyID).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return
	}

	var contentJSON model.JSON
	if req.Content != nil {
		raw, _ := json.Marshal(req.Content)
		contentJSON = model.JSON(raw)
	} else {
		// 默认空 JSON 对象，避免 NOT NULL 约束错误
		contentJSON = model.JSON(json.RawMessage(`{}`))
	}

	record := model.Record{
		BabyID:    req.BabyID,
		FamilyID:  req.FamilyID,
		Type:      model.RecordType(req.Type),
		StartedAt: startedAt,
		Note:      req.Note,
		Content:   contentJSON,
		CreatedBy: "admin",
	}
	if endedAt != nil {
		record.EndedAt = sql.NullTime{Time: *endedAt, Valid: true}
	}

	if err := mysql.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "创建成功",
		Data: gin.H{
			"id":         record.ID,
			"baby_id":    record.BabyID,
			"family_id":  record.FamilyID,
			"type":       string(record.Type),
			"started_at": record.StartedAt.Format("2006-01-02 15:04:05"),
			"note":       record.Note,
		},
	})
}

func parseTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", s)
}

// UpdateRecord updates a record.
// PUT /api/v1/admin/records/:id
func UpdateRecord(c *gin.Context) {
	id := c.Param("id")
	var record model.Record
	if err := mysql.DB.Where("id = ?", id).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "记录不存在"})
		return
	}

	var req struct {
		Type      string                 `json:"type" binding:"omitempty,oneof=feeding sleep diaper growth"`
		BabyID    string                 `json:"baby_id" binding:"omitempty"`
		FamilyID  string                 `json:"family_id" binding:"omitempty"`
		StartedAt string                 `json:"started_at" binding:"omitempty"`
		Note      string                 `json:"note"`
		Content   map[string]interface{} `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	if req.Type != "" {
		record.Type = model.RecordType(req.Type)
	}
	if req.BabyID != "" {
		record.BabyID = req.BabyID
	}
	if req.FamilyID != "" {
		record.FamilyID = req.FamilyID
	}
	if req.StartedAt != "" {
		startedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.StartedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "时间格式错误"})
			return
		}
		record.StartedAt = startedAt
	}
	if req.Content != nil {
		raw, _ := json.Marshal(req.Content)
		record.Content = model.JSON(raw)
	}
	record.Note = req.Note

	if err := mysql.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
		Data: gin.H{
			"id":         record.ID,
			"baby_id":    record.BabyID,
			"family_id":  record.FamilyID,
			"type":       string(record.Type),
			"started_at": record.StartedAt.Format("2006-01-02 15:04:05"),
			"note":       record.Note,
		},
	})
}

// DeleteRecord deletes a record.
// DELETE /api/v1/admin/records/:id
func DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Where("id = ?", id).Delete(&model.Record{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已删除"})
}

func recordContentToMap(content model.JSON) map[string]interface{} {
	var payload map[string]interface{}
	if len(content) == 0 {
		return nil
	}
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil
	}
	return payload
}
