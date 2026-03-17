package goutil

import "math"

// ToFixed rounds a float64 to the given decimal precision.
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}

// Decimal2 rounds a float64 to 2 decimal places.
func Decimal2(x float64) float64 {
	return math.Round(x*100) / 100
}
