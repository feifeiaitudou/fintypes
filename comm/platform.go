package comm

import (
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/commpkg/dsa/stringz"
	"github.com/shawnwyckoff/commpkg/sys/clock"
	"strings"
)

type (
	Platform string

	PlatformInfo struct {
		Support  []AssetType
		OpenDate clock.Date
	}
)

var (
	NASDAQOpenDate, _ = clock.NewDate(1971, 2, 4)
	NYSEOpenDate, _   = clock.NewDate(1817, 3, 8)
	AMEXOpenDate, _   = clock.NewDate(1971, 2, 8) // https://www.loc.gov/rr/business/amex/amex.html
	SSEOpenDate, _    = clock.NewDate(1990, 12, 19)
	SZSEOpenDate, _   = clock.NewDate(1991, 7, 3)
	HKEXOpenDate, _   = clock.NewDate(1891, 2, 3) // https://www.ximalaya.com/shangye/22958651/178058751
	TwOpenDate, _     = clock.NewDate(1962, 2, 9)

	PlatformUnknown = Platform("")
	PlatformOpen    = enrollPlatform("open", PlatformInfo{Support: []AssetType{AssetTypeCoin, AssetTypeStock, AssetTypeIndex, AssetTypeMetal}, OpenDate: 0})  // fake platform for internet finance data like gold price
	PlatformIndex   = enrollPlatform("index", PlatformInfo{Support: []AssetType{AssetTypeCoin, AssetTypeStock, AssetTypeIndex, AssetTypeMetal}, OpenDate: 0}) // fake platform for all indexes

	CryptoCompare = enrollPlatform("CryptoCompare", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	DataHub       = enrollPlatform("DataHub", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	YahooFinance  = enrollPlatform("YahooFinance", PlatformInfo{Support: []AssetType{AssetTypeCoin, AssetTypeStock}, OpenDate: 0})
	GoogleFinance = enrollPlatform("GoogleFinance", PlatformInfo{Support: []AssetType{AssetTypeCoin, AssetTypeStock}, OpenDate: 0})

	Coinbase = enrollPlatform("Coinbase", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Bitfinex = enrollPlatform("Bitfinex", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Bitstamp = enrollPlatform("Bitstamp", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Kraken   = enrollPlatform("Kraken", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Binance  = enrollPlatform("Binance", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Bittrex  = enrollPlatform("Bittrex", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})
	Gemini   = enrollPlatform("Gemini", PlatformInfo{Support: []AssetType{AssetTypeCoin}, OpenDate: 0})

	Nasdaq = enrollPlatform("Nasdaq", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: NASDAQOpenDate})
	Nyse   = enrollPlatform("Nyse", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: NYSEOpenDate})
	Amex   = enrollPlatform("Amex", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: AMEXOpenDate}) // belongs to NYSE now
	Szse   = enrollPlatform("Szse", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: SZSEOpenDate}) // Shen Zhen Stock Exchange
	Sse    = enrollPlatform("Sse", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: SSEOpenDate})   // Shanghai Stock Exchange
	Hkex   = enrollPlatform("Hkex", PlatformInfo{Support: []AssetType{AssetTypeStock}, OpenDate: HKEXOpenDate}) // Hong Kong Exchange

	allPlatformInfos = map[Platform]PlatformInfo{}
)

func enrollPlatform(name string, info PlatformInfo) Platform {
	allPlatformInfos[Platform(name)] = info
	return Platform(name)
}

func (p *Platform) Info() PlatformInfo {
	info, ok := allPlatformInfos[*p]
	if !ok {
		return PlatformInfo{}
	}
	return info
}

func (p Platform) String() string {
	return string(p)
}

// Check whether supported exchange or n ot.
func (p Platform) IsSupported() bool {
	for _, v := range AllSupportedExs {
		if v.String() == strings.ToLower(p.String()) {
			return true
		}
	}
	return false
}

func (p Platform) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}

func (p *Platform) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = stringz.RemoveHead(str, 1)
	str = stringz.RemoveTail(str, 1)
	plt, err := ParsePlatform(str)
	if err != nil {
		return err
	}
	*p = plt
	return nil
}

func ParsePlatform(name string) (Platform, error) {
	for plt := range allPlatformInfos {
		if strings.ToLower(name) == strings.ToLower(plt.String()) {
			return plt, nil
		}
	}

	return PlatformUnknown, errors.Errorf("PlatformUnknown platform %s", name)
}

func RemoveDuplicatePlatforms(platforms []Platform) []Platform {
	var ss []string
	for _, v := range platforms {
		ss = append(ss, v.String())
	}
	ss = stringz.RemoveDuplicate(ss)
	var res []Platform
	for _, v := range ss {
		res = append(res, Platform(v))
	}
	return res
}

func AllStockExchanges() []Platform {
	var res []Platform
	for plt, info := range allPlatformInfos {
		for _, v := range info.Support {
			if v == AssetTypeStock {
				res = append(res, plt)
			}
		}
	}
	return res
}
