// Package google internal/webhook/google/lookup.go
package google

// NotificationType is the Google RTDN subscription code.
type NotificationType int

const (
	NotificationRecovered NotificationType = 1
	NotificationRenewed   NotificationType = 2
	NotificationCancelled NotificationType = 3
	NotificationPurchased NotificationType = 4
	NotificationOnHold    NotificationType = 5
	NotificationRestarted NotificationType = 7
	NotificationPaused    NotificationType = 10
	NotificationRevoked   NotificationType = 12
	NotificationExpired   NotificationType = 13
)
