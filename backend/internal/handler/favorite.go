package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	svcCtx *svc.ServiceContext
}

func NewFavoriteHandler(svcCtx *svc.ServiceContext) *FavoriteHandler {
	return &FavoriteHandler{
		svcCtx: svcCtx,
	}
}

// AddFavorite 添加收藏
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var req struct {
		GoodsID int64 `json:"goods_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	favoriteLogic := logic.NewFavoriteLogic(c, h.svcCtx)

	if err := favoriteLogic.AddFavorite(userID.(int64), req.GoodsID); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "收藏成功"})
}

// RemoveFavorite 取消收藏
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	goodsIDStr := c.Param("goodsId")
	goodsID, err := strconv.ParseInt(goodsIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "商品ID无效")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	favoriteLogic := logic.NewFavoriteLogic(c, h.svcCtx)

	if err := favoriteLogic.RemoveFavorite(userID.(int64), goodsID); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "取消收藏成功"})
}

// CheckFavorite 检查是否已收藏
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	goodsIDStr := c.Query("goods_id")
	if goodsIDStr == "" {
		response.Error(c, response.ErrInvalidParam, "商品ID不能为空")
		return
	}

	goodsID, err := strconv.ParseInt(goodsIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "商品ID无效")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	favoriteLogic := logic.NewFavoriteLogic(c, h.svcCtx)

	isFavorited, err := favoriteLogic.CheckFavorite(userID.(int64), goodsID)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"is_favorited": isFavorited})
}

// GetFavoriteList 获取收藏列表
func (h *FavoriteHandler) GetFavoriteList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	favoriteLogic := logic.NewFavoriteLogic(c, h.svcCtx)

	list, total, err := favoriteLogic.GetFavoriteList(userID.(int64), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetFavoriteCount 获取收藏数量
func (h *FavoriteHandler) GetFavoriteCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	favoriteLogic := logic.NewFavoriteLogic(c, h.svcCtx)

	count, err := favoriteLogic.GetFavoriteCount(userID.(int64))
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}
