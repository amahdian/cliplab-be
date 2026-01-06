package rocksolid

// ReelData represents the top-level structure of the Instagram Reel JSON.
type ReelData struct {
	ID                          string                    `json:"id"`
	Shortcode                   string                    `json:"shortcode"`
	ThumbnailSrc                string                    `json:"thumbnail_src"`
	Dimensions                  Dimensions                `json:"dimensions"`
	GatingInfo                  interface{}               `json:"gating_info"`
	FactCheckOverallRating      interface{}               `json:"fact_check_overall_rating"`
	FactCheckInformation        interface{}               `json:"fact_check_information"`
	SensitivityFrictionInfo     interface{}               `json:"sensitivity_friction_info"`
	SharingFrictionInfo         SharingFrictionInfo       `json:"sharing_friction_info"`
	MediaOverlayInfo            interface{}               `json:"media_overlay_info"`
	MediaPreview                string                    `json:"media_preview"`
	DisplayURL                  string                    `json:"display_url"`
	DisplayResources            []DisplayResource         `json:"display_resources"`
	AccessibilityCaption        interface{}               `json:"accessibility_caption"`
	HasAudio                    bool                      `json:"has_audio"`
	VideoURL                    string                    `json:"video_url"`
	VideoViewCount              int                       `json:"video_view_count"`
	VideoPlayCount              int                       `json:"video_play_count"`
	EncodingStatus              interface{}               `json:"encoding_status"`
	IsPublished                 bool                      `json:"is_published"`
	ProductType                 string                    `json:"product_type"`
	Title                       string                    `json:"title"`
	VideoDuration               float64                   `json:"video_duration"`
	ClipsMusicAttributionInfo   ClipsMusicAttributionInfo `json:"clips_music_attribution_info"`
	IsVideo                     bool                      `json:"is_video"`
	UpcomingEvent               interface{}               `json:"upcoming_event"`
	EdgeMediaToTaggedUser       EdgeMediaToTaggedUser     `json:"edge_media_to_tagged_user"`
	Owner                       Owner                     `json:"owner"`
	EdgeMediaToCaption          EdgeMediaToCaption        `json:"edge_media_to_caption"`
	CanSeeInsightsAsBrand       bool                      `json:"can_see_insights_as_brand"`
	CaptionIsEdited             bool                      `json:"caption_is_edited"`
	HasRankedComments           bool                      `json:"has_ranked_comments"`
	LikeAndViewCountsDisabled   bool                      `json:"like_and_view_counts_disabled"`
	EdgeMediaToParentComment    EdgeMediaToParentComment  `json:"edge_media_to_parent_comment"`
	CommentsDisabled            bool                      `json:"comments_disabled"`
	CommentingDisabledForViewer bool                      `json:"commenting_disabled_for_viewer"`
	TakenAtTimestamp            int64                     `json:"taken_at_timestamp"`
	EdgeMediaPreviewLike        EdgeCount                 `json:"edge_media_preview_like"`
	EdgeMediaToSponsorUser      interface{}               `json:"edge_media_to_sponsor_user"`
	IsAffiliate                 bool                      `json:"is_affiliate"`
	IsPaidPartnership           bool                      `json:"is_paid_partnership"`
	Location                    interface{}               `json:"location"`
	NftAssetInfo                interface{}               `json:"nft_asset_info"`
	ViewerHasLiked              bool                      `json:"viewer_has_liked"`
	ViewerHasSaved              bool                      `json:"viewer_has_saved"`
	ViewerHasSavedToCollection  bool                      `json:"viewer_has_saved_to_collection"`
	ViewerInPhotoOfYou          bool                      `json:"viewer_in_photo_of_you"`
	ViewerCanReshare            bool                      `json:"viewer_can_reshare"`
	IsAd                        bool                      `json:"is_ad"`
	EdgeWebMediaToRelatedMedia  interface{}               `json:"edge_web_media_to_related_media"`
	CoauthorProducers           []CoauthorProducer        `json:"coauthor_producers"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type SharingFrictionInfo struct {
	ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
	BloksAppURL               interface{} `json:"bloks_app_url"`
}

type DisplayResource struct {
	Src          string `json:"src"`
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
}

type ClipsMusicAttributionInfo struct {
	ArtistName            string `json:"artist_name"`
	SongName              string `json:"song_name"`
	UsesOriginalAudio     bool   `json:"uses_original_audio"`
	ShouldMuteAudio       bool   `json:"should_mute_audio"`
	ShouldMuteAudioReason string `json:"should_mute_audio_reason"`
	AudioID               string `json:"audio_id"`
}

type EdgeMediaToTaggedUser struct {
	Edges []TaggedUserEdge `json:"edges"`
}

type TaggedUserEdge struct {
	Node TaggedUserNode `json:"node"`
}

type TaggedUserNode struct {
	User User    `json:"user"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	ID   string  `json:"id"`
}

type User struct {
	FullName           string `json:"full_name"`
	FollowedByViewer   bool   `json:"followed_by_viewer"`
	ID                 string `json:"id"`
	IsVerified         bool   `json:"is_verified"`
	ProfilePicURL      string `json:"profile_pic_url"`
	Username           string `json:"username"`
	BlockedByViewer    bool   `json:"blocked_by_viewer,omitempty"`
	RestrictedByViewer bool   `json:"restricted_by_viewer,omitempty"`
}

type Owner struct {
	ID                        string    `json:"id"`
	Username                  string    `json:"username"`
	IsVerified                bool      `json:"is_verified"`
	ProfilePicURL             string    `json:"profile_pic_url"`
	BlockedByViewer           bool      `json:"blocked_by_viewer"`
	RestrictedByViewer        bool      `json:"restricted_by_viewer"`
	FollowedByViewer          bool      `json:"followed_by_viewer"`
	FullName                  string    `json:"full_name"`
	HasBlockedViewer          bool      `json:"has_blocked_viewer"`
	IsEmbedsDisabled          bool      `json:"is_embeds_disabled"`
	IsPrivate                 bool      `json:"is_private"`
	IsUnpublished             bool      `json:"is_unpublished"`
	RequestedByViewer         bool      `json:"requested_by_viewer"`
	PassTieringRecommendation bool      `json:"pass_tiering_recommendation"`
	EdgeOwnerToTimelineMedia  EdgeCount `json:"edge_owner_to_timeline_media"`
	EdgeFollowedBy            EdgeCount `json:"edge_followed_by"`
}

type EdgeCount struct {
	Count int64 `json:"count"`
}

type EdgeMediaToCaption struct {
	Edges []CaptionEdge `json:"edges"`
}

type CaptionEdge struct {
	Node CaptionNode `json:"node"`
}

type CaptionNode struct {
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
	ID        string `json:"id"`
}

type EdgeMediaToParentComment struct {
	Count    int64         `json:"count"`
	PageInfo PageInfo      `json:"page_info"`
	Edges    []CommentEdge `json:"edges"`
}

type CommentEdge struct {
	Node CommentNode `json:"node"`
}

type CommentNode struct {
	ID                   string               `json:"id"`
	Text                 string               `json:"text"`
	CreatedAt            int64                `json:"created_at"`
	DidReportAsSpam      bool                 `json:"did_report_as_spam"`
	Owner                User                 `json:"owner"`
	ViewerHasLiked       bool                 `json:"viewer_has_liked"`
	EdgeLikedBy          EdgeCount            `json:"edge_liked_by"`
	IsRestrictedPending  bool                 `json:"is_restricted_pending"`
	EdgeThreadedComments EdgeThreadedComments `json:"edge_threaded_comments,omitempty"`
}

type EdgeThreadedComments struct {
	Count    int           `json:"count"`
	PageInfo PageInfo      `json:"page_info"`
	Edges    []CommentEdge `json:"edges"`
}

type PageInfo struct {
	HasNextPage bool    `json:"has_next_page"`
	EndCursor   *string `json:"end_cursor"`
}

type CoauthorProducer struct {
	ID            string `json:"id"`
	IsVerified    bool   `json:"is_verified"`
	ProfilePicURL string `json:"profile_pic_url"`
	Username      string `json:"username"`
}

type Reels struct {
	Reels           []Reel `json:"reels"`
	PaginationToken string `json:"pagination_token"`
}

type Reel struct {
	Node   ReelNode `json:"node"`
	Cursor string   `json:"cursor"`
}

type ReelNode struct {
	Media    Media  `json:"media"`
	Typename string `json:"__typename"`
}

type Media struct {
	Pk                         string         `json:"pk"`
	ID                         string         `json:"id"`
	Code                       string         `json:"code"` // The shortcode
	MediaOverlayInfo           interface{}    `json:"media_overlay_info"`
	BoostedStatus              interface{}    `json:"boosted_status"`
	BoostUnavailableIdentifier interface{}    `json:"boost_unavailable_identifier"`
	BoostUnavailableReason     interface{}    `json:"boost_unavailable_reason"`
	User                       UserReference  `json:"user"`
	ProductType                string         `json:"product_type"`
	PlayCount                  int64          `json:"play_count"`
	ViewCount                  int64          `json:"view_count"`
	LikeAndViewCountsDisabled  bool           `json:"like_and_view_counts_disabled"`
	CommentCount               int64          `json:"comment_count"`
	LikeCount                  int64          `json:"like_count"`
	Audience                   interface{}    `json:"audience"`
	ClipsTabPinnedUserIds      []string       `json:"clips_tab_pinned_user_ids"`
	HasViewsFetching           bool           `json:"has_views_fetching"`
	MediaType                  int            `json:"media_type"`
	CarouselMedia              interface{}    `json:"carousel_media"`
	ImageVersions2             ImageVersions2 `json:"image_versions2"`
	Preview                    *string        `json:"preview"`
	OriginalHeight             int            `json:"original_height"`
	OriginalWidth              int            `json:"original_width"`
}

type UserReference struct {
	Pk string `json:"pk"`
	ID string `json:"id"`
}

type ImageVersions2 struct {
	Candidates []ImageCandidate `json:"candidates"`
}

type ImageCandidate struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
