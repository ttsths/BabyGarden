package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRecord 创建记录
// @Summary 创建记录
// @Description 创建喂养/睡眠/尿布/排泄/测温/成长记录
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body CreateRecordRequest true "记录信息"
// @Success 200 {object} model.Response{data=RecordResponse}
// @Router /api/v1/record [post]
func CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}

	baby, _, err := loadBabyForRecord(c, req.BabyID)
	if err != nil {
		return
	}

	startedAt, err := time.Parse(time.RFC3339, req.StartedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "开始时间格式错误"})
		return
	}
	var endedAt sql.NullTime
	if req.EndedAt != nil {
		parsed, err := time.Parse(time.RFC3339, *req.EndedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束时间格式错误"})
			return
		}
		if parsed.Before(startedAt) {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束时间不能早于开始时间"})
			return
		}
		endedAt = sql.NullTime{Time: parsed, Valid: true}
	}

	recordType := model.RecordType(req.Type)
	contentJSON, err := normalizeRecordContent(recordType, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: err.Error()})
		return
	}

	userID := middleware.GetUserIDOrZero(c)
	record := model.Record{
		BabyID:    baby.ID,
		FamilyID:  baby.FamilyID,
		Type:      recordType,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		Content:   contentJSON,
		Note:      req.Note,
		CreatedBy: userID,
	}

	if err := mysql.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建失败"})
		return
	}

	contentMap := recordContentToMap(record.Content)
	response := RecordResponse{
		ID:        record.ID,
		BabyID:    record.BabyID,
		Type:      string(record.Type),
		StartedAt: record.StartedAt.Format(time.RFC3339),
		EndedAt:   formatNullTime(record.EndedAt),
		Content:   contentMap,
	}
	if record.Type == model.RecordTypeFeeding {
		hoursSince, err := hoursSinceLastFeeding(record.BabyID, record.StartedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "计算喂养间隔失败"})
			return
		}
		response.HoursSinceLastFeed = hoursSince
	}
	if record.Type == model.RecordTypeSleep && record.EndedAt.Valid {
		duration := record.EndedAt.Time.Sub(record.StartedAt).Hours()
		response.DurationHours = &duration
	}

	publishSyncEvent(record.FamilyID, "record_created", response)
	go dispatchRecordPush("record_created", &record, contentMap)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "创建成功",
		Data: response,
	})
}

// ListRecords 获取记录列表
// @Summary 获取记录列表
// @Description 分页获取指定宝宝的记录列表
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param baby_id query string true "宝宝ID"
// @Param type query string false "记录类型: feeding/sleep/diaper/excretion/temperature/growth"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页数量，默认20" default(20)
// @Success 200 {object} model.Response{data=model.ListResponse{list=[]RecordResponse}}
// @Router /api/v1/record [get]
func ListRecords(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	recordType := c.Query("type")
	date := c.Query("date")
	page := parsePage(c.DefaultQuery("page", "1"))
	pageSize := parsePageSize(c.DefaultQuery("page_size", "20"))

	query := mysql.DB.Model(&model.Record{}).Where("baby_id = ?", babyID)
	if recordType != "" {
		query = query.Where("type = ?", recordType)
	}
	if date != "" {
		start, err := time.ParseInLocation("2006-01-02", date, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
			return
		}
		end := start.Add(24 * time.Hour)
		query = query.Where("started_at >= ? AND started_at < ?", start, end)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	var records []model.Record
	if err := query.Order("started_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	list := make([]RecordResponse, 0, len(records))
	for _, item := range records {
		list = append(list, RecordResponse{
			ID:        item.ID,
			BabyID:    item.BabyID,
			Type:      string(item.Type),
			StartedAt: item.StartedAt.Format(time.RFC3339),
			EndedAt:   formatNullTime(item.EndedAt),
			Content:   recordContentToMap(item.Content),
		})
	}
	totalPages := calcTotalPages(total, pageSize)
	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: list,
			Pagination: model.Pagination{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}

// GetRecord 获取记录详情
// @Summary 获取记录详情
// @Description 获取指定记录的详细信息
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "记录ID"
// @Success 200 {object} model.Response{data=RecordDetailResponse}
// @Router /api/v1/record/{id} [get]
func GetRecord(c *gin.Context) {
	record, _, err := loadRecordWithMemberAccess(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: RecordDetailResponse{
			ID:        record.ID,
			BabyID:    record.BabyID,
			Type:      string(record.Type),
			StartedAt: record.StartedAt.Format(time.RFC3339),
			EndedAt:   formatNullTime(record.EndedAt),
			Content:   recordContentToMap(record.Content),
			Note:      record.Note,
			CreatedBy: record.CreatedBy,
			CreatedAt: record.CreatedAt.Format(time.RFC3339),
		},
	})
}

// UpdateRecord 更新记录
// @Summary 更新记录
// @Description 更新指定记录的信息
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "记录ID"
// @Param data body UpdateRecordRequest true "更新信息"
// @Success 200 {object} model.Response{data=RecordResponse}
// @Router /api/v1/record/{id} [put]
func UpdateRecord(c *gin.Context) {
	record, member, err := loadRecordWithMemberAccess(c)
	if err != nil {
		return
	}
	if !canEditRecord(middleware.GetUserIDOrZero(c), member, record) {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权修改记录"})
		return
	}

	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}

	if req.StartedAt != nil {
		parsed, err := time.Parse(time.RFC3339, *req.StartedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "开始时间格式错误"})
			return
		}
		if record.EndedAt.Valid && record.EndedAt.Time.Before(parsed) {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "开始时间不能晚于结束时间"})
			return
		}
		record.StartedAt = parsed
	}
	if req.EndedAt != nil {
		parsed, err := time.Parse(time.RFC3339, *req.EndedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束时间格式错误"})
			return
		}
		if parsed.Before(record.StartedAt) {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束时间不能早于开始时间"})
			return
		}
		record.EndedAt = sql.NullTime{Time: parsed, Valid: true}
	}
	if req.Content != nil {
		contentJSON, err := normalizeRecordContent(record.Type, req.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: err.Error()})
			return
		}
		record.Content = contentJSON
	}
	if req.Note != nil {
		record.Note = *req.Note
	}

	if err := mysql.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "更新失败"})
		return
	}

	updatePayload := map[string]interface{}{
		"id":         record.ID,
		"baby_id":    record.BabyID,
		"family_id":  record.FamilyID,
		"type":       string(record.Type),
		"started_at": record.StartedAt.Format(time.RFC3339),
		"ended_at":   formatNullTime(record.EndedAt),
		"content":    recordContentToMap(record.Content),
		"note":       record.Note,
	}
	publishSyncEvent(record.FamilyID, "record_updated", updatePayload)
	go dispatchRecordPush("record_updated", record, recordContentToMap(record.Content))

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "更新成功",
	})
}

// DeleteRecord 删除记录
// @Summary 删除记录
// @Description 删除指定记录（软删除）
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "记录ID"
// @Success 200 {object} model.Response
// @Router /api/v1/record/{id} [delete]
func DeleteRecord(c *gin.Context) {
	record, member, err := loadRecordWithMemberAccess(c)
	if err != nil {
		return
	}
	if !canEditRecord(middleware.GetUserIDOrZero(c), member, record) {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权删除记录"})
		return
	}

	if err := mysql.DB.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "删除失败"})
		return
	}

	deletePayload := map[string]interface{}{
		"id":        record.ID,
		"baby_id":   record.BabyID,
		"family_id": record.FamilyID,
		"type":      string(record.Type),
	}
	publishSyncEvent(record.FamilyID, "record_deleted", deletePayload)
	go dispatchRecordPush("record_deleted", record, recordContentToMap(record.Content))

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "删除成功",
	})
}

// 请求响应结构

type CreateRecordRequest struct {
	BabyID    string                 `json:"baby_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Type      string                 `json:"type" binding:"required,oneof=feeding sleep diaper excretion temperature growth" example:"feeding"`
	StartedAt string                 `json:"started_at" binding:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2024-03-08T10:00:00Z"`
	EndedAt   *string                `json:"ended_at,omitempty" example:"2024-03-08T10:15:00Z"`
	Content   map[string]interface{} `json:"content" binding:"required"`
	Note      string                 `json:"note,omitempty" example:"备注信息"`
}

type RecordResponse struct {
	ID        string                 `json:"id"`
	BabyID    string                 `json:"baby_id"`
	Type      string                 `json:"type"`
	StartedAt string                 `json:"started_at"`
	EndedAt   *string                `json:"ended_at,omitempty"`
	Content   map[string]interface{} `json:"content"`
	// 喂养记录的间隔小时数（无上一次记录时为空）
	HoursSinceLastFeed *int `json:"hours_since_last_feed,omitempty"`
	// 睡眠记录时长（小时）
	DurationHours *float64 `json:"duration_hours,omitempty"`
}

type RecordDetailResponse struct {
	ID        string                 `json:"id"`
	BabyID    string                 `json:"baby_id"`
	Type      string                 `json:"type"`
	StartedAt string                 `json:"started_at"`
	EndedAt   *string                `json:"ended_at,omitempty"`
	Content   map[string]interface{} `json:"content"`
	Note      string                 `json:"note,omitempty"`
	CreatedBy string                 `json:"created_by"`
	CreatedAt string                 `json:"created_at"`
}

type UpdateRecordRequest struct {
	StartedAt *string                `json:"started_at,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndedAt   *string                `json:"ended_at,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Content   map[string]interface{} `json:"content,omitempty"`
	Note      *string                `json:"note,omitempty"`
}

func loadBabyForRecord(c *gin.Context, babyID string) (*model.Baby, *model.FamilyMember, error) {
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID错误"})
		return nil, nil, errors.New("invalid baby id")
	}
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return nil, nil, errors.New("unauthorized")
	}

	var baby model.Baby
	if err := mysql.DB.Where("id = ?", babyID).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return nil, nil, err
	}

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", baby.FamilyID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权访问"})
		return nil, nil, err
	}
	return &baby, &member, nil
}

func loadRecordWithMemberAccess(c *gin.Context) (*model.Record, *model.FamilyMember, error) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "记录ID错误"})
		return nil, nil, errors.New("invalid record id")
	}
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return nil, nil, errors.New("unauthorized")
	}

	var record model.Record
	if err := mysql.DB.Where("id = ?", recordID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "记录不存在"})
		return nil, nil, err
	}

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", record.FamilyID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权访问"})
		return nil, nil, err
	}
	return &record, &member, nil
}

func normalizeRecordContent(recordType model.RecordType, content map[string]interface{}) (model.JSON, error) {
	if content == nil {
		return nil, errors.New("记录内容不能为空")
	}
	raw, err := json.Marshal(content)
	if err != nil {
		return nil, errors.New("记录内容格式错误")
	}

	switch recordType {
	case model.RecordTypeFeeding:
		var payload model.FeedingContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Type == "" {
			return nil, errors.New("喂养记录内容不完整")
		}
	case model.RecordTypeSleep:
		var payload model.SleepContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Quality == "" || payload.Location == "" {
			return nil, errors.New("睡眠记录内容不完整")
		}
	case model.RecordTypeDiaper:
		var payload model.DiaperContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Type == "" {
			return nil, errors.New("尿布记录内容不完整")
		}
	case model.RecordTypeExcretion:
		var payload model.ExcretionContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Type == "" {
			return nil, errors.New("排泄记录内容不完整")
		}
	case model.RecordTypeTemperature:
		var payload model.TemperatureContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Value <= 0 {
			return nil, errors.New("测温记录内容不完整")
		}
	case model.RecordTypeGrowth:
		var payload model.GrowthContent
		if err := json.Unmarshal(raw, &payload); err != nil || payload.Weight <= 0 || payload.Height <= 0 {
			return nil, errors.New("成长记录内容不完整")
		}
	default:
		return nil, errors.New("记录类型不支持")
	}
	return model.JSON(raw), nil
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

func hoursSinceLastFeeding(babyID string, startedAt time.Time) (*int, error) {
	var last model.Record
	if err := mysql.DB.Where("baby_id = ? AND type = ? AND started_at < ?", babyID, model.RecordTypeFeeding, startedAt).Order("started_at desc").First(&last).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	hours := int(startedAt.Sub(last.StartedAt).Hours())
	if hours < 0 {
		hours = 0
	}
	return &hours, nil
}

func formatNullTime(value sql.NullTime) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format(time.RFC3339)
	return &formatted
}

func parsePage(value string) int {
	page, err := strconv.Atoi(value)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func parsePageSize(value string) int {
	size, err := strconv.Atoi(value)
	if err != nil || size <= 0 || size > 100 {
		return 20
	}
	return size
}

func calcTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

func canEditRecord(userID string, member *model.FamilyMember, record *model.Record) bool {
	if member != nil && member.IsAdmin() {
		return true
	}
	return record != nil && record.CreatedBy == userID
}
