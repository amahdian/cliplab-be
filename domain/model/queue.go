package model

type PostQueueData struct {
	Id       string         `json:"id"`
	Url      string         `json:"url"`
	Platform SocialPlatform `json:"platform"`
}
