package svc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/clients/gemini"
	"github.com/amahdian/cliplab-be/clients/scrapecreators"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc/utils"
	"github.com/google/uuid"
	"github.com/pemistahl/lingua-go"
	"github.com/samber/lo"
)

type QueueSvc interface {
	ProcessRequest(url string, requestId uuid.UUID, platform model.SocialPlatform) error
	RenewPost(url string, requestId uuid.UUID, platform model.SocialPlatform) error
}

type postQueueSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs

	GeminiClient  gemini.Client
	ScraperClient scrapecreators.Client
	historySvc    HistorySvc
	creditSvc     CreditSvc
}

func newPostQueueSvc(ctx context.Context, stg storage.PgStorage, envs *env.Envs, geminiClient gemini.Client, scraperClient scrapecreators.Client, historySvc HistorySvc, creditSvc CreditSvc) QueueSvc {
	return &postQueueSvc{
		ctx:           ctx,
		stg:           stg,
		envs:          envs,
		GeminiClient:  geminiClient,
		ScraperClient: scraperClient,
		historySvc:    historySvc,
		creditSvc:     creditSvc,
	}
}

func (s *postQueueSvc) ProcessRequest(url string, requestId uuid.UUID, platform model.SocialPlatform) (err error) {
	logger.Debug("Processing post:", url)

	request, err := s.stg.AnalyzeRequest(s.ctx).FindById(requestId)
	if err != nil {
		return err
	}

	// Handle history and potential refund on background failure
	defer func() {
		status := model.UsageStatusCompleted
		respLog := request.LlmResponse

		if err != nil {
			status = model.UsageStatusFailed
			respLog = err.Error()

			// If it's a permanent failure, refund the credits
			// We can check if it's a registered user
			if request.UserId != nil {
				credits := model.GetCreditRule(model.CreditKeyReelAnalyze).Amount
				_ = s.creditSvc.AddCredits(*request.UserId, credits)
			}
		}

		_ = s.historySvc.UpdateUsageByRequestID(requestId, status, respLog)
	}()

	var post *model.Post
	if request.PostId != nil {
		post, err = s.stg.Post(s.ctx).FindByHashId(*request.PostId)
		if err != nil {
			return err
		}
	} else {
		shortcode := utils.GetInstagramShortcode(request.Link)
		post = &model.Post{
			ID:         shortcode,
			Link:       request.Link,
			UserName:   "",
			UserAnchor: "",
		}

		err = s.stg.Post(s.ctx).CreateOne(post)
		if err != nil {
			return err
		}

		request.PostId = lo.ToPtr(post.ID)
	}

	request.Status = model.RequestStatusProcessing
	_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)

	switch platform {
	case model.PlatformInstagram:
		var reelDto *scrapecreators.ReelData
		var otherReelsDto *scrapecreators.ReelsResponse
		reelDto, otherReelsDto, err = s.renewInstagramScrap(post)
		if err != nil {
			request.FailReason = lo.ToPtr(err.Error())
			request.Status = model.RequestStatusFailed
			_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)
			return err
		} else {
			err = s.processInstagramScrap(request, reelDto, otherReelsDto)
			if err != nil {
				request.Status = model.RequestStatusFailed
				request.FailReason = lo.ToPtr(err.Error())
				_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)
				return err
			}
		}
	}

	request.FailReason = nil
	request.Status = model.RequestStatusCompleted
	_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)

	return nil
}

func (s *postQueueSvc) RenewPost(url string, requestId uuid.UUID, platform model.SocialPlatform) (err error) {
	logger.Debug("Renew post stats:", url)

	request, err := s.stg.AnalyzeRequest(s.ctx).FindById(requestId)
	if err != nil {
		return err
	}

	// Handle history and potential refund
	defer func() {
		status := model.UsageStatusCompleted
		respLog := ""

		if err != nil {
			status = model.UsageStatusFailed
			respLog = err.Error()
			// For analysis, we might want to refund if background renew fails
			if request.UserId != nil {
				credits := model.GetCreditRule(model.CreditKeyReelAnalyze).Amount
				_ = s.creditSvc.AddCredits(*request.UserId, credits)
			}
		}

		_ = s.historySvc.UpdateUsageByRequestID(requestId, status, respLog)
	}()

	request.Status = model.RequestStatusProcessing
	_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)

	post, err := s.stg.Post(s.ctx).FindByHashId(*request.PostId)
	if err != nil {
		return err
	}

	switch platform {
	case model.PlatformInstagram:
		_, _, err = s.renewInstagramScrap(post)
	}

	if err != nil {
		request.Status = model.RequestStatusFailed
		request.FailReason = lo.ToPtr(err.Error())
	} else {
		request.Status = model.RequestStatusCompleted
	}

	_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)

	return err
}

func (s *postQueueSvc) processInstagramScrap(request *model.AnalyzeRequest, reelDto *scrapecreators.ReelData, otherReelsDto *scrapecreators.ReelsResponse) error {
	detector := lingua.NewLanguageDetectorBuilder().FromAllSpokenLanguages().Build()

	// 2. Get Advanced Video Analysis from Gemini
	llmRequest, llmResponse, analysis, err := s.getInstagramVideoAnalysis(*reelDto, *otherReelsDto, detector)
	request.LlmRequest = llmRequest
	request.LlmResponse = llmResponse
	_ = s.stg.AnalyzeRequest(s.ctx).UpdateOne(request, false)

	if err != nil {
		failReason := fmt.Sprintf("Gemini analysis failed: %s", err.Error())
		return errs.Wrapf(err, "Gemini analysis failed: %s", failReason)
	}
	if analysis == nil {
		return errs.Newf(errs.Internal, nil, "Gemini analysis is empty")
	}

	// 3. Prepare PostContent slice
	var contents []*model.PostContent

	// A. Add Caption (from Instagram reelDto)
	if len(reelDto.EdgeMediaToCaption.Edges) > 0 {
		lang, ok := detector.DetectLanguageOf(reelDto.EdgeMediaToCaption.Edges[0].Node.Text)
		langStr := lang.IsoCode639_1().String()
		if !ok {
			langStr = ""
		}
		contents = append(contents, &model.PostContent{
			PostID:   *request.PostId,
			Type:     model.ContentCaption,
			Text:     reelDto.EdgeMediaToCaption.Edges[0].Node.Text,
			Language: langStr,
		})
	}

	// C. Add Transcript (from Gemini Segments)
	if len(analysis.Content.Segments) > 0 {
		for _, seg := range analysis.Content.Segments {
			lang, ok := detector.DetectLanguageOf(seg.Content)
			langStr := lang.IsoCode639_1().String()
			if !ok {
				langStr = ""
			}

			contents = append(contents, &model.PostContent{
				PostID:   *request.PostId,
				Type:     model.ContentTranscript,
				Text:     seg.Content,
				Language: langStr,
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

	hashtags := make([]string, len(analysis.Publish.Hashtags))
	for i, hashtag := range analysis.Publish.Hashtags {
		hashtags[i] = strings.ReplaceAll(hashtag, "#", "")
	}

	// 5. Save Post Analysis
	postAnalysis := &model.PostAnalysis{
		PostId:            *request.PostId,
		BigIdea:           analysis.Summary.BigIdea,
		WhyViral:          analysis.Summary.WhyViral,
		AudienceSentiment: analysis.Summary.AudienceSentiment,
		SentimentScore:    analysis.Summary.SentimentScore,
		Verdict:           analysis.Summary.Verdict,
		Strengths:         analysis.Analysis.Strengths,
		Weaknesses:        analysis.Analysis.Weaknesses,
		HookIdeas:         analysis.Remix.HookIdeas,
		ScriptIdeas:       analysis.Remix.ScriptIdeas,
		Hashtags:          hashtags,
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
			ctaScore = m.Score
			continue
		}
	}

	postAnalysis.ViralScore = int(getViralScore(
		topicScore, hookScore, pacingScore, valueDeliveryScore,
		shareabilityScore, ctaScore,
		analysis.Analysis.Scope.Confidence, analysis.Analysis.Scope.Level))

	if err = s.stg.PostAnalysis(s.ctx).CreateOne(postAnalysis); err != nil {
		logger.Error("Failed to save post analysis:", err)
	}

	return nil
}

func (s *postQueueSvc) renewInstagramScrap(post *model.Post) (*scrapecreators.ReelData, *scrapecreators.ReelsResponse, error) {
	dto, err := s.ScraperClient.GetInstagramPost(post.Link)
	if err != nil {
		return nil, nil, err
	}

	// 1. Update initial post info
	post.UserName = dto.Owner.FullName
	post.UserAnchor = dto.Owner.Username
	post.UserProfileLink = fmt.Sprintf("https://instagram.com/%s", dto.Owner.Username)
	post.PostDate = time.Unix(dto.TakenAtTimestamp, 0)
	post.ImageURL = &dto.ThumbnailSrc
	post.VideoURL = &dto.VideoURL

	post.LikeCount = dto.EdgeMediaPreviewLike.Count
	post.CommentCount = dto.EdgeMediaToParentComment.Count
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

	_ = s.stg.Post(s.ctx).UpsertOne(post, false)

	if dto.VideoURL == "" {
		return dto, nil, nil
	}

	otherReelsDto, err := s.ScraperClient.GetInstagramPageReels(post.UserAnchor)
	if err != nil {
		return nil, nil, err
	}

	// 1.5 Save channel history
	totalLike := int64(0)
	totalComment := int64(0)
	totalViewCount := int64(0)
	totalPlayCount := int64(0)
	for _, reel := range otherReelsDto.Items {
		totalLike += reel.Media.LikeCount
		totalComment += reel.Media.CommentCount
		totalPlayCount += reel.Media.PlayCount
	}

	avgLikes := int64(0)
	avgComments := int64(0)
	avgViews := int64(0)
	avgPlays := int64(0)
	if len(otherReelsDto.Items) > 0 {
		avgLikes = totalLike / int64(len(otherReelsDto.Items))
		avgComments = totalComment / int64(len(otherReelsDto.Items))
		avgViews = totalViewCount / int64(len(otherReelsDto.Items))
		avgPlays = totalPlayCount / int64(len(otherReelsDto.Items))
	}

	_ = s.stg.ChannelHistory(s.ctx).CreateOne(&model.ChannelHistory{
		ChannelID:         *post.ChannelId,
		FollowersCount:    dto.Owner.EdgeFollowedBy.Count,
		MediaCount:        dto.Owner.EdgeOwnerToTimelineMedia.Count,
		AverageLikes:      float64(avgLikes),
		AverageComments:   float64(avgComments),
		AverageVideoViews: float64(avgViews),
		AverageVideoPlays: float64(avgPlays),
	})

	_ = s.stg.Post(s.ctx).UpdateOne(post, false)

	return dto, otherReelsDto, nil
}

func (s *postQueueSvc) getInstagramVideoAnalysis(dto scrapecreators.ReelData, otherReelsDto scrapecreators.ReelsResponse, detector lingua.LanguageDetector) (string, string, *gemini.AnalysisResponse, error) {

	// 2. We usually analyze the main video (first one) or the longest one.
	// Instagram carousels might have multiple videos, but for MVP we process the primary one.
	targetVideo := dto.VideoURL
	caption := ""
	language := "US"
	if len(dto.EdgeMediaToCaption.Edges) > 0 {
		caption = dto.EdgeMediaToCaption.Edges[0].Node.Text
		lang, ok := detector.DetectLanguageOf(dto.EdgeMediaToCaption.Edges[0].Node.Text)
		if ok {
			language = lang.IsoCode639_1().String()
		}
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

	er := float64(dto.EdgeMediaPreviewLike.Count+dto.EdgeMediaToParentComment.Count) / float64(dto.Owner.EdgeFollowedBy.Count)
	videoStats := map[string]float64{
		"like_count":      float64(dto.EdgeMediaPreviewLike.Count),
		"comment_count":   float64(dto.EdgeMediaToParentComment.Count),
		"play_count":      float64(dto.VideoPlayCount),
		"engagement_rate": er * 100,
	}

	totalLike := int64(0)
	totalComment := int64(0)
	totalPlayCount := int64(0)
	for _, reel := range otherReelsDto.Items {
		totalLike += reel.Media.LikeCount
		totalComment += reel.Media.CommentCount
		totalPlayCount += reel.Media.PlayCount
	}

	averageStats := map[string]float64{}
	if len(otherReelsDto.Items) > 0 && totalLike > 0 {
		averageStats = map[string]float64{
			"follower_count":          float64(dto.Owner.EdgeFollowedBy.Count),
			"average_like_count":      float64(totalLike) / float64(len(otherReelsDto.Items)),
			"average_comment_count":   float64(totalComment) / float64(len(otherReelsDto.Items)),
			"average_play_count":      float64(totalPlayCount) / float64(len(otherReelsDto.Items)),
			"average_engagement_rate": (float64(totalLike+totalComment) / float64(int64(len(otherReelsDto.Items))*dto.Owner.EdgeFollowedBy.Count)) * 100,
		}
	}

	logger.Infof("Starting AI analysis for video: %s", targetVideo)

	// 3. Call the Gemini client using the direct URL
	// We pass "Instagram" as the platform context
	llmRequest, llmResponse, result, err := s.GeminiClient.AnalyzeVideo(
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
		return llmRequest, llmResponse, nil, errs.Newf(errs.Internal, err, "failed to analyze video content")
	}

	logger.Infof("Finished analysis for video: %s", targetVideo)

	return llmRequest, llmResponse, result, nil
}

func getViralScore(topicScore, hookScore, pacingScore, valueDeliveryScore, shareabilityScore, ctaScore, scopeConfidence int, scope string) float64 {
	t := float64(topicScore)
	h := float64(hookScore)
	p := float64(pacingScore)
	v := float64(valueDeliveryScore)
	s := float64(shareabilityScore)
	c := float64(ctaScore)

	gateMultiplier := 1.0
	if t < 60 || s < 60 {
		gateMultiplier = 0.6
	}

	scopeMultiplier := 1.0

	if scopeConfidence >= 70 {
		switch scope {
		case "Local":
			scopeMultiplier = 0.75
		case "National":
			scopeMultiplier = 0.9
		case "Global":
			scopeMultiplier = 1.0
		}
	}

	// Base weighted score
	score := (h*0.25 + t*0.2 + p*0.15 + v*0.15 + s*0.15 + c*0.1) * gateMultiplier * scopeMultiplier

	if c > 90 && s < 70 {
		score *= 0.85
	}

	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}
