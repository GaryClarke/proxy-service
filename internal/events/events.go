package events

import (
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/lookup"
)

// SubscriptionEvent represents the lookup data for a subscription event.
type SubscriptionEvent struct {
	Name             string  `json:"name"`
	SubStatus        *string `json:"sub_status"`
	Category         string  `json:"category"`
	NotificationType string  `json:"notification_type"`
	SubType          *string `json:"sub_type"`
	InTrial          *bool   `json:"in_trial"`
}

// GetLookupData reads the lookup file for the specified provider (e.g. "apple" or "google")
// and returns the parsed event definitions.
func GetLookupData(provider string) (map[string]SubscriptionEvent, error) {
	filename := fmt.Sprintf("%s.json", provider)
	data, err := lookup.Data.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read lookup file %s: %w", filename, err)
	}
	var events map[string]SubscriptionEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lookup data from %s: %w", filename, err)
	}
	return events, nil
}
