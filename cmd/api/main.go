package main

import (
	"fmt"
	"log"
	"net/http"
	routes "projeto-api/cmd/api/routes"
	service "projeto-api/cmd/service"
)

const webPort = "8185"

type Config struct{}

func main() {
	// app := Config{}
	log.Printf("Starting api service on port %s\n", webPort)

	go service.ConsumeMessages()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes.SetupRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
	fmt.Println("Hello Word")
}
