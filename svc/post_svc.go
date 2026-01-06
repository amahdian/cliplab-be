package svc

import (
	"context"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/amahdian/cliplab-be/svc/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PostSvc interface {
	AddPostToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error)
	GetPostById(id string) (*resp.PostResponse, error)
}

type postSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs

	fileSvc     FileSvc
	RedisClient *redis.Client
}

func newPostSvc(
	ctx context.Context,
	stg storage.PgStorage,
	envs *env.Envs,
	redisClient *redis.Client,
	fileSvc FileSvc) PostSvc {
	return &postSvc{
		ctx:         ctx,
		stg:         stg,
		envs:        envs,
		RedisClient: redisClient,
		fileSvc:     fileSvc,
	}
}

func (s *postSvc) AddPostToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error) {
	platform := detectSocialMediaID(url)
	if platform != model.PlatformInstagram {
		return nil, errs.Newf(errs.InvalidArgument, nil, "unsupported platform, we only support Instagram for now")
	}

	estimatedTime := getEstimatedTimeByPlatform(platform)

	shortcode := utils.GetInstagramShortcode(url.String())
	post, err := s.stg.Post(s.ctx).FindByHashId(shortcode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrapf(err, "failed to find post by hash id %s", shortcode)
	}

	if post.ID != "" && post.Status != model.PostStatusFailed {
		if post.Status == model.PostStatusCompleted {
			estimatedTime = 0
		}
		return &resp.PostQueueResponse{
			Id:            post.ID,
			EstimatedTime: estimatedTime,
		}, nil
	}

	post = &model.Post{
		ID:     shortcode,
		UserIP: ip.String(),
		Link:   url.String(),
		Status: model.PostStatusPending,
	}

	if user.Id != uuid.Nil {
		post.UserId = &user.Id
	}

	if err := s.stg.Post(s.ctx).CreateOne(post); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to save post")
	}

	jsonData, err := json.Marshal(&model.PostQueueData{
		Id:       post.ID,
		Url:      post.Link,
		Platform: platform,
	})
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to marshal post data")
	}
	if err = s.RedisClient.LPush(s.ctx, global.RedisPostQueue, jsonData).Err(); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to publish post")
	}

	return &resp.PostQueueResponse{
		Id:            post.ID,
		EstimatedTime: estimatedTime,
	}, nil
}

func (s *postSvc) GetPostById(id string) (*resp.PostResponse, error) {
	p, err := s.stg.Post(s.ctx).FindByHashId(id)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find post by id")
	}

	post := *p
	if post.Status == model.PostStatusFailed {
		return nil, errs.Newf(errs.Internal, nil, "%s", *post.FailReason)
	}

	if post.Status != model.PostStatusCompleted {
		return &resp.PostResponse{
			Status:   post.Status,
			Platform: model.PlatformInstagram,
		}, nil
	}

	channel, err := s.stg.Channel(s.ctx).FindByHandler(post.UserAnchor)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find channel")
	}

	er := (float64(post.LikeCount+post.CommentCount) / float64(channel.LastHistory.FollowersCount)) * 100
	avgER := (float64(channel.LastHistory.AverageLikes+channel.LastHistory.AverageComments) / float64(channel.LastHistory.FollowersCount)) * 100

	res := &resp.PostResponse{
		Platform:              model.PlatformInstagram,
		Status:                post.Status,
		UserLink:              lo.ToPtr(post.UserProfileLink),
		UserHandler:           lo.ToPtr(post.UserName),
		ImageUrl:              post.ImageURL,
		VideoUrl:              post.VideoURL,
		LikeCount:             post.LikeCount,
		CommentCount:          post.CommentCount,
		ViewCount:             post.VideoPlayCount,
		PostDate:              post.PostDate,
		EngagementRate:        er,
		AverageLikeCount:      channel.LastHistory.AverageLikes,
		AverageCommentCount:   channel.LastHistory.AverageComments,
		AverageViewCount:      channel.LastHistory.AverageVideoPlays,
		AverageEngagementRate: avgER,
	}

	contents, err := s.stg.PostContent(s.ctx).ListByPostId(post.ID)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to list contents")
	}
	analysis, err := s.stg.PostAnalysis(s.ctx).FindByPostId(post.ID)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to get post analysis")
	}

	res.Analysis = analysis

	for _, content := range contents {
		switch content.Type {
		case model.ContentTranscript:
			metaData := content.Metadata.(*model.SegmentPostContentMetadata)
			res.Segments = append(res.Segments, &resp.PostContentSegmentResponse{
				PostContentResponse: &resp.PostContentResponse{
					Content:  content.Text,
					Language: content.Language,
				},
				Timestamp: metaData.Timestamp,
				Emotion:   metaData.Emotion,
				Speaker:   metaData.Speaker,
			})
		case model.ContentCaption:
			res.Caption = &resp.PostContentResponse{
				Content:  content.Text,
				Language: content.Language,
			}
		}
	}

	return res, nil
}

func detectSocialMediaID(url url.URL) model.SocialPlatform {
	text := strings.TrimSpace(url.String())
	text = strings.Split(text, "?")[0]
	text = strings.TrimSuffix(text, "/")

	// Patterns capture the ID as the first group
	youtubeRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([\w-]+)`)
	instagramRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?instagram\.com/(?:[^/]+/)?(?:p|reels?|tv)/([A-Za-z0-9_-]+)`)
	tiktokRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?tiktok\.com/@[\w.-]+/video/(\d+)`)
	twitterRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:twitter\.com|x\.com)/\w+/status/(\d+)`)

	if match := youtubeRegex.FindStringSubmatch(text); len(match) > 1 {
		return model.PlatformYouTube
	}
	if match := instagramRegex.FindStringSubmatch(text); len(match) > 1 {
		return model.PlatformInstagram
	}
	if match := tiktokRegex.FindStringSubmatch(text); len(match) > 1 {
		return model.PlatformTikTok
	}
	if match := twitterRegex.FindStringSubmatch(text); len(match) > 1 {
		return model.PlatformTwitter
	}

	return model.PlatformUnknown
}

func getEstimatedTimeByPlatform(platform model.SocialPlatform) int {
	switch platform {
	case model.PlatformInstagram, model.PlatformTikTok, model.PlatformTwitter:
		return 50
	case model.PlatformYouTube:
		return 120
	default:
		return 0
	}
}
