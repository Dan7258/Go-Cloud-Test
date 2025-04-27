package main

import (
	"cloud/config_handler"
	"cloud/load_balancer"
	"cloud/models"
	"net/http"
)

func main() {
	models.InitENV()
	models.InitDB()
	config := new(configHandler.Config)
	config.Init()
	lb := new(loadBalancer.LoadBalancer)
	lb.Init(*config)
	http.HandleFunc("/", lb.ServeProxy)
	http.ListenAndServe(":"+config.Port, nil)
}
