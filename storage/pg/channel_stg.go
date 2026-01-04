package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type ChannelStg struct {
	crudStg[*model.Channel]
}

func NewChannelStg(ses *ormSession) *ChannelStg {
	return &ChannelStg{
		crudStg: crudStg[*model.Channel]{db: ses.db},
	}
}
