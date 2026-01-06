package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"gorm.io/gorm"
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
	err := s.db.
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		First(res, "handler = ?", handler).Error
	if err != nil {
		return nil, err
	}

	if len(res.Histories) > 0 {
		res.LastHistory = res.Histories[0]
	}

	return res, nil
}
