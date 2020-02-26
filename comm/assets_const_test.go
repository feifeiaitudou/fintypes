package comm

import (
	"fmt"
	"github.com/shawnwyckoff/commpkg/apputil/test"
	"testing"
)

func TestAllFiats(t *testing.T) {
	list := AllFiats()
	fmt.Println(len(list))
	fmt.Println(list)
}

func TestStableCoinsByFiat(t *testing.T) {
	coins := StableCoinsByFiat(USD)
	if len(coins) < 3 {
		test.PrintlnExit(t, "StableCoinsByFiat error")
	}
}
