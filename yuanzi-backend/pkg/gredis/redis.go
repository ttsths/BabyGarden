package gredis

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()
var ErrRedisNotConnected = errors.New("redis: not connected")

// connected 检查 Redis 是否已连接
func connected() bool {
	return RedisClient != nil
}

// IsConnected 公开的连接检查，供外部使用
func IsConnected() bool {
	return connected()
}

// Setup 初始化 Redis 连接
func Setup() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.Database,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		logger.Error("Redis connect failed", logger.Err(err))
		RedisClient = nil
		return
	}
	logger.Info("Redis connected successfully",
		logger.String("host", config.GlobalConfig.Redis.Host),
		logger.Int("port", config.GlobalConfig.Redis.Port),
	)
}

// === String 操作 ===

// Set 设置键值
func Set(key string, value interface{}, expiration int) error {
	if !connected() {
		return ErrRedisNotConnected
	}
	ttl := time.Duration(expiration) * time.Second
	return RedisClient.Set(Ctx, key, value, ttl).Err()
}

// Get 获取键值
func Get(key string) (string, error) {
	if !connected() {
		return "", ErrRedisNotConnected
	}
	return RedisClient.Get(Ctx, key).Result()
}

// Del 删除键
func Del(key string) error {
	if !connected() {
		return ErrRedisNotConnected
	}
	return RedisClient.Del(Ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	if !connected() {
		return false, ErrRedisNotConnected
	}
	n, err := RedisClient.Exists(Ctx, key).Result()
	return n > 0, err
}

// SetEx 设置键值并指定过期时间（秒）
func SetEx(key string, value interface{}, seconds int) error {
	if !connected() {
		return ErrRedisNotConnected
	}
	return RedisClient.SetEx(Ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// Incr 自增
func Incr(key string) (int64, error) {
	if !connected() {
		return 0, ErrRedisNotConnected
	}
	return RedisClient.Incr(Ctx, key).Result()
}

// IncrBy 自增指定值
func IncrBy(key string, value int64) (int64, error) {
	if !connected() {
		return 0, ErrRedisNotConnected
	}
	return RedisClient.IncrBy(Ctx, key, value).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	if !connected() {
		return 0, ErrRedisNotConnected
	}
	return RedisClient.Decr(Ctx, key).Result()
}

// Expire 设置过期时间
func Expire(key string, seconds int) error {
	if !connected() {
		return ErrRedisNotConnected
	}
	return RedisClient.Expire(Ctx, key, time.Duration(seconds)*time.Second).Err()
}

// === Hash 操作 ===

// HSet 设置哈希字段值
func HSet(key string, field string, value interface{}) error {
	return RedisClient.HSet(Ctx, key, field, value).Err()
}

// HGet 获取哈希字段值
func HGet(key string, field string) (string, error) {
	return RedisClient.HGet(Ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(key string) (map[string]string, error) {
	return RedisClient.HGetAll(Ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) (int64, error) {
	return RedisClient.HDel(Ctx, key, fields...).Result()
}

// HExists 检查哈希字段是否存在
func HExists(key string, field string) (bool, error) {
	return RedisClient.HExists(Ctx, key, field).Result()
}

// === List 操作 ===

// LPush 向列表头部添加元素
func LPush(key string, values ...interface{}) error {
	return RedisClient.LPush(Ctx, key, values...).Err()
}

// RPush 向列表尾部添加元素
func RPush(key string, values ...interface{}) error {
	return RedisClient.RPush(Ctx, key, values...).Err()
}

// LPop 从列表头部弹出元素
func LPop(key string) (string, error) {
	return RedisClient.LPop(Ctx, key).Result()
}

// RPop 从列表尾部弹出元素
func RPop(key string) (string, error) {
	return RedisClient.RPop(Ctx, key).Result()
}

// LLen 获取列表长度
func LLen(key string) (int64, error) {
	return RedisClient.LLen(Ctx, key).Result()
}

// LRange 获取列表指定范围的元素
func LRange(key string, start, end int64) ([]string, error) {
	return RedisClient.LRange(Ctx, key, start, end).Result()
}

// LIndex 获取列表指定索引的元素
func LIndex(key string, index int64) (string, error) {
	return RedisClient.LIndex(Ctx, key, index).Result()
}

// LTrim 修剪列表
func LTrim(key string, start, end int64) error {
	return RedisClient.LTrim(Ctx, key, start, end).Err()
}

// === 验证码专用操作 ===

// SetCode 设置验证码（5分钟过期）
func SetCode(phone string, code string) error {
	key := fmt.Sprintf("code:%s", phone)
	return SetEx(key, code, 300)
}

// GetCode 获取验证码
func GetCode(phone string) (string, error) {
	key := fmt.Sprintf("code:%s", phone)
	return Get(key)
}

// DelCode 删除验证码
func DelCode(phone string) error {
	key := fmt.Sprintf("code:%s", phone)
	return Del(key)
}

// IsCodeExists 检查验证码是否存在
func IsCodeExists(phone string) (bool, error) {
	key := fmt.Sprintf("code:%s", phone)
	return Exists(key)
}

// === Token 黑名单操作 ===

// AddTokenToBlacklist 将 Token 加入黑名单（基于 JTI）
func AddTokenToBlacklist(jti string, ttl time.Duration) error {
	key := fmt.Sprintf("token:blacklist:%s", jti)
	return SetEx(key, "1", int(ttl.Seconds())+60)
}

// IsTokenBlacklisted 检查 Token 是否在黑名单中
func IsTokenBlacklisted(jti string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", jti)
	return Exists(key)
}

// RemoveTokenFromBlacklist 从黑名单移除 Token
func RemoveTokenFromBlacklist(jti string) error {
	key := fmt.Sprintf("token:blacklist:%s", jti)
	return Del(key)
}

// Publish 发布消息
func Publish(channel string, message string) error {
	if !connected() {
		return ErrRedisNotConnected
	}
	return RedisClient.Publish(Ctx, channel, message).Err()
}

// Subscribe 订阅频道
func Subscribe(channel string) *redis.PubSub {
	if !connected() {
		return nil
	}
	return RedisClient.Subscribe(Ctx, channel)
}
