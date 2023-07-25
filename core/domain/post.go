package domain

import (
	"github.com/gdgpf/posts-api/core/dto"
	"github.com/gin-gonic/gin"
)

type Post struct {
	ID          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Description string   `json:"description" db:"description"`
	Latitude    *float32 `json:"latitude" db:"latitude"`
	Longitude   *float32 `json:"longitude" db:"longitude"`
	Likes       int      `json:"likes" db:"likes"`
	Dislikes    int      `json:"dislikes" db:"dislikes"`
}

type PostRepository interface {
	Create(post *dto.PostRequestBody) (*Post, error)
	Fetch() ([]Post, error)
	Like(int) (*Post, error)
	Dislike(int) (*Post, error)
}

type PostUseCase interface {
	Create(post *dto.PostRequestBody) (*Post, error)
	Fetch() ([]Post, error)
	Like(int) (*Post, error)
	Dislike(int) (*Post, error)
}

type PostHTTPService interface {
	Create(*gin.Context)
	Fetch(*gin.Context)
	Like(*gin.Context)
	Dislike(*gin.Context)
}
