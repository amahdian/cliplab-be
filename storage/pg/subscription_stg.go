package pg

import (
	"errors"

	"github.com/amahdian/cliplab-be/domain/model"
	"gorm.io/gorm"
)

type SubscriptionStg struct {
	crudStg[*model.Subscription]
}

func NewSubscriptionStg(ses *ormSession) *SubscriptionStg {
	return &SubscriptionStg{
		crudStg: crudStg[*model.Subscription]{db: ses.db},
	}
}

func (stg *SubscriptionStg) FindByPaddleID(paddleID string) (*model.Subscription, error) {
	var sub model.Subscription
	err := stg.db.Where("paddle_subscription_id = ?", paddleID).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}
