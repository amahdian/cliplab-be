package storage

import (
	"net"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
)

type AnalyzeRequestStorage interface {
	PgCrudStorage[*model.AnalyzeRequest]

	FindByUrl(url string) (*model.AnalyzeRequest, error)
	ListByPostId(id string) ([]*model.AnalyzeRequest, error)
	CountByIpAndDate(ip net.IP, date time.Time) (int64, error)
}
