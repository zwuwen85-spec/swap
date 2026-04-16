package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExchangeLogic 交换逻辑层
type ExchangeLogic struct {
	db *gorm.DB
}

// NewExchangeLogic 创建交换逻辑层
func NewExchangeLogic(db *gorm.DB) *ExchangeLogic {
	return &ExchangeLogic{db: db}
}

// Create 发起交换请求
func (l *ExchangeLogic) Create(ctx context.Context, initiatorID int64, req *model.ExchangeCreateDTO) (*model.Exchange, error) {
	// 1. 检查目标商品是否存在
	var goods model.Goods
	if err := l.db.WithContext(ctx).Where("id = ? AND status = 1", req.GoodsID).First(&goods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在或已下架")
		}
		return nil, err
	}

	// 2. 不能给自己的商品发起交换
	if goods.UserID == initiatorID {
		return nil, errors.New("不能给自己的商品发起交换")
	}

	// 3. 如果是交换，需要提供我的商品
	if req.Type == 2 && req.MyGoodsID == 0 {
		return nil, errors.New("交换类型必须提供要交换的商品")
	}

	// 4. 如果提供了我的商品，验证是否属于发起人
	if req.MyGoodsID > 0 {
		var myGoods model.Goods
		if err := l.db.WithContext(ctx).
			Where("id = ? AND user_id = ? AND status = 1", req.MyGoodsID, initiatorID).
			First(&myGoods).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("你的商品不存在或已下架")
			}
			return nil, err
		}
	}

	// 5. 检查是否已有待处理的交换请求
	var count int64
	if err := l.db.WithContext(ctx).Model(&model.Exchange{}).
		Where("initiator_id = ? AND goods_id = ? AND status IN (0, 1)",
			initiatorID, req.GoodsID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("已有待处理的交换请求")
	}

	// 6. 创建交换请求
	exchange := &model.Exchange{
		InitiatorID: initiatorID,
		TargetID:    goods.UserID,
		GoodsID:     req.GoodsID,
		MyGoodsID:   req.MyGoodsID,
		Type:        req.Type,
		Message:     req.Message,
		Status:      0, // 待处理
	}

	if err := l.db.WithContext(ctx).Create(exchange).Error; err != nil {
		return nil, err
	}

	// 7. 发送通知给卖家
	notificationLogic := NewNotificationLogic(l.db)
	_ = notificationLogic.CreateExchangeNotification(ctx, goods.UserID, exchange.ID, goods.Title, "create")

	return exchange, nil
}

// GetList 获取交换列表
func (l *ExchangeLogic) GetList(ctx context.Context, userID int64, req *model.ExchangeListRequest) ([]*model.Exchange, int64, error) {
	var exchanges []*model.Exchange
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Exchange{})

	// 筛选类型
	switch req.Type {
	case "incoming":
		// 收到的请求（我是目标用户）
		query = query.Where("target_id = ?", userID)
	case "outgoing":
		// 发起的请求（我是发起人）
		query = query.Where("initiator_id = ?", userID)
	case "all":
		// 全部（我是发起人或目标人）
		query = query.Where("initiator_id = ? OR target_id = ?", userID, userID)
	default:
		// 默认显示全部
		query = query.Where("initiator_id = ? OR target_id = ?", userID, userID)
	}

	// 状态筛选
	if req.Status >= 0 && req.Status <= 4 {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	// 查询列表（预加载关联数据）
	if err := query.
		Preload("Initiator").
		Preload("Target").
		Preload("Goods").
		Preload("MyGoods").
		Order("create_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&exchanges).Error; err != nil {
		return nil, 0, err
	}

	return exchanges, total, nil
}

// GetDetail 获取交换详情
func (l *ExchangeLogic) GetDetail(ctx context.Context, exchangeID int64) (*model.Exchange, error) {
	var exchange model.Exchange
	if err := l.db.WithContext(ctx).
		Preload("Initiator").
		Preload("Target").
		Preload("Goods").
		Preload("MyGoods").
		Where("id = ?", exchangeID).
		First(&exchange).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("交换请求不存在")
		}
		return nil, err
	}

	return &exchange, nil
}

// Handle 处理交换请求（接受/拒绝/取消/完成）
func (l *ExchangeLogic) Handle(ctx context.Context, userID int64, req *model.HandleExchangeDTO) error {
	// 1. 获取交换请求
	var exchange model.Exchange
	if err := l.db.WithContext(ctx).Where("id = ?", req.ExchangeID).First(&exchange).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("交换请求不存在")
		}
		return err
	}

	// 2. 根据操作类型处理
	switch req.Action {
	case "accept":
		return l.accept(ctx, userID, &exchange)
	case "reject":
		return l.reject(ctx, userID, &exchange, req.RejectReason)
	case "cancel":
		return l.cancel(ctx, userID, &exchange)
	case "complete":
		return l.complete(ctx, userID, &exchange)
	default:
		return errors.New("无效的操作")
	}
}

// accept 接受交换
func (l *ExchangeLogic) accept(ctx context.Context, userID int64, exchange *model.Exchange) error {
	// 检查权限：只有目标用户可以接受
	if exchange.TargetID != userID {
		return errors.New("无权操作此交换请求")
	}

	// 检查状态
	if !exchange.CanBeHandled() {
		return errors.New("该请求已被处理")
	}

	// 更新状态
	if err := l.db.WithContext(ctx).
		Model(exchange).
		Updates(map[string]interface{}{
			"status":     1, // 已接受
			"update_time": exchange.UpdateTime,
		}).Error; err != nil {
		return err
	}

	// 获取商品标题并发送通知
	var goods model.Goods
	l.db.WithContext(ctx).Select("title").Where("id = ?", exchange.GoodsID).First(&goods)
	notificationLogic := NewNotificationLogic(l.db)
	_ = notificationLogic.CreateExchangeNotification(ctx, exchange.InitiatorID, exchange.ID, goods.Title, "accept")

	return nil
}

// reject 拒绝交换
func (l *ExchangeLogic) reject(ctx context.Context, userID int64, exchange *model.Exchange, reason string) error {
	// 检查权限：只有目标用户可以拒绝
	if exchange.TargetID != userID {
		return errors.New("无权操作此交换请求")
	}

	// 检查状态
	if !exchange.CanBeHandled() {
		return errors.New("该请求已被处理")
	}

	// 更新状态
	if err := l.db.WithContext(ctx).
		Model(exchange).
		Updates(map[string]interface{}{
			"status":        2, // 已拒绝
			"reject_reason": reason,
			"update_time":   exchange.UpdateTime,
		}).Error; err != nil {
		return err
	}

	// 获取商品标题并发送通知
	var goods model.Goods
	l.db.WithContext(ctx).Select("title").Where("id = ?", exchange.GoodsID).First(&goods)
	notificationLogic := NewNotificationLogic(l.db)
	_ = notificationLogic.CreateExchangeNotification(ctx, exchange.InitiatorID, exchange.ID, goods.Title, "reject")

	return nil
}

// cancel 取消交换
func (l *ExchangeLogic) cancel(ctx context.Context, userID int64, exchange *model.Exchange) error {
	// 检查权限：只有发起人可以取消
	if exchange.InitiatorID != userID {
		return errors.New("无权操作此交换请求")
	}

	// 检查状态
	if !exchange.CanBeCancelled() {
		return errors.New("该请求无法取消")
	}

	// 更新状态
	if err := l.db.WithContext(ctx).
		Model(exchange).
		Updates(map[string]interface{}{
			"status":      3, // 已取消
			"update_time": exchange.UpdateTime,
		}).Error; err != nil {
		return err
	}

	return nil
}

// complete 完成交换
func (l *ExchangeLogic) complete(ctx context.Context, userID int64, exchange *model.Exchange) error {
	// 检查权限：发起人或目标用户都可以完成
	if exchange.InitiatorID != userID && exchange.TargetID != userID {
		return errors.New("无权操作此交换请求")
	}

	// 检查状态
	if exchange.Status != 1 {
		return errors.New("只有已接受的交换才能完成")
	}

	// 获取目标商品信息，根据商品的type决定最终状态
	var targetGoods model.Goods
	if err := l.db.WithContext(ctx).Select("type").Where("id = ?", exchange.GoodsID).First(&targetGoods).Error; err != nil {
		return err
	}

	// 根据商品类型确定商品状态
	// 商品type=1(售卖) → 已售出(status=2)
	// 商品type=2(交换) → 已交换(status=3)
	// 商品type=3(均可) → 根据交换类型：购买→已售出，交换→已交换
	var goodsStatus int8
	fmt.Printf("[DEBUG] 完成交换 - 商品ID: %d, 商品Type: %d, 交换Type: %d\n",
		exchange.GoodsID, targetGoods.Type, exchange.Type)

	if targetGoods.Type == 1 {
		goodsStatus = 2 // 售卖商品 → 已售出
		fmt.Printf("[DEBUG] 售卖商品完成 → 已售出(status=2)\n")
	} else if targetGoods.Type == 2 {
		goodsStatus = 3 // 交换商品 → 已交换
		fmt.Printf("[DEBUG] 交换商品完成 → 已交换(status=3)\n")
	} else {
		// type=3(均可)，根据交换请求类型
		goodsStatus = 2 // 默认已售出
		if exchange.Type == 2 {
			goodsStatus = 3 // 交换请求 → 已交换
		}
		fmt.Printf("[DEBUG] 均可商品完成 → status=%d\n", goodsStatus)
	}

	// 更新目标商品状态
	now := time.Now().Unix()
	// 使用 UpdateColumns 而不是 Updates，避免触发 BeforeUpdate 钩子
	// 这样就不会更新 update_time，只更新 sold_time 和 status
	if err := l.db.WithContext(ctx).
		Model(&model.Goods{}).
		Where("id = ?", exchange.GoodsID).
		UpdateColumns(map[string]interface{}{
			"status":    goodsStatus,
			"sold_time": now,
		}).Error; err != nil {
		return err
	}

	// 如果是物物交换，也更新我的商品状态为已交换
	if exchange.MyGoodsID > 0 {
		fmt.Printf("[DEBUG] 物物交换 - 我的商品ID: %d → 已交换(status=3)\n", exchange.MyGoodsID)
		if err := l.db.WithContext(ctx).
			Model(&model.Goods{}).
			Where("id = ?", exchange.MyGoodsID).
			UpdateColumns(map[string]interface{}{
				"status":    3, // 我的商品肯定是已交换
				"sold_time": now,
			}).Error; err != nil {
			return err
		}
	}

	// 更新交换状态
	if err := l.db.WithContext(ctx).
		Model(exchange).
		Updates(map[string]interface{}{
			"status":        4, // 已完成
			"complete_time": exchange.CompleteTime,
			"update_time":   exchange.UpdateTime,
		}).Error; err != nil {
		return err
	}

	// 获取商品标题并发送通知给双方
	var goods model.Goods
	l.db.WithContext(ctx).Select("title").Where("id = ?", exchange.GoodsID).First(&goods)
	notificationLogic := NewNotificationLogic(l.db)
	_ = notificationLogic.CreateExchangeNotification(ctx, exchange.InitiatorID, exchange.ID, goods.Title, "complete")

	return nil
}

// GetPendingCount 获取待处理交换数量
func (l *ExchangeLogic) GetPendingCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := l.db.WithContext(ctx).Model(&model.Exchange{}).
		Where("target_id = ? AND status = 0", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
