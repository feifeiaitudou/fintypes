package comm

var (
	IndexDJI, _  = NewAsset("i.dji")  // Dow Jones Industrial Average (^DJI)
	IndexSP, _   = NewAsset("i.sp")   // S&P 500 (^GSPC)
	IndexIXIC, _ = NewAsset("i.ixic") // NASDAQ Composite (^IXIC)
	IndexNYA, _  = NewAsset("i.nya")  // NYSE Composite (^NYA)
	IndexRUT, _  = NewAsset("i.rut")  // Russell 2000 (^RUT)
	IndexVIX, _  = NewAsset("i.vix")  // CBOE Volatility Index (^VIX)
	IndexSHH, _  = NewAsset("i.shh")  // SSE Composite Index (000001.SS), 上证指数
)

func AllIndexes() []Asset {
	return []Asset{
		IndexDJI,
	}
}

// index is special trade pair
func IndexToPairExt(index Asset) PairExt {
	_, ok := index.ToIndex()
	if !ok {
		return PairExtErr
	}
	switch index {
	case IndexDJI:
		return NewPair("DJI", "USD").SetPlatform(PlatformIndex)
	case IndexSP:
		return NewPair("SP", "USD").SetPlatform(PlatformIndex)
	case IndexIXIC:
		return NewPair("IXIC", "USD").SetPlatform(PlatformIndex)
	case IndexNYA:
		return NewPair("NYA", "USD").SetPlatform(PlatformIndex)
	case IndexRUT:
		return NewPair("RUT", "USD").SetPlatform(PlatformIndex)
	case IndexVIX:
		return NewPair("VIX", "USD").SetPlatform(PlatformIndex)
	case IndexSHH:
		return NewPair("SHH", "CNY").SetPlatform(PlatformIndex)
	}
	return PairExtErr
}
