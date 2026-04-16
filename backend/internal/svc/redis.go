package svc

import (
	"campus-swap-shop/pkg/utils"
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// initRedis 初始化Redis
func initRedis(v *viper.Viper, log *zap.Logger) (*utils.RedisClient, error) {
	host := v.GetString("redis.host")
	port := v.GetString("redis.port")
	password := v.GetString("redis.password")
	db := v.GetInt("redis.db")
	poolSize := v.GetInt("redis.pool_size")

	// 构建地址
	addr := host + ":" + port

	// 创建Redis客户端
	redisClient := utils.NewRedisClient(addr, password, db)

	// 测试连接
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Error("Redis连接失败", zap.Error(err))
		return nil, err
	}

	log.Info("Redis连接成功",
		zap.String("host", host),
		zap.Int("db", db),
		zap.Int("pool_size", poolSize),
	)

	return redisClient, nil
}
