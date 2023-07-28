package postusecase

import (
	"github.com/gdgpf/posts-api/core/domain"
)

func (usecase usecase) Like(postID int) (*domain.Post, error) {
	return usecase.repository.Like(postID)
}
