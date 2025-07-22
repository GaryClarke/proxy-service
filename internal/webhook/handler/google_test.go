package handler

import (
	"testing"

	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/stretchr/testify/assert"
)

func TestGoogleHandler_Supports(t *testing.T) {
	h := NewGoogleHandler()

	cases := []struct {
		name  string
		notif string
		want  bool
	}{
		{"exact match", GoogleNotification, true},
		{"wrong notif", "FooBar", false},
		{"apple notif", AppleNotification, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := webhook.Webhook{
				Meta: webhook.Meta{Notification: tc.notif},
			}
			got := h.supports(wh)
			assert.Equal(t, tc.want, got)
		})
	}
}
