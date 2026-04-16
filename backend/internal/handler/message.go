package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"campus-swap-shop/pkg/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageLogic *logic.MessageLogic
	wsManager    *websocket.WebSocketManager
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(serviceCtx *svc.ServiceContext, wsManager *websocket.WebSocketManager) *MessageHandler {
	return &MessageHandler{
		messageLogic: logic.NewMessageLogic(serviceCtx.DB),
		wsManager:    wsManager,
	}
}

// SendMessage 发送消息（HTTP接口，用于离线消息）
func (h *MessageHandler) SendMessage(c *gin.Context) {
	// 从上下文获取用户ID
	senderID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.SendMessageDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 发送消息
	message, err := h.messageLogic.SendMessage(c.Request.Context(), senderID.(int64), req.ReceiverID, req.Content, req.Type, req.GoodsID)
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	// 如果接收者在线，通过WebSocket推送
	if h.wsManager.IsOnline(req.ReceiverID) {
		wsMessage := &websocket.Message{
			Type:      "message",
			SenderID:  message.SenderID,
			ReceiverID: message.ReceiverID,
			Content:   message.Content,
			Timestamp: message.CreateTime,
		}

		sent := h.wsManager.SendToUser(req.ReceiverID, wsMessage)
		if sent {
			// 实时推送成功，标记为已读
			h.messageLogic.MarkAsRead(c.Request.Context(), message.SenderID, message.ReceiverID)
		}
	}

	response.SuccessWithMessage(c, "发送成功", gin.H{
		"message_id": message.ID,
	})
}

// GetMessageList 获取消息列表
func (h *MessageHandler) GetMessageList(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	targetIDStr := c.Query("user_id")
	if targetIDStr == "" {
		response.Error(c, response.CodeParamError, "缺少聊天对象ID")
		return
	}

	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError, "用户ID格式错误")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 获取消息列表
	list, total, err := h.messageLogic.GetMessageList(c.Request.Context(), userID.(int64), targetID, page, pageSize)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	// 标记为已读
	if total > 0 {
		// 判断哪个是发送者
		var senderID int64
		if userID.(int64) > targetID {
			senderID = targetID
		} else {
			senderID = userID.(int64)
		}

		h.messageLogic.MarkAsRead(c.Request.Context(), senderID, userID.(int64))
	}

	if page < 1 {
		page = 1
	}

	response.PageSuccess(c, list, total, page, pageSize)
}

// GetConversations 获取会话列表
func (h *MessageHandler) GetConversations(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	conversations, err := h.messageLogic.GetConversations(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, conversations)
}

// GetUnreadCount 获取未读消息数量
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	count, err := h.messageLogic.GetUnreadCount(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, count)
}
