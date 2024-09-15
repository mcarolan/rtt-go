package transformations

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"rtt/matrix"
	"rtt/shared"
	"rtt/tuple"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

var matrixVariableName = `([A-Z]+)`
var tupleVariableName = `([a-z]+[0-9]*)`

var root2DivisionPattern = `√2\/(\d+)`
var root2Division = regexp.MustCompile(root2DivisionPattern)

var complexDecimal = `([0-9\.√\-\/]+)`

func parseComplexDecimal(s string) (float64, error) {
	sign := 1.0
	remaining := s
	if s[0] == '-' {
		sign = -1
		remaining = remaining[1:]
	}

	match := root2Division.FindStringSubmatch(remaining)
	if match != nil {
		divisor, err := strconv.ParseFloat(match[1], 64)

		if err != nil {
			return 0, err
		}
		return math.Sqrt2 / divisor * sign, nil
	}

	f, err := strconv.ParseFloat(remaining, 64)
	if err != nil {
		return 0, err
	}
	return f * sign, nil
}

func parseComplexXYZ(xString, yString, zString string) (float64, float64, float64, error) {
	x, err := parseComplexDecimal(xString)
	if err != nil {
		return 0, 0, 0, err
	}
	y, err := parseComplexDecimal(yString)
	if err != nil {
		return 0, 0, 0, err
	}
	z, err := parseComplexDecimal(zString)
	if err != nil {
		return 0, 0, 0, err
	}
	return x, y, z, nil
}

func aPoint(ctx context.Context, variable string, xString, yString, zString string) (context.Context, error) {
	x, y, z, err := parseComplexXYZ(xString, yString, zString)
	if err != nil {
		return ctx, err
	}
	p := tuple.Point(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := tuple.Vector(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aTranslation(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	t := Translation(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, t), nil
}

func aScaling(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	t := Scaling(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, t), nil
}

func aShearing(ctx context.Context, variable string, xy, xz, yx, yz, zx, zy float64) (context.Context, error) {
	t := Shearing(xy, xz, yx, yz, zx, zy)
	return context.WithValue(ctx, variables{name: variable}, t), nil
}

func aRotation(ctx context.Context, variable, over string, value float64) (context.Context, error) {
	if over == "x" {
		return context.WithValue(ctx, variables{name: variable}, RotationX(math.Pi/value)), nil
	} else if over == "y" {
		return context.WithValue(ctx, variables{name: variable}, RotationY(math.Pi/value)), nil
	} else if over == "z" {
		return context.WithValue(ctx, variables{name: variable}, RotationZ(math.Pi/value)), nil
	} else {
		return ctx, fmt.Errorf("Unknown component %s", over)
	}
}

func transformationConstructors(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, complexDecimal, complexDecimal, complexDecimal)
	sc.Step(regex, aPoint)
	regex = fmt.Sprintf(`^(.+) ← translation\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	sc.Step(regex, aTranslation)
	regex = fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	sc.Step(regex, aVector)
	regex = fmt.Sprintf(`^(.+) ← scaling\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	sc.Step(regex, aScaling)
	regex = `^(.+) ← rotation_(.)\(π \/ (\d+)\)$`
	sc.Step(regex, aRotation)
	regex = fmt.Sprintf(`^(.+) ← shearing\(%s, %s, %s, %s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	sc.Step(regex, aShearing)
}

func assertMultiplyComparePoint(ctx context.Context, a, b string, xStr, yStr, zStr string) (context.Context, error) {
	x, y, z, err := parseComplexXYZ(xStr, yStr, zStr)

	if err != nil {
		return ctx, err
	}

	expected := tuple.Point(x, y, z)

	aMatrix := ctx.Value(variables{name: a}).(*matrix.Matrix)
	bPoint := ctx.Value(variables{name: b}).(*tuple.Tuple)

	result := aMatrix.MultiplyTuple(bPoint)

	if !tuple.CompareTuple(expected, result) {
		return ctx, fmt.Errorf("%+v was not %+v", result, expected)
	}

	return ctx, nil
}

func assertMultiplyCompareVector(ctx context.Context, a, b string, x, y, z float64) (context.Context, error) {
	expected := tuple.Vector(x, y, z)

	aMatrix := ctx.Value(variables{name: a}).(*matrix.Matrix)
	bPoint := ctx.Value(variables{name: b}).(*tuple.Tuple)

	result := aMatrix.MultiplyTuple(bPoint)

	if *expected != *result {
		return ctx, fmt.Errorf("%+v was not %+v", result, expected)
	}

	return ctx, nil
}

func assignInverse(ctx context.Context, destination, input string) (context.Context, error) {
	inputMatrix := ctx.Value(variables{name: input}).(*matrix.Matrix)

	result, _ := inputMatrix.Invert()

	return context.WithValue(ctx, variables{name: destination}, result), nil
}

func transformationAssertions(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) \* (.+) = point\(%s, %s, %s\)$`, complexDecimal, complexDecimal, complexDecimal)
	sc.Step(regex, assertMultiplyComparePoint)
	regex = fmt.Sprintf(`^(.+) \* (.+) = vector\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	sc.Step(regex, assertMultiplyCompareVector)
}

func matrixTupleMultiplication(ctx context.Context, destination, aName, bName string) context.Context {
	a := ctx.Value(variables{name: aName}).(*matrix.Matrix)
	b := ctx.Value(variables{name: bName}).(*tuple.Tuple)

	t := a.MultiplyTuple(b)

	return context.WithValue(ctx, variables{name: destination}, t)
}

func matrixMultiplication(ctx context.Context, destination, aName, bName, cName string) context.Context {
	a := ctx.Value(variables{name: aName}).(*matrix.Matrix)
	b := ctx.Value(variables{name: bName}).(*matrix.Matrix)
	c := ctx.Value(variables{name: cName}).(*matrix.Matrix)

	t := a.Multiply(b.Multiply(c))

	return context.WithValue(ctx, variables{name: destination}, t)
}

func transformationAssignments(sc *godog.ScenarioContext) {
	regex := `^(.+) = inverse\((.+)\)$`
	sc.Step(regex, assignInverse)
	regex = fmt.Sprintf(`^%s = %s \* %s$`, tupleVariableName, matrixVariableName, tupleVariableName)
	sc.Step(regex, matrixTupleMultiplication)
	regex = fmt.Sprintf(`^%s (=|!=) %s$`, tupleVariableName, tupleVariableName)
	sc.Step(regex, assertTupleEquals)
	regex = fmt.Sprintf(`^(.+) ← %s \* %s \* %s$`, matrixVariableName, matrixVariableName, matrixVariableName)
	sc.Step(regex, matrixMultiplication)
}

func assertTupleEquals(ctx context.Context, aName, operator, bName string) (context.Context, error) {
	a := ctx.Value(variables{name: aName}).(*tuple.Tuple)
	b := ctx.Value(variables{name: bName}).(*tuple.Tuple)

	if operator == "=" && !tuple.CompareTuple(a, b) {
		return ctx, fmt.Errorf("Error %s != %s!", aName, bName)
	}

	if operator == "!=" && tuple.CompareTuple(a, b) {
		return ctx, fmt.Errorf("Error %s = %s!", aName, bName)
	}

	return ctx, nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	transformationConstructors(sc)
	transformationAssignments(sc)
	transformationAssertions(sc)
}

func TestFeature(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/transformations.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
