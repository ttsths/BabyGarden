package admin

import (
	"net/http"
	"time"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AIUsageSummaryItem struct {
	Period       string `json:"period" gorm:"column:period"`
	Provider     string `json:"provider" gorm:"column:provider"`
	Requests     int64  `json:"requests" gorm:"column:requests"`
	InputTokens  int    `json:"input_tokens" gorm:"column:input_tokens"`
	OutputTokens int    `json:"output_tokens" gorm:"column:output_tokens"`
	CachedTokens int    `json:"cached_tokens" gorm:"column:cached_tokens"`
	TotalTokens  int    `json:"total_tokens" gorm:"column:total_tokens"`
}

type AIUsageOverview struct {
	TodayTotalTokens int `json:"today_total_tokens"`
	WeekTotalTokens  int `json:"week_total_tokens"`
	MonthTotalTokens int `json:"month_total_tokens"`
}

// GetAIUsage returns paginated AI token usage logs.
// GET /api/v1/admin/ai/usage
func GetAIUsage(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var logs []model.AIUsageLog

	query := applyAIUsageFilters(mysql.DB.Model(&model.AIUsageLog{}), c)
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询AI使用记录失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: logs,
			Pagination: model.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	})
}

// GetAIUsageSummary returns aggregated AI usage by day/week/month.
// GET /api/v1/admin/ai/usage/summary
func GetAIUsageSummary(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	days := atoi(c.DefaultQuery("days", "30"))
	if days <= 0 || days > 365 {
		days = 30
	}

	now := time.Now()
	start := now.AddDate(0, 0, -days+1)
	query := mysql.DB.Model(&model.AIUsageLog{}).Where("created_at >= ?", start)
	query = applyAIUsageDimensionFilters(query, c)

	periodExpr := "DATE_FORMAT(created_at, '%Y-%m-%d')"
	switch period {
	case "week":
		periodExpr = "DATE_FORMAT(DATE_SUB(created_at, INTERVAL WEEKDAY(created_at) DAY), '%Y-%m-%d')"
	case "month":
		periodExpr = "DATE_FORMAT(created_at, '%Y-%m')"
	}

	var items []AIUsageSummaryItem
	if err := query.
		Select(periodExpr + " AS period, provider, COUNT(*) AS requests, COALESCE(SUM(input_tokens), 0) AS input_tokens, COALESCE(SUM(output_tokens), 0) AS output_tokens, COALESCE(SUM(cached_tokens), 0) AS cached_tokens, COALESCE(SUM(total_tokens), 0) AS total_tokens").
		Group("period, provider").
		Order("period ASC, provider ASC").
		Scan(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询AI使用汇总失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: gin.H{
		"period": period,
		"items":  items,
	}})
}

// GetUserAIUsage returns one user's AI usage detail and overview.
// GET /api/v1/admin/ai/usage/:userId
func GetUserAIUsage(c *gin.Context) {
	userID := c.Param("userId")
	page, pageSize := paginate(c)

	query := applyAIUsageFilters(mysql.DB.Model(&model.AIUsageLog{}).Where("user_id = ?", userID), c)
	var total int64
	query.Count(&total)

	var logs []model.AIUsageLog
	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询用户AI使用详情失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"overview": aiUsageOverview(mysql.DB.Model(&model.AIUsageLog{}).Where("user_id = ? AND status = ?", userID, "success")),
			"logs": model.ListResponse{
				List: logs,
				Pagination: model.Pagination{
					Page:     page,
					PageSize: pageSize,
					Total:    total,
				},
			},
		},
	})
}

func applyAIUsageFilters(query *gorm.DB, c *gin.Context) *gorm.DB {
	query = applyAIUsageDimensionFilters(query, c)
	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			query = query.Where("created_at < ?", t.AddDate(0, 0, 1))
		}
	}
	return query
}

func applyAIUsageDimensionFilters(query *gorm.DB, c *gin.Context) *gorm.DB {
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if familyID := c.Query("family_id"); familyID != "" {
		query = query.Where("family_id = ?", familyID)
	}
	if provider := c.Query("provider"); provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if requestType := c.Query("request_type"); requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	return query
}

func aiUsageOverview(query *gorm.DB) AIUsageOverview {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()+6)%7)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	return AIUsageOverview{
		TodayTotalTokens: sumTotalTokens(query.Session(&gorm.Session{}), todayStart),
		WeekTotalTokens:  sumTotalTokens(query.Session(&gorm.Session{}), weekStart),
		MonthTotalTokens: sumTotalTokens(query.Session(&gorm.Session{}), monthStart),
	}
}

func sumTotalTokens(query *gorm.DB, start time.Time) int {
	var total int
	_ = query.Where("created_at >= ?", start).Select("COALESCE(SUM(total_tokens), 0)").Scan(&total).Error
	return total
}
