package main

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/config"
	"github.com/garyclarke/proxy-service/internal/event/forwarder"
	"github.com/garyclarke/proxy-service/internal/segment"
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
	// 1) set up logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 2) wire up the real Segment client
	segmentClient, err := segment.NewClient(cfg.SegmentKey, cfg.SegmentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("segment client init: %w", err)
	}

	// 3) build your forwarders using that client
	appleForwarders := []forwarder.EventForwarder{
		forwarder.NewAppleSubscriptionStartForwarder(segmentClient),
		forwarder.NewAppleSubscriptionTrackForwarder(segmentClient),
	}

	appleHandler := handler.NewAppleHandler(appleForwarders)
	delegator := handler.NewHandlerDelegator(appleHandler)

	return &application{
		config:           cfg,
		logger:           logger,
		handlerDelegator: delegator,
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
