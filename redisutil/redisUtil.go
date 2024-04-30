package redisutil

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// RedisClient redis客户端结构体
type RedisClient struct {
	// redis客户端
	redisClient *redis.Client
}

// NewRedisClient 创建redis客户端
// @param addr redis地址
// @param port redis端口
// @param password redis密码
// @param db redis数据库
// @return client,isConnect redis客户端,连接是否成功
func NewRedisClient(addr string, port int, password string, db int) (client *RedisClient, isConnect bool) {
	client = &RedisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", addr, port),
			Password: password,
			DB:       db,
		}),
	}
	// 检测连接
	if client.isConnect() {
		return client, true
	}
	return nil, false
}

// isConnect 检测连接
// @return bool 连接是否成功
func (client *RedisClient) isConnect() bool {
	_, err := client.redisClient.Ping().Result()
	if err == nil {
		return true
	}
	fmt.Println("ping error", err.Error())
	return false
}

// GetString 根据key获取string值
// @param key 键
// @return string 值
func (client *RedisClient) GetString(key string) string {
	result, err := client.redisClient.Get(key).Result()
	if err == nil {
		return result
	}
	return ""
}

// HasKey 检测key是否存在
// @param key 键
// @return bool key是否存在
func (client *RedisClient) HasKey(key string) bool {
	return client.GetString(key) != ""
}

// GetInt 根据key获取int值
// @param key 键
// @return int 值
func (client *RedisClient) GetInt(key string) int {
	result, err := client.redisClient.Get(key).Int()
	if err == nil {
		return result
	}
	return 0
}

// GetInt64 根据key获取int64值
// @param key 键
// @return int64 值
func (client *RedisClient) GetInt64(key string) int64 {
	result, err := client.redisClient.Get(key).Int64()
	if err == nil {
		return result
	}
	return 0
}

// Set 设置值
// @param key 键
// @param value 值
// @return bool 设置是否成功
func (client *RedisClient) Set(key string, value interface{}) bool {
	result, err := client.redisClient.Set(key, value, 0).Result()
	if err == nil && result == "OK" {
		return true
	}
	return false
}

// SetWithExpire 设置值并设置过期时间
// @param key 键
// @param value 值
// @param expire 过期时间(秒)
// @return bool 设置是否成功
func (client *RedisClient) SetWithExpire(key string, value interface{}, expire int64) bool {
	result, err := client.redisClient.Set(key, value, time.Second*time.Duration(expire)).Result()
	if err == nil && result == "OK" {
		return true
	}
	return false
}
