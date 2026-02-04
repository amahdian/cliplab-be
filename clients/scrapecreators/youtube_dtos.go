package scrapecreators

import "time"

// YouTubeVideoResponse represents the response from the YouTube video detail endpoint.
type YouTubeVideoResponse struct {
	Success            bool               `json:"success"`
	CreditsRemaining   int                `json:"credits_remaining"`
	ID                 string             `json:"id"`
	Thumbnail          string             `json:"thumbnail"`
	URL                string             `json:"url"`
	PublishDate        time.Time          `json:"publishDate"`
	Type               string             `json:"type"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	CommentCountText   string             `json:"commentCountText"`
	CommentCountInt    int                `json:"commentCountInt"`
	LikeCountText      string             `json:"likeCountText"`
	LikeCountInt       *int               `json:"likeCountInt"`
	ViewCountText      string             `json:"viewCountText"`
	ViewCountInt       int64              `json:"viewCountInt"`
	PublishDateText    string             `json:"publishDateText"`
	Collaborators      []interface{}      `json:"collaborators"`
	Channel            YouTubeChannel     `json:"channel"`
	Chapters           []interface{}      `json:"chapters"`
	WatchNextVideos    []YoutubeWatchNext `json:"watchNextVideos"`
	Keywords           []string           `json:"keywords"`
	Genre              string             `json:"genre"`
	DurationMs         int                `json:"durationMs"`
	DurationFormatted  string             `json:"durationFormatted"`
	CaptionTracks      []interface{}      `json:"captionTracks"`
	Transcript         interface{}        `json:"transcript"`
	TranscriptOnlyText interface{}        `json:"transcript_only_text"`
}

type YouTubeChannel struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Handle string `json:"handle"`
	Title  string `json:"title"`
}

type YoutubeWatchNext struct {
	ID                string         `json:"id"`
	Title             string         `json:"title"`
	Thumbnail         string         `json:"thumbnail"`
	Channel           YouTubeChannel `json:"channel"`
	PublishedTimeText string         `json:"publishedTimeText"`
	PublishedTime     time.Time      `json:"publishedTime"`
	PublishDateText   string         `json:"publishDateText"`
	PublishDate       time.Time      `json:"publishDate"`
	ViewCountText     string         `json:"viewCountText"`
	ViewCountInt      int64          `json:"viewCountInt"`
	LengthText        string         `json:"lengthText"`
	LengthInSeconds   int            `json:"lengthInSeconds"`
	VideoURL          string         `json:"videoUrl"`
}
