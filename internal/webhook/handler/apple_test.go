package handler

import (
	"github.com/garyclarke/proxy-service/internal/assert"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"testing"
)

func TestDecodeSubscriptionWebhook_ValidPayload(t *testing.T) {
	// Build a valid payload using the shared test utility function.
	payload := testutil.BuildSubNotesPayload(
		"SUBSCRIBED",
		"INITIAL_BUY",
		"false",
		"100000123456789",
		"uk.co.bbc.goodfood2",
		1,
	)

	// Decode the payload.
	sub, err := DecodeSubscriptionWebhook(payload)

	// Assert
	assert.NilFatalError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, sub.Brand, brand.GF)
}

//func TestDecodeSubscriptionWebhook_InvalidPayload(t *testing.T) {
//	invalidPayload := `{"invalid": "data"}`
//	_, err := decoder.DecodeSubscriptionWebhook(invalidPayload)
//	if err == nil {
//		t.Fatal("expected error for invalid payload, got nil")
//	}
//}
