package admin

import (
	"net/http"
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
