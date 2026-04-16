package logic

import (
	"campus-swap-shop/internal/model"
	"campus-swap-shop/pkg/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

// stringPtr 返回字符串的指针，如果字符串为空则返回 nil
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// UserLogic 用户逻辑层
type UserLogic struct {
	db *gorm.DB
}

// NewUserLogic 创建用户逻辑层
func NewUserLogic(db *gorm.DB) *UserLogic {
	return &UserLogic{db: db}
}

// Register 用户注册
func (l *UserLogic) Register(ctx context.Context, req *model.UserRegisterDTO) (*model.User, error) {
	// 1. 检查用户名是否已存在
	var count int64
	if err := l.db.WithContext(ctx).Table("user").Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	if req.Email != "" {
		if err := l.db.WithContext(ctx).Table("user").Where("email = ?", req.Email).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 3. 检查手机号是否已存在
	if req.Phone != "" {
		if err := l.db.WithContext(ctx).Table("user").Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 4. 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. 创建用户
	user := &model.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Status:    1,
		Nickname:  stringPtr(req.Nickname),
		Email:     stringPtr(req.Email),
		Phone:     stringPtr(req.Phone),
		StudentID: stringPtr(req.StudentID),
	}

	if err := l.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	// 6. 返回用户信息（不包含密码）
	user.Password = ""
	return user, nil
}

// Login 用户登录
func (l *UserLogic) Login(ctx context.Context, req *model.UserLoginDTO) (*model.User, error) {
	// 1. 查询用户
	var user model.User
	if err := l.db.WithContext(ctx).Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 2. 检查账号状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用或冻结")
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 4. 返回用户信息（不包含密码）
	user.Password = ""
	return &user, nil
}

// GetUserInfoByID 根据ID获取用户信息
func (l *UserLogic) GetUserInfoByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	if err := l.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

// UpdateUserInfo 更新用户信息
func (l *UserLogic) UpdateUserInfo(ctx context.Context, userID int64, req *model.UserUpdateDTO) error {
	// 构建更新数据
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender >= 0 && req.Gender <= 2 {
		updates["gender"] = req.Gender
	}
	if req.Phone != "" {
		// 检查手机号是否被其他用户使用
		var count int64
		if err := l.db.WithContext(ctx).Table("user").
			Where("phone = ? AND id != ?", req.Phone, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已被使用")
		}
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		// 检查邮箱是否被其他用户使用
		var count int64
		if err := l.db.WithContext(ctx).Table("user").
			Where("email = ? AND id != ?", req.Email, userID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已被使用")
		}
		updates["email"] = req.Email
	}
	if req.StudentID != "" {
		updates["student_id"] = req.StudentID
	}
	if req.School != "" {
		updates["school"] = req.School
	}
	if req.Major != "" {
		updates["major"] = req.Major
	}
	if req.QQ != "" {
		updates["qq"] = req.QQ
	}
	if req.WeChat != "" {
		updates["we_chat"] = req.WeChat
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的信息")
	}

	// 更新
	if err := l.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// ChangePassword 修改密码
func (l *UserLogic) ChangePassword(ctx context.Context, userID int64, req *model.ChangePasswordDTO) error {
	// 1. 获取用户当前密码
	var user model.User
	if err := l.db.WithContext(ctx).Select("password").Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// 2. 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 3. 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 4. 更新密码
	if err := l.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}

// UpdateLastLogin 更新最后登录信息
func (l *UserLogic) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	return l.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_time": utils.GetCurrentTimestamp(),
			"last_login_ip":   ip,
		}).Error
}
