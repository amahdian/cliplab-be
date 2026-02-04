package scrapecreators

// TikTokVideoResponse represents the response from the TikTok video detail endpoint.
type TikTokVideoResponse struct {
	Success          bool              `json:"success"`
	CreditsRemaining int               `json:"credits_remaining"`
	AwemeDetail      TikTokAwemeDetail `json:"aweme_detail"`
}

type TikTokProfileVideosResponse struct {
	Success          bool                `json:"success"`
	CreditsRemaining int                 `json:"credits_remaining"`
	MaxCursor        int64               `json:"max_cursor"`
	HasMore          int                 `json:"has_more"` // 1 or 0
	StatusCode       int                 `json:"status_code"`
	AwemeList        []TikTokAwemeDetail `json:"aweme_list"`
}

type TikTokAwemeDetail struct {
	AddedSoundMusicInfo            *TikTokMusic        `json:"added_sound_music_info"`
	AigcInfo                       TikTokAigcInfo      `json:"aigc_info"`
	AllowGift                      bool                `json:"allow_gift"`
	Anchors                        interface{}         `json:"anchors"`
	AnchorsExtras                  string              `json:"anchors_extras"`
	AnimatedImageInfo              TikTokAnimatedInfo  `json:"animated_image_info"`
	Author                         TikTokAuthor        `json:"author"`
	AuthorUserID                   int64               `json:"author_user_id"`
	AwemeAcl                       TikTokAwemeAcl      `json:"aweme_acl"`
	AwemeID                        string              `json:"aweme_id"`
	AwemeType                      int                 `json:"aweme_type"`
	Banners                        interface{}         `json:"banners"`
	BehindTheSongMusicIds          interface{}         `json:"behind_the_song_music_ids"`
	BehindTheSongVideoMusicIds     interface{}         `json:"behind_the_song_video_music_ids"`
	BodydanceScore                 int                 `json:"bodydance_score"`
	BrandedContentAccounts         interface{}         `json:"branded_content_accounts"`
	C2paInfo                       TikTokC2paInfo      `json:"c2pa_info"`
	CcTemplateInfo                 TikTokTemplateInfo  `json:"cc_template_info"`
	ChaList                        []TikTokChallenge   `json:"cha_list"`
	ChallengePosition              interface{}         `json:"challenge_position"`
	CmtSwt                         bool                `json:"cmt_swt"`
	CollectStat                    int                 `json:"collect_stat"`
	CommentConfig                  TikTokCommentConfig `json:"comment_config"`
	CommentTopbarInfo              interface{}         `json:"comment_topbar_info"`
	CommerceConfigData             interface{}         `json:"commerce_config_data"`
	CommerceInfo                   TikTokCommerceInfo  `json:"commerce_info"`
	ContentDesc                    string              `json:"content_desc"`
	ContentDescExtra               []interface{}       `json:"content_desc_extra"`
	ContentLevel                   int                 `json:"content_level"`
	ContentModel                   TikTokContentModel  `json:"content_model"`
	ContentOriginalType            int                 `json:"content_original_type"`
	ContentSizeType                int                 `json:"content_size_type"`
	ContentType                    string              `json:"content_type"`
	CoverLabels                    interface{}         `json:"cover_labels"`
	CreateTime                     int64               `json:"create_time"`
	CreationInfo                   TikTokCreationInfo  `json:"creation_info"`
	CreatorAiComment               TikTokAiComment     `json:"creator_ai_comment"`
	Desc                           string              `json:"desc"`
	DescLanguage                   string              `json:"desc_language"`
	DisableSearchTrendingBar       bool                `json:"disable_search_trending_bar"`
	Distance                       string              `json:"distance"`
	DistributeType                 int                 `json:"distribute_type"`
	EcosystemPerceptionEnhancement interface{}         `json:"ecosystem_perception_enhancement"`
	FollowUpPublishFromID          int                 `json:"follow_up_publish_from_id"`
	Geofencing                     interface{}         `json:"geofencing"`
	GeofencingRegions              interface{}         `json:"geofencing_regions"`
	GreenScreenMaterials           interface{}         `json:"green_screen_materials"`
	GroupID                        string              `json:"group_id"`
	GroupIDList                    interface{}         `json:"group_id_list"`
	HasDanmaku                     bool                `json:"has_danmaku"`
	HasPromoteEntry                int                 `json:"has_promote_entry"`
	HasVsEntry                     bool                `json:"has_vs_entry"`
	HaveDashboard                  bool                `json:"have_dashboard"`
	HybridLabel                    interface{}         `json:"hybrid_label"`
	ImageInfos                     interface{}         `json:"image_infos"`
	InteractPermission             TikTokPermissions   `json:"interact_permission"`
	InteractionStickers            []interface{}       `json:"interaction_stickers"`
	IsAds                          bool                `json:"is_ads"`
	IsDescriptionTranslatable      bool                `json:"is_description_translatable"`
	IsHashTag                      int                 `json:"is_hash_tag"`
	IsNffOrNr                      bool                `json:"is_nff_or_nr"`
	IsOnThisDay                    int                 `json:"is_on_this_day"`
	IsPaidContent                  bool                `json:"is_paid_content"`
	IsPgcshow                      bool                `json:"is_pgcshow"`
	IsPreview                      int                 `json:"is_preview"`
	IsRelieve                      bool                `json:"is_relieve"`
	IsTextStickerTranslatable      bool                `json:"is_text_sticker_translatable"`
	IsTitleTranslatable            bool                `json:"is_title_translatable"`
	IsTop                          int                 `json:"is_top"`
	IsVr                           bool                `json:"is_vr"`
	ItemCommentSettings            int                 `json:"item_comment_settings"`
	ItemDuet                       int                 `json:"item_duet"`
	ItemReact                      int                 `json:"item_react"`
	ItemStitch                     int                 `json:"item_stitch"`
	LabelTop                       TikTokImage         `json:"label_top"`
	LabelTopText                   interface{}         `json:"label_top_text"`
	LongVideo                      interface{}         `json:"long_video"`
	MainArchCommon                 string              `json:"main_arch_common"`
	MaskInfos                      []interface{}       `json:"mask_infos"`
	MemeRegInfo                    interface{}         `json:"meme_reg_info"`
	MiscInfo                       string              `json:"misc_info"`
	MufCommentInfoV2               interface{}         `json:"muf_comment_info_v2"`
	Music                          TikTokMusic         `json:"music"`
	MusicBeginTimeInMs             int                 `json:"music_begin_time_in_ms"`
	MusicEndTimeInMs               int                 `json:"music_end_time_in_ms"`
	MusicSelectedFrom              string              `json:"music_selected_from"`
	MusicTitleStyle                int                 `json:"music_title_style"`
	MusicVolume                    string              `json:"music_volume"`
	NeedTrimStep                   bool                `json:"need_trim_step"`
	NeedVsEntry                    bool                `json:"need_vs_entry"`
	NicknamePosition               interface{}         `json:"nickname_position"`
	NoSelectedMusic                bool                `json:"no_selected_music"`
	OperatorBoostInfo              interface{}         `json:"operator_boost_info"`
	OriginCommentIds               interface{}         `json:"origin_comment_ids"`
	OriginVolume                   string              `json:"origin_volume"`
	OriginalClientText             interface{}         `json:"original_client_text"`
	PaidContentInfo                interface{}         `json:"paid_content_info"`
	PickedUsers                    []interface{}       `json:"picked_users"`
	PlaylistBlocked                bool                `json:"playlist_blocked"`
	PoiReTagSignal                 int                 `json:"poi_re_tag_signal"`
	Position                       interface{}         `json:"position"`
	PreventDownload                bool                `json:"prevent_download"`
	ProductsInfo                   interface{}         `json:"products_info"`
	Promote                        TikTokPromote       `json:"promote"`
	PromoteCapcutToggle            int                 `json:"promote_capcut_toggle"`
	PromoteIconText                string              `json:"promote_icon_text"`
	PromoteToast                   string              `json:"promote_toast"`
	PromoteToastKey                string              `json:"promote_toast_key"`
	QuestionList                   interface{}         `json:"question_list"`
	QuickReplyEmojis               []string            `json:"quick_reply_emojis"`
	Rate                           int                 `json:"rate"`
	ReferenceTtsVoiceIds           interface{}         `json:"reference_tts_voice_ids"`
	ReferenceVoiceFilterIds        interface{}         `json:"reference_voice_filter_ids"`
	Region                         string              `json:"region"`
	RiskInfos                      TikTokRiskInfos     `json:"risk_infos"`
	SearchHighlight                interface{}         `json:"search_highlight"`
	ShareInfo                      TikTokShareInfo     `json:"share_info"`
	ShareURL                       string              `json:"share_url"`
	ShootTabName                   string              `json:"shoot_tab_name"`
	SmartSearchInfo                TikTokSearchInfo    `json:"smart_search_info"`
	SocialInteractionBlob          interface{}         `json:"social_interaction_blob"`
	SolariaProfile                 TikTokSolaria       `json:"solaria_profile"`
	SortLabel                      string              `json:"sort_label"`
	Statistics                     TikTokStatistics    `json:"statistics"`
	Status                         TikTokStatus        `json:"status"`
	SupportDanmaku                 bool                `json:"support_danmaku"`
	SurveyInfo                     interface{}         `json:"survey_info"`
	TakoBubbleInfo                 TikTokTakoInfo      `json:"tako_bubble_info"`
	TextExtra                      []interface{}       `json:"text_extra"`
	TextStickerMajorLang           string              `json:"text_sticker_major_lang"`
	TitleLanguage                  string              `json:"title_language"`
	TnsUeFeedInfo                  interface{}         `json:"tns_ue_feed_info"`
	TtecSuggestWords               interface{}         `json:"ttec_suggest_words"`
	TtsVoiceIds                    interface{}         `json:"tts_voice_ids"`
	TttProductRecallType           int                 `json:"ttt_product_recall_type"`
	UniqidPosition                 interface{}         `json:"uniqid_position"`
	UpvoteInfo                     TikTokUpvoteInfo    `json:"upvote_info"`
	UpvotePreload                  TikTokUpvotePreload `json:"upvote_preload"`
	UsedFullSong                   bool                `json:"used_full_song"`
	UserDigged                     int                 `json:"user_digged"`
	Video                          TikTokVideo         `json:"video"`
	VideoControl                   TikTokVideoControl  `json:"video_control"`
	VideoLabels                    []interface{}       `json:"video_labels"`
	VideoText                      []interface{}       `json:"video_text"`
	VisualSearchInfo               interface{}         `json:"visual_search_info"`
	VoiceFilterIds                 interface{}         `json:"voice_filter_ids"`
	WithPromotionalMusic           bool                `json:"with_promotional_music"`
	WithoutWatermark               bool                `json:"without_watermark"`
	IsAd                           bool                `json:"is_ad"`
	IsEligibleForCommission        bool                `json:"is_eligible_for_commission"`
	IsPaidPartnership              bool                `json:"is_paid_partnership"`
	CreateTimeUtc                  string              `json:"create_time_utc"`
	URL                            string              `json:"url"`
	ShopProductURL                 interface{}         `json:"shop_product_url"`
}

type TikTokAigcInfo struct {
	AigcLabelType int  `json:"aigc_label_type"`
	CreatedByAi   bool `json:"created_by_ai"`
}

type TikTokAnimatedInfo struct {
	Effect int `json:"effect"`
	Type   int `json:"type"`
}

type TikTokAuthor struct {
	UID                 string      `json:"uid"`
	Nickname            string      `json:"nickname"`
	UniqueID            string      `json:"unique_id"`
	Signature           string      `json:"signature"`
	AvatarThumb         TikTokImage `json:"avatar_thumb"`
	AvatarMedium        TikTokImage `json:"avatar_medium"`
	AvatarLarger        TikTokImage `json:"avatar_larger"`
	FollowerCount       int64       `json:"follower_count"`
	FollowingCount      int64       `json:"following_count"`
	TotalFavorited      int64       `json:"total_favorited"`
	AwemeCount          int         `json:"aweme_count"`
	VerificationType    int         `json:"verification_type"`
	SecUID              string      `json:"sec_uid"`
	IsStar              bool        `json:"is_star"`
	Region              string      `json:"region"`
	Language            string      `json:"language"`
	YoutubeChannelId    string      `json:"youtube_channel_id"`
	YoutubeChannelTitle string      `json:"youtube_channel_title"`
}

type TikTokAwemeAcl struct {
	DownloadGeneral TikTokAclDetail `json:"download_general"`
	DownloadMask    TikTokAclDetail `json:"download_mask_panel"`
	ShareGeneral    TikTokAclDetail `json:"share_general"`
}

type TikTokAclDetail struct {
	Code      int    `json:"code"`
	Extra     string `json:"extra"`
	Mute      bool   `json:"mute"`
	ShowType  int    `json:"show_type"`
	Transcode int    `json:"transcode"`
	ToastMsg  string `json:"toast_msg"`
}

type TikTokC2paInfo struct {
	AigcPercentageType int  `json:"aigc_percentage_type"`
	IsCapcut           bool `json:"is_capcut"`
	IsTiktok           bool `json:"is_tiktok"`
}

type TikTokTemplateInfo struct {
	AuthorName string `json:"author_name"`
	ClipCount  int    `json:"clip_count"`
	Desc       string `json:"desc"`
}

type TikTokChallenge struct {
	CID     string `json:"cid"`
	ChaName string `json:"cha_name"`
	Desc    string `json:"desc"`
	Schema  string `json:"schema"`
	Type    int    `json:"type"`
	SubUint int    `json:"sub_type"`
}

type TikTokCommentConfig struct {
	QuickComment TikTokQuickComment `json:"quick_comment"`
}

type TikTokQuickComment struct {
	Enabled  bool `json:"enabled"`
	RecLevel int  `json:"rec_level"`
}

type TikTokCommerceInfo struct {
	AdvPromotable   bool   `json:"adv_promotable"`
	OrganicLogExtra string `json:"organic_log_extra"`
}

type TikTokContentModel struct {
	CustomBiz   TikTokCustomBiz   `json:"custom_biz"`
	StandardBiz TikTokStandardBiz `json:"standard_biz"`
}

type TikTokCustomBiz struct {
	AwemeTrace string `json:"aweme_trace"`
}

type TikTokStandardBiz struct {
	CreatorAnalytics interface{} `json:"creator_analytics"`
}

type TikTokCreationInfo struct {
	CreationUsedFunctions []string `json:"creation_used_functions"`
}

type TikTokAiComment struct {
	EligibleVideo bool `json:"eligible_video"`
	HasAiTopic    bool `json:"has_ai_topic"`
}

type TikTokPermissions struct {
	Duet   int `json:"duet"`
	Upvote int `json:"upvote"`
	Stitch int `json:"stitch"`
}

type TikTokPromote struct {
	Extra string `json:"extra"`
}

type TikTokRiskInfos struct {
	RiskSink bool `json:"risk_sink"`
	Warn     bool `json:"warn"`
}

type TikTokShareInfo struct {
	ShareDesc  string `json:"share_desc"`
	ShareURL   string `json:"share_url"`
	ShareTitle string `json:"share_title"`
}

type TikTokSearchInfo struct {
	AskTakoIntentType int `json:"ask_tako_intent_type"`
}

type TikTokSolaria struct {
	Profile string `json:"profile"`
}

type TikTokStatistics struct {
	CollectCount  int64 `json:"collect_count"`
	CommentCount  int64 `json:"comment_count"`
	DiggCount     int64 `json:"digg_count"`
	PlayCount     int64 `json:"play_count"`
	ShareCount    int64 `json:"share_count"`
	DownloadCount int64 `json:"download_count"`
	ForwardCount  int64 `json:"forward_count"`
}

type TikTokStatus struct {
	AllowComment  bool `json:"allow_comment"`
	AllowShare    bool `json:"allow_share"`
	IsDelete      bool `json:"is_delete"`
	IsProhibited  bool `json:"is_prohibited"`
	PrivateStatus int  `json:"private_status"`
	Reviewed      int  `json:"reviewed"`
}

type TikTokTakoInfo struct {
	TakoBubbleEnable bool `json:"tako_bubble_enable"`
}

type TikTokUpvoteInfo struct {
	UserUpvoted bool `json:"user_upvoted"`
}

type TikTokUpvotePreload struct {
	NeedPullUpvoteInfo bool `json:"need_pull_upvote_info"`
}

type TikTokVideo struct {
	Height       int             `json:"height"`
	Width        int             `json:"width"`
	Duration     int             `json:"duration"`
	Ratio        string          `json:"ratio"`
	PlayAddr     TikTokImage     `json:"play_addr"`
	Cover        TikTokImage     `json:"cover"`
	DynamicCover TikTokImage     `json:"dynamic_cover"`
	BitRate      []TikTokBitRate `json:"bit_rate"`
}

type TikTokBitRate struct {
	BitRate  int         `json:"bit_rate"`
	FPS      int         `json:"fps"`
	GearName string      `json:"gear_name"`
	PlayAddr TikTokImage `json:"play_addr"`
}

type TikTokVideoControl struct {
	AllowDownload bool `json:"allow_download"`
	AllowDuet     bool `json:"allow_duet"`
	AllowReact    bool `json:"allow_react"`
	AllowStitch   bool `json:"allow_stitch"`
}

type TikTokImage struct {
	Height  int      `json:"height"`
	Width   int      `json:"width"`
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

type TikTokMusic struct {
	ID          interface{} `json:"id"`
	IDStr       string      `json:"id_str"`
	Title       string      `json:"title"`
	Author      string      `json:"author"`
	Album       string      `json:"album"`
	Duration    int         `json:"duration"`
	PlayURL     TikTokImage `json:"play_url"`
	CoverThumb  TikTokImage `json:"cover_thumb"`
	CoverMedium TikTokImage `json:"cover_medium"`
	CoverLarge  TikTokImage `json:"cover_large"`
}
