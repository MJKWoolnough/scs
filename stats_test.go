package scs

import (
	"bytes"
	"testing"
)

func TestStats(t *testing.T) {
	t.Parallel()

	type cs map[string]StatChange

	tests := []struct {
		changeStats cs
		expecting   string
	}{
		{
			cs{
				"cmd_get": Increment,
			},
			"cmd_get 1\r\ncmd_set 0\r\nget_hits 0\r\nget_misses 0\r\ndelete_hits 0\r\ndelete_misses 0\r\ncurr_items 0\r\nlimit_items 65535\r\n",
		},
		{
			cs{
				"get_hits":      SetStat(2),
				"cmd_get":       Decrement,
				"get_misses":    SetStat(3),
				"cmd_set":       Increment,
				"delete_hits":   SetStat(4),
				"delete_misses": SetStat(5),
			},
			"cmd_get 0\r\ncmd_set 1\r\nget_hits 2\r\nget_misses 3\r\ndelete_hits 4\r\ndelete_misses 5\r\ncurr_items 0\r\nlimit_items 65535\r\n",
		},
	}
	s := NewStore()
	for n, test := range tests {
		for k, v := range test.changeStats {
			s.SetStat(k, v)
		}
		w := new(bytes.Buffer)
		Stats(s, nil, w, "")
		result := w.String()
		if result != test.expecting {
			t.Errorf("test %d:\nexpecting: %s\ngot: %s", n+1, test.expecting, result)
		}
	}
}
