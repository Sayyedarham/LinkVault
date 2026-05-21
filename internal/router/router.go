package router

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/hphp/linkvault/internal/handler"
	"github.com/hphp/linkvault/internal/middleware"
)

func New(db *dynamodb.Client, tableName string) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handler.Health)

		// Bookmarks (protected)
		bookmarkHandler := handler.NewBookmarkHandler(db, tableName)
		bookmarks := v1.Group("/bookmarks")
		bookmarks.Use(middleware.AuthMiddleware())
		{
			bookmarks.POST("", bookmarkHandler.Create)
			bookmarks.GET("", bookmarkHandler.List)
			bookmarks.DELETE("/:id", bookmarkHandler.Delete)
		}
	}
	return r
}