package time_measure

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bygui86/go-k8s-probes/commons"
)

/*
	GetDelta returns the delta between two time.Time in the time.Duration specified.
	If time.Duration is not supported, the func will return seconds per default.
*/
func GetDelta(start, end time.Time, format time.Duration) float64 {
	switch format {
	case time.Nanosecond:
		return float64(end.Sub(start).Nanoseconds())
	case time.Microsecond:
		return float64(end.Sub(start).Microseconds())
	case time.Millisecond:
		return float64(end.Sub(start).Milliseconds())
	case time.Second:
		return end.Sub(start).Seconds()
	case time.Minute:
		return end.Sub(start).Minutes()
	case time.Hour:
		return end.Sub(start).Hours()
	}
	return end.Sub(start).Seconds()
}

// unix nanos

func GetDeltaInSec(start, end int64) float64 {
	delta, _ := decimal.
		NewFromInt(end - start).
		Div(decimal.NewFromInt(commons.SecondDivider)).
		Float64()
	return delta
}

func GetDeltaInMillisec(start, end int64) float64 {
	delta, _ := decimal.
		NewFromInt(end - start).
		Div(decimal.NewFromInt(commons.MillisecDivider)).
		Float64()
	return delta
}

func GetDeltaInNanos(start, end int64) float64 {
	delta := end - start
	deltaF, _ := decimal.NewFromInt(delta).Float64()
	return deltaF
}
