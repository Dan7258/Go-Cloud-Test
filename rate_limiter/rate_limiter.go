// Package rateLimiter предоставляет функционал для ограничения частоты запросов с использованием алгоритма Token Bucket.
package rateLimiter

import (
	"cloud/logger"
	"cloud/models"
	"encoding/json"
	"time"
)

// CallClient уменьшает количество доступных токенов на 1 и обновляет время последнего запроса.
// Возвращает true, если запрос успешен (токены есть), иначе false.
func (tb *TokenBucket) CallClient() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	if tb.CurrentTokenCount == 0 {
		return false
	}
	tb.CurrentTokenCount--
	tb.LastCall = time.Now()
	return true
}

// CreateNewClient создает нового клиента с переданным clientID и инициализирует его данные.
// Добавляет нового клиента в базу данных.
func (tb *TokenBucket) CreateNewClient(clientID string) error {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.RateLimits.ClientID = clientID
	tb.RateLimits.RatePerSec = tokenBucketConfig.RatePerSec
	tb.RateLimits.Capacity = tokenBucketConfig.Capacity
	tb.CurrentTokenCount = tokenBucketConfig.Capacity
	tb.LastCall = time.Now()
	err := models.CreateClient(tb.RateLimits)
	if err != nil {
		logger.PrintError("Ошибка при установке данных в БД: " + err.Error())
		return err
	}
	logger.PrintInfo("Новый клиент добавлен в БД")
	return nil
}

// GetClientDataFromDB загружает данные клиента из базы данных по его clientID.
// Обновляет количество токенов до максимальной емкости и устанавливает время последнего запроса.
func (tb *TokenBucket) GetClientDataFromDB(clientID string) error {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	var err error
	tb.RateLimits, err = models.GetClient(clientID)
	tb.CurrentTokenCount = tb.RateLimits.Capacity
	tb.LastCall = time.Now()
	if err != nil {
		logger.PrintError("Ошибка при получении данных из БД: " + err.Error())
		return err
	}
	logger.PrintInfo("Данные клиента получены из БД")
	return nil
}

// GetClientDataFromRedis загружает данные клиента из Redis.
// Если данные не найдены или произошла ошибка, возвращает false.
func (tb *TokenBucket) GetClientDataFromRedis(clientID string) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	data, err := models.GetDataFromRedis(clientID)
	if err != nil || data == nil {
		logger.PrintWarning("Ошибка при получении данных из Redis")
		return false
	}
	json.Unmarshal(data, tb)
	logger.PrintInfo("Данные клиента получены из Redis")
	return true
}

// SetClientDataInRedis сохраняет данные клиента в Redis.
// Возвращает false в случае ошибки.
func (tb *TokenBucket) SetClientDataInRedis() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	data, err := json.Marshal(tb)
	if err != nil {
		logger.PrintError("Ошибка при установке данных в Redis: " + err.Error())
		return false
	}
	err = models.SetDataInRedis(tb.RateLimits.ClientID, data, time.Hour)
	if err != nil {
		return false
	}
	logger.PrintInfo("Данные клиента установлены в Redis")
	return true
}

// UpdateClientDataByKeysInRedis обновляет данные всех клиентов в Redis, добавляя токены, основываясь на времени последнего запроса.
func UpdateClientDataByKeysInRedis() {
	keys, err := models.GetAllKeysFromRedis()
	if err != nil {
		return
	}
	for _, key := range keys {
		tb := new(TokenBucket)
		tb.GetClientDataFromRedis(key)
		diff := int(time.Now().Sub(tb.LastCall).Seconds())
		tb.CurrentTokenCount = (diff*tb.RateLimits.RatePerSec + tb.CurrentTokenCount)
		if tb.CurrentTokenCount > tb.RateLimits.Capacity {
			tb.CurrentTokenCount = tb.RateLimits.Capacity
		}
		tb.LastCall = time.Now()
		tb.SetClientDataInRedis()
	}
	logger.PrintInfo("Пополнили токены")
}

// StartTokenTicker запускает тикер, который обновляет данные клиентов в Redis каждые 30 секунд.
func StartTokenTicker() {
	tiker := time.NewTicker(30 * time.Second)
	defer tiker.Stop()
	for {
		select {
		case <-tiker.C:
			UpdateClientDataByKeysInRedis()
		}
	}
}
