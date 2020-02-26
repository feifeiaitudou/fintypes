package comm

import "testing"

func TestFixNameSymbol(t *testing.T) {
	type item struct {
		src      string
		expected string
	}

	items := []item{
		item{src: "globalboost-y", expected: "globalboost-y"},
		item{src: "Alt.Estate token", expected: "alt-estate-token"},
		item{src: "EtherDelta Token", expected: "etherdelta-token"},
		item{src: "I/O Coin", expected: "i-o-coin"},
		item{src: "Carboneum [C8] Token", expected: "carboneum-token"},
		item{src: "BLOC.MONEY", expected: "bloc-money"},
		item{src: "COMSA [XEM]", expected: "comsa"},
		item{src: "  COMSA [XEM]", expected: "comsa"},
		item{src: "  COMSA [XEM]  ", expected: "comsa"},
	}

	for _, v := range items {
		if got := FixCoinName(v.src); got != v.expected {
			t.Errorf("src (%s), expected (%s), but get (%s)", v.src, v.expected, got)
			return
		}
	}
}
