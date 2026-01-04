package utils

import "regexp"

func GetInstagramShortcode(urlStr string) string {
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
