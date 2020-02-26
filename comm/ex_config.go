package comm

import (
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/sys/clock"
	"time"
)

/*
note: all timestamps are UTC + 0 timezone
*/

var (
	ErrFunctionNotSupported = errors.Errorf("function not supported")
	AllSupportedExs         = []Platform{Binance}
)

type (
	ExConfig struct {
		Name                   Platform
		Email                  string // optional
		MaxDepth               int
		PairDelimiter          string // pair separator
		PairDelimiterLeftTail  []string
		PairDelimiterRightHead []string
		PairNormalOrder        bool // whether is ISO order unit first, quote second
		PairUpperCase          bool
		PairsSeparator         string // separator between multiple pairs in api request
		Periods                map[Period]string
		TradeStatus            map[TradeStatus]string
		TradeTypes             map[TradeTypeSide]string
		RateLimit              time.Duration
		FillRateLimit          time.Duration
		KlineRateLimit         time.Duration
		MakerFee               decimals.Decimal
		TakerFee               decimals.Decimal
		WithdrawalFees         map[string]decimals.Decimal
		MarketEnabled          map[Market]bool
		Clock                  clock.Clock
		IsBackTestEx           bool
		TradeBeginTime         time.Time
	}
)

func (cc ExConfig) SupportedPeriods() []Period {
	var r []Period
	for p := range cc.Periods {
		r = append(r, p)
	}
	return r
}

func (cc ExConfig) SpotEnabled() bool {
	enabled, ok := cc.MarketEnabled[MarketSpot]
	return ok && enabled
}

func (cc ExConfig) MarginEnabled() bool {
	enabled, ok := cc.MarketEnabled[MarketMargin]
	return ok && enabled
}

func (cc ExConfig) FutureEnabled() bool {
	enabled, ok := cc.MarketEnabled[MarketFuture]
	return ok && enabled
}

func (cc ExConfig) PerpEnabled() bool {
	enabled, ok := cc.MarketEnabled[MarketPerp]
	return ok && enabled
}

func (cc ExConfig) MinPeriod() Period {
	minPeriod := PeriodError

	for k := range cc.Periods {
		if minPeriod == PeriodError {
			minPeriod = k
		} else {
			if k.ToSeconds() < minPeriod.ToSeconds() {
				minPeriod = k
			}
		}
	}
	return minPeriod
}
