package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global"
)

type DatasetId string

const (
	DatasetIdInstagram DatasetId = "gd_lk5ns7kz21pck8jpis"
	DatasetIdTwitter   DatasetId = "gd_lwxkxvnf1cynvib9co"
	DatasetIdTiktok    DatasetId = "gd_lu702nij2f790tmv9h"
	DatasetIdYoutube   DatasetId = "gd_lk56epmy2i5g7lzu0k"
)

type BrightDataClient interface {
	TriggerScrape(postUrl, postId string, platform model.SocialPlatform) error
}

type brightDataClientImpl struct {
	BaseUrl    string
	Token      string
	ApiHost    string
	HTTPClient *http.Client
}

func NewBrightDataClient(baseUrl, token, apiHost string) BrightDataClient {
	return &brightDataClientImpl{
		BaseUrl: baseUrl,
		Token:   token,
		ApiHost: apiHost,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *brightDataClientImpl) TriggerScrape(postUrl, postId string, platform model.SocialPlatform) error {
	endpoint := fmt.Sprintf("%s/datasets/v3/trigger", c.BaseUrl)
	webhook := fmt.Sprintf("%s%s/webhook/brightdata/%s", c.ApiHost, global.ApiPrefix, postId)

	datasetId := getDatasetId(platform)

	// Construct query string
	q := url.Values{}
	q.Set("dataset_id", string(datasetId))
	q.Set("endpoint", webhook)
	q.Set("format", "json")
	q.Set("uncompressed_webhook", "true")
	q.Set("include_errors", "true")

	// Request body
	payload := make([]map[string]string, 0)
	payload = append(payload, map[string]string{"url": postUrl})

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint+"?"+q.Encode(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("BrightData returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func getDatasetId(platform model.SocialPlatform) DatasetId {
	switch platform {
	case model.PlatformInstagram:
		return DatasetIdInstagram
	case model.PlatformTwitter:
		return DatasetIdTwitter
	case model.PlatformTikTok:
		return DatasetIdTiktok
	case model.PlatformYouTube:
		return DatasetIdYoutube
	default:
		return ""
	}
}
