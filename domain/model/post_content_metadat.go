package model

type PostContentMetadata interface {
	MetadataType() PostContentType
}

type SegmentPostContentMetadata struct {
	Timestamp string `json:"timestamp"`
	Speaker   string `json:"speaker"`
	Emotion   string `json:"emotion"`
}

func (s *SegmentPostContentMetadata) MetadataType() PostContentType {
	return ContentTranscript
}
