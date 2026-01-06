package svc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/clients/gemini"
	"github.com/amahdian/cliplab-be/clients/rocksolid"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc/utils"
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

	GeminiClient  gemini.Client
	ScraperClient rocksolid.Client
}

func newPostQueueSvc(ctx context.Context, stg storage.PgStorage, envs *env.Envs, geminiClient gemini.Client, scraperClient rocksolid.Client) QueueSvc {
	return &postQueueSvc{
		ctx:           ctx,
		stg:           stg,
		envs:          envs,
		GeminiClient:  geminiClient,
		ScraperClient: scraperClient,
	}
}

func (s *postQueueSvc) ProcessPost(url, id string, platform model.SocialPlatform) error {
	logger.Debug("Processing post:", url)

	post, err := s.stg.Post(s.ctx).FindByHashId(id)
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
	if post.Status != model.PostStatusPending {
		return
	}

	shortcode := utils.GetInstagramShortcode(post.Link)

	dto, err := s.ScraperClient.GetInstagramPost(shortcode)
	if err != nil {
		post.FailReason = lo.ToPtr(err.Error())
		_ = s.stg.Post(s.ctx).UpdateOne(post, false)
		return
	}

	// 1. Update initial post info
	post.Status = model.PostStatusProcessing
	post.UserName = dto.Owner.FullName
	post.UserAnchor = dto.Owner.Username
	post.UserProfileLink = fmt.Sprintf("https://instagram.com/%s", dto.Owner.Username)
	post.PostDate = time.Unix(dto.TakenAtTimestamp, 0)
	post.ImageURL = &dto.ThumbnailSrc
	post.VideoURL = &dto.VideoURL

	post.LikeCount = int64(dto.EdgeMediaPreviewLike.Count)
	post.CommentCount = int64(dto.EdgeMediaToParentComment.Count)
	post.VideoViewCount = int64(dto.VideoViewCount)
	post.VideoPlayCount = int64(dto.VideoPlayCount)

	if post.ChannelId == nil {
		channel, _ := s.stg.Channel(s.ctx).FindByHandler(dto.Owner.Username)
		if channel == nil {
			channel = &model.Channel{
				FullName: dto.Owner.FullName,
				Handler:  dto.Owner.Username,
				Platform: model.PlatformInstagram,
			}
			_ = s.stg.Channel(s.ctx).CreateOne(channel)
		}

		post.ChannelId = &channel.ID
	}

	_ = s.stg.Post(s.ctx).UpdateOne(post, false)

	if dto.VideoURL == "" {
		return
	}

	otherVideos, err := s.ScraperClient.GetInstagramPageReels(post.UserAnchor)
	if err != nil {
		post.FailReason = lo.ToPtr(err.Error())
		_ = s.stg.Post(s.ctx).UpdateOne(post, false)
		return
	}

	// 1.5 Save channel history
	totalLike := int64(0)
	totalComment := int64(0)
	totalViewCount := int64(0)
	totalPlayCount := int64(0)
	for _, reel := range otherVideos.Reels {
		totalLike += reel.Node.Media.LikeCount
		totalComment += reel.Node.Media.CommentCount
		totalPlayCount += reel.Node.Media.PlayCount
	}

	avgLikes := int64(0)
	avgComments := int64(0)
	avgViews := int64(0)
	avgPlays := int64(0)
	if len(otherVideos.Reels) > 0 {
		avgLikes = totalLike / int64(len(otherVideos.Reels))
		avgComments = totalComment / int64(len(otherVideos.Reels))
		avgViews = totalViewCount / int64(len(otherVideos.Reels))
		avgPlays = totalPlayCount / int64(len(otherVideos.Reels))
	}

	_ = s.stg.ChannelHistory(s.ctx).CreateOne(&model.ChannelHistory{
		ChannelID:         *post.ChannelId,
		FollowersCount:    dto.Owner.EdgeFollowedBy.Count,
		MediaCount:        dto.Owner.EdgeOwnerToTimelineMedia.Count,
		AverageLikes:      avgLikes,
		AverageComments:   avgComments,
		AverageVideoViews: avgViews,
		AverageVideoPlays: avgPlays,
	})

	detector := lingua.NewLanguageDetectorBuilder().FromAllSpokenLanguages().Build()

	// 2. Get Advanced Video Analysis from Gemini
	analysis, err := s.getInstagramVideoAnalysis(*dto, *otherVideos, detector)
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

	// A. Add Caption (from Instagram DTO)
	if len(dto.EdgeMediaToCaption.Edges) > 0 {
		lang, _ := detector.DetectLanguageOf(dto.EdgeMediaToCaption.Edges[0].Node.Text)
		contents = append(contents, &model.PostContent{
			PostID:   post.ID,
			Type:     model.ContentCaption,
			Text:     dto.EdgeMediaToCaption.Edges[0].Node.Text,
			Language: lang.IsoCode639_1().String(),
			Status:   model.PostStatusCompleted,
		})
	}

	// C. Add Transcript (from Gemini Segments)
	if len(analysis.Content.Segments) > 0 {
		for _, seg := range analysis.Content.Segments {
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

	// 4. Save all contents to database
	if len(contents) > 0 {
		if err = s.stg.PostContent(s.ctx).CreateMany(contents); err != nil {
			logger.Error("Failed to save post contents:", err)
		}
	}

	// 5. Save Post Analysis
	postAnalysis := &model.PostAnalysis{
		PostId:            post.ID,
		BigIdea:           analysis.Summary.BigIdea,
		WhyViral:          analysis.Summary.WhyViral,
		AudienceSentiment: analysis.Summary.AudienceSentiment,
		SentimentScore:    analysis.Summary.SentimentScore,
		Strengths:         analysis.Analysis.Strengths,
		Weaknesses:        analysis.Analysis.Weaknesses,
		HookIdeas:         analysis.Remix.HookIdeas,
		ScriptIdeas:       analysis.Remix.ScriptIdeas,
		Hashtags:          analysis.Publish.Hashtags,
		Captions: model.PostAnalysisCaptions{
			Casual:       analysis.Publish.Captions.Casual,
			Professional: analysis.Publish.Captions.Professional,
			Viral:        analysis.Publish.Captions.Viral,
		},
	}

	topicScore := 0
	hookScore := 0
	pacingScore := 0
	valueDeliveryScore := 0
	shareabilityScore := 0
	ctaScore := 0

	for _, m := range analysis.Analysis.Metrics {
		postAnalysis.Metrics = append(postAnalysis.Metrics, model.PostAnalysisMetric{
			Label:       m.Label,
			Score:       m.Score,
			Explanation: m.Explanation,
			Suggestion:  m.Suggestion,
		})

		if strings.Contains(strings.ToLower(m.Label), "top") {
			topicScore = m.Score
			continue
		}
		if strings.Contains(strings.ToLower(m.Label), "hook") {
			hookScore = m.Score
			continue
		}
		if strings.Contains(strings.ToLower(m.Label), "pacing") {
			pacingScore = m.Score
			continue
		}
		if strings.Contains(strings.ToLower(m.Label), "value") {
			valueDeliveryScore = m.Score
			continue
		}
		if strings.Contains(strings.ToLower(m.Label), "share") {
			shareabilityScore = m.Score
			continue
		}
		if strings.Contains(strings.ToLower(m.Label), "cta") {
			shareabilityScore = m.Score
			continue
		}
	}

	postAnalysis.ViralScore = int(getViralScore(topicScore, hookScore, pacingScore, valueDeliveryScore, shareabilityScore, ctaScore))

	if err = s.stg.PostAnalysis(s.ctx).CreateOne(postAnalysis); err != nil {
		logger.Error("Failed to save post analysis:", err)
	}

	// 6. Finalize Post status
	post.Status = model.PostStatusCompleted
	_ = s.stg.Post(s.ctx).UpdateOne(post, false)
}

func (s *postQueueSvc) getInstagramVideoAnalysis(dto rocksolid.ReelData, otherVideos rocksolid.Reels, detector lingua.LanguageDetector) (*gemini.AnalysisResponse, error) {

	// 2. We usually analyze the main video (first one) or the longest one.
	// Instagram carousels might have multiple videos, but for MVP we process the primary one.
	targetVideo := dto.VideoURL
	caption := ""
	language := "US"
	if len(dto.EdgeMediaToCaption.Edges) > 0 {
		caption = dto.EdgeMediaToCaption.Edges[0].Node.Text
		lang, _ := detector.DetectLanguageOf(dto.EdgeMediaToCaption.Edges[0].Node.Text)
		language = lang.IsoCode639_1().String()
	}

	coauthors := make([]string, len(dto.CoauthorProducers))
	for i, producer := range dto.CoauthorProducers {
		coauthors[i] = producer.Username
	}

	comments := make([]string, len(dto.EdgeMediaToParentComment.Edges))
	for i, edge := range dto.EdgeMediaToParentComment.Edges {
		comments[i] = edge.Node.Text
	}
	publishedAt := time.Unix(dto.TakenAtTimestamp, 0)

	er := float64(dto.EdgeMediaPreviewLike.Count+dto.EdgeMediaToParentComment.Count) / float64(dto.Owner.EdgeOwnerToTimelineMedia.Count)
	videoStats := map[string]float64{
		"like_count":      float64(dto.EdgeMediaPreviewLike.Count),
		"comment_count":   float64(dto.EdgeMediaToParentComment.Count),
		"play_count":      float64(dto.VideoPlayCount),
		"engagement_rate": er * 100,
	}

	totalLike := int64(0)
	totalComment := int64(0)
	totalViewCount := int64(0)
	totalPlayCount := int64(0)
	for _, reel := range otherVideos.Reels {
		totalLike += reel.Node.Media.LikeCount
		totalComment += reel.Node.Media.CommentCount
		totalViewCount += reel.Node.Media.ViewCount
		totalPlayCount += reel.Node.Media.PlayCount
	}

	averageStats := map[string]float64{}
	if len(otherVideos.Reels) > 0 && totalLike > 0 {
		averageStats = map[string]float64{
			"follower_count":          float64(dto.Owner.EdgeFollowedBy.Count),
			"average_like_count":      float64(totalLike) / float64(len(otherVideos.Reels)),
			"average_comment_count":   float64(totalComment) / float64(len(otherVideos.Reels)),
			"average_play_count":      float64(totalPlayCount) / float64(len(otherVideos.Reels)),
			"average_engagement_rate": (float64(totalLike+totalComment) / float64(int64(len(otherVideos.Reels))*dto.Owner.EdgeFollowedBy.Count)) * 100,
		}
	}

	logger.Infof("Starting AI analysis for video: %s", targetVideo)

	// 3. Call the Gemini client using the direct URL
	// We pass "Instagram" as the platform context
	result, err := s.GeminiClient.AnalyzeVideo(
		model.PlatformInstagram,
		targetVideo,
		caption,
		coauthors,
		comments,
		videoStats,
		averageStats,
		publishedAt,
		language)
	if err != nil {
		logger.Errorf("Failed to analyze video via Gemini: %v", err)
		return nil, errs.Newf(errs.Internal, err, "failed to analyze video content")
	}

	logger.Infof("Finished analysis for video: %s", targetVideo)

	return result, nil
}

func getViralScore(topicScore, hookScore, pacingScore, valueDeliveryScore, shareabilityScore, ctaScore int) float64 {
	return (float64(hookScore) * 0.25) +
		(float64(topicScore) * 0.2) +
		(float64(pacingScore) * 0.15) +
		(float64(valueDeliveryScore) * 0.15) +
		(float64(shareabilityScore) * 0.15) +
		(float64(ctaScore) * 0.1)
}
