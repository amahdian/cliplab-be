package req

import "github.com/amahdian/cliplab-be/domain/model"

type EngagementRateRequest struct {
	URL      string               `form:"url" binding:"required"`
	Platform model.SocialPlatform `form:"platform"`
}
