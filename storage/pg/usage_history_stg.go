package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/domain/model/common"
	"github.com/google/uuid"
)

type UsageHistoryStg struct {
	crudStg[*model.UsageHistory]
}

func NewUsageHistoryStg(ses *ormSession) *UsageHistoryStg {
	return &UsageHistoryStg{
		crudStg: crudStg[*model.UsageHistory]{db: ses.db},
	}
}

func (stg *UsageHistoryStg) FindByUserId(userId uuid.UUID, pagination *common.Pagination) (histories []*model.UsageHistory, err error) {
	err = stg.db.
		Scopes(stg.withPagination(pagination)).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&histories).
		Error

	return
}
