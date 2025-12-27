package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm/schema"
)

// PgCrudStorage is the base storage class that provides common functionalities which all stores can benefit from.
// Please add your common storage logic here.
type PgCrudStorage[M schema.Tabler] interface {
	CreateOne(model M) error
	CreateMany(models []M) error
	CreateManyWithAssociation(models []M, saveAssociations bool) error
	CreateInBatches(models []M) error

	FindById(id uuid.UUID) (model M, err error)
	ListByIds(ids []uuid.UUID) (models []M, err error)

	UpdateOne(model M, updateZeroValues bool) error
	UpdateMany(models []M) error

	UpsertOne(model M, saveAssociations bool) error
	UpdatePartial(model M, returnUpdated bool) error
	UpsertMany(models []M) error
	UpsertManyWithAssociation(models []M, saveAssociations bool) error
	UpsertInBatches(models []M) error

	ExistsById(id uuid.UUID) (exists bool, err error)

	DeleteOne(model M) error
	DeleteMany(models []M) error
	DeleteById(id uuid.UUID) error
	DeleteByIds(ids []uuid.UUID) error

	ListAll() (models []M, err error)
}
