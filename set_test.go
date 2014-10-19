package scs

import (
	"bytes"
	"strings"
	"testing"
)

func TestSet(t *testing.T) {
	t.Parallel()

	eightkb := make([]byte, 8*1024)
	tests := []struct {
		key, value, result string
	}{
		{"test1", "value1", "STORED\r\n"},
		{"test1 test2", "", "ERROR invalid character ' ' at position 5 in set arguments\r\n"},
		{"test2", "value2", "STORED\r\n"},
		{
			"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
				"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
				"0123456789012345678901234567890123456789012345678",
			"value3",
			"STORED\r\n",
		},
		{
			"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
				"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
				"01234567890123456789012345678901234567890123456789",
			"",
			"ERROR key length (250) exceeded maxiumum allowed (249)\r\n",
		},
		{"largeTest1", string(eightkb), "ERROR value length (8192) exceeded maxiumum allowed (8191)\r\n"},
		{"largeTest2", string(eightkb[:len(eightkb)-1]), "STORED\r\n"},
	}

	s := NewStore()

	for n, test := range tests {
		w := new(bytes.Buffer)
		Set(s, strings.NewReader(test.value+commandDelim), w, test.key)
		result := w.String()
		if result != test.result {
			t.Errorf("test %d:\nexpecting: %s\ngot: %s", n+1, test.result, result)
		} else if test.result[:5] != "ERROR" {
			if r, _ := s.Get(test.key); r != test.value {
				t.Errorf("test %d:\nexpecting stored value: %q\ngot: %q", n+1, test.value, r)
			}
		} else if _, ok := s.Get(test.key); ok {
			t.Errorf("test %d: expecting no value, got one", n+1)
		}
	}
}
