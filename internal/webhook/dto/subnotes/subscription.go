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
	EventTimeMillis       string                 `json:"eventTimeMillis"`       // Android -> eventTimeMillis, iOS signedDate
	DeveloperNotification *DeveloperNotification `json:"developerNotification"` // Android only
	SubscriptionPurchase  *SubscriptionPurchase  `json:"subscriptionPurchase"`  // Android only
	JwsRenewalInfo        *JwsRenewalInfo        `json:"jwsRenewalInfo"`        // iOS only
	JwsTransaction        *JwsTransaction        `json:"jwsTransaction"`        // iOS only
	ServerData            *ServerData            `json:"serverData"`            // iOS only
	AirshipClaim          *string                `json:"airshipClaim"`          // Only if available
	AirshipChannelID      *string                `json:"airshipChannelId"`      // Only if available
	Brand                 brand.Brand            `json:"-"`
}

type DeveloperNotification struct {
	Version                  string                    `json:"version"`
	PackageName              string                    `json:"packageName"`
	EventTimeMillis          string                    `json:"eventTimeMillis"`
	SubscriptionNotification *SubscriptionNotification `json:"subscriptionNotification"`
}

// SubscriptionPurchase corresponds to Google’s SubscriptionPurchaseV2.
type SubscriptionPurchase struct {
	LatestOrderID      string                         `json:"latestOrderId"`
	PaymentState       string                         `json:"paymentState"`
	PackageName        string                         `json:"packageName"`
	ProductID          string                         `json:"productId"`
	PurchaseTimeMillis string                         `json:"purchaseTimeMillis"`
	PurchaseState      int                            `json:"purchaseState"`
	PurchaseToken      string                         `json:"purchaseToken"`
	Acknowledged       bool                           `json:"acknowledged"`
	LineItems          []SubscriptionPurchaseLineItem `json:"lineItems"`
	AutoRenewing       bool                           `json:"autoRenewing"`
}

// SubscriptionPurchaseLineItem represents each element in the lineItems array.
type SubscriptionPurchaseLineItem struct {
	ExpiryTime       string                   `json:"expiryTime"`
	ProductID        string                   `json:"productId"`
	AutoRenewingPlan AutoRenewingPlan         `json:"autoRenewingPlan"`
	OfferDetails     SubscriptionOfferDetails `json:"offerDetails"`
}

// AutoRenewingPlan matches the nested autoRenewingPlan block.
type AutoRenewingPlan struct {
	AutoRenewEnabled bool `json:"autoRenewEnabled"`
}

// SubscriptionOfferDetails matches the nested offerDetails block.
type SubscriptionOfferDetails struct {
	BasePlanID string   `json:"basePlanId"`
	OfferID    *string  `json:"offerId"`
	OfferTags  []string `json:"offerTags"`
}

// SubscriptionNotification carries the fields we care about for subscriptions.
type SubscriptionNotification struct {
	Version          string `json:"version"`
	NotificationType int    `json:"notificationType"`
	PurchaseToken    string `json:"purchaseToken"`
	SubscriptionId   string `json:"subscriptionId"`
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
