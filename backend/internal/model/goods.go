package model

import (
	"encoding/json"
	"time"
)

// Category 分类模型
type Category struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"size:50;not null"`
	ParentID   int64  `json:"parent_id" gorm:"default:0"`
	Level      int8   `json:"level" gorm:"default:1"`
	Icon       string `json:"icon" gorm:"size:255"`
	Sort       int    `json:"sort" gorm:"default:0"`
	Status     int8   `json:"status" gorm:"default:1"` // 0禁用 1启用
	CreateTime int64  `json:"create_time" gorm:"not null"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "category"
}

// Goods 商品模型
type Goods struct {
	ID             int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string   `json:"title" gorm:"size:100;not null"`
	Description    string   `json:"description" gorm:"type:text"`
	CategoryID     int64     `json:"category_id" gorm:"not null"`
	Category       *Category `json:"category" gorm:"foreignKey:CategoryID"`
	UserID         int64     `json:"user_id" gorm:"not null"`
	User           *User     `json:"user" gorm:"foreignKey:UserID"`
	Type           int8     `json:"type" gorm:"default:1"`              // 1售卖 2交换 3均可
	Price          float64  `json:"price" gorm:"type:decimal(10,2)"`
	OriginalPrice  float64  `json:"original_price" gorm:"type:decimal(10,2)"`
	Images         string   `json:"images" gorm:"type:json"`             // JSON数组
	Condition      int8     `json:"condition" gorm:"default:1"`          // 1全新 2九成新 3八成新 4七成新
	Status         int8     `json:"status" gorm:"default:0"`             // 0下架 1在售 2已售 3已交换
	ViewCount      int      `json:"view_count" gorm:"default:0"`
	FavoriteCount  int      `json:"favorite_count" gorm:"default:0"`
	ExchangeCount  int      `json:"exchange_count" gorm:"default:0"`
	Tags           string   `json:"tags" gorm:"size:200"`
	Location       string   `json:"location" gorm:"size:100"`
	Latitude       float64  `json:"latitude" gorm:"type:decimal(10,7)"`
	Longitude      float64  `json:"longitude" gorm:"type:decimal(10,7)"`
	CreateTime     int64    `json:"create_time" gorm:"not null"`
	UpdateTime     int64    `json:"update_time"`
	SoldTime       int64    `json:"sold_time"`
}

// TableName 指定表名
func (Goods) TableName() string {
	return "goods"
}

// GetImages 获取图片数组
func (g *Goods) GetImages() []string {
	if g.Images == "" {
		return []string{}
	}
	var images []string
	json.Unmarshal([]byte(g.Images), &images)
	return images
}

// SetImages 设置图片数组
func (g *Goods) SetImages(images []string) {
	data, _ := json.Marshal(images)
	g.Images = string(data)
}

// BeforeCreate GORM钩子：创建前
func (g *Goods) BeforeCreate() error {
	if g.CreateTime == 0 {
		g.CreateTime = time.Now().Unix()
	}
	if g.Status == 0 {
		g.Status = 1 // 默认在售
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (g *Goods) BeforeUpdate() error {
	g.UpdateTime = time.Now().Unix()
	return nil
}

// GoodsCreateDTO 创建商品DTO
type GoodsCreateDTO struct {
	Title         string   `json:"title" binding:"required,min=1,max=100"`
	Description   string   `json:"description"`
	CategoryID    int64    `json:"category_id" binding:"required"`
	Type          int8     `json:"type" binding:"required,min=1,max=3"`
	Price         float64  `json:"price" binding:"required_if=Type 1"`
	OriginalPrice float64  `json:"original_price"`
	Images        []string `json:"images" binding:"required,min=1,max=9"`
	Condition     int8     `json:"condition" binding:"required,min=1,max=4"`
	Tags          string   `json:"tags"`
	Location      string   `json:"location"`
}

// GoodsListRequest 商品列表请求
type GoodsListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=50"`
	CategoryID int64  `form:"category_id"`
	Type       int8   `form:"type" binding:"omitempty,min=0,max=3"`
	Condition  int8   `form:"condition" binding:"omitempty,min=0,max=4"`
	Status     int8   `form:"status" binding:"omitempty,min=0,max=3"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort" binding:"omitempty,oneof=time_desc time_asc price_desc price_asc"`
}

// GoodsUpdateDTO 更新商品DTO
type GoodsUpdateDTO struct {
	ID         int64    `json:"id" binding:"required"`
	Title      string   `json:"title" binding:"omitempty,min=1,max=100"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	Status      int8     `json:"status" binding:"omitempty,min=0,max=3"`
}

// GoodsDetailResponse 商品详情响应
type GoodsDetailResponse struct {
	*Goods
	IsFavorited bool `json:"is_favorited"` // 当前用户是否已收藏
}
