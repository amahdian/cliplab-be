package dtos

type InstagramItem struct {
	Urls              []*InstagramUrlInfo `json:"urls"`
	Meta              *InstagramMeta      `json:"meta"`
	PictureURL        string              `json:"pictureUrl"`
	PictureURLWrapped string              `json:"pictureUrlWrapped"`
	Service           string              `json:"service"`
}

type InstagramUrlInfo struct {
	URL       string `json:"url"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
}

type InstagramMeta struct {
	Title        string        `json:"title"`
	SourceURL    string        `json:"sourceUrl"`
	Shortcode    string        `json:"shortcode"`
	CommentCount int           `json:"commentCount"`
	LikeCount    int           `json:"likeCount"`
	TakenAt      int64         `json:"takenAt"` // Unix timestamp
	Comments     []interface{} `json:"comments"`
	Username     string        `json:"username"`
}
