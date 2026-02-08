package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type SubscriptionStorage interface {
	PgCrudStorage[*model.Subscription]
	FindByPaddleID(paddleID string) (*model.Subscription, error)
}
