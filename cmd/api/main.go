package main

import (
	"github.com/garyclarke/proxy-service/internal/config"
	"log"
	"os"
)

func main() {
	// Load configuration (flags + env vars)
	cfg := config.Load()

	// Initialize app with config
	app, err := newApplication(cfg)
	if err != nil {
		log.Fatalf("init failed: %v", err)
	}

	// Serve the app
	if err := app.Serve(); err != nil {
		app.logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
