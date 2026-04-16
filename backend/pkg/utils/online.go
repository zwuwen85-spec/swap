package utils

import (
	"context"
	"fmt"
)

// OnlineStatusManager 在线状态管理器
type OnlineStatusManager struct {
	redis *RedisClient
}

// NewOnlineStatusManager 创建在线状态管理器
func NewOnlineStatusManager(redis *RedisClient) *OnlineStatusManager {
	return &OnlineStatusManager{
		redis: redis,
	}
}

// SetUserOnline 设置用户在线
func (m *OnlineStatusManager) SetUserOnline(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("online:user:%d", userID)
	return m.redis.SetEx(ctx, key, "1", 300) // 5分钟过期
}

// SetUserOffline 设置用户离线
func (m *OnlineStatusManager) SetUserOffline(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("online:user:%d", userID)
	return m.redis.Del(ctx, key)
}

// IsUserOnline 检查用户是否在线
func (m *OnlineStatusManager) IsUserOnline(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("online:user:%d", userID)
	exists, err := m.redis.Exists(ctx, key)
	return exists > 0, err
}

// GetOnlineUsers 获取所有在线用户
func (m *OnlineStatusManager) GetOnlineUsers(ctx context.Context) ([]int64, error) {
	keys, err := m.redis.Keys(ctx, "online:user:*").Result()
	if err != nil {
		return nil, err
	}

	users := make([]int64, 0, len(keys))
	for _, key := range keys {
		var userID int64
		_, err := fmt.Sscanf(key, "online:user:%d", &userID)
		if err == nil {
			users = append(users, userID)
		}
	}

	return users, nil
}

// RefreshOnlineStatus 刷新在线状态（心跳时调用）
func (m *OnlineStatusManager) RefreshOnlineStatus(ctx context.Context, userID int64) error {
	return m.SetUserOnline(ctx, userID)
}
