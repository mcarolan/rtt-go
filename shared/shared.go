package shared

import (
	"math"
)

func CompareFloat(a, b float64) bool {
	epsilon := 0.00001
	return math.Abs(a-b) < epsilon
}
