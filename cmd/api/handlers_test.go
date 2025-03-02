package main

import (
	"bytes"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/assert"
	"io"
	"net/http"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	// Create a new instance of our application struct.
	app := newTestApplication(t)

	// create a new test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained in
	// the ts.URL field. We can use this along with the ts.Client().Get() method
	// to make a GET /healthcheck request against the test server. This returns a
	// http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rs.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, string(body), `{"environment":"testing","status":"available","version":"1.0.0"}`)
}

func TestWebhookHandler_AppleScenarios(t *testing.T) {
	// Define a set of test cases for various Apple subscription scenarios.
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		//expectedIdentify interface{} // Replace with your concrete type
		//expectedTrack    interface{} // Replace with your concrete type
		// Optionally, expected error message etc.
	}{
		{
			name:           "Subscription Start - Initial Buy, not in trial",
			payload:        getPayload("SUBSCRIBED", "INITIAL_BUY", "false", 1, "100000123456789", "uk.co.bbc.goodfood2"),
			expectedStatus: http.StatusNoContent,
			// expectedIdentify: expectedIdentifySubStartModel(true), // Your helper function
			// expectedTrack:    expectedTrackModel("SUB_EVENT_STARTED", "Subscriber", false),
		},
	}

	// Create your test server using the application's routes.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.postJSON(t, "/webhook", tt.payload)

			// Check the HTTP status code.
			if code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, code)
			}

			// Next, you can assert that your mock segment client (or any other dependency)
			// was called with the expected identify and track models.
			//
			// For instance, if your segment client is a dependency in your application and
			// you can inspect its state after the request, you could do:
			//
			// gotIdentify := app.segmentClient.GetLatestIdentifyModel()
			// assert.Equal(t, gotIdentify, tt.expectedIdentify)
			//
			// gotTrack := app.segmentClient.GetLatestTrackModel()
			// assert.Equal(t, gotTrack, tt.expectedTrack)
			//
			// These functions would be part of your fake/mock segment client that you inject
			// into your application for testing.
		})
	}
}

// getPayload returns a JSON payload as a string.
// You can adjust the parameters as needed.
func getPayload(
	notificationType string,
	subType string,
	inTrial string,
	autoRenewStatus int,
	transactionId string,
	brand string,
) string {
	// Note: The inner JSON payload is a string literal,
	// so you need to escape the inner quotes.
	return fmt.Sprintf(`{
    "payload": "{\"payload\":{\"subscription\":{\"properties\":{\"identityId\":\"909365ca-7d09-457d-8f5e-433885a25573\",\"email\":\"example-apple@email.com\",
    \"memberStatus\":\"active\",\"entitlement\":\"APP - iOS\",\"timePeriod\":\"Annual\",\"currency\":\"USD\",
    \"startDate\":\"2106170000\",\"endDate\":\"2107170000\",\"promotionalOfferApplied\":%s,\"productName\":\"exampleSku\",
    \"productId\":\"testAppleSku-1709199168\",\"platform\":\"ios\",\"client\":null,\"orderId\":null,
    \"originalTransactionId\":\"%s\",\"version\":\"1.0\"},\"eventTimeMillis\":\"1623936000000\",\"developerNotification\":null,
    \"subscriptionPurchase\":null,\"jwsRenewalInfo\":{\"environment\":\"Sandbox\",\"originalTransactionId\":\"%s\",
    \"productId\":\"exampleSku\",\"signedDate\":1623936000000,\"recentSubscriptionStartDate\":1623936000000,
    \"autoRenewProductId\":null,\"autoRenewStatus\":%d,\"expirationIntent\":null,\"gracePeriodExpiresDate\":null,
    \"isInBillingRetryPeriod\":false,\"offerIdentifier\":null,\"offerType\":null,\"priceIncreaseStatus\":null},
    \"jwsTransaction\":{\"bundleId\":\"%s\",\"environment\":\"Sandbox\",\"originalTransactionId\":\"%s\",
    \"productId\":\"testAppleSku\",\"expiresDate\":1623936000000,\"inAppOwnershipType\":\"PURCHASED\",
    \"originalPurchaseDate\":1623936000000,\"purchaseDate\":1623936000000,\"signedDate\":1623936000000,
    \"transactionId\":\"100000123456789\",\"type\":\"%s\",\"quantity\":0,\"appAccountToken\":\"\",
    \"offerType\":null,\"offerIdentifier\":null,\"isUpgraded\":null,\"revocationDate\":null,\"revocationReason\":null,
    \"subscriptionGroupIdentifier\":null,\"webOrderLineItemId\":null},\"serverData\":{\"notificationType\":\"%s\",
    \"subType\":\"%s\",\"notificationUUID\":\"12345678-1234-1234-1234-123456789012\",\"notificationVersion\":\"1.0\",
    \"appAppleId\":null,\"bundleId\":\"com.example\",\"bundleVersion\":\"1.0\"},\"airshipClaim\":\"test-airship-id\",
    \"airshipChannelId\":\"test-airship-channel\"}},\"notificationType\":\"mobile-purchase\"}",
    "meta": {
        "client": "Sub Notifications API",
        "notification": "AppleIAPNotification"
    }
}`, inTrial, transactionId, transactionId, autoRenewStatus, brand, transactionId, subType, notificationType, subType)
}
