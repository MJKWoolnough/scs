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

func TestProcessFull(t *testing.T) {
	t.Parallel()

	input := "set sushi\r\n" +
		"delicious\r\n" +
		"set topcoder\r\n" +
		"fun\r\n" +
		"get sushi topcoder\r\n" +
		"get sushi\r\n" +
		"delete sushi\r\n" +
		"get sushi\r\n" +
		"get topcoder\r\n" +
		"get topcoder sushi\r\n" +
		"delete sushi\r\n" +
		"stats\r\n" +
		"quit\r\n"
	output := "STORED\r\n" +
		"STORED\r\n" +
		"VALUE sushi\r\n" +
		"delicious\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"VALUE sushi\r\n" +
		"delicious\r\n" +
		"END\r\n" +
		"DELETED\r\n" +
		"END\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"NOT_FOUND\r\n" +
		"cmd_get 7\r\n" +
		"cmd_set 2\r\n" +
		"get_hits 5\r\n" +
		"get_misses 2\r\n" +
		"delete_hits 1\r\n" +
		"delete_misses 1\r\n" +
		"curr_items 1\r\n" +
		"limit_items 65535\r\n" +
		"END\r\n"
	s := NewStore()
	w := new(bytes.Buffer)
	s.Process(strings.NewReader(input), w)
	result := w.String()
	if result != output {
		t.Errorf("expecting: %s\ngot: %s", output, result)
	}
}
