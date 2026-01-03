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
