// Package track contains helpers for converting payloads into Segment Track messages.
package track

import "github.com/segmentio/analytics-go"

// TrackConverter converts a struct into an analytics.Track message.
type TrackConverter interface {
	// ToTrack converts the implementation into an analytics.Track message,
	// which can then be enqueued for sending to Segment.
	ToTrack() analytics.Track
}
