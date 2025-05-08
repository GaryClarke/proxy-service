package segment

import "github.com/segmentio/analytics-go"

// analyticsClient wraps the real analytics-go client to implement our Client interface.
type analyticsClient struct {
	client analytics.Client
}

// NewClient constructs a Segment client pointing at the given endpoint.
func NewClient(writeKey, endpoint string) (Client, error) {
	c, err := analytics.NewWithConfig(writeKey, analytics.Config{
		Endpoint: endpoint,
	})
	if err != nil {
		return nil, err
	}
	return &analyticsClient{client: c}, nil
}

// Identify sends an Identify call to Segment.
func (a *analyticsClient) Identify(msg analytics.Identify) error {
	return a.client.Enqueue(msg)
}

// Track sends a Track call to Segment.
func (a *analyticsClient) Track(msg analytics.Track) error {
	return a.client.Enqueue(msg)
}
