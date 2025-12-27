package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/google/uuid"
)

type PostContentStorage interface {
	PgCrudStorage[*model.PostContent]

	ListByPostId(postId uuid.UUID) ([]*model.PostContent, error)
}
