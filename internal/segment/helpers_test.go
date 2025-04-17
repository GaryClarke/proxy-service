package segment_test

import (
	"testing"

	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/stretchr/testify/assert"
)

func TestSetIfNotEmpty(t *testing.T) {
	traits := &segment.Traits{
		Traits: make(map[string]interface{}),
	}

	traits.SetIfNotEmpty("non_empty_string", "hello")
	traits.SetIfNotEmpty("empty_string", "")
	traits.SetIfNotEmpty("nil_value", nil)
	traits.SetIfNotEmpty("number", 42)
	traits.SetIfNotEmpty("boolean", false)

	assert.Equal(t, "hello", traits.Traits["non_empty_string"])
	assert.NotContains(t, traits.Traits, "empty_string")
	assert.NotContains(t, traits.Traits, "nil_value")
	assert.Equal(t, 42, traits.Traits["number"])
	assert.Equal(t, false, traits.Traits["boolean"])
}
