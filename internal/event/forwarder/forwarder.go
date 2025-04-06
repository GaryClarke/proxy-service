package forwarder

import "github.com/garyclarke/proxy-service/internal/event"

type EventForwarder interface {
	Supports(event *event.SubscriptionEvent) bool
	Forward(event *event.SubscriptionEvent) error
}
