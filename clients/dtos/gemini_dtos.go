package dtos

type AnalysisResponse struct {
	Summary       string   `json:"summary"`
	KeyPoints     []string `json:"key_points"`
	Hook          string   `json:"hook"`
	TrendMetadata string   `json:"trend_metadata"`
	Giveaway      struct {
		IsDetected   bool   `json:"is_detected"`
		Prize        string `json:"prize"`
		Requirements string `json:"requirements"`
		Deadline     string `json:"deadline"`
	} `json:"giveaway"`
	Segments []struct {
		Speaker      string `json:"speaker"`
		Timestamp    string `json:"timestamp"`
		Content      string `json:"content"`
		Language     string `json:"language"`
		LanguageCode string `json:"language_code"`
		Emotion      string `json:"emotion"`
	} `json:"segments"`
}
