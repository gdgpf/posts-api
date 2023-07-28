package rest

import (
	"github.com/gdgpf/posts-api/di"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

func InitializeRest(configHTTPServerPort string, database *sqlx.DB, validate *validator.Validate) {
	router := gin.Default()

	postService := di.ConfigPostDI(database, validate)

	router.GET("/post", postService.Fetch)
	router.POST("/post", postService.Create)
	router.POST("/post/:id/like", postService.Like)
	router.POST("/post/:id/dislike", postService.Dislike)

	router.Run(":" + configHTTPServerPort)
}
