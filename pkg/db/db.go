package db

import (
	"context"
	"log"
	"time"

	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"gorm.io/gorm"

	extraClausePlugin "github.com/WinterYukky/gorm-extra-clause-plugin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	gormLogger "gorm.io/gorm/logger"
)

type LogLevel string

const (
	LogLevelSilent LogLevel = "silent"
	LogLevelError           = "error"
	LogLevelWarn            = "warn"
	LogLevelInfo            = "info"
)

func BuildMongoConnection(dsn string, logLevel LogLevel) (*mongo.Client, error) {
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			if logLevel == LogLevelInfo {
				log.Printf("Started: %s %s\n", evt.CommandName, evt.Command)
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			if logLevel == LogLevelInfo {
				log.Printf("Succeeded: %s %s\n", evt.CommandName, evt.Reply)
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			if logLevel != LogLevelSilent {
				log.Printf("Failed: %s %s\n", evt.CommandName, evt.Failure)
			}
		},
	}
	opts := options.Client().ApplyURI(dsn).SetMonitor(monitor)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}
	logger.Info("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func OpenGormDb(dsn string, logLevel LogLevel) (*gorm.DB, error) {
	logger.Infof("trying to open connection to postgres database with dsn: %q", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel(logLevel)),
	})
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "Failed to open postgres connection: %v", err)
	}

	err = db.Use(extraClausePlugin.New())
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "Failed to use extra clause plugin: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "Failed to get underlying sql db from gorm: %v", err)
	}

	err = sqlDb.Ping()
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "Failed to ping postgres database: %v", err)
	}

	logger.Info("successfully opened database connection")
	return db, nil
}

func gormLogLevel(level LogLevel) gormLogger.LogLevel {
	switch level {
	case LogLevelSilent:
		return gormLogger.Silent
	case LogLevelError:
		return gormLogger.Error
	case LogLevelWarn:
		return gormLogger.Warn
	case LogLevelInfo:
		return gormLogger.Info
	}
	return gormLogger.Error
}
