package comm

import (
	"github.com/shawnwyckoff/commpkg/apputil/test"
	"github.com/shawnwyckoff/commpkg/sys/clock"
	"testing"
	"time"
)

func TestDurationToPeriod(t *testing.T) {
	tm1 := time.Date(2019, 8, 1, 1, 15, 0, 0, time.UTC)
	tm2 := time.Date(2019, 8, 1, 2, 15, 0, 0, time.UTC)
	p := DurationToPeriod(tm2.Sub(tm1))
	if p != Period1Hour {
		t.Errorf("should get Period1Hour, but %s got", p.String())
		return
	}

	tm1 = time.Date(2019, 8, 1, 1, 15, 0, 0, time.UTC)
	tm2 = time.Date(2019, 8, 1, 3, 15, 0, 0, time.UTC)
	p = DurationToPeriod(tm2.Sub(tm1))
	if p != Period2Hour {
		t.Errorf("should get Period1Hour, but %s got", p.String())
		return
	}
}

func TestRoundPeriodEarlier(t *testing.T) {
	tm1 := time.Date(2019, 8, 1, 1, 16, 20, 100, time.UTC)
	tm2 := time.Date(2019, 1, 6, 0, 0, 0, 0, time.UTC) // sunday
	tm3 := time.Date(2019, 9, 1, 9, 0, 0, 0, clock.TimeZoneAsiaShanghai)
	tm4 := time.Date(2019, 9, 1, 8, 0, 0, 0, clock.TimeZoneAsiaShanghai)
	tm5 := time.Date(2019, 9, 1, 7, 0, 0, 0, clock.TimeZoneAsiaShanghai)
	cl := test.NewCaseList()
	cl.New().Input(tm1).Input(Period1Min).Expect("2019-08-01 01:16:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period3Min).Expect("2019-08-01 01:15:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period5Min).Expect("2019-08-01 01:15:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period15Min).Expect("2019-08-01 01:15:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period30Min).Expect("2019-08-01 01:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period1Hour).Expect("2019-08-01 01:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period2Hour).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period4Hour).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period6Hour).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period8Hour).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period12Hour).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period1Day).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period1Week).Expect("2019-07-28 00:00:00 +0000 UTC")
	cl.New().Input(tm1).Input(Period1MonthFUZZY).Expect("2019-08-01 00:00:00 +0000 UTC")
	cl.New().Input(tm2).Input(Period1Week).Expect("2018-12-31 00:00:00 +0000 UTC")
	// 下面这3个例子很重要，至于为啥是08:00:00 +0800 CST，可以参考Go标准库time.Round，Round实际上是在UTC基础上的Round
	cl.New().Input(tm3).Input(Period1Day).Expect("2019-09-01 08:00:00 +0800 CST")
	cl.New().Input(tm4).Input(Period1Day).Expect("2019-09-01 08:00:00 +0800 CST")
	cl.New().Input(tm5).Input(Period1Day).Expect("2019-08-31 08:00:00 +0800 CST")

	pc := PeriodRoundConfig{WeekBegin: time.Sunday, UseLocalZeroOClockAsDayBeginning: false}
	for _, v := range cl.Get() {
		inTime := v.Inputs[0].(time.Time)
		inPeriod := v.Inputs[1].(Period)
		expectString := v.Expects[0].(string)
		got := RoundPeriodEarlier(inTime, inPeriod, pc)
		if got.String() != expectString {
			t.Errorf("RoundPeriodEarlier(%s, %s) got %s, but %s expected", inTime, inPeriod, got.String(), expectString)
			return
		}
	}
}
