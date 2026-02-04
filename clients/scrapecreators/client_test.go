package scrapecreators

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetInstagramPost(t *testing.T) {
	// Try to load env files from the cliplab-be directory
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	assert.NotNil(t, c)

	// Using a real post for integration testing
	res, err := c.GetInstagramPost("https://www.instagram.com/reel/DR5GmXVjeEV/")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "DR5GmXVjeEV", res.Shortcode)
}

func TestGetInstagramPageReels(t *testing.T) {
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	res, err := c.GetInstagramPageReels("instagram")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.NotEmpty(t, res.Items)
}

func TestGetInstagramPagePosts(t *testing.T) {
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		token = os.Getenv("SCRAP_CREATORS_TOKEN")
	}

	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	res, err := c.GetInstagramPagePosts("instagram")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.NotEmpty(t, res.Items)
}

func TestGetTikTokVideo(t *testing.T) {
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	res, err := c.GetTikTokVideo("https://www.tiktok.com/@evanjrocha/video/7581126950450744598")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "7581126950450744598", res.AwemeID)
	assert.NotZero(t, res.Statistics.PlayCount)
}

func TestGetTikTokProfileVideos(t *testing.T) {
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	res, err := c.GetTikTokProfileVideos("evanjrocha")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.NotEmpty(t, res.AwemeList)
}

func TestGetYoutubeVideo(t *testing.T) {
	_ = godotenv.Load("../../.env.dev", "../../.env")

	token := os.Getenv("SCRAPE_CREATORS_TOKEN")
	if token == "" {
		t.Skip("SCRAPE_CREATORS_TOKEN not set")
	}

	c := NewClient(token)
	res, err := c.GetYouTubeVideo("https://www.youtube.com/shorts/IB5cfAszMZY")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.Equal(t, "IB5cfAszMZY", res.ID)
	assert.NotZero(t, res.ViewCountInt)
}
