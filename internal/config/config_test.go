package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garyclarke/proxy-service/internal/config"
)

func TestLoadWithArgs_Defaults(t *testing.T) {
	// Ensure no PORT or key in environment
	os.Unsetenv("PORT")
	os.Unsetenv("SEGMENT_SUBSCRIPTION_WRITE_KEY")
	os.Setenv("SEGMENT_SUBSCRIPTION_WRITE_KEY", "key123")

	cfg, err := config.LoadWithArgs([]string{})
	require.NoError(t, err)

	assert.Equal(t, 4000, cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "key123", cfg.SegmentKey)
	assert.Equal(t, "https://events.eu1.segmentapis.com", cfg.SegmentEndpoint)
}

func TestLoadWithArgs_FlagsOverrideEnv(t *testing.T) {
	os.Setenv("PORT", "5555")
	os.Setenv("SEGMENT_SUBSCRIPTION_WRITE_KEY", "envkey")
	os.Setenv("SEGMENT_ENDPOINT", "https://foo.bar")

	args := []string{
		"-port", "8888",
		"-segment-key", "flagkey",
		"-segment-endpoint", "https://baz.qux",
	}
	cfg, err := config.LoadWithArgs(args)
	require.NoError(t, err)

	assert.Equal(t, 8888, cfg.Port)
	assert.Equal(t, "flagkey", cfg.SegmentKey)
	assert.Equal(t, "https://baz.qux", cfg.SegmentEndpoint)
}

func TestLoadWithArgs_InvalidPortEnv(t *testing.T) {
	os.Setenv("PORT", "not-an-int")
	os.Setenv("SEGMENT_SUBSCRIPTION_WRITE_KEY", "somekey")

	cfg, err := config.LoadWithArgs([]string{})
	require.NoError(t, err)

	// invalid PORT should fall back to default
	assert.Equal(t, 4000, cfg.Port)
}

func TestLoadWithArgs_MissingRequired(t *testing.T) {
	os.Unsetenv("SEGMENT_SUBSCRIPTION_WRITE_KEY")
	_, err := config.LoadWithArgs([]string{})
	assert.Error(t, err, "expected an error when the write key is missing")
}
