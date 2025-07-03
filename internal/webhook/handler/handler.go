// Package handler provides abstractions for webhook dispatching.
package handler

import (
	"context"
	"github.com/garyclarke/proxy-service/internal/webhook"
)

// WebhookHandler is the interface that each handler implements.
type WebhookHandler interface {
	supports(wh webhook.Webhook) bool
	handle(ctx context.Context, wh webhook.Webhook) error
}

// Delegator holds a slice of WebhookHandler.
type Delegator struct {
	handlers []WebhookHandler
}

// NewHandlerDelegator constructs a new HandlerDelegator.
func NewHandlerDelegator(handlers ...WebhookHandler) *Delegator {
	return &Delegator{
		handlers: handlers,
	}
}

// Delegate iterates over the handlers and calls Handle on those that support the webhook.
func (d *Delegator) Delegate(ctx context.Context, wh webhook.Webhook) error {
	for _, handler := range d.handlers {
		if handler.supports(wh) {
			if err := handler.handle(ctx, wh); err != nil {
				return err // or log/accumulate errors as needed
			}
		}
	}
	return nil
}
