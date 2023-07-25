package di

import (
	"github.com/gdgpf/posts-api/adapter/http/postservice"
	"github.com/gdgpf/posts-api/adapter/postgres/postrepository"
	"github.com/gdgpf/posts-api/core/domain"
	"github.com/gdgpf/posts-api/core/usecase/postusecase"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func ConfigPostDI(database *sqlx.DB, validator *validator.Validate) domain.PostHTTPService {
	postRepository := postrepository.New(database)
	postUseCase := postusecase.New(postRepository)
	postService := postservice.New(postUseCase, validator)

	return postService
}
