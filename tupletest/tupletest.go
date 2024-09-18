package tupletest

import (
	"context"
	"fmt"
	"rtt/sharedtest"
	"rtt/tuple"

	"github.com/cucumber/godog"
)

func doConstructPoint(ctx context.Context, variable string, xString, yString, zString string) (context.Context, error) {
	x, y, z, err := sharedtest.ParseXYZ(xString, yString, zString)
	if err != nil {
		return ctx, err
	}
	p := tuple.Point(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, p), nil
}

func AddConstructPoint(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doConstructPoint)
}

func doConstructVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := tuple.Vector(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, p), nil
}

func AddConstructVector(sc *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	sc.Step(regex, doConstructVector)
}
