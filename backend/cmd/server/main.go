package main

import (
	"campus-swap-shop/internal/handler"
	"campus-swap-shop/internal/middleware"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"campus-swap-shop/pkg/websocket"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化服务上下文
	serviceCtx, err := svc.NewServiceContext()
	if err != nil {
		panic(fmt.Sprintf("初始化服务上下文失败: %v", err))
	}

	// Gin设置
	mode := serviceCtx.Config.GetString("server.mode")
	gin.SetMode(mode)
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery(serviceCtx.Logger))

	// 静态文件服务（上传的文件）
	r.Static("/uploads", serviceCtx.Config.GetString("upload.path"))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// API路由
	v1 := r.Group("/api/v1")
	{
		// 初始化WebSocket管理器
		wsManager := websocket.NewWebSocketManager()

		// 初始化Handler
		userHandler := handler.NewUserHandler(serviceCtx)
		goodsHandler := handler.NewGoodsHandler(serviceCtx)
		uploadHandler := handler.NewUploadHandler(serviceCtx)
		exchangeHandler := handler.NewExchangeHandler(serviceCtx)
		messageHandler := handler.NewMessageHandler(serviceCtx, wsManager)
		favoriteHandler := handler.NewFavoriteHandler(serviceCtx)
		commentHandler := handler.NewCommentHandler(serviceCtx)
		notificationHandler := handler.NewNotificationHandler(serviceCtx)
		reportHandler := handler.NewReportHandler(serviceCtx)

		// 公开接口：用户相关
		user := v1.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/:id", userHandler.GetUserInfoByID)
			user.GET("/:id/goods", userHandler.GetUserGoods)
		}

		// 公开接口：商品相关
		goods := v1.Group("/goods")
		{
			goods.GET("/list", goodsHandler.GetList)
			goods.GET("/detail", goodsHandler.GetDetail)
			goods.GET("/search", goodsHandler.Search)
		}

		// 公开接口：分类
		category := v1.Group("/category")
		{
			category.GET("/list", goodsHandler.GetCategories)
		}

		// 需要认证的接口：用户相关
		auth := v1.Group("")
		auth.Use(middleware.Auth(serviceCtx.JwtSecret))
		{
			user := auth.Group("/user")
			{
				user.GET("/info", userHandler.GetUserInfo)
				user.PUT("/info", userHandler.UpdateUserInfo)
				user.POST("/change-password", userHandler.ChangePassword)
				user.POST("/avatar", uploadHandler.UploadImage)
			}

			// 需要认证的接口：商品相关
			goods := auth.Group("/goods")
			{
				goods.POST("/create", goodsHandler.Create)
				goods.PUT("/update", goodsHandler.Update)
				goods.DELETE("/delete", goodsHandler.Delete)
				goods.GET("/my", goodsHandler.GetMyGoods)
			}

			// 需要认证的接口：上传
			upload := auth.Group("/upload")
			{
				upload.POST("/image", uploadHandler.UploadImage)
				upload.POST("/goods-images", uploadHandler.UploadGoodsImages)
			}

			// 需要认证的接口：交换相关
			exchange := auth.Group("/exchange")
			{
				exchange.POST("/create", exchangeHandler.Create)
				exchange.POST("/handle", exchangeHandler.Handle)
				exchange.GET("/list", exchangeHandler.GetList)
				exchange.GET("/detail", exchangeHandler.GetDetail)
				exchange.GET("/pending-count", exchangeHandler.GetPendingCount)
			}

			// 需要认证的接口：消息相关
			message := auth.Group("/message")
			{
				message.GET("/list", messageHandler.GetMessageList)
				message.GET("/conversations", messageHandler.GetConversations)
				message.GET("/unread-count", messageHandler.GetUnreadCount)
				message.GET("/online-status", messageHandler.CheckOnline)
				message.POST("/send", messageHandler.SendMessage)
			}

			// 需要认证的接口：收藏相关
			favorite := auth.Group("/favorite")
			{
				favorite.POST("/add", favoriteHandler.AddFavorite)
				favorite.DELETE("/remove/:goodsId", favoriteHandler.RemoveFavorite)
				favorite.GET("/check", favoriteHandler.CheckFavorite)
				favorite.GET("/list", favoriteHandler.GetFavoriteList)
				favorite.GET("/count", favoriteHandler.GetFavoriteCount)
			}

			// 需要认证的接口：评论相关
			comment := auth.Group("/comment")
			{
				comment.POST("/create", commentHandler.CreateComment)
				comment.DELETE("/:id", commentHandler.DeleteComment)
				comment.GET("/my", commentHandler.GetUserComments)
				comment.GET("/received", commentHandler.GetReceivedComments)
			}

			// 需要认证的接口：通知相关
			notification := auth.Group("/notification")
			{
				notification.GET("/list", notificationHandler.GetNotificationList)
				notification.GET("/unread-count", notificationHandler.GetUnreadCount)
				notification.PUT("/read/:id", notificationHandler.MarkAsRead)
				notification.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notification.DELETE("/:id", notificationHandler.DeleteNotification)
				notification.DELETE("/clear-read", notificationHandler.ClearReadNotifications)
			}

			// 需要认证的接口：举报相关
			report := auth.Group("/report")
			{
				report.POST("/create", reportHandler.CreateReport)
				report.GET("/my", reportHandler.GetMyReports)
				report.DELETE("/:id", reportHandler.CancelReport)

				// 管理员接口（需要管理员权限，这里简化为需要登录）
				report.GET("/list", reportHandler.GetReportList)
				report.GET("/detail/:id", reportHandler.GetReportDetail)
				report.PUT("/handle/:id", reportHandler.HandleReport)
				report.GET("/pending-count", reportHandler.GetPendingCount)
			}
		}

		// 公开接口：评论相关
		comment := v1.Group("/comment")
		{
			comment.GET("/list", commentHandler.GetCommentList)
			comment.GET("/rating", commentHandler.GetGoodsRating)
		}

		// WebSocket接口
		ws := v1.Group("/ws")
		ws.Use(middleware.Auth(serviceCtx.Config.GetString("jwt.secret")))
		{
			ws.GET("/chat", func(c *gin.Context) {
				// WebSocket连接需要JWT认证
				userID, exists := c.Get("user_id")
				if !exists || userID == nil {
					c.JSON(200, gin.H{"code": 10002, "message": "未登录"})
					return
				}
				wsManager.HandleWebSocket(c)
			})
		}
	}

	// 启动服务器
	port := serviceCtx.Config.GetString("server.port")
	serviceCtx.Logger.Info("服务器启动",
		zap.String("port", port),
		zap.String("mode", mode),
	)

	if err := r.Run(":" + port); err != nil {
		serviceCtx.Logger.Fatal("服务器启动失败", zap.Error(err))
	}
}
