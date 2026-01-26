package model

import (
	"time"

	"github.com/google/uuid"
)

type PostAnalysis struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PostId string    `json:"postId" `

	ViralScore        int         `json:"viralScore"`
	BigIdea           string      `json:"bigIdea"`
	WhyViral          string      `json:"whyViral"`
	AudienceSentiment string      `json:"audienceSentiment"`
	SentimentScore    int         `json:"sentimentScore"`
	Verdict           PostVerdict `json:"verdict"`

	Metrics    []PostAnalysisMetric `json:"metrics" gorm:"serializer:json"`
	Strengths  []string             `json:"strengths" gorm:"serializer:json"`
	Weaknesses []string             `json:"weaknesses" gorm:"serializer:json"`

	HookIdeas   []string `json:"hookIdeas" gorm:"serializer:json"`
	ScriptIdeas []string `json:"scriptIdeas" gorm:"serializer:json"`

	Captions PostAnalysisCaptions `json:"captions" gorm:"serializer:json"`
	Hashtags []string             `json:"hashtags" gorm:"serializer:json"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Post *Post `json:"-" gorm:"foreignKey:PostId;references:ID"`
}

type PostAnalysisMetric struct {
	Label       string `json:"label"`
	Score       int    `json:"score"`
	Explanation string `json:"explanation"`
	Suggestion  string `json:"suggestion"`
}

type PostAnalysisCaptions struct {
	Casual       string `json:"casual"`
	Professional string `json:"professional"`
	Viral        string `json:"viral"`
}

type PostVerdict struct {
	Status    string `json:"status"`
	Reasoning string `json:"reasoning"`
}

func (*PostAnalysis) TableName() string {
	return "post_analyses"
}
