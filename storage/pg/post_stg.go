package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostStg struct {
	crudStg[*model.Post]
}

func NewPostStg(ses *ormSession) *PostStg {
	return &PostStg{
		crudStg: crudStg[*model.Post]{db: ses.db},
	}
}

func (s *PostStg) FindByUrl(url string) (*model.Post, error) {
	res := &model.Post{}
	err := s.db.First(res, "link = ?", url).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
