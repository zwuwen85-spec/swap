package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

// stringValue 安全地获取字符串指针的值
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// MessageLogic 消息逻辑层
type MessageLogic struct {
	db *gorm.DB
}

// NewMessageLogic 创建消息逻辑层
func NewMessageLogic(db *gorm.DB) *MessageLogic {
	return &MessageLogic{db: db}
}

// SendMessage 发送消息
func (l *MessageLogic) SendMessage(ctx context.Context, senderID, receiverID int64, content string, msgType int8, goodsID int64) (*model.Message, error) {
	// 1. 验证接收者是否存在
	var user model.User
	if err := l.db.WithContext(ctx).Where("id = ? AND status = 1", receiverID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 2. 不能给自己发消息
	if senderID == receiverID {
		return nil, errors.New("不能给自己发送消息")
	}

	// 3. 验证商品ID（如果是商品卡片类型）
	if msgType == 3 && goodsID > 0 {
		var goods model.Goods
		if err := l.db.WithContext(ctx).Where("id = ?", goodsID).First(&goods).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("商品不存在")
			}
			return nil, err
		}
	}

	// 4. 创建消息
	message := &model.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		GoodsID:    goodsID,
		IsRead:     0,
	}

	if err := l.db.WithContext(ctx).Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessageList 获取消息列表
func (l *MessageLogic) GetMessageList(ctx context.Context, userID1, userID2 int64, page, pageSize int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	// 构建查询：两个用户之间的消息
	query := l.db.WithContext(ctx).Model(&model.Message{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	// 查询列表
	if err := query.
		Preload("Sender").
		Order("create_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetConversations 获取会话列表（最近聊天的用户）
func (l *MessageLogic) GetConversations(ctx context.Context, userID int64) ([]*model.Conversation, error) {
	type Result struct {
		UserID     int64
		UpdateTime int64
	}

	var results []*Result

	// 查询最近的对话
	query := `
		SELECT
			CASE
				WHEN sender_id = ? THEN receiver_id
				ELSE sender_id
			END AS user_id,
			MAX(create_time) as update_time
		FROM message
		WHERE sender_id = ? OR receiver_id = ?
		GROUP BY 1
		ORDER BY update_time DESC
	`

	if err := l.db.WithContext(ctx).Raw(query, userID, userID, userID).Scan(&results).Error; err != nil {
		return nil, err
	}

	// 查询用户信息和未读数
	conversations := make([]*model.Conversation, 0, len(results))
	for _, result := range results {
		// 查询用户信息
		var user model.User
		if err := l.db.WithContext(ctx).Where("id = ?", result.UserID).First(&user).Error; err != nil {
			continue
		}

		// 查询最后一条消息内容
		var lastMsg model.Message
		l.db.WithContext(ctx).
			Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, result.UserID, result.UserID, userID).
			Order("id DESC").
			First(&lastMsg)

		// 查询未读数
		var unreadCount int64
		l.db.WithContext(ctx).Model(&model.Message{}).
			Where("sender_id = ? AND receiver_id = ? AND is_read = 0", result.UserID, userID).
			Count(&unreadCount)

		// 截取最后一条消息的前50个字符
		lastMessageContent := lastMsg.Content
		runes := []rune(lastMessageContent)
		if len(runes) > 50 {
			lastMessageContent = string(runes[:50]) + "..."
		}

		conversations = append(conversations, &model.Conversation{
			UserID:      user.ID,
			Username:    user.Username,
			Nickname:    stringValue(user.Nickname),
			Avatar:      stringValue(user.Avatar),
			LastMessage: lastMessageContent,
			UnreadCount: unreadCount,
			UpdateTime:  result.UpdateTime,
		})
	}

	return conversations, nil
}

// MarkAsRead 标记消息为已读
func (l *MessageLogic) MarkAsRead(ctx context.Context, senderID, receiverID int64) error {
	return l.db.WithContext(ctx).Model(&model.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = 0", senderID, receiverID).
		Updates(map[string]interface{}{
			"is_read": 1,
		}).Error
}

// GetUnreadCount 获取未读消息总数
func (l *MessageLogic) GetUnreadCount(ctx context.Context, userID int64) (*model.UnreadCountResponse, error) {
	// 总未读数
	var total int64
	if err := l.db.WithContext(ctx).Model(&model.Message{}).
		Where("receiver_id = ? AND is_read = 0", userID).
		Count(&total).Error; err != nil {
		return nil, err
	}

	// 每个用户的未读数
	type Result struct {
		SenderID    int64
		UnreadCount int64
	}

	var results []*Result
	if err := l.db.WithContext(ctx).Model(&model.Message{}).
		Select("sender_id, COUNT(*) as unread_count").
		Where("receiver_id = ? AND is_read = 0", userID).
		Group("sender_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	userUnread := make(map[int64]int64)
	for _, result := range results {
		userUnread[result.SenderID] = result.UnreadCount
	}

	return &model.UnreadCountResponse{
		Total:      total,
		UserUnread: userUnread,
	}, nil
}
