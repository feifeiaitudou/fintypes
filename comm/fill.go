package comm

import (
	"github.com/shawnwyckoff/commpkg/apputil/errorz"
	"github.com/shawnwyckoff/commpkg/dsa/decimals"
	"github.com/shawnwyckoff/commpkg/sys/clock"
	"time"
)

type (
	Fill struct {
		Id    int64            `json:"Id" bson:"_id"`
		Time  time.Time        `json:"Time" bson:"Time"`
		Price decimals.Decimal `json:"Price" bson:"Price"`
		UnitQty   decimals.Decimal `json:"UnitQty" bson:"UnitQty"`
		Side string `json:"Side" bson:"Side"` // "buy", "sell", "auction"...
	}
	FillOption struct {
		BeginTime    time.Time
		TimeDuration time.Duration

		BeginId int64
		IdLimit int64
	}
)

func (fo FillOption) VerifyBinance() error {
	if clock.BeforeEqual(fo.BeginTime, clock.ZeroTime) &&
		fo.BeginId <= 0 {
		return errorz.Errorf("both beginTime(%s) and beginId[%d] are invalid", fo.BeginTime.String(), fo.BeginId)
	}

	if fo.BeginTime.After(clock.ZeroTime) {
		if clock.BeforeEqual(fo.BeginTime, BTCGenesisBlockTime) {
			return errorz.Errorf("invalid begin time(%s)", fo.BeginTime)
		}
		if fo.TimeDuration <= 0 || fo.TimeDuration > time.Hour {
			return errorz.Errorf("invalid time duration(%s), [1s, 1hour] limited by API", fo.TimeDuration)
		}
	}

	if fo.BeginId > 0 {
		if fo.IdLimit > 1000 {
			return errorz.Errorf("invalid IdLimit[%d], [1, 1000] limited by API", fo.IdLimit)
		}
	}
	return nil
}
