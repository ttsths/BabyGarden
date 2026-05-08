package admin

import (
	"net/http"
	"time"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

	"github.com/gin-gonic/gin"
)

// GetStatsOverview returns system overview statistics.
// GET /api/v1/admin/stats/overview
func GetStatsOverview(c *gin.Context) {
	var userCount, familyCount, babyCount, photoCount, recordCount int64

	mysql.DB.Model(&model.User{}).Where("status > 0").Count(&userCount)
	mysql.DB.Model(&model.Family{}).Count(&familyCount)
	mysql.DB.Model(&model.Baby{}).Count(&babyCount)
	mysql.DB.Model(&model.Photo{}).Where("status = ?", model.PhotoStatusActive).Count(&photoCount)
	mysql.DB.Model(&model.Record{}).Count(&recordCount)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"users":    userCount,
			"families": familyCount,
			"babies":   babyCount,
			"photos":   photoCount,
			"records":  recordCount,
		},
	})
}

// GetDailyStats returns daily new user/baby/record counts for last 30 days.
// Queries are parallelized to avoid cumulative latency exceeding CF Workers timeout.
// GET /api/v1/admin/stats/daily
func GetDailyStats(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if v := atoi(d); v > 0 && v <= 90 {
			days = v
		}
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	startDate := todayStart.AddDate(0, 0, -(days - 1))
	endDate := startDate.AddDate(0, 0, days)

	type dateCount struct {
		Date string `gorm:"column:date"`
		Cnt  int64  `gorm:"column:cnt"`
	}

	type queryResult struct {
		counts []dateCount
		err    error
	}

	// Parallelize the 3 aggregation queries to reduce worst-case latency
	// from 3×T (sequential) to max(T) (parallel).
	ctx := c.Request.Context()
	userCh := make(chan queryResult, 1)
	babyCh := make(chan queryResult, 1)
	recordCh := make(chan queryResult, 1)

	go func() {
		var rows []dateCount
		err := mysql.DB.WithContext(ctx).
			Model(&model.User{}).
			Select("DATE(created_at) as date, COUNT(*) as cnt").
			Where("created_at >= ? AND created_at < ?", startDate, endDate).
			Group("DATE(created_at)").
			Scan(&rows).Error
		userCh <- queryResult{rows, err}
	}()

	go func() {
		var rows []dateCount
		err := mysql.DB.WithContext(ctx).
			Model(&model.Baby{}).
			Select("DATE(created_at) as date, COUNT(*) as cnt").
			Where("created_at >= ? AND created_at < ?", startDate, endDate).
			Group("DATE(created_at)").
			Scan(&rows).Error
		babyCh <- queryResult{rows, err}
	}()

	go func() {
		var rows []dateCount
		err := mysql.DB.WithContext(ctx).
			Model(&model.Record{}).
			Select("DATE(created_at) as date, COUNT(*) as cnt").
			Where("created_at >= ? AND created_at < ?", startDate, endDate).
			Group("DATE(created_at)").
			Scan(&rows).Error
		recordCh <- queryResult{rows, err}
	}()

	// Collect results; all 3 queries run concurrently so total wait ≈ max(query time)
	userRes, babyRes, recordRes := <-userCh, <-babyCh, <-recordCh

	if userRes.err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询用户统计失败"})
		return
	}
	if babyRes.err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询宝宝统计失败"})
		return
	}
	if recordRes.err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询记录统计失败"})
		return
	}

	userCounts := userRes.counts
	babyCounts := babyRes.counts
	recordCounts := recordRes.counts

	userMap := make(map[string]int64, len(userCounts))
	for _, r := range userCounts {
		userMap[r.Date] = r.Cnt
	}
	babyMap := make(map[string]int64, len(babyCounts))
	for _, r := range babyCounts {
		babyMap[r.Date] = r.Cnt
	}
	recordMap := make(map[string]int64, len(recordCounts))
	for _, r := range recordCounts {
		recordMap[r.Date] = r.Cnt
	}

	type dailyItem struct {
		Date       string `json:"date"`
		NewUsers   int64  `json:"new_users"`
		NewBabies  int64  `json:"new_babies"`
		NewRecords int64  `json:"new_records"`
	}

	items := make([]dailyItem, days)
	for i := 0; i < days; i++ {
		date := todayStart.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		items[i] = dailyItem{
			Date:       dateStr,
			NewUsers:   userMap[dateStr],
			NewBabies:  babyMap[dateStr],
			NewRecords: recordMap[dateStr],
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: items,
	})
}
