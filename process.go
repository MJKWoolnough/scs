package scs

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	QuitCmd         = "quit"
	StatConnections = "connections"
)

func (s *Store) Process(r io.Reader, w io.Writer) {
	r = bufio.NewReader(r)
	for {
		line, err := ReadToDelim(r, []byte(commandDelim))
		if err != nil {
			s.logFunc(err)
			WriteError(w, err)
			return
		}
		if line == commandDelim {
			continue
		}
		line = strings.TrimSuffix(line, commandDelim)
		cmd := strings.SplitN(line, " ", 2)
		cmdStr := cmd[0]
		args := ""
		if len(cmd) > 1 {
			args = cmd[1]
		}
		if cmdStr == QuitCmd {
			return
		}
		cmdFn := s.commands[cmdStr]
		if cmdFn == nil {
			WriteError(w, NoCommand(cmdStr))
		} else {
			cmdFn(s, r, w, args)
		}
	}
	s.SetStat(StatConnections, Increment)
}

func ReadToDelim(r io.Reader, delim []byte) (string, error) {
	read := make([]byte, 0, 32)
	char := make([]byte, 1)
	for !bytes.HasSuffix(read, delim) {
		_, err := r.Read(char)
		if err != nil {
			return "", err
		}
		read = append(read, char[0])
	}
	return string(read), nil
}

//Errors

type NoCommand string

func (n NoCommand) Error() string {
	return "command not found: " + string(n)
}
