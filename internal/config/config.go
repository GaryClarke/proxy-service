package config

import (
	"flag"
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

func Load() Config {
	// 1) In dev, read .env into os.Environ
	_ = godotenv.Load()

	var cfg Config

	// 2) Declare your flags, defaulting from env-vars
	flag.IntVar(&cfg.Port, "port",
		getEnvInt("PORT", 4000),
		"HTTP server port")
	flag.StringVar(&cfg.Env, "env",
		getEnv("ENV", "development"),
		"Environment (development|staging|production)")
	flag.StringVar(&cfg.SegmentKey, "segment-key",
		getEnv("SEGMENT_SUBSCRIPTION_WRITE_KEY", ""),
		"Segment write key")
	flag.StringVar(&cfg.SegmentEndpoint, "segment-endpoint",
		getEnv("SEGMENT_ENDPOINT", "https://events.eu1.segmentapis.com"),
		"Segment API endpoint")
	flag.Parse()

	// 3) Validate required settings
	if cfg.SegmentKey == "" {
		log.Fatal("SEGMENT_SUBSCRIPTION_WRITE_KEY or -segment-key must be set")
	}

	// todo Later generate this automatically at build time
	cfg.Version = "1.0.0"

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
