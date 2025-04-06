package event_test

import (
	"testing"

	"github.com/garyclarke/proxy-service/internal/assert"
)

func TestGetLookupData_Apple(t *testing.T) {
	data, err := event.GetLookupData("apple")
	assert.NilFatalError(t, err)
	assert.NotNil(t, data)

	// Check correct number of items in apple map
	assert.Equal(t, len(data), 20)

	// Verify that a specific event key exists and has expected values.
	event, ok := data["SUBSCRIBED|INITIAL_BUY|false"]
	if !ok {
		t.Fatal("expected key 'SUBSCRIBED|INITIAL_BUY|false' not found")
	}
	assert.Equal(t, event.Name, "subscription_started")
	assert.Equal(t, event.Category, "CATEGORY_START")
	assert.Equal(t, *event.SubStatus, "First time subscriber")
}

func TestGetLookupData_Google(t *testing.T) {
	data, err := event.GetLookupData("google")
	assert.NilFatalError(t, err)
	assert.NotNil(t, data)

	// Check correct number of items in google map
	assert.Equal(t, len(data), 14)

	// Verify that a specific event key exists and has expected values.
	event, ok := data["PURCHASED|null|false"]
	if !ok {
		t.Fatal("expected key 'PURCHASED|null|false' not found")
	}
	assert.Equal(t, event.Name, "subscription_started")
	assert.Equal(t, event.Category, "CATEGORY_START")
	assert.Equal(t, *event.SubStatus, "Subscriber")
}
