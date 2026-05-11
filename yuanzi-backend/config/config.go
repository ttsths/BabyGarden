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
	ProviderChain []string          `mapstructure:"provider_chain"`
	TimeoutSec    int               `mapstructure:"timeout_seconds"`
	MaxRetries    int               `mapstructure:"max_retries_per_provider"`
	EnableFallback bool             `mapstructure:"enable_fallback"`
	Safety        AISafetyConfig    `mapstructure:"safety"`
	Quota         AIQuotaConfig     `mapstructure:"quota"`
	GrokAI        GrokAIConfig      `mapstructure:"grokai"`
	Cloudflare    CloudflareAIConfig `mapstructure:"cloudflare"`
	DashScope     DashScopeAIConfig `mapstructure:"dashscope"`
	CLIProxyAPI   CLIProxyAIConfig  `mapstructure:"cliproxyapi"`

	// 旧字段（保留兼容）
	DashScopeAPIKey    string `mapstructure:"dashscope_api_key"`
	NLSAppKey          string `mapstructure:"nls_app_key"`
	NLSEndpoint        string `mapstructure:"nls_endpoint"`
	NLSAccessKeyID     string `mapstructure:"nls_access_key_id"`
	NLSAccessKeySecret string `mapstructure:"nls_access_key_secret"`
}

type AISafetyConfig struct {
	SystemPrompt     string `mapstructure:"system_prompt"`
	MaxPromptChars   int    `mapstructure:"max_prompt_chars"`
	MaxOutputTokens  int    `mapstructure:"max_output_tokens"`
}

type AIQuotaConfig struct {
	RedisPrefix          string `mapstructure:"redis_prefix"`
	PerUserDailyLimit    int    `mapstructure:"per_user_daily_limit"`
	PerFamilyDailyLimit  int    `mapstructure:"per_family_daily_limit"`
}

type GrokAIConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	BaseURL        string `mapstructure:"base_url"`
	APIKey         string `mapstructure:"api_key"`
	Model          string `mapstructure:"model"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type CloudflareAIConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	AccountID         string `mapstructure:"account_id"`
	APIToken          string `mapstructure:"api_token"`
	GatewayID         string `mapstructure:"gateway_id"`
	UseGateway        bool   `mapstructure:"use_gateway"`
	Model             string `mapstructure:"model"`
	DailyNeuronBudget int    `mapstructure:"daily_neuron_budget"`
	HardNeuronBudget  int    `mapstructure:"hard_neuron_budget"`
	TimeoutSeconds    int    `mapstructure:"timeout_seconds"`
	CacheTTLSeconds   int    `mapstructure:"cache_ttl_seconds"`
}

type DashScopeAIConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	APIKey         string `mapstructure:"api_key"`
	Model          string `mapstructure:"model"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type CLIProxyAIConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	BaseURL        string `mapstructure:"base_url"`
	APIKey         string `mapstructure:"api_key"`
	Model          string `mapstructure:"model"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
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

	viper.SetDefault("ai.provider_chain", []string{"grokai", "cloudflare_workers_ai", "dashscope", "cliproxyapi"})
	viper.SetDefault("ai.timeout_seconds", 45)
	viper.SetDefault("ai.max_retries_per_provider", 1)
	viper.SetDefault("ai.enable_fallback", true)
	viper.SetDefault("ai.safety.system_prompt", "你是 BabyGarden 的育儿助手。回答应谨慎、实用、非诊断。涉及疾病、用药、急症时建议及时咨询医生或急诊。")
	viper.SetDefault("ai.safety.max_prompt_chars", 12000)
	viper.SetDefault("ai.safety.max_output_tokens", 1200)
	viper.SetDefault("ai.quota.redis_prefix", "babygarden:ai")
	viper.SetDefault("ai.quota.per_user_daily_limit", 50)
	viper.SetDefault("ai.quota.per_family_daily_limit", 200)
	viper.SetDefault("ai.grokai.enabled", false)
	viper.SetDefault("ai.grokai.model", "grok-4.1-fast")
	viper.SetDefault("ai.grokai.timeout_seconds", 45)
	viper.SetDefault("ai.cloudflare.enabled", false)
	viper.SetDefault("ai.cloudflare.use_gateway", false)
	viper.SetDefault("ai.cloudflare.model", "@cf/moonshotai/kimi-k2.6")
	viper.SetDefault("ai.cloudflare.daily_neuron_budget", 9000)
	viper.SetDefault("ai.cloudflare.hard_neuron_budget", 9800)
	viper.SetDefault("ai.cloudflare.timeout_seconds", 45)
	viper.SetDefault("ai.cloudflare.cache_ttl_seconds", 300)
	viper.SetDefault("ai.dashscope.enabled", true)
	viper.SetDefault("ai.dashscope.model", "qwen-plus")
	viper.SetDefault("ai.dashscope.timeout_seconds", 45)
	viper.SetDefault("ai.cliproxyapi.enabled", false)
	viper.SetDefault("ai.cliproxyapi.model", "gpt-5.5")
	viper.SetDefault("ai.cliproxyapi.timeout_seconds", 60)
	viper.SetDefault("ai.nls_access_key_id", "")
	viper.SetDefault("ai.nls_access_key_secret", "")

	viper.SetDefault("push.apns.team_id", "")
	viper.SetDefault("push.apns.key_id", "")
	viper.SetDefault("push.apns.bundle_id", "")
	viper.SetDefault("push.apns.key_path", "")
	viper.SetDefault("push.apns.use_sandbox", true)
	viper.SetDefault("push.apns.timeout_seconds", 10)
}
