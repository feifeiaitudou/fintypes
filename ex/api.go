package ex

import (
	"github.com/shawnwyckoff/commpkg/apputil/errorz"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/sys/clock"
	. "github.com/shawnwyckoff/fintypes/comm"
	. "github.com/shawnwyckoff/foxs/frame"
	"time"
)

type (
	// Ex interface
	Ex interface {
		// exchange custom settings
		Config() *ExConfig

		// get all supported pairs, min trade amount...
		GetMarketInfo() (*MarketInfo, error)

		// get account info includes all currency balances.
		GetAccount() (*Account, error)

		// get open order books
		GetDepth(market Market, target Pair, limit int) (*Depth, error)

		// get all ticks
		GetTicks() (map[PairExt]Tick, error)

		// get candle bars
		GetKline(market Market, target Pair, period Period, since *time.Time) (*Kline, error)

		// get exchange filled trades history
		GetFills(market Market, target Pair, fromId *int64, limit int) ([]Fill, error)

		// margin account borrowable
		GetBorrowable(asset string) (decimals.Decimal, error)

		// margin account borrow
		Borrow(asset string, amount decimals.Decimal) error

		// margin account repay
		Repay(asset string, amount decimals.Decimal) error

		// transfer between spot and margin account
		Transfer(asset string, amount decimals.Decimal, target Market) error

		// limit-buy, limit-sell, market-buy, market-sell
		// when market-buy/market-sell, price will be ignored
		// amount: always unit amount, not quote amount, whether trade type is buy or sell.
		Trade(market Market, target Pair, t TradeTypeSide, amount, price decimals.Decimal) (*OrderId, error)

		// get all my history orders' info
		GetAllOrders(market Market, target Pair) ([]Order, error)

		// get all my unfinished orders' info
		GetOpenOrders(market Market, target Pair) ([]Order, error)

		// get order info by id
		GetOrder(id OrderId) (*Order, error)

		// cancel unfinished order by id
		CancelOrder(id OrderId) error
	}
)

// email is required in living trading, but not required in kline spider
func NewEx(name Platform, apiKey, apiSecret, proxy string, c clock.Clock, email string) (Ex, error) {
	switch name {
	default:
		return nil, errorz.Errorf("unsupported exchange(%s)", name.String())
	}
}
