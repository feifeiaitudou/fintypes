package frame

import (
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/fintypes/comm"
	"time"
)

type (
	KDot struct {
		Time       time.Time        `json:"Time" bson:"_id" csv:"Time"`                                // statistic begin time
		Open       decimals.Decimal `json:"Open,omitempty" bson:"Open,omitempty" csv:"Open,omitempty"` // open price in USD
		Low        decimals.Decimal `json:"Low,omitempty" bson:"Low,omitempty" csv:"Low,omitempty"`
		High       decimals.Decimal `json:"High,omitempty" bson:"High,omitempty" csv:"High,omitempty"`
		Close      decimals.Decimal `json:"Close,omitempty" bson:"Close,omitempty" csv:"Close,omitempty"`
		Volume     decimals.Decimal `json:"Volume,omitempty" bson:"Volume,omitempty" csv:"Volume,omitempty"` // volume in quote asset, always, it is fiat in stock, it is USD(s)/BTC/ETH... in crypto currency
		indicators map[string]float64
	}

	// KDot items for Global
	Kline struct {
		Pair   comm.PairExt
		Period comm.Period
		Items  []KDot
		sorted bool
	}
)

func (k *Kline) Len() int {
	return len(k.Items)
}
