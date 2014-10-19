package scs

type StatChange func(int) int

// Add i to the stat given to.
func Add(i int) StatChange {
	return func(j int) int {
		return i + j
	}
}

// SetStat sets the value of the stat to the given value.
func SetStat(i int) StatChange {
	return func(_ int) int {
		return i
	}
}

// Increment increases the stat by one. Same as calling Add(1).
func Increment(i int) int {
	return i + 1
}

// Decrement decreases the stat by one. Same as calling Add(-1).
func Decrement(i int) int {
	return i - 1
}

// Zero resets the stat. Same as calling SetStat(0).
func Zero(i int) int {
	return 0
}
