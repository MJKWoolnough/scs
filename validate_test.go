package scs

import (
	"reflect"
	"testing"
)

func TestValidateBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input, valid string
		expected     error
	}{
		{"teststring", "abcdefghijklmnopqrstuvwxyz", nil},
		{"teststring", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", ValidationError{'t', 0}},
		{"example args[with] {valid} ...punctuation... @marks@", validBytes, nil},
		{"example args[with] {valid} ...punctuation... @marks@", validArgBytes, ValidationError{' ', 7}},
	}

	for n, test := range tests {
		if err := ValidateBytes(test.input, test.valid); !reflect.DeepEqual(err, test.expected) {
			t.Errorf("test %d: expecting error: %v\ngot: %v", n+1, test.expected, err)
		}
	}

}
