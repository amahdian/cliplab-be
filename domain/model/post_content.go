package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type PostContentType string

const (
	ContentCaption    PostContentType = "caption"
	ContentTranscript PostContentType = "transcript"
)

type PostContent struct {
	ID       uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PostID   string          `json:"postId"`
	Type     PostContentType `json:"type"`
	Language string          `json:"language"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Text string `json:"text"`

	MetadataRaw string              `json:"-" gorm:"serializer:json;column:metadata;"`
	Metadata    PostContentMetadata `json:"metadata" gorm:"-"`
	Embedding   *pgvector.Vector    `json:"embedding"`

	Post *Post `json:"post" gorm:"foreignKey:PostID;references:ID"`
}

func (*PostContent) TableName() string {
	return "post_contents"
}

func (p *PostContent) BeforeSave(_ *gorm.DB) (err error) {
	if p.Metadata == nil {
		return nil
	}

	bytes, err := json.Marshal(p.Metadata)
	if err != nil {
		return err
	}
	p.MetadataRaw = string(bytes)
	return nil
}

func (p *PostContent) AfterFind(_ *gorm.DB) (err error) {
	if p.MetadataRaw == "" {
		return nil
	}

	var contentValue PostContentMetadata
	switch p.Type {
	case ContentTranscript:
		contentValue = &SegmentPostContentMetadata{}
	}
	if err := json.Unmarshal([]byte(p.MetadataRaw), contentValue); err != nil {
		return err
	}
	p.Metadata = contentValue

	return nil
}
