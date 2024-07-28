package tuple

import (
	"context"
	"errors"
	"fmt"
	"math"
	"rtt/shared"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

func aTuple(ctx context.Context, variable string, x, y, z, w float64) (context.Context, error) {
	t := Tuple{
		x,
		y,
		z,
		w,
	}

	return context.WithValue(ctx, variables{name: variable}, &t), nil
}

func aPoint(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := Point(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := Vector(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aColor(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := Color(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aNormalized(ctx context.Context, variable, in string) (context.Context, error) {
	t := ctx.Value(variables{name: in}).(*Tuple)
	t = t.Normalize()
	return context.WithValue(ctx, variables{name: variable}, t), nil
}

func add(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*Tuple)
	rightTuple := ctx.Value(variables{name: right}).(*Tuple)
	t := leftTuple.Add(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func subtract(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*Tuple)
	rightTuple := ctx.Value(variables{name: right}).(*Tuple)
	t := leftTuple.Subtract(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func mul(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*Tuple)
	rightTuple := ctx.Value(variables{name: right}).(*Tuple)
	t := leftTuple.Hadamard(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func negate(ctx context.Context, assignee, in string) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*Tuple)
	t := inTuple.Negate()
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func scalar_mul(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*Tuple)
	t := inTuple.ScalarMultiply(scalar)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func scalar_div(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*Tuple)
	t := inTuple.ScalarDiv(scalar)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func compareTuple(ctx context.Context, variable string, x, y, z, w float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := Tuple{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}

	if *actual != t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareVector(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := Vector(x, y, z)

	if *actual != *t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func comparePoint(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := Point(x, y, z)

	if *actual != *t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareColor(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := Color(x, y, z)

	if !CompareTuple(actual, t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareMag(ctx context.Context, variable string, sqrt string, expected float64) error {
	actual, ok := ctx.Value(variables{variable}).(*Tuple)

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
	l := ctx.Value(variables{name: left}).(*Tuple)
	r := ctx.Value(variables{name: right}).(*Tuple)

	if l.Dot(r) != expected {
		return fmt.Errorf("dot(%+v, %+v) was %f, not %f", l, r, l.Dot(r), expected)
	}

	return nil
}

func compareCross(ctx context.Context, left, right, expected string) error {
	l := ctx.Value(variables{name: left}).(*Tuple)
	r := ctx.Value(variables{name: right}).(*Tuple)
	e := ctx.Value(variables{name: expected}).(*Tuple)

	if *l.Cross(r) != *e {
		return fmt.Errorf("cross(%+v, %+v) was %+v, not %+v", l, r, l.Cross(r), e)
	}

	return nil
}

func compareNormalize(ctx context.Context, variable string, expected string) error {
	in := ctx.Value(variables{variable}).(*Tuple)
	expectedTuple := ctx.Value(variables{name: expected}).(*Tuple)

	if !CompareTuple(in.Normalize(), expectedTuple) {
		return fmt.Errorf("%+v normalized was %+v, not %+v", in, in.Normalize(), expectedTuple)
	}

	return nil
}

func aComponentEquals(ctx context.Context, variable string, component string, value float64) error {
	tuple, ok := ctx.Value(variables{variable}).(*Tuple)

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
	tuple, ok := ctx.Value(variables{variable}).(*Tuple)

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
	tuple, ok := ctx.Value(variables{variable}).(*Tuple)

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
	regex := fmt.Sprintf(`^(.+) ← tuple\(%s, %s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aTuple)

	regex = fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aPoint)

	regex = fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aVector)

	regex = fmt.Sprintf(`^(.+) ← color\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aColor)

	ctx.Step(`^(.+) ← normalize\((.+)\)$`, aNormalized)
}

func TupleAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+)\.(x|y|z|w|red|green|blue) = %s$`, shared.Decimal)
	ctx.Step(regex, aComponentEquals)

	ctx.Step(`^(.) is (not )?a point$`, aPointCheck)
	ctx.Step(`^(.) is (not )?a vector$`, aVectorCheck)

	regex = fmt.Sprintf(`^(.+) = tuple\(%s, %s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, compareTuple)

	regex = fmt.Sprintf(`^(.+) = vector\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, compareVector)

	regex = fmt.Sprintf(`^(.+) = point\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, comparePoint)

	regex = fmt.Sprintf(`^(.+) = color\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, compareColor)

	regex = fmt.Sprintf(`^magnitude\((.+)\) = (√)?%s$`, shared.Decimal)
	ctx.Step(regex, compareMag)

	ctx.Step(`^normalize\((.+)\) = (.+)$`, compareNormalize)

	regex = fmt.Sprintf(`^dot\((.+), (.+)\) = %s$`, shared.Decimal)
	ctx.Step(regex, compareDot)

	ctx.Step(`^cross\((.+), (.+)\) = (.+)$`, compareCross)
}

func TupleAssignments(ctx *godog.ScenarioContext) {
	ctx.Step(`^(.+) = (.+) \+ (.+)$`, add)
	ctx.Step(`^(.+) = (.+) \- (.+)$`, subtract)
	ctx.Step(`^(.+) = \-(.+)$`, negate)
	regex := fmt.Sprintf(`^(.+) = (.+) \* %s$`, shared.Decimal)
	ctx.Step(regex, scalar_mul)
	ctx.Step(`^(.+) = (.+) \* (.+)$`, mul)
	regex = fmt.Sprintf(`^(.+) = (.+) \/ %s$`, shared.Decimal)
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
