package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostContentStorage interface {
	PgCrudStorage[*model.PostContent]

	ListByPostId(postId string) ([]*model.PostContent, error)
}
