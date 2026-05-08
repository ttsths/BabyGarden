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

// GetDailyStats returns daily new user/baby/record counts for last N days.
// GET /api/v1/admin/stats/daily?days=30
//
// Optimized: 3 GROUP BY queries instead of 3*N individual count queries.
// Old approach took 90 round-trips for 30 days; new takes 3.
func GetDailyStats(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if v := atoi(d); v > 0 && v <= 90 {
			days = v
		}
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -(days - 1))
	endDate := startDate.AddDate(0, 0, days).Add(-1 * time.Second)

	type dailyItem struct {
		Date       string `json:"date"`
		NewUsers   int64  `json:"new_users"`
		NewBabies  int64  `json:"new_babies"`
		NewRecords int64  `json:"new_records"`
	}

	// Single GROUP BY query per model — 3 queries instead of 3*days.
	countByDate := func(modelPtr any, col string) map[string]int64 {
		type row struct {
			Date string `gorm:"column:date"`
			Cnt  int64  `gorm:"column:cnt"`
		}
		var rows []row
		mysql.DB.Model(modelPtr).
			Select("DATE("+col+") as date, COUNT(*) as cnt").
			Where(col+" >= ? AND "+col+" <= ?", startDate, endDate).
			Group("DATE(" + col + ")").
			Find(&rows)
		m := make(map[string]int64, len(rows))
		for _, r := range rows {
			m[r.Date] = r.Cnt
		}
		return m
	}

	userByDate := countByDate(&model.User{}, "created_at")
	babyByDate := countByDate(&model.Baby{}, "created_at")
	recordByDate := countByDate(&model.Record{}, "created_at")

	items := make([]dailyItem, days)
	for i := 0; i < days; i++ {
		dateStr := startDate.AddDate(0, 0, i).Format("2006-01-02")
		items[i] = dailyItem{
			Date:       dateStr,
			NewUsers:   userByDate[dateStr],
			NewBabies:  babyByDate[dateStr],
			NewRecords: recordByDate[dateStr],
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: items,
	})
}
