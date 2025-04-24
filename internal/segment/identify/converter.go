// Package identify contains helpers for converting payloads into Segment Identify messages.
package identify

import "github.com/segmentio/analytics-go"

// IdentifyConverter is an interface for events that can be converted into a Segment identify message.
type IdentifyConverter interface {
	// ToIdentify converts the implementation to an analytics.Identify message,
	// which can then be enqueued for sending to Segment.
	ToIdentify() analytics.Identify
}
