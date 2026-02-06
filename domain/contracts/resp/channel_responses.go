package resp

import "time"

type ChannelEngagementResponse struct {
	ProfileImage          string           `json:"profileImage"`
	ProfileUrl            string           `json:"profileUrl"`
	ProfileDescriptor     string           `json:"profileDescriptor"`
	FollowersCount        int64            `json:"followersCount"`
	FollowingCount        int64            `json:"followingCount"`
	PostsCount            int64            `json:"postsCount"`
	AverageLikeCount      float64          `json:"averageLikeCount"`
	AverageCommentCount   float64          `json:"averageCommentCount"`
	AveragePlayCount      float64          `json:"averagePlayCount"`
	AverageShareCount     float64          `json:"averageShareCount"`
	AverageEngagementRate float64          `json:"averageEngagementRate"`
	LatestPostPublishDate time.Time        `json:"latestPostPublishDate"`
	Breakdown             []*PostBreakdown `json:"breakdown"`
	Credits               *int             `json:"credits,omitempty"`
}

type PostBreakdown struct {
	Image          string    `json:"image"`
	Link           string    `json:"link"`
	PublishDate    time.Time `json:"publishDate"`
	LikeCount      int64     `json:"likeCount"`
	CommentCount   int64     `json:"commentCount"`
	ViewCount      int64     `json:"viewCount"`
	PlayCount      int64     `json:"playCount"`
	ShareCount     int64     `json:"shareCount"`
	EngagementRate float64   `json:"engagementRate"`
}
