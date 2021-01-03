package time_measure

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bygui86/go-k8s-probes/commons"
)

func StartTimeMeasure() *TimeMeasure {
	return &TimeMeasure{
		start:      time.Now(),
		stop:       time.Time{},
		delta:      -1,
		deltaNanos: -1,
	}
}

func (t *TimeMeasure) StopTimeMeasure() {
	t.stop = time.Now()
	t.delta = t.stop.Sub(t.start)
	t.deltaNanos = t.stop.UnixNano() - t.start.UnixNano()
}

func (t *TimeMeasure) StopAndLogTimeMeasure(messagePrefix string, loggingFunc func(format string, args ...interface{})) {
	t.stop = time.Now()
	t.delta = t.stop.Sub(t.start)
	t.deltaNanos = t.stop.UnixNano() - t.start.UnixNano()
	loggingFunc(messagePrefix+" took %s seconds / %s milliseconds / %d nanos",
		t.GetDeltaInSec().String(), t.GetDeltaInMil().String(), t.GetDeltaInNanos())
}

//	WARN: This method is made on purpose all-in-line to avoid useless allocations and resources usage
func ShortTimeMeasure(messagePrefix string, loggingFunc func(format string, args ...interface{})) func() {
	start := time.Now()
	return func() {
		stop := time.Now().UnixNano()
		loggingFunc(messagePrefix+" took %s seconds / %s milliseconds / %d nanos",
			decimal.NewFromInt(stop-start.UnixNano()).Div(decimal.NewFromInt(commons.SecondDivider)),
			decimal.NewFromInt(stop-start.UnixNano()).Div(decimal.NewFromInt(commons.MillisecDivider)),
			stop,
		)
	}
}
