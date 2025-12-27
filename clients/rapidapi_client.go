package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amahdian/cliplab-be/clients/dtos"
	"github.com/pkg/errors"
)

type RapidApiClient interface {
	GetInstagramPost(shortcode string) ([]*dtos.InstagramItem, error)
}

type rapidApiClient struct {
	Token      string
	HTTPClient *http.Client
}

func NewRapidApiClient(token string) RapidApiClient {
	return &rapidApiClient{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *rapidApiClient) GetInstagramPost(shortcode string) ([]*dtos.InstagramItem, error) {
	requestBody := map[string]interface{}{
		"shortcode": shortcode,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	endpoint := "https://instagram120.p.rapidapi.com/api/instagram/mediaByShortcode"
	resp, err := c.doPost(endpoint, body, "instagram120.p.rapidapi.com", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 8. Handle non-successful status codes
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("scrap request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 9. Decode the JSON response
	var scrapResp []*dtos.InstagramItem
	if err := json.NewDecoder(resp.Body).Decode(&scrapResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode transcription response")
	}

	return scrapResp, nil
}

func (c *rapidApiClient) doPost(endpoint string, body []byte, host string, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("x-rapidapi-key", c.Token)
	req.Header.Set("x-rapidapi-host", host)

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
