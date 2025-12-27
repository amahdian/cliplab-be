package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/amahdian/cliplab-be/clients"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/pkg/db"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/server/router"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/storage/pg"
	"github.com/amahdian/cliplab-be/svc"
	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/amahdian/cliplab-be/version"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Server struct {
	Envs *env.Envs

	GeminiClient   clients.GeminiClient
	RapidApiClient clients.RapidApiClient
	RedisClient    *redis.Client

	Authenticator auth.Authenticator

	StorageConfig svc.StorageConfig
	PgStorage     storage.PgStorage
	Svc           svc.Svc
	Router        *router.Router
}

func NewServer(envs *env.Envs) (*Server, error) {
	s := &Server{
		Envs: envs,
	}
	if err := s.setupLogger(); err != nil {
		return nil, errors.Wrap(err, "Failed to initialize logger.")
	}
	if err := s.migrateDb(); err != nil {
		return nil, errors.Wrap(err, "failed to migrate the db")
	}
	if err := s.setupStorage(); err != nil {
		return nil, errors.Wrap(err, "Failed to setup storage.")
	}
	if err := s.setupInfrastructure(); err != nil {
		return nil, errors.Wrap(err, "Failed to initialize the infrastructure.")
	}
	s.setupServices()
	s.setupRouter()

	go runQueue(s.RedisClient, s.Svc)

	return s, nil
}

func (s *Server) Run() error {
	//defer func(s *Server) {
	//	err := s.Close()
	//	if err != nil {
	//		logger.Errorf("Failed to gracefully shutdown the server and release resources: %v", err)
	//	}
	//}(s)

	logger.Infof("Starting HTTP RESTFul server on port: %s", s.Envs.Server.HttpPort)

	err := s.Router.Run(fmt.Sprintf(":%s", s.Envs.Server.HttpPort))
	return err
}

func (s *Server) Close() error {
	if err := logger.Close(); err != nil {
		logger.Errorf("Failed to close/sync the logger: %v", err) // can it actually log itself?
		return err
	}
	return nil
}

func (s *Server) setupLogger() error {
	logger.ConfigureFromEnvs(s.Envs)
	return nil
}

func (s *Server) migrateDb() error {
	err := pg.EnsureDatabaseExists(s.Envs.Db.Dsn)
	if err != nil {
		return errors.Wrap(err, "failed to create database")
	}
	migrationsDir := fmt.Sprintf("file://%s/migrations", s.Envs.Server.AssetsDir)
	migrator, err := migrate.New(migrationsDir, s.Envs.Db.Dsn)
	if err != nil {
		return errors.Wrap(err, "failed to open new migration instance")
	}

	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("all migrations are already applied")
			return nil
		}
		return errors.Wrap(err, "failed to run migrations")
	}

	logger.Info("applied migrations to the db")

	return nil
}

func (s *Server) setupInfrastructure() error {
	if err := s.setupAuthenticator(); err != nil {
		return err
	}
	if err := s.setupGPTClient(); err != nil {
		return err
	}
	if err := s.setupBrightDataClient(); err != nil {
		return err
	}
	if err := s.setupRedis(); err != nil {
		return err
	}
	if err := s.setupFileStorage(); err != nil {
		return err
	}
	return nil
}

func (s *Server) setupStorage() error {
	logLevelEnv := strings.ToLower(s.Envs.Db.LogLevel)
	logLevel := db.LogLevel(logLevelEnv)
	pgDb, err := db.OpenGormDb(s.Envs.Db.Dsn, logLevel)
	if err != nil {
		return errors.Wrap(err, "Failed to open gorm connection.")
	}
	s.PgStorage = pg.NewStg(pgDb)

	return nil
}

func (s *Server) setupServices() {
	s.Svc = svc.NewSvc(
		s.PgStorage,
		s.Envs,
		s.GeminiClient,
		s.RapidApiClient,
		s.RedisClient,
		s.StorageConfig,
	)
}

func (s *Server) setupRouter() {
	s.Router = router.NewRouter(
		s.Svc,
		s.Envs,
		s.Authenticator)
}

func (s *Server) setupAuthenticator() error {
	s.Authenticator = auth.NewAuthenticator(s.Envs, s.PgStorage)
	return nil
}

func (s *Server) setupGPTClient() error {
	client := clients.NewGeminiClient(s.Envs.Gemini.ClientHost, s.Envs.Gemini.Token)
	s.GeminiClient = client
	return nil
}

func (s *Server) setupBrightDataClient() error {
	client := clients.NewRapidApiClient(s.Envs.RapidApi.Token)
	s.RapidApiClient = client
	return nil
}

func (s *Server) setupRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     s.Envs.Redis.Address,
		Password: s.Envs.Redis.Password,
		DB:       s.Envs.Redis.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return errors.Wrap(err, "failed to open redis connection")
	}

	s.RedisClient = client
	return nil
}

func (s *Server) setupFileStorage() error {
	if s.Envs.FileStorage.Bypass {
		logger.Warn("file storage is not configured")
		return nil
	}

	// store the aws config for later use
	s.StorageConfig = svc.StorageConfig{
		AwsConfig: &aws.Config{
			Credentials:      credentials.NewStaticCredentials(s.Envs.FileStorage.AccessKey, s.Envs.FileStorage.SecretKey, ""),
			Endpoint:         aws.String(s.Envs.FileStorage.Endpoint),
			Region:           aws.String(s.Envs.FileStorage.Region),
			S3ForcePathStyle: aws.Bool(true),
		},
		ProductName: version.AppName,
	}
	return nil
}
