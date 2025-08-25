package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	log.Printf("Loading Configuration ...\n")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Register Routes
	mux := router.NewRouter(cfg)

	// Add Middleware
	mw := middleware.RateLimit(mux, cfg)

	log.Printf("Starting API Gateway on port %d\n", cfg.Server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mw)
}
