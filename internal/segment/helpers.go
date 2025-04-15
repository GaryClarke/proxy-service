package segment

import "github.com/segmentio/analytics-go"

// Traits embeds analytics.Traits.
type Traits struct {
	analytics.Traits
}

// SetIfNotEmpty sets a key-value pair only if the value is not empty (or nil), and returns the Traits to enable chaining.
func (t *Traits) SetIfNotEmpty(field string, value interface{}) *Traits {
	switch v := value.(type) {
	case string:
		if v != "" {
			t.Traits[field] = value
		}
	case nil:
		// Do nothing for nil
	default:
		t.Traits[field] = value
	}
	return t
}
