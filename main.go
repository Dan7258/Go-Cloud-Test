package main

import (
	"cloud/config_handler"
	"cloud/controllers"
	"cloud/load_balancer"
	"cloud/models"
	"net/http"
)

func main() {
	models.InitENV()
	models.InitDB()
	models.InitRDB()
	config := new(configHandler.Config)
	config.Init()
	lb := new(loadBalancer.LoadBalancer)
	lb.Init(*config)
	http.HandleFunc("/clients/", controllers.ClientHandler)
	//http.HandleFunc("/", lb.ServeProxy)
	http.ListenAndServe(":"+config.Port, nil)
}
