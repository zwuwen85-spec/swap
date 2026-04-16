package handler

import (
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploadPath string
	maxSize    int64
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(serviceCtx *svc.ServiceContext) *UploadHandler {
	maxSize := serviceCtx.Config.GetInt64("upload.max_size")
	// 默认最大5MB
	if maxSize == 0 {
		maxSize = 5 * 1024 * 1024
	}
	return &UploadHandler{
		uploadPath: serviceCtx.Config.GetString("upload.path"),
		maxSize:    maxSize,
	}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取用户ID（可选）
	userID, _ := c.Get("user_id")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, response.CodeParamError, "请选择文件")
		return
	}

	// 验证文件大小
	if file.Size > h.maxSize {
		response.Error(c, response.CodeParamError, fmt.Sprintf("文件大小不能超过%dMB", h.maxSize/(1024*1024)))
		return
	}

	// 验证文件类型
	if !h.isAllowedImageType(file) {
		response.Error(c, response.CodeParamError, "只支持jpg、jpeg、png、gif格式的图片")
		return
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := h.generateFileName(userID, ext)

	// 创建目录
	date := time.Now().Format("2006-01-02")
	dir := filepath.Join(h.uploadPath, date)
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, response.CodeServerError, "创建目录失败")
		return
	}

	// 保存文件
	filePath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Error(c, response.CodeServerError, "保存文件失败")
		return
	}

	// 返回URL
	url := fmt.Sprintf("/uploads/%s/%s", date, filename)

	response.SuccessWithMessage(c, "上传成功", gin.H{
		"url": url,
	})
}

// UploadGoodsImages 上传商品图片（多图）
func (h *UploadHandler) UploadGoodsImages(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	// 获取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, response.CodeParamError, "获取文件失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, response.CodeParamError, "请选择文件")
		return
	}

	// 最多9张图片
	if len(files) > 9 {
		response.Error(c, response.CodeParamError, "最多上传9张图片")
		return
	}

	var urls []string
	date := time.Now().Format("2006-01-02")
	dir := filepath.Join(h.uploadPath, date)
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, response.CodeServerError, "创建目录失败")
		return
	}

	// 保存每个文件
	for _, fileHeader := range files {
		// 验证文件大小
		if fileHeader.Size > h.maxSize {
			response.Error(c, response.CodeParamError,
				fmt.Sprintf("文件 %s 大小超过限制", fileHeader.Filename))
			return
		}

		// 验证文件类型
		if !h.isAllowedImageType(fileHeader) {
			response.Error(c, response.CodeParamError,
				fmt.Sprintf("文件 %s 格式不支持", fileHeader.Filename))
			return
		}

		// 生成文件名
		ext := filepath.Ext(fileHeader.Filename)
		filename := h.generateFileName(userID, ext)

		// 保存文件
		filePath := filepath.Join(dir, filename)
		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			response.Error(c, response.CodeServerError, "保存文件失败")
			return
		}

		url := fmt.Sprintf("/uploads/%s/%s", date, filename)
		urls = append(urls, url)
	}

	response.SuccessWithMessage(c, "上传成功", gin.H{
		"urls": urls,
	})
}

// isAllowedImageType 检查是否是允许的图片类型
func (h *UploadHandler) isAllowedImageType(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}
	return allowedExts[ext]
}

// generateFileName 生成文件名
func (h *UploadHandler) generateFileName(userID interface{}, ext string) string {
	// 格式：{userID}_{timestamp}_{randomUUID}.ext
	var userIDStr string
	if userID != nil {
		userIDStr = strconv.FormatInt(userID.(int64), 10)
	} else {
		userIDStr = "0"
	}

	timestamp := time.Now().Unix()
	randomUUID := uuid.New().String()[:8]

	return fmt.Sprintf("%s_%d_%s%s", userIDStr, timestamp, randomUUID, ext)
}

// ValidateImage 验证图片文件内容（可选）
func (h *UploadHandler) ValidateImage(file multipart.File) error {
	// 读取前512字节用于检测类型
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// 这里可以添加更详细的文件内容验证
	// 例如检查图片的真实格式、尺寸等

	return nil
}

// DeleteFile 删除文件
func (h *UploadHandler) DeleteFile(filePath string) error {
	// 构建完整路径
	fullPath := filepath.Join(h.uploadPath, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	// 删除文件
	return os.Remove(fullPath)
}

// SaveFromURL 从URL保存图片（可选，用于后续扩展）
func (h *UploadHandler) SaveFromURL(url string, userID int64) (string, error) {
	// TODO: 实现从URL下载图片并保存
	// 这可以用于用户输入图片URL的场景
	return "", nil
}

// GetFileInfo 获取文件信息
func (h *UploadHandler) GetFileInfo(filePath string) (os.FileInfo, error) {
	fullPath := filepath.Join(h.uploadPath, filePath)
	return os.Stat(fullPath)
}

// ReadFile 读取文件
func (h *UploadHandler) ReadFile(filePath string) (*os.File, error) {
	fullPath := filepath.Join(h.uploadPath, filePath)
	return os.Open(fullPath)
}

// CopyFile 复制文件
func (h *UploadHandler) CopyFile(src, dst string) error {
	sourceFile, err := h.ReadFile(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destDir := filepath.Dir(filepath.Join(h.uploadPath, dst))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(filepath.Join(h.uploadPath, dst))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
