package svc

import (
	"context"

	"github.com/amahdian/cliplab-be/clients/gemini"
	"github.com/amahdian/cliplab-be/clients/scrapecreators"
	"github.com/redis/go-redis/v9"

	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/aws/aws-sdk-go/aws"
)

type Svc interface {
	NewUserSvc(ctx context.Context) UserSvc
	NewPostSvc(ctx context.Context) AnalyzeSvc
	NewFileSvc(ctx context.Context) FileSvc
	NewQueueSvc(ctx context.Context) QueueSvc
	NewWebSocketSvc(ctx context.Context) WebSocketSvc
	NewChannelSvc(ctx context.Context) ChannelSvc
	NewCreditSvc(ctx context.Context) CreditSvc
	NewBillingSvc(ctx context.Context) BillingSvc
	NewHistorySvc(ctx context.Context) HistorySvc
}

type StorageConfig struct {
	AwsConfig   *aws.Config
	ProductName string
}

type svcImpl struct {
	pgStg         storage.PgStorage
	Envs          *env.Envs
	geminiClient  gemini.Client
	scraperClient scrapecreators.Client
	redisClient   *redis.Client
	storageConfig StorageConfig
}

func NewSvc(
	pgStg storage.PgStorage,
	envs *env.Envs,
	geminiClient gemini.Client,
	scraperClient scrapecreators.Client,
	redisClient *redis.Client,
	storageConfig StorageConfig) Svc {

	return &svcImpl{
		pgStg,
		envs,
		geminiClient,
		scraperClient,
		redisClient,
		storageConfig,
	}
}

func (s *svcImpl) NewUserSvc(ctx context.Context) UserSvc {
	return newUserSvc(ctx, s.pgStg, s.Envs)
}

func (s *svcImpl) NewFileSvc(ctx context.Context) FileSvc {
	return newFileSvc(ctx, s.storageConfig)
}

func (s *svcImpl) NewPostSvc(ctx context.Context) AnalyzeSvc {
	return newAnalyzeSvc(ctx, s.pgStg, s.Envs, s.redisClient, s.NewFileSvc(ctx), s.NewCreditSvc(ctx), s.NewHistorySvc(ctx))
}

func (s *svcImpl) NewQueueSvc(ctx context.Context) QueueSvc {
	return newPostQueueSvc(ctx, s.pgStg, s.Envs, s.geminiClient, s.scraperClient, s.NewHistorySvc(ctx), s.NewCreditSvc(ctx))
}

func (s *svcImpl) NewWebSocketSvc(ctx context.Context) WebSocketSvc {
	return NewWebSocketSvc()
}

func (s *svcImpl) NewChannelSvc(ctx context.Context) ChannelSvc {
	return newChannelSvc(ctx, s.pgStg, s.scraperClient, s.NewCreditSvc(ctx), s.NewHistorySvc(ctx))
}

func (s *svcImpl) NewCreditSvc(ctx context.Context) CreditSvc {
	return newCreditSvc(ctx, s.pgStg)
}

func (s *svcImpl) NewBillingSvc(ctx context.Context) BillingSvc {
	return newBillingSvc(ctx, s.pgStg, s.Envs, s.NewCreditSvc(ctx))
}

func (s *svcImpl) NewHistorySvc(ctx context.Context) HistorySvc {
	return newHistorySvc(ctx, s.pgStg)
}
