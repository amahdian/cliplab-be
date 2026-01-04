package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type ChannelHistoryStorage interface {
	PgCrudStorage[*model.ChannelHistory]
}
