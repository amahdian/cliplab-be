package svc

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/clients/scrapecreators"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/amahdian/cliplab-be/svc/utils"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type ChannelSvc interface {
	GetChannelEngagement(urlStr string, platform model.SocialPlatform, user *auth.UserInfo) (*resp.ChannelEngagementResponse, error)
}

type channelSvc struct {
	ctx           context.Context
	stg           storage.PgStorage
	scraperClient scrapecreators.Client
	creditSvc     CreditSvc
}

func newChannelSvc(ctx context.Context, stg storage.PgStorage, scraperClient scrapecreators.Client, creditSvc CreditSvc) ChannelSvc {
	return &channelSvc{
		ctx:           ctx,
		stg:           stg,
		scraperClient: scraperClient,
		creditSvc:     creditSvc,
	}
}

func (s *channelSvc) GetChannelEngagement(urlStr string, platform model.SocialPlatform, user *auth.UserInfo) (*resp.ChannelEngagementResponse, error) {
	var err error
	if platform == "" || platform == model.PlatformUnknown {
		u, err := url.Parse(urlStr)
		if err == nil {
			platform = utils.DetectSocialMediaID(*u)
		}
	}

	if platform == model.PlatformUnknown {
		return nil, errs.Newf(errs.InvalidArgument, nil, "unknown social media platform")
	}

	handle := s.extractHandle(urlStr, platform)
	if handle == "" {
		return nil, errs.Newf(errs.InvalidArgument, nil, "could not extract handle from url")
	}

	var remainingCredits *int
	canSeeBreakdown := false
	if user != nil {
		if credits, err := s.creditSvc.CheckAndDeduct(user.Id, model.CreditKeyEngagementBreakdown); err == nil {
			remainingCredits = &credits
			canSeeBreakdown = true
		} else {
			// If deduction fails (e.g. no credits), still try to get current balance for display
			if credits, err := s.creditSvc.GetBalance(user.Id); err == nil {
				remainingCredits = &credits
			}
		}
	}

	// Check DB for existing channel data
	channel, _ := s.stg.Channel(s.ctx).FindByHandler(handle)
	if channel != nil && channel.LastHistory != nil && time.Since(channel.LastHistory.CreatedAt) < 2*24*time.Hour {
		posts, _ := s.stg.Post(s.ctx).ListBy("channel_id = ?", channel.ID)
		return &resp.ChannelEngagementResponse{
			ProfileImage:          channel.LastHistory.ProfileImage,
			ProfileUrl:            s.generateProfileUrl(handle, platform),
			ProfileDescriptor:     channel.LastHistory.ProfileDescriptor,
			FollowersCount:        channel.LastHistory.FollowersCount,
			FollowingCount:        channel.LastHistory.FollowingCount,
			PostsCount:            channel.LastHistory.MediaCount,
			AverageLikeCount:      channel.LastHistory.AverageLikes,
			AverageCommentCount:   channel.LastHistory.AverageComments,
			AveragePlayCount:      channel.LastHistory.AverageVideoPlays,
			AverageEngagementRate: channel.LastHistory.AverageEngagementRate,
			LatestPostPublishDate: channel.LastHistory.LatestPostPublishDate,
			Breakdown: func() []*resp.PostBreakdown {
				if canSeeBreakdown {
					return s.mapPostsToBreakdown(posts, channel.LastHistory.FollowersCount)
				}
				return []*resp.PostBreakdown{}
			}(),
			Credits: remainingCredits,
		}, nil
	}

	var result *resp.ChannelEngagementResponse
	var posts []*model.Post
	switch platform {
	case model.PlatformInstagram:
		posts, result, err = s.getInstagramEngagement(handle)
	case model.PlatformTikTok:
		posts, result, err = s.getTikTokEngagement(handle)
	case model.PlatformTwitter:
		posts, result, err = s.getTwitterEngagement(handle)
	default:
		return nil, errs.Newf(errs.InvalidArgument, nil, "engagement rate calculation not supported for this platform yet")
	}

	if err != nil {
		return nil, err
	}

	// Store data in DB
	err = s.storeChannelData(handle, platform, result, posts)
	if err != nil {
		// Log error but return result anyway? Usually yes for caching.
		fmt.Printf("failed to store channel data: %v\n", err)
	}

	if canSeeBreakdown {
		result.Breakdown = s.mapPostsToBreakdown(posts, result.FollowersCount)
	} else {
		result.Breakdown = []*resp.PostBreakdown{}
	}

	result.Credits = remainingCredits

	return result, nil
}

func (s *channelSvc) storeChannelData(handle string, platform model.SocialPlatform, data *resp.ChannelEngagementResponse, posts []*model.Post) error {
	channel, err := s.stg.Channel(s.ctx).FindByHandler(handle)
	if err != nil {
		return err
	}

	if channel == nil {
		channel = &model.Channel{
			FullName:          handle, // Fallback if name not in data
			Handler:           handle,
			Platform:          platform,
			ProfileImage:      data.ProfileImage,
			ProfileDescriptor: data.ProfileDescriptor,
		}
		err = s.stg.Channel(s.ctx).CreateOne(channel)
		if err != nil {
			return err
		}
	} else {
		channel.ProfileImage = data.ProfileImage
		channel.ProfileDescriptor = data.ProfileDescriptor
		err = s.stg.Channel(s.ctx).UpdateOne(channel, false)
		if err != nil {
			return err
		}
	}

	history := &model.ChannelHistory{
		ChannelID:             channel.ID,
		FollowersCount:        data.FollowersCount,
		FollowingCount:        data.FollowingCount,
		MediaCount:            data.PostsCount,
		ProfileImage:          data.ProfileImage,
		ProfileDescriptor:     data.ProfileDescriptor,
		AverageLikes:          data.AverageLikeCount,
		AverageComments:       data.AverageCommentCount,
		AverageVideoPlays:     data.AveragePlayCount,
		AverageEngagementRate: data.AverageEngagementRate,
		LatestPostPublishDate: data.LatestPostPublishDate,
	}

	err = s.stg.ChannelHistory(s.ctx).CreateOne(history)
	if err != nil {
		return err
	}

	// Store posts
	if len(posts) > 0 {
		for _, p := range posts {
			p.ChannelId = &channel.ID
		}
		err = s.stg.Post(s.ctx).UpsertMany(posts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *channelSvc) generateProfileUrl(handle string, platform model.SocialPlatform) string {
	switch platform {
	case model.PlatformInstagram:
		return fmt.Sprintf("https://www.instagram.com/%s", handle)
	case model.PlatformTikTok:
		return fmt.Sprintf("https://www.tiktok.com/@%s", handle)
	case model.PlatformTwitter:
		return fmt.Sprintf("https://x.com/%s", handle)
	}
	return ""
}

func (s *channelSvc) getInstagramEngagement(handle string) ([]*model.Post, *resp.ChannelEngagementResponse, error) {
	var profile *scrapecreators.InstagramProfileResponse
	var instaPosts *scrapecreators.PostsResponse

	g, _ := errgroup.WithContext(s.ctx)
	g.Go(func() (err error) {
		profile, err = s.scraperClient.GetInstagramProfile(handle)
		return
	})
	g.Go(func() (err error) {
		instaPosts, err = s.scraperClient.GetInstagramPagePosts(handle)
		return
	})
	if err := g.Wait(); err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch instagram reels")
	}

	if instaPosts != nil && len(instaPosts.Items) == 0 {
		return nil, nil, errs.Newf(errs.NotFound, nil, "no posts found for this user")
	}

	owner := profile.Data.User
	var posts []*model.Post
	var totalLikes, totalComments, playCounts, videoCount int64
	for _, item := range instaPosts.Items {
		totalLikes += item.LikeCount
		totalComments += item.CommentCount
		playCounts += item.PlayCount
		if item.PlayCount != 0 {
			videoCount++
		}

		format := model.PostFormatImage
		if len(item.VideoVersions) > 0 {
			format = model.PostFormatVideo
		}

		post := &model.Post{
			ID:               item.Code,
			Link:             fmt.Sprintf("https://www.instagram.com/p/%s/", item.Code),
			ImageURL:         &item.DisplayUri,
			Format:           format,
			UserName:         owner.Username,
			UserAnchor:       owner.FullName,
			UserProfileLink:  fmt.Sprintf("https://www.instagram.com/%s/", owner.Username),
			UserProfileImage: owner.ProfilePicUrlHd,
			LikeCount:        item.LikeCount,
			CommentCount:     item.CommentCount,
			VideoViewCount:   item.PlayCount,
			PostDate:         time.Unix(item.TakenAt, 0),
		}
		if len(item.VideoVersions) > 0 {
			post.VideoURL = &item.VideoVersions[0].Url
		}

		posts = append(posts, post)
	}

	count := int64(len(profile.Data.User.EdgeOwnerToTimelineMedia.Edges))
	avgLikes := float64(totalLikes) / float64(count)
	avgComments := float64(totalComments) / float64(count)
	avgPlays := float64(0)
	if videoCount > 0 {
		avgPlays = float64(playCounts) / float64(videoCount)
	}

	followers := owner.EdgeFollowedBy.Count
	var avgER float64
	if followers > 0 {
		avgER = ((avgLikes + avgComments) / float64(followers)) * 100
	}

	firstItem := profile.Data.User.EdgeOwnerToTimelineMedia.Edges[0]
	return posts, &resp.ChannelEngagementResponse{
		ProfileImage:          owner.ProfilePicUrlHd,
		ProfileDescriptor:     owner.Biography,
		ProfileUrl:            fmt.Sprintf("https://www.instagram.com/%s", handle),
		FollowersCount:        followers,
		FollowingCount:        owner.EdgeFollow.Count,
		PostsCount:            owner.EdgeOwnerToTimelineMedia.Count,
		AverageLikeCount:      avgLikes,
		AverageCommentCount:   avgComments,
		AveragePlayCount:      avgPlays,
		AverageEngagementRate: avgER,
		LatestPostPublishDate: time.Unix(firstItem.Node.TakenAtTimestamp, 0),
	}, nil
}

func (s *channelSvc) getTikTokEngagement(handle string) ([]*model.Post, *resp.ChannelEngagementResponse, error) {
	res, err := s.scraperClient.GetTikTokProfileVideos(handle)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch tiktok videos")
	}

	if len(res.AwemeList) == 0 {
		return nil, nil, errs.Newf(errs.NotFound, nil, "no videos found for this user")
	}

	first := res.AwemeList[0]
	var posts []*model.Post
	var totalLikes, totalComments, totalPlays int64
	for _, item := range res.AwemeList {
		totalLikes += item.Statistics.DiggCount
		totalComments += item.Statistics.CommentCount
		totalPlays += item.Statistics.PlayCount

		coverImg := ""
		if len(item.Video.Cover.URLList) > 0 {
			coverImg = item.Video.Cover.URLList[0]
		}
		videoUrl := ""
		if len(item.Video.PlayAddr.URLList) > 0 {
			videoUrl = item.Video.PlayAddr.URLList[0]
		}

		profileImg := ""
		if len(item.Author.AvatarLarger.URLList) > 0 {
			profileImg = item.Author.AvatarLarger.URLList[0]
		}

		posts = append(posts, &model.Post{
			ID:               item.AwemeID,
			Link:             item.URL,
			ImageURL:         &coverImg,
			VideoURL:         &videoUrl,
			Format:           model.PostFormatVideo,
			UserName:         item.Author.UniqueID,
			UserAnchor:       item.Author.Nickname,
			UserProfileLink:  fmt.Sprintf("https://www.tiktok.com/@%s", item.Author.UniqueID),
			UserProfileImage: profileImg,
			LikeCount:        item.Statistics.DiggCount,
			CommentCount:     item.Statistics.CommentCount,
			VideoPlayCount:   item.Statistics.PlayCount,
			PostDate:         time.Unix(item.CreateTime, 0),
		})
	}

	count := int64(len(res.AwemeList))
	avgLikes := float64(totalLikes) / float64(count)
	avgComments := float64(totalComments) / float64(count)
	avgPlays := float64(totalPlays) / float64(count)

	followers := first.Author.FollowerCount
	var avgER float64
	if followers > 0 {
		avgER = ((avgLikes + avgComments) / float64(followers)) * 100
	}

	profileImg := ""
	if len(first.Author.AvatarLarger.URLList) > 0 {
		profileImg = first.Author.AvatarLarger.URLList[0]
	}

	return posts, &resp.ChannelEngagementResponse{
		ProfileImage:          profileImg,
		ProfileUrl:            fmt.Sprintf("https://www.tiktok.com/@%s", handle),
		FollowersCount:        followers,
		FollowingCount:        first.Author.FollowingCount,
		PostsCount:            int64(first.Author.AwemeCount),
		AverageLikeCount:      avgLikes,
		AverageCommentCount:   avgComments,
		AveragePlayCount:      avgPlays,
		AverageEngagementRate: avgER,
		LatestPostPublishDate: time.Unix(first.CreateTime, 0),
	}, nil
}

func (s *channelSvc) getTwitterEngagement(handle string) ([]*model.Post, *resp.ChannelEngagementResponse, error) {
	res, err := s.scraperClient.GetUserTweets(handle)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch twitter tweets")
	}

	if len(res.Tweets) == 0 {
		return nil, nil, errs.Newf(errs.NotFound, nil, "no tweets found for this user")
	}

	first := res.Tweets[0]
	var posts []*model.Post
	var totalLikes, totalReplies, totalRetweets int64
	for _, item := range res.Tweets {
		totalLikes += int64(item.Legacy.FavoriteCount)
		totalReplies += int64(item.Legacy.ReplyCount)
		totalRetweets += int64(item.Legacy.RetweetCount)

		userResult := item.Core.UserResults.Result
		var imageURL, videoURL *string
		format := model.PostFormatText
		if len(item.Legacy.ExtendedEntities.Media) > 0 {
			m := item.Legacy.ExtendedEntities.Media[0]
			imageURL = &m.MediaURLHttps
			if m.Type == "video" || m.Type == "animated_gif" {
				format = model.PostFormatVideo
				// Video URL is usually deeper in media[0].video_info.variants, but not in DTO yet.
			} else {
				format = model.PostFormatImage
			}
		}

		postDate, _ := s.parseTwitterDate(item.Legacy.CreatedAt)

		posts = append(posts, &model.Post{
			ID:               item.RestID,
			Link:             item.URL,
			ImageURL:         imageURL,
			VideoURL:         videoURL,
			Format:           format,
			UserName:         userResult.Core.ScreenName,
			UserAnchor:       userResult.Core.Name,
			UserProfileLink:  fmt.Sprintf("https://x.com/%s", userResult.Core.ScreenName),
			UserProfileImage: userResult.Avatar.ImageURL,
			LikeCount:        int64(item.Legacy.FavoriteCount),
			CommentCount:     int64(item.Legacy.ReplyCount),
			VideoViewCount:   0, // Needs string conversion from item.Views.Count if possible
			PostDate:         postDate,
		})
	}

	count := int64(len(res.Tweets))
	avgLikes := float64(totalLikes) / float64(count)
	avgComments := float64(totalReplies) / float64(count)
	avgRetweets := float64(totalRetweets) / float64(count)

	userResult := first.Core.UserResults.Result
	followers := int64(userResult.Legacy.FollowersCount)
	var avgER float64
	if followers > 0 {
		avgER = (avgLikes + avgComments + avgRetweets/float64(followers)) * 100
	}

	latestDate, _ := s.parseTwitterDate(first.Legacy.CreatedAt)

	return posts, &resp.ChannelEngagementResponse{
		ProfileImage:          userResult.Avatar.ImageURL,
		ProfileUrl:            fmt.Sprintf("https://x.com/%s", handle),
		ProfileDescriptor:     userResult.Legacy.Description,
		FollowersCount:        followers,
		FollowingCount:        int64(userResult.Legacy.FriendsCount),
		PostsCount:            int64(userResult.Legacy.StatusesCount),
		AverageLikeCount:      avgLikes,
		AverageCommentCount:   avgComments,
		AveragePlayCount:      0,
		AverageShareCount:     avgRetweets,
		AverageEngagementRate: avgER,
		LatestPostPublishDate: latestDate,
	}, nil
}

func (s *channelSvc) extractHandle(urlStr string, platform model.SocialPlatform) string {
	text := strings.TrimSpace(urlStr)
	text = strings.Split(text, "?")[0]
	text = strings.TrimSuffix(text, "/")

	// If it doesn't look like a URL but we have a platform, treat it as a direct handle
	if !strings.Contains(text, ".") && !strings.Contains(text, "/") {
		return strings.TrimPrefix(text, "@")
	}

	var re *regexp.Regexp
	switch platform {
	case model.PlatformYouTube:
		re = regexp.MustCompile(`youtube\.com/(?:@|c/|channel/)?([\w.-]+)`)
	case model.PlatformInstagram:
		if strings.Contains(text, "/reels/") || strings.Contains(text, "/reel/") || strings.Contains(text, "/p/") || strings.Contains(text, "/tv/") {
			return ""
		}
		re = regexp.MustCompile(`instagram\.com/([\w.-]+)`)
	case model.PlatformTikTok:
		re = regexp.MustCompile(`tiktok\.com/@([\w.-]+)`)
	case model.PlatformTwitter:
		re = regexp.MustCompile(`(?:twitter\.com|x\.com)/([\w.-]+)`)
	}

	if re != nil {
		if re.MatchString(text) {
			match := re.FindStringSubmatch(text)
			if len(match) > 1 {
				if platform == model.PlatformTwitter && strings.Contains(match[1], "/status/") {
					return strings.Split(match[1], "/")[0]
				}
				return match[1]
			}
		}
	}

	// Fallback: if it's a simple word and platform is known, it's likely a handle
	if platform != model.PlatformUnknown {
		handleRegex := regexp.MustCompile(`^@?[\w.-]+$`)
		if handleRegex.MatchString(text) {
			return strings.TrimPrefix(text, "@")
		}
	}

	return ""
}

func (s *channelSvc) parseTwitterDate(createdAt string) (time.Time, error) {
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	return time.Parse(layout, createdAt)
}

func (s *channelSvc) mapPostsToBreakdown(posts []*model.Post, followers int64) []*resp.PostBreakdown {
	var breakdown []*resp.PostBreakdown
	for _, p := range posts {
		var er float64
		if followers > 0 {
			er = (float64(p.LikeCount+p.CommentCount) / float64(followers)) * 100
		}
		img := ""
		if p.ImageURL != nil {
			img = *p.ImageURL
		}
		breakdown = append(breakdown, &resp.PostBreakdown{
			Image:          img,
			Link:           p.Link,
			PublishDate:    p.PostDate,
			LikeCount:      p.LikeCount,
			CommentCount:   p.CommentCount,
			ViewCount:      p.VideoViewCount,
			PlayCount:      p.VideoPlayCount,
			EngagementRate: er,
		})
	}
	return breakdown
}
