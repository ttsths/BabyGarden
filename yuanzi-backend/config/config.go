package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OSS      OSSConfig      `mapstructure:"oss"`
	R2       R2Config       `mapstructure:"r2"`
	Storage  StorageConfig  `mapstructure:"storage"`
	SMS      SMSConfig      `mapstructure:"sms"`
	AI       AIConfig       `mapstructure:"ai"`
	Push     PushConfig     `mapstructure:"push"`
}

type ServerConfig struct {
	RunMode      string `mapstructure:"run_mode"`
	HttpPort     int    `mapstructure:"http_port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Type        string `mapstructure:"type"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	TablePrefix string `mapstructure:"table_prefix"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	ExpireDuration time.Duration `mapstructure:"expire_duration"`
	RefreshDays    int           `mapstructure:"refresh_days"`
}

type OSSConfig struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	CdnDomain       string `mapstructure:"cdn_domain"`
	CallbackSecret  string `mapstructure:"callback_secret"`
}

type R2Config struct {
	AccountID       string `mapstructure:"account_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Bucket          string `mapstructure:"bucket"`
	PublicURL       string `mapstructure:"public_url"`
}

type StorageConfig struct {
	Provider string `mapstructure:"provider"` // "oss" or "r2"
}

type SMSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	SignName        string `mapstructure:"sign_name"`
	TemplateCode    string `mapstructure:"template_code"`
}

type AIConfig struct {
	DashScopeAPIKey    string `mapstructure:"dashscope_api_key"`
	NLSAppKey          string `mapstructure:"nls_app_key"`
	NLSEndpoint        string `mapstructure:"nls_endpoint"`
	NLSAccessKeyID     string `mapstructure:"nls_access_key_id"`
	NLSAccessKeySecret string `mapstructure:"nls_access_key_secret"`
}

type PushConfig struct {
	APNs APNsConfig `mapstructure:"apns"`
}

type APNsConfig struct {
	TeamID         string `mapstructure:"team_id"`
	KeyID          string `mapstructure:"key_id"`
	BundleID       string `mapstructure:"bundle_id"`
	KeyPath        string `mapstructure:"key_path"`
	UseSandbox     bool   `mapstructure:"use_sandbox"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

var GlobalConfig Config

func Setup() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 读取环境变量（自动覆盖配置文件中的值）
	// SetEnvKeyReplacer 让 R2_ACCOUNT_ID 这样的环境变量能映射到 r2.account_id
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found, using defaults and env vars: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("storage.provider", "r2")

	viper.SetDefault("server.run_mode", "debug")
	viper.SetDefault("server.http_port", 8080)
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)

	viper.SetDefault("database.type", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.name", "yuanzi")
	viper.SetDefault("database.max_idle_conn", 10)
	viper.SetDefault("database.max_open_conn", 100)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.database", 0)

	viper.SetDefault("jwt.secret", "yuanzi-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_duration", 2*time.Hour)
	viper.SetDefault("jwt.refresh_days", 7)

	viper.SetDefault("ai.nls_access_key_id", "")
	viper.SetDefault("ai.nls_access_key_secret", "")

	viper.SetDefault("push.apns.team_id", "")
	viper.SetDefault("push.apns.key_id", "")
	viper.SetDefault("push.apns.bundle_id", "")
	viper.SetDefault("push.apns.key_path", "")
	viper.SetDefault("push.apns.use_sandbox", true)
	viper.SetDefault("push.apns.timeout_seconds", 10)
}
