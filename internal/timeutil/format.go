package timeutil

import "time"

// ParseSubscriptionDate parses a yyMMddHHmm timestamp (e.g. "2106170000")
// and returns an RFC3339 pointer or nil on error.
func ParseSubscriptionDate(v string) *string {
	t, err := time.Parse("0601021504", v)
	if err != nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
