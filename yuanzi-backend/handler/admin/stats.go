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

	var userCounts, babyCounts, recordCounts []dateCount

	mysql.DB.Model(&model.User{}).
		Select("DATE(created_at) as date, COUNT(*) as cnt").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Group("DATE(created_at)").
		Scan(&userCounts)

	mysql.DB.Model(&model.Baby{}).
		Select("DATE(created_at) as date, COUNT(*) as cnt").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Group("DATE(created_at)").
		Scan(&babyCounts)

	mysql.DB.Model(&model.Record{}).
		Select("DATE(created_at) as date, COUNT(*) as cnt").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Group("DATE(created_at)").
		Scan(&recordCounts)

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
