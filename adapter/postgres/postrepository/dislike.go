package postrepository

import (
	"github.com/gdgpf/posts-api/core/domain"
)

func (repository repository) Dislike(postID int) (*domain.Post, error) {
	postCreated := domain.Post{}

	err := repository.database.QueryRowx(
		"UPDATE post SET dislikes = dislikes + 1 WHERE id = $1 returning *;",
		postID,
	).StructScan(&postCreated)

	if err != nil {
		return nil, err
	}

	return &postCreated, nil
}
