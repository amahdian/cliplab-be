package pg

import (
	"net"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
)

type AnalyzeRequestStg struct {
	crudStg[*model.AnalyzeRequest]
}

func NewAnalyzeRequest(ses *ormSession) *AnalyzeRequestStg {
	return &AnalyzeRequestStg{
		crudStg: crudStg[*model.AnalyzeRequest]{db: ses.db},
	}
}

func (s *AnalyzeRequestStg) FindByUrl(url string) (*model.AnalyzeRequest, error) {
	res := &model.AnalyzeRequest{}
	err := s.db.First(res, "link = ?", url).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AnalyzeRequestStg) ListByPostId(id string) ([]*model.AnalyzeRequest, error) {
	var res []*model.AnalyzeRequest
	err := s.db.Find(&res, "post_id = ?", id).Order("updated_at desc").Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AnalyzeRequestStg) CountByIpAndDate(ip net.IP, date time.Time) (int64, error) {
	var count int64

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	err := s.db.
		Model(&model.Post{}).
		Where("user_ip = ? AND updated_at >= ? AND updated_at < ?", ip, start, end).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
