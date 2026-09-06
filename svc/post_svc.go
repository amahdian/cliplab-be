package svc

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
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
	Image(imgURL string) ([]byte, string, error)
}

type analyzeSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs

	fileSvc     FileSvc
	RedisClient *redis.Client
	creditSvc   CreditSvc
	historySvc  HistorySvc
}

func newAnalyzeSvc(
	ctx context.Context,
	stg storage.PgStorage,
	envs *env.Envs,
	redisClient *redis.Client,
	fileSvc FileSvc,
	creditSvc CreditSvc,
	historySvc HistorySvc) AnalyzeSvc {
	return &analyzeSvc{
		ctx:         ctx,
		stg:         stg,
		envs:        envs,
		RedisClient: redisClient,
		fileSvc:     fileSvc,
		creditSvc:   creditSvc,
		historySvc:  historySvc,
	}
}

func (s *analyzeSvc) AddRequestToAnalyzeQueue(url url.URL, user *auth.UserInfo, ip net.IP) (*resp.PostQueueResponse, error) {
	platform := utils.DetectSocialMediaID(url)
	if platform != model.PlatformInstagram {
		return nil, errs.Newf(errs.InvalidArgument, nil, "unsupported platform, we only support Instagram reels for now")
	}

	estimatedTime := getEstimatedTimeByPlatform(platform)

	shortcode := utils.GetInstagramShortcode(url.String())
	post, err := s.stg.Post(s.ctx).FindByHashId(shortcode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrapf(err, "failed to find post by hash id %s", shortcode)
	}

	var result *resp.PostQueueResponse
	var analyzeRequest *model.AnalyzeRequest
	var creditsUsed int
	var remainingCredits int

	// Handle history and credit reversal
	defer func() {
		if user == nil || user.Id == uuid.Nil {
			return
		}

		if err != nil {
			// Refund credits if there was an error during queuing phase
			if creditsUsed > 0 {
				_ = s.creditSvc.AddCredits(user.Id, creditsUsed)
				remainingCredits += creditsUsed
			}
			// Log failure history immediately
			_ = s.historySvc.StoreUsage(user.Id, model.ToolVideoAnalysis, creditsUsed, remainingCredits, url.String(), platform, "analysis", model.UsageStatusFailed, nil, err.Error())
		} else if analyzeRequest != nil {
			// Success in queuing: store the "Processing" record
			_ = s.historySvc.StoreUsage(user.Id, model.ToolVideoAnalysis, creditsUsed, remainingCredits, url.String(), platform, "analysis", model.UsageStatusProcessing, &analyzeRequest.ID, "")
		}
	}()

	now := time.Now()

	if user == nil || user.Id == uuid.Nil {
		err = errs.Newf(errs.PermissionDenied, nil, "payment required")
		return nil, err
	}

	// Deduct credits for registered users
	if rem, creditErr := s.creditSvc.CheckAndDeduct(user.Id, model.CreditKeyReelAnalyze); creditErr == nil {
		remainingCredits = rem
		creditsUsed = model.GetCreditRule(model.CreditKeyReelAnalyze).Amount
	} else {
		err = creditErr
		return nil, err
	}

	if analyzeRequest == nil {
		analyzeRequest = &model.AnalyzeRequest{
			ID:     uuid.New(),
			Link:   url.String(),
			UserId: &user.Id,
			UserIP: string(ip),
			Status: model.RequestStatusPending,
		}

		if post != nil {
			analyzeRequest.PostId = &post.ID
		}
	}

	if analyzeRequest.Status == model.RequestStatusCompleted {
		if post.UpdatedAt.After(now.Add(-1200 * time.Hour)) {
			result = &resp.PostQueueResponse{
				RequestId:     analyzeRequest.ID.String(),
				EstimatedTime: 0,
			}
			return result, nil
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

	result = &resp.PostQueueResponse{
		RequestId:     analyzeRequest.ID.String(),
		EstimatedTime: estimatedTime,
	}
	return result, nil
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

func (s *analyzeSvc) Image(imgURL string) ([]byte, string, error) {
	u, err := url.Parse(imgURL)
	if err != nil {
		return nil, "", errs.Newf(errs.InvalidArgument, err, "invalid image url")
	}

	// Instagram and Facebook CDN hosts validation
	isValidHost := false
	validHosts := []string{"cdninstagram.com", "fbcdn.net", "instagram.com"}
	for _, host := range validHosts {
		if strings.HasSuffix(u.Host, host) {
			isValidHost = true
			break
		}
	}

	if !isValidHost {
		return nil, "", errs.Newf(errs.InvalidArgument, nil, "not a valid instagram image host")
	}

	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, "", errs.Newf(errs.Internal, err, "failed to fetch image")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errs.Newf(errs.Internal, nil, "failed to fetch image, status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errs.Newf(errs.Internal, err, "failed to read image body")
	}

	return data, resp.Header.Get("Content-Type"), nil
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
