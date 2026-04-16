package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/jwt"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	jwtManager  *jwt.JWTManager
	userLogic   *logic.UserLogic
	goodsLogic  *logic.GoodsLogic
}

// NewUserHandler 创建用户处理器
func NewUserHandler(serviceCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		jwtManager: jwt.NewJWTManager(serviceCtx.JwtSecret),
		userLogic:  logic.NewUserLogic(serviceCtx.DB),
		goodsLogic: logic.NewGoodsLogic(serviceCtx.DB),
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserRegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 注册
	user, err := h.userLogic.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeUserAlreadyExist, err.Error())
		return
	}

	// 生成Token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, response.CodeServerError, "生成Token失败")
		return
	}

	// 更新最后登录信息
	h.userLogic.UpdateLastLogin(c.Request.Context(), user.ID, c.ClientIP())

	response.SuccessWithMessage(c, "注册成功", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"token":    token,
		"expire":   7 * 24 * 3600, // 7天
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 登录
	user, err := h.userLogic.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound, err.Error())
		} else if err.Error() == "密码错误" {
			response.Error(c, response.CodePasswordError, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	// 生成Token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, response.CodeServerError, "生成Token失败")
		return
	}

	// 更新最后登录信息
	h.userLogic.UpdateLastLogin(c.Request.Context(), user.ID, c.ClientIP())

	response.SuccessWithMessage(c, "登录成功", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"token":    token,
		"expire":   7 * 24 * 3600, // 7天
	})
}

// GetUserInfo 获取当前用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	// 获取用户信息
	user, err := h.userLogic.GetUserInfoByID(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, response.CodeUserNotFound, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 更新用户信息
	if err := h.userLogic.UpdateUserInfo(c.Request.Context(), userID.(int64), &req); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 修改密码
	if err := h.userLogic.ChangePassword(c.Request.Context(), userID.(int64), &req); err != nil {
		if err.Error() == "原密码错误" {
			response.Error(c, response.CodePasswordError, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// UploadAvatar 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, response.CodeParamError, "请选择文件")
		return
	}

	// 验证文件大小（最大5MB）
	if file.Size > 5*1024*1024 {
		response.Error(c, response.CodeParamError, "文件大小不能超过5MB")
		return
	}

	// TODO: 实现文件上传到OSS或本地存储
	// 这里暂时返回一个模拟的URL
	avatarURL := "/uploads/avatar/" + strconv.FormatInt(userID.(int64), 10) + ".jpg"

	// 更新用户头像
	if err := h.userLogic.UpdateUserInfo(c.Request.Context(), userID.(int64), &model.UserUpdateDTO{
		Avatar: avatarURL,
	}); err != nil {
		response.Error(c, response.CodeServerError, "更新头像失败")
		return
	}

	response.SuccessWithMessage(c, "上传成功", gin.H{
		"url": avatarURL,
	})
}

// GetUserInfoByID 获取指定用户信息（公开接口）
func (h *UserHandler) GetUserInfoByID(c *gin.Context) {
	// 从URL参数获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError, "用户ID格式错误")
		return
	}

	// 获取用户信息
	user, err := h.userLogic.GetUserInfoByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeUserNotFound, err.Error())
		return
	}

	response.Success(c, user)
}

// GetUserGoods 获取指定用户的商品列表（公开接口）
func (h *UserHandler) GetUserGoods(c *gin.Context) {
	// 从URL参数获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError, "用户ID格式错误")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 获取用户的商品列表（只获取在售商品）
	list, total, err := h.goodsLogic.GetUserGoods(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.PageSuccess(c, list, total, page, pageSize)
}
