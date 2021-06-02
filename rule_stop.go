package techan

import "github.com/sdcoffey/big"

type (
	stopLossOrGainRule struct {
		Indicator
		tolerance big.Decimal
		mode      cmpMode
	}
	cmpMode int
)

const (
	lte cmpMode = iota
	gte
)

// NewStopLossRule returns a new rule that is satisfied when the given loss tolerance (a percentage) is met or exceeded.
// Loss tolerance should be a value between -1 and 1.
func NewStopLossRule(series *TimeSeries, lossTolerance float64) Rule {
	return stopLossOrGainRule{
		Indicator: NewClosePriceIndicator(series),
		tolerance: big.NewDecimal(lossTolerance),
		mode:      lte,
	}
}
func NewStopGainRule(series *TimeSeries, lossTolerance float64) Rule {
	return stopLossOrGainRule{
		Indicator: NewClosePriceIndicator(series),
		tolerance: big.NewDecimal(lossTolerance),
		mode:      gte,
	}
}

func (slr stopLossOrGainRule) IsSatisfied(index int, record *TradingRecord) bool {
	if !record.CurrentPosition().IsOpen() {
		return false
	}

	openPrice := record.CurrentPosition().EntranceOrder().Price
	loss := slr.Indicator.Calculate(index).Div(openPrice).Sub(big.ONE)
	if slr.mode == lte {
		return loss.LTE(slr.tolerance)
	} else {
		return loss.GTE(slr.tolerance)
	}

}
