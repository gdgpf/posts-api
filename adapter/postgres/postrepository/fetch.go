package postrepository

import (
	"github.com/gdgpf/posts-api/core/domain"
)

func (repository repository) Fetch() ([]domain.Post, error) {
	posts := []domain.Post{}

	err := repository.database.Select(
		&posts,
		"SELECT * FROM post",
	)

	if err != nil {
		return nil, err
	}

	return posts, nil
}
