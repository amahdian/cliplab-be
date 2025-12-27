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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PostSvc interface {
	AddPostToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error)
	GetPostById(id uuid.UUID) (*resp.PostResponse, error)
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

	post, err := s.stg.Post(s.ctx).FindByUrl(url.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrapf(err, "failed to find post by url %s", url.String())
	}
	if post != nil && post.Status != model.PostStatusFailed {
		if post.Status == model.PostStatusCompleted {
			estimatedTime = 0
		}
		return &resp.PostQueueResponse{
			Id:            post.ID,
			EstimatedTime: estimatedTime,
		}, nil
	}

	post = &model.Post{
		UserIP:   ip.String(),
		Link:     url.String(),
		Platform: platform,
		Status:   model.PostStatusPending,
	}

	if user.Id != uuid.Nil {
		post.UserId = &user.Id
	}

	if err := s.stg.Post(s.ctx).CreateOne(post); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to save post")
	}

	jsonData, err := json.Marshal(&model.PostQueueData{
		Id:       post.ID.String(),
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

func (s *postSvc) GetPostById(id uuid.UUID) (*resp.PostResponse, error) {
	post, err := s.stg.Post(s.ctx).FindById(id)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to find post by id")
	}

	if post.Status == model.PostStatusFailed {
		return nil, errs.Newf(errs.Internal, nil, *post.FailReason)
	}

	if post.Status != model.PostStatusCompleted {
		return &resp.PostResponse{
			Status:   post.Status,
			Platform: model.PlatformInstagram,
		}, nil
	}

	res := &resp.PostResponse{
		Status:   post.Status,
		Platform: post.Platform,
		UserLink: lo.ToPtr(post.UserProfileLink),
		ImageUrl: post.ImageURL,
		VideoUrl: post.VideoURL,
	}

	contents, err := s.stg.PostContent(s.ctx).ListByPostId(post.ID)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to list contents")
	}

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
		case model.ContentKeyPoint:
			res.KeyPoints = append(res.KeyPoints, &resp.PostContentResponse{
				Content:  content.Text,
				Language: content.Language,
			})
		case model.ContentSummary:
			res.Summary = &resp.PostContentResponse{
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
		return 15
	case model.PlatformYouTube:
		return 30
	default:
		return 0
	}
}
