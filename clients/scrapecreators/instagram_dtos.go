package scrapecreators

// Response represents the top-level structure of the ScrapeCreators Instagram Post response.
type Response struct {
	Success          bool   `json:"success"`
	CreditsRemaining int    `json:"credits_remaining"`
	Data             Data   `json:"data"`
	Status           string `json:"status"`
}

type Data struct {
	XdtShortcodeMedia ReelData `json:"xdt_shortcode_media"`
}

// ReelData represents the Instagram Reel data.
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
	MediaPreview                *string                   `json:"media_preview"`
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
	BlockedByViewer    *bool  `json:"blocked_by_viewer,omitempty"`
	RestrictedByViewer *bool  `json:"restricted_by_viewer,omitempty"`
}

type Owner struct {
	ID                        string      `json:"id"`
	Username                  string      `json:"username"`
	IsVerified                bool        `json:"is_verified"`
	ProfilePicURL             string      `json:"profile_pic_url"`
	BlockedByViewer           bool        `json:"blocked_by_viewer"`
	RestrictedByViewer        interface{} `json:"restricted_by_viewer"`
	FollowedByViewer          bool        `json:"followed_by_viewer"`
	FullName                  string      `json:"full_name"`
	HasBlockedViewer          bool        `json:"has_blocked_viewer"`
	IsEmbedsDisabled          bool        `json:"is_embeds_disabled"`
	IsPrivate                 bool        `json:"is_private"`
	IsUnpublished             bool        `json:"is_unpublished"`
	RequestedByViewer         bool        `json:"requested_by_viewer"`
	PassTieringRecommendation bool        `json:"pass_tiering_recommendation"`
	EdgeOwnerToTimelineMedia  EdgeCount   `json:"edge_owner_to_timeline_media"`
	EdgeFollowedBy            EdgeCount   `json:"edge_followed_by"`
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

// --- Reels Response DTOs ---

// ReelsResponse represents the response from the Instagram User Reels endpoint.
type ReelsResponse struct {
	Success          bool       `json:"success"`
	CreditsRemaining int        `json:"credits_remaining"`
	Items            []ReelItem `json:"items"`
}

type ReelItem struct {
	Media ReelMedia `json:"media"`
}

type ReelMedia struct {
	StrongID                                             string                   `json:"strong_id__"`
	Fbid                                                 int64                    `json:"fbid"`
	DeletedReason                                        int                      `json:"deleted_reason"`
	IsUnifiedVideo                                       bool                     `json:"is_unified_video"`
	ClientCacheKey                                       string                   `json:"client_cache_key"`
	CollaboratorEditEligibility                          bool                     `json:"collaborator_edit_eligibility"`
	CommentThreadingEnabled                              bool                     `json:"comment_threading_enabled"`
	CommercialityStatus                                  string                   `json:"commerciality_status"`
	IntegrityReviewDecision                              string                   `json:"integrity_review_decision"`
	IsReshareOfTextPostAppMediaInIg                      bool                     `json:"is_reshare_of_text_post_app_media_in_ig"`
	IsVisualReplyCommenterNoticeEnabled                  bool                     `json:"is_visual_reply_commenter_notice_enabled"`
	ShareCountDisabled                                   bool                     `json:"share_count_disabled"`
	Pk                                                   int64                    `json:"pk"`
	ID                                                   string                   `json:"id"`
	HasDelayedMetadata                                   bool                     `json:"has_delayed_metadata"`
	MezqlToken                                           string                   `json:"mezql_token"`
	ShouldRequestAds                                     bool                     `json:"should_request_ads"`
	HasPrivatelyLiked                                    bool                     `json:"has_privately_liked"`
	IsQuietPost                                          bool                     `json:"is_quiet_post"`
	SubtypeNameForREST                                   string                   `json:"subtype_name_for_REST__"`
	PlayCount                                            int64                    `json:"play_count"`
	IgPlayCount                                          int64                    `json:"ig_play_count"`
	AreRemixesCrosspostable                              bool                     `json:"are_remixes_crosspostable"`
	IsThirdPartyDownloadsEligible                        bool                     `json:"is_third_party_downloads_eligible"`
	HasAudio                                             bool                     `json:"has_audio"`
	VideoDuration                                        float64                  `json:"video_duration"`
	IsDashEligible                                       int                      `json:"is_dash_eligible"`
	ImageVersions2                                       ImageVersions2           `json:"image_versions2"`
	IgMediaSharingDisabled                               bool                     `json:"ig_media_sharing_disabled"`
	MediaCroppingInfo                                    MediaCroppingInfo        `json:"media_cropping_info"`
	MediaType                                            int                      `json:"media_type"`
	OriginalWidth                                        int                      `json:"original_width"`
	OriginalHeight                                       int                      `json:"original_height"`
	OrganicTrackingToken                                 string                   `json:"organic_tracking_token"`
	Caption                                              *ReelCaption             `json:"caption"`
	CoauthorProducers                                    []interface{}            `json:"coauthor_producers"`
	HasTaggedUsers                                       bool                     `json:"has_tagged_users"`
	ClipsMetadata                                        ClipsMetadata            `json:"clips_metadata"`
	VideoVersions                                        []VideoVersion           `json:"video_versions"`
	VideoDashManifest                                    string                   `json:"video_dash_manifest"`
	NumberOfQualities                                    int                      `json:"number_of_qualities"`
	LikeCount                                            int64                    `json:"like_count"`
	CommentCount                                         int64                    `json:"comment_count"`
	TakenAt                                              int64                    `json:"taken_at"`
	UserTags                                             *ReelUserTags            `json:"usertags"`
	PhotoOfYou                                           bool                     `json:"photo_of_you"`
	CanSeeInsightsAsBrand                                bool                     `json:"can_see_insights_as_brand"`
	DisplayUri                                           string                   `json:"display_uri"`
	IsInProfileGrid                                      bool                     `json:"is_in_profile_grid"`
	User                                                 ReelUser                 `json:"user"`
	Owner                                                ReelUser                 `json:"owner"`
	ProductType                                          string                   `json:"product_type"`
	IsPaidPartnership                                    bool                     `json:"is_paid_partnership"`
	LoggingInfoToken                                     string                   `json:"logging_info_token"`
	SharingFrictionInfo                                  SharingFrictionInfo      `json:"sharing_friction_info"`
	CanViewerReshare                                     bool                     `json:"can_viewer_reshare"`
	CanViewerSave                                        bool                     `json:"can_viewer_save"`
	CaptionIsEdited                                      bool                     `json:"caption_is_edited"`
	Code                                                 string                   `json:"code"`
	DeviceTimestamp                                      int64                    `json:"device_timestamp"`
	Url                                                  string                   `json:"url"`
	FilterType                                           int                      `json:"filter_type"`
	HasSharedToFb                                        int                      `json:"has_shared_to_fb"`
	IsPostLiveClipsMedia                                 bool                     `json:"is_post_live_clips_media"`
	ProfileGridThumbnailFittingStyle                     string                   `json:"profile_grid_thumbnail_fitting_style"`
	HasHighRiskGenAiInformTreatment                      bool                     `json:"has_high_risk_gen_ai_inform_treatment"`
	EnableMediaNotesProduction                           bool                     `json:"enable_media_notes_production"`
	HasLiked                                             bool                     `json:"has_liked"`
	HasViewsFetching                                     bool                     `json:"has_views_fetching"`
	IsCommentsGifComposerEnabled                         bool                     `json:"is_comments_gif_composer_enabled"`
	IsOrganicProductTaggingEligible                      bool                     `json:"is_organic_product_tagging_eligible"`
	MediaReposterBottomsheetEnabled                      bool                     `json:"media_reposter_bottomsheet_enabled"`
	ReportInfo                                           ReportInfo               `json:"report_info"`
	CommentLikesEnabled                                  bool                     `json:"comment_likes_enabled"`
	ProductSuggestions                                   []interface{}            `json:"product_suggestions"`
	CommentInformTreatment                               InformTreatment          `json:"comment_inform_treatment"`
	IgbioProduct                                         interface{}              `json:"igbio_product"`
	IsOpenToPublicSubmission                             bool                     `json:"is_open_to_public_submission"`
	IsSocialUfiDisabled                                  bool                     `json:"is_social_ufi_disabled"`
	OpenCarouselShowFollowButton                         bool                     `json:"open_carousel_show_follow_button"`
	TimelinePinnedUserIds                                interface{}              `json:"timeline_pinned_user_ids"`
	FbUserTags                                           interface{}              `json:"fb_user_tags"`
	CoauthorProducerCanSeeOrganicInsights                bool                     `json:"coauthor_producer_can_see_organic_insights"`
	InvitedCoauthorProducers                             []interface{}            `json:"invited_coauthor_producers"`
	ProfileGridControlEnabled                            bool                     `json:"profile_grid_control_enabled"`
	IsArtistPick                                         bool                     `json:"is_artist_pick"`
	BoostUnavailableIdentifier                           *string                  `json:"boost_unavailable_identifier"`
	BoostUnavailableReason                               *string                  `json:"boost_unavailable_reason"`
	BoostUnavailableReasonV2                             []interface{}            `json:"boost_unavailable_reason_v2"`
	SubscribeCtaVisible                                  bool                     `json:"subscribe_cta_visible"`
	CreativeConfig                                       interface{}              `json:"creative_config"`
	IsCutoutStickerAllowed                               bool                     `json:"is_cutout_sticker_allowed"`
	CutoutStickerInfo                                    []interface{}            `json:"cutout_sticker_info"`
	IsTaggedMediaSharedToViewerProfileGrid               bool                     `json:"is_tagged_media_shared_to_viewer_profile_grid"`
	ShouldShowAuthorPogForTaggedMediaSharedToProfileGrid bool                     `json:"should_show_author_pog_for_tagged_media_shared_to_profile_grid"`
	MetaAiSuggestedPrompts                               []interface{}            `json:"meta_ai_suggested_prompts"`
	GenAiChatWithAiCtaInfo                               interface{}              `json:"gen_ai_chat_with_ai_cta_info"`
	CanReply                                             bool                     `json:"can_reply"`
	IsEligibleContentForPostRollAd                       bool                     `json:"is_eligible_content_for_post_roll_ad"`
	IsPhotoCommentsComposerEnabledForAuthor              bool                     `json:"is_photo_comments_composer_enabled_for_author"`
	HideViewAllCommentEntrypoint                         bool                     `json:"hide_view_all_comment_entrypoint"`
	EligibleInsightsEntrypoints                          string                   `json:"eligible_insights_entrypoints"`
	HasMoreComments                                      bool                     `json:"has_more_comments"`
	MaxNumVisiblePreviewComments                         int                      `json:"max_num_visible_preview_comments"`
	ExploreHideComments                                  bool                     `json:"explore_hide_comments"`
	HiddenLikesStringVariant                             int                      `json:"hidden_likes_string_variant"`
	VideoStickerLocales                                  []string                 `json:"video_sticker_locales"`
	IsEligibleForAutodub                                 bool                     `json:"is_eligible_for_autodub"`
	IsEligibleForAutodubUpsell                           bool                     `json:"is_eligible_for_autodub_upsell"`
	MediaAttributionsData                                []interface{}            `json:"media_attributions_data"`
	CreatorProductLinkInfos                              []interface{}            `json:"creator_product_link_infos"`
	IsEligibleForPoe                                     bool                     `json:"is_eligible_for_poe"`
	IsEligibleForOrganicEagerRefresh                     bool                     `json:"is_eligible_for_organic_eager_refresh"`
	ShopRoutingUserId                                    *string                  `json:"shop_routing_user_id"`
	FeaturedProducts                                     []interface{}            `json:"featured_products"`
	IsReuseAllowed                                       bool                     `json:"is_reuse_allowed"`
	GenAiDetectionMethod                                 *GenAiDetectionMethod    `json:"gen_ai_detection_method"`
	WearableAttributionInfo                              *WearableAttributionInfo `json:"wearable_attribution_info"`
}

type GenAiDetectionMethod struct {
	DetectionMethod string `json:"detection_method"`
}

type ReportInfo struct {
	HasViewerSubmittedReport bool `json:"has_viewer_submitted_report"`
}

type InformTreatment struct {
	ActionType                interface{} `json:"action_type"`
	ShouldHaveInformTreatment bool        `json:"should_have_inform_treatment"`
	Text                      string      `json:"text"`
	Url                       interface{} `json:"url"`
}

type ImageVersions2 struct {
	AdditionalCandidates AdditionalCandidates `json:"additional_candidates"`
	Candidates           []Candidate          `json:"candidates"`
}

type AdditionalCandidates struct {
	FirstFrame     Candidate `json:"first_frame"`
	IgtvFirstFrame Candidate `json:"igtv_first_frame"`
}

type Candidate struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type MediaCroppingInfo struct {
	FourByThreeCrop Crop `json:"four_by_three_crop"`
}

type Crop struct {
	CropLeft   float64 `json:"crop_left"`
	CropRight  float64 `json:"crop_right"`
	CropTop    float64 `json:"crop_top"`
	CropBottom float64 `json:"crop_bottom"`
}

type ReelCaption struct {
	StrongID     string   `json:"strong_id__"`
	CreatedAt    int64    `json:"created_at"`
	CreatedAtUtc int64    `json:"created_at_utc"`
	Pk           string   `json:"pk"`
	MediaID      string   `json:"media_id"`
	Text         string   `json:"text"`
	User         ReelUser `json:"user"`
}

type ReelUser struct {
	StrongID                       string `json:"strong_id__"`
	Pk                             string `json:"pk"`
	PkID                           string `json:"pk_id"`
	ID                             int64  `json:"id"`
	FullName                       string `json:"full_name"`
	Username                       string `json:"username"`
	IsPrivate                      bool   `json:"is_private"`
	IsVerified                     bool   `json:"is_verified"`
	ProfilePicId                   string `json:"profile_pic_id"`
	ProfilePicUrl                  string `json:"profile_pic_url"`
	FbidV2                         string `json:"fbid_v2"`
	FeedPostReshareDisabled        bool   `json:"feed_post_reshare_disabled"`
	IsUnpublished                  bool   `json:"is_unpublished"`
	ThirdPartyDownloadsEnabled     int    `json:"third_party_downloads_enabled"`
	CanSeeQuietPostAttribution     bool   `json:"can_see_quiet_post_attribution"`
	AccountType                    int    `json:"account_type"`
	ShowAccountTransparencyDetails bool   `json:"show_account_transparency_details"`
	TransparencyProductEnabled     bool   `json:"transparency_product_enabled"`
	IsActiveOnTextPostApp          bool   `json:"is_active_on_text_post_app"`
}

type ClipsMetadata struct {
	ClipsCreationEntryPoint string      `json:"clips_creation_entry_point"`
	AudioType               string      `json:"audio_type"`
	MusicInfo               interface{} `json:"music_info"`
	OriginalSoundInfo       interface{} `json:"original_sound_info"`
	MashupInfo              MashupInfo  `json:"mashup_info"`
}

type MashupInfo struct {
	CanToggleMashupsAllowed bool `json:"can_toggle_mashups_allowed"`
	HasBeenMashedUp         bool `json:"has_been_mashed_up"`
	IsLightWeightCheck      bool `json:"is_light_weight_check"`
	MashupsAllowed          bool `json:"mashups_allowed"`
}

type VideoVersion struct {
	Type   int    `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}

type ReelUserTags struct {
	In []ReelUserTag `json:"in"`
}

type ReelUserTag struct {
	Position []float64 `json:"position"`
	User     ReelUser  `json:"user"`
}

// --- Page Posts Response DTOs ---

type PostsResponse struct {
	Success                   bool        `json:"success"`
	CreditsRemaining          int         `json:"credits_remaining"`
	NextMaxID                 string      `json:"next_max_id"`
	User                      PostUser    `json:"user"`
	AutoLoadMoreEnabled       bool        `json:"auto_load_more_enabled"`
	Status                    string      `json:"status"`
	ProfileGridItems          interface{} `json:"profile_grid_items"`
	ProfileGridItemsCursor    interface{} `json:"profile_grid_items_cursor"`
	PinnedProfileGridItemsIDs interface{} `json:"pinned_profile_grid_items_ids"`
	SpecialEmptyState         interface{} `json:"special_empty_state"`
	NumResults                int         `json:"num_results"`
	MoreAvailable             bool        `json:"more_available"`
	Items                     []ReelMedia `json:"items"`
}

type PostUser struct {
	StrongID               string `json:"strong_id__"`
	Pk                     string `json:"pk"`
	PkID                   string `json:"pk_id"`
	FullName               string `json:"full_name"`
	ProfileGridDisplayType string `json:"profile_grid_display_type"`
	ID                     string `json:"id"`
	Username               string `json:"username"`
	IsPrivate              bool   `json:"is_private"`
	IsVerified             bool   `json:"is_verified"`
	ProfilePicId           string `json:"profile_pic_id"`
	ProfilePicUrl          string `json:"profile_pic_url"`
	IsActiveOnTextPostApp  bool   `json:"is_active_on_text_post_app"`
}

type WearableAttributionInfo struct {
	AttributionCtaActionUrl   string    `json:"attribution_cta_action_url"`
	AttributionCtaText        string    `json:"attribution_cta_text"`
	AttributionIconUrl        string    `json:"attribution_icon_url"`
	AttributionSubtitle       string    `json:"attribution_subtitle"`
	AttributionTitle          string    `json:"attribution_title"`
	AttributionTopIconUrl     string    `json:"attribution_top_icon_url"`
	AttributionType           string    `json:"attribution_type"`
	IsWearableMediaProducer   bool      `json:"is_wearable_media_producer"`
	PivotPageCtaLabel         string    `json:"pivot_page_cta_label"`
	PivotPageCtaUrl           string    `json:"pivot_page_cta_url"`
	PivotPageDescription      string    `json:"pivot_page_description"`
	PivotPageImageUrl         string    `json:"pivot_page_image_url"`
	PivotPageTitle            string    `json:"pivot_page_title"`
	ReelsPillAttributionTitle string    `json:"reels_pill_attribution_title"`
	PivotPageMicroUserDict    *PostUser `json:"pivot_page_micro_user_dict"`
}
