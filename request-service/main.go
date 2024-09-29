package main

import (
	"fmt"
	"go-compiler/request-service/pkg/router"
	"net/http"
)

func main() {
	appRouter := router.GetRouter()
	port := ":8080"

	httpServer := &http.Server{
		Addr:    port,
		Handler: appRouter,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port", port)

}
