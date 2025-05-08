package segment

import "github.com/segmentio/analytics-go"

// Client defines the methods our forwarders use to send events.
type Client interface {
	Identify(msg analytics.Identify) error
	Track(msg analytics.Track) error
}
