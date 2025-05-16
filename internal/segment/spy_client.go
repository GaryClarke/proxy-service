package segment

import (
	analytics "github.com/segmentio/analytics-go"
)

// SpyClient records Identify and Track calls.
type SpyClient struct {
	Identifies []analytics.Identify
	Tracks     []analytics.Track
}

func (s *SpyClient) Identify(msg analytics.Identify) error {
	s.Identifies = append(s.Identifies, msg)
	return nil
}

func (s *SpyClient) Track(msg analytics.Track) error {
	s.Tracks = append(s.Tracks, msg)
	return nil
}
