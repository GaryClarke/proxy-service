package testutil

import (
	"fmt"
	"strings"
)

// BuildAppleWebhook returns a JSON payload string for testing purposes,
// containing only the inner payload (i.e. the JSON that maps to SubscriptionPayload).
func BuildAppleWebhook(
	notificationType, subType, transactionId, brand string,
	inTrial bool,
	autoRenewStatus int,
) string {
	inner := fmt.Sprintf(`{"payload":{"subscription":{"properties":{"identityId":"%s","email":"example-apple@email.com","memberStatus":"active","entitlement":"APP - iOS","timePeriod":"Annual","currency":"USD","startDate":"2106170000","endDate":"2107170000","promotionalOfferApplied":%t,"productName":"exampleSku","productId":"testAppleSku-1709199168","platform":"ios","client":null,"orderId":null,"originalTransactionId":"%s","version":"1.0"},"eventTimeMillis":"1623936000000","developerNotification":null,"subscriptionPurchase":null,"jwsRenewalInfo":{"environment":"Sandbox","originalTransactionId":"%s","productId":"exampleSku","signedDate":1623936000000,"recentSubscriptionStartDate":1623936000000,"autoRenewProductId":null,"autoRenewStatus":%d,"expirationIntent":null,"gracePeriodExpiresDate":null,"isInBillingRetryPeriod":false,"offerIdentifier":null,"offerType":null,"priceIncreaseStatus":null},"jwsTransaction":{"bundleId":"%s","environment":"Sandbox","originalTransactionId":"%s","productId":"testAppleSku","expiresDate":1623936000000,"inAppOwnershipType":"PURCHASED","originalPurchaseDate":1623936000000,"purchaseDate":1623936000000,"signedDate":1623936000000,"transactionId":"100000123456789","type":"%s","quantity":0,"appAccountToken":"","offerType":null,"offerIdentifier":null,"isUpgraded":null,"revocationDate":null,"revocationReason":null,"subscriptionGroupIdentifier":null,"webOrderLineItemId":null},"serverData":{"notificationType":"%s","subType":"%s","notificationUUID":"12345678-1234-1234-1234-123456789012","notificationVersion":"1.0","appAppleId":null,"bundleId":"com.example","bundleVersion":"1.0"},"airshipClaim":"%s","airshipChannelId":"%s"}},"notificationType":"mobile-purchase"}`,
		TestUserID, inTrial, transactionId, transactionId, autoRenewStatus, brand, transactionId, subType, notificationType, subType, TestAirshipID, TestAirshipChannelID)

	escaped := strings.ReplaceAll(inner, `"`, `\"`)

	return fmt.Sprintf(`{"payload":"%s","meta":{"client":"Sub Notifications API","notification":"AppleIAPNotification"}}`, escaped)
}

// BuildAppleInnerPayload returns the JSON that maps directly to
// subnotes.SubscriptionPayload (no escaping, no meta).
func BuildAppleInnerPayload(
	notificationType, subType, transactionID, brand string,
	inTrial bool,
	autoRenewStatus int,
) string {
	return fmt.Sprintf(`{
      "payload": {
        "subscription": {
          "properties": {
            "identityId":    "%s",
            "email":         "example-apple@email.com",
            "memberStatus":  "active",
            "entitlement":   "APP - iOS",
            "timePeriod":    "Annual",
            "currency":      "USD",
            "startDate":     "2106170000",
            "endDate":       "2107170000",
            "promotionalOfferApplied": %t,
            "productName":   "exampleSku",
            "productId":     "testAppleSku-1709199168",
            "platform":      "ios",
            "originalTransactionId": "%s"
          },
          "eventTimeMillis": "1623936000000",
          "jwsRenewalInfo": {
            "environment":                "Sandbox",
            "originalTransactionId":      "%s",
            "productId":                  "exampleSku",
            "signedDate":                 1623936000000,
            "recentSubscriptionStartDate":1623936000000,
            "autoRenewStatus":            %d
          },
          "jwsTransaction": {
            "bundleId":      "%s",
            "originalTransactionId": "%s",
            "productId":     "testAppleSku",
            "transactionId": "100000123456789",
            "type":          "%s"
          },
          "serverData": {
            "notificationType": "%s",
            "subType":          "%s"
          },
          "airshipClaim":     "%s",
          "airshipChannelId": "%s"
        }
      }
    }`,
		TestUserID, inTrial,
		transactionID, transactionID, autoRenewStatus,
		brand, transactionID, subType,
		notificationType, subType,
		TestAirshipID, TestAirshipChannelID,
	)
}

func BuildGoogleWebhook(
	notificationType int,
	inTrial bool,
	autoRenewing bool,
	subscriptionID, brand string,
) string {
	inner := fmt.Sprintf(`{"payload":{"subscription":{"properties":{"identityId":"%s","email":"example@email.com","memberStatus":"active","entitlement":"APP - Android","timePeriod":"Annual","currency":"USD","startDate":"2021-06-17T00:00:00.000Z","endDate":"2021-07-17T00:00:00.000Z","promotionalOfferApplied":%t,"productName":"exampleSku","productId":"testGoogleSku-1709199168","platform":"android","client":null,"orderId":"%s","originalTransactionId":null,"version":null},"eventTimeMillis":"1709199142000","developerNotification":{"version":"1.0","packageName":"%s","eventTimeMillis":"1623936000000","subscriptionNotification":{"version":"1.0","notificationType":%d,"purchaseToken":"token","subscriptionId":"%s"}},"subscriptionPurchase":{"latestOrderId":"%s","paymentState":"2","packageName":"%s","productId":"testGoogleSku","purchaseTimeMillis":"1623936000000","purchaseState":0,"purchaseToken":"token","acknowledged":false,"lineItems":[{"expiryTime":"2024-05-14T09:23:27.845Z","productId":"goodfood_a_annual_30_freetrial","autoRenewingPlan":{"autoRenewEnabled":false},"offerDetails":{"basePlanId":"goodfood-a-annual-30-freetrial1","offerId":null,"offerTags":null}},{"expiryTime":"2026-05-14T09:23:27.845Z","productId":"goodfood_a_annual_30_freetrial","autoRenewingPlan":{"autoRenewEnabled":%t},"offerDetails":{"basePlanId":"goodfood-a-annual-30-freetrial1","offerId":null,"offerTags":null}}]},"jwsRenewalInfo":null,"jwsTransaction":null,"serverData":null,"airshipClaim":"%s","airshipChannelId":"%s"}},"notificationType":"mobile-purchase"}`,
		TestUserID, inTrial, subscriptionID, brand, notificationType, subscriptionID,
		subscriptionID, brand, autoRenewing, TestAirshipID, TestAirshipChannelID)

	escaped := strings.ReplaceAll(inner, `"`, `\"`)

	return fmt.Sprintf(`{"payload":"%s","meta":{"client":"Sub Notifications API","notification":"GoogleIAPNotification"}}`, escaped)
}

// BuildGoogleInnerPayload returns the raw JSON that maps directly to
// subnotes.SubscriptionPayload (no escaping, no meta).
func BuildGoogleInnerPayload(
	notificationType int,
	inTrial, autoRenew bool,
	subscriptionID, brand string,
) string {
	return fmt.Sprintf(`{
  "payload": {
    "subscription": {
      "properties": {
        "identityId":              "%s",
        "email":                   "example@email.com",
        "memberStatus":            "active",
        "entitlement":             "APP - Android",
        "timePeriod":              "Annual",
        "currency":                "USD",
        "startDate":               "2021-06-17T00:00:00.000Z",
        "endDate":                 "2021-07-17T00:00:00.000Z",
        "promotionalOfferApplied": %t,
        "productName":             "exampleSku",
        "productId":               "testGoogleSku-1709199168",
        "platform":                "android",
        "orderId":                 "%s"
      },
      "eventTimeMillis": "1709199142000",
      "developerNotification": {
        "version":                 "1.0",
        "packageName":             "%s",
        "eventTimeMillis":         "1623936000000",
        "subscriptionNotification": {
          "version":          "1.0",
          "notificationType": %d,
          "purchaseToken":    "token",
          "subscriptionId":   "%s"
        }
      },
      "subscriptionPurchase": {
        "latestOrderId":      "%s",
        "paymentState":       "2",
        "packageName":        "%s",
        "productId":          "testGoogleSku",
        "purchaseTimeMillis": "1623936000000",
        "purchaseState":      0,
        "purchaseToken":      "token",
        "acknowledged":       false,
        "lineItems": [
          {
            "expiryTime":       "2024-05-14T09:23:27.845Z",
            "productId":        "goodfood_a_annual_30_freetrial",
            "autoRenewingPlan": { "autoRenewEnabled": false },
            "offerDetails": {
              "basePlanId": "goodfood-a-annual-30-freetrial1",
              "offerId":    null,
              "offerTags":  null
            }
          },
          {
            "expiryTime":       "2026-05-14T09:23:27.845Z",
            "productId":        "goodfood_a_annual_30_freetrial",
            "autoRenewingPlan": { "autoRenewEnabled": %t },
            "offerDetails": {
              "basePlanId": "goodfood-a-annual-30-freetrial1",
              "offerId":    null,
              "offerTags":  null
            }
          }
        ]
      },
      "airshipClaim":     "%s",
      "airshipChannelId": "%s"
    }
  }
}`, TestUserID, inTrial, subscriptionID,
		brand, notificationType, subscriptionID,
		subscriptionID, brand, autoRenew,
		TestAirshipID, TestAirshipChannelID)
}
