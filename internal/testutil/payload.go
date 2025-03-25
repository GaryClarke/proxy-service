package testutil

import "fmt"

// BuildSubNotesPayload returns a JSON payload string for testing purposes,
// containing only the inner payload (i.e. the JSON that maps to SubscriptionPayload).
func BuildSubNotesPayload(
	notificationType, subType, transactionId, brand string,
	inTrial bool,
	autoRenewStatus int,
) string {
	return fmt.Sprintf(`{"payload":{"subscription":{
  "properties":{
    "identityId": "909365ca-7d09-457d-8f5e-433885a25573",
    "email": "example-apple@email.com",
    "memberStatus": "active",
    "entitlement": "APP - iOS",
    "timePeriod": "Annual",
    "currency": "USD",
    "startDate": "2106170000",
    "endDate": "2107170000",
    "promotionalOfferApplied": %t,
    "productName": "exampleSku",
    "productId": "testAppleSku-1709199168",
    "platform": "ios",
    "client": null,
    "orderId": null,
    "originalTransactionId": "%s",
    "version": "1.0"
  },
  "eventTimeMillis": "1623936000000",
  "developerNotification": null,
  "subscriptionPurchase": null,
  "jwsRenewalInfo": {
    "environment": "Sandbox",
    "originalTransactionId": "%s",
    "productId": "exampleSku",
    "signedDate": 1623936000000,
    "recentSubscriptionStartDate": 1623936000000,
    "autoRenewProductId": null,
    "autoRenewStatus": %d,
    "expirationIntent": null,
    "gracePeriodExpiresDate": null,
    "isInBillingRetryPeriod": false,
    "offerIdentifier": null,
    "offerType": null,
    "priceIncreaseStatus": null
  },
  "jwsTransaction": {
    "bundleId": "%s",
    "environment": "Sandbox",
    "originalTransactionId": "%s",
    "productId": "testAppleSku",
    "expiresDate": 1623936000000,
    "inAppOwnershipType": "PURCHASED",
    "originalPurchaseDate": 1623936000000,
    "purchaseDate": 1623936000000,
    "signedDate": 1623936000000,
    "transactionId": "100000123456789",
    "type": "%s",
    "quantity": 0,
    "appAccountToken": "",
    "offerType": null,
    "offerIdentifier": null,
    "isUpgraded": null,
    "revocationDate": null,
    "revocationReason": null,
    "subscriptionGroupIdentifier": null,
    "webOrderLineItemId": null
  },
  "serverData": {
    "notificationType": "%s",
    "subType": "%s",
    "notificationUUID": "12345678-1234-1234-1234-123456789012",
    "notificationVersion": "1.0",
    "appAppleId": null,
    "bundleId": "com.example",
    "bundleVersion": "1.0"
  },
  "airshipClaim": "test-airship-id",
  "airshipChannelId": "test-airship-channel"
}}}`,
		inTrial,          // properties.promotionalOfferApplied
		transactionId,    // properties.originalTransactionId
		transactionId,    // jwsRenewalInfo.originalTransactionId
		autoRenewStatus,  // jwsRenewalInfo.autoRenewStatus
		brand,            // jwsTransaction.bundleId
		transactionId,    // jwsTransaction.originalTransactionId
		subType,          // jwsTransaction.type
		notificationType, // serverData.notificationType
		subType,          // serverData.subType
	)
}
