package svc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/amahdian/cliplab-be/clients"
	"github.com/amahdian/cliplab-be/clients/dtos"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/google/uuid"
	"github.com/pemistahl/lingua-go"
	"github.com/samber/lo"
)

type QueueSvc interface {
	ProcessPost(url, id string, platform model.SocialPlatform) error
}

type postQueueSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs

	GeminiClient   clients.GeminiClient
	RapidApiClient clients.RapidApiClient
}

func newPostQueueSvc(ctx context.Context, stg storage.PgStorage, envs *env.Envs, geminiClient clients.GeminiClient, rapidApiClient clients.RapidApiClient) QueueSvc {
	return &postQueueSvc{
		ctx:            ctx,
		stg:            stg,
		envs:           envs,
		GeminiClient:   geminiClient,
		RapidApiClient: rapidApiClient,
	}
}

func (s *postQueueSvc) ProcessPost(url, id string, platform model.SocialPlatform) error {
	logger.Debug("Processing post:", url)

	postId, _ := uuid.Parse(id)

	post, err := s.stg.Post(s.ctx).FindById(postId)
	if err != nil {
		return err
	}

	switch platform {
	case model.PlatformInstagram:
		s.processInstagramScrap(post)
	}

	return nil
}

func (s *postQueueSvc) processInstagramScrap(post *model.Post) {
	shortcode := getInstagramShortcode(post.Link)
	dtos, err := s.RapidApiClient.GetInstagramPost(shortcode)
	if err != nil {
		logger.Error(err)
		return
	}

	dto := dtos[0]

	if post.Status != model.PostStatusPending {
		return
	}
	if len(dto.Urls) == 0 {
		return
	}

	// 1. Update initial post info
	post.Status = model.PostStatusProcessing
	post.UserName = dto.Meta.Username
	post.UserAnchor = dto.Meta.Username
	post.UserProfileLink = fmt.Sprintf("https://instagram.com/%s", dto.Meta.Username)
	post.PostDate = time.Unix(dto.Meta.TakenAt, 0)
	post.ImageURL = &dto.PictureURL
	post.VideoURL = &dto.Urls[0].URL
	_ = s.stg.Post(s.ctx).UpdateOne(post, false)

	// 2. Get Advanced Video Analysis from Gemini
	analysis, err := s.getInstagramVideoAnalysis(*dto)
	if err != nil {
		failReason := fmt.Sprintf("Gemini analysis failed: %s", err.Error())
		logger.Error(failReason)

		post.Status = model.PostStatusFailed
		post.FailReason = lo.ToPtr(failReason)
		_ = s.stg.Post(s.ctx).UpdateOne(post, false)
		return
	}
	if analysis == nil {
		failReason := "Gemini analysis is empty"
		logger.Error(failReason)

		post.FailReason = lo.ToPtr(failReason)
		post.Status = model.PostStatusFailed
		return
	}

	// 3. Prepare PostContent slice
	var contents []*model.PostContent
	detector := lingua.NewLanguageDetectorBuilder().FromAllSpokenLanguages().Build()

	// A. Add Caption (from Instagram DTO)
	if dto.Meta.Title != "" {
		lang, _ := detector.DetectLanguageOf(dto.Meta.Title)
		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentCaption,
			Text:     dto.Meta.Title,
			Language: lang.IsoCode639_1().String(),
			Status:   model.PostStatusCompleted,
		})
	}

	// B. Add Summary (from Gemini)
	if analysis.Summary != "" {
		// We know Gemini followed our language instruction for summary
		lang, _ := detector.DetectLanguageOf(analysis.Summary)

		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentSummary,
			Text:     analysis.Summary,
			Language: lang.IsoCode639_1().String(),
			Status:   model.PostStatusCompleted,
		})
	}

	// C. Add Transcript (from Gemini Segments)
	if len(analysis.Segments) > 0 {
		for _, seg := range analysis.Segments {
			lang, _ := detector.DetectLanguageOf(seg.Content)

			contents = append(contents, &model.PostContent{
				PostID:   post.ID,
				Type:     model.ContentTranscript,
				Text:     seg.Content,
				Language: lang.IsoCode639_1().String(),
				Status:   model.PostStatusCompleted,
				Metadata: &model.SegmentPostContentMetadata{
					Timestamp: seg.Timestamp,
					Speaker:   seg.Speaker,
					Emotion:   seg.Emotion,
				},
			})
		}
	}

	// D. Add Trend Metadata
	if analysis.TrendMetadata != "" {
		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentTrendMetadata,
			Text:     analysis.TrendMetadata,
			Language: string(model.LanguageEnglish),
			Status:   model.PostStatusCompleted,
		})
	}

	// E. Add Giveaway
	if analysis.Giveaway.IsDetected {
		giveawayText := fmt.Sprintf("Prize: %s\nRequirements: %s\nDeadline: %s", analysis.Giveaway.Prize, analysis.Giveaway.Requirements, analysis.Giveaway.Deadline)

		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentTrendMetadata,
			Text:     giveawayText,
			Language: string(model.LanguageEnglish),
			Status:   model.PostStatusCompleted,
			Metadata: &model.GiveawayPostContentMetadata{
				Prize:        analysis.Giveaway.Prize,
				Deadline:     analysis.Giveaway.Deadline,
				Requirements: analysis.Giveaway.Requirements,
			},
		})
	}

	// F. Add Key Points
	if len(analysis.KeyPoints) > 0 {
		for _, point := range analysis.KeyPoints {
			lang, _ := detector.DetectLanguageOf(analysis.Summary)

			contents = append(contents, &model.PostContent{
				PostID:   post.ID,
				Type:     model.ContentKeyPoint,
				Text:     point,
				Language: lang.IsoCode639_1().String(),
				Status:   model.PostStatusCompleted,
			})
		}
	}

	// F. Add Key Points
	if analysis.Hook != "" && len(analysis.Segments) > 0 {
		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentKeyPoint,
			Text:     analysis.Hook,
			Language: analysis.Segments[0].Language,
			Status:   model.PostStatusCompleted,
		})
	}

	// 4. Save all contents to database
	if len(contents) > 0 {
		if err = s.stg.PostContent(s.ctx).CreateMany(contents); err != nil {
			logger.Error("Failed to save post contents:", err)
		}
	}

	// 5. Finalize Post status
	post.Status = model.PostStatusCompleted
	_ = s.stg.Post(s.ctx).UpdateOne(post, false)
}

func (s *postQueueSvc) getInstagramVideoAnalysis(dto dtos.InstagramItem) (*dtos.AnalysisResponse, error) {

	// 2. We usually analyze the main video (first one) or the longest one.
	// Instagram carousels might have multiple videos, but for MVP we process the primary one.
	targetVideo := dto.Urls[0]

	logger.Infof("Starting AI analysis for video: %s", targetVideo.URL)

	// 3. Call the Gemini client using the direct URL
	// We pass "Instagram" as the platform context
	result, err := s.GeminiClient.AnalyzeVideo("Instagram", targetVideo.URL)
	if err != nil {
		logger.Errorf("Failed to analyze video via Gemini: %v", err)
		return nil, errs.Newf(errs.Internal, err, "failed to analyze video content")
	}

	return result, nil
}

func getInstagramShortcode(urlStr string) string {
	// regex explanation:
	// reels? -> matches 'reel' or 'reels'
	// ([A-Za-z0-9_-]+) -> Capture Group 1: the actual shortcode
	re := regexp.MustCompile(`instagram\.com/(?:[^/]+/)?(?:p|reels?|tv)/([A-Za-z0-9_-]+)`)

	match := re.FindStringSubmatch(urlStr)

	// match[0] is the full string that matched
	// match[1] is the first capture group (our shortcode)
	if len(match) > 1 {
		return match[1]
	}

	return ""
}
