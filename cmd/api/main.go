package main

import (
	"flag"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/event/forwarder"
	"github.com/garyclarke/proxy-service/internal/webhook/handler"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Later generate this automatically at build time, but for now just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt reads the environment variable named by key into an int.
// If the variable is not set or cannot be parsed, it returns the fallback value.
func getEnvInt(key string, fallback int) int {
	s := os.Getenv(key)
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("warning: invalid integer for %s=%q, using default %d", key, s, fallback)
		return fallback
	}
	return v
}

// Define a config struct to hold all the configuration settings for the application.
type config struct {
	port            int
	env             string
	debugMode       bool
	segmentKey      string
	segmentEndpoint string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config           config
	logger           *slog.Logger
	handlerDelegator *handler.Delegator
}

func main() {
	// load .env for local/dev
	_ = godotenv.Load()

	// Declare an instance of the config struct.
	var cfg config

	// Flags whose defaults pull from the environment
	flag.IntVar(&cfg.port, "port", getEnvInt("PORT", 4000), "API server port")
	flag.StringVar(&cfg.env, "env", getEnv("ENV", "development"), "Environment")
	flag.StringVar(&cfg.segmentKey, "segment-key", getEnv("SEGMENT_SUBSCRIPTION_WRITE_KEY", ""), "Segment write key")
	flag.StringVar(&cfg.segmentEndpoint, "segment-endpoint",
		getEnv("SEGMENT_ENDPOINT", "https://events.eu1.segmentapis.com"),
		"Segment API endpoint")
	flag.Parse()

	// Validate
	if cfg.segmentKey == "" {
		log.Fatal("SEGMENT_SUBSCRIPTION_WRITE_KEY or -segment-key must be set")
	}

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
