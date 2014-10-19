package scs

import (
	"io"
	"strings"
)

var (
	validGetBytes = validArgBytes + "\t\n\v\f\r\x85\xa0" //any ascii whitespace is allowed as a separator
	getPrefix     = "VALUE "
	getEnd        = "END"
)

const (
	StatCmdGet         = "cmd_get"
	StatCmdGetHits     = "get_hits"
	StatCmdGetMisses   = "get_misses"
	dataLengthEstimate = 30
)

// Get retrieves numerous key/values from the store and prints them to the Writer
func Get(s *Store, _ io.Reader, w io.Writer, args string) {
	if len(args) == 0 {
		WriteError(w, NoArgs{})
		return
	}
	err := ValidateBytes(args, validBytes)
	if err != nil {
		WriteError(w, InvalidChars{err, "in get arguments"})
		return
	}
	keys := strings.Fields(args)

	toWrite := make([]byte, 0, len(keys)*(len(getPrefix)+2*len(commandDelim)+dataLengthEstimate)+len(getEnd)+len(args)+len(commandDelim)) //Estimate for initial storage
	for _, key := range keys {
		data, ok := s.Get(key)
		if ok {
			toWrite = append(toWrite, []byte(getPrefix+key+commandDelim+data+commandDelim)...)
			defer s.SetStat(StatCmdGetHits, Increment)
		} else {
			defer s.SetStat(StatCmdGetMisses, Increment)
		}
	}
	toWrite = append(toWrite, []byte(getEnd+commandDelim)...)
	_, err = w.Write(toWrite)
	if err != nil {
		s.logFunc(err)
	}
	s.SetStat(StatCmdGet, Add(len(keys)))
}

//Errors

type NoArgs struct{}

func (NoArgs) Error() string {
	return "no arguments given"
}
