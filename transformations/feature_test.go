package transformations

import (
	"context"
	"fmt"
	"math"
	"rtt/matrix"
	"rtt/sharedtest"
	"rtt/tuple"
	"rtt/tupletest"
	"testing"

	"github.com/cucumber/godog"
)

var matrixVariableName = `([A-Z]+)`
var tupleVariableName = `([a-z]+[0-9]*)`

func aTranslation(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	t := Translation(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, t), nil
}

func aScaling(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	t := Scaling(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, t), nil
}

func aShearing(ctx context.Context, variable string, xy, xz, yx, yz, zx, zy float64) (context.Context, error) {
	t := Shearing(xy, xz, yx, yz, zx, zy)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, t), nil
}

func aRotation(ctx context.Context, variable, over string, value float64) (context.Context, error) {
	if over == "x" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, RotationX(math.Pi/value)), nil
	} else if over == "y" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, RotationY(math.Pi/value)), nil
	} else if over == "z" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, RotationZ(math.Pi/value)), nil
	} else {
		return ctx, fmt.Errorf("Unknown component %s", over)
	}
}

func transformationConstructors(sc *godog.ScenarioContext) {
	tupletest.AddConstructPoint(sc)
	tupletest.AddConstructVector(sc)
	regex := fmt.Sprintf(`^(.+) ← translation\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, aTranslation)
	regex = fmt.Sprintf(`^(.+) ← scaling\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, aScaling)
	regex = `^(.+) ← rotation_(.)\(π \/ (\d+)\)$`
	sc.Step(regex, aRotation)
	regex = fmt.Sprintf(`^(.+) ← shearing\(%s, %s, %s, %s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, aShearing)
}

func assertMultiplyComparePoint(ctx context.Context, a, b string, xStr, yStr, zStr string) (context.Context, error) {
	x, y, z, err := sharedtest.ParseXYZ(xStr, yStr, zStr)

	if err != nil {
		return ctx, err
	}

	expected := tuple.Point(x, y, z)

	aMatrix := ctx.Value(sharedtest.Variables{Name: a}).(*matrix.Matrix)
	bPoint := ctx.Value(sharedtest.Variables{Name: b}).(*tuple.Tuple)

	result := aMatrix.MultiplyTuple(bPoint)

	if !tuple.CompareTuple(expected, result) {
		return ctx, fmt.Errorf("%+v was not %+v", result, expected)
	}

	return ctx, nil
}

func assertMultiplyCompareVector(ctx context.Context, a, b string, x, y, z float64) (context.Context, error) {
	expected := tuple.Vector(x, y, z)

	aMatrix := ctx.Value(sharedtest.Variables{Name: a}).(*matrix.Matrix)
	bPoint := ctx.Value(sharedtest.Variables{Name: b}).(*tuple.Tuple)

	result := aMatrix.MultiplyTuple(bPoint)

	if *expected != *result {
		return ctx, fmt.Errorf("%+v was not %+v", result, expected)
	}

	return ctx, nil
}

func assignInverse(ctx context.Context, destination, input string) (context.Context, error) {
	inputMatrix := ctx.Value(sharedtest.Variables{Name: input}).(*matrix.Matrix)

	result, _ := inputMatrix.Invert()

	return context.WithValue(ctx, sharedtest.Variables{Name: destination}, result), nil
}

func transformationAssertions(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) \* (.+) = point\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, assertMultiplyComparePoint)
	regex = fmt.Sprintf(`^(.+) \* (.+) = vector\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, assertMultiplyCompareVector)
}

func matrixTupleMultiplication(ctx context.Context, destination, aName, bName string) context.Context {
	a := ctx.Value(sharedtest.Variables{Name: aName}).(*matrix.Matrix)
	b := ctx.Value(sharedtest.Variables{Name: bName}).(*tuple.Tuple)

	t := a.MultiplyTuple(b)

	return context.WithValue(ctx, sharedtest.Variables{Name: destination}, t)
}

func matrixMultiplication(ctx context.Context, destination, aName, bName, cName string) context.Context {
	a := ctx.Value(sharedtest.Variables{Name: aName}).(*matrix.Matrix)
	b := ctx.Value(sharedtest.Variables{Name: bName}).(*matrix.Matrix)
	c := ctx.Value(sharedtest.Variables{Name: cName}).(*matrix.Matrix)

	t := a.Multiply(b.Multiply(c))

	return context.WithValue(ctx, sharedtest.Variables{Name: destination}, t)
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
	a := ctx.Value(sharedtest.Variables{Name: aName}).(*tuple.Tuple)
	b := ctx.Value(sharedtest.Variables{Name: bName}).(*tuple.Tuple)

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
