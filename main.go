package main

import (
	"cloud/config_handler"
	"cloud/load_balancer"
	"net/http"
)

func main() {
	config := new(configHandler.Config)
	config.Init()
	lb := new(loadBalancer.LoadBalancer)
	lb.Init(*config)
	http.HandleFunc("/", lb.ServeProxy)
	http.ListenAndServe(":"+config.Port, nil)
}
