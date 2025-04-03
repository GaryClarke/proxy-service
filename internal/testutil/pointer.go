package testutil

// PtrStr returns a pointer to the given string.
func PtrStr(s string) *string {
	return &s
}

// PtrBool returns a pointer to the given bool.
func PtrBool(b bool) *bool {
	return &b
}
