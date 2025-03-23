package subnotes

import "github.com/garyclarke/proxy-service/internal/brand"

// SubscriptionPayload mirrors the nested structure of the inner payload.
type SubscriptionPayload struct {
	Payload struct {
		Subscription Subscription `json:"subscription"`
	} `json:"payload"`
}

// Subscription represents the subscription details extracted from the webhook payload.
// Sub notes nests everything related to the subscription inside here
type Subscription struct {
	Properties            SubscriptionProperties `json:"properties"`
	EventTimeMillis       int64                  `json:"eventTimeMillis"`
	DeveloperNotification *DeveloperNotification `json:"developerNotification"`
	SubscriptionPurchase  *SubscriptionPurchase  `json:"subscriptionPurchase"`
	JwsRenewalInfo        *JwsRenewalInfo        `json:"jwsRenewalInfo"`
	JwsTransaction        *JwsTransaction        `json:"jwsTransaction"`
	ServerData            *ServerData            `json:"serverData"`
	AirshipClaim          *string                `json:"airshipClaim"`
	AirshipChannelID      *string                `json:"airshipChannelId"`
	Brand                 brand.Brand            `json:"-"`
}

type DeveloperNotification struct {
}

type SubscriptionPurchase struct {
}

type SubscriptionProperties struct {
	IdentityID              *string `json:"identityId"`
	Email                   *string `json:"email"`
	MemberStatus            string  `json:"memberStatus"`
	Entitlement             string  `json:"entitlement"`
	TimePeriod              string  `json:"timePeriod"`
	Currency                *string `json:"currency"`
	StartDate               string  `json:"startDate"`
	EndDate                 string  `json:"endDate"`
	PromotionalOfferApplied bool    `json:"promotionalOfferApplied"`
	ProductName             string  `json:"productName"`
	ProductID               *string `json:"productId"`
	Platform                *string `json:"platform"`
	Client                  *string `json:"client"`
	OrderID                 *string `json:"orderId"`
	OfferID                 *string `json:"offerId"`
	FreeTrial               *bool   `json:"freeTrial"`
	OriginalTransactionID   *string `json:"originalTransactionId"`
	Version                 *string `json:"version"`
}

// JwsRenewalInfo represents the decoded payload for JWS renewal information.
type JwsRenewalInfo struct {
	Environment                 string  `json:"environment"`                 // non-nullable
	OriginalTransactionID       string  `json:"originalTransactionId"`       // non-nullable
	ProductID                   string  `json:"productId"`                   // non-nullable
	SignedDate                  int64   `json:"signedDate"`                  // non-nullable
	RecentSubscriptionStartDate int64   `json:"recentSubscriptionStartDate"` // non-nullable
	AutoRenewProductID          *string `json:"autoRenewProductId"`          // nullable
	AutoRenewStatus             *int    `json:"autoRenewStatus"`             // nullable
	ExpirationIntent            *int    `json:"expirationIntent"`            // nullable
	GracePeriodExpiresDate      *int    `json:"gracePeriodExpiresDate"`      // nullable (using int for consistency with PHP)
	IsInBillingRetryPeriod      bool    `json:"isInBillingRetryPeriod"`      // non-nullable; defaults to false
	OfferIdentifier             *string `json:"offerIdentifier"`             // nullable
	OfferType                   *int    `json:"offerType"`                   // nullable
	PriceIncreaseStatus         *int    `json:"priceIncreaseStatus"`         // nullable
}

type JwsTransaction struct {
	BundleID                    string  `json:"bundleId"`
	Environment                 string  `json:"environment"`
	OriginalTransactionID       string  `json:"originalTransactionId"`
	ProductID                   string  `json:"productId"`
	ExpiresDate                 int64   `json:"expiresDate"`
	InAppOwnershipType          string  `json:"inAppOwnershipType"`
	OriginalPurchaseDate        int64   `json:"originalPurchaseDate"`
	PurchaseDate                int64   `json:"purchaseDate"`
	SignedDate                  int64   `json:"signedDate"`
	TransactionID               string  `json:"transactionId"`
	Type                        string  `json:"type"`
	Quantity                    int     `json:"quantity"`
	AppAccountToken             string  `json:"appAccountToken"`
	OfferType                   *string `json:"offerType"`
	OfferIdentifier             *string `json:"offerIdentifier"`
	IsUpgraded                  *bool   `json:"isUpgraded"`
	RevocationDate              *int64  `json:"revocationDate"`
	RevocationReason            *string `json:"revocationReason"`
	SubscriptionGroupIdentifier *string `json:"subscriptionGroupIdentifier"`
	WebOrderLineItemID          *string `json:"webOrderLineItemId"`
}

type ServerData struct {
	NotificationType    string  `json:"notificationType"`
	SubType             *string `json:"subType"`
	NotificationUUID    string  `json:"notificationUUID"`
	NotificationVersion string  `json:"notificationVersion"`
	AppAppleID          *string `json:"appAppleId"`
	BundleID            string  `json:"bundleId"`
	BundleVersion       string  `json:"bundleVersion"`
}
