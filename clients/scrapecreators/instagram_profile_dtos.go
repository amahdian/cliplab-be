package scrapecreators

type InstagramProfileResponse struct {
	Success          bool                 `json:"success"`
	CreditsRemaining int                  `json:"credits_remaining"`
	Data             InstagramProfileData `json:"data"`
}

type InstagramProfileData struct {
	User InstagramProfileUser `json:"user"`
}

type InstagramProfileUser struct {
	AiAgentOwnerUsername           interface{}           `json:"ai_agent_owner_username"`
	AiAgentType                    interface{}           `json:"ai_agent_type"`
	Biography                      string                `json:"biography"`
	BioLinks                       []BioLink             `json:"bio_links"`
	FbProfileBiolink               interface{}           `json:"fb_profile_biolink"`
	BiographyWithEntities          BiographyWithEntities `json:"biography_with_entities"`
	BlockedByViewer                bool                  `json:"blocked_by_viewer"`
	RestrictedByViewer             interface{}           `json:"restricted_by_viewer"`
	CountryBlock                   bool                  `json:"country_block"`
	EimuId                         string                `json:"eimu_id"`
	ExternalUrl                    string                `json:"external_url"`
	ExternalUrlLinkshimmed         string                `json:"external_url_linkshimmed"`
	EdgeFollowedBy                 EdgeCount             `json:"edge_followed_by"`
	Fbid                           string                `json:"fbid"`
	FollowedByViewer               bool                  `json:"followed_by_viewer"`
	EdgeFollow                     EdgeCount             `json:"edge_follow"`
	FollowsViewer                  bool                  `json:"follows_viewer"`
	FullName                       string                `json:"full_name"`
	GroupMetadata                  interface{}           `json:"group_metadata"`
	HasArEffects                   bool                  `json:"has_ar_effects"`
	HasClips                       bool                  `json:"has_clips"`
	HasGuides                      bool                  `json:"has_guides"`
	HasChannel                     bool                  `json:"has_channel"`
	HasBlockedViewer               bool                  `json:"has_blocked_viewer"`
	HighlightReelCount             int                   `json:"highlight_reel_count"`
	HasOnboardedToTextPostApp      bool                  `json:"has_onboarded_to_text_post_app"`
	HasRequestedViewer             bool                  `json:"has_requested_viewer"`
	HideLikeAndViewCounts          bool                  `json:"hide_like_and_view_counts"`
	ID                             string                `json:"id"`
	IsBusinessAccount              bool                  `json:"is_business_account"`
	IsProfessionalAccount          bool                  `json:"is_professional_account"`
	IsSupervisionEnabled           bool                  `json:"is_supervision_enabled"`
	IsGuardianOfViewer             bool                  `json:"is_guardian_of_viewer"`
	IsSupervisedByViewer           bool                  `json:"is_supervised_by_viewer"`
	IsSupervisedUser               bool                  `json:"is_supervised_user"`
	IsEmbedsDisabled               bool                  `json:"is_embeds_disabled"`
	IsJoinedRecently               bool                  `json:"is_joined_recently"`
	GuardianId                     interface{}           `json:"guardian_id"`
	BusinessAddressJson            interface{}           `json:"business_address_json"`
	BusinessContactMethod          string                `json:"business_contact_method"`
	BusinessEmail                  *string               `json:"business_email"`
	BusinessPhoneNumber            *string               `json:"business_phone_number"`
	BusinessCategoryName           *string               `json:"business_category_name"`
	OverallCategoryName            *string               `json:"overall_category_name"`
	CategoryEnum                   *string               `json:"category_enum"`
	CategoryName                   *string               `json:"category_name"`
	IsPrivate                      bool                  `json:"is_private"`
	IsVerified                     bool                  `json:"is_verified"`
	IsVerifiedByMv4b               bool                  `json:"is_verified_by_mv4b"`
	IsRegulatedC18                 bool                  `json:"is_regulated_c18"`
	EdgeMutualFollowedBy           EdgeMutualFollowedBy  `json:"edge_mutual_followed_by"`
	PinnedChannelsListCount        int                   `json:"pinned_channels_list_count"`
	ProfilePicUrl                  string                `json:"profile_pic_url"`
	ProfilePicUrlHd                string                `json:"profile_pic_url_hd"`
	RequestedByViewer              bool                  `json:"requested_by_viewer"`
	ShouldShowCategory             bool                  `json:"should_show_category"`
	ShouldShowPublicContacts       bool                  `json:"should_show_public_contacts"`
	ShowAccountTransparencyDetails bool                  `json:"show_account_transparency_details"`
	ShowTextPostAppBadge           bool                  `json:"show_text_post_app_badge"`
	RemoveMessageEntrypoint        bool                  `json:"remove_message_entrypoint"`
	TransparencyLabel              interface{}           `json:"transparency_label"`
	TransparencyProduct            interface{}           `json:"transparency_product"`
	Username                       string                `json:"username"`
	Pronouns                       []string              `json:"pronouns"`
	EdgeFelixVideoTimeline         ProfileTimeline       `json:"edge_felix_video_timeline"`
	EdgeOwnerToTimelineMedia       ProfileTimeline       `json:"edge_owner_to_timeline_media"`
	EdgeSavedMedia                 ProfileTimeline       `json:"edge_saved_media"`
	EdgeMediaCollections           ProfileTimeline       `json:"edge_media_collections"`
}

type BioLink struct {
	Title    string `json:"title"`
	LynxUrl  string `json:"lynx_url"`
	Url      string `json:"url"`
	LinkType string `json:"link_type"`
}

type BiographyWithEntities struct {
	RawText  string        `json:"raw_text"`
	Entities []interface{} `json:"entities"`
}

type EdgeMutualFollowedBy struct {
	Count int64         `json:"count"`
	Edges []interface{} `json:"edges"`
}

type ProfileTimeline struct {
	Count    int64          `json:"count"`
	PageInfo PageInfo       `json:"page_info"`
	Edges    []TimelineEdge `json:"edges"`
}

type TimelineEdge struct {
	Node TimelineNode `json:"node"`
}

type TimelineNode struct {
	Typename                         string                `json:"__typename"`
	ID                               string                `json:"id"`
	Shortcode                        string                `json:"shortcode"`
	Dimensions                       Dimensions            `json:"dimensions"`
	DisplayURL                       string                `json:"display_url"`
	EdgeMediaToTaggedUser            EdgeMediaToTaggedUser `json:"edge_media_to_tagged_user"`
	FactCheckOverallRating           interface{}           `json:"fact_check_overall_rating"`
	FactCheckInformation             interface{}           `json:"fact_check_information"`
	GatingInfo                       interface{}           `json:"gating_info"`
	SharingFrictionInfo              SharingFrictionInfo   `json:"sharing_friction_info"`
	MediaOverlayInfo                 interface{}           `json:"media_overlay_info"`
	MediaPreview                     *string               `json:"media_preview"`
	Owner                            ProfileTimelineOwner  `json:"owner"`
	IsVideo                          bool                  `json:"is_video"`
	HasUpcomingEvent                 bool                  `json:"has_upcoming_event"`
	AccessibilityCaption             interface{}           `json:"accessibility_caption"`
	DashInfo                         *DashInfo             `json:"dash_info,omitempty"`
	HasAudio                         *bool                 `json:"has_audio,omitempty"`
	TrackingToken                    string                `json:"tracking_token"`
	VideoURL                         *string               `json:"video_url,omitempty"`
	VideoViewCount                   *int64                `json:"video_view_count,omitempty"`
	EdgeMediaToCaption               EdgeMediaToCaption    `json:"edge_media_to_caption"`
	EdgeMediaToComment               EdgeCount             `json:"edge_media_to_comment"`
	CommentsDisabled                 bool                  `json:"comments_disabled"`
	TakenAtTimestamp                 int64                 `json:"taken_at_timestamp"`
	EdgeLikedBy                      EdgeCount             `json:"edge_liked_by"`
	EdgeMediaPreviewLike             EdgeCount             `json:"edge_media_preview_like"`
	Location                         interface{}           `json:"location"`
	NftAssetInfo                     interface{}           `json:"nft_asset_info"`
	ThumbnailSrc                     string                `json:"thumbnail_src"`
	ThumbnailTallSrc                 *string               `json:"thumbnail_tall_src,omitempty"`
	ThumbnailResources               []ThumbnailResource   `json:"thumbnail_resources"`
	FelixProfileGridCrop             *Crop                 `json:"felix_profile_grid_crop,omitempty"`
	TallProfileGridCrop              *Crop                 `json:"tall_profile_grid_crop,omitempty"`
	ProfileGridThumbnailFittingStyle string                `json:"profile_grid_thumbnail_fitting_style"`
	CoauthorProducers                []CoauthorProducer    `json:"coauthor_producers"`
	PinnedForUsers                   []interface{}         `json:"pinned_for_users"`
	ViewerCanReshare                 bool                  `json:"viewer_can_reshare"`
	LikeAndViewCountsDisabled        bool                  `json:"like_and_view_counts_disabled"`
	EncodingStatus                   interface{}           `json:"encoding_status"`
	IsPublished                      bool                  `json:"is_published"`
	ProductType                      string                `json:"product_type"`
	Title                            string                `json:"title"`
	VideoDuration                    *float64              `json:"video_duration,omitempty"`
}

type ProfileTimelineOwner struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type DashInfo struct {
	IsDashEligible    bool   `json:"is_dash_eligible"`
	VideoDashManifest string `json:"video_dash_manifest"`
	NumberOfQualities int    `json:"number_of_qualities"`
}

type ThumbnailResource struct {
	Src          string `json:"src"`
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
}
