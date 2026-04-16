package svc

import (
	"campus-swap-shop/pkg/logger"
	"campus-swap-shop/pkg/utils"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config   *viper.Viper
	DB       *gorm.DB
	Redis    *utils.RedisClient
	Logger   *zap.Logger
	JwtSecret string
}

// NewServiceContext 创建服务上下文
func NewServiceContext() (*ServiceContext, error) {
	// 加载配置
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// 初始化日志
	log := logger.NewLogger(*config)

	// 初始化MySQL
	db, err := initDB(config, log)
	if err != nil {
		log.Fatal("初始化MySQL失败", zap.Error(err))
		return nil, err
	}

	// 初始化Redis
	redisClient, err := initRedis(config, log)
	if err != nil {
		log.Fatal("初始化Redis失败", zap.Error(err))
		return nil, err
	}

	return &ServiceContext{
		Config:    config,
		DB:        db,
		Redis:     redisClient,
		Logger:    log,
		JwtSecret: config.GetString("jwt.secret"),
	}, nil
}

// loadConfig 加载配置文件
func loadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}
