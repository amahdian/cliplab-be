package svc

import (
	"context"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

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

type AnalyzeSvc interface {
	AddRequestToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error)
	GetAnalyzeResult(id uuid.UUID) (*resp.AnalyzeResult, error)
}

type analyzeSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs

	fileSvc     FileSvc
	RedisClient *redis.Client
}

func newAnalyzeSvc(
	ctx context.Context,
	stg storage.PgStorage,
	envs *env.Envs,
	redisClient *redis.Client,
	fileSvc FileSvc) AnalyzeSvc {
	return &analyzeSvc{
		ctx:         ctx,
		stg:         stg,
		envs:        envs,
		RedisClient: redisClient,
		fileSvc:     fileSvc,
	}
}

func (s *analyzeSvc) AddRequestToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error) {
	platform := detectSocialMediaID(url)
	if platform != model.PlatformInstagram {
		return nil, errs.Newf(errs.InvalidArgument, nil, "unsupported platform, we only support Instagram reels for now")
	}

	estimatedTime := getEstimatedTimeByPlatform(platform)

	shortcode := utils.GetInstagramShortcode(url.String())
	post, err := s.stg.Post(s.ctx).FindByHashId(shortcode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrapf(err, "failed to find post by hash id %s", shortcode)
	}

	now := time.Now()

	if post.ID == "" && user.Id == uuid.Nil {
		// check the rate limit
		requestCount, err := s.stg.AnalyzeRequest(s.ctx).CountByIpAndDate(ip, now)
		if err == nil && requestCount >= 2 {
			return nil, errs.Newf(errs.PermissionDenied, nil, "payment required")
		}
	}

	var analyzeRequest *model.AnalyzeRequest

	// get request for the post
	var requests []*model.AnalyzeRequest
	if post.ID != "" {
		requests, err = s.stg.AnalyzeRequest(s.ctx).ListByPostId(post.ID)
		if err != nil {
			return nil, errs.Wrapf(err, "failed to find request by post id %s", post.ID)
		}
	}

	for _, request := range requests {
		if user.Id != uuid.Nil && request.UserId != nil && *request.UserId == user.Id {
			analyzeRequest = request
			break
		}
		if user.Id == uuid.Nil && request.UserId == nil && request.UserIP == ip.String() {
			analyzeRequest = request
			break
		}
	}
	if analyzeRequest == nil {
		analyzeRequest = &model.AnalyzeRequest{
			UserIP: ip.String(),
			Link:   url.String(),
			Status: model.RequestStatusPending,
		}
		if user.Id != uuid.Nil {
			analyzeRequest.UserId = &user.Id
		}
	}

	if len(requests) > 0 {
		analyzeRequest.Status = requests[0].Status
	}
	if post.ID != "" {
		analyzeRequest.PostId = lo.ToPtr(post.ID)
	}

	if err := s.stg.AnalyzeRequest(s.ctx).UpsertOne(analyzeRequest, false); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to save analyze requet")
	}

	if analyzeRequest.Status == model.RequestStatusCompleted {
		if post.UpdatedAt.After(now.Add(-1200 * time.Hour)) {
			return &resp.PostQueueResponse{
				RequestId:     analyzeRequest.ID.String(),
				EstimatedTime: 0,
			}, nil
		} else {
			jsonData, _ := json.Marshal(&model.PostQueueData{
				Id:       analyzeRequest.ID,
				PostId:   analyzeRequest.PostId,
				Url:      analyzeRequest.Link,
				Platform: platform,
			})
			s.RedisClient.LPush(s.ctx, global.RedisPostRenewQueue, jsonData)

			return &resp.PostQueueResponse{
				RequestId:     post.ID,
				EstimatedTime: 10,
			}, nil
		}
	}

	jsonData, err := json.Marshal(&model.PostQueueData{
		Id:       analyzeRequest.ID,
		PostId:   analyzeRequest.PostId,
		Url:      analyzeRequest.Link,
		Platform: platform,
	})
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to marshal post data")
	}
	if err = s.RedisClient.LPush(s.ctx, global.RedisPostQueue, jsonData).Err(); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to publish post")
	}

	return &resp.PostQueueResponse{
		RequestId:     post.ID,
		EstimatedTime: estimatedTime,
	}, nil
}

func (s *analyzeSvc) GetAnalyzeResult(id uuid.UUID) (*resp.AnalyzeResult, error) {
	r, err := s.stg.AnalyzeRequest(s.ctx).FindById(id)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find post by id")
	}

	request := *r
	if request.Status == model.RequestStatusFailed {
		return nil, errs.Newf(errs.Internal, nil, "Failed to analyze the post. Please try again later.")
	}

	if request.Status != model.RequestStatusCompleted {
		return &resp.AnalyzeResult{
			Status:   request.Status,
			Platform: model.PlatformInstagram,
		}, nil
	}

	p, err := s.stg.Post(s.ctx).FindByHashId(*request.PostId)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find post by hash id")
	}

	post := *p
	channel, err := s.stg.Channel(s.ctx).FindByHandler(post.UserAnchor)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find channel")
	}

	er := (float64(post.LikeCount+post.CommentCount) / float64(channel.LastHistory.FollowersCount)) * 100
	avgER := (float64(channel.LastHistory.AverageLikes+channel.LastHistory.AverageComments) / float64(channel.LastHistory.FollowersCount)) * 100

	res := &resp.AnalyzeResult{
		Platform:              model.PlatformInstagram,
		Status:                request.Status,
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
	//instagramRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?instagram\.com/(?:[^/]+/)?(?:p|reels?|tv)/([A-Za-z0-9_-]+)`)
	instagramRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?instagram\.com/(?:reels?|reel)/([A-Za-z0-9_-]+)`)
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
		return 60
	case model.PlatformYouTube:
		return 120
	default:
		return 0
	}
}
