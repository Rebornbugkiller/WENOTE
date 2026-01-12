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
				notes.POST("/batch/delete", noteHandler.BatchDelete)
				notes.POST("/batch/restore", noteHandler.BatchRestore)
				notes.POST("/batch/move", noteHandler.BatchMove)
			}

			tagHandler := handler.NewTagHandler()
			tags := authorized.Group("/tags")
			{
				tags.GET("", tagHandler.List)
				tags.POST("", tagHandler.Create)
				tags.PATCH("/:id", tagHandler.Update)
				tags.DELETE("/:id", tagHandler.Delete)
			}
		}
	}

	return r
}
