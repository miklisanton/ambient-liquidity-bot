package utils

import (
	"math"
	"testing"
)

func TestTickToPrice(t *testing.T) {
    price := 3473.07
    dec0 := 18
    dec1 := 6
    tick := PriceToTick(price, dec0, dec1)
    if math.Abs(TickToPrice(tick, dec0, dec1) - price) > 0.1 {
        t.Fatalf("error converting tick to price tick calculated = %d", tick)
    }
    t.Logf("tick: %d, price: %f", tick, price)

}

