package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type ChannelHistoryStg struct {
	crudStg[*model.ChannelHistory]
}

func NewChannelHistoryStg(ses *ormSession) *ChannelHistoryStg {
	return &ChannelHistoryStg{
		crudStg: crudStg[*model.ChannelHistory]{db: ses.db},
	}
}
