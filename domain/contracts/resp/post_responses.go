package resp

import (
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/google/uuid"
)

type PostQueueResponse struct {
	Id            uuid.UUID `json:"id"`
	EstimatedTime int       `json:"estimatedTime"`
}

type PostContentResponse struct {
	Content  string `json:"content"`
	Language string `json:"language"`
}

type PostContentSegmentResponse struct {
	*PostContentResponse `json:",inline"`
	Timestamp            string `json:"timestamp"`
	Emotion              string `json:"emotion"`
	Speaker              string `json:"speaker"`
}

type PostResponse struct {
	Status model.PostStatus `json:"status"`

	ImageUrl *string              `json:"imageUrl,omitempty"`
	VideoUrl *string              `json:"videoUrl,omitempty"`
	Caption  *string              `json:"caption,omitempty"`
	UserLink *string              `json:"userLink,omitempty"`
	Platform model.SocialPlatform `json:"platform"`

	Hook      *string                       `json:"hook,omitempty"`
	Segments  []*PostContentSegmentResponse `json:"segments,omitempty"`
	KeyPoints []*PostContentResponse        `json:"keyPoints,omitempty"`
	Summary   *PostContentResponse          `json:"summary,omitempty"`
}
