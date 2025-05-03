package main

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/config"
	"github.com/garyclarke/proxy-service/internal/event/forwarder"
	"github.com/garyclarke/proxy-service/internal/webhook/handler"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config           config.Config
	logger           *slog.Logger
	handlerDelegator *handler.Delegator
}

func main() {
	// Load configuration (flags + env vars)
	cfg := config.Load()

	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize the WebhookHandlers and initialize the delegator
	appleHandler := handler.NewAppleHandler([]forwarder.EventForwarder{&forwarder.AppleSubscriptionStartForwarder{}})

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config:           cfg,
		logger:           logger,
		handlerDelegator: handler.NewHandlerDelegator(appleHandler),
	}

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
