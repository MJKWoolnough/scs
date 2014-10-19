package scs

import (
	"bytes"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args, result string
	}{
		{"", "ERROR no arguments given\r\n"},
		{"testData1", "VALUE testData1\r\ndata1 abc\r\nEND\r\n"},
		{"testData2", "VALUE testData2\r\ndata2 def\r\nEND\r\n"},
		{"testData1 testData2", "VALUE testData1\r\ndata1 abc\r\nVALUE testData2\r\ndata2 def\r\nEND\r\n"},
		{"testData1 testData3 testData2", "VALUE testData1\r\ndata1 abc\r\nVALUE testData3\r\ndata3 ghi\r\nVALUE testData2\r\ndata2 def\r\nEND\r\n"},
		{"wrong testData2 nokey testData3 notthere", "VALUE testData2\r\ndata2 def\r\nVALUE testData3\r\ndata3 ghi\r\nEND\r\n"},
		{"testData1 £", "ERROR invalid character '\\u00a3' at position 10 in get arguments\r\n"},
	}

	s := NewStore()

	s.Set("testData1", "data1 abc")
	s.Set("testData2", "data2 def")
	s.Set("testData3", "data3 ghi")

	for n, test := range tests {
		w := new(bytes.Buffer)
		Get(s, nil, w, test.args)
		result := w.String()
		if result != test.result {
			t.Errorf("test %d:\nexpecting: %s\ngot: %s", n+1, test.result, result)
		}
	}
}
