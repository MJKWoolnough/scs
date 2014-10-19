package scs

// Option type to set options in the store upon creation
type Option func(*Store)

// Limit sets the maximum number of items the store will hold
func Limit(l uint64) Option {
	return func(s *Store) {
		s.limit = l
	}
}

// CaseInsensitiveKeys sets the store to ignore case when setting, getting and deleting keys
func CaseInsensitiveKeys() Option {
	return func(s *Store) {
		if !s.locked {
			s.caseInsensitiveKeys = true
		}
	}
}

// LogFunc sets a custom logging function for the stores' errors
func LogFunc(l func(error)) Option {
	return func(s *Store) {
		if !s.locked {
			s.logFunc = l
		}
	}
}

// RegisterCommands allows the adding of custom commands to the store. cmds is a map of command
// names to functions.
//
// Adding this option stops the deault commands from being applied. To have the default commands
// in addition to the custom ones, also supply the RegisterDefaultOption option to NewStore.
func RegisterCommands(cmds map[string]Cmd) Option {
	return func(s *Store) {
		if !s.locked {
			for str, cmd := range cmds {
				s.commands[str] = cmd
			}
		}
	}
}

// RegisterDefaultCommands applies the default get, set, delete and stats commands to the store.
//
// Useful when adding custom functions as the default will not be applied when custom commands are
// supplied.
func RegisterDefaultCommands() Option {
	return RegisterCommands(defaultCommands)
}
