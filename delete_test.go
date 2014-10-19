package scs

import (
	"bytes"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key, result string
	}{
		{"test1", "DELETED\r\n"},
		{"test1", "NOT_FOUND\r\n"},
		{"test2", "DELETED\r\n"},
		{"test invalid", "ERROR invalid character ' ' at position 4 in delete arguments\r\n"},
	}

	s := NewStore()

	s.Set("test1", "data1")
	s.Set("test2", "data2")

	for n, test := range tests {
		w := new(bytes.Buffer)
		Delete(s, nil, w, test.key)
		result := w.String()
		if result != test.result {
			t.Errorf("test %d: expecting: %s\ngot: %s", n+1, test.result, result)
		}
	}
}
