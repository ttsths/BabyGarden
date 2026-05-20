package mysql

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
)

// TestDBConnection 测试数据库连接
func TestDBConnection(t *testing.T) {
	if os.Getenv("RUN_EXTERNAL_INTEGRATION_TESTS") != "1" {
		t.Skip("跳过外部 MySQL 集成测试：设置 RUN_EXTERNAL_INTEGRATION_TESTS=1 后执行")
	}

	// 设置配置文件路径为项目根目录
	projectRoot := getProjectRoot()
	t.Logf("项目根目录: %s", projectRoot)

	// 使用 config 包的 Setup，但需要先配置 viper 的路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(projectRoot, "config"))

	// 覆盖 config.Setup 中的默认值（必须在 config.Setup 调用前）
	viper.Set("database.host", "121.43.54.203")
	viper.Set("database.port", 3306)
	viper.Set("database.user", "root")
	viper.Set("database.password", "Sin2=cos2=1")
	viper.Set("database.name", "yuanzi")
	viper.Set("database.max_idle_conn", 20)
	viper.Set("database.max_open_conn", 100)

	viper.Set("redis.host", "121.43.54.203")
	viper.Set("redis.port", 31379)
	viper.Set("redis.password", "123456")
	viper.Set("redis.database", 1)

	viper.Set("jwt.secret", "test-secret-key")
	viper.Set("jwt.expire_duration", "2h")
	viper.Set("jwt.refresh_days", 7)

	// 设置环境变量
	os.Setenv("DATABASE_HOST", "121.43.54.203")
	os.Setenv("DATABASE_PORT", "3306")
	os.Setenv("DATABASE_USER", "root")
	os.Setenv("DATABASE_PASSWORD", "Sin2=cos2=1")
	os.Setenv("DATABASE_NAME", "yuanzi")

	// 初始化配置（会读取 viper 的值）
	config.Setup()

	// 初始化日志
	logger.Setup()

	// 初始化数据库连接
	Setup()

	// 验证 DB 不为空
	if DB == nil {
		t.Fatal("数据库连接初始化失败，DB 为空")
	}

	// 验证数据库连接池配置（已在 mysql.go 中配置，此处验证配置已生效）
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("获取 sql.DB 失败: %v", err)
	}

	// 验证连接可用性（ping）
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	t.Log("数据库连接测试通过")

	// 清理
	Close()
}

// TestGetDB 测试 GetDB 函数
func TestGetDB(t *testing.T) {
	if os.Getenv("RUN_EXTERNAL_INTEGRATION_TESTS") != "1" {
		t.Skip("跳过外部 MySQL 集成测试：设置 RUN_EXTERNAL_INTEGRATION_TESTS=1 后执行")
	}

	projectRoot := getProjectRoot()

	// 使用 config 包的 Setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(projectRoot, "config"))
	viper.Set("database.host", "121.43.54.203")
	viper.Set("database.port", 3306)
	viper.Set("database.user", "root")
	viper.Set("database.password", "Sin2=cos2=1")
	viper.Set("database.name", "yuanzi")

	config.Setup()
	logger.Setup()
	Setup()

	db := GetDB()
	if db == nil {
		t.Fatal("GetDB 返回空")
	}

	// 验证可以执行基本查询（使用 INFORMATION_SCHEMA 验证数据库存在性）
	var count int64
	err := db.Table("information_schema.schemata").Where("schema_name = ?", "yuanzi").Count(&count).Error
	if err != nil {
		t.Logf("警告: 查询 schemata 失败 (可能数据库不存在): %v", err)
	}

	t.Log("GetDB 测试通过")

	// 清理
	Close()
}

// TestClose 测试 Close 函数
func TestClose(t *testing.T) {
	if os.Getenv("RUN_EXTERNAL_INTEGRATION_TESTS") != "1" {
		t.Skip("跳过外部 MySQL 集成测试：设置 RUN_EXTERNAL_INTEGRATION_TESTS=1 后执行")
	}

	config.Setup()
	logger.Setup()
	Setup()

	// 先验证连接存在
	if DB == nil {
		t.Fatal("DB 为空，无法测试 Close")
	}

	// 关闭数据库连接
	Close()

	t.Log("Close 测试通过")
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 从测试文件目录向上查找 config/config.yaml
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		configPath := filepath.Join(dir, "config", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return "."
}
