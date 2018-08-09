package util

import (
	"time"
)

// Time is a helper routine that allocates a new time value
// to store v and returns a pointer to it.
func Time(v time.Time) *time.Time {
	return &v
}
