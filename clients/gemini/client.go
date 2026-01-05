package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/pkg/errors"
)

type Client interface {
	AnalyzeVideo(
		platform model.SocialPlatform,
		videoURL string,
		caption string,
		coauthors []string,
		comments []string,
		stats map[string]float64,
		averageStats map[string]float64,
		publishedAt time.Time,
		targetRegion string,
	) (*AnalysisResponse, error)
}

type client struct {
	BaseUrl    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseUrl, token string) Client {
	return &client{
		BaseUrl: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			// Increased timeout for streaming
			Timeout: 10 * time.Minute,
		},
	}
}

func (c *client) AnalyzeVideo(
	platform model.SocialPlatform,
	videoURL string,
	caption string,
	coauthors []string,
	comments []string,
	stats map[string]float64,
	averageStats map[string]float64,
	publishedAt time.Time,
	targetRegion string,
) (*AnalysisResponse, error) {

	currentTime := time.Now().Format(time.RFC1123)
	pubTime := publishedAt.Format(time.RFC1123)
	statsContext, _ := json.Marshal(stats)
	averageStatsContext, _ := json.Marshal(averageStats)

	promptText := fmt.Sprintf(`
Act as a world-class Viral Content Strategist. Analyze this %s video and its metadata.

[CONTEXT DATA]:
- Video Caption: %s
- Co-Authors: %s
- Viedo Engagement Stats: %s
- Page Average Engagement Stats: %s
- Audience Comments: %s
- Timing: Published at %s (Current Time: %s)
- Target Region: %s

[SEARCH STRATEGY - MANDATORY]:
Use the Google Search tool to perform a "Contextual Deep Dive":
1. ENTITY CHECK: Search for names, songs, or celebrities mentioned in the video or caption. Check their current viral status.
2. EVENT CORRELATION: Search for major events (concerts, movie releases, trending news) happening in %s around %s that could be driving this video's virality.
3. WAVE DETECTION: Is this video riding a cultural wave? (e.g., a trending sound, a challenge, or a viral news topic).

[SCORING LOGIC FOR TOPIC}:
- If the video is linked to a massive current event (e.g., a Taylor Swift concert happening NOW): Give a 90+ Topic Score.
- If the topic is evergreen but the "Angle" is new: Give 70-80.
- If the video uses a trending song that is just starting to peak: Boost the Topic Score.

[OUTPUT STRUCTURE - JSON ONLY]:
1. summary:
   - big_idea: The core message/value proposition.
   - why_viral: If stats are high, explain the psychological trigger. If not, explain what was missing.
   - audience_sentiment: A deep dive text analysis into how people reacted.
   - sentiment_score: A number (0-100) where 0 is extremely negative, 50 is neutral/mixed, and 100 is extremely positive based on comments.

2. content:
   - hook: The exact opening line or visual action that captures attention (Starting 3 seconds of video).
   - segments: Provide a detailed array of objects, where each object MUST contain:
     "speaker": (Identify the speaker, e.g., "Speaker 1", "Host", or Name if known),
     "timestamp": (Format as [MM:SS]),
     "content": (The actual transcribed text in the original language),
     "emotion": (Choose ONLY from: "happy", "sad", "angry", "neutral" based on the speaker's tone).

3. analysis:
   - Provide a list of "metrics" where each item must contain: 
     "label": (The name of the metric: Hook Strength, Topic Potential, Pacing, Value Delivery, Shareability, or CTA),
     "score": (A number 0-100),
     "explanation": (Why this score was given based on visual and search data),
     "suggestion": (Specific, actionable advice to improve this specific metric for the next video).
   - strengths: 3 or more string item.
   - weaknesses: 3 or more string item.

4. remix:
   - hook_ideas: 3 fresh ways to start the same video.
   - script_ideas: 3 alternative angles to re-tell this story.

5. pulish:
   - captions: Provide Object with: casual, professional, viral fields. (its not an array just one object)
   - hashtags: 5-10 highly relevant and trending hashtags.

IMPORTANT: 
- For visual analysis, focus on frame changes and on-screen text.
- Topic score must be based on "Web Search" data to check if the subject is currently trending.
- Return ONLY the raw JSON object. Do not include any conversational text or explanations outside the JSON block.

[TIERED VIRALITY SCALE]:
- Nano Accounts (<1K Followers): Viral status is ONLY achieved if Views > 5,000 AND Views > (Followers * 10). Below 5k views is considered "Network Reach" (friends & family), not virality.
- Micro Accounts (1K - 10K Followers): Viral status starts if Views > (Followers * 5).
- Mid-Large Accounts (>100K Followers): Viral status starts if Views > (Followers * 1.5) OR if the ER is 2x the average of their last 5 posts.

JUDGMENT RULE: 
Do not give a high "Viral Score" to very small accounts just because their views-to-follower ratio is high, unless they have broken the 5,000 views barrier. Crossing the "Stranger's Feed" (Explore/Reels Tab) is the true mark of virality.`,
		platform, caption, strings.Join(coauthors, "|"), statsContext, averageStatsContext, strings.Join(comments, "|"), pubTime, currentTime, targetRegion, targetRegion, pubTime)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"file_data": map[string]string{
							"file_uri":  videoURL,
							"mime_type": "video/mp4",
						},
						"video_metadata": map[string]interface{}{"fps": 0.5},
					},
					{"text": promptText},
				},
			},
		},
		"tools": []map[string]interface{}{
			{"google_search": map[string]interface{}{}},
		},
		"generation_config": map[string]interface{}{
			//"response_mime_type": "application/json",
			"temperature": 0.2,
			"response_schema": map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"summary": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"big_idea":           map[string]interface{}{"type": "STRING"},
							"why_viral":          map[string]interface{}{"type": "STRING"},
							"audience_sentiment": map[string]interface{}{"type": "STRING"},
							"sentiment_score":    map[string]interface{}{"type": "INTEGER", "description": "0 to 100 scale of audience sentiment"},
						},
					},
					"content": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"hook": map[string]interface{}{"type": "STRING"},
							"segments": map[string]interface{}{
								"type": "ARRAY",
								"items": map[string]interface{}{
									"type": "OBJECT",
									"properties": map[string]interface{}{
										"speaker":   map[string]interface{}{"type": "STRING"},
										"timestamp": map[string]interface{}{"type": "STRING"},
										"content":   map[string]interface{}{"type": "STRING"},
										"emotion":   map[string]interface{}{"type": "STRING", "enum": []string{"happy", "sad", "angry", "neutral"}},
									},
								},
							},
						},
					},
					"analysis": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"metrics": map[string]interface{}{
								"type": "ARRAY",
								"items": map[string]interface{}{
									"type": "OBJECT",
									"properties": map[string]interface{}{
										"label":       map[string]interface{}{"type": "STRING"},
										"score":       map[string]interface{}{"type": "INTEGER"},
										"explanation": map[string]interface{}{"type": "STRING"},
										"suggestion":  map[string]interface{}{"type": "STRING"},
									},
								},
							},
							"strengths":  map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
							"weaknesses": map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
						},
					},
					"remix": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"hook_ideas":   map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
							"script_ideas": map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
						},
					},
					"publish": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"captions": map[string]interface{}{
								"type": "OBJECT",
								"properties": map[string]interface{}{
									"casual":       map[string]interface{}{"type": "STRING"},
									"professional": map[string]interface{}{"type": "STRING"},
									"viral":        map[string]interface{}{"type": "STRING"},
								},
							},
							"hashtags": map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
						},
					},
				},
				"required": []string{"summary", "content", "analysis", "remix", "publish"},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	endpoint := "/v1beta/models/gemini-2.5-flash:generateContent"
	resp, err := c.doPost(endpoint, body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Parse the raw Google response
	var googleResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode google response")
	}

	if len(googleResp.Candidates) == 0 || len(googleResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty response from Gemini API")
	}

	// 5. Unmarshal the actual JSON string from the response part into our struct
	var finalResult AnalysisResponse
	actualJson := googleResp.Candidates[0].Content.Parts[0].Text

	actualJson = strings.TrimPrefix(actualJson, "```json")
	actualJson = strings.TrimSuffix(actualJson, "```")
	actualJson = strings.TrimSpace(actualJson)

	if err := json.Unmarshal([]byte(actualJson), &finalResult); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal analysis result")
	}

	return &finalResult, nil
}

func (c *client) doPost(endpoint string, body []byte, headers map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseUrl, endpoint)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set required headers for all requests
	req.Header.Set("x-goog-api-key", c.Token)

	// Only set Content-Type if not already provided in headers
	if headers == nil || headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set any additional headers (this will override Content-Type if provided)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	return resp, nil
}
