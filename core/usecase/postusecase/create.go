package postusecase

import (
	"github.com/gdgpf/posts-api/core/domain"
	"github.com/gdgpf/posts-api/core/dto"
)

func (usecase usecase) Create(post *dto.PostRequestBody) (*domain.Post, error) {
	return usecase.repository.Create(post)
}
