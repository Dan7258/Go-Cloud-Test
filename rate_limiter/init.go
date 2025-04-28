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
	Capacity      int
	RatePerSecond int
}

var tokenBucketConfig TokenBucketConfig

func (tb *TokenBucket) Init(config configHandler.Config) {
	tokenBucketConfig.Capacity = config.Capacity
	tokenBucketConfig.RatePerSecond = config.RatePerSecond
	logger.PrintInfo("Получены данные для Rate-Limiting")
}
