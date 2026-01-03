package pg

import (
	"errors"

	"github.com/amahdian/cliplab-be/domain/model"
	"gorm.io/gorm"
)

type UserStg struct {
	crudStg[*model.User]
}

func NewUserStg(ses *ormSession) *UserStg {
	return &UserStg{
		crudStg: crudStg[*model.User]{db: ses.db},
	}
}

func (stg *UserStg) FindByEmail(email string) (user *model.User, err error) {
	err = stg.db.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return
}

func (stg *UserStg) FindByProvider(provider model.Provider, providerID string) (*model.User, error) {
	var user model.User
	result := stg.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, result.Error
}
