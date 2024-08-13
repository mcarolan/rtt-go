package matrix

import (
	"context"
	"errors"
	"fmt"
	"rtt/shared"
	"rtt/tuple"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

var matrixVariableName = `([A-Z]+)`
var tupleVariableName = `([a-z]+)`

func aMatrix(ctx context.Context, size int32, name string, table *godog.Table) (context.Context, error) {
	var m *Matrix = nil

	if size == 4 {
		m = Matrix4()
	} else if size == 2 {
		m = Matrix2()
	} else if size == 3 {
		m = Matrix3()
	} else {
		return ctx, errors.ErrUnsupported
	}

	width := len(table.Rows[0].Cells)
	values := make([]float64, len(table.Rows)*width)

	for x := 0; x < width; x++ {
		for y := 0; y < len(table.Rows); y++ {
			value, _ := strconv.ParseFloat(table.Rows[y].Cells[x].Value, 64)
			values[y*width+x] = value
		}
	}
	m.values = values

	return context.WithValue(ctx, variables{name}, m), nil
}

func aTuple(ctx context.Context, variable string, x, y, z, w float64) (context.Context, error) {
	t := tuple.Tuple{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}

	return context.WithValue(ctx, variables{name: variable}, &t), nil
}

func assertComponent(ctx context.Context, name string, x, y int, value float64) (context.Context, error) {
	m := ctx.Value(variables{name}).(*Matrix)

	c := m.At(x, y)

	if !shared.CompareFloat(c, value) {
		return ctx, fmt.Errorf("Expected %f found %f", value, c)
	}

	return ctx, nil
}

func matrixMultiplication(ctx context.Context, destination, aName, bName string) context.Context {
	a := ctx.Value(variables{name: aName}).(*Matrix)
	b := ctx.Value(variables{name: bName}).(*Matrix)

	m := a.Multiply(b)

	return context.WithValue(ctx, variables{name: destination}, m)
}

func matrixTranspose(ctx context.Context, destination, input string) context.Context {
	a := ctx.Value(variables{name: input}).(*Matrix)

	m := a.Transpose()

	return context.WithValue(ctx, variables{name: destination}, m)
}

func matrixTupleMultiplication(ctx context.Context, destination, aName, bName string) context.Context {
	a := ctx.Value(variables{name: aName}).(*Matrix)
	b := ctx.Value(variables{name: bName}).(*tuple.Tuple)

	t := a.MultiplyTuple(b)

	return context.WithValue(ctx, variables{name: destination}, t)
}

func assignIdMatrix(ctx context.Context, destination string) context.Context {
	return context.WithValue(ctx, variables{name: destination}, Identity)
}

func assertMatrixEquals(ctx context.Context, aName, operator, bName string) (context.Context, error) {
	a := ctx.Value(variables{name: aName}).(*Matrix)
	b := ctx.Value(variables{name: bName}).(*Matrix)

	if len(a.values) != len(b.values) {
		return ctx, fmt.Errorf("Left has %d values, right has %d values!", len(a.values), len(b.values))
	}

	if operator == "=" && !a.Equals(b) {
		return ctx, fmt.Errorf("Expected %s = %s", aName, bName)
	}

	if operator == "!=" && a.Equals(b) {
		return ctx, fmt.Errorf("Expected %s != %s", aName, bName)
	}

	return ctx, nil
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

func assertDeterminant(ctx context.Context, name string, expectedResult float64) (context.Context, error) {
	a := ctx.Value(variables{name}).(*Matrix)

	if !shared.CompareFloat(a.Determinant(), expectedResult) {
		return ctx, fmt.Errorf("Error expected %f actual %f!", expectedResult, a.Determinant())
	}

	return ctx, nil
}

func assertMinor(ctx context.Context, name string, y, x int, expectedResult float64) (context.Context, error) {
	a := ctx.Value(variables{name}).(*Matrix)

	if !shared.CompareFloat(a.Minor(y, x), expectedResult) {
		return ctx, fmt.Errorf("Error expected %f actual %f!", expectedResult, a.Minor(y, x))
	}

	return ctx, nil
}

func assertCofactor(ctx context.Context, name string, y, x int, expectedResult float64) (context.Context, error) {
	a := ctx.Value(variables{name}).(*Matrix)

	if !shared.CompareFloat(a.Cofactor(y, x), expectedResult) {
		return ctx, fmt.Errorf("Error expected %f actual %f!", expectedResult, a.Cofactor(y, x))
	}

	return ctx, nil
}

func assertSubmatrix(ctx context.Context, name string, row, col int, expectedResult string) (context.Context, error) {
	a := ctx.Value(variables{name}).(*Matrix)
	expected := ctx.Value(variables{name: expectedResult}).(*Matrix)

	if !a.Submatrix(row, col).Equals(expected) {
		return ctx, fmt.Errorf("Error %s.submatrix(%d, %d) was not %s", name, row, col, expectedResult)
	}

	return ctx, nil
}

func assertInvertible(ctx context.Context, name, not string) (context.Context, error) {
	a := ctx.Value(variables{name}).(*Matrix)

	if not == "" && !a.IsInvertible() {
		return ctx, fmt.Errorf("Expected %s to be invertible", name)
	}

	if not == "not " && a.IsInvertible() {
		return ctx, fmt.Errorf("Expected %s to not be invertible", name)
	}

	return ctx, nil
}

func assignSubmatrix(ctx context.Context, destination, input string, row, col int) context.Context {
	a := ctx.Value(variables{name: input}).(*Matrix)
	m := a.Submatrix(row, col)
	return context.WithValue(ctx, variables{name: destination}, m)
}

func matrixConstructors(ctx *godog.ScenarioContext) {
	ctx.Step(`^the following (\d+)x(?:\d+) matrix (.+):$`, aMatrix)
	regex := fmt.Sprintf(`^(.+) ← tuple\(%s, %s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aTuple)
	regex = fmt.Sprintf(`^%s ← submatrix\(%s, %s, %s\)$`, matrixVariableName, matrixVariableName, shared.PosInt, shared.PosInt)
	ctx.Step(regex, assignSubmatrix)
}

func matrixAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^%s\[%s,%s\] = %s$`, matrixVariableName, shared.PosInt, shared.PosInt, shared.Decimal)
	ctx.Step(regex, assertComponent)
	regex = fmt.Sprintf(`^%s (=|!=) %s$`, matrixVariableName, matrixVariableName)
	ctx.Step(regex, assertMatrixEquals)
	regex = fmt.Sprintf(`^%s (=|!=) %s$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, assertTupleEquals)
	regex = fmt.Sprintf(`^%s = transpose\(%s\)$`, matrixVariableName, matrixVariableName)
	ctx.Step(regex, matrixTranspose)
	regex = fmt.Sprintf(`^determinant\(%s\) = %s$`, matrixVariableName, shared.Decimal)
	ctx.Step(regex, assertDeterminant)
	regex = fmt.Sprintf(`submatrix\(%s, %s, %s\) = %s`, matrixVariableName, shared.PosInt, shared.PosInt, matrixVariableName)
	ctx.Step(regex, assertSubmatrix)
	regex = fmt.Sprintf(`^minor\(%s, %s, %s\) = %s$`, matrixVariableName, shared.PosInt, shared.PosInt, shared.Decimal)
	ctx.Step(regex, assertMinor)
	regex = fmt.Sprintf(`^cofactor\(%s, %s, %s\) = %s$`, matrixVariableName, shared.PosInt, shared.PosInt, shared.Decimal)
	ctx.Step(regex, assertCofactor)
	regex = fmt.Sprintf(`^%s is (not )?invertible$`, matrixVariableName)
	ctx.Step(regex, assertInvertible)
}

func matrixAssignments(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^%s = %s \* %s$`, matrixVariableName, matrixVariableName, matrixVariableName)
	ctx.Step(regex, matrixMultiplication)
	regex = fmt.Sprintf(`^%s = %s \* %s$`, tupleVariableName, matrixVariableName, tupleVariableName)
	ctx.Step(regex, matrixTupleMultiplication)
	ctx.Step(fmt.Sprintf(`^%s = identity_matrix$`, matrixVariableName), assignIdMatrix)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	matrixConstructors(ctx)
	matrixAssignments(ctx)
	matrixAssertions(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/matrices.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
