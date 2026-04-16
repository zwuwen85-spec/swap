package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password       string    `json:"-" gorm:"size:100;not null"` // 密码不在JSON中返回
	Email          *string   `json:"email" gorm:"uniqueIndex;size:100"`
	Phone          *string   `json:"phone" gorm:"uniqueIndex;size:20"`
	Avatar         *string   `json:"avatar" gorm:"size:255"`
	Nickname       *string   `json:"nickname" gorm:"size:50"`
	Gender         int8      `json:"gender" gorm:"default:0"` // 0未知 1男 2女
	StudentID      *string   `json:"student_id" gorm:"size:50"`
	School         *string   `json:"school" gorm:"size:100"`
	Major          *string   `json:"major" gorm:"size:100"`
	QQ             *string   `json:"qq" gorm:"size:20"`
	WeChat         *string   `json:"wechat" gorm:"column:we_chat;size:50"`
	Status         int8      `json:"status" gorm:"default:1"` // 0禁用 1正常 2冻结
	CreditScore    int       `json:"credit_score" gorm:"default:100"`
	CreateTime     int64     `json:"create_time" gorm:"not null"`
	UpdateTime     *int64    `json:"update_time"`
	LastLoginTime  *int64    `json:"last_login_time"`
	LastLoginIP    *string   `json:"last_login_ip" gorm:"size:50"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// BeforeCreate GORM钩子：创建前
func (u *User) BeforeCreate() error {
	if u.CreateTime == 0 {
		u.CreateTime = time.Now().Unix()
	}
	if u.Status == 0 {
		u.Status = 1
	}
	if u.CreditScore == 0 {
		u.CreditScore = 100
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (u *User) BeforeUpdate() error {
	now := time.Now().Unix()
	u.UpdateTime = &now
	return nil
}

// UserLoginDTO 用户登录DTO
type UserLoginDTO struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// UserRegisterDTO 用户注册DTO
type UserRegisterDTO struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	Nickname  string `json:"nickname" binding:"omitempty,max=50"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,len=11"`
	StudentID string `json:"student_id" binding:"omitempty,max=50"`
}

// UserUpdateDTO 用户更新DTO
type UserUpdateDTO struct {
	Nickname  string `json:"nickname" binding:"omitempty,max=50"`
	Avatar    string `json:"avatar" binding:"omitempty,max=255"`
	Gender    int8   `json:"gender" binding:"omitempty,min=0,max=2"`
	Phone     string `json:"phone" binding:"omitempty,max=11"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	StudentID string `json:"student_id" binding:"omitempty,max=50"`
	School    string `json:"school" binding:"omitempty,max=100"`
	Major     string `json:"major" binding:"omitempty,max=100"`
	QQ        string `json:"qq" binding:"omitempty,max=20"`
	WeChat    string `json:"wechat" binding:"omitempty,max=50"`
}

// ChangePasswordDTO 修改密码DTO
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}
