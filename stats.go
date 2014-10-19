package scs

import (
	"bytes"
	"io"
	"strconv"
)

var allStats = [...]string{"cmd_get", "cmd_set", "get_hits", "get_misses", "delete_hits", "delete_misses"}

const (
	StatCurrItems  = "curr_items"
	StatLimitItems = "limit_items"
	StatCmdStats   = "cmd_stats"
	StatsCmdSuffix = "END"
)

// Stats prints store usage statistics to the Writer
func Stats(s *Store, _ io.Reader, w io.Writer, _ string) {
	b := new(bytes.Buffer)
	for _, stat := range allStats {
		b.Write([]byte(stat + " " + strconv.Itoa(s.ReadStat(stat)) + commandDelim))
	}
	b.Write([]byte(StatCurrItems + " " + strconv.Itoa(len(s.data)) + commandDelim + StatLimitItems + " " + strconv.FormatUint(s.limit, 10) + commandDelim + StatsCmdSuffix + commandDelim))
	_, err := b.WriteTo(w)
	if err != nil {
		s.logFunc(err)
	}
	s.SetStat(StatCmdStats, Increment)
}
