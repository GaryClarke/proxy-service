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

// newApplication initializes the application struct.
func newApplication(cfg config.Config) (*application, error) {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize the WebhookHandlers and initialize the delegator
	appleHandler := handler.NewAppleHandler([]forwarder.EventForwarder{&forwarder.AppleSubscriptionStartForwarder{}})

	return &application{
		config:           cfg,
		logger:           logger,
		handlerDelegator: handler.NewHandlerDelegator(appleHandler),
	}, nil
}

// Serve starts the HTTP server on the configured port.
func (app *application) Serve() error {
	addr := fmt.Sprintf(":%d", app.config.Port)
	app.logger.Info("starting server", "addr", addr, "env", app.config.Env)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	return srv.ListenAndServe()
}
