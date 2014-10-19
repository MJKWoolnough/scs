package scs

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
)

const (
	DefaultLimit = 65535
	commandDelim = "\r\n"
)

type Cmd func(*Store, io.Reader, io.Writer, string)

func noopLogFunc(e error) {}

var defaultCommands = map[string]Cmd{
	"set":    Set,
	"get":    Get,
	"delete": Delete,
	"stats":  Stats,
}

type Store struct {
	caseInsensitiveKeys bool
	locked              bool
	data                map[string]string
	commands            map[string]Cmd
	stats               map[string]int
	dataMutex           sync.RWMutex
	statMutex           sync.RWMutex
	logFunc             func(error)
	requestStop         chan struct{}
	stop                chan struct{}
	limit               uint64
}

// NewStore returns a new Store struct with the various Options applied.
//
// Options can only be applied at this time.
// Store.Stop() should be called when finished with, regardless of whether
// a server was started.
func NewStore(o ...Option) *Store {
	s := &Store{
		data:        make(map[string]string),
		commands:    make(map[string]Cmd),
		stats:       make(map[string]int),
		logFunc:     noopLogFunc,
		requestStop: make(chan struct{}),
		stop:        make(chan struct{}),
		limit:       DefaultLimit,
	}
	for _, opt := range o {
		opt(s)
	}
	s.locked = true
	if len(s.commands) == 0 {
		s.commands = defaultCommands
	}
	go interruptHandler(s)
	return s
}

func interruptHandler(s *Store) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		s.logFunc(InterruptError{})
	case <-s.requestStop:
	}
	signal.Stop(c)
	close(s.requestStop)
	s.requestStop = nil
	close(s.stop)
}

// SetStat modifies the requested stat with the supplied StatChange function.
func (s *Store) SetStat(stat string, sc StatChange) {
	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	s.stats[stat] = sc(s.stats[stat])
}

// Stats return a full list of stat keys.
func (s *Store) Stats() []string {
	s.statMutex.RLock()
	defer s.statMutex.RUnlock()
	stats := make([]string, len(s.stats))
	for stat := range s.stats {
		stats = append(stats, stat)
	}
	return stats
}

// ReadStat returns the requested stat
func (s *Store) ReadStat(stat string) int {
	s.statMutex.RLock()
	defer s.statMutex.RUnlock()
	return s.stats[stat]
}

// Set adds a key/value to the store. Returns a Full error when limit is reached.
func (s *Store) Set(name, data string) error {
	if s.caseInsensitiveKeys {
		name = strings.ToLower(name)
	}
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()
	if uint64(len(s.data)) >= s.limit {
		return Full(s.limit)
	}
	s.data[name] = data
	return nil
}

// Get retrieves a key/value from the store. Returned bool indicates existence of key
// in case of empty string.
func (s *Store) Get(name string) (string, bool) {
	if s.caseInsensitiveKeys {
		name = strings.ToLower(name)
	}
	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()
	d, ok := s.data[name]
	return d, ok
}

// Delete removes a key/value from the store. Returned bool indicates previous existence
// of key.
func (s *Store) Delete(name string) bool {
	if s.caseInsensitiveKeys {
		name = strings.ToLower(name)
	}
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()
	if _, ok := s.data[name]; !ok {
		return false
	}
	delete(s.data, name)
	return true
}

// Stop shuts down any running server and stops any more from starting.
func (s *Store) Stop() {
	if s.requestStop != nil {
		s.requestStop <- struct{}{}
	}
}

//Errors

// Full is an error which is returned when trying to add to a store that is full
type Full int

func (Full) Error() string {
	return "cannot set, storage full"
}

type InterruptError struct{}

func (InterruptError) Error() string {
	return "interrupt received - quitting"
}
