package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"
)

var adminTestSetupOnce sync.Once

func setupAdminTestDB(t *testing.T) {
	t.Helper()
	adminTestSetupOnce.Do(func() {
		projectRoot := getAdminProjectRoot(t)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(filepath.Join(projectRoot, "config"))
		config.Setup()
		logger.Setup()
		mysql.Setup()
		gredis.Setup()
		gin.SetMode(gin.TestMode)
	})
}

func getAdminProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}
	for i := 0; i < 5; i++ {
		configPath := filepath.Join(dir, "config", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("找不到 config/config.yaml")
	return ""
}

// --- Test helpers ---

func uniquePhone(prefix string) string {
	seed := time.Now().UnixNano() + int64(rand.Intn(1000))
	return fmt.Sprintf("13%09d", (seed+int64(len(prefix))*97)%1000000000)
}

func createTestUser(t *testing.T, phone, nickname string, isAdmin int8) model.User {
	t.Helper()
	user := model.User{Phone: phone, Nickname: nickname, Status: 1, IsAdmin: isAdmin}
	if err := mysql.DB.Create(&user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	return user
}

func createTestFamily(t *testing.T, createdBy string, name string) model.Family {
	t.Helper()
	code := fmt.Sprintf("T%07d", rand.Intn(9999999))
	family := model.Family{Name: name, InviteCode: code, CreatedBy: createdBy, StorageLimit: 104857600}
	if err := mysql.DB.Create(&family).Error; err != nil {
		t.Fatalf("创建测试家庭失败: %v", err)
	}
	return family
}

func createTestFamilyMember(t *testing.T, familyID, userID string, role model.FamilyRole) model.FamilyMember {
	t.Helper()
	member := model.FamilyMember{FamilyID: familyID, UserID: userID, Role: role, JoinedAt: time.Now()}
	if err := mysql.DB.Create(&member).Error; err != nil {
		t.Fatalf("创建测试成员失败: %v", err)
	}
	return member
}

// generateAdminToken creates a JWT token for an admin user (for middleware-based tests).
func generateAdminToken(t *testing.T, userID, phone string) string {
	t.Helper()
	claims := middlewareClaims{
		UserID:  userID,
		Phone:   phone,
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "yuanzi-admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		t.Fatalf("生成测试 Token 失败: %v", err)
	}
	return tokenString
}

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("JSON 编码失败: %v", err)
	}
	return data
}

func decodeJSON(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("JSON 解码失败: %v body=%s", err, string(body))
	}
}

// Cleanup helpers

func cleanupBabyRecords(t *testing.T, familyIDs ...string) {
	t.Helper()
	if len(familyIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.Record{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.Baby{}).Error
	_ = mysql.DB.Where("family_id IN ?", familyIDs).Delete(&model.FamilyMember{}).Error
	_ = mysql.DB.Where("id IN ?", familyIDs).Delete(&model.Family{}).Error
}

func cleanupUsers(t *testing.T, userIDs ...string) {
	t.Helper()
	if len(userIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("id IN ?", userIDs).Delete(&model.User{}).Error
}

func mustJSON(t *testing.T, v interface{}) model.JSON {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("JSON 编码失败: %v", err)
	}
	return model.JSON(data)
}

func bytesBody(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

func dbConn() *gorm.DB {
	return mysql.DB
}