package rateLimiter

import (
	configHandler "cloud/config_handler"
	"cloud/logger"
	"cloud/models"
	"sync"
	"time"
)

type TokenBucket struct {
	RateLimits        models.RateLimits
	CurrentTokenCount int
	LastCall          time.Time
	mutex             sync.Mutex
}

type TokenBucketConfig struct {
	Capacity   int
	RatePerSec int
}

var tokenBucketConfig TokenBucketConfig

func InitTBConfig(config configHandler.Config) {
	tokenBucketConfig.Capacity = config.Capacity
	tokenBucketConfig.RatePerSec = config.RatePerSec
	logger.PrintInfo("Получены данные для Rate-Limiting")
}
