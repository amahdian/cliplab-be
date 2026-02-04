package scrapecreators

// TwitterTweet represents the core data of a tweet.
type TwitterTweet struct {
	Typename       string       `json:"__typename"`
	RestID         string       `json:"rest_id"`
	Core           TwitterCore  `json:"core"`
	EditControl    interface{}  `json:"edit_control"`
	IsTranslatable bool         `json:"is_translatable"`
	Views          TwitterViews `json:"views"`
	Source         string       `json:"source"`
	Legacy         TweetLegacy  `json:"legacy"`
	URL            string       `json:"url"` // Often present in lists
}

// TwitterTweetResponse represents the response from the Twitter tweet detail endpoint.
type TwitterTweetResponse struct {
	Success          bool `json:"success"`
	CreditsRemaining int  `json:"credits_remaining"`
	TwitterTweet
}

// TwitterUserTweetsResponse represents the response from the Twitter user tweets endpoint.
type TwitterUserTweetsResponse struct {
	Success          bool           `json:"success"`
	CreditsRemaining int            `json:"credits_remaining"`
	Tweets           []TwitterTweet `json:"tweets"`
}

type TwitterCore struct {
	UserResults TwitterUserResults `json:"user_results"`
}

type TwitterUserResults struct {
	Result TwitterUserResult `json:"result"`
}

type TwitterUserResult struct {
	Typename       string            `json:"__typename"`
	ID             string            `json:"id"`
	RestID         string            `json:"rest_id"`
	Avatar         TwitterAvatar     `json:"avatar"`
	Core           TwitterUserCore   `json:"core"`
	IsBlueVerified bool              `json:"is_blue_verified"`
	Legacy         TwitterUserLegacy `json:"legacy"`
}

type TwitterAvatar struct {
	ImageURL string `json:"image_url"`
}

type TwitterUserCore struct {
	CreatedAt  string `json:"created_at"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type TwitterUserLegacy struct {
	Description      string `json:"description"`
	FollowersCount   int    `json:"followers_count"`
	FriendsCount     int    `json:"friends_count"`
	StatusesCount    int    `json:"statuses_count"`
	FavouritesCount  int    `json:"favourites_count"`
	MediaCount       int    `json:"media_count"`
	ProfileBannerURL string `json:"profile_banner_url"`
}

type TwitterViews struct {
	Count string `json:"count"`
	State string `json:"state"`
}

type TweetLegacy struct {
	BookmarkCount     int             `json:"bookmark_count"`
	CreatedAt         string          `json:"created_at"`
	ConversationIDStr string          `json:"conversation_id_str"`
	Entities          TwitterEntities `json:"entities"`
	ExtendedEntities  TwitterEntities `json:"extended_entities"`
	FavoriteCount     int             `json:"favorite_count"`
	FullText          string          `json:"full_text"`
	ReplyCount        int             `json:"reply_count"`
	RetweetCount      int             `json:"retweet_count"`
	QuoteCount        int             `json:"quote_count"`
	UserIDStr         string          `json:"user_id_str"`
	IDStr             string          `json:"id_str"`
	PossiblySensitive bool            `json:"possibly_sensitive"`
}

type TwitterEntities struct {
	Hashtags []interface{}  `json:"hashtags"`
	Media    []TwitterMedia `json:"media"`
	Urls     []interface{}  `json:"urls"`
	Symbols  []interface{}  `json:"symbols"`
}

type TwitterMedia struct {
	DisplayURL    string `json:"display_url"`
	ExpandedURL   string `json:"expanded_url"`
	IDStr         string `json:"id_str"`
	MediaURLHttps string `json:"media_url_https"`
	Type          string `json:"type"`
	URL           string `json:"url"`
}
