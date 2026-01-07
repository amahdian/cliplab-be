package rocksolid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type Client interface {
	GetInstagramPost(shortcode string) (*ReelData, error)
	GetInstagramPageReels(shortcode string) (*Reels, error)
}

type client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) Client {
	return &client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *client) GetInstagramPost(shortcode string) (*ReelData, error) {
	endpoint := "https://auto-poster.co.uk/yt_api/get_media_data_v2.php"
	url := fmt.Sprintf("%s?media_code=%s", endpoint, shortcode)
	resp, err := c.doGet(url, nil)
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
	var scrapResp *ReelData
	if err := json.NewDecoder(resp.Body).Decode(&scrapResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode transcription response")
	}

	return scrapResp, nil
}

func (c *client) GetInstagramPageReels(username string) (*Reels, error) {
	form := url.Values{}
	form.Set("username_or_url", username)
	form.Set("amount", "30")
	form.Set("pagination_token", "")

	endpoint := "https://auto-poster.co.uk/yt_api/get_ig_user_reels.php"
	resp, err := c.doPost(
		endpoint,
		[]byte(form.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"scrap request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var scrapResp Reels
	if err := json.NewDecoder(resp.Body).Decode(&scrapResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &scrapResp, nil
}

func (c *client) doPost(
	endpoint string,
	body []byte,
	headers map[string]string,
) (*http.Response, error) {

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("AP_API_KEY", c.Token)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	return resp, nil
}

func (c *client) doGet(endpoint string, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("AP_API_KEY", c.Token)

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
