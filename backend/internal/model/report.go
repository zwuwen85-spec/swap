package model

import (
	"time"

	"gorm.io/gorm"
)

// Report 举报模型
type Report struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ReporterID   int64  `json:"reporter_id" gorm:"not null;index:idx_reporter_id"`
	TargetType   int8   `json:"target_type" gorm:"not null;index:idx_target"` // 1商品 2用户 3评论
	TargetID     int64  `json:"target_id" gorm:"not null;index:idx_target"`
	Reason       string `json:"reason" gorm:"type:varchar(200);not null"`
	Description  string `json:"description" gorm:"type:text"`
	Status       int8   `json:"status" gorm:"not null;default:0;index:idx_status"` // 0待处理 1已处理 2已驳回
	HandlerID    int64  `json:"handler_id" gorm:"index"`
	HandleResult string `json:"handle_result" gorm:"type:text"`
	HandleTime   int64  `json:"handle_time"`
	CreateTime   int64  `json:"create_time" gorm:"not null;index:idx_create_time"`

	// 关联
	Reporter *User `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Handler  *User `json:"handler,omitempty" gorm:"foreignKey:HandlerID"`

	// 目标对象（根据target_type动态关联）
	Goods   *Goods   `json:"goods,omitempty" gorm:"foreignKey:TargetID;references:ID"`
	User    *User    `json:"target_user,omitempty" gorm:"foreignKey:TargetID;references:ID"`
	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:TargetID;references:ID"`
}

// TableName 指定表名
func (Report) TableName() string {
	return "report"
}

// BeforeCreate 创建前钩子
func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.CreateTime == 0 {
		r.CreateTime = time.Now().Unix()
	}
	return nil
}

// ReportCreateDTO 创建举报DTO
type ReportCreateDTO struct {
	TargetType  int8   `json:"target_type" binding:"required,min=1,max=3"`
	TargetID    int64  `json:"target_id" binding:"required"`
	Reason      string `json:"reason" binding:"required,max=200"`
	Description string `json:"description" binding:"max=1000"`
}

// ReportListRequest 举报列表请求
type ReportListRequest struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=50"`
	TargetType int8   `form:"target_type"`
	Status     int8   `form:"status"`
}

// ReportHandleDTO 处理举报DTO
type ReportHandleDTO struct {
	Status       int8   `json:"status" binding:"required,min=1,max=2"` // 1已处理 2已驳回
	HandleResult string `json:"handle_result" binding:"max=1000"`
}

// 举报类型常量
const (
	ReportTypeGoods   = 1 // 举报商品
	ReportTypeUser    = 2 // 举报用户
	ReportTypeComment = 3 // 举报评论
)

// 举报状态常量
const (
	ReportStatusPending   = 0 // 待处理
	ReportStatusProcessed = 1 // 已处理
	ReportStatusRejected  = 2 // 已驳回
)

// 举报原因
var (
	ReportReasons = map[int8][]string{
		1: {"虚假信息", "违规商品", "价格欺诈", "图片不符", "其他"},
		2: {"恶意行为", "欺诈行为", "骚扰用户", "虚假身份", "其他"},
		3: {"违法违规", "恶意攻击", "广告 spam", "不实信息", "其他"},
	}
)
