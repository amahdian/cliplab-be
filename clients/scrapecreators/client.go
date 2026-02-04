package scrapecreators

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
	GetInstagramPost(postURL string) (*ReelData, error)
	GetInstagramPageReels(username string) (*ReelsResponse, error)
	GetInstagramPagePosts(username string) (*PostsResponse, error)
	GetTikTokVideo(videoURL string) (*TikTokAwemeDetail, error)
	GetTikTokProfileVideos(username string) (*TikTokProfileVideosResponse, error)
	GetYouTubeVideo(videoURL string) (*YouTubeVideoResponse, error)
	GetTweet(tweetURL string) (*TwitterTweetResponse, error)
	GetUserTweets(username string) (*TwitterUserTweetsResponse, error)
	GetInstagramProfile(username string) (*InstagramProfileResponse, error)
}

type client struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(token string) Client {
	return &client{
		Token:   token,
		BaseURL: "https://api.scrapecreators.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) GetInstagramPost(postURL string) (*ReelData, error) {
	endpoint := fmt.Sprintf("%s/v1/instagram/post", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("url", postURL)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result Response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators API returned success=false: %s", result.Status)
	}

	return &result.Data.XdtShortcodeMedia, nil
}

func (c *client) GetInstagramPageReels(username string) (*ReelsResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/instagram/user/reels", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("handle", username)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators reels request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result ReelsResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode reels response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators reels API returned success=false")
	}

	return &result, nil
}

func (c *client) GetInstagramPagePosts(username string) (*PostsResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/instagram/user/posts", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("handle", username)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators posts request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result PostsResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode posts response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators posts API returned success=false")
	}

	return &result, nil
}

func (c *client) GetTikTokVideo(videoURL string) (*TikTokAwemeDetail, error) {
	endpoint := fmt.Sprintf("%s/v2/tiktok/video", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("url", videoURL)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators tiktok request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result TikTokVideoResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode tiktok response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators tiktok API returned success=false")
	}

	return &result.AwemeDetail, nil
}

func (c *client) GetTikTokProfileVideos(username string) (*TikTokProfileVideosResponse, error) {
	endpoint := fmt.Sprintf("%s/v2/tiktok/user/videos", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("handle", username)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators tiktok profile videos request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result TikTokProfileVideosResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode tiktok profile videos response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators tiktok profile videos API returned success=false")
	}

	return &result, nil
}

func (c *client) GetYouTubeVideo(videoURL string) (*YouTubeVideoResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/youtube/video", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("url", videoURL)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators youtube request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result YouTubeVideoResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode youtube response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators youtube API returned success=false")
	}

	return &result, nil
}

func (c *client) GetTweet(tweetURL string) (*TwitterTweetResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/twitter/tweet", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("url", tweetURL)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators twitter request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result TwitterTweetResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode twitter response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators twitter API returned success=false")
	}

	return &result, nil
}

func (c *client) GetUserTweets(username string) (*TwitterUserTweetsResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/twitter/user/tweets", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("handle", username)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators twitter user tweets request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result TwitterUserTweetsResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode twitter response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators twitter API returned success=false")
	}

	return &result, nil
}

func (c *client) GetInstagramProfile(username string) (*InstagramProfileResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/instagram/user/info", c.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	q := u.Query()
	q.Set("handle", username)
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"scrapecreators profile request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var result InstagramProfileResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode profile response")
	}

	if !result.Success {
		return nil, fmt.Errorf("scrapecreators profile API returned success=false")
	}

	return &result, nil
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

	req.Header.Set("x-api-key", c.Token)

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

	req.Header.Set("x-api-key", c.Token)

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
