package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/google/uuid"
)

type UserStorage interface {
	PgCrudStorage[*model.User]
	FindByEmail(email string) (*model.User, error)
	FindById(id uuid.UUID) (*model.User, error)
	FindByProvider(provider model.Provider, providerID string) (*model.User, error)
}
