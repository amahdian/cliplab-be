package gemini

type AnalysisResponse struct {
	Summary struct {
		BigIdea           string `json:"big_idea"`
		WhyViral          string `json:"why_viral"`
		AudienceSentiment string `json:"audience_sentiment"`
		SentimentScore    int    `json:"sentiment_score"`
	} `json:"summary"`
	Content struct {
		Hook     string `json:"hook"`
		Segments []struct {
			Speaker      string `json:"speaker"`
			Timestamp    string `json:"timestamp"`
			Content      string `json:"content"`
			Language     string `json:"language"`
			LanguageCode string `json:"language_code"`
			Emotion      string `json:"emotion"`
		} `json:"segments"`
	} `json:"content"`
	Analysis struct {
		Scope struct {
			Level      string `json:"level"`
			Confidence int    `json:"confidence"`
		} `json:"scope"`
		Metrics []struct {
			Label       string `json:"label"`
			Score       int    `json:"score"`
			Explanation string `json:"explanation"`
			Suggestion  string `json:"suggestion"`
		} `json:"metrics"`
		Strengths  []string `json:"strengths"`
		Weaknesses []string `json:"weaknesses"`
	} `json:"analysis"`
	Remix struct {
		HookIdeas   []string `json:"hook_ideas"`
		ScriptIdeas []string `json:"script_ideas"`
	} `json:"remix"`
	Publish struct {
		Captions struct {
			Casual       string `json:"casual"`
			Professional string `json:"professional"`
			Viral        string `json:"viral"`
		} `json:"captions"`
		Hashtags []string `json:"hashtags"`
	} `json:"publish"`
}
