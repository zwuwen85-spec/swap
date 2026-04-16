package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端
type RedisClient struct {
	*redis.Client
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{rdb}
}

// Set 设置键值
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

// SetEx 设置键值并指定过期时间（秒）
func (r *RedisClient) SetEx(ctx context.Context, key string, value interface{}, expiration int64) error {
	return r.Client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

// Get 获取值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Incr 递增
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration int) error {
	return r.Client.Expire(ctx, key, 0).Err()
}
