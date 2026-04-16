package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CommentLogic struct {
	db *gorm.DB
}

func NewCommentLogic(db *gorm.DB) *CommentLogic {
	return &CommentLogic{db: db}
}

// CreateComment 创建评论
func (l *CommentLogic) CreateComment(ctx context.Context, userID int64, req *model.CommentCreateDTO) (*model.Comment, error) {
	// 1. 验证商品是否存在
	var goods model.Goods
	err := l.db.WithContext(ctx).Where("id = ?", req.GoodsID).First(&goods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在")
		}
		return nil, err
	}

	// 2. 如果是回复评论，验证父评论是否存在
	targetUserID := goods.UserID // 默认被评论者是卖家
	if req.ParentID > 0 {
		var parentComment model.Comment
		err := l.db.WithContext(ctx).Where("id = ? AND goods_id = ?", req.ParentID, req.GoodsID).First(&parentComment).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			return nil, err
		}
		// 回复评论时，被评论者是父评论的发布者
		targetUserID = parentComment.UserID
	}

	// 3. 创建评论
	comment := &model.Comment{
		GoodsID:      req.GoodsID,
		UserID:       userID,
		TargetUserID: targetUserID,
		Content:      req.Content,
		Rating:       req.Rating,
		ParentID:     req.ParentID,
		Status:       1, // 默认显示
	}

	if err := l.db.WithContext(ctx).Create(comment).Error; err != nil {
		return nil, err
	}

	// 发送通知（只有主评论才通知卖家，回复不额外通知）
	if req.ParentID == 0 {
		// 获取评论者信息
		var commenter model.User
		l.db.WithContext(ctx).Select("username").Where("id = ?", userID).First(&commenter)

		notificationLogic := NewNotificationLogic(l.db)
		_ = notificationLogic.CreateCommentNotification(ctx, targetUserID, commenter.Username, goods.Title)
	}

	return comment, nil
}

// GetCommentList 获取商品评论列表
func (l *CommentLogic) GetCommentList(ctx context.Context, goodsID int64, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// 只查询主评论（parent_id = 0）
	query := l.db.WithContext(ctx).Model(&model.Comment{}).
		Where("goods_id = ? AND parent_id = 0 AND status = 1", goodsID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询评论列表（包含用户信息和回复）
	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = 1").Order("create_time ASC").Preload("User")
		}).
		Preload("Replies.User").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetCommentByID 根据ID获取评论
func (l *CommentLogic) GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error) {
	var comment model.Comment
	err := l.db.WithContext(ctx).
		Preload("User").
		Preload("Goods").
		Where("id = ?", commentID).
		First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func (l *CommentLogic) DeleteComment(ctx context.Context, commentID, userID int64) error {
	// 验证评论是否属于该用户
	var comment model.Comment
	err := l.db.WithContext(ctx).Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在或无权删除")
		}
		return err
	}

	// 软删除：将状态设为0
	if err := l.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("status", 0).Error; err != nil {
		return err
	}

	// 同时删除该评论的所有回复
	l.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("parent_id = ?", commentID).
		Update("status", 0)

	return nil
}

// GetGoodsRating 获取商品平均评分
func (l *CommentLogic) GetGoodsRating(ctx context.Context, goodsID int64) (float64, int64, error) {
	var result struct {
		AvgRating float64
		Count     int64
	}

	err := l.db.WithContext(ctx).
		Model(&model.Comment{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("goods_id = ? AND parent_id = 0 AND status = 1", goodsID).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.AvgRating, result.Count, nil
}

// GetUserComments 获取用户的评论历史
func (l *CommentLogic) GetUserComments(ctx context.Context, userID int64, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Comment{}).
		Where("user_id = ? AND parent_id = 0 AND status = 1", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("Goods").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetReceivedComments 获取用户收到的评论（作为卖家）
func (l *CommentLogic) GetReceivedComments(ctx context.Context, userID int64, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Comment{}).
		Where("target_user_id = ? AND parent_id = 0 AND status = 1", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("User").
		Preload("Goods").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
