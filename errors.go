package scs

import "io"

var (
	errString = []byte{'E', 'R', 'R', 'O', 'R', ' '}
)

func WriteError(w io.Writer, err error) error {
	_, err = w.Write(append(append(errString, err.Error()...), commandDelim...))
	return err
}

//Generic Errors

type InvalidChars struct {
	v  error
	in string
}

func (i InvalidChars) Error() string {
	return i.v.Error() + " " + i.in
}
