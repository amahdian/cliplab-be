package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostAnalysisStg struct {
	crudStg[*model.PostAnalysis]
}

func NewPostAnalysisStg(ses *ormSession) *PostAnalysisStg {
	return &PostAnalysisStg{
		crudStg: crudStg[*model.PostAnalysis]{db: ses.db},
	}
}

func (s *PostAnalysisStg) FindByPostId(id string) (*model.PostAnalysis, error) {
	res := &model.PostAnalysis{}
	err := s.db.First(res, "post_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
