package rateLimiter

import (
	"cloud/models"
	"encoding/json"
	"sync"
	"time"
)

type TokenBucket struct {
	RateLimits        models.RateLimits
	CurrentTokenCount int
	LastCall          time.Time
	mutex             sync.Mutex
}

func (tb *TokenBucket) CallClient() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	diff := int(time.Now().Sub(tb.LastCall).Seconds())
	tb.CurrentTokenCount = (diff*tb.RateLimits.RatePerSecond + tb.CurrentTokenCount) % tb.RateLimits.Capacity
	if tb.CurrentTokenCount == 0 {
		return false
	}
	tb.CurrentTokenCount--
	tb.LastCall = time.Now()
	return true
}

func (tb *TokenBucket) CreateNewClient(clientID string) error {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.RateLimits.ClientID = clientID
	tb.RateLimits.RatePerSecond = 2
	tb.RateLimits.Capacity = 100
	tb.CurrentTokenCount = 100
	tb.LastCall = time.Now()
	err := models.CreateClient(tb.RateLimits)
	return err
}

func (tb *TokenBucket) GetClientDataFromDB(clientID string) error {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	var err error
	tb.RateLimits, err = models.GetClient(clientID)
	tb.CurrentTokenCount = tb.RateLimits.Capacity
	tb.LastCall = time.Now()
	return err
}

func (tb *TokenBucket) GetClientDataFromRedis(clientID string) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	data, err := models.GetDataFromRedis(clientID)
	if err != nil || data == nil {
		return false
	}
	json.Unmarshal(data, tb)
	return true
}

func (tb *TokenBucket) SetClientDataInRedis() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	data, err := json.Marshal(tb)
	if err != nil {
		return false
	}
	err = models.SetDataInRedis(tb.RateLimits.ClientID, data, time.Hour)
	if err != nil {
		return false
	}
	return true
}
