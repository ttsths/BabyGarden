package gredis

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
)

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
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

// setupTestRedis 初始化测试 Redis 客户端
func setupTestRedis() {
	// 设置配置文件路径
	projectRoot := getProjectRoot()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(projectRoot, "config"))

	// 覆盖配置
	viper.Set("redis.host", "121.43.54.203")
	viper.Set("redis.port", 31379)
	viper.Set("redis.password", "123456")
	viper.Set("redis.database", 1)

	config.Setup()
	logger.Setup()
	Setup()
}

// TestSetup 测试 Redis 连接初始化
func TestSetup(t *testing.T) {
	setupTestRedis()

	if RedisClient == nil {
		t.Fatal("Redis 客户端初始化失败")
	}

	// 验证连接可用性
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Redis 连接测试失败: %v", err)
	}

	t.Log("Redis 连接测试通过")
}

// TestSetGet 测试 Set/Get 操作
func TestSetGet(t *testing.T) {
	setupTestRedis()

	key := "test_key"
	value := "test_value"
	expiration := 60 // 秒

	// 测试 Set
	err := Set(key, value, expiration)
	if err != nil {
		t.Fatalf("Set 操作失败: %v", err)
	}

	// 测试 Get
	val, err := Get(key)
	if err != nil {
		t.Fatalf("Get 操作失败: %v", err)
	}

	if val != value {
		t.Errorf("值不匹配: 期望=%s, 实际=%s", value, val)
	}

	// 清理
	_ = Del(key)

	t.Log("Set/Get 测试通过")
}

// TestSetEx 测试 SetEx 操作
func TestSetEx(t *testing.T) {
	setupTestRedis()

	key := "test_setex_key"
	value := "test_setex_value"
	expiration := 5 // 秒

	// 测试 SetEx
	err := SetEx(key, value, expiration)
	if err != nil {
		t.Fatalf("SetEx 操作失败: %v", err)
	}

	// 立即获取验证
	val, err := Get(key)
	if err != nil {
		t.Fatalf("Get 操作失败: %v", err)
	}

	if val != value {
		t.Errorf("值不匹配: 期望=%s, 实际=%s", value, val)
	}

	// 等待过期
	time.Sleep(6 * time.Second)

	// 验证已过期
	_, err = Get(key)
	if err == nil {
		t.Error("SetEx 过期功能未生效，键仍然存在")
	}

	t.Log("SetEx 测试通过")
}

// TestDel 测试 Del 操作
func TestDel(t *testing.T) {
	setupTestRedis()

	key := "test_del_key"
	value := "test_del_value"

	// 先设置
	_ = SetEx(key, value, 60)

	// 验证存在
	exists, err := Exists(key)
	if err != nil {
		t.Fatalf("Exists 操作失败: %v", err)
	}
	if !exists {
		t.Fatal("键不存在，无法测试 Del")
	}

	// 删除
	err = Del(key)
	if err != nil {
		t.Fatalf("Del 操作失败: %v", err)
	}

	// 验证已删除
	exists, err = Exists(key)
	if err != nil {
		t.Fatalf("Exists 操作失败: %v", err)
	}
	if exists {
		t.Error("Del 后键仍然存在")
	}

	t.Log("Del 测试通过")
}

// TestExists 测试 Exists 操作
func TestExists(t *testing.T) {
	setupTestRedis()

	key := "test_exists_key"

	// 测试不存在的键
	exists, err := Exists(key)
	if err != nil {
		t.Fatalf("Exists 操作失败: %v", err)
	}
	if exists {
		t.Error("新键应该不存在")
	}

	// 设置键
	_ = SetEx(key, "value", 60)

	// 测试存在的键
	exists, err = Exists(key)
	if err != nil {
		t.Fatalf("Exists 操作失败: %v", err)
	}
	if !exists {
		t.Error("已存在的键应该返回 true")
	}

	// 清理
	_ = Del(key)

	t.Log("Exists 测试通过")
}

// TestHash 测试 Hash 操作
func TestHash(t *testing.T) {
	setupTestRedis()

	key := "test_hash"
	field := "name"
	value := "test_value"

	// HSet
	err := HSet(key, field, value)
	if err != nil {
		t.Fatalf("HSet 操作失败: %v", err)
	}

	// HGet
	val, err := HGet(key, field)
	if err != nil {
		t.Fatalf("HGet 操作失败: %v", err)
	}

	if val != value {
		t.Errorf("值不匹配: 期望=%s, 实际=%s", value, val)
	}

	// HGetAll
	hashMap, err := HGetAll(key)
	if err != nil {
		t.Fatalf("HGetAll 操作失败: %v", err)
	}

	if len(hashMap) != 1 || hashMap[field] != value {
		t.Errorf("HGetAll 返回值不正确: %v", hashMap)
	}

	// HDel 删除字段
	deleted, err := HDel(key, field)
	if err != nil {
		t.Fatalf("HDel 操作失败: %v", err)
	}
	if deleted != 1 {
		t.Errorf("删除的字段数不正确: 期望=1, 实际=%d", deleted)
	}

	// 清理
	_ = Del(key)

	t.Log("Hash 测试通过")
}

// TestList 测试 List 操作
func TestList(t *testing.T) {
	setupTestRedis()

	key := "test_list"
	values := []interface{}{"item1", "item2", "item3"}

	// LPush
	err := LPush(key, values...)
	if err != nil {
		t.Fatalf("LPush 操作失败: %v", err)
	}

	// LLen
	length, err := LLen(key)
	if err != nil {
		t.Fatalf("LLen 操作失败: %v", err)
	}
	if length != int64(len(values)) {
		t.Errorf("列表长度不匹配: 期望=%d, 实际=%d", len(values), length)
	}

	// LRANGE
	items, err := LRange(key, 0, -1)
	if err != nil {
		t.Fatalf("LRange 操作失败: %v", err)
	}

	// LPush 后进先出，原始顺序是 [item1, item2, item3]，压入后是 [item3, item2, item1]
	if len(items) != len(values) {
		t.Errorf("列表项数量不匹配: 期望=%d, 实际=%d", len(values), len(items))
	}
	// 验证所有项都存在（顺序可能不同）
	foundItems := make(map[string]bool)
	for _, item := range items {
		foundItems[item] = true
	}
	for _, v := range values {
		if !foundItems[v.(string)] {
			t.Errorf("列表缺少项: %v", v)
		}
	}

	// LPop
	poped, err := LPop(key)
	if err != nil {
		t.Fatalf("LPop 操作失败: %v", err)
	}
	// LPush 是后进先出，item3 最后压入，最先弹出
	if poped != "item3" {
		t.Errorf("弹出的值不匹配（LPush 后进先出）: 期望=item3, 实际=%s", poped)
	}

	// 清理
	_ = Del(key)

	t.Log("List 测试通过")
}

// TestIncrDecr 测试自增自减操作
func TestIncrDecr(t *testing.T) {
	setupTestRedis()

	key := "test_counter"

	// 初始值
	err := SetEx(key, "0", 60)
	if err != nil {
		t.Fatalf("SetEx 操作失败: %v", err)
	}

	// Incr
	newVal, err := Incr(key)
	if err != nil {
		t.Fatalf("Incr 操作失败: %v", err)
	}
	if newVal != 1 {
		t.Errorf("自增后值不正确: 期望=1, 实际=%d", newVal)
	}

	// Decr
	newVal, err = Decr(key)
	if err != nil {
		t.Fatalf("Decr 操作失败: %v", err)
	}
	if newVal != 0 {
		t.Errorf("自减后值不正确: 期望=0, 实际=%d", newVal)
	}

	// 清理
	_ = Del(key)

	t.Log("Incr/Decr 测试通过")
}

// TestExpire 测试过期时间设置
func TestExpire(t *testing.T) {
	setupTestRedis()

	key := "test_expire_key"
	value := "test_expire_value"

	// 设置键（无过期时间）
	err := Set(key, value, 0)
	if err != nil {
		t.Fatalf("Set 操作失败: %v", err)
	}

	// 设置过期时间
	err = Expire(key, 5)
	if err != nil {
		t.Fatalf("Expire 操作失败: %v", err)
	}

	// 等待过期
	time.Sleep(6 * time.Second)

	// 验证已过期
	_, err = Get(key)
	if err == nil {
		t.Error("Expire 后键仍然存在")
	}

	t.Log("Expire 测试通过")
}

// TestCodeOps 测试验证码专用操作
func TestCodeOps(t *testing.T) {
	setupTestRedis()

	phone := "13800138000"
	code := "123456"

	// 设置验证码
	err := SetCode(phone, code)
	if err != nil {
		t.Fatalf("SetCode 操作失败: %v", err)
	}

	// 获取验证码
	storedCode, err := GetCode(phone)
	if err != nil {
		t.Fatalf("GetCode 操作失败: %v", err)
	}

	if storedCode != code {
		t.Errorf("验证码不匹配: 期望=%s, 实际=%s", code, storedCode)
	}

	// 检查验证码是否存在
	exists, err := IsCodeExists(phone)
	if err != nil {
		t.Fatalf("IsCodeExists 操作失败: %v", err)
	}
	if !exists {
		t.Error("验证码应该存在")
	}

	// 删除验证码
	err = DelCode(phone)
	if err != nil {
		t.Fatalf("DelCode 操作失败: %v", err)
	}

	// 验证已删除
	storedCode, err = GetCode(phone)
	if err == nil || storedCode != "" {
		t.Error("DelCode 后验证码仍然存在")
	}

	t.Log("验证码操作测试通过")
}

// TestTokenBlacklist 测试 Token 黑名单操作
func TestTokenBlacklist(t *testing.T) {
	setupTestRedis()

	jti := "test_jti_12345"
	ttl := 300 * time.Second // 5分钟

	// 添加到黑名单
	err := AddTokenToBlacklist(jti, ttl)
	if err != nil {
		t.Fatalf("AddTokenToBlacklist 操作失败: %v", err)
	}

	// 检查是否在黑名单中
	blacklisted, err := IsTokenBlacklisted(jti)
	if err != nil {
		t.Fatalf("IsTokenBlacklisted 操作失败: %v", err)
	}
	if !blacklisted {
		t.Error("Token 应该在黑名单中")
	}

	// 从黑名单移除
	err = RemoveTokenFromBlacklist(jti)
	if err != nil {
		t.Fatalf("RemoveTokenFromBlacklist 操作失败: %v", err)
	}

	// 再次检查
	blacklisted, err = IsTokenBlacklisted(jti)
	if err != nil {
		t.Fatalf("IsTokenBlacklisted 操作失败: %v", err)
	}
	if blacklisted {
		t.Error("Token 移除后仍在黑名单中")
	}

	t.Log("Token 黑名单测试通过")
}
