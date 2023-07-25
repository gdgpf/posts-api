package postusecase

import (
	"github.com/gdgpf/posts-api/core/domain"
)

func (usecase usecase) Dislike(postID int) (*domain.Post, error) {
	return usecase.repository.Dislike(postID)
}
