package test

import "time"

// Clock is a clock that returns always the same value
type Clock struct {
	Value time.Time
}

// Now returns the Clock value
func (c Clock) Now() time.Time {
	return c.Value
}
