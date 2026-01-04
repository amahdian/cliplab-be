package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type ChannelStorage interface {
	PgCrudStorage[*model.Channel]
}
