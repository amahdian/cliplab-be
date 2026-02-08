package svc

import (
	"context"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/google/uuid"
)

type CreditSvc interface {
	CheckAndDeduct(userId uuid.UUID, ruleKey model.Tool) (int, error)
	GetBalance(userId uuid.UUID) (int, error)
	AddCredits(userId uuid.UUID, amount int) error
}

type creditSvc struct {
	ctx context.Context
	stg storage.PgStorage
}

func newCreditSvc(ctx context.Context, stg storage.PgStorage) CreditSvc {
	return &creditSvc{
		ctx: ctx,
		stg: stg,
	}
}

func (s *creditSvc) CheckAndDeduct(userId uuid.UUID, ruleKey model.Tool) (int, error) {
	rule := model.GetCreditRule(ruleKey)
	if rule == nil {
		return 0, errs.Newf(errs.InvalidArgument, nil, "invalid credit rule key: %s", ruleKey)
	}

	user, err := s.stg.User(s.ctx).FindById(userId)
	if err != nil {
		return 0, errs.Wrapf(err, "failed to find user")
	}
	if user == nil {
		return 0, errs.Newf(errs.NotFound, nil, "user not found")
	}

	// If cost is 0, just return current balance
	if rule.Amount <= 0 {
		return user.Credits, nil
	}

	if user.Credits < rule.Amount {
		return user.Credits, errs.Newf(errs.PermissionDenied, nil, "insufficient credits")
	}

	// Deduct and update
	user.Credits -= rule.Amount
	if err := s.stg.User(s.ctx).UpdateOne(user, false); err != nil {
		return user.Credits + rule.Amount, errs.Wrapf(err, "failed to deduct credits")
	}

	return user.Credits, nil
}

func (s *creditSvc) GetBalance(userId uuid.UUID) (int, error) {
	user, err := s.stg.User(s.ctx).FindById(userId)
	if err != nil {
		return 0, errs.Wrapf(err, "failed to find user")
	}
	if user == nil {
		return 0, errs.Newf(errs.NotFound, nil, "user not found")
	}
	return user.Credits, nil
}

func (s *creditSvc) AddCredits(userId uuid.UUID, amount int) error {
	user, err := s.stg.User(s.ctx).FindById(userId)
	if err != nil {
		return errs.Wrapf(err, "failed to find user")
	}
	if user == nil {
		return errs.Newf(errs.NotFound, nil, "user not found")
	}

	user.Credits += amount
	if err := s.stg.User(s.ctx).UpdateOne(user, false); err != nil {
		return errs.Wrapf(err, "failed to add credits")
	}

	return nil
}
