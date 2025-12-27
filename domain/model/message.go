package model

type StreamEventType string

const (
	StreamEventTypeResponse = "response"
)

type StreamedMessage struct {
	Type    StreamEventType `json:"type"`
	Content string          `json:"content"`
}
