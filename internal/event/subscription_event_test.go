package event_test

import (
	"testing"

	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/stretchr/testify/assert"
)

func TestGetLookupData_Apple(t *testing.T) {
	data, err := event.GetLookupData("apple")
	assert.NoError(t, err, "expected no error when loading Apple lookup data")
	assert.NotNil(t, data, "expected non‑nil map for Apple lookup data")
	assert.Equal(t, 20, len(data), "expected 20 items in Apple lookup data")

	// Verify that a specific event key exists and has expected values.
	subEvent, ok := data["SUBSCRIBED|INITIAL_BUY|false"]
	assert.True(t, ok, "expected key 'SUBSCRIBED|INITIAL_BUY|false' to be present")
	assert.Equal(t, "subscription_started", subEvent.Name)
	assert.Equal(t, "CATEGORY_START", subEvent.Category)
	if assert.NotNil(t, subEvent.SubStatus, "SubStatus pointer should not be nil") {
		assert.Equal(t, "First time subscriber", *subEvent.SubStatus)
	}
}

func TestGetLookupData_Google(t *testing.T) {
	data, err := event.GetLookupData("google")
	assert.NoError(t, err, "expected no error when loading Google lookup data")
	assert.NotNil(t, data, "expected non‑nil map for Google lookup data")
	assert.Equal(t, 14, len(data), "expected 14 items in Google lookup data")

	// Verify that a specific event key exists and has expected values.
	subEvent, ok := data["PURCHASED|null|false"]
	assert.True(t, ok, "expected key 'PURCHASED|null|false' to be present")
	assert.Equal(t, "subscription_started", subEvent.Name)
	assert.Equal(t, "CATEGORY_START", subEvent.Category)
	if assert.NotNil(t, subEvent.SubStatus, "SubStatus pointer should not be nil") {
		assert.Equal(t, "Subscriber", *subEvent.SubStatus)
	}
}
