package resp

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type PostQueueResponse struct {
	Id            string `json:"id"`
	EstimatedTime int    `json:"estimatedTime"`
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
	UserLink *string              `json:"userLink,omitempty"`
	Platform model.SocialPlatform `json:"platform"`

	Caption  *PostContentResponse          `json:"caption,omitempty"`
	Segments []*PostContentSegmentResponse `json:"segments,omitempty"`

	Analysis *model.PostAnalysis `json:"analysis"`
}
