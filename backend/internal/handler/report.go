package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	svcCtx *svc.ServiceContext
}

func NewReportHandler(serviceCtx *svc.ServiceContext) *ReportHandler {
	return &ReportHandler{
		svcCtx: serviceCtx,
	}
}

// CreateReport 创建举报
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req model.ReportCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	reporterID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)

	report, err := reportLogic.CreateReport(c.Request.Context(), reporterID.(int64), &req)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"report_id": report.ID,
		"message":   "举报成功，我们会尽快处理",
	})
}

// GetReportList 获取举报列表（管理员）
func (h *ReportHandler) GetReportList(c *gin.Context) {
	var req model.ReportListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)
	reports, total, err := reportLogic.GetReportList(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      reports,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetMyReports 获取我的举报列表
func (h *ReportHandler) GetMyReports(c *gin.Context) {
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

	reporterID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)

	reports, total, err := reportLogic.GetMyReports(c.Request.Context(), reporterID.(int64), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// HandleReport 处理举报（管理员）
func (h *ReportHandler) HandleReport(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseInt(reportIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "举报ID无效")
		return
	}

	var req model.ReportHandleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParam, err.Error())
		return
	}

	handlerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)

	if err := reportLogic.HandleReport(c.Request.Context(), reportID, handlerID.(int64), &req); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	statusText := "已驳回"
	if req.Status == 1 {
		statusText = "已处理"
	}

	response.Success(c, gin.H{"message": "举报" + statusText})
}

// GetReportDetail 获取举报详情（管理员）
func (h *ReportHandler) GetReportDetail(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseInt(reportIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "举报ID无效")
		return
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)
	report, err := reportLogic.GetReportByID(c.Request.Context(), reportID)
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, report)
}

// GetPendingCount 获取待处理举报数量（管理员）
func (h *ReportHandler) GetPendingCount(c *gin.Context) {
	reportLogic := logic.NewReportLogic(h.svcCtx.DB)
	count, err := reportLogic.GetPendingCount(c.Request.Context())
	if err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}

// CancelReport 撤销举报
func (h *ReportHandler) CancelReport(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseInt(reportIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrInvalidParam, "举报ID无效")
		return
	}

	reporterID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	reportLogic := logic.NewReportLogic(h.svcCtx.DB)

	if err := reportLogic.CancelReport(c.Request.Context(), reportID, reporterID.(int64)); err != nil {
		response.Error(c, response.ErrServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "撤销成功"})
}
