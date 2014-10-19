package scs

import (
	"strconv"
	"strings"
)

var (
	validBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'\"*+-/=?^_`{|}~()<>[]:;@,. "
	validArgBytes = validBytes[:len(validBytes)-2] //No Space
)

func ValidateBytes(data, valid string) error {
	for p, c := range data {
		if strings.IndexRune(valid, c) == -1 {
			return ValidationError{c, p}
		}
	}
	return nil
}

//Errors

type ValidationError struct {
	c rune
	n int
}

func (v ValidationError) Error() string {
	return "invalid character " + strconv.QuoteRuneToASCII(v.c) + " at position " + strconv.Itoa(v.n)
}
