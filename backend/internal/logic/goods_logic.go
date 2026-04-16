package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// GoodsLogic 商品逻辑层
type GoodsLogic struct {
	db *gorm.DB
}

// NewGoodsLogic 创建商品逻辑层
func NewGoodsLogic(db *gorm.DB) *GoodsLogic {
	return &GoodsLogic{db: db}
}

// Create 创建商品
func (l *GoodsLogic) Create(ctx context.Context, userID int64, req *model.GoodsCreateDTO) (*model.Goods, error) {
	// 1. 验证分类是否存在
	var category model.Category
	if err := l.db.WithContext(ctx).Where("id = ? AND status = 1", req.CategoryID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	// 2. 验证价格（售卖类型必须要有价格）
	if req.Type == 1 && req.Price <= 0 {
		return nil, errors.New("售卖商品必须设置价格")
	}

	// 3. 创建商品
	goods := &model.Goods{
		Title:         req.Title,
		Description:   req.Description,
		CategoryID:    req.CategoryID,
		UserID:        userID,
		Type:          req.Type,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Condition:     req.Condition,
		Tags:          req.Tags,
		Location:      req.Location,
		Status:        1, // 默认在售
	}

	// 设置图片
	goods.SetImages(req.Images)

	// 手动设置创建时间（确保 BeforeCreate 钩子未执行时的回退）
	if goods.CreateTime == 0 {
		goods.CreateTime = time.Now().Unix()
	}

	if err := l.db.WithContext(ctx).Create(goods).Error; err != nil {
		return nil, err
	}

	return goods, nil
}

// GetList 获取商品列表
func (l *GoodsLogic) GetList(ctx context.Context, req *model.GoodsListRequest) ([]*model.Goods, int64, error) {
	var goods []*model.Goods
	var total int64

	// 构建查询
	query := l.db.WithContext(ctx).Model(&model.Goods{})

	// 筛选条件
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.Type > 0 && req.Type <= 3 {
		query = query.Where("type = ?", req.Type)
	}

	if req.Condition > 0 && req.Condition <= 4 {
		query = query.Where("`condition` = ?", req.Condition)
	}

	// 状态筛选：只显示在售商品（status=1），除非明确指定其他状态
	// 注意：Status的零值是0，但我们不将0视为"未设置"，而是视为"下架"
	// 因此我们需要检查前端是否真的传递了status参数
	// 为了简化，这里我们默认只显示在售商品，如果需要查看其他状态需要显式传递
	query = query.Where("status = 1")

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
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

	// 排序
	orderBy := "create_time DESC"
	switch req.Sort {
	case "time_asc":
		orderBy = "create_time ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "price_asc":
		orderBy = "price ASC"
	}

	// 查询列表（预加载用户和分类信息）
	if err := query.
		Preload("User").
		Preload("Category").
		Order(orderBy).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&goods).Error; err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// GetDetail 获取商品详情
func (l *GoodsLogic) GetDetail(ctx context.Context, goodsID int64) (*model.Goods, error) {
	var goods model.Goods
	if err := l.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Where("id = ?", goodsID).
		First(&goods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在")
		}
		return nil, err
	}

	// 增加浏览次数
	if err := l.db.WithContext(ctx).
		Model(&model.Goods{}).
		Where("id = ?", goodsID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		// 记录错误但不影响返回商品详情
		// 可以考虑添加日志
	}

	// 修复create_time为0的问题（旧数据兼容）
	if goods.CreateTime == 0 {
		// 优先使用update_time
		if goods.UpdateTime > 0 {
			goods.CreateTime = goods.UpdateTime
		} else if goods.User != nil && goods.User.CreateTime > 0 {
			// 如果商品没有时间戳，使用用户的创建时间作为参考
			// 因为商品必须在用户创建后才能创建
			goods.CreateTime = goods.User.CreateTime
		}
		// 如果都没有，保持create_time为0，前端会显示"暂无"
	}

	return &goods, nil
}

// GetMyGoods 获取我的商品
func (l *GoodsLogic) GetMyGoods(ctx context.Context, userID int64, page, pageSize, status int) ([]*model.Goods, int64, error) {
	var goods []*model.Goods
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Goods{}).Where("user_id = ?", userID)

	// 状态筛选
	if status >= 0 && status <= 3 {
		query = query.Where("status = ?", status)
	}

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

	// 查询列表（预加载用户和分类信息）
	if err := query.
		Preload("User").
		Preload("Category").
		Order("create_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&goods).Error; err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// Update 更新商品
func (l *GoodsLogic) Update(ctx context.Context, userID, goodsID int64, req *model.GoodsUpdateDTO) error {
	// 1. 检查商品是否存在且属于当前用户
	var goods model.Goods
	if err := l.db.WithContext(ctx).Where("id = ? AND user_id = ?", goodsID, userID).First(&goods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权操作")
		}
		return err
	}

	// 2. 构建更新数据
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if len(req.Images) > 0 {
		goods.SetImages(req.Images)
		updates["images"] = goods.Images
	}
	if req.Status >= 0 && req.Status <= 3 {
		updates["status"] = req.Status
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的信息")
	}

	// 手动设置更新时间
	updates["update_time"] = time.Now().Unix()

	// 3. 更新
	if err := l.db.WithContext(ctx).Model(&goods).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// Delete 删除商品
func (l *GoodsLogic) Delete(ctx context.Context, userID, goodsID int64) error {
	// 检查商品是否存在且属于当前用户
	var goods model.Goods
	if err := l.db.WithContext(ctx).Where("id = ? AND user_id = ?", goodsID, userID).First(&goods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权操作")
		}
		return err
	}

	// 删除商品（软删除）
	if err := l.db.WithContext(ctx).Delete(&goods).Error; err != nil {
		return err
	}

	return nil
}

// IncreaseViewCount 增加浏览次数
func (l *GoodsLogic) IncreaseViewCount(ctx context.Context, goodsID int64) error {
	return l.db.WithContext(ctx).
		Model(&model.Goods{}).
		Where("id = ?", goodsID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

// GetCategories 获取分类列表
func (l *GoodsLogic) GetCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category

	// 查询所有启用的分类
	if err := l.db.WithContext(ctx).
		Where("status = 1").
		Order("sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Search 搜索商品
func (l *GoodsLogic) Search(ctx context.Context, keyword string, page, pageSize int) ([]*model.Goods, int64, error) {
	var goods []*model.Goods
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Goods{}).
		Where("status = 1 AND title LIKE ?", "%"+keyword+"%")

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

	// 查询列表（预加载分类信息）
	if err := query.
		Preload("Category").
		Order("create_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&goods).Error; err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// GetUserGoods 获取指定用户的商品列表
func (l *GoodsLogic) GetUserGoods(ctx context.Context, userID int64, page, pageSize int) ([]*model.Goods, int64, error) {
	var goods []*model.Goods
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Goods{}).Where("user_id = ? AND status = 1", userID)

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

	// 查询列表（预加载用户和分类信息）
	if err := query.
		Preload("User").
		Preload("Category").
		Order("create_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&goods).Error; err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}
