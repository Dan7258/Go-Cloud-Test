package loadBalancer

import (
	"cloud/config_handler"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backends  []*url.URL
	currentId uint64
	mutex     sync.Mutex
}

func (lb *LoadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	currentHost := lb.GetNextBackend()
	proxy := httputil.NewSingleHostReverseProxy(currentHost)
	proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) GetNextBackend() *url.URL {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.currentId++
	log.Println("Используем: ", lb.backends[lb.currentId%uint64(len(lb.backends))])
	return lb.backends[lb.currentId%uint64(len(lb.backends))]
}

func (lb *LoadBalancer) Init(config configHandler.Config) {
	lb.currentId = 0
	lb.backends = make([]*url.URL, len(config.Backends))
	for i, _ := range lb.backends {
		lb.backends[i], _ = url.Parse("http://" + config.Backends[i])
	}
}
