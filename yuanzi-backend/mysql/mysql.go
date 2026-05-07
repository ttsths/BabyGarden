package mysql

import (
	"fmt"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Setup 初始化数据库连接
func Setup() {
	cfg := config.GlobalConfig.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s&readTimeout=10s&writeTimeout=10s&tls=skip-verify&allowCleartextPasswords=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		logger.Fatal("Failed to connect to database", logger.Err(err))
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB", logger.Err(err))
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

	logger.Info("Database connected successfully",
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.String("database", cfg.Name),
	)
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}
