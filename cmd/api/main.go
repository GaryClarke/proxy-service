package main

import (
	"flag"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/webhook/handler"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Later generate this automatically at build time, but for now just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for the application.
type config struct {
	port      int
	env       string
	debugMode bool
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config           config
	logger           *slog.Logger
	handlerDelegator *handler.Delegator
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config:           cfg,
		logger:           logger,
		handlerDelegator: handler.NewHandlerDelegator(),
	}

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
