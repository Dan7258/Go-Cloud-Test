// Package loadBalancer реализует базовый механизм балансировки нагрузки между серверами-бэкендами.
package loadBalancer

import (
	"cloud/logger"
	"cloud/models"
	"cloud/rate_limiter"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// ServeProxy обрабатывает входящие HTTP-запросы, перенаправляя их на доступный сервер-бэкенд.
// Проверяет доступность серверов, применяет ограничение по количеству запросов для клиента (rate limiting).
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
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Попытка получить данные клиента из Redis
		ok := tokenBucket.GetClientDataFromRedis(ip)
		if !ok {
			var err error
			if models.ThsClientExists(ip) {
				err = tokenBucket.GetClientDataFromDB(ip)
			} else {
				err = tokenBucket.CreateNewClient(ip)
			}
			if err != nil {
				logger.PrintWarning(err.Error())
			} else {
				tokenBucket.SetClientDataInRedis()
			}
		}

		ok = tokenBucket.CallClient()
		if ok {
			// Прокси-запрос на доступный сервер
			proxy := httputil.NewSingleHostReverseProxy(currentHost)
			proxy.ServeHTTP(w, r)
		} else {
			// Превышено ограничение запросов
			logger.SendError(w, http.StatusTooManyRequests, "Куда летим? Слишком много запросов!")
		}
		tokenBucket.SetClientDataInRedis()
	} else {
		// Нет доступных серверов
		logger.SendError(w, http.StatusServiceUnavailable, "Серверам плохо, попробуйте позже")
		logger.PrintError("Нет доступных серверов")
	}
}

// GetNextBackend возвращает следующий сервер-бэкенд для обработки запроса по алгоритму Round-Robin.
// Гарантирует потокобезопасность с помощью мьютекса.
func (lb *LoadBalancer) GetNextBackend() *url.URL {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.currentId = (lb.currentId + 1) % uint64(len(lb.backends))
	return lb.backends[lb.currentId]
}

// Ping проверяет доступность указанного сервера-бэкенда.
// Выполняет HTTP-запрос с таймаутом и возвращает true при успешном подключении.
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
	defer resp.Body.Close()

	msg := fmt.Sprintf("Соединение с сервером %s установлено", backend)
	logger.PrintInfo(msg)
	return true
}
