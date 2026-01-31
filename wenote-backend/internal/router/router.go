package router

import (
	"wenote-backend/config"
	"wenote-backend/internal/handler"
	"wenote-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(config.GlobalConfig.Server.Mode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	// 静态文件服务（图片等）
	r.Static("/uploads", "./uploads")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Cloud Note API is running",
		})
	})

	v1 := r.Group("/api/v1")
	{
		authHandler := handler.NewAuthHandler()
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		authorized := v1.Group("")
		authorized.Use(middleware.JWTAuth())
		{
			userHandler := handler.NewUserHandler()
			users := authorized.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PATCH("/me", userHandler.UpdateProfile)
				users.POST("/me/password", userHandler.ChangePassword)
				users.DELETE("/me", userHandler.DeleteAccount)
			}

			notebookHandler := handler.NewNotebookHandler()
			notebooks := authorized.Group("/notebooks")
			{
				notebooks.GET("/default", notebookHandler.GetDefault)
				notebooks.GET("", notebookHandler.List)
				notebooks.POST("", notebookHandler.Create)
				notebooks.GET("/:id", notebookHandler.GetByID)
				notebooks.PATCH("/:id", notebookHandler.Update)
				notebooks.DELETE("/:id", notebookHandler.Delete)
			}

			noteHandler := handler.NewNoteHandler()
			notes := authorized.Group("/notes")
			{
				notes.GET("", noteHandler.List)
				notes.GET("/trash", noteHandler.ListDeleted)
				notes.POST("", noteHandler.Create)
				notes.GET("/:id", noteHandler.GetByID)
				notes.PATCH("/:id", noteHandler.Update)
				notes.DELETE("/:id", noteHandler.Delete)
				notes.POST("/:id/restore", noteHandler.Restore)
				notes.PUT("/:id/tags", noteHandler.UpdateTags)
				notes.PUT("/:id/tags/apply-suggestions", noteHandler.ApplySuggestedTags)
				notes.POST("/:id/ai/generate", noteHandler.GenerateSummaryAndTags)
				notes.POST("/ai/assist", noteHandler.AIAssist) // AI写作助手
				notes.POST("/batch/delete", noteHandler.BatchDelete)
				notes.POST("/batch/restore", noteHandler.BatchRestore)
				notes.POST("/batch/move", noteHandler.BatchMove)
				notes.DELETE("/trash", noteHandler.EmptyTrash)
				// 附件相关路由
				notes.POST("/:id/attachments", handler.NewAttachmentHandler().UploadImage)
				notes.GET("/:id/attachments", handler.NewAttachmentHandler().GetAttachments)
				}

			// 附件删除路由
			attachmentHandler := handler.NewAttachmentHandler()
			attachments := authorized.Group("/attachments")
			{
				attachments.DELETE("/:id", attachmentHandler.DeleteAttachment)
			}

		tagHandler := handler.NewTagHandler()
		tags := authorized.Group("/tags")
		{
			tags.GET("", tagHandler.List)
			tags.POST("", tagHandler.Create)
			tags.PATCH("/:id", tagHandler.Update)
			tags.DELETE("/:id", tagHandler.Delete)
		}

			// 统计数据路由
			statsHandler := handler.NewStatsHandler()
			stats := authorized.Group("/stats")
			{
				stats.GET("/overview", statsHandler.GetOverview)
				stats.GET("/trend", statsHandler.GetTrendData)
				stats.GET("/tags", statsHandler.GetTagStats)
				stats.GET("/notebooks", statsHandler.GetNotebookStats)
			}

			// 游戏化路由
			gamificationHandler := handler.NewGamificationHandler()
			gamification := authorized.Group("/gamification")
			{
				gamification.GET("/status", gamificationHandler.GetStatus)
				gamification.GET("/achievements", gamificationHandler.GetAchievements)
				gamification.POST("/goal", gamificationHandler.UpdateGoal)
				gamification.GET("/report", gamificationHandler.GetReport)
				gamification.POST("/achievements/:id/notify", gamificationHandler.MarkAchievementNotified)
			}
	}
}

	return r
}
