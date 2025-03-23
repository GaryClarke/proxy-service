package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("Got %v; want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got %q; expected to contain %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got %v; expected: nil error", actual)
	}
}

func NilFatalError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Fatalf("got %v; expected: nil error", actual)
	}
}

func NotNil(t *testing.T, actual any) {
	t.Helper()
	// Use reflect to catch cases where actual is a typed nil.
	if actual == nil || (reflect.ValueOf(actual).Kind() == reflect.Ptr && reflect.ValueOf(actual).IsNil()) {
		t.Errorf("expected non-nil value, got nil")
	}
}
