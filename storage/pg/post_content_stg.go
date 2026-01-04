package pg

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostContentStg struct {
	crudStg[*model.PostContent]
}

func NewPostContentStg(ses *ormSession) *PostContentStg {
	return &PostContentStg{
		crudStg: crudStg[*model.PostContent]{db: ses.db},
	}
}

func (s *PostContentStg) ListByPostId(postId string) ([]*model.PostContent, error) {
	var list []*model.PostContent
	err := s.db.Where("post_id = ?", postId).
		Find(&list).
		Order("created_at asc").
		Error
	return list, err
}
