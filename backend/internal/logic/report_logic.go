package logic

import (
	"campus-swap-shop/internal/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ReportLogic struct {
	db *gorm.DB
}

func NewReportLogic(db *gorm.DB) *ReportLogic {
	return &ReportLogic{db: db}
}

// CreateReport 创建举报
func (l *ReportLogic) CreateReport(ctx context.Context, reporterID int64, req *model.ReportCreateDTO) (*model.Report, error) {
	// 1. 验证目标对象是否存在
	switch req.TargetType {
	case model.ReportTypeGoods:
		var goods model.Goods
		if err := l.db.WithContext(ctx).Where("id = ?", req.TargetID).First(&goods).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("商品不存在")
			}
			return nil, err
		}
	case model.ReportTypeUser:
		var user model.User
		if err := l.db.WithContext(ctx).Where("id = ?", req.TargetID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("用户不存在")
			}
			return nil, err
		}
	case model.ReportTypeComment:
		var comment model.Comment
		if err := l.db.WithContext(ctx).Where("id = ?", req.TargetID).First(&comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("评论不存在")
			}
			return nil, err
		}
	default:
		return nil, errors.New("无效的举报类型")
	}

	// 2. 检查是否已经举报过（待处理或已处理的举报）
	var count int64
	l.db.WithContext(ctx).Model(&model.Report{}).
		Where("reporter_id = ? AND target_type = ? AND target_id = ? AND status IN (0, 1)",
			reporterID, req.TargetType, req.TargetID).
		Count(&count)
	if count > 0 {
		return nil, errors.New("您已经举报过该对象，请勿重复举报")
	}

	// 3. 创建举报
	report := &model.Report{
		ReporterID:  reporterID,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      model.ReportStatusPending,
	}

	if err := l.db.WithContext(ctx).Create(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetReportList 获取举报列表（管理员）
func (l *ReportLogic) GetReportList(ctx context.Context, req *model.ReportListRequest) ([]*model.Report, int64, error) {
	var reports []*model.Report
	var total int64

	// 构建查询
	query := l.db.WithContext(ctx).Model(&model.Report{})

	// 类型筛选
	if req.TargetType > 0 && req.TargetType <= 3 {
		query = query.Where("target_type = ?", req.TargetType)
	}

	// 状态筛选
	if req.Status >= 0 && req.Status <= 2 {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("create_time DESC").
		Limit(req.PageSize).
		Offset(offset).
		Preload("Reporter").
		Preload("Handler").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	// 根据target_type加载目标对象
	for _, report := range reports {
		switch report.TargetType {
		case model.ReportTypeGoods:
			l.db.WithContext(ctx).Preload("User").First(&report.Goods, report.TargetID)
		case model.ReportTypeUser:
			l.db.WithContext(ctx).First(&report.User, report.TargetID)
		case model.ReportTypeComment:
			l.db.WithContext(ctx).Preload("User").First(&report.Comment, report.TargetID)
		}
	}

	return reports, total, nil
}

// GetMyReports 获取我的举报列表
func (l *ReportLogic) GetMyReports(ctx context.Context, userID int64, page, pageSize int) ([]*model.Report, int64, error) {
	var reports []*model.Report
	var total int64

	query := l.db.WithContext(ctx).Model(&model.Report{}).Where("reporter_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("Handler").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	// 加载目标对象
	for _, report := range reports {
		switch report.TargetType {
		case model.ReportTypeGoods:
			l.db.WithContext(ctx).Select("id, title").First(&report.Goods, report.TargetID)
		case model.ReportTypeUser:
			l.db.WithContext(ctx).Select("id, username, nickname").First(&report.User, report.TargetID)
		case model.ReportTypeComment:
			l.db.WithContext(ctx).Select("id, content").First(&report.Comment, report.TargetID)
		}
	}

	return reports, total, nil
}

// HandleReport 处理举报（管理员）
func (l *ReportLogic) HandleReport(ctx context.Context, reportID, handlerID int64, req *model.ReportHandleDTO) error {
	// 1. 获取举报
	var report model.Report
	err := l.db.WithContext(ctx).Where("id = ?", reportID).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("举报不存在")
		}
		return err
	}

	// 2. 检查状态
	if report.Status != model.ReportStatusPending {
		return errors.New("该举报已被处理")
	}

	// 3. 更新举报状态
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status":        req.Status,
		"handler_id":    handlerID,
		"handle_result": req.HandleResult,
		"handle_time":   now,
	}

	if err := l.db.WithContext(ctx).
		Model(&report).
		Updates(updates).Error; err != nil {
		return err
	}

	// 4. 如果举报成立，自动处理目标对象
	if req.Status == model.ReportStatusProcessed {
		switch report.TargetType {
		case model.ReportTypeGoods:
			// 下架商品
			l.db.WithContext(ctx).
				Model(&model.Goods{}).
				Where("id = ?", report.TargetID).
				Update("status", 0) // 0表示下架
		case model.ReportTypeUser:
			// 封禁用户
			l.db.WithContext(ctx).
				Model(&model.User{}).
				Where("id = ?", report.TargetID).
				Update("status", 2) // 2表示冻结
		case model.ReportTypeComment:
			// 隐藏评论
			l.db.WithContext(ctx).
				Model(&model.Comment{}).
				Where("id = ?", report.TargetID).
				Update("status", 0) // 0表示隐藏
		}
	}

	return nil
}

// GetReportByID 获取举报详情
func (l *ReportLogic) GetReportByID(ctx context.Context, reportID int64) (*model.Report, error) {
	var report model.Report
	err := l.db.WithContext(ctx).
		Preload("Reporter").
		Preload("Handler").
		Where("id = ?", reportID).
		First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("举报不存在")
		}
		return nil, err
	}

	// 加载目标对象
	switch report.TargetType {
	case model.ReportTypeGoods:
		l.db.WithContext(ctx).Preload("User").First(&report.Goods, report.TargetID)
	case model.ReportTypeUser:
		l.db.WithContext(ctx).First(&report.User, report.TargetID)
	case model.ReportTypeComment:
		l.db.WithContext(ctx).Preload("User").First(&report.Comment, report.TargetID)
	}

	return &report, nil
}

// GetPendingCount 获取待处理举报数量（管理员）
func (l *ReportLogic) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := l.db.WithContext(ctx).
		Model(&model.Report{}).
		Where("status = ?", model.ReportStatusPending).
		Count(&count).Error
	return count, err
}

// CancelReport 撤销举报（仅限待处理状态）
func (l *ReportLogic) CancelReport(ctx context.Context, reportID, reporterID int64) error {
	result := l.db.WithContext(ctx).
		Where("id = ? AND reporter_id = ? AND status = ?", reportID, reporterID, model.ReportStatusPending).
		Delete(&model.Report{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("举报不存在或无法撤销")
	}

	return nil
}
