package tupletest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"rtt/sharedtest"
	"rtt/tuple"
	"testing"

	"github.com/cucumber/godog"
)

func aTuple(ctx context.Context, variable string, x, y, z, w float64) (context.Context, error) {
	t := tuple.Tuple{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, &t), nil
}

func aNormalized(ctx context.Context, variable, in string) (context.Context, error) {
	t := ctx.Value(sharedtest.Variables{Name: in}).(*tuple.Tuple)
	t = t.Normalize()
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, t), nil
}

func aReflect(ctx context.Context, variable, inName, normalName string) (context.Context, error) {
	in := ctx.Value(sharedtest.Variables{Name: inName}).(*tuple.Tuple)
	normal := ctx.Value(sharedtest.Variables{Name: normalName}).(*tuple.Tuple)
	result := in.Reflect(normal)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func add(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(sharedtest.Variables{Name: left}).(*tuple.Tuple)
	rightTuple := ctx.Value(sharedtest.Variables{Name: right}).(*tuple.Tuple)
	t := leftTuple.Add(rightTuple)
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func subtract(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(sharedtest.Variables{Name: left}).(*tuple.Tuple)
	rightTuple := ctx.Value(sharedtest.Variables{Name: right}).(*tuple.Tuple)
	t := leftTuple.Subtract(rightTuple)
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func mul(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(sharedtest.Variables{Name: left}).(*tuple.Tuple)
	rightTuple := ctx.Value(sharedtest.Variables{Name: right}).(*tuple.Tuple)
	t := leftTuple.Hadamard(rightTuple)
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func negate(ctx context.Context, assignee, in string) (context.Context, error) {
	inTuple := ctx.Value(sharedtest.Variables{Name: in}).(*tuple.Tuple)
	t := inTuple.Negate()
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func scalar_mul(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(sharedtest.Variables{Name: in}).(*tuple.Tuple)
	t := inTuple.ScalarMultiply(scalar)
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func scalar_div(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(sharedtest.Variables{Name: in}).(*tuple.Tuple)
	t := inTuple.ScalarDiv(scalar)
	return context.WithValue(ctx, sharedtest.Variables{Name: assignee}, t), nil
}

func compareMag(ctx context.Context, variable string, sqrt string, expected float64) error {
	actual, ok := ctx.Value(sharedtest.Variables{variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if sqrt == "√" {
		expected = math.Sqrt(expected)
	}

	if actual.Magnitude() != expected {
		return fmt.Errorf("%+v magnitude was %f, not %f", actual, actual.Magnitude(), expected)
	}

	return nil
}

func compareDot(ctx context.Context, left, right string, expected float64) error {
	l := ctx.Value(sharedtest.Variables{Name: left}).(*tuple.Tuple)
	r := ctx.Value(sharedtest.Variables{Name: right}).(*tuple.Tuple)

	if l.Dot(r) != expected {
		return fmt.Errorf("dot(%+v, %+v) was %f, not %f", l, r, l.Dot(r), expected)
	}

	return nil
}

func compareCross(ctx context.Context, left, right, expected string) error {
	l := ctx.Value(sharedtest.Variables{Name: left}).(*tuple.Tuple)
	r := ctx.Value(sharedtest.Variables{Name: right}).(*tuple.Tuple)
	e := ctx.Value(sharedtest.Variables{Name: expected}).(*tuple.Tuple)

	if *l.Cross(r) != *e {
		return fmt.Errorf("cross(%+v, %+v) was %+v, not %+v", l, r, l.Cross(r), e)
	}

	return nil
}

func aComponentEquals(ctx context.Context, variable string, component string, value float64) error {
	tuple, ok := ctx.Value(sharedtest.Variables{variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple [%s] is not set (will check component [%s] for value [%f])", variable, component, value)
	}

	actual := 0.0
	switch component {
	case "x":
		actual = tuple.X
	case "y":
		actual = tuple.Y
	case "z":
		actual = tuple.Z
	case "w":
		actual = tuple.W
	case "red":
		actual = tuple.Red()
	case "green":
		actual = tuple.Green()
	case "blue":
		actual = tuple.Blue()
	default:
		return fmt.Errorf("Unknown component '%s'", component)
	}

	if actual != value {
		return fmt.Errorf("Expected value %f for component '%s', actual was %f", value, component, actual)
	}

	return nil
}

func aPointCheck(ctx context.Context, variable string, notA string) error {
	tuple, ok := ctx.Value(sharedtest.Variables{variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if !tuple.IsPoint() && notA == "" {
		return errors.New("aTuple is not a point")
	}

	if tuple.IsPoint() && notA == "not a" {
		return errors.New("aTuple is a point")
	}

	return nil
}

func aVectorCheck(ctx context.Context, variable string, notA string) error {
	tuple, ok := ctx.Value(sharedtest.Variables{variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if !tuple.IsVector() && notA == "" {
		return errors.New("aTuple is not a vector")
	}

	if tuple.IsVector() && notA == "not a " {
		return errors.New("aTuple is a vector")
	}

	return nil
}

func TupleConstructors(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← tuple\(%s, %s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aTuple)

	AddConstructPoint(ctx)
	AddConstructVector(ctx)
	AddConstructColor(ctx)

	ctx.Step(`^(.+) ← normalize\((.+)\)$`, aNormalized)

	regex = fmt.Sprintf(`^(.+) ← reflect\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aReflect)
}

func TupleAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+)\.(x|y|z|w|red|green|blue) = %s$`, sharedtest.Decimal)
	ctx.Step(regex, aComponentEquals)

	ctx.Step(`^(.) is (not )?a point$`, aPointCheck)
	ctx.Step(`^(.) is (not )?a vector$`, aVectorCheck)

	AddCompareColor(ctx)
	AddComparePoint(ctx)
	AddCompareTuple(ctx)
	AddCompareVector(ctx)
	AddCompareNormalize(ctx)

	regex = fmt.Sprintf(`^magnitude\((.+)\) = (√)?%s$`, sharedtest.Decimal)
	ctx.Step(regex, compareMag)

	regex = fmt.Sprintf(`^dot\((.+), (.+)\) = %s$`, sharedtest.Decimal)
	ctx.Step(regex, compareDot)

	ctx.Step(`^cross\((.+), (.+)\) = (.+)$`, compareCross)
}

func TupleAssignments(ctx *godog.ScenarioContext) {
	ctx.Step(`^(.+) = (.+) \+ (.+)$`, add)
	ctx.Step(`^(.+) = (.+) \- (.+)$`, subtract)
	ctx.Step(`^(.+) = \-(.+)$`, negate)
	regex := fmt.Sprintf(`^(.+) = (.+) \* %s$`, sharedtest.Decimal)
	ctx.Step(regex, scalar_mul)
	ctx.Step(`^(.+) = (.+) \* (.+)$`, mul)
	regex = fmt.Sprintf(`^(.+) = (.+) \/ %s$`, sharedtest.Decimal)
	ctx.Step(regex, scalar_div)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	TupleConstructors(ctx)
	TupleAssertions(ctx)
	TupleAssignments(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/tuples.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
