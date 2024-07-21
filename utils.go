package main

import "math"

func compare_float(a, b float64) bool {
	epsilon := 0.00001
	return math.Abs(a-b) < epsilon
}

func compare_tuple(a, b *tuple) bool {
	return compare_float(a.x, b.x) && compare_float(a.y, b.y) && compare_float(a.z, b.z) && compare_float(a.w, b.w)
}
