package comm

import (
	"github.com/shawnwyckoff/commpkg/apputil/errorz"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/dsa/stringz"
)

type (
	// 币安的永续合约说明，可以辅助完善PairInfo
	// https://binance.zendesk.com/hc/zh-cn/articles/360033161972-%E5%90%88%E7%BA%A6%E7%BB%86%E5%88%99
	PairInfo struct {
		// spot/margin only
		//MinNotional decimals.Decimal // notional = price * amount
		Enabled        bool
		UnitPrecision  int              // FIXME precision of price & amount?
		QuotePrecision int              // FIXME precision of price & amount?
		LotMin         decimals.Decimal // mostly it is the same as LostStep
		LotStep        decimals.Decimal // min trade amount in unit, same as amount in Trade()
		// contract only
		MaintMarginPercent    decimals.Decimal
		RequiredMarginPercent decimals.Decimal
	}

	MarketInfo struct {
		Infos map[PairExt]PairInfo
	}
)

// all trading pairs, if some pair exist in different markets(spot,margin...), only one be kept
func (mi *MarketInfo) Pairs() []Pair {
	var res []Pair
	resMap := map[Pair]bool{}
	for pe := range mi.Infos {
		resMap[pe.Pair()] = true
	}
	for k := range resMap {
		res = append(res, k)
	}
	return res
}

func (mi *MarketInfo) PairExts() []PairExt {
	var res []PairExt
	for pe := range mi.Infos {
		res = append(res, pe)
	}

	return res
}

func (mi *MarketInfo) PairExtAts(platform Platform) []PairExt {
	var res []PairExt
	for pe := range mi.Infos {
		res = append(res, pe.SetPlatform(platform))
	}
	return res
}

func (mi *MarketInfo) PairsIncludeFilter(includeSymbols []string, excludeSymbols []string) []Pair {
	includeSymbols = stringz.ToUpper(includeSymbols)
	excludeSymbols = stringz.ToUpper(excludeSymbols)

	var all []Pair
	for _, p := range mi.Pairs() {
		all = append(all, p)
	}

	var res []Pair
	for _, p := range all {
		if len(includeSymbols) > 0 {
			if !stringz.Contains(includeSymbols, p.Quote()) && !stringz.Contains(includeSymbols, p.Unit()) {
				continue
			}
		}

		if len(excludeSymbols) > 0 {
			if stringz.Contains(excludeSymbols, p.Quote()) || stringz.Contains(excludeSymbols, p.Unit()) {
				continue
			}
		}

		res = append(res, p)
	}

	return res
}

func (mi *MarketInfo) PairsAllowedFilter(allowedSymbols []string) []Pair {
	allowedSymbols = stringz.ToUpper(allowedSymbols)
	if len(allowedSymbols) == 0 {
		return nil
	}

	var all []Pair
	for _, p := range mi.Pairs() {
		all = append(all, p)
	}

	var res []Pair
	for _, p := range all {
		if stringz.Contains(allowedSymbols, p.Quote()) && stringz.Contains(allowedSymbols, p.Unit()) {
			res = append(res, p)
		}
	}

	return res
}

func (mi *MarketInfo) AvailableSymbols() []string {
	var res []string
	tmp := map[string]bool{}
	for _, pe := range mi.PairExts() {
		tmp[pe.Pair().Unit()] = true
		tmp[pe.Pair().Quote()] = true
	}
	for symbol := range tmp {
		res = append(res, symbol)
	}
	return res
}

func (mi *MarketInfo) SupportMargin(pair Pair) bool {
	info, ok := mi.Infos[pair.SetMarket(MarketMargin)]
	if !ok {
		return false
	}
	return ok && info.Enabled
}

func (mi *MarketInfo) Verify() error {
	for pe := range mi.Infos {
		if pe.HasPeriod() || pe.HasPlatform() || pe.HasMarket() == false {
			return errorz.Errorf("invalid PairExt(%s)", pe.String())
		}
	}
	return nil
}
