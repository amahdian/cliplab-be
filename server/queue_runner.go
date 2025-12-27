package server

import (
	"context"
	"encoding/json"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/svc"
	"github.com/redis/go-redis/v9"
)

func runQueue(rdb *redis.Client, svc svc.Svc) {
	logger.Info("Subscriber running in background...")

	for {
		res, err := rdb.BRPop(context.Background(), 0,
			global.RedisPostQueue,
		).Result()

		if err != nil {
			logger.Error("error running queue pop", err)
			continue
		}
		channel, payload := res[0], res[1]
		handleMessage(channel, payload, svc)
	}
}

func handleMessage(channel, payload string, svc svc.Svc) {
	qSvc := svc.NewQueueSvc(context.Background())

	switch channel {
	case global.RedisPostQueue:
		var data model.PostQueueData
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			logger.Error("invalid queue payload:", err)
			return
		}
		err := qSvc.ProcessPost(data.Url, data.Id, data.Platform)
		if err != nil {
			logger.Error("Failed to process post:", err)
		}
	default:
		logger.Error("Unhandled channel %s: %s", channel, payload)
	}
}
