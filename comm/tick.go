package comm

import (
	"encoding/json"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"time"
)

// In some exchange, it won't return High/Low/Volume, even returns Time & Last only, like binance
// They will be filled with -1
type Tick struct {
	Time   time.Time        `json:"time"`
	Last   decimals.Decimal `json:"last"` // 最新成交价, binance has Last only
	Buy    decimals.Decimal `json:"buy"`  // 买一价
	Sell   decimals.Decimal `json:"sell"` // 卖一价
	High   decimals.Decimal `json:"high"` // 最高价
	Low    decimals.Decimal `json:"low"`  // 最低价
	Volume decimals.Decimal `json:"vol"`  // 最近的24小时成交量
}

func (t Tick) String() string {
	buf, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(buf)
}

func TicksToPrices(ticks map[PairExt]Tick) map[PairExt]decimals.Decimal {
	r := make(map[PairExt]decimals.Decimal)
	for pair, tick := range ticks {
		r[pair] = tick.Last
	}
	return r
}
