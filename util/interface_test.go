package util

import (
	"testing"

	"github.com/kingstonduy/go-core/errorx"
)

func TestCheckNilInterface(t *testing.T) {
	var errx *errorx.Error
	var err error = errx

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "Nil interface",
			input:    nil,
			expected: true,
		},
		{
			name:     "Nil pointer",
			input:    (*int)(nil),
			expected: true,
		},
		{
			name:     "Nil slice",
			input:    ([]int)(nil),
			expected: true,
		},
		{
			name:     "Nil map",
			input:    (map[string]int)(nil),
			expected: true,
		},
		{
			name:     "Nil function",
			input:    (func())(nil),
			expected: true,
		},
		{
			name:     "Nil interface type",
			input:    (interface{})(nil),
			expected: true,
		},
		{
			name:     "Non-nil pointer",
			input:    new(int),
			expected: false,
		},
		{
			name:     "Non-nil slice",
			input:    []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "Non-nil map",
			input:    map[string]int{"key": 1},
			expected: false,
		},
		{
			name:     "Non-nil function",
			input:    func() {},
			expected: false,
		},
		{
			name:     "Non-nil interface with value",
			input:    interface{}(42),
			expected: false,
		},
		{
			name:     "String value",
			input:    "hello",
			expected: false,
		},
		{
			name:     "Errorx interface",
			input:    err,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckNilInterface(tt.input)
			if result != tt.expected {
				t.Errorf("CheckNilInterface(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
