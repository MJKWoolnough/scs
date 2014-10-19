package scs

import (
	"bytes"
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input, output string
	}{
		{"quit\r\n", ""},
		{"", "ERROR EOF\r\n"},
		{"get unknown\r\nquit\r\n", "END\r\n"},
		{"set myVar\r\ndataHere\r\nget myVar\r\nquit\r\n", "STORED\r\nVALUE myVar\r\ndataHere\r\nEND\r\n"},
		{"set myVar2\r\ndataHere2\r\nget myVar2\r\ndelete myVar2\r\nquit\r\n", "STORED\r\nVALUE myVar2\r\ndataHere2\r\nEND\r\nDELETED\r\n"},
		{"delete myVar3\r\nquit\r\n", "NOT_FOUND\r\n"},
	}

	s := NewStore(Limit(5))

	for n, test := range tests {
		w := new(bytes.Buffer)
		s.Process(strings.NewReader(test.input), w)
		result := w.String()
		if result != test.output {
			t.Errorf("test %d:\nexpecting: %s\ngot: %s", n+1, test.output, result)
		}
	}
}
