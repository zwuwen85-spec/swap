package logic

import (
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"context"
	"errors"

	"gorm.io/gorm"
)

type FavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AddFavorite 添加收藏
func (l *FavoriteLogic) AddFavorite(userID int64, goodsID int64) error {
	// 检查商品是否存在
	var goods model.Goods
	err := l.svcCtx.DB.Where("id = ? AND status = 1", goodsID).First(&goods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或已下架")
		}
		return err
	}

	// 检查是否已收藏
	var count int64
	l.svcCtx.DB.Model(&model.Favorite{}).
		Where("user_id = ? AND goods_id = ?", userID, goodsID).
		Count(&count)
	if count > 0 {
		return errors.New("已经收藏过了")
	}

	// 创建收藏
	favorite := &model.Favorite{
		UserID:  userID,
		GoodsID: goodsID,
	}
	if err := l.svcCtx.DB.Create(favorite).Error; err != nil {
		return err
	}

	// 更新商品收藏数
	l.svcCtx.DB.Model(&model.Goods{}).
		Where("id = ?", goodsID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1"))

	return nil
}

// RemoveFavorite 取消收藏
func (l *FavoriteLogic) RemoveFavorite(userID int64, goodsID int64) error {
	// 删除收藏记录
	result := l.svcCtx.DB.Where("user_id = ? AND goods_id = ?", userID, goodsID).Delete(&model.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到收藏记录")
	}

	// 更新商品收藏数
	l.svcCtx.DB.Model(&model.Goods{}).
		Where("id = ?", goodsID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1"))

	return nil
}

// CheckFavorite 检查是否已收藏
func (l *FavoriteLogic) CheckFavorite(userID int64, goodsID int64) (bool, error) {
	var count int64
	err := l.svcCtx.DB.Model(&model.Favorite{}).
		Where("user_id = ? AND goods_id = ?", userID, goodsID).
		Count(&count).Error
	return count > 0, err
}

// GetFavoriteList 获取收藏列表
func (l *FavoriteLogic) GetFavoriteList(userID int64, page, pageSize int) ([]*model.Favorite, int64, error) {
	var favorites []*model.Favorite
	var total int64

	// 计算总数
	if err := l.svcCtx.DB.Model(&model.Favorite{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询收藏列表（包含商品信息）
	offset := (page - 1) * pageSize
	err := l.svcCtx.DB.Where("user_id = ?", userID).
		Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("Goods").
		Find(&favorites).Error

	if err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// GetFavoriteCount 获取收藏数量
func (l *FavoriteLogic) GetFavoriteCount(userID int64) (int64, error) {
	var count int64
	err := l.svcCtx.DB.Model(&model.Favorite{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
