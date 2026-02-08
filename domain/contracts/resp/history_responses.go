package resp

import (
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/google/uuid"
)

type UserHistoryResponse struct {
	ID               uuid.UUID            `json:"id"`
	ToolName         model.Tool           `json:"toolName"`
	ToolDisplayName  string               `json:"toolDisplayName"`
	CreditsUsed      int                  `json:"creditsUsed"`
	RemainingCredits int                  `json:"remainingCredits"`
	Status           model.UsageStatus    `json:"status"`
	Details          string               `json:"details"`
	Platform         model.SocialPlatform `json:"platform"`
	CreatedAt        time.Time            `json:"createdAt"`
}

func MapUsageHistoryToResponse(h *model.UsageHistory) *UserHistoryResponse {
	return &UserHistoryResponse{
		ID:               h.ID,
		ToolName:         h.ToolName,
		ToolDisplayName:  h.ToolName.DisplayName(),
		CreditsUsed:      h.CreditsUsed,
		RemainingCredits: h.RemainingCredits,
		Status:           h.Status,
		Details:          h.Details,
		Platform:         h.Platform,
		CreatedAt:        h.CreatedAt,
	}
}

func MapUsageHistoriesToResponse(histories []*model.UsageHistory) []*UserHistoryResponse {
	res := make([]*UserHistoryResponse, len(histories))
	for i, h := range histories {
		res[i] = MapUsageHistoryToResponse(h)
	}
	return res
}
