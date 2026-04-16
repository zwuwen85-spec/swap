package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCommentHandler(serviceCtx *svc.ServiceContext) *CommentHandler {
	return &CommentHandler{
		svcCtx: serviceCtx,
	}
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req model.CommentCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)

	comment, err := commentLogic.CreateComment(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"comment_id": comment.ID,
		"message":    "评论成功",
	})
}

// GetCommentList 获取商品评论列表
func (h *CommentHandler) GetCommentList(c *gin.Context) {
	var req model.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)
	comments, total, err := commentLogic.GetCommentList(c.Request.Context(), req.GoodsID, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "评论ID无效")
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)

	if err := commentLogic.DeleteComment(c.Request.Context(), commentID, userID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetGoodsRating 获取商品评分
func (h *CommentHandler) GetGoodsRating(c *gin.Context) {
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

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)
	avgRating, count, err := commentLogic.GetGoodsRating(c.Request.Context(), goodsID)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"avg_rating": avgRating,
		"count":      count,
	})
}

// GetUserComments 获取我的评论
func (h *CommentHandler) GetUserComments(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)

	comments, total, err := commentLogic.GetUserComments(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetReceivedComments 获取收到的评论
func (h *CommentHandler) GetReceivedComments(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	commentLogic := logic.NewCommentLogic(h.svcCtx.DB)

	comments, total, err := commentLogic.GetReceivedComments(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
