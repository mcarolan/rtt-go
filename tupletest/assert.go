package tupletest

import (
	"context"
	"fmt"
	"rtt/sharedtest"
	"rtt/tuple"

	"github.com/cucumber/godog"
)

func doCompareTuple(ctx context.Context, variable string, x, y, z, w float64) error {
	actual, ok := ctx.Value(sharedtest.Variables{Name: variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := tuple.Tuple{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}

	if !tuple.CompareTuple(actual, &t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func AddCompareTuple(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) = tuple\(%s, %s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doCompareTuple)
}

func doCompareVector(ctx context.Context, variable string, xString, yString, zString string) error {
	actual, ok := ctx.Value(sharedtest.Variables{Name: variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	x, y, z, err := sharedtest.ParseXYZ(xString, yString, zString)
	if err != nil {
		return err
	}

	t := tuple.Vector(x, y, z)

	if !tuple.CompareTuple(actual, t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func AddCompareVector(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) = vector\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doCompareVector)
}

func doComparePoint(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(sharedtest.Variables{Name: variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := tuple.Point(x, y, z)

	if !tuple.CompareTuple(actual, t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func AddComparePoint(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) = point\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doComparePoint)
}

func doCompareColor(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(sharedtest.Variables{Name: variable}).(*tuple.Tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := tuple.Color(x, y, z)

	if !tuple.CompareTuple(actual, t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func AddCompareColor(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) = color\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doCompareColor)
}

func doCompareNormalize(ctx context.Context, variable string, expected string) error {
	in := ctx.Value(sharedtest.Variables{Name: variable}).(*tuple.Tuple)
	expectedTuple := ctx.Value(sharedtest.Variables{Name: expected}).(*tuple.Tuple)

	if !tuple.CompareTuple(in.Normalize(), expectedTuple) {
		return fmt.Errorf("%+v normalized was %+v, not %+v", in, in.Normalize(), expectedTuple)
	}

	return nil
}

func AddCompareNormalize(sc *godog.ScenarioContext) {
	sc.Step(`^normalize\((.+)\) = (.+)$`, doCompareNormalize)
	sc.Step(`^(.+) = normalize\((.+)\)$`, doCompareNormalize)
}
