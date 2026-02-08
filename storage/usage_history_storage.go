package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/domain/model/common"
	"github.com/google/uuid"
)

type UsageHistoryStorage interface {
	PgCrudStorage[*model.UsageHistory]
	FindByUserId(userId uuid.UUID, pagination *common.Pagination) ([]*model.UsageHistory, error)
}
