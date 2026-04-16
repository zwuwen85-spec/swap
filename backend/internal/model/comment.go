package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID            int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	GoodsID       int64  `json:"goods_id" gorm:"not null;index:idx_goods_id"`
	UserID        int64  `json:"user_id" gorm:"not null;index:idx_user_id"`
	TargetUserID  int64  `json:"target_user_id" gorm:"index"`
	Content       string `json:"content" gorm:"type:text;not null"`
	Rating        int8   `json:"rating" gorm:"not null;default:5"`
	ParentID      int64  `json:"parent_id" gorm:"default:0"`
	Status        int8   `json:"status" gorm:"not null;default:1"`
	CreateTime    int64  `json:"create_time" gorm:"not null;index:idx_create_time"`

	// 关联
	User       *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	TargetUser *User    `json:"target_user,omitempty" gorm:"foreignKey:TargetUserID"`
	Goods      *Goods   `json:"goods,omitempty" gorm:"foreignKey:GoodsID"`
	Parent     *Comment `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies    []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comment"
}

// BeforeCreate 创建前钩子
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.CreateTime == 0 {
		c.CreateTime = time.Now().Unix()
	}
	return nil
}

// CommentCreateDTO 创建评论DTO
type CommentCreateDTO struct {
	GoodsID      int64  `json:"goods_id" binding:"required"`
	Content      string `json:"content" binding:"required,max=500"`
	Rating       int8   `json:"rating" binding:"required,min=1,max=5"`
	ParentID     int64  `json:"parent_id"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	GoodsID int64 `form:"goods_id" binding:"required"`
	Page    int   `form:"page" binding:"min=1"`
	PageSize int   `form:"page_size" binding:"min=1,max=50"`
}

// CommentReplyDTO 回复评论DTO
type CommentReplyDTO struct {
	Content   string `json:"content" binding:"required,max=500"`
	ParentID  int64  `json:"parent_id" binding:"required"`
}
