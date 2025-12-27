package dtos

import "time"

type InstagramContentType string

const (
	ContentTypeVideo InstagramContentType = "Video"
	ContentTypePhoto InstagramContentType = "Photo"
)

type InstagramPost struct {
	URL                 string                       `json:"url"`
	UserPosted          string                       `json:"user_posted"`
	Description         string                       `json:"description"`
	NumComments         int                          `json:"num_comments"`
	DatePosted          time.Time                    `json:"date_posted"`
	Likes               int                          `json:"likes"`
	Photos              []string                     `json:"photos"`
	Videos              []string                     `json:"videos"`
	LatestComments      []*InstgarmLatestComment     `json:"latest_comments"`
	PostID              string                       `json:"post_id"`
	Shortcode           string                       `json:"shortcode"`
	ContentType         string                       `json:"content_type"`
	Pk                  string                       `json:"pk"`
	ContentID           string                       `json:"content_id"`
	EngagementScoreView int                          `json:"engagement_score_view"`
	Thumbnail           string                       `json:"thumbnail"`
	ProductType         string                       `json:"product_type"`
	VideoPlayCount      int                          `json:"video_play_count"`
	Followers           int                          `json:"followers"`
	PostsCount          int                          `json:"posts_count"`
	ProfileImageLink    string                       `json:"profile_image_link"`
	IsVerified          bool                         `json:"is_verified"`
	IsPaidPartnership   bool                         `json:"is_paid_partnership"`
	PartnershipDetails  *InstagramPartnershipDetails `json:"partnership_details"`
	UserPostedID        string                       `json:"user_posted_id"`
	PostContent         []*InstagramPostContent      `json:"post_content"`
	Audio               InstagramAudio               `json:"audio"`
	ProfileURL          string                       `json:"profile_url"`
	VideosDuration      []*InstagramVideoDuration    `json:"videos_duration"`
	Images              []interface{}                `json:"images"`   // Using interface{} as it's an empty array in the sample
	AltText             interface{}                  `json:"alt_text"` // Using interface{} as it's null in the sample
	PhotosNumber        int                          `json:"photos_number"`
	Timestamp           time.Time                    `json:"timestamp"`
	Input               *InstagramInput              `json:"input"`
}

type InstgarmLatestComment struct {
	Comments       string `json:"comments"`
	UserCommenting string `json:"user_commenting"`
	Likes          int    `json:"likes"`
	ProfilePicture string `json:"profile_picture"`
}

type InstagramPartnershipDetails struct {
	ProfileID  interface{} `json:"profile_id"`  // Using interface{} as it's null in the sample
	Username   interface{} `json:"username"`    // Using interface{} as it's null in the sample
	ProfileURL interface{} `json:"profile_url"` // Using interface{} as it's null in the sample
}

type InstagramPostContent struct {
	Index int                  `json:"index"`
	Type  InstagramContentType `json:"type"`
	URL   string               `json:"url"`
	ID    string               `json:"id"`
}

type InstagramAudio struct {
	AudioAssetID       string `json:"audio_asset_id"`
	OriginalAudioTitle string `json:"original_audio_title"`
	IGArtistUsername   string `json:"ig_artist_username"`
	IGArtistID         string `json:"ig_artist_id"`
}

type InstagramVideoDuration struct {
	URL           string  `json:"url"`
	VideoDuration float64 `json:"video_duration"`
}

type InstagramInput struct {
	URL string `json:"url"`
}
