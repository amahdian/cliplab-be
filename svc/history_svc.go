package svc

import (
	"context"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/domain/model/common"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/google/uuid"
)

type HistorySvc interface {
	StoreUsage(userId uuid.UUID, toolName model.Tool, creditsUsed int, remainingCredits int, details string, platform model.SocialPlatform, usageType string, status model.UsageStatus, requestID *uuid.UUID, responseLog string) error
	UpdateUsageByRequestID(requestID uuid.UUID, status model.UsageStatus, responseLog string) error
	GetUserHistory(userId uuid.UUID, pagination *common.Pagination) ([]*model.UsageHistory, error)
}

type historySvc struct {
	ctx context.Context
	stg storage.PgStorage
}

func newHistorySvc(ctx context.Context, stg storage.PgStorage) HistorySvc {
	return &historySvc{
		ctx: ctx,
		stg: stg,
	}
}

func (s *historySvc) StoreUsage(
	userId uuid.UUID,
	toolName model.Tool,
	creditsUsed int,
	remainingCredits int,
	details string,
	platform model.SocialPlatform,
	usageType string,
	status model.UsageStatus,
	requestID *uuid.UUID,
	responseLog string,
) error {
	history := &model.UsageHistory{
		UserId:           userId,
		ToolName:         toolName,
		CreditsUsed:      creditsUsed,
		RemainingCredits: remainingCredits,
		Status:           status,
		Details:          details,
		Platform:         platform,
		Type:             usageType,
		RequestID:        requestID,
		ResponseLog:      responseLog,
	}

	if err := s.stg.UsageHistory(s.ctx).CreateOne(history); err != nil {
		return errs.Wrapf(err, "failed to store usage history")
	}

	return nil
}

func (s *historySvc) UpdateUsageByRequestID(requestID uuid.UUID, status model.UsageStatus, responseLog string) error {
	histories, err := s.stg.UsageHistory(s.ctx).ListBy("request_id = ?", requestID)
	if err != nil {
		return errs.Wrapf(err, "failed to find usage history by request id")
	}

	if len(histories) == 0 {
		return nil // Or return error? If it's a guest or skipped, maybe fine.
	}

	history := histories[0]
	history.Status = status
	history.ResponseLog = responseLog

	if err := s.stg.UsageHistory(s.ctx).UpdateOne(history, false); err != nil {
		return errs.Wrapf(err, "failed to update usage history")
	}

	return nil
}

func (s *historySvc) GetUserHistory(userId uuid.UUID, pagination *common.Pagination) ([]*model.UsageHistory, error) {
	histories, err := s.stg.UsageHistory(s.ctx).FindByUserId(userId, pagination)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to fetch user history")
	}
	return histories, nil
}
