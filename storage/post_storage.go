package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostStorage interface {
	PgCrudStorage[*model.Post]

	FindByUrl(url string) (*model.Post, error)
}
