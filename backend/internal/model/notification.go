package model

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64  `json:"user_id" gorm:"not null;index:idx_user_id;index:idx_user_read"`
	Type       int8   `json:"type" gorm:"not null;index:idx_type"` // 1系统 2交换 3评论 4点赞
	Title      string `json:"title" gorm:"type:varchar(100);not null"`
	Content    string `json:"content" gorm:"type:text"`
	Link       string `json:"link" gorm:"type:varchar(255)"`
	IsRead     int8   `json:"is_read" gorm:"not null;default:0;index:idx_is_read"`
	CreateTime int64  `json:"create_time" gorm:"not null;index:idx_create_time"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notification"
}

// BeforeCreate 创建前钩子
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.CreateTime == 0 {
		n.CreateTime = time.Now().Unix()
	}
	return nil
}

// NotificationListRequest 通知列表请求
type NotificationListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=50"`
	Type     int8   `form:"type"`
	IsRead   *int8  `form:"is_read"` // 使用指针支持不传参（查询全部）
}

// NotificationCreateDTO 创建通知DTO
type NotificationCreateDTO struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Type    int8   `json:"type" binding:"required,min=1,max=4"`
	Title   string `json:"title" binding:"required,max=100"`
	Content string `json:"content" binding:"max=1000"`
	Link    string `json:"link" binding:"max=255"`
}

// 通知类型常量
const (
	NotificationTypeSystem    = 1 // 系统通知
	NotificationTypeExchange  = 2 // 交换通知
	NotificationTypeComment   = 3 // 评论通知
	NotificationTypeLike      = 4 // 点赞通知（暂未实现）
)
