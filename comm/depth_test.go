package comm

import (
	"encoding/json"
	"github.com/shawnwyckoff/commpkg/apputil/test"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/dsa/jsons"
	"testing"
	"time"
)

// example: LTC/USDT
func newTestDepth(sort bool) Depth {
	r := &Depth{}
	r.Sells = append(r.Sells, OrderBook{Price: decimals.NewFromInt(13), Amount: decimals.NewFromInt(1)})
	r.Sells = append(r.Sells, OrderBook{Price: decimals.NewFromInt(12), Amount: decimals.NewFromInt(1)})
	r.Sells = append(r.Sells, OrderBook{Price: decimals.NewFromInt(11), Amount: decimals.NewFromInt(1)})
	r.Buys = append(r.Buys, OrderBook{Price: decimals.NewFromInt(10), Amount: decimals.NewFromInt(1)})
	r.Buys = append(r.Buys, OrderBook{Price: decimals.NewFromInt(8), Amount: decimals.NewFromInt(1)})
	r.Buys = append(r.Buys, OrderBook{Price: decimals.NewFromInt(9), Amount: decimals.NewFromInt(1)})
	if sort {
		r.Sort()
	}
	return *r
}

func TestDepth_MarketBuyDetect(t *testing.T) {
	d := newTestDepth(true)

	dt := d.MarketBuyDetect(decimals.NewFromInt(11))
	expected := DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(1),
		QuoteDealAmount: decimals.NewFromInt(11),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(11),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetect error1")
		return
	}

	dt = d.MarketBuyDetect(decimals.NewFromInt(23))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(23),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(12),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetect error2")
		return
	}

	dt = d.MarketBuyDetect(decimals.NewFromInt(36))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(3),
		QuoteDealAmount: decimals.NewFromInt(36),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(13),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetect error3")
		return
	}

	dt = d.MarketBuyDetect(decimals.NewFromInt(37))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(3),
		QuoteDealAmount: decimals.NewFromInt(36),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      PricePiercing,
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetect error4")
		return
	}
}

func TestDepth_MarketSellDetect(t *testing.T) {
	d := newTestDepth(true)

	dt := d.MarketSellDetect(decimals.NewFromInt(1))
	expected := DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(1),
		QuoteDealAmount: decimals.NewFromInt(10),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(10),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetect error1")
		return
	}

	dt = d.MarketSellDetect(decimals.NewFromInt(2))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(19),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(9),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetect error2")
		return
	}

	dt = d.MarketSellDetect(decimals.NewFromInt(3))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(3),
		QuoteDealAmount: decimals.NewFromInt(27),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(8),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetect error3")
		return
	}

	dt = d.MarketSellDetect(decimals.NewFromInt(4))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(3),
		QuoteDealAmount: decimals.NewFromInt(27),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      PricePiercing,
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetect error4")
		return
	}
}

func TestDepth_MarketBuyDetectEx(t *testing.T) {
	d := newTestDepth(true)

	dt := d.MarketBuyDetectEx(decimals.NewFromInt(11), decimals.NewFromFloat64(0.1), 8, decimals.NewFromInt(1))
	expected := DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(1),
		QuoteDealAmount: decimals.NewFromInt(11),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(11),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetectEx error1")
		return
	}

	dt = d.MarketBuyDetectEx(decimals.NewFromInt(23), decimals.NewFromFloat64(0.1), 8, decimals.NewFromInt(1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(23),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(12),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetectEx error2")
		return
	}

	dt = d.MarketBuyDetectEx(decimals.NewFromInt(36), decimals.NewFromFloat64(0.1), 8, decimals.NewFromInt(1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(23),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(12),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetectEx error3")
		return
	}

	dt = d.MarketBuyDetectEx(decimals.NewFromInt(37), decimals.NewFromFloat64(0.1), 8, decimals.NewFromInt(1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(23),
		Price1:          decimals.NewFromInt(11),
		PriceReach:      decimals.NewFromInt(12),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketBuyDetectEx error4")
		return
	}

	depthString := `{"Time":"2019-02-18T11:44:00Z","Sells":[{"Price":"0.001436","Amount":"6502773"},{"Price":"0.001436","Amount":"6502773"}],"Buys":[{"Price":"0.0014276","Amount":"6502773"},{"Price":"0.001426","Amount":"6502773"}]}`
	if err := json.Unmarshal([]byte(depthString), &d); err != nil {
		t.Error(err)
		return
	}
	dt = d.MarketBuyDetectEx(decimals.NewFromInt(10000), decimals.NewFromFloat64(0.1), 22, decimals.NewFromFloat64(0.000001))
	// dt.UnitDealAmount.Mul(dt.PriceReach).String() == "10000.0000000000000000000308"
	// 必须带上精度，否则算下来要大于原始输入quoteAmount了
	if dt.UnitDealAmount.Mul(dt.PriceReach).WithPrec(8).EqualInt(10000) == false {
		test.PrintlnExit(t, "DetectedQuote %s(%s x %s) != Original Quote Amount 10000", dt.UnitDealAmount.Mul(dt.PriceReach).String(), dt.UnitDealAmount.String(), dt.PriceReach.String())
	}

	precision := 8
	lot := decimals.NewFromInt(1)
	dt = d.MarketBuyDetectEx(decimals.NewFromInt(10000), decimals.NewFromFloat64(0.1), precision, lot)
	if dt.UnitDealAmount.Mul(dt.PriceReach).WithPrec(precision).String() != "9999.999568" {
		t.Errorf("MarketBuyDetectEx2 error2")
	}

	// 下面这个测试用例涉及精度问题
	d = Depth{
		Time: time.Now(),
		DepthRawData: DepthRawData{
			Sells: OrderBookList{OrderBook{Price: decimals.NewFromFloat64(6273.03627303), Amount: decimals.NewFromInt(10000000000000)}},
			Buys:  OrderBookList{OrderBook{Price: decimals.NewFromFloat64(6273.02372697), Amount: decimals.NewFromInt(10000000000000)}},
		},
	}
	quoteAmount, err := decimals.NewFromString("5141.73808735821834851056117427")
	test.Assert(t, err)
	precision = 8
	lot = decimals.NewFromFloat64(0.000001)
	dt = d.MarketBuyDetectEx(quoteAmount, decimals.NewFromFloat64(0.01), precision, lot)
	if dt.UnitDealAmount.Mul(dt.PriceReach).Trunc(precision, lot.Float64()).String() != "5141.73181899" {
		test.PrintlnExit(t, "DetectedQuote %s(%s x %s) != Expected Quote Amount 5141.73181941", dt.UnitDealAmount.Mul(dt.PriceReach).String(), dt.UnitDealAmount.String(), dt.PriceReach.String())
	}
	if dt.UnitDealAmount.Mul(dt.PriceReach).GreaterThan(quoteAmount) {
		test.PrintlnExit(t, "DetectedQuote %s(%s x %s) > Original Quote Amount %s", dt.UnitDealAmount.Mul(dt.PriceReach).String(), dt.UnitDealAmount.String(), dt.PriceReach.String(), quoteAmount)
	}
}

func TestDepth_MarketSellDetectEx(t *testing.T) {
	d := newTestDepth(true)

	dt := d.MarketSellDetectEx(decimals.NewFromInt(1), decimals.NewFromFloat64(0.1))
	expected := DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(1),
		QuoteDealAmount: decimals.NewFromInt(10),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(10),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetectEx error1")
		return
	}

	dt = d.MarketSellDetectEx(decimals.NewFromInt(2), decimals.NewFromFloat64(0.1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(19),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(9),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetectEx error2")
		return
	}

	dt = d.MarketSellDetectEx(decimals.NewFromInt(3), decimals.NewFromFloat64(0.1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(19),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(9),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetectEx error3")
		return
	}

	dt = d.MarketSellDetectEx(decimals.NewFromInt(4), decimals.NewFromFloat64(0.1))
	expected = DepthTolerance{
		UnitDealAmount:  decimals.NewFromInt(2),
		QuoteDealAmount: decimals.NewFromInt(19),
		Price1:          decimals.NewFromInt(10),
		PriceReach:      decimals.NewFromInt(9),
	}
	if !dt.Equal(expected) {
		t.Errorf("MarketSellDetectEx error4")
		return
	}
}

func TestDepth_Sort(t *testing.T) {
	withOrder := newTestDepth(true)
	outOfOrder := newTestDepth(false)
	outOfOrder.Sort()
	if jsons.MarshalStringDefault(withOrder, false) != jsons.MarshalStringDefault(outOfOrder, false) {
		test.PrintlnExit(t, "2 depths should equal")
	}
}
