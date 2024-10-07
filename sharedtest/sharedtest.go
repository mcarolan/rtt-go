package sharedtest

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var Decimal = `([0-9\.√\-\/]+)`
var PosInt = `(\d+)`
var MatrixVariableName = `([A-Z]+)`
var TupleVariableName = `([a-z]+[0-9]*)`

var rootDivisionPattern = fmt.Sprintf(`√%s\/(\d+)`, Decimal)
var rootDivision = regexp.MustCompile(rootDivisionPattern)

type Variables struct{ Name string }

func parseDecimal(s string) (float64, error) {
	sign := 1.0
	remaining := s
	if s[0] == '-' {
		sign = -1
		remaining = remaining[1:]
	}

	match := rootDivision.FindStringSubmatch(remaining)
	if match != nil {
		root, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, err
		}

		divisor, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, err
		}
		return (math.Sqrt(root) / divisor) * sign, nil
	}

	f, err := strconv.ParseFloat(remaining, 64)
	if err != nil {
		return 0, err
	}
	return f * sign, nil
}

func ParseXYZ(xString, yString, zString string) (float64, float64, float64, error) {
	x, err := parseDecimal(xString)
	if err != nil {
		return 0, 0, 0, err
	}
	y, err := parseDecimal(yString)
	if err != nil {
		return 0, 0, 0, err
	}
	z, err := parseDecimal(zString)
	if err != nil {
		return 0, 0, 0, err
	}
	return x, y, z, nil
}
