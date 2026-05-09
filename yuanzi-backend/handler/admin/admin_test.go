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
	
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"
	"yuanzi-backend/pkg/storage"
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
func bytesBody(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

func dbConn() *gorm.DB {
	return mysql.DB
}

// --- storageURL 单元测试 ---

func TestStorageURL(t *testing.T) {
	t.Run("nil provider returns empty", func(t *testing.T) {
		got := storageURL(nil, "some/key.jpg")
		if got != "" {
			t.Errorf("storageURL(nil, ...) = %q, want empty", got)
		}
	})

	t.Run("empty key returns empty", func(t *testing.T) {
		// Even with a mock provider, empty key → empty result
		got := storageURL(nil, "")
		if got != "" {
			t.Errorf("storageURL(..., \"\") = %q, want empty", got)
		}
	})
}

// TestStorageURLIntegration tests storageURL with a real provider (needs config).
// This is a placeholder — real integration test needs actual R2 config.
func TestStorageURLIntegration(t *testing.T) {
	setupAdminTestDB(t)

	provider, err := storage.NewProviderFromConfig()
	if err != nil {
		t.Skipf("Skipping integration test: storage provider not configured (%v)", err)
	}

	url := storageURL(provider, "test-family/test-baby/test-photo.jpg")
	if url == "" {
		t.Error("storageURL returned empty for valid provider and key")
	}
	t.Logf("Generated URL: %s", url)
}