package comm

import "github.com/shawnwyckoff/commpkg/apputil/errorz"

var (
	allMarkets []Market

	MarketError  Market = ""
	MarketSpot          = enrollNewMarket("spot")   // 现货
	MarketMargin        = enrollNewMarket("margin") // 现货杠杆
	MarketFuture        = enrollNewMarket("future") // 交割合约
	MarketPerp          = enrollNewMarket("perp")   // perpetual swap （永续合约）
)

func enrollNewMarket(name string) Market {
	res := Market(name)
	allMarkets = append(allMarkets, res)
	return res
}

func (m Market) Verify() error {
	for _, v := range allMarkets {
		if m == v {
			return nil
		}
	}
	return errorz.Errorf("invalid Market(%s)", m)
}

func ParseMarket(s string) (Market, error) {
	if err := Market(s).Verify(); err != nil {
		return MarketError, err
	}
	return Market(s), nil
}
