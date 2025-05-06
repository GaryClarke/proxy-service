package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DebugMode       bool
	Env             string
	Port            int
	SegmentKey      string
	SegmentEndpoint string
	Version         string
}

// LoadWithArgs parses flags out of the provided args slice,
// falling back to env-vars (and defaults) when a flag isn’t set.
func LoadWithArgs(args []string) (Config, error) {
	// 1) In dev, load .env (no-op if missing)
	_ = godotenv.Load()

	var cfg Config
	fs := flag.NewFlagSet("proxy-service", flag.ContinueOnError)

	// 2) Register exactly the flags you want, with env-var defaults
	fs.IntVar(&cfg.Port, "port",
		getEnvInt("PORT", 4000),
		"HTTP server port")
	fs.StringVar(&cfg.Env, "env",
		getEnv("ENV", "development"),
		"Environment (development|staging|production)")
	fs.StringVar(&cfg.SegmentKey, "segment-key",
		getEnv("SEGMENT_SUBSCRIPTION_WRITE_KEY", ""),
		"Segment write key")
	fs.StringVar(&cfg.SegmentEndpoint, "segment-endpoint",
		getEnv("SEGMENT_ENDPOINT", "https://events.eu1.segmentapis.com"),
		"Segment API endpoint")

	// 3) Parse only the args passed in (no -test.* flags or any global flags)
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	// 4) Validate required
	if cfg.SegmentKey == "" {
		return Config{}, fmt.Errorf("SEGMENT_SUBSCRIPTION_WRITE_KEY or -segment-key must be set")
	}

	// 5) Populate Version (or later via ldflags)
	cfg.Version = "1.0.0"
	return cfg, nil
}

// Load is the convenience wrapper for real runs.
// It parses os.Args[1:] and fatals on error.
func Load() Config {
	cfg, err := LoadWithArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// getEnv returns the value of the environment variable named by key.
// If the variable is empty or not set, it returns fallback.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt reads an integer env-var or returns the provided fallback.
// Logs a warning if the value is present but not a valid int.
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
