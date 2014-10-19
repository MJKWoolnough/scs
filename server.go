package scs

import (
	"net"
	"sync"
	"time"
)

const serverWaitTimeOut = time.Second

type DeadlineListener interface {
	net.Listener
	SetDeadline(time.Time) error
}

type TimeoutError interface {
	error
	Timeout() bool
}

func (s *Store) Serve(l DeadlineListener) {
	var w sync.WaitGroup
Server:
	for {
		select {
		case <-s.stop:
			break Server
		default:
		}
		l.SetDeadline(time.Now().Add(serverWaitTimeOut))
		conn, err := l.Accept()
		if err != nil {
			if o, ok := err.(TimeoutError); ok && o.Timeout() {
				continue
			}
			s.logFunc(err)
		}
		w.Add(1)
		go func() {
			defer w.Done()
			s.Process(conn, conn)
			if err = conn.Close(); err != nil {
				s.logFunc(err)
			}
		}()
	}
	w.Wait()
}
