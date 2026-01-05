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

func (s *ChannelStg) FindByHandler(handler string) (*model.Channel, error) {
	res := &model.Channel{}
	err := s.db.First(res, "handler = ?", handler).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
