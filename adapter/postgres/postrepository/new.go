package postrepository

import (
	"github.com/gdgpf/posts-api/core/domain"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	database *sqlx.DB
}

func New(database *sqlx.DB) domain.PostRepository {
	return &repository{
		database: database,
	}
}
