package storage

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostAnalysisStorage interface {
	PgCrudStorage[*model.PostAnalysis]
}
