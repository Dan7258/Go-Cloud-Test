package loadBalancer

import (
	configHandler "cloud/config_handler"
	"cloud/logger"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backends  []*url.URL
	currentId uint64
	mutex     sync.Mutex
}

func (lb *LoadBalancer) Init(config configHandler.Config) {
	lb.currentId = 0
	lb.backends = make([]*url.URL, len(config.Backends))
	for i, _ := range lb.backends {
		lb.backends[i], _ = url.Parse("http://" + config.Backends[i])
	}
	logger.PrintInfo("Получены данные для балансировщика нагрузки")
}
