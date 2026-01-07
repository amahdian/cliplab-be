package pg

import (
	"net"
	"time"

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

func (s *PostStg) CountByIpAndDate(ip net.IP, date time.Time) (int64, error) {
	var count int64

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	err := s.db.
		Model(&model.Post{}).
		Where("user_ip = ? AND created_at >= ? AND created_at < ?", ip, start, end).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
