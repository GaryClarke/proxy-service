package testutil

// PtrStr returns a pointer to the given string.
func PtrStr(s string) *string {
	return &s
}

// PtrBool returns a pointer to the given bool.
func PtrBool(b bool) *bool {
	return &b
}

// PtrInt returns a pointer to the given int.
func PtrInt(i int) *int {
	return &i
}

// PtrBoolFromInt returns a pointer to a bool if ip is 0 or 1, otherwise nil.
// If ip is nil, returns nil.
func PtrBoolFromInt(ip *int) *bool {
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
