package track

import "github.com/segmentio/analytics-go"

// TrackConverter is an interface for payloads that can be converted into a Segment track message.
type TrackConverter interface {
	// ToTrack converts the implementation into an analytics.Track message,
	// which can then be enqueued for sending to Segment.
	ToTrack() analytics.Track
}
