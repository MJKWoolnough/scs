package scs

import "io"

const (
	deleteSucessful     = "DELETED"
	deleteUnsucessful   = "NOT_FOUND"
	StatCmdDelete       = "cmd_delete"
	StatCmdDeleteHit    = "delete_hits"
	StatCmdDeleteMisses = "delete_misses"
)

// Delete removes a key from the store if it exists. Indicates success on Writer.
func Delete(s *Store, _ io.Reader, w io.Writer, args string) {
	err := ValidateBytes(args, validArgBytes)
	if err != nil {
		WriteError(w, InvalidChars{err, "in delete arguments"})
		return
	}
	s.SetStat(StatCmdDelete, Increment)
	if s.Delete(args) {
		s.SetStat(StatCmdDeleteHit, Increment)
		w.Write([]byte(deleteSucessful + commandDelim))
	} else {
		s.SetStat(StatCmdDeleteMisses, Increment)
		w.Write([]byte(deleteUnsucessful + commandDelim))
	}
}
