// Package rateLimiter предоставляет функционал для ограничения частоты запросов с использованием алгоритма Token Bucket.
package rateLimiter

import (
	"cloud/config_handler"
	"cloud/logger"
	"cloud/models"
	"sync"
	"time"
)

// TokenBucket представляет собой структуру, которая хранит информацию о токенах для клиента.
// Используется для ограничения частоты запросов с помощью алгоритма "Token Bucket".
type TokenBucket struct {
	RateLimits        models.RateLimits // Данные о лимитах для клиента.
	CurrentTokenCount int               // Текущее количество доступных токенов.
	LastCall          time.Time         // Время последнего запроса клиента.
	mutex             sync.Mutex        // Мьютекс для синхронизации доступа к данным.
}

// TokenBucketConfig хранит конфигурацию для всех клиентов в отношении лимита токенов.
type TokenBucketConfig struct {
	Capacity   int // Максимальное количество токенов, которые могут быть в корзине.
	RatePerSec int // Количество токенов, добавляемых в корзину каждую секунду.
}

// tokenBucketConfig глобальная переменная, хранящая шаблонную конфигурацию для всех клиентов.
var tokenBucketConfig TokenBucketConfig

// InitTBConfig инициализирует конфигурацию для алгоритма Rate Limiting,
// используя параметры из конфигурационного файла приложения.
// Настройка лимита и емкости корзины токенов.
func InitTBConfig(config configHandler.Config) {
	tokenBucketConfig.Capacity = config.Capacity
	tokenBucketConfig.RatePerSec = config.RatePerSec
	logger.PrintInfo("Получены данные для Rate-Limiting")
}
