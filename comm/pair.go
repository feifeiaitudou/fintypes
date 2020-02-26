package comm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/commpkg/apputil/errorz"
	"github.com/shawnwyckoff/commpkg/dsa/stringz"
	"sort"
	"strings"
)

/*
Pair ISO standard format: Unit/Quote
Unit: Base Currency / Unit Currency (amount is 1 when exchange these two currencies)
Quote: Quote Currency / Pricing Currency / Secondary Currency
Example: USD/CNY = 6.6
means 1USD = 6.6CNY, 以USD为基准或者1个单位，需要6.6CNY的报价，所以USD是基准货币，CNY是报价货币

Notice:
ok、binance、hitbtc... has correct Unit / Quote order, but Bittrex is not, it's order is quote/unit.
*/

// TODO test PairExt
type (
	Pair    string // BTC/USDT
	PairExt string // BTC/USD.1min.swap.okex BTCG2020/USD.1min.future.CME BTC/USDT.1min.spot.binance BTC/USDT.1day.sse BTC/USDT.binance
)

const (
	PairErr          = Pair("")
	PairExtErr       = PairExt("")
	pairExtDelimiter = "."
)

var (
	peCache             = map[PairExt]PairExt{}
	commonPairDelimiter = []string{"/", "-", "_", "|", "+", ":", "#", ",", ".", "\\", "*"}
)

func parsePairWithOptions(s string, delimiters []string, leftTailDelimiters, rightHeaderDelimiters []string, ISOOrder bool) (pair Pair, unit, quote string, err error) {
	if s == "" {
		return Pair(""), "", "", errors.Errorf("nil input pair")
	}
	if len(delimiters) == 0 {
		delimiters = commonPairDelimiter
	}
	s = strings.Replace(s, " ", "", -1)
	s = strings.ToUpper(s)
	first := ""
	second := ""

	// Like LTC-ETH...
	for _, pd := range delimiters {
		if pd == "" {
			continue
		}
		ss := strings.Split(s, pd)
		if len(ss) != 2 {
			continue
		}
		first = ss[0]
		second = ss[1]
		break
	}

	//specialDelimiters := []string{"BULL"}

	// special delimiters which will append unit, on binance there are pairs like XRPBULL/BUSD ETHBULL/USDT
	// 需要跟在左侧Unit的特殊分隔符，币安中存在XRPBULL/BUSD ETHBULL/USDT等交易对
	for _, sd := range leftTailDelimiters {
		ss := strings.Split(s, sd)
		ss = stringz.RemoveByValue(ss, "") // this is absolutely necessary, if s start with "BULL", it is an empty string at 0 of ss
		if len(ss) != 2 {
			continue
		}
		first = ss[0] + sd
		second = ss[1]
		break
	}

	// special delimiters which will be quote header, no examples for now
	// 需要在右侧Quote头部的特殊分隔符，目前还不需要
	for _, sd := range rightHeaderDelimiters {
		ss := strings.Split(s, sd)
		ss = stringz.RemoveByValue(ss, "") // this is absolutely necessary, if s start with "BULL", it is an empty string at 0 of ss
		if len(ss) != 2 {
			continue
		}
		first = ss[0]
		second = sd + ss[1]
		break
	}

	// Like LTCETH...
	if first == "" || second == "" {
		for _, v := range AllQuoteAssets() {
			/*quoteSymbol := v.Symbol()
			/*if v.TypeSide() == AssetTypeFiat {
				quoteSymbol = v.Name()
			}*/
			if stringz.StartWith(s, v.TradeSymbol()) {
				first = v.TradeSymbol()
				second = stringz.RemoveHead(s, len(v.TradeSymbol()))
				break
			}

			if stringz.EndWith(s, v.TradeSymbol()) {
				first = stringz.RemoveTail(s, len(v.TradeSymbol()))
				second = v.TradeSymbol()
				break
			}
		}
	}

	if first == "" || second == "" {
		return Pair(""), "", "", errors.Errorf("invalid trade pair %s", s)
	}

	if ISOOrder {
		return Pair(fmt.Sprintf("%s/%s", strings.ToUpper(first), strings.ToUpper(second))), strings.ToUpper(first), strings.ToUpper(second), nil
	} else {
		return Pair(fmt.Sprintf("%s/%s", strings.ToUpper(second), strings.ToUpper(first))), strings.ToUpper(second), strings.ToUpper(first), nil

	}
}

// parse ISO standard pair string
func ParsePair(s string) (Pair, error) {
	pair, _, _, err := parsePairWithOptions(s, []string{"/"}, nil, nil, true)
	return pair, err
}

func ParsePairCustom(s string, config *ExConfig) (Pair, error) {
	if config != nil {
		if pair, _, _, err := parsePairWithOptions(s, []string{config.PairDelimiter}, config.PairDelimiterLeftTail, config.PairDelimiterRightHead, config.PairNormalOrder); err != nil {
			return Pair(""), err
		} else {
			return pair, nil
		}
	} else {
		if pair, _, _, err := parsePairWithOptions(s, nil, nil, nil, true); err != nil {
			return Pair(""), err
		} else {
			return pair, nil
		}
	}
}

func NewPair(unit, quote string) Pair {
	if unit == "" || quote == "" {
		return PairErr
	}
	p, err := ParsePair(fmt.Sprintf("%s/%s", strings.ToUpper(unit), strings.ToUpper(quote)))
	if err != nil {
		return PairErr
	}
	return p
}

func (p Pair) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.FormatISO() + `"`), nil
}

func (p Pair) FormatISO() string {
	return p.Format("/", true)
}

func (p Pair) String() string {
	return p.FormatISO()
}

func (p Pair) Format(split string, upper bool) string {
	_, unit, quote, err := parsePairWithOptions(string(p), nil, nil, nil, true)
	if err != nil {
		return ""
	}
	if upper {
		unit = strings.ToUpper(unit)
		quote = strings.ToUpper(quote)
		split = strings.ToUpper(split)
	} else {
		unit = strings.ToLower(unit)
		quote = strings.ToLower(quote)
		split = strings.ToLower(split)
	}
	return fmt.Sprintf("%s%s%s", unit, split, quote)
}

func (p Pair) First() string {
	return p.Unit()
}

func (p Pair) Second() string {
	return p.Quote()
}

// get first asset in ISO trade pair format
func (p Pair) Unit() string {
	_, unit, _, err := parsePairWithOptions(string(p), nil, nil, nil, true)
	if err != nil {
		return ""
	}
	return strings.ToUpper(unit)
}

// get first asset in ISO trade pair format
func (p Pair) Quote() string {
	_, _, quote, err := parsePairWithOptions(string(p), nil, nil, nil, true)
	if err != nil {
		return ""
	}
	return strings.ToUpper(quote)
}

func (p Pair) CustomFormat(config *ExConfig) string {
	delimiter, normalOrder, upperCase := config.PairDelimiter, config.PairNormalOrder, config.PairUpperCase
	first := ""
	second := ""
	if normalOrder {
		if upperCase {
			first = strings.ToUpper(p.Unit())
			second = strings.ToUpper(p.Quote())
			delimiter = strings.ToUpper(delimiter)
		} else {
			first = strings.ToLower(p.Unit())
			second = strings.ToLower(p.Quote())
			delimiter = strings.ToLower(delimiter)
		}
	} else {
		if upperCase {
			first = strings.ToUpper(p.Quote())
			second = strings.ToUpper(p.Unit())
			delimiter = strings.ToUpper(delimiter)
		} else {
			first = strings.ToLower(p.Quote())
			second = strings.ToLower(p.Unit())
			delimiter = strings.ToLower(delimiter)
		}
	}
	s := fmt.Sprintf("%s%s%s", first, delimiter, second)
	if config.PairUpperCase {
		return strings.ToUpper(s)
	} else {
		return strings.ToLower(s)
	}
}

func (p Pair) Verify() error {
	if _, _, _, err := parsePairWithOptions(string(p), nil, nil, nil, true); err != nil {
		return err
	}
	return nil
}

func (p Pair) SetPeriod(period Period) PairExt {
	return PairExt(fmt.Sprintf("%s%s%s", p.String(), pairExtDelimiter, period.String()))
}

func (p Pair) SetPlatform(platform Platform) PairExt {
	return PairExt(fmt.Sprintf("%s%s%s", p.String(), pairExtDelimiter, platform.String()))
}

func (p Pair) SetMarket(market Market) PairExt {
	return PairExt(fmt.Sprintf("%s%s%s", p.String(), pairExtDelimiter, market))
}

// TODO 也许需要提高性能
func ParsePairExtString(s string) (Pair, *Period, *Market, *Platform, error) {
	defErr := errorz.Errorf("invalid PairExt(%s)", s)

	ss := strings.Split(s, pairExtDelimiter)
	if len(ss) <= 0 || len(ss) >= 5 {
		return PairErr, nil, nil, nil, defErr
	}
	if err := Pair(ss[0]).Verify(); err != nil {
		return PairErr, nil, nil, nil, defErr
	}

	var resPeriod *Period = nil
	var resMarket *Market = nil
	var resPlatform *Platform = nil
	for i := 1; i < len(ss); i++ {
		period, err := ParsePeriod(ss[i])
		if err == nil {
			if resPeriod != nil { // 重复出现了，这是异常
				return PairErr, nil, nil, nil, defErr
			} else {
				resPeriod = &period
			}
		}

		market, err := ParseMarket(ss[i])
		if err != nil {
			if resMarket != nil { // 重复出现了，这是异常
				return PairErr, nil, nil, nil, defErr
			} else {
				resMarket = &market
			}
		}

		platform, err := ParsePlatform(ss[i])
		if err == nil {
			if resPlatform != nil { // 重复出现了，这是异常
				return PairErr, nil, nil, nil, defErr
			} else {
				resPlatform = &platform
			}
		}
	}

	return Pair(ss[0]), resPeriod, resMarket, resPlatform, nil
}

func ParsePairExt(s string) (PairExt, error) {
	pair, period, market, platform, err := ParsePairExtString(s)
	if err != nil {
		return PairExtErr, err
	}
	return NewPairExt(pair, period, market, platform), nil
}

func NewPairExt(p Pair, period *Period, market *Market, platform *Platform) PairExt {
	if period == nil && market == nil && platform == nil {
		return PairExtErr
	}
	// 虽然都是有效指针，但指向的内容全部是错误代码，也是不允许的
	if (period != nil && *period == PeriodError) && (market != nil && *market == MarketError) && (platform != nil && *platform == PlatformUnknown) {
		return PairExtErr
	}
	s := p.String()
	if period != nil {
		s += pairExtDelimiter + period.String()
	}
	if market != nil {
		s += pairExtDelimiter + string(*market)
	}
	if platform != nil {
		s += pairExtDelimiter + platform.String()
	}
	return PairExt(s)
}

func (pe PairExt) Verify() error {
	_, _, _, _, err := ParsePairExtString(string(pe))
	if err != nil {
		return err
	}
	return nil
}

func (pe PairExt) Complete() bool {
	_, period, market, platform, err := ParsePairExtString(string(pe))
	return err == nil && period != nil && market != nil && platform != nil
}

func (pe PairExt) HasPeriod() bool {
	_, period, _, _, err := ParsePairExtString(string(pe))
	return err == nil && period != nil
}

func (pe PairExt) HasMarket() bool {
	_, _, market, _, err := ParsePairExtString(string(pe))
	return err == nil && market != nil
}

func (pe PairExt) HasPlatform() bool {
	_, _, _, platform, err := ParsePairExtString(string(pe))
	return err == nil && platform != nil
}

func (pe PairExt) SetPeriod(newPeriod Period) PairExt {
	pair, period, market, platform, err := ParsePairExtString(string(pe))
	if err != nil {
		return PairExtErr
	}
	period = &newPeriod
	return NewPairExt(pair, period, market, platform)
}

func (pe PairExt) SetPlatform(newPlatform Platform) PairExt {
	pair, period, market, platform, err := ParsePairExtString(string(pe))
	if err != nil {
		return PairExtErr
	}
	platform = &newPlatform
	return NewPairExt(pair, period, market, platform)
}

func (pe PairExt) SetMarket(newMarket Market) PairExt {
	pair, period, market, platform, err := ParsePairExtString(string(pe))
	if err != nil {
		return PairExtErr
	}
	market = &newMarket
	return NewPairExt(pair, period, market, platform)
}

func (pe PairExt) Pair() Pair {
	pair, _, _, _, err := ParsePairExtString(string(pe))
	if err != nil {
		return PairErr
	}
	return pair
}

func (pe PairExt) Period() Period {
	_, period, _, _, err := ParsePairExtString(string(pe))
	if err != nil || period == nil {
		return PeriodError
	}
	return *period
}

func (pe PairExt) Platform() Platform {
	_, _, _, platform, err := ParsePairExtString(string(pe))
	if err != nil || platform == nil {
		return PlatformUnknown
	}
	return *platform
}

func (pe PairExt) Market() Market {
	_, _, market, _, err := ParsePairExtString(string(pe))
	if err != nil || market == nil {
		return MarketError
	}
	return *market
}

func (pe PairExt) PairMarket() PairExt {
	pair, _, market, _, err := ParsePairExtString(string(pe))
	if err != nil || market == nil {
		return PairExtErr
	}
	return pair.SetMarket(*market)
}

func (pe PairExt) String() string {
	return string(pe)
}

func PairExtsSort(src []PairExt) {
	var ss []string
	for _, v := range src {
		ss = append(ss, v.String())
	}
	sort.Strings(ss)

	src = nil
	for _, v := range ss {
		src = append(src, PairExt(v))
	}
}

func PairExtsInclude(src []PairExt, find PairExt) bool {
	for _, v := range src {
		if v == find {
			return true
		}
	}
	return false
}

func PairExtsEqual(a []PairExt, b []PairExt) bool {
	var sa []string
	var sb []string
	for _, v := range a {
		sa = append(sa, v.String())
	}
	for _, v := range b {
		sb = append(sb, v.String())
	}
	sort.Strings(sa)
	sort.Strings(sb)
	return strings.Join(sa, ",") == strings.Join(sb, ",")
}

// find same Pairs between exchanges
func FindSamePairs(pairs map[Platform][]Pair) map[Pair][]Platform {
	result := make(map[Pair][]Platform)
	for scanEx, scanExPairs := range pairs {
		for _, scanPair := range scanExPairs {
			existedExs := result[scanPair]
			existedExs = append(existedExs, scanEx)
			result[scanPair] = existedExs
		}
	}

	for pair, platforms := range result {
		if len(platforms) <= 1 {
			delete(result, pair)
		}
	}
	return result
}
