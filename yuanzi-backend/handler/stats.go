package handler

import (
	"encoding/json"
	"errors"
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

// GetStatsSummary 按日/周/月/自定义区间统计首页趋势。
// @Summary 区间统计
// @Description 获取每日睡眠、白天单次睡眠、喝奶量趋势
// @Tags 统计
// @Accept json
// @Produce json
// @Security Bearer
// @Param baby_id query string true "宝宝ID"
// @Param range query string false "day/week/month/custom" default(week)
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Success 200 {object} model.Response{data=StatsSummaryResponse}
// @Router /api/v1/stats/summary [get]
func GetStatsSummary(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	rangeType := c.DefaultQuery("range", "week")
	startDate, endDate, err := parseStatsRange(rangeType, c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: err.Error()})
		return
	}

	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days <= 0 || days > 366 {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "统计区间需在1到366天之间"})
		return
	}

	dates := make([]string, 0, days)
	sleepHours := make([]float64, days)
	daytimeSleep := make([]float64, days)
	milkAmount := make([]int, days)
	daytimeSleepCounts := make([]int, days)
	for i := 0; i < days; i++ {
		dates = append(dates, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}

	var records []model.Record
	queryEnd := endDate.Add(24 * time.Hour)
	if err := mysql.DB.Where("baby_id = ? AND started_at >= ? AND started_at < ?", babyID, startDate, queryEnd).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	for _, record := range records {
		idx := int(record.StartedAt.In(time.Local).Sub(startDate).Hours() / 24)
		if idx < 0 || idx >= days {
			continue
		}
		switch record.Type {
		case model.RecordTypeFeeding:
			var payload model.FeedingContent
			_ = json.Unmarshal(record.Content, &payload)
			if payload.Amount > 0 {
				milkAmount[idx] += payload.Amount
			}
		case model.RecordTypeSleep:
			if !record.EndedAt.Valid {
				continue
			}
			duration := record.EndedAt.Time.Sub(record.StartedAt).Hours()
			if duration < 0 {
				continue
			}
			sleepHours[idx] += duration
			hour := record.StartedAt.In(time.Local).Hour()
			if hour >= 6 && hour < 18 {
				daytimeSleep[idx] += duration
				daytimeSleepCounts[idx]++
			}
		}
	}

	totalSleep := 0.0
	totalMilk := 0
	totalDaytimeSleep := 0.0
	daytimeDays := 0
	for i := range dates {
		totalSleep += sleepHours[i]
		totalMilk += milkAmount[i]
		if daytimeSleepCounts[i] > 0 {
			daytimeSleep[i] = daytimeSleep[i] / float64(daytimeSleepCounts[i])
			totalDaytimeSleep += daytimeSleep[i]
			daytimeDays++
		}
	}

	summary := StatsSummaryTotals{
		AvgDailySleepHours: round2(totalSleep / float64(days)),
		AvgDailyMilkAmount: totalMilk / days,
	}
	if daytimeDays > 0 {
		summary.AvgDaytimeSingleSleepHours = round2(totalDaytimeSleep / float64(daytimeDays))
	}
	for i := range sleepHours {
		sleepHours[i] = round2(sleepHours[i])
		daytimeSleep[i] = round2(daytimeSleep[i])
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: StatsSummaryResponse{
		Range:                   rangeType,
		StartDate:               startDate.Format("2006-01-02"),
		EndDate:                 endDate.Format("2006-01-02"),
		Dates:                   dates,
		DailyAvgSleepHours:      sleepHours,
		DaytimeSingleSleepHours: daytimeSleep,
		DailyAvgMilkAmount:      milkAmount,
		Summary:                 summary,
	}})
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

func parseStatsRange(rangeType, startValue, endValue string) (time.Time, time.Time, error) {
	today, err := parseStatsDate(endValue)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	switch rangeType {
	case "day":
		if startValue != "" {
			start, err := parseStatsDate(startValue)
			return start, start, err
		}
		return today, today, nil
	case "week", "":
		return today.AddDate(0, 0, -6), today, nil
	case "month":
		start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)
		return start, today, nil
	case "custom":
		if startValue == "" || endValue == "" {
			return time.Time{}, time.Time{}, errors.New("自定义区间需要 start_date 和 end_date")
		}
		start, err := parseStatsDate(startValue)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("开始日期格式错误")
		}
		end, err := parseStatsDate(endValue)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("结束日期格式错误")
		}
		if end.Before(start) {
			return time.Time{}, time.Time{}, errors.New("结束日期不能早于开始日期")
		}
		return start, end, nil
	default:
		return time.Time{}, time.Time{}, errors.New("统计区间类型错误")
	}
}

func round2(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
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

type StatsSummaryResponse struct {
	Range                   string             `json:"range"`
	StartDate               string             `json:"start_date"`
	EndDate                 string             `json:"end_date"`
	Dates                   []string           `json:"dates"`
	DailyAvgSleepHours      []float64          `json:"daily_avg_sleep_hours"`
	DaytimeSingleSleepHours []float64          `json:"daytime_single_sleep_hours"`
	DailyAvgMilkAmount      []int              `json:"daily_avg_milk_amount"`
	Summary                 StatsSummaryTotals `json:"summary"`
}

type StatsSummaryTotals struct {
	AvgDailySleepHours         float64 `json:"avg_daily_sleep_hours"`
	AvgDaytimeSingleSleepHours float64 `json:"avg_daytime_single_sleep_hours"`
	AvgDailyMilkAmount         int     `json:"avg_daily_milk_amount"`
}
