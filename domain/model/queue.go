package model

import "github.com/google/uuid"

type PostQueueData struct {
	Id       uuid.UUID      `json:"id"`
	PostId   *string        `json:"postId"`
	Url      string         `json:"url"`
	Platform SocialPlatform `json:"platform"`
}
