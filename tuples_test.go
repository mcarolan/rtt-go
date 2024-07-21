package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/cucumber/godog"
)

type tuples struct{ variable string }

var decimal = `(-?\d+(?:\.\d+)?)`

func aTuple(ctx context.Context, variable string, x, y, z, w float64) (context.Context, error) {
	t := tuple{
		x,
		y,
		z,
		w,
	}

	return context.WithValue(ctx, tuples{variable}, t), nil
}

func aPoint(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := point(x, y, z)
	return context.WithValue(ctx, tuples{variable}, p), nil
}

func aVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := vector(x, y, z)
	return context.WithValue(ctx, tuples{variable}, p), nil
}

func aNormalized(ctx context.Context, variable, in string) (context.Context, error) {
	t := ctx.Value(tuples{variable: in}).(*tuple)
	t = t.normalize()
	return context.WithValue(ctx, tuples{variable}, t), nil
}

func add(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(tuples{variable: left}).(tuple)
	rightTuple := ctx.Value(tuples{variable: right}).(tuple)
	t := leftTuple.add(&rightTuple)
	return context.WithValue(ctx, tuples{variable: assignee}, t), nil
}

func subtract(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(tuples{variable: left}).(*tuple)
	rightTuple := ctx.Value(tuples{variable: right}).(*tuple)
	t := leftTuple.subtract(rightTuple)
	return context.WithValue(ctx, tuples{variable: assignee}, t), nil
}

func negate(ctx context.Context, assignee, in string) (context.Context, error) {
	inTuple := ctx.Value(tuples{variable: in}).(tuple)
	t := inTuple.negate()
	return context.WithValue(ctx, tuples{variable: assignee}, t), nil
}

func scalar_mul(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(tuples{variable: in}).(tuple)
	t := inTuple.scalar_multiply(scalar)
	return context.WithValue(ctx, tuples{variable: assignee}, t), nil
}

func scalar_div(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(tuples{variable: in}).(tuple)
	t := inTuple.scalar_div(scalar)
	return context.WithValue(ctx, tuples{variable: assignee}, t), nil
}

func compareTuple(ctx context.Context, variable string, x, y, z, w float64) error {
	actual, ok := ctx.Value(tuples{variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := tuple{
		x: x,
		y: y,
		z: z,
		w: w,
	}

	if *actual != t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareVector(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(tuples{variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := vector(x, y, z)

	if *actual != *t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func comparePoint(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(tuples{variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := point(x, y, z)

	if *actual != *t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareMag(ctx context.Context, variable string, sqrt string, expected float64) error {
	actual, ok := ctx.Value(tuples{variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if sqrt == "√" {
		expected = math.Sqrt(expected)
	}

	if actual.magnitude() != expected {
		return fmt.Errorf("%+v magnitude was %f, not %f", actual, actual.magnitude(), expected)
	}

	return nil
}

func compareDot(ctx context.Context, left, right string, expected float64) error {
	l := ctx.Value(tuples{variable: left}).(*tuple)
	r := ctx.Value(tuples{variable: right}).(*tuple)

	if l.dot(r) != expected {
		return fmt.Errorf("dot(%+v, %+v) was %f, not %f", l, r, l.dot(r), expected)
	}

	return nil
}

func compareCross(ctx context.Context, left, right, expected string) error {
	l := ctx.Value(tuples{variable: left}).(*tuple)
	r := ctx.Value(tuples{variable: right}).(*tuple)
	e := ctx.Value(tuples{variable: expected}).(*tuple)

	if *l.cross(r) != *e {
		return fmt.Errorf("cross(%+v, %+v) was %+v, not %+v", l, r, l.cross(r), e)
	}

	return nil
}

func compareNormalize(ctx context.Context, variable string, expected string) error {
	in := ctx.Value(tuples{variable}).(*tuple)
	expectedTuple := ctx.Value(tuples{variable: expected}).(*tuple)

	if !compare_tuple(in.normalize(), expectedTuple) {
		return fmt.Errorf("%+v normalized was %+v, not %+v", in, in.normalize(), expectedTuple)
	}

	return nil
}

func aComponentEquals(ctx context.Context, variable string, component string, value float64) error {
	tuple, ok := ctx.Value(tuples{variable}).(tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	actual := 0.0
	switch component {
	case "x":
		actual = tuple.x
	case "y":
		actual = tuple.y
	case "z":
		actual = tuple.z
	case "w":
		actual = tuple.w
	default:
		return fmt.Errorf("Unknown component '%s'", component)
	}

	if actual != value {
		return fmt.Errorf("Expected value %f for component '%s', actual was %f", value, component, actual)
	}

	return nil
}

func aPointCheck(ctx context.Context, variable string, notA string) error {
	tuple, ok := ctx.Value(tuples{variable}).(tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if !tuple.is_point() && notA == "" {
		return errors.New("aTuple is not a point")
	}

	if tuple.is_point() && notA == "not a" {
		return errors.New("aTuple is a point")
	}

	return nil
}

func aVectorCheck(ctx context.Context, variable string, notA string) error {
	tuple, ok := ctx.Value(tuples{variable}).(tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	if !tuple.is_vector() && notA == "" {
		return errors.New("aTuple is not a vector")
	}

	if tuple.is_vector() && notA == "not a " {
		return errors.New("aTuple is a vector")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← tuple\(%s, %s, %s, %s\)$`, decimal, decimal, decimal, decimal)
	ctx.Step(regex, aTuple)

	regex = fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, aPoint)

	regex = fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, aVector)

	ctx.Step(`^(.+) ← normalize\((.+)\)$`, aNormalized)

	regex = fmt.Sprintf(`^(.)\.(x|y|z|w) = %s$`, decimal)
	ctx.Step(regex, aComponentEquals)

	ctx.Step(`^(.) is (not )?a point$`, aPointCheck)
	ctx.Step(`^(.) is (not )?a vector$`, aVectorCheck)

	regex = fmt.Sprintf(`^(.+) = tuple\(%s, %s, %s, %s\)$`, decimal, decimal, decimal, decimal)
	ctx.Step(regex, compareTuple)

	regex = fmt.Sprintf(`^(.+) = vector\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, compareVector)

	regex = fmt.Sprintf(`^(.+) = point\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, comparePoint)

	regex = fmt.Sprintf(`^magnitude\((.+)\) = (√)?%s$`, decimal)
	ctx.Step(regex, compareMag)

	ctx.Step(`^normalize\((.+)\) = (.+)$`, compareNormalize)

	regex = fmt.Sprintf(`^dot\((.+), (.+)\) = %s$`, decimal)
	ctx.Step(regex, compareDot)

	ctx.Step(`^cross\((.+), (.+)\) = (.+)$`, compareCross)

	ctx.Step(`^(.+) = (.+) \+ (.+)$`, add)
	ctx.Step(`^(.+) = (.+) \- (.+)$`, subtract)
	ctx.Step(`^(.+) = \-(.+)$`, negate)
	regex = fmt.Sprintf(`^(.+) = (.+) \* %s$`, decimal)
	ctx.Step(regex, scalar_mul)
	regex = fmt.Sprintf(`^(.+) = (.+) \/ %s$`, decimal)
	ctx.Step(regex, scalar_div)

	ctx.Step(regex, aComponentEquals)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
