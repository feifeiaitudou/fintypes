package comm

/*
一个交易有多个属性，包括：

TradeType 市价单或者限价单（必要）
market/limit

TradeTypeSide 做空还是做多（必要）
short/long

TradeIntent这个交易的主观目的（注解）
open/reduce/add/close

TradeIncome这个交易（reduce或者close时）是否盈利（注解）
stop-loss/stop-profit
*/

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/dsa/jsons"
	"github.com/shawnwyckoff/commpkg/dsa/stringz"
	"strings"
	"time"
)

type (
	TradeStatus string

	OrderId string

	TradeTypeSide string // 把TradeType和TradeSide结合到一起，减少了参数个数

	TradeIntent string

	TradeIncome string

	Order struct {
		Id         OrderId
		Time       time.Time
		Pair       Pair
		TypeSide   TradeTypeSide
		Price      decimals.Decimal
		Amount     decimals.Decimal // initial total amount in unit, unit always
		Status     TradeStatus
		AvgPrice   decimals.Decimal // Binance貌似不提供AvgPrice
		DealAmount decimals.Decimal // filled amount in unit, NOT quote,PaperEx在撮合的时候是这么理解的，如果以后要改，也要修正paperEx
		Fee        decimals.Decimal // Binance貌似不提供Fee
	}
)

const (
	OrderIdDelimiter = ":"

	TradeStatusError           TradeStatus = ""
	TradeStatusNew             TradeStatus = "new"
	TradeStatusPartiallyFilled TradeStatus = "partially_filled"
	TradeStatusFilled          TradeStatus = "filled"
	TradeStatusCanceled        TradeStatus = "canceled"
	TradeStatusCanceling       TradeStatus = "canceling"
	TradeStatusRejected        TradeStatus = "rejected"
	TradeStatusExpired         TradeStatus = "expired"

	TradeTypeSideError      TradeTypeSide = ""
	TradeTypeSideLimitSell  TradeTypeSide = "limit-sell"  // 做空
	TradeTypeSideMarketSell TradeTypeSide = "market-sell" // 做空
	TradeTypeSideLimitBuy   TradeTypeSide = "limit-buy"   // 做多
	TradeTypeSideMarketBuy  TradeTypeSide = "market-buy"  // 做多

	TradeIntentError  TradeIntent = ""
	TradeIntentOpen   TradeIntent = "open"   // 开仓进场
	TradeIntentReduce TradeIntent = "reduce" // 减仓
	TradeIntentAdd    TradeIntent = "add"    // 加仓
	TradeIntentClose  TradeIntent = "close"  // 平仓离场

	TradeIncomeError  TradeIncome = ""
	TradeIncomeLoss   TradeIncome = "loss"   // 止损
	TradeIncomeProfit TradeIncome = "profit" // 止盈
)

func (ts TradeStatus) String() string {
	return string(ts)
}

func (ts TradeStatus) End() bool {
	return ts == TradeStatusFilled || ts == TradeStatusCanceled || ts == TradeStatusRejected || ts == TradeStatusExpired
}

func NewOrderId(market Market, pair Pair, strId string) OrderId {
	return OrderId(fmt.Sprintf("%s%s%s%s%s", market, OrderIdDelimiter, pair.String(), OrderIdDelimiter, strId))
}

func (id OrderId) Market() Market {
	ss := strings.Split(string(id), OrderIdDelimiter)
	if len(ss) != 3 {
		return MarketError
	}
	accType := Market(ss[0])
	if accType != MarketSpot && accType != MarketMargin {
		return MarketError
	}
	return accType
}

func (id OrderId) Pair() Pair {
	ss := strings.Split(string(id), OrderIdDelimiter)
	if len(ss) != 3 {
		return PairErr
	}
	p, err := ParsePair(ss[1])
	if err != nil {
		return PairErr
	}
	return p
}

func (id OrderId) StrId() string {
	ss := strings.Split(string(id), OrderIdDelimiter)
	if len(ss) != 3 {
		return ""
	}
	return ss[2]
}

func (id OrderId) Verify() error {
	errInvalidOrderId := errors.Errorf(`invalid OrderId(%s)`, string(id))
	if id.Market() == MarketError {
		return errInvalidOrderId
	}
	if id.Pair() == PairErr {
		return errInvalidOrderId
	}
	if id.StrId() == "" {
		return errInvalidOrderId
	}
	return nil
}

func (id OrderId) String() string {
	return string(id)
}

func (id OrderId) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

func (id *OrderId) UnmarshalJSON(b []byte) error {
	errInvalidOrderId := errors.Errorf(`invalid OrderId(%s)`, string(b))

	s := string(b)
	s = stringz.RemoveHead(s, 1)
	s = stringz.RemoveTail(s, 1)

	oi := OrderId(s)
	if oi.Market() == MarketError {
		return errInvalidOrderId
	}
	if oi.Pair() == PairErr {
		return errInvalidOrderId
	}
	if oi.StrId() == "" {
		return errInvalidOrderId
	}

	*id = oi
	return nil
}

func (od Order) String() string {
	return jsons.MarshalStringDefault(od, false)
}

func (tt TradeTypeSide) String() string {
	return string(tt)
}

func (tt TradeTypeSide) CustomFormat(config *ExConfig) string {
	for k, v := range config.TradeTypes {
		if k == tt {
			return v
		}
	}
	return tt.String()
}

func (tt TradeTypeSide) IsLimit() bool {
	return tt == TradeTypeSideLimitBuy || tt == TradeTypeSideLimitSell
}

func (tt TradeTypeSide) IsMarket() bool {
	return tt == TradeTypeSideMarketBuy || tt == TradeTypeSideMarketSell
}

func (tt TradeTypeSide) IsBuy() bool {
	return tt == TradeTypeSideLimitBuy || tt == TradeTypeSideMarketBuy
}

func (tt TradeTypeSide) IsSell() bool {
	return tt == TradeTypeSideLimitSell || tt == TradeTypeSideMarketSell
}

func (tt TradeTypeSide) Verify() error {
	if tt != TradeTypeSideLimitSell && tt != TradeTypeSideLimitBuy && tt != TradeTypeSideMarketSell && tt != TradeTypeSideMarketBuy {
		return errors.Errorf("invalid Side(%s)", string(tt))
	}
	return nil
}
