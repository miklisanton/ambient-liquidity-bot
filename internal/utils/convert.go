package utils

import (
    "math"
)

func TickToPrice(tick, dec0, dec1 int) float64 {
    p := math.Pow(1.0001, float64(tick))
    pNormalized := p * math.Pow(10, float64(dec1-dec0))
    return 1 / pNormalized
}

func PriceToTick(price float64, dec0, dec1 int) int {
    pNormalized := 1 / price
    p := pNormalized / math.Pow(10, float64(dec1-dec0))
    return int(math.Log(p) / math.Log(1.0001))
}

