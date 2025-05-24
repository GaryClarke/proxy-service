package ptr

// Str returns a pointer to the given string.
func Str(s string) *string { return &s }

// Bool returns a pointer to the given bool.
func Bool(b bool) *bool { return &b }

// Int returns a pointer to the given int.
func Int(i int) *int { return &i }

// BoolFromIntPtr maps a *int (0 or 1) to a *bool; nil or out‐of‐range → nil.
func BoolFromIntPtr(ip *int) *bool {
	if ip == nil {
		return nil
	}
	switch *ip {
	case 0:
		b := false
		return &b
	case 1:
		b := true
		return &b
	default:
		return nil
	}
}
