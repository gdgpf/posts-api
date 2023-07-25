package postservice

import (
	"github.com/gdgpf/posts-api/core/domain"
	"github.com/go-playground/validator/v10"
)

type service struct {
	usecase   domain.PostUseCase
	validator *validator.Validate
}

func New(usecase domain.PostUseCase, validator *validator.Validate) domain.PostHTTPService {
	return &service{
		usecase:   usecase,
		validator: validator,
	}
}
