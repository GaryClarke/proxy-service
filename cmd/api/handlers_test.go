package main

import (
	"bytes"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	testUserID           = "909365ca-7d09-457d-8f5e-433885a25573"
	testBrandCode        = "gf"
	testSubscriptionID   = "100000123456789"
	testAirshipChannelID = "test-airship-channel"
	testAirshipID        = "test-airship-id"
	testCurrency         = "USD"
	testFrequency        = "Annual"
	testProductName      = "exampleSku"
	testRenewalDate      = "2021-07-17T00:00:00Z"
	testStartDate        = "2021-06-17T00:00:00Z"
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
		expectedTrack    analytics.Track
		// Optionally, expected error message etc.
	}{
		{
			name: "Subscription Start|Initial Buy|Not in trial",
			payload: getApplePayload(
				"SUBSCRIBED",
				"INITIAL_BUY",
				testSubscriptionID,
				"uk.co.bbc.goodfood2",
				false,
				1,
			),
			expectedStatus: http.StatusNoContent,
			expectedIdentify: expectedIdentifySubStartModel(
				testUserID,           // UserID:Subscription.Properties.IdentityID
				testBrandCode,        // Derived
				testUserID,           // AccountGuid:Subscription.Properties.IdentityID
				testSubscriptionID,   // SubscriptionID:Subscription.JwsTransaction.OriginalTransactionID
				true,                 // Derived
				true,                 // AutoRenewEnabled:Subscription.JwsRenewalInfo.AutoRenewStatus
				testAirshipChannelID, // AirshipChannelID:Subscription.AirshipChannelID
				testAirshipID,        // AirshipID:Subscription.AirshipClaim
			),
			expectedTrack: expectedTrackModel(
				// event, subStatus, inTrial, autoRenew, subType, platform
				"subscription_started",  // lookup event name
				"First time subscriber", // lookup subStatus
				false,                   // inTrial
				true,                    // autoRenew
				"INITIAL_BUY",           // subType (notification subtype)
				"ios",                   // platform (from payload JSON)
			),
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
			assert.Len(t, spy.Tracks, 1)
			gotTrack := spy.Tracks[0]
			assert.Equal(t, tt.expectedTrack, gotTrack)
		})
	}
}

func TestWebhookHandler_GoogleScenarios(t *testing.T) {
	// Define a set of test cases for various Google subscription scenarios.
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		//expectedIdentify analytics.Identify
		//expectedTrack    analytics.Track
		// Optionally, expected error message etc.
	}{
		{
			name: "Subscription Start|Not in trial",
			payload: getGooglePayload(
				4,     // notificationType = PURCHASED
				false, // inTrial = false
				true,  // autoRenewing = true (default Google behavior)
				testSubscriptionID,
				"uk.co.bbc.goodfood2",
			),
			expectedStatus: http.StatusNoContent,
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

			// 3) assert HTTP response
			assert.Equal(t, tt.expectedStatus, code)
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

func expectedTrackModel(
	event string,
	subStatus string,
	inTrial bool,
	autoRenew bool,
	subType string,
	platform string,
) analytics.Track {
	// Build the traits
	traits := segment.NewFieldBuilder().
		SetIfNotEmpty(fmt.Sprintf("acc_%s_guid", testBrandCode), testUserID).
		SetIfNotEmpty(fmt.Sprintf("app_%s_sub_id", testBrandCode), testSubscriptionID).
		SetIfNotEmpty(fmt.Sprintf("%s_airship_channel_id", testBrandCode), ptr.Str(testAirshipChannelID)).
		SetIfNotEmpty(fmt.Sprintf("acc_%s_airship_id", testBrandCode), ptr.Str(testAirshipID))

	// Build the properties
	props := segment.NewFieldBuilder().
		SetIfNotEmpty("airship_channel_id", ptr.Str(testAirshipChannelID)).
		SetIfNotEmpty("airship_id", ptr.Str(testAirshipID)).
		SetIfNotEmpty("sub_id", testSubscriptionID).
		SetIfNotEmpty("brand_code", testBrandCode).
		SetIfNotEmpty("sub_type", ptr.Str(subType)).
		SetIfNotEmpty("platform", &platform).
		SetIfNotEmpty("sub_auto_renew_status", &autoRenew).
		SetIfNotEmpty("sub_currency", ptr.Str(testCurrency)).
		SetIfNotEmpty("sub_frequency", ptr.Str(testFrequency)).
		SetIfNotEmpty("sub_in_trial", &inTrial).
		SetIfNotEmpty("sub_product_name", ptr.Str(testProductName)).
		SetIfNotEmpty("sub_renewal_date", ptr.Str(testRenewalDate)).
		SetIfNotEmpty("sub_start_date", ptr.Str(testStartDate)).
		SetIfNotEmpty("sub_status", &subStatus).
		SetIfNotEmpty("sub_with_offer", &inTrial)

	// Assemble context
	ctx := &analytics.Context{
		Extra: map[string]interface{}{
			"brand_code": testBrandCode,
			"traits":     traits.ToMap(),
		},
	}

	return analytics.Track{
		Event:      event,
		UserId:     testUserID,
		Properties: props.Properties(),
		Context:    ctx,
	}
}

// getApplePayload returns a JSON payload string for testing purposes.
func getApplePayload(
	notificationType, subType, transactionId, brand string,
	inTrial bool,
	autoRenewStatus int,
) string {
	inner := fmt.Sprintf(`{"payload":{"subscription":{"properties":{"identityId":"%s","email":"example-apple@email.com","memberStatus":"active","entitlement":"APP - iOS","timePeriod":"Annual","currency":"USD","startDate":"2106170000","endDate":"2107170000","promotionalOfferApplied":%t,"productName":"exampleSku","productId":"testAppleSku-1709199168","platform":"ios","client":null,"orderId":null,"originalTransactionId":"%s","version":"1.0"},"eventTimeMillis":"1623936000000","developerNotification":null,"subscriptionPurchase":null,"jwsRenewalInfo":{"environment":"Sandbox","originalTransactionId":"%s","productId":"exampleSku","signedDate":1623936000000,"recentSubscriptionStartDate":1623936000000,"autoRenewProductId":null,"autoRenewStatus":%d,"expirationIntent":null,"gracePeriodExpiresDate":null,"isInBillingRetryPeriod":false,"offerIdentifier":null,"offerType":null,"priceIncreaseStatus":null},"jwsTransaction":{"bundleId":"%s","environment":"Sandbox","originalTransactionId":"%s","productId":"testAppleSku","expiresDate":1623936000000,"inAppOwnershipType":"PURCHASED","originalPurchaseDate":1623936000000,"purchaseDate":1623936000000,"signedDate":1623936000000,"transactionId":"100000123456789","type":"%s","quantity":0,"appAccountToken":"","offerType":null,"offerIdentifier":null,"isUpgraded":null,"revocationDate":null,"revocationReason":null,"subscriptionGroupIdentifier":null,"webOrderLineItemId":null},"serverData":{"notificationType":"%s","subType":"%s","notificationUUID":"12345678-1234-1234-1234-123456789012","notificationVersion":"1.0","appAppleId":null,"bundleId":"com.example","bundleVersion":"1.0"},"airshipClaim":"%s","airshipChannelId":"%s"}},"notificationType":"mobile-purchase"}`,
		testUserID, inTrial, transactionId, transactionId, autoRenewStatus, brand, transactionId, subType, notificationType, subType, testAirshipID, testAirshipChannelID)

	escaped := strings.ReplaceAll(inner, `"`, `\"`)

	return fmt.Sprintf(`{"payload":"%s","meta":{"client":"Sub Notifications API","notification":"AppleIAPNotification"}}`, escaped)
}

// getGooglePayload returns a JSON payload string for Google tests.
func getGooglePayload(
	notificationType int,
	inTrial bool,
	autoRenewing bool,
	subscriptionID, brand string,
) string {
	inner := fmt.Sprintf(`{"payload":{"subscription":{"properties":{"identityId":"%s","email":"example@email.com","memberStatus":"active","entitlement":"APP - Android","timePeriod":"Annual","currency":"USD","startDate":"2021-06-17T00:00:00.000Z","endDate":"2021-07-17T00:00:00.000Z","promotionalOfferApplied":%t,"productName":"exampleSku","productId":"testGoogleSku-1709199168","platform":"android","client":null,"orderId":"%s","originalTransactionId":null,"version":null},"eventTimeMillis":"1709199142000","developerNotification":{"version":"1.0","packageName":"%s","eventTimeMillis":"1623936000000","subscriptionNotification":{"version":"1.0","notificationType":%d,"purchaseToken":"token","subscriptionId":"%s"}},"subscriptionPurchase":{"latestOrderId":"%s","paymentState":"2","packageName":"%s","productId":"testGoogleSku","purchaseTimeMillis":"1623936000000","purchaseState":0,"purchaseToken":"token","acknowledged":false,"lineItems":[{"expiryTime":"2024-05-14T09:23:27.845Z","productId":"goodfood_a_annual_30_freetrial","autoRenewingPlan":{"autoRenewEnabled":false},"offerDetails":{"basePlanId":"goodfood-a-annual-30-freetrial1","offerId":null,"offerTags":null}},{"expiryTime":"2026-05-14T09:23:27.845Z","productId":"goodfood_a_annual_30_freetrial","autoRenewingPlan":{"autoRenewEnabled":%t},"offerDetails":{"basePlanId":"goodfood-a-annual-30-freetrial1","offerId":null,"offerTags":null}}]},"jwsRenewalInfo":null,"jwsTransaction":null,"serverData":null,"airshipClaim":"%s","airshipChannelId":"%s"}},"notificationType":"mobile-purchase"}`,
		testUserID, inTrial, subscriptionID, brand, notificationType, subscriptionID,
		subscriptionID, brand, autoRenewing, testAirshipID, testAirshipChannelID)

	escaped := strings.ReplaceAll(inner, `"`, `\"`)

	return fmt.Sprintf(`{"payload":"%s","meta":{"client":"Sub Notifications API","notification":"GoogleIAPNotification"}}`, escaped)
}
