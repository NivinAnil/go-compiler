package main

import (
	logger "go-compiler/common/pkg/utils"
	"go-compiler/notification-service/pkg"
	"net/http"
)

func main() {

	//init logger

	StartService()
}

func StartService() {
	log := logger.GetLogger()
	r := pkg.GetRouter()
	port := ":8080"

	httpServer := &http.Server{
		Addr:    port,
		Handler: r,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("listen: ", port)
	}
	log.Info("Service is up and Ready")
}
