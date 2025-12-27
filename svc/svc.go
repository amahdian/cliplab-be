package svc

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/amahdian/cliplab-be/clients"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/aws/aws-sdk-go/aws"
)

type Svc interface {
	NewUserSvc(ctx context.Context) UserSvc
	NewPostSvc(ctx context.Context) PostSvc
	NewFileSvc(ctx context.Context) FileSvc
	NewQueueSvc(ctx context.Context) QueueSvc
	NewWebSocketSvc(ctx context.Context) WebSocketSvc
}

type StorageConfig struct {
	AwsConfig   *aws.Config
	ProductName string
}

type svcImpl struct {
	pgStg          storage.PgStorage
	Envs           *env.Envs
	geminiClient   clients.GeminiClient
	rapidApiClient clients.RapidApiClient
	redisClient    *redis.Client
	storageConfig  StorageConfig
}

func NewSvc(
	pgStg storage.PgStorage,
	envs *env.Envs,
	geminiClient clients.GeminiClient,
	rapidApiClient clients.RapidApiClient,
	redisClient *redis.Client,
	storageConfig StorageConfig) Svc {

	return &svcImpl{
		pgStg,
		envs,
		geminiClient,
		rapidApiClient,
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

func (s *svcImpl) NewPostSvc(ctx context.Context) PostSvc {
	return newPostSvc(ctx, s.pgStg, s.Envs, s.redisClient, s.NewFileSvc(ctx))
}

func (s *svcImpl) NewQueueSvc(ctx context.Context) QueueSvc {
	return newPostQueueSvc(ctx, s.pgStg, s.Envs, s.geminiClient, s.rapidApiClient)
}

func (s *svcImpl) NewWebSocketSvc(ctx context.Context) WebSocketSvc {
	return NewWebSocketSvc()
}
