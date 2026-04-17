package model

import (
	"time"
)

// Message 聊天消息模型
type Message struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	SenderID   int64  `json:"sender_id" gorm:"not null"`
	ReceiverID int64  `json:"receiver_id" gorm:"not null"`
	Content    string `json:"content" gorm:"type:text;not null"`
	Type       int8   `json:"type" gorm:"default:1"`           // 1文本 2图片 3商品卡片
	GoodsID    int64  `json:"goods_id,omitempty"`
	IsRead     int8   `json:"is_read" gorm:"default:0"`       // 0未读 1已读
	CreateTime int64  `json:"create_time" gorm:"not null"`

	// 关联
	Sender *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "message"
}

// BeforeCreate GORM钩子：创建前
func (m *Message) BeforeCreate() error {
	if m.CreateTime == 0 {
		m.CreateTime = time.Now().Unix()
	}
	return nil
}

// MessageListRequest 消息列表请求
type MessageListRequest struct {
	UserID  int64 `form:"user_id" binding:"required"` // 聊天对象ID
	Page    int   `form:"page" binding:"omitempty,min=1"`
	PageSize int  `form:"page_size" binding:"omitempty,min=1,max=50"`
}

// UnreadCountResponse 未读消息数量响应
type UnreadCountResponse struct {
	Total      int64            `json:"total"`
	UserUnread map[int64]int64 `json:"user_unread"` // 每个用户的未读数
}

// Conversation 会话
type Conversation struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	LastMessage string `json:"last_message"`
	UnreadCount int64  `json:"unread_count"`
	UpdateTime   int64  `json:"update_time"`
}

// SendMessageDTO 发送消息DTO
type SendMessageDTO struct {
	ReceiverID int64 `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=500"`
	Type       int8   `json:"type" binding:"omitempty,min=1,max=3"`
	GoodsID    int64  `json:"goods_id"`
}
