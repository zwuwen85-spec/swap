package model

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 收藏模型
type Favorite struct {
	ID         int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64 `json:"user_id" gorm:"not null;index:idx_user_id;uniqueIndex:uk_user_goods"`
	GoodsID    int64 `json:"goods_id" gorm:"not null;uniqueIndex:uk_user_goods"`
	CreateTime int64 `json:"create_time" gorm:"not null;index:idx_create_time"`

	// 关联
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Goods *Goods `json:"goods,omitempty" gorm:"foreignKey:GoodsID"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorite"
}

// BeforeCreate 创建前钩子
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.CreateTime == 0 {
		f.CreateTime = time.Now().Unix()
	}
	return nil
}

// FavoriteDTO 收藏请求DTO
type FavoriteDTO struct {
	GoodsID int64 `json:"goods_id" binding:"required"`
}
