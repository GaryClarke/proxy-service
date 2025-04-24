// Package segment provides helpers for building Segment payloads,
// wrapping github.com/segmentio/analytics-go types and adding
// "set only if non-empty" semantics for Traits.
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
	case *string:
		if v != nil && *v != "" {
			t.Traits[field] = *v
		}
	case bool:
		t.Traits[field] = v
	case *bool:
		if v != nil {
			t.Traits[field] = *v
		}
	default:
		if value != nil {
			t.Traits[field] = value
		}
	}
	return t
}
