package utils

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/amahdian/cliplab-be/domain/model"
)

func DetectSocialMediaID(u url.URL) model.SocialPlatform {
	text := strings.TrimSpace(u.String())
	text = strings.Split(text, "?")[0]
	text = strings.TrimSuffix(text, "/")

	// Post Patterns
	youtubeRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([\w-]+)`)
	instagramRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?instagram\.com/(?:reels?|reel|p|tv)/([A-Za-z0-9_-]+)`)
	tiktokRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?tiktok\.com/@[\w.-]+/video/(\d+)`)
	twitterRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:twitter\.com|x\.com)/\w+/status/(\d+)`)

	// Profile Patterns
	youtubeProfileRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?youtube\.com/(?:@|c/|channel/)?([\w.-]+)`)
	instagramProfileRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?instagram\.com/([\w.-]+)`)
	tiktokProfileRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?tiktok\.com/@([\w.-]+)`)
	twitterProfileRegex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:twitter\.com|x\.com)/([\w.-]+)`)

	if youtubeRegex.MatchString(text) || youtubeProfileRegex.MatchString(text) {
		return model.PlatformYouTube
	}
	if instagramRegex.MatchString(text) || instagramProfileRegex.MatchString(text) {
		return model.PlatformInstagram
	}
	if tiktokRegex.MatchString(text) || tiktokProfileRegex.MatchString(text) {
		return model.PlatformTikTok
	}
	if twitterRegex.MatchString(text) || twitterProfileRegex.MatchString(text) {
		return model.PlatformTwitter
	}

	return model.PlatformUnknown
}
