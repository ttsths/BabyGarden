package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetDailyStats 日统计
// @Summary 日统计
// @Description 获取指定日期的喂养/睡眠/排泄统计
// @Tags 统计
// @Accept json
// @Produce json
// @Security Bearer
// @Param baby_id query string true "宝宝ID"
// @Param date query string false "日期(YYYY-MM-DD)"
// @Success 200 {object} model.Response{data=DailyStatsResponse}
// @Router /api/v1/stats/daily [get]
func GetDailyStats(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	date, err := parseStatsDate(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
		return
	}
	start := date
	end := start.Add(24 * time.Hour)

	var records []model.Record
	if err := mysql.DB.Where("baby_id = ? AND started_at >= ? AND started_at < ?", babyID, start, end).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	stats := DailyStatsResponse{}
	for _, record := range records {
		switch record.Type {
		case model.RecordTypeFeeding:
			stats.Feeding.Count++
			var payload model.FeedingContent
			_ = json.Unmarshal(record.Content, &payload)
			stats.Feeding.TotalAmount += payload.Amount
		case model.RecordTypeSleep:
			stats.Sleep.Count++
			if record.EndedAt.Valid {
				stats.Sleep.TotalHours += record.EndedAt.Time.Sub(record.StartedAt).Hours()
			}
		case model.RecordTypeDiaper:
			stats.Diaper.Count++
		}
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: stats})
}

// GetWeeklyStats 周统计
// @Summary 周统计
// @Description 获取近7天喂养/睡眠/排泄统计
// @Tags 统计
// @Accept json
// @Produce json
// @Security Bearer
// @Param baby_id query string true "宝宝ID"
// @Success 200 {object} model.Response{data=WeeklyStatsResponse}
// @Router /api/v1/stats/weekly [get]
func GetWeeklyStats(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	endDate, err := parseStatsDate(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
		return
	}
	startDate := endDate.AddDate(0, 0, -6)
	end := endDate.Add(24 * time.Hour)

	dates := make([]string, 0, 7)
	feeding := make([]int, 7)
	sleep := make([]float64, 7)
	diaper := make([]int, 7)
	for i := 0; i < 7; i++ {
		day := startDate.AddDate(0, 0, i)
		dates = append(dates, day.Format("2006-01-02"))
	}

	var records []model.Record
	if err := mysql.DB.Where("baby_id = ? AND started_at >= ? AND started_at < ?", babyID, startDate, end).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	for _, record := range records {
		idx := int(record.StartedAt.In(time.Local).Sub(startDate).Hours() / 24)
		if idx < 0 || idx >= 7 {
			continue
		}
		switch record.Type {
		case model.RecordTypeFeeding:
			feeding[idx]++
		case model.RecordTypeSleep:
			if record.EndedAt.Valid {
				sleep[idx] += record.EndedAt.Time.Sub(record.StartedAt).Hours()
			}
		case model.RecordTypeDiaper:
			diaper[idx]++
		}
	}

	response := WeeklyStatsResponse{
		Dates:   dates,
		Feeding: feeding,
		Sleep:   sleep,
		Diaper:  diaper,
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: response})
}

func parseStatsDate(value string) (time.Time, error) {
	if value == "" {
		now := time.Now().In(time.Local)
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

type DailyStatsResponse struct {
	Feeding FeedingStats `json:"feeding"`
	Sleep   SleepStats   `json:"sleep"`
	Diaper  DiaperStats  `json:"diaper"`
}

type FeedingStats struct {
	Count       int `json:"count"`
	TotalAmount int `json:"total_amount"`
}

type SleepStats struct {
	Count      int     `json:"count"`
	TotalHours float64 `json:"total_hours"`
}

type DiaperStats struct {
	Count int `json:"count"`
}

type WeeklyStatsResponse struct {
	Dates   []string  `json:"dates"`
	Feeding []int     `json:"feeding"`
	Sleep   []float64 `json:"sleep"`
	Diaper  []int     `json:"diaper"`
}
