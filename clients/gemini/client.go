package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/errs"
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
Act as a senior %s Content & Growth Analyst with a critical, data-driven mindset.
Your task is to analyze the provided %s honestly and precisely.
Do NOT hype. Do NOT add extra sections. Do NOT invent fields.

You MUST return ONLY a valid JSON object that EXACTLY matches the schema provided below.
Any deviation, extra field, missing field, or rewording of keys is NOT allowed.

[CONTEXT DATA]:
- Video Caption: %s
- Co-Authors: %s
- Viedo Engagement Stats: %s
- Page Average Engagement Stats (last posts): %s
- Audience Comments (sample): %s
- Timing: Published at %s (Current Time: %s)
- Target Region: %s

--------------------------------
MANDATORY ANALYSIS RULES
--------------------------------

1. BASELINE COMPARISON (CRITICAL)
- Always evaluate this post relative to the page’s own historical averages.
- If Views, Likes, or Comments are BELOW page average, you MUST reflect this negatively
  in scores, explanations, and weaknesses.

2. ENGAGEMENT QUALITY
- Distinguish between:
  - CTA-driven comments (e.g. repeated single-word comments)
  - Organic engagement (opinions, emotional reactions, discussion)
- High comment count alone does NOT mean virality.

3. VALUE CLARITY
- If the video mainly validates emotions and delays real value to an external offer,
  reflect this clearly in Value Delivery scoring and suggestions.

4. VIRALITY HONESTY
- Do NOT label content as viral unless it meaningfully exceeds page averages
  or clearly penetrates non-follower feeds.
- Funnel effectiveness ≠ Virality.

5. TOPIC SCORING
- When searching for current trends or waves, ALWAYS take into account:
  - Publish Time %s
  - Target Region %s
  - Any recognizable personalities or celebrities detected in the video frames
- Only assign Topic Score 90+ if Web Search confirms a current cultural wave/event relevant to that time and region.
- If no wave is detected, classify topic as Evergreen/Saturated and score conservatively (≤80).

--------------------------------
SEARCH STRATEGY (REQUIRED)
--------------------------------

- Analyze the video frames to detect any recognizable personalities or celebrities.
- Use Google Search (or another reliable source) to verify:
  1. Whether the topic is currently trending **at the time of publish** (%s) in the **target region** (%s).
  2. Whether similar narratives, challenges, or sounds are peaking or declining in that region/time.
  3. If the content is riding a sound, challenge, or event wave relevant to the region/time.
- If no wave is detected → classify as Evergreen/Saturated.

--------------------------------
OUTPUT FORMAT (STRICT)
--------------------------------

Return ONLY the following JSON structure.
Use clear, concise language inside values.
Scores must be realistic and justified.

{
  "summary": {
    "big_idea": "The core message/value proposition.",
    "why_viral": "Explain whether it truly went viral or what psychological trigger it relies on instead.",
    "audience_sentiment": "Deep analysis of how the audience emotionally and cognitively reacted.",
    "sentiment_score": 0-100
  },
  "content": {
    "hook": "Detailed analysis of the first 3 seconds including visual hook, verbal hook, and any pattern interruption.",
    "segments": [
      {
        "speaker": "Identity (e.g. Creator, Narrator)",
        "timestamp": "[MM:SS]",
        "content": "Transcribed or summarized spoken content",
        "emotion": "happy | sad | angry | neutral | anxious | hopeful"
      }
    ]
  },
  "analysis": {
    "metrics": [
      {
        "label": "Hook Strength | Topic Potential | Pacing | Value Delivery | Shareability | CTA",
        "score": 0-100,
        "explanation": "Data-backed rationale for the score.",
        "suggestion": "Specific and actionable improvement."
      }
    ],
    "strengths": [
      "Clear, concrete strengths based on data and structure"
    ],
    "weaknesses": [
      "Clear, concrete weaknesses based on performance and saturation"
    ]
  },
  "remix": {
    "hook_ideas": [
      "3 alternative opening hooks that are sharper or more disruptive"
    ],
    "script_ideas": [
      "3 alternative script angles or narratives"
    ]
  },
  "publish": {
    "captions": {
      "casual": "Conversational caption",
      "professional": "Clean, authority-based caption",
      "viral": "Short, punchy, curiosity-driven caption"
    },
    "hashtags": [
      "5–10 relevant and currently popular hashtags for the region"
    ]
  }
}

--------------------------------
FINAL INSTRUCTIONS
--------------------------------

- Output ONLY valid JSON.
- Do NOT include markdown, explanations, or commentary outside the JSON.
- Do NOT add, rename, or remove fields.
- Be analytical, not motivational.
- Assume the reader will use this output programmatically.`,
		platform, platform, caption, strings.Join(coauthors, "|"),
		statsContext, averageStatsContext, strings.Join(comments, "|"),
		pubTime, currentTime, targetRegion, pubTime, targetRegion, pubTime, targetRegion,
	)

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
		var apiErr struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error.Code != 0 {
			return nil, errs.Newf(errs.Internal, err, apiErr.Error.Message)
		}

		return nil, errors.New("google api returned no candidates and no error details")
	}

	// 5. Unmarshal the actual JSON string from the response part into our struct
	var finalResult AnalysisResponse
	actualJson := googleResp.Candidates[0].Content.Parts[0].Text

	// Clean up markdown markers if present
	actualJson = strings.ReplaceAll(actualJson, "```json", "")
	actualJson = strings.ReplaceAll(actualJson, "```", "")
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
