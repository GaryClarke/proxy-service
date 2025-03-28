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
		"100000123456789",
		"uk.co.bbc.goodfood2",
		false,
		1,
	)

	// Decode the payload.
	sub, err := decodeSubscriptionWebhook(payload)

	// Assert
	assert.NilFatalError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, sub.Brand, brand.GF)

	// Assert some values in nested structs.
	assert.Equal(t, sub.Properties.MemberStatus, "active")
	assert.Equal(t, sub.Properties.Entitlement, "APP - iOS")
	assert.Equal(t, sub.Properties.TimePeriod, "Annual")

	// Assert a value in jwsTransaction.
	assert.NotNil(t, sub.JwsTransaction)
	assert.Equal(t, sub.JwsTransaction.BundleID, "uk.co.bbc.goodfood2")
	assert.Equal(t, sub.JwsTransaction.Type, "INITIAL_BUY")

	// Assert a value in jwsRenewalInfo.
	assert.NotNil(t, sub.JwsRenewalInfo)
	assert.Equal(t, sub.JwsRenewalInfo.Environment, "Sandbox")

	// Assert a value in serverData.
	assert.NotNil(t, sub.ServerData)
	assert.Equal(t, sub.ServerData.NotificationType, "SUBSCRIBED")
	assert.Equal(t, *sub.ServerData.SubType, "INITIAL_BUY")
}

func TestDecodeSubscriptionWebhook_InvalidPayload(t *testing.T) {
	invalidPayload := `{"invalid": "data"}`
	_, err := decodeSubscriptionWebhook(invalidPayload)
	if err == nil {
		t.Fatal("expected error for invalid payload, got nil")
	}
}
