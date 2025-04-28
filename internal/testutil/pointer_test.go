package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtrBoolFromInt(t *testing.T) {
	zero := 0
	one := 1
	two := 2

	tests := []struct {
		name string
		in   *int
		want *bool
	}{
		{
			name: "nil input yields nil",
			in:   nil,
			want: nil,
		},
		{
			name: "0 input yields pointer to false",
			in:   &zero,
			want: PtrBool(false),
		},
		{
			name: "1 input yields pointer to true",
			in:   &one,
			want: PtrBool(true),
		},
		{
			name: "other value yields nil",
			in:   &two,
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PtrBoolFromInt(tc.in)
			if tc.want == nil {
				assert.Nil(t, got)
			} else {
				// not nil, so must point to the same boolean value
				assert.NotNil(t, got)
				assert.Equal(t, *tc.want, *got)
			}
		})
	}
}
