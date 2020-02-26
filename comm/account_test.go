package comm

import (
	"fmt"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"testing"
)

func TestNewTestAccount(t *testing.T) {
	ticks := map[PairExt]decimals.Decimal{}
	ticks[PairExt("BTC/USDT.spot")] = decimals.NewFromInt(5000)
	ticks[PairExt("ETH/USDT.spot")] = decimals.NewFromInt(500)
	ticks[PairExt("IOTX/ETH.spot")] = decimals.NewFromFloat64(0.005)
	acc, err := NewTestAccount([]string{"BTC", "USDT", "ETH", "IOTX"}, decimals.NewFromInt(10000), ticks)
	if err != nil {
		t.Error(err)
		return
	}
	if !acc.Spot["BTC"].Free.Equal(decimals.NewFromFloat64(0.5)) {
		t.Errorf("NewTestAccount error1")
		return
	}
	if !acc.Spot["ETH"].Free.Equal(decimals.NewFromInt(5)) {
		t.Errorf("NewTestAccount error2")
		return
	}
	if !acc.Spot["IOTX"].Free.Equal(decimals.NewFromInt(1000)) {
		t.Errorf("NewTestAccount error3")
		return
	}
	if !acc.Spot["USDT"].Free.Equal(decimals.NewFromInt(2500)) {
		t.Errorf("NewTestAccount error4")
		return
	}
}

func TestAccount_TotalInUSD(t *testing.T) {
	ticks := map[PairExt]decimals.Decimal{}
	ticks[PairExt("BTC/USDT.spot")] = decimals.NewFromInt(5000)
	ticks[PairExt("ETH/USDT.spot")] = decimals.NewFromInt(500)
	ticks[PairExt("IOTX/ETH.spot")] = decimals.NewFromFloat64(0.005)
	acc, err := NewTestAccount([]string{"BTC", "USDT", "ETH", "IOTX"}, decimals.NewFromInt(10000), ticks)
	if err != nil {
		t.Error(err)
		return
	}
	total, err := acc.TotalInUSD(ticks)
	if err != nil {
		t.Error(err)
		return
	}
	if total == nil || !total.Free.EqualInt(10000) {
		t.Errorf("TotalInUSD error")
		return
	}
}

func TestBalance_Add(t *testing.T) {
	blc := Balance{Free: decimals.NewFromInt(1000)}
	blc.Add(blc)
	fmt.Println(blc)
}
