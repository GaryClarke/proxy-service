package main

import (
	"bytes"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	// Create a new instance of our application struct.
	app := newTestApplication(t, false, &segment.SpyClient{})

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

	// Even though Close() can return an error, we deliberately ignore it here.
	// In this context (a test closing an HTTP response body), there isn’t any
	// meaningful recovery we could do if Close failed—our only goal is to free
	// the underlying resources. Using a tiny deferred closure lets us explicitly
	// discard the error (_ = …) without adding linter comments.
	defer func() { _ = rs.Body.Close() }()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	assert.Equal(
		t,
		`{"environment":"testing","status":"available","version":"1.0.0"}`,
		string(body),
	)
}

func TestWebhookHandler_AppleScenarios(t *testing.T) {
	// Define a set of test cases for various Apple subscription scenarios.
	tests := []struct {
		name             string
		payload          string
		expectedStatus   int
		expectedIdentify analytics.Identify
		//expectedTrack    interface{} // Replace with your concrete type
		// Optionally, expected error message etc.
	}{
		{
			name: "Subscription Start|Initial Buy|Not in trial",
			payload: getPayload(
				"SUBSCRIBED",
				"INITIAL_BUY",
				"100000123456789",
				"uk.co.bbc.goodfood2",
				false,
				1,
			),
			expectedStatus: http.StatusNoContent,
			expectedIdentify: expectedIdentifySubStartModel(
				"909365ca-7d09-457d-8f5e-433885a25573", // UserID:Subscription.Properties.IdentityID
				"gf",                                   // Derived
				"909365ca-7d09-457d-8f5e-433885a25573", // AccountGuid:Subscription.Properties.IdentityID
				"100000123456789",                      // SubscriptionID:Subscription.JwsTransaction.OriginalTransactionID
				true,                                   // Derived
				true,                                   // AutoRenewEnabled:Subscription.JwsRenewalInfo.AutoRenewStatus
				"test-airship-channel",                 // AirshipChannelID:Subscription.AirshipChannelID
				"test-airship-id",                      // AirshipID:Subscription.AirshipClaim
			),
			// expectedTrack:    expectedTrackModel("SUB_EVENT_STARTED", "Subscriber", false),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// 1) fresh spy, app, and test server per test
			spy := &segment.SpyClient{}
			app := newTestApplication(t, false, spy)
			ts := newTestServer(t, app.routes())
			defer ts.Close()

			// 2) exercise endpoint
			code, _, _ := ts.postJSON(t, "/webhook", tt.payload)

			// Check the HTTP status code.
			assert.Equal(t, tt.expectedStatus, code)

			// 3) assert Identify was called correctly
			assert.Len(t, spy.Identifies, 1)
			got := spy.Identifies[0]
			assert.Equal(t, tt.expectedIdentify, got)

			// 4) assert Track was called correctly

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

func expectedIdentifySubStartModel(
	userID, brandCode, accountGuid, subscriptionID string,
	subscribed, autoRenew bool,
	airshipChannelID, airshipID string,
) analytics.Identify {
	// replicate the logic in ToIdentify()
	fb := segment.NewFieldBuilder()
	fb.SetIfNotEmpty(fmt.Sprintf("acc_%s_guid", brandCode), accountGuid)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_sub", brandCode), subscribed)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_sub_id", brandCode), subscriptionID)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_auto_renew_status", brandCode), autoRenew)
	fb.SetIfNotEmpty(fmt.Sprintf("%s_airship_channel_id", brandCode), airshipChannelID)
	fb.SetIfNotEmpty(fmt.Sprintf("acc_%s_airship_id", brandCode), airshipID)

	ctx := &analytics.Context{
		Extra: map[string]interface{}{
			"brand_code": brandCode,
		},
	}

	return analytics.Identify{
		UserId:  userID,
		Traits:  fb.Traits(),
		Context: ctx,
	}
}

// getPayload returns a JSON payload string for testing purposes.
func getPayload(
	notificationType, subType, transactionId, brand string,
	inTrial bool,
	autoRenewStatus int,
) string {
	return fmt.Sprintf(`{"payload": "{\"payload\":{\"subscription\":{\"properties\":{\"identityId\":\"909365ca-7d09-457d-8f5e-433885a25573\",\"email\":\"example-apple@email.com\",\"memberStatus\":\"active\",\"entitlement\":\"APP - iOS\",\"timePeriod\":\"Annual\",\"currency\":\"USD\",\"startDate\":\"2106170000\",\"endDate\":\"2107170000\",\"promotionalOfferApplied\":%t,\"productName\":\"exampleSku\",\"productId\":\"testAppleSku-1709199168\",\"platform\":\"ios\",\"client\":null,\"orderId\":null,\"originalTransactionId\":\"%s\",\"version\":\"1.0\"},\"eventTimeMillis\":\"1623936000000\",\"developerNotification\":null,\"subscriptionPurchase\":null,\"jwsRenewalInfo\":{\"environment\":\"Sandbox\",\"originalTransactionId\":\"%s\",\"productId\":\"exampleSku\",\"signedDate\":1623936000000,\"recentSubscriptionStartDate\":1623936000000,\"autoRenewProductId\":null,\"autoRenewStatus\":%d,\"expirationIntent\":null,\"gracePeriodExpiresDate\":null,\"isInBillingRetryPeriod\":false,\"offerIdentifier\":null,\"offerType\":null,\"priceIncreaseStatus\":null},\"jwsTransaction\":{\"bundleId\":\"%s\",\"environment\":\"Sandbox\",\"originalTransactionId\":\"%s\",\"productId\":\"testAppleSku\",\"expiresDate\":1623936000000,\"inAppOwnershipType\":\"PURCHASED\",\"originalPurchaseDate\":1623936000000,\"purchaseDate\":1623936000000,\"signedDate\":1623936000000,\"transactionId\":\"100000123456789\",\"type\":\"%s\",\"quantity\":0,\"appAccountToken\":\"\",\"offerType\":null,\"offerIdentifier\":null,\"isUpgraded\":null,\"revocationDate\":null,\"revocationReason\":null,\"subscriptionGroupIdentifier\":null,\"webOrderLineItemId\":null},\"serverData\":{\"notificationType\":\"%s\",\"subType\":\"%s\",\"notificationUUID\":\"12345678-1234-1234-1234-123456789012\",\"notificationVersion\":\"1.0\",\"appAppleId\":null,\"bundleId\":\"com.example\",\"bundleVersion\":\"1.0\"},\"airshipClaim\":\"test-airship-id\",\"airshipChannelId\":\"test-airship-channel\"}},\"notificationType\":\"mobile-purchase\"}", "meta": {"client": "Sub Notifications API","notification": "AppleIAPNotification"}}`,
		inTrial, transactionId, transactionId, autoRenewStatus, brand, transactionId, subType, notificationType, subType)
}
