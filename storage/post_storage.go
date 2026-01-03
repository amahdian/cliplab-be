package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostStorage interface {
	PgCrudStorage[*model.Post]

	FindByHashId(id string) (*model.Post, error)
	FindByUrl(url string) (*model.Post, error)
}
