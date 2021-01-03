package time_measure

import (
	"time"
)

type TimeMeasure struct {
	start      time.Time
	stop       time.Time
	delta      time.Duration
	deltaNanos int64
}
