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

type GiveawayPostContentMetadata struct {
	Requirements string `json:"requirements"`
	Deadline     string `json:"deadline"`
	Prize        string `json:"prize"`
}

func (s *GiveawayPostContentMetadata) MetadataType() PostContentType {
	return ContentGiveaway
}
