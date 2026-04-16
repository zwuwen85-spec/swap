package svc

import (
	"campus-swap-shop/internal/model"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// initDB 初始化数据库
func initDB(v *viper.Viper, log *zap.Logger) (*gorm.DB, error) {
	host := v.GetString("mysql.host")
	port := v.GetString("mysql.port")
	database := v.GetString("mysql.database")
	username := v.GetString("mysql.username")
	password := v.GetString("mysql.password")
	maxIdleConns := v.GetInt("mysql.max_idle_conns")
	maxOpenConns := v.GetInt("mysql.max_open_conns")
	connMaxLifetime := v.GetInt("mysql.conn_max_lifetime")

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)

	// GORM配置 - 使用默认日志
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Error("连接MySQL失败", zap.Error(err))
		return nil, err
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("获取数据库连接失败", zap.Error(err))
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		log.Error("数据库连接测试失败", zap.Error(err))
		return nil, err
	}

	log.Info("MySQL连接成功",
		zap.String("host", host),
		zap.String("database", database),
	)

	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Goods{},
		&model.Exchange{},
		&model.Message{},
		&model.Favorite{},
		&model.Comment{},
		&model.Notification{},
		&model.Report{},
	)
	if err != nil {
		log.Error("数据库迁移失败", zap.Error(err))
		return nil, err
	}

	log.Info("数据库表结构迁移成功")

	return db, nil
}
