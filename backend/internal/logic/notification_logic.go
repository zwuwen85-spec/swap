package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type NotificationLogic struct {
	db *gorm.DB
}

func NewNotificationLogic(db *gorm.DB) *NotificationLogic {
	return &NotificationLogic{db: db}
}

// CreateNotification 创建通知
func (l *NotificationLogic) CreateNotification(ctx context.Context, req *model.NotificationCreateDTO) error {
	notification := &model.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Content: req.Content,
		Link:    req.Link,
		IsRead:  0,
	}

	return l.db.WithContext(ctx).Create(notification).Error
}

// CreateExchangeNotification 创建交换通知（辅助方法）
func (l *NotificationLogic) CreateExchangeNotification(ctx context.Context, userID int64, exchangeID int64, goodsTitle string, action string) error {
	var title string
	var content string

	switch action {
	case "create":
		title = "新的交换请求"
		content = "有人想要交换您的商品：" + goodsTitle
	case "accept":
		title = "交换请求已接受"
		content = "卖家接受了您的交换请求：" + goodsTitle
	case "reject":
		title = "交换请求已拒绝"
		content = "卖家拒绝了您的交换请求：" + goodsTitle
	case "complete":
		title = "交换已完成"
		content = "您的交换已完成：" + goodsTitle
	default:
		return nil
	}

	return l.CreateNotification(ctx, &model.NotificationCreateDTO{
		UserID:  userID,
		Type:    model.NotificationTypeExchange,
		Title:   title,
		Content: content,
		Link:    "/exchange/" + string(rune(exchangeID)),
	})
}

// CreateCommentNotification 创建评论通知（辅助方法）
func (l *NotificationLogic) CreateCommentNotification(ctx context.Context, targetUserID int64, commenterName string, goodsTitle string) error {
	return l.CreateNotification(ctx, &model.NotificationCreateDTO{
		UserID:  targetUserID,
		Type:    model.NotificationTypeComment,
		Title:   "新评论通知",
		Content: commenterName + " 评论了您的商品：" + goodsTitle,
	})
}

// GetNotificationList 获取通知列表
func (l *NotificationLogic) GetNotificationList(ctx context.Context, userID int64, req *model.NotificationListRequest) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var total int64

	// 构建查询
	query := l.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)

	// 类型筛选
	if req.Type > 0 && req.Type <= 4 {
		query = query.Where("type = ?", req.Type)
	}

	// 已读筛选
	if req.IsRead != nil {
		query = query.Where("is_read = ?", *req.IsRead)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("create_time DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetUnreadCount 获取未读通知数量
func (l *NotificationLogic) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := l.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = 0", userID).
		Count(&count).Error
	return count, err
}

// MarkAsRead 标记单条通知为已读
func (l *NotificationLogic) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	result := l.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", 1)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("通知不存在")
	}

	return nil
}

// MarkAllAsRead 标记所有通知为已读
func (l *NotificationLogic) MarkAllAsRead(ctx context.Context, userID int64) error {
	return l.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = 0", userID).
		Update("is_read", 1).Error
}

// DeleteNotification 删除通知
func (l *NotificationLogic) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	result := l.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&model.Notification{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("通知不存在")
	}

	return nil
}

// ClearReadNotifications 清空已读通知
func (l *NotificationLogic) ClearReadNotifications(ctx context.Context, userID int64) error {
	return l.db.WithContext(ctx).
		Where("user_id = ? AND is_read = 1", userID).
		Delete(&model.Notification{}).Error
}
