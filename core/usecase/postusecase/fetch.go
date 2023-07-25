package postusecase

import (
	"github.com/gdgpf/posts-api/core/domain"
)

func (usecase usecase) Fetch() ([]domain.Post, error) {
	return usecase.repository.Fetch()
}
