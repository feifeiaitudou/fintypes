package comm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/dsa/stringz"
	"strings"
)

const (
	AccountAddressNull AccountAddress = ""
)

type (
	// sub account
	// total = free + locked
	// net = free + locked - borrowed - interest
	Balance struct {
		Free     decimals.Decimal `json:"Free"`
		Locked   decimals.Decimal `json:"Locked"`
		Borrowed decimals.Decimal `json:"Borrowed"` // only available in margin account
		Interest decimals.Decimal `json:"Interest"` // only available in margin account
	}

	/*
		huobi的futureAccount成员
		AccountRights float64 //账户权益
		KeepDeposit   float64 //保证金
		ProfitReal    float64 //已实现盈亏
		ProfitUnreal  float64
		RiskRate      float64 //保证金率
	*/
	// 合约账户余额
	ContractBalance struct {
		InitialMargin         decimals.Decimal // 当前的仓位占用保证金，你下好单之后哪怕不操作这个值也会动态变化
		MaintMargin           decimals.Decimal // maintenance margin 维持保证金
		MarginBalance         decimals.Decimal // 保证金余额
		MaxWithdrawAmount     decimals.Decimal // 最大可提现
		PositionInitialMargin decimals.Decimal // 当前的仓位占用保证金
		UnrealizedProfit      decimals.Decimal // 未实现的盈亏
		WalletBalance         decimals.Decimal // 钱包余额
	}

	Account struct {
		Spot   map[string]Balance         `json:"Spot"`   // 现货
		Margin map[string]Balance         `json:"Margin"` // 杠杆
		Future map[string]ContractBalance `json:"Future"` // 期货合约
		Perp   map[string]ContractBalance `json:"Perp"`   // 永续合约
	}

	// example: buffett@gmail.com@binance
	AccountAddress string

	Market string
)

func NewAccountAddress(email string, platform Platform) AccountAddress {
	return AccountAddress(fmt.Sprintf("%s@%s", email, platform.String()))
}

func (tai AccountAddress) String() string {
	return string(tai)
}

func (tai AccountAddress) Email() string {
	ss := strings.Split(string(tai), "@")
	if len(ss) != 3 {
		return ""
	}
	return ss[0] + "@" + ss[1]
}

func (tai AccountAddress) Platform() Platform {
	ss := strings.Split(string(tai), "@")
	if len(ss) != 3 {
		return PlatformUnknown
	}
	return Platform(ss[2])
}

func ParseAccountAddress(s string) (AccountAddress, error) {
	if s == "" {
		return AccountAddressNull, nil
	}
	ss := strings.Split(s, "@")
	if len(ss) != 3 {
		return "", errors.Errorf("invalid AccountAddress(%s)", s)
	}
	_, err := ParsePlatform(ss[2])
	if err != nil {
		return "", err
	}
	return AccountAddress(s), nil
}

func (tai AccountAddress) MarshalJSON() ([]byte, error) {
	return []byte(`"` + tai.String() + `"`), nil
}

func (tai *AccountAddress) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = stringz.RemoveHead(str, 1)
	str = stringz.RemoveTail(str, 1)
	decTAI, err := ParseAccountAddress(str)
	if err != nil {
		return err
	}
	*tai = decTAI
	return nil
}

func NewEmptyAccount() *Account {
	acc := Account{Spot: map[string]Balance{}, Margin: map[string]Balance{}, Future: map[string]ContractBalance{}}
	return &acc
}

func NewTestAccount(assets []string, totalInUSD decimals.Decimal, tickers map[PairExt]decimals.Decimal) (*Account, error) {
	if len(assets) == 0 {
		return nil, errors.Errorf("no assets")
	}

	eachUSD := totalInUSD.DivInt(len(assets))
	r := NewEmptyAccount()

	for _, v := range assets {
		price, err := getPriceInUSD(MarketSpot, v, tickers)
		if err != nil {
			return nil, err
		}
		r.SetSpot(v, eachUSD.Div(price), decimals.Zero)
	}
	return r, nil
}

func (sa Balance) ToRepay() decimals.Decimal {
	return sa.Borrowed.Add(sa.Interest)
}

func (sa Balance) Available() decimals.Decimal {
	return sa.Free
}

func (sa Balance) IsZero() bool {
	return sa.Free.IsZero() && sa.Locked.IsZero() && sa.Borrowed.IsZero() && sa.Interest.IsZero()
}

func (sa Balance) Total() decimals.Decimal {
	return sa.Free.Add(sa.Locked)
}

func (sa Balance) Net() decimals.Decimal {
	return sa.Free.Add(sa.Locked).Sub(sa.Borrowed).Sub(sa.Interest)
}

// TODO: 目前的实现可能不严谨
func (sa ContractBalance) IsZero() bool {
	return sa.InitialMargin.IsZero() && sa.MaintMargin.IsZero() && sa.MarginBalance.IsZero() && sa.WalletBalance.IsZero()
}

// Add function will change content of receiver, so pointer required
func (sa *Balance) Add(toAdd Balance) {
	sa.Free = sa.Free.Add(toAdd.Free)
	sa.Borrowed = sa.Borrowed.Add(toAdd.Borrowed)
	sa.Locked = sa.Locked.Add(toAdd.Locked)
	sa.Interest = sa.Interest.Add(toAdd.Interest)
}

func (a *Account) AssetInTotal(asset string) Balance {
	spotAsset := a.AssetInSpot(asset)
	marginAsset := a.AssetInMargin(asset)
	r := Balance{
		Free:     spotAsset.Free.Add(marginAsset.Free),
		Locked:   spotAsset.Locked.Add(marginAsset.Locked),
		Borrowed: spotAsset.Borrowed.Add(marginAsset.Borrowed),
		Interest: spotAsset.Interest.Add(marginAsset.Interest),
	}
	return r
}

func (a *Account) AssetInSpot(asset string) Balance {
	if sa, ok := a.Spot[asset]; ok {
		return sa
	}
	return Balance{}
}

func (a *Account) AssetInMargin(asset string) Balance {
	if sa, ok := a.Margin[asset]; ok {
		return sa
	}
	return Balance{}
}

func (a *Account) SetSpot(asset string, free, lockedAmount decimals.Decimal) {
	a.Spot[asset] = Balance{Free: free, Locked: lockedAmount}
}

func (a *Account) SetMargin(asset string, free, lockedAmount, borrowedAmount decimals.Decimal) {
	a.Margin[asset] = Balance{Free: free, Locked: lockedAmount, Borrowed: borrowedAmount}
}

func (a *Account) String() string {
	return a.JsonString()
}

func (a *Account) JsonString() string {
	res, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(res)
}

func (a *Account) TextString() string {
	return fmt.Sprintf("%+v", a)
}

func (a *Account) Add(toAdd Account) {
	newAccount := joinAccounts(*a, toAdd)
	*a = *newAccount
}

func (a *Account) TransferToMargin(unit string, transferAmount decimals.Decimal) error {
	if a.Spot[unit].Available().LessThan(transferAmount) {
		return errors.Errorf("available amount(%s) of unit(%s) less than transfer amount(%s)", a.Spot[unit].Available().String(), unit, transferAmount.String())
	}
	spotSA := a.Spot[unit]
	spotSA.Free = spotSA.Free.Sub(transferAmount)
	a.Spot[unit] = spotSA
	marginSA := a.Margin[unit]
	marginSA.Free = marginSA.Free.Add(transferAmount)
	a.Margin[unit] = marginSA
	return nil
}

func (a *Account) TransferToSpot(unit string, transferAmount decimals.Decimal) error {
	if a.Margin[unit].Available().LessThan(transferAmount) {
		return errors.Errorf("available amount(%s) of unit(%s) less than transfer amount(%s)", a.Margin[unit].Available().String(), unit, transferAmount.String())
	}
	marginSA := a.Margin[unit]
	marginSA.Free = marginSA.Free.Sub(transferAmount)
	a.Margin[unit] = marginSA
	spotSA := a.Spot[unit]
	spotSA.Free = spotSA.Free.Add(transferAmount)
	a.Spot[unit] = spotSA
	return nil
}

// TODO 添加合约相关的余额
// 计算的用途是统计回报的相关指标
func (a *Account) TotalInUSD(ticks map[PairExt]decimals.Decimal) (*Balance, error) {
	total := &Balance{}

	for unit, subAcc := range a.Spot {
		price, err := getPriceInUSD(MarketSpot, unit, ticks)
		if err != nil {
			return nil, err
		}
		total.Free = total.Free.Add(subAcc.Free.Mul(price))
		total.Locked = total.Locked.Add(subAcc.Locked.Mul(price))
	}

	for unit, subAcc := range a.Margin {
		price, err := getPriceInUSD(MarketMargin, unit, ticks)
		if err != nil {
			return nil, err
		}
		total.Free = total.Free.Add(subAcc.Free.Mul(price))
		total.Locked = total.Locked.Add(subAcc.Locked.Mul(price))
		total.Borrowed = total.Borrowed.Add(subAcc.Borrowed.Mul(price))
		total.Interest = total.Interest.Add(subAcc.Interest.Mul(price))
	}

	return total, nil
}

func (a *Account) FiatNStableBalance() Balance {
	return Balance{}
}

func joinAccounts(account ...Account) *Account {
	r := NewEmptyAccount()

	// spot
	for _, acc := range account {
		for asset, sa := range acc.Spot {
			existed := r.Spot[asset]
			existed.Free = existed.Free.Add(sa.Free)
			existed.Locked = existed.Locked.Add(sa.Locked)
			r.Spot[asset] = existed
		}
	}

	// margin
	for _, acc := range account {
		for asset, sa := range acc.Margin {
			existed := r.Margin[asset]
			existed.Free = existed.Free.Add(sa.Free)
			existed.Locked = existed.Locked.Add(sa.Locked)
			existed.Borrowed = existed.Borrowed.Add(sa.Borrowed)
			existed.Interest = existed.Interest.Add(sa.Interest)
			r.Margin[asset] = existed
		}
	}

	return r
}

func getPriceInUSD(market Market, unit string, tickers map[PairExt]decimals.Decimal) (decimals.Decimal, error) {
	// "USDT/USD", USDT itself
	if strings.ToUpper(unit) == "USDT" {
		return decimals.One, nil
	}

	var usdCoins []string
	for _, usdCoin := range StableCoinsByFiat(USD) {
		usdCoins = append(usdCoins, usdCoin.Symbol())
	}
	usdCoins = append(usdCoins, "USD")

	// "***/USDT"...
	for _, usdCoin := range usdCoins {
		p := NewPair(unit, usdCoin)
		price, ok := tickers[p.SetMarket(market)]
		if ok {
			return price, nil
		}
	}

	// "***/BTC"...
	for _, quoteCoin := range AllQuoteCoins() {
		p := NewPair(unit, quoteCoin.Symbol())
		unitQuotePrice, ok := tickers[p.SetMarket(market)]
		if ok {
			for _, usdCoin := range usdCoins {
				p2 := NewPair(quoteCoin.Symbol(), usdCoin)
				quoteUsdPrice, ok2 := tickers[p2.SetMarket(market)]
				if ok2 {
					return unitQuotePrice.Mul(quoteUsdPrice), nil
				}
			}
		}
	}

	return decimals.Zero, errors.Errorf("can't get USD balance for %s", unit)
}



/*
func (a *Account) Get(id Market, assetName string) (amount, lockedAmount decimals.Decimal) {
	if id == MarketSpot {
		if sa, ok := a.Spot[assetName]; !ok {
			return decimals.Zero, decimals.Zero
		} else {
			return sa.Amount, sa.Locked
		}
	}
	if id == MarketMargin {
		if sa, ok := a.Margin[assetName]; !ok {
			return decimals.Zero, decimals.Zero
		} else {
			return sa.Amount, sa.Locked
		}
	}
	return decimals.Zero, decimals.Zero
}*/
