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
	type dailyItem struct {
		Date       string `json:"date"`
		NewUsers   int64  `json:"new_users"`
		NewBabies  int64  `json:"new_babies"`
		NewRecords int64  `json:"new_records"`
	}

	items := make([]dailyItem, days)
	for i := 0; i < days; i++ {
		day := now.AddDate(0, 0, -i)
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
		end := start.Add(24 * time.Hour)
		dateStr := start.Format("2006-01-02")

		var nc, bc, rc int64
		mysql.DB.Model(&model.User{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&nc)
		mysql.DB.Model(&model.Baby{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&bc)
		mysql.DB.Model(&model.Record{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&rc)

		items[i] = dailyItem{
			Date:       dateStr,
			NewUsers:   nc,
			NewBabies:  bc,
			NewRecords: rc,
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: items,
	})
}
