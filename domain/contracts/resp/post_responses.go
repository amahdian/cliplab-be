package resp

import (
	"time"

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

	ImageUrl    *string              `json:"imageUrl,omitempty"`
	VideoUrl    *string              `json:"videoUrl,omitempty"`
	UserLink    *string              `json:"userLink,omitempty"`
	UserHandler *string              `json:"userHandler,omitempty"`
	Platform    model.SocialPlatform `json:"platform"`

	Caption  *PostContentResponse          `json:"caption,omitempty"`
	Segments []*PostContentSegmentResponse `json:"segments,omitempty"`

	LikeCount      int64     `json:"likeCount,omitempty"`
	CommentCount   int64     `json:"commentCount,omitempty"`
	ViewCount      int64     `json:"viewCount,omitempty"`
	EngagementRate float64   `json:"engagementRate,omitempty"`
	PostDate       time.Time `json:"postDate,omitempty"`

	AverageLikeCount      int64   `json:"averageLikeCount,omitempty"`
	AverageCommentCount   int64   `json:"averageCommentCount,omitempty"`
	AverageViewCount      int64   `json:"averageViewCount,omitempty"`
	AverageEngagementRate float64 `json:"averageEngagementRate,omitempty"`

	Analysis *model.PostAnalysis `json:"analysis,omitempty"`
}
