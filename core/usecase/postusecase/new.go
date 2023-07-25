package postusecase

import "github.com/gdgpf/posts-api/core/domain"

type usecase struct {
	repository domain.PostRepository
}

func New(repository domain.PostRepository) domain.PostUseCase {
	return &usecase{
		repository: repository,
	}
}
