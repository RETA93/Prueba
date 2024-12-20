package main

import (
	"fmt"
	"go-project/config"
	"go-project/handlers"
	"go-project/logger"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Initialize configuration and logger
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	// Define routes
	http.HandleFunc("/api/example", handlers.ExampleHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Start the server
	log.Info("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Server failed: %s", err))
	}
}
