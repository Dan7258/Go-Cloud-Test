package loadBalancer

import (
	"cloud/logger"
	"cloud/models"
	"cloud/rate_limiter"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (lb *LoadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	currentHost := lb.GetNextBackend()
	counter := 0
	for ; counter < len(lb.backends); counter++ {
		currentHost = lb.GetNextBackend()
		if lb.Ping(currentHost) {
			counter = len(lb.backends)
			continue
		}
		currentHost = nil
	}
	if currentHost != nil {
		tokenBucket := new(rateLimiter.TokenBucket)
		ok := tokenBucket.GetClientDataFromRedis(r.RemoteAddr)
		if !ok {
			var err error
			if models.ThsClientExists(r.RemoteAddr) {
				err = tokenBucket.GetClientDataFromDB(r.RemoteAddr)
			} else {
				err = tokenBucket.CreateNewClient(r.RemoteAddr)
			}
			if err != nil {
				logger.PrintWarning(err.Error())
			}
		}
		ok = tokenBucket.CallClient()
		if ok {
			proxy := httputil.NewSingleHostReverseProxy(currentHost)
			proxy.ServeHTTP(w, r)
		} else {
			logger.SendError(w, http.StatusTooManyRequests, "Куда летим? Слишком много запросов!")
		}
		tokenBucket.SetClientDataInRedis()

	} else {
		logger.SendError(w, http.StatusServiceUnavailable, "Серверам плохо, попробуйте позже")
		logger.PrintError("Нет доступных серверов")
	}
}

func (lb *LoadBalancer) GetNextBackend() *url.URL {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.currentId = (lb.currentId + 1) % uint64(len(lb.backends))
	return lb.backends[lb.currentId]
}

func (lb *LoadBalancer) Ping(backend *url.URL) bool {
	timeout := time.Second * 3
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(backend.String())
	logger.PrintInfo("Проверка состояния сервера: " + backend.String())
	if err != nil {
		msg := fmt.Sprintf("Сервер %s не отвечает: %v", backend, err)
		logger.PrintWarning(msg)
		return false
	}
	resp.Body.Close()
	msg := fmt.Sprintf("Соединение с сервером %s установлено", backend)
	logger.PrintInfo(msg)
	return true
}
