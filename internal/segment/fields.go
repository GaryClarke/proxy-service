package segment

import "github.com/segmentio/analytics-go"

// FieldBuilder is a generic builder for maps used in both Traits and Properties.
type FieldBuilder struct {
	data map[string]interface{}
}

// NewFieldBuilder creates a new builder with an empty map.
func NewFieldBuilder() *FieldBuilder {
	return &FieldBuilder{data: make(map[string]interface{})}
}

// SetIfNotEmpty only adds the field if v is non-nil (and non-empty for strings).
// It returns itself so you can chain calls.
func (b *FieldBuilder) SetIfNotEmpty(key string, v interface{}) *FieldBuilder {
	switch val := v.(type) {
	case string:
		if val != "" {
			b.data[key] = val
		}
	case *string:
		if val != nil && *val != "" {
			b.data[key] = *val
		}
	case bool:
		b.data[key] = val
	case *bool:
		if val != nil {
			b.data[key] = *val
		}
	default:
		if v != nil {
			b.data[key] = v
		}
	}
	return b
}

// Traits exports the built map as analytics.Traits.
func (b *FieldBuilder) Traits() analytics.Traits {
	return analytics.Traits(b.data)
}

// Properties exports the built map as analytics.Properties.
func (b *FieldBuilder) Properties() analytics.Properties {
	return analytics.Properties(b.data)
}

func (b *FieldBuilder) ToMap() map[string]interface{} {
	return b.data
}
