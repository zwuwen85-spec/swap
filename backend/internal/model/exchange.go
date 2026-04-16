package model

import (
	"time"
)

// Exchange 交换模型
type Exchange struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	InitiatorID  int64  `json:"initiator_id" gorm:"not null"`
	TargetID     int64  `json:"target_id" gorm:"not null"`
	GoodsID      int64  `json:"goods_id" gorm:"not null"`
	MyGoodsID    int64  `json:"my_goods_id"`           // 用于物物交换
	Type         int8   `json:"type" gorm:"default:1"`  // 1购买 2交换
	Message      string `json:"message" gorm:"type:text"`
	Status       int8   `json:"status" gorm:"default:0"` // 0待处理 1已接受 2已拒绝 3已取消 4已完成
	RejectReason string `json:"reject_reason" gorm:"size:200"`
	CreateTime   int64  `json:"create_time" gorm:"not null"`
	UpdateTime   int64  `json:"update_time"`
	CompleteTime int64  `json:"complete_time"`

	// 关联
	Initiator *User  `json:"initiator,omitempty" gorm:"foreignKey:InitiatorID"`
	Target    *User  `json:"target,omitempty" gorm:"foreignKey:TargetID"`
	Goods     *Goods `json:"goods,omitempty" gorm:"foreignKey:GoodsID"`
	MyGoods   *Goods `json:"my_goods,omitempty" gorm:"foreignKey:MyGoodsID"`
}

// TableName 指定表名
func (Exchange) TableName() string {
	return "exchange"
}

// BeforeCreate GORM钩子：创建前
func (e *Exchange) BeforeCreate() error {
	if e.CreateTime == 0 {
		e.CreateTime = time.Now().Unix()
	}
	if e.Status == 0 {
		e.Status = 0 // 默认待处理
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (e *Exchange) BeforeUpdate() error {
	e.UpdateTime = time.Now().Unix()

	// 如果状态变为已完成，记录完成时间
	if e.Status == 4 && e.CompleteTime == 0 {
		e.CompleteTime = time.Now().Unix()
	}

	return nil
}

// ExchangeCreateDTO 创建交换请求DTO
type ExchangeCreateDTO struct {
	GoodsID   int64  `json:"goods_id" binding:"required"`
	Type      int8   `json:"type" binding:"required,min=1,max=2"`
	MyGoodsID int64  `json:"my_goods_id"` // 交换时必填
	Message   string `json:"message"`
}

// ExchangeListRequest 交换列表请求
type ExchangeListRequest struct {
	Type   string `form:"type" binding:"omitempty,oneof=incoming outgoing all"` // incoming/outgoing/all
	Status int8   `form:"status" binding:"omitempty,min=0,max=4"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	PageSize int  `form:"page_size" binding:"omitempty,min=1,max=50"`
}

// HandleExchangeDTO 处理交换DTO
type HandleExchangeDTO struct {
	ExchangeID   int64  `json:"exchange_id" binding:"required"`
	Action       string `json:"action" binding:"required,oneof=accept reject cancel complete"`
	RejectReason string `json:"reject_reason"`
}

// GetStatusText 获取状态文本
func (e *Exchange) GetStatusText() string {
	statusMap := map[int8]string{
		0: "待处理",
		1: "已接受",
		2: "已拒绝",
		3: "已取消",
		4: "已完成",
	}
	return statusMap[e.Status]
}

// GetTypeText 获取类型文本
func (e *Exchange) GetTypeText() string {
	typeMap := map[int8]string{
		1: "购买",
		2: "交换",
	}
	return typeMap[e.Type]
}

// CanBeHandled 检查是否可以被处理
func (e *Exchange) CanBeHandled() bool {
	return e.Status == 0 // 只有待处理的可以处理
}

// CanBeCancelled 检查是否可以被取消
func (e *Exchange) CanBeCancelled() bool {
	return e.Status == 0 // 只有待处理的可以取消
}

// IsCompleted 检查是否已完成
func (e *Exchange) IsCompleted() bool {
	return e.Status == 4
}

// IsRejected 检查是否已拒绝
func (e *Exchange) IsRejected() bool {
	return e.Status == 2
}

// IsCancelled 检查是否已取消
func (e *Exchange) IsCancelled() bool {
	return e.Status == 3
}
