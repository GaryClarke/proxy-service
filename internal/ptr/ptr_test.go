package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolFromIntPtr(t *testing.T) {
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
			want: Bool(false),
		},
		{
			name: "1 input yields pointer to true",
			in:   &one,
			want: Bool(true),
		},
		{
			name: "other value yields nil",
			in:   &two,
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolFromIntPtr(tc.in)
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

func TestStr(t *testing.T) {
	s := "hello"
	got := Str(s)
	assert.NotNil(t, got)
	assert.Equal(t, s, *got)
}

func TestBool(t *testing.T) {
	b := true
	got := Bool(b)
	assert.NotNil(t, got)
	assert.Equal(t, b, *got)
}

func TestInt(t *testing.T) {
	i := 42
	got := Int(i)
	assert.NotNil(t, got)
	assert.Equal(t, i, *got)
}
