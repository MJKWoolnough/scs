package scs

import (
	"io"
	"strconv"
	"strings"
)

const (
	setStored        = "STORED"
	StatCmdSet       = "cmd_set"
	StatCmdSetHits   = "set_hits"
	StatCmdSetMisses = "set_misses"
	keyLimit         = 249
	valueLimit       = 8*1024 - 1 // 8KB - 1
)

// Set adds a key/value pair to the store, overwriting existing data.
func Set(s *Store, r io.Reader, w io.Writer, args string) {
	err := ValidateBytes(args, validArgBytes)
	if err != nil {
		WriteError(w, InvalidChars{err, "in set arguments"})
		return
	}
	if len(args) > keyLimit {
		WriteError(w, ExceededKeyLengthLimit(len(args)))
		return
	}
	data, err := ReadToDelim(r, []byte(commandDelim))
	if err != nil {
		s.logFunc(err)
		WriteError(w, err)
	}
	data = strings.TrimSuffix(data, commandDelim)
	if len(data) > valueLimit {
		WriteError(w, ExceededValueLengthLimit(len(data)))
		return
	}
	err = s.Set(args, data)
	if err != nil {
		WriteError(w, err)
		s.SetStat(StatCmdSetMisses, Increment)
	} else {
		w.Write([]byte(setStored + commandDelim))
		s.SetStat(StatCmdSetHits, Increment)
	}
	s.SetStat(StatCmdSet, Increment)
}

//Errors

type ExceededKeyLengthLimit int

func (e ExceededKeyLengthLimit) Error() string {
	return "key length (" + strconv.Itoa(int(e)) + ") exceeded maxiumum allowed (" + strconv.Itoa(keyLimit) + ")"
}

type ExceededValueLengthLimit int

func (e ExceededValueLengthLimit) Error() string {
	return "value length (" + strconv.Itoa(int(e)) + ") exceeded maxiumum allowed (" + strconv.Itoa(valueLimit) + ")"
}
