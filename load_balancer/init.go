// Package loadBalancer реализует базовый механизм балансировки нагрузки между серверами-бэкендами.
package loadBalancer

import (
	configHandler "cloud/config_handler"
	"cloud/logger"
	"net/url"
	"sync"
)

// LoadBalancer представляет структуру балансировщика нагрузки.
// Хранит список серверов-бэкендов и текущее состояние для выбора сервера.
type LoadBalancer struct {
	backends  []*url.URL // Список доступных серверов-бэкендов.
	currentId uint64     // Идентификатор для отслеживания текущего выбранного сервера.
	mutex     sync.Mutex // Мьютекс для безопасной работы в многопоточной среде.
}

// Init инициализирует балансировщик нагрузки на основе переданной конфигурации.
// Преобразует строки адресов бэкендов в объекты URL.
func (lb *LoadBalancer) Init(config configHandler.Config) {
	lb.currentId = 0
	lb.backends = make([]*url.URL, len(config.Backends))
	for i := range lb.backends {
		lb.backends[i], _ = url.Parse("http://" + config.Backends[i])
	}
	logger.PrintInfo("Получены данные для балансировщика нагрузки")
}
