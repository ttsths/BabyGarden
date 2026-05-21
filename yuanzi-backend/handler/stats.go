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

	stats := buildDailyStats(records)

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

	response, err := buildRangeStats(babyID, startDate, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: response})
}

// GetMonthlyStats 月统计
func GetMonthlyStats(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	anchor, err := parseStatsDate(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
		return
	}
	start := time.Date(anchor.Year(), anchor.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	response, err := buildRangeStats(babyID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: response})
}

// GetRangeStats 自定义时间区间统计
func GetRangeStats(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	start, err := parseStatsDate(c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "开始日期格式错误"})
		return
	}
	endDate, err := parseStatsDate(c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束日期格式错误"})
		return
	}
	if endDate.Before(start) {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "结束日期不能早于开始日期"})
		return
	}
	end := endDate.Add(24 * time.Hour)
	response, err := buildRangeStats(babyID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
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
	Feeding     FeedingStats     `json:"feeding"`
	Sleep       SleepStats       `json:"sleep"`
	Diaper      DiaperStats      `json:"diaper"`
	Temperature TemperatureStats `json:"temperature"`
}

type FeedingStats struct {
	Count         int     `json:"count"`
	TotalAmount   int     `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

type SleepStats struct {
	Count                int     `json:"count"`
	TotalHours           float64 `json:"total_hours"`
	AverageDurationHours float64 `json:"average_duration_hours"`
	DaytimeSingleHours   float64 `json:"daytime_single_hours"`
}

type DiaperStats struct {
	Count int `json:"count"`
}

type TemperatureStats struct {
	Count  int     `json:"count"`
	Latest float64 `json:"latest"`
}

type WeeklyStatsResponse struct {
	Dates              []string  `json:"dates"`
	Feeding            []int     `json:"feeding"`
	Sleep              []float64 `json:"sleep"`
	Diaper             []int     `json:"diaper"`
	DailyAverageSleep  []float64 `json:"daily_average_sleep_hours"`
	DaytimeSingleSleep []float64 `json:"daytime_single_sleep_hours"`
	DailyAverageMilk   []float64 `json:"daily_average_milk_amount"`
	TemperatureLatest  []float64 `json:"temperature_latest"`
}

func buildDailyStats(records []model.Record) DailyStatsResponse {
	stats := DailyStatsResponse{}
	var sleepDurations []float64
	var feedingAmounts []int
	for _, record := range records {
		switch record.Type {
		case model.RecordTypeFeeding:
			stats.Feeding.Count++
			var payload model.FeedingContent
			_ = json.Unmarshal(record.Content, &payload)
			stats.Feeding.TotalAmount += payload.Amount
			if payload.Amount > 0 {
				feedingAmounts = append(feedingAmounts, payload.Amount)
			}
		case model.RecordTypeSleep:
			stats.Sleep.Count++
			if record.EndedAt.Valid {
				duration := record.EndedAt.Time.Sub(record.StartedAt).Hours()
				stats.Sleep.TotalHours += duration
				sleepDurations = append(sleepDurations, duration)
				if isDaytimeSleep(record.StartedAt) && duration > stats.Sleep.DaytimeSingleHours {
					stats.Sleep.DaytimeSingleHours = duration
				}
			}
		case model.RecordTypeDiaper:
			stats.Diaper.Count++
		case model.RecordTypeTemp:
			var payload model.TemperatureContent
			_ = json.Unmarshal(record.Content, &payload)
			stats.Temperature.Count++
			stats.Temperature.Latest = payload.Value
		}
	}
	if len(sleepDurations) > 0 {
		stats.Sleep.AverageDurationHours = stats.Sleep.TotalHours / float64(len(sleepDurations))
	}
	if len(feedingAmounts) > 0 {
		stats.Feeding.AverageAmount = float64(stats.Feeding.TotalAmount) / float64(len(feedingAmounts))
	}
	return stats
}

func buildRangeStats(babyID string, start, end time.Time) (WeeklyStatsResponse, error) {
	dayCount := int(end.Sub(start).Hours() / 24)
	if dayCount < 1 {
		dayCount = 1
	}
	dates := make([]string, 0, dayCount)
	buckets := make([][]model.Record, dayCount)
	for i := 0; i < dayCount; i++ {
		day := start.AddDate(0, 0, i)
		dates = append(dates, day.Format("2006-01-02"))
	}

	var records []model.Record
	if err := mysql.DB.Where("baby_id = ? AND started_at >= ? AND started_at < ?", babyID, start, end).Find(&records).Error; err != nil {
		return WeeklyStatsResponse{}, err
	}
	for _, record := range records {
		idx := int(record.StartedAt.In(time.Local).Sub(start).Hours() / 24)
		if idx >= 0 && idx < dayCount {
			buckets[idx] = append(buckets[idx], record)
		}
	}

	response := WeeklyStatsResponse{
		Dates:              dates,
		Feeding:            make([]int, dayCount),
		Sleep:              make([]float64, dayCount),
		Diaper:             make([]int, dayCount),
		DailyAverageSleep:  make([]float64, dayCount),
		DaytimeSingleSleep: make([]float64, dayCount),
		DailyAverageMilk:   make([]float64, dayCount),
		TemperatureLatest:  make([]float64, dayCount),
	}
	for i, records := range buckets {
		stats := buildDailyStats(records)
		response.Feeding[i] = stats.Feeding.Count
		response.Sleep[i] = stats.Sleep.TotalHours
		response.Diaper[i] = stats.Diaper.Count
		response.DailyAverageSleep[i] = stats.Sleep.AverageDurationHours
		response.DaytimeSingleSleep[i] = stats.Sleep.DaytimeSingleHours
		response.DailyAverageMilk[i] = stats.Feeding.AverageAmount
		response.TemperatureLatest[i] = stats.Temperature.Latest
	}
	return response, nil
}

func isDaytimeSleep(startedAt time.Time) bool {
	hour := startedAt.In(time.Local).Hour()
	return hour >= 6 && hour < 20
}
