package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amahdian/cliplab-be/clients/dtos"
	"github.com/pkg/errors"
)

type GeminiClient interface {
	AnalyzeVideo(platform string, videoURL string) (*dtos.AnalysisResponse, error)
}

type geminiClient struct {
	BaseUrl    string
	Token      string
	HTTPClient *http.Client
}

func NewGeminiClient(baseUrl, token string) GeminiClient {
	return &geminiClient{
		BaseUrl: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			// Increased timeout for streaming
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *geminiClient) AnalyzeVideo(platform string, videoURL string) (*dtos.AnalysisResponse, error) {
	// 1. Define the dynamic prompt
	promptText := fmt.Sprintf(`This video is from %s. Analyze the content and provide a structured response:
1. Summary & Key Points: Must be in the SAME language as the video. Focus on actual content, not visuals.
2. Hooks: Identify the opening hook in the video's language.
3. Giveaway: If detected, extract details in English.
4. Trend Analysis Raw Data: Highly condensed summary of topics/keywords in English for future trend detection.
5. Transcription: Provide timestamps and speaker detection.`, platform)

	// 2. Build the Gemini API request body
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"file_data": map[string]string{"file_uri": videoURL, "mime_type": "video/mp4"}},
					{"text": promptText},
				},
			},
		},
		"generation_config": map[string]interface{}{
			"response_mime_type": "application/json",
			"response_schema": map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"summary":        map[string]interface{}{"type": "STRING"},
					"key_points":     map[string]interface{}{"type": "ARRAY", "items": map[string]string{"type": "STRING"}},
					"hook":           map[string]interface{}{"type": "STRING"},
					"trend_metadata": map[string]interface{}{"type": "STRING"},
					"giveaway": map[string]interface{}{
						"type": "OBJECT",
						"properties": map[string]interface{}{
							"is_detected":  map[string]interface{}{"type": "BOOLEAN"},
							"prize":        map[string]interface{}{"type": "STRING"},
							"requirements": map[string]interface{}{"type": "STRING"},
							"deadline":     map[string]interface{}{"type": "STRING"},
						},
						"required": []string{"is_detected"},
					},
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
				"required": []string{"summary", "key_points", "hook", "giveaway", "trend_metadata", "segments"},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	// 3. Call the internal doPost method
	// Endpoint for Gemini 1.5 Flash
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
	var finalResult dtos.AnalysisResponse
	actualJson := googleResp.Candidates[0].Content.Parts[0].Text
	if err := json.Unmarshal([]byte(actualJson), &finalResult); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal analysis result")
	}

	return &finalResult, nil
}

func (c *geminiClient) doPost(endpoint string, body []byte, headers map[string]string) (*http.Response, error) {
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
