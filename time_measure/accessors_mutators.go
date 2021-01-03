package time_measure

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bygui86/go-k8s-probes/commons"
)

func (t *TimeMeasure) GetStart() time.Time {
	return t.start
}

func (t *TimeMeasure) GetStop() time.Time {
	return t.stop
}

func (t *TimeMeasure) GetDelta() time.Duration {
	return t.delta
}

func (t *TimeMeasure) GetDeltaInSec() decimal.Decimal {
	return decimal.NewFromInt(t.deltaNanos).Div(decimal.NewFromInt(commons.SecondDivider))
}

func (t *TimeMeasure) GetDeltaInMil() decimal.Decimal {
	return decimal.NewFromInt(t.deltaNanos).Div(decimal.NewFromInt(commons.MillisecDivider))
}

func (t *TimeMeasure) GetDeltaInNanos() int64 {
	return t.deltaNanos
}
