package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
}

func NewNotificationHandler(serviceCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{
		svcCtx: serviceCtx,
	}
}

// GetNotificationList 获取通知列表
func (h *NotificationHandler) GetNotificationList(c *gin.Context) {
	var req model.NotificationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	notifications, total, err := notificationLogic.GetNotificationList(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      notifications,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetUnreadCount 获取未读通知数量
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	count, err := notificationLogic.GetUnreadCount(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}

// MarkAsRead 标记单条通知为已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "通知ID无效")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	if err := notificationLogic.MarkAsRead(c.Request.Context(), notificationID, userID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "标记成功"})
}

// MarkAllAsRead 标记所有通知为已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	if err := notificationLogic.MarkAllAsRead(c.Request.Context(), userID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "全部标记为已读"})
}

// DeleteNotification 删除通知
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "通知ID无效")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	if err := notificationLogic.DeleteNotification(c.Request.Context(), notificationID, userID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ClearReadNotifications 清空已读通知
func (h *NotificationHandler) ClearReadNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	notificationLogic := logic.NewNotificationLogic(h.svcCtx.DB)

	if err := notificationLogic.ClearReadNotifications(c.Request.Context(), userID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "清空成功"})
}

// CreateNotification 创建通知（内部接口，用于其他模块调用）
func (h *NotificationHandler) CreateNotification(ctx int, userID int64, notifType int8, title, content, link string) error {
	// 这个接口一般不直接暴露给前端
	// 而是在后端其他逻辑中调用
	// 例如：交换请求被创建时，自动发送通知给卖家
	return nil
}
