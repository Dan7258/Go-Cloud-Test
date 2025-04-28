package main

import (
	"cloud/config_handler"
	"cloud/controllers"
	"cloud/load_balancer"
	"cloud/models"
	"cloud/rate_limiter"
	"net/http"
)

// Main запускает все необходимые функции.
func main() {
	config := new(configHandler.Config)
	config.Init()
	models.InitENV()
	models.InitDB()
	models.InitRDB()
	lb := new(loadBalancer.LoadBalancer)
	lb.Init(*config)
	rateLimiter.InitTBConfig(*config)
	go rateLimiter.StartTokenTicker()
	http.HandleFunc("/clients/", controllers.ClientHandler)
	http.HandleFunc("/", lb.ServeProxy)
	http.ListenAndServe(":"+config.Port, nil)
}
