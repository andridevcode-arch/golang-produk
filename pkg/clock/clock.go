package clock

import (
	"time"
)

func Now() time.Time {
	return time.Now().UTC()
}

func Unix() int64 {
	return time.Now().Unix()
}
