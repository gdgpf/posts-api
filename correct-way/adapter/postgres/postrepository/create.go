package postrepository

import (
	"github.com/gdgpf/posts-api/core/domain"
	"github.com/gdgpf/posts-api/core/dto"
)

func (repository repository) Create(post *dto.PostRequestBody) (*domain.Post, error) {
	postCreated := domain.Post{}

	err := repository.database.QueryRowx(
		"INSERT INTO post (title, description, latitude, longitude) VALUES ($1, $2, $3, $4) returning *",
		post.Title, post.Description, post.Latitude, post.Longitude,
	).StructScan(&postCreated)

	if err != nil {
		return nil, err
	}

	return &postCreated, nil
}
