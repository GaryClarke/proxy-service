package segment_test

import (
	"github.com/garyclarke/proxy-service/internal/ptr"
	"testing"

	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"

	"github.com/garyclarke/proxy-service/internal/segment"
)

func TestFieldBuilder_SetIfNotEmpty_and_Exports(t *testing.T) {
	fb := segment.NewFieldBuilder().
		// strings
		SetIfNotEmpty("s1", "hello").
		SetIfNotEmpty("s2", "").
		// *string
		SetIfNotEmpty("p1", ptr.Str("world")).
		SetIfNotEmpty("p2", (*string)(nil)).
		// bools
		SetIfNotEmpty("b1", true).
		SetIfNotEmpty("b2", (*bool)(nil)).
		// ints (default case)
		SetIfNotEmpty("i1", 42).
		SetIfNotEmpty("i2", nil)

	// As Traits
	traits := fb.Traits()
	assert.Equal(t, "hello", traits["s1"])
	assert.NotContains(t, traits, "s2")
	assert.Equal(t, "world", traits["p1"])
	assert.NotContains(t, traits, "p2")
	assert.Equal(t, true, traits["b1"])
	assert.NotContains(t, traits, "b2")
	assert.Equal(t, 42, traits["i1"])
	assert.NotContains(t, traits, "i2")

	// As Properties
	props := fb.Properties()
	// same assertions
	assert.Equal(t, analytics.Properties(traits), props)
}
