package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GoodsHandler 商品处理器
type GoodsHandler struct {
	goodsLogic *logic.GoodsLogic
	svcCtx     *svc.ServiceContext
}

// NewGoodsHandler 创建商品处理器
func NewGoodsHandler(serviceCtx *svc.ServiceContext) *GoodsHandler {
	return &GoodsHandler{
		goodsLogic: logic.NewGoodsLogic(serviceCtx.DB),
		svcCtx:     serviceCtx,
	}
}

// Create 发布商品
func (h *GoodsHandler) Create(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.GoodsCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 创建商品
	goods, err := h.goodsLogic.Create(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "发布成功", gin.H{
		"goods_id": goods.ID,
	})
}

// GetList 获取商品列表
func (h *GoodsHandler) GetList(c *gin.Context) {
	var req model.GoodsListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 获取列表
	list, total, err := h.goodsLogic.GetList(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
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

	response.PageSuccess(c, list, total, page, pageSize)
}

// GetDetail 获取商品详情
func (h *GoodsHandler) GetDetail(c *gin.Context) {
	goodsIDStr := c.Query("id")
	if goodsIDStr == "" {
		response.Error(c, response.CodeParamError, "缺少商品ID")
		return
	}

	goodsID, err := strconv.ParseInt(goodsIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError, "商品ID格式错误")
		return
	}

	// 获取详情
	goods, err := h.goodsLogic.GetDetail(c.Request.Context(), goodsID)
	if err != nil {
		if err.Error() == "商品不存在" {
			response.Error(c, response.CodeGoodsNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	// 增加浏览次数
	go h.goodsLogic.IncreaseViewCount(c.Request.Context(), goodsID)

	// 检查是否收藏（如果已登录）
	isFavorited := false
	if userID, exists := c.Get("user_id"); exists {
		favoriteLogic := logic.NewFavoriteLogic(c.Request.Context(), h.svcCtx)
		if favorited, err := favoriteLogic.CheckFavorite(userID.(int64), goodsID); err == nil {
			isFavorited = favorited
		}
	}

	response.Success(c, gin.H{
		"goods":        goods,
		"is_favorited": isFavorited,
	})
}

// GetMyGoods 获取我的商品
func (h *GoodsHandler) GetMyGoods(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	// 获取参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	statusStr := c.DefaultQuery("status", "-1")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	status, _ := strconv.Atoi(statusStr)

	// 获取列表
	list, total, err := h.goodsLogic.GetMyGoods(c.Request.Context(), userID.(int64), page, pageSize, status)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	if page < 1 {
		page = 1
	}

	response.PageSuccess(c, list, total, page, pageSize)
}

// Update 更新商品
func (h *GoodsHandler) Update(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.GoodsUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 更新商品
	if err := h.goodsLogic.Update(c.Request.Context(), userID.(int64), req.ID, &req); err != nil {
		if err.Error() == "商品不存在或无权操作" {
			response.Error(c, response.CodeNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除商品
func (h *GoodsHandler) Delete(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 删除商品
	if err := h.goodsLogic.Delete(c.Request.Context(), userID.(int64), req.ID); err != nil {
		if err.Error() == "商品不存在或无权操作" {
			response.Error(c, response.CodeNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetCategories 获取分类列表
func (h *GoodsHandler) GetCategories(c *gin.Context) {
	categories, err := h.goodsLogic.GetCategories(c.Request.Context())
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	// 直接返回扁平列表，前端根据需要自行构建树形结构
	response.Success(c, categories)
}

// Search 搜索商品
func (h *GoodsHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.Error(c, response.CodeParamError, "请输入搜索关键词")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 搜索
	list, total, err := h.goodsLogic.Search(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	if page < 1 {
		page = 1
	}

	response.PageSuccess(c, list, total, page, pageSize)
}

