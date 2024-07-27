package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

var decimal = `(-?\d+(?:\.\d+)?)`
var posint = `(\d+)`

func aCanvas(ctx context.Context, variable string, width, height int32) (context.Context, error) {
	c := new_canvas(width, height)
	return context.WithValue(ctx, variables{name: variable}, c), nil
}

func aTuple(ctx context.Context, variable string, x, y, z, w float64) (context.Context, error) {
	t := tuple{
		x,
		y,
		z,
		w,
	}

	return context.WithValue(ctx, variables{name: variable}, &t), nil
}

func aPoint(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := point(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := vector(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aColor(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := color(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aNormalized(ctx context.Context, variable, in string) (context.Context, error) {
	t := ctx.Value(variables{name: in}).(*tuple)
	t = t.normalize()
	return context.WithValue(ctx, variables{name: variable}, t), nil
}

func add(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*tuple)
	rightTuple := ctx.Value(variables{name: right}).(*tuple)
	t := leftTuple.add(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func subtract(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*tuple)
	rightTuple := ctx.Value(variables{name: right}).(*tuple)
	t := leftTuple.subtract(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func mul(ctx context.Context, assignee, left, right string) (context.Context, error) {
	leftTuple := ctx.Value(variables{name: left}).(*tuple)
	rightTuple := ctx.Value(variables{name: right}).(*tuple)
	t := leftTuple.hadamard(rightTuple)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func negate(ctx context.Context, assignee, in string) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*tuple)
	t := inTuple.negate()
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func scalar_mul(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*tuple)
	t := inTuple.scalar_multiply(scalar)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func scalar_div(ctx context.Context, assignee, in string, scalar float64) (context.Context, error) {
	inTuple := ctx.Value(variables{name: in}).(*tuple)
	t := inTuple.scalar_div(scalar)
	return context.WithValue(ctx, variables{name: assignee}, t), nil
}

func compareTuple(ctx context.Context, variable string, x, y, z, w float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*tuple)

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
	actual, ok := ctx.Value(variables{name: variable}).(*tuple)

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
	actual, ok := ctx.Value(variables{name: variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := point(x, y, z)

	if *actual != *t {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareColor(ctx context.Context, variable string, x, y, z float64) error {
	actual, ok := ctx.Value(variables{name: variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple %s is not set", variable)
	}

	t := color(x, y, z)

	if !compare_tuple(actual, t) {
		return fmt.Errorf("%+v was not %+v", actual, t)
	}

	return nil
}

func compareMag(ctx context.Context, variable string, sqrt string, expected float64) error {
	actual, ok := ctx.Value(variables{variable}).(*tuple)

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
	l := ctx.Value(variables{name: left}).(*tuple)
	r := ctx.Value(variables{name: right}).(*tuple)

	if l.dot(r) != expected {
		return fmt.Errorf("dot(%+v, %+v) was %f, not %f", l, r, l.dot(r), expected)
	}

	return nil
}

func compareCross(ctx context.Context, left, right, expected string) error {
	l := ctx.Value(variables{name: left}).(*tuple)
	r := ctx.Value(variables{name: right}).(*tuple)
	e := ctx.Value(variables{name: expected}).(*tuple)

	if *l.cross(r) != *e {
		return fmt.Errorf("cross(%+v, %+v) was %+v, not %+v", l, r, l.cross(r), e)
	}

	return nil
}

func compareNormalize(ctx context.Context, variable string, expected string) error {
	in := ctx.Value(variables{variable}).(*tuple)
	expectedTuple := ctx.Value(variables{name: expected}).(*tuple)

	if !compare_tuple(in.normalize(), expectedTuple) {
		return fmt.Errorf("%+v normalized was %+v, not %+v", in, in.normalize(), expectedTuple)
	}

	return nil
}

func aComponentEquals(ctx context.Context, variable string, component string, value float64) error {
	tuple, ok := ctx.Value(variables{variable}).(*tuple)

	if !ok {
		return fmt.Errorf("tuple [%s] is not set (will check component [%s] for value [%f])", variable, component, value)
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
	case "red":
		actual = tuple.red()
	case "green":
		actual = tuple.green()
	case "blue":
		actual = tuple.blue()
	default:
		return fmt.Errorf("Unknown component '%s'", component)
	}

	if actual != value {
		return fmt.Errorf("Expected value %f for component '%s', actual was %f", value, component, actual)
	}

	return nil
}

func aCanvasComponentEquals(ctx context.Context, variable string, component string, value int32) error {
	canvas, ok := ctx.Value(variables{variable}).(*canvas)

	if !ok {
		return fmt.Errorf("tuple [%s] is not set (will check component [%s] for value [%d])", variable, component, value)
	}

	actual := 0
	switch component {
	case "width":
		actual = int(canvas.width)
	case "height":
		actual = int(canvas.height)
	default:
		return fmt.Errorf("Unknown component '%s'", component)
	}

	if actual != int(value) {
		return fmt.Errorf("Expected value %d for component '%s', actual was %d", value, component, actual)
	}

	return nil
}

func aPointCheck(ctx context.Context, variable string, notA string) error {
	tuple, ok := ctx.Value(variables{variable}).(*tuple)

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
	tuple, ok := ctx.Value(variables{variable}).(*tuple)

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

func everyPixelCheck(ctx context.Context, variable, color string) error {
	canvas := ctx.Value(variables{variable}).(*canvas)
	c := ctx.Value(variables{name: color}).(*tuple)

	for _, p := range canvas.pixels {
		if !compare_tuple(&p, c) {
			return errors.New("pixel check failed")
		}
	}
	return nil
}

func everyPixelSet(ctx context.Context, variable, color string) error {
	canvas := ctx.Value(variables{variable}).(*canvas)
	c := ctx.Value(variables{name: color}).(*tuple)

	for x := 0; x < int(canvas.width); x++ {
		for y := 0; y < int(canvas.height); y++ {
			canvas.write_pixel(int32(x), int32(y), c)
		}
	}

	return nil
}

func write_pixel(ctx context.Context, canvas_var string, x, y int32, color_var string) {
	canvas := ctx.Value(variables{name: canvas_var}).(*canvas)
	color := ctx.Value(variables{name: color_var}).(*tuple)
	canvas.write_pixel(x, y, color)
}

func pixel_at(ctx context.Context, canvas_var string, x, y int32, color_var string) error {
	canvas := ctx.Value(variables{name: canvas_var}).(*canvas)
	color := ctx.Value(variables{name: color_var}).(*tuple)
	actual := canvas.pixel_at(x, y)

	if !compare_tuple(color, actual) {
		return fmt.Errorf("pixel at %d, %d was %+v not %+v", x, y, actual, color)
	}
	return nil
}

func lines_are(ctx context.Context, from, to int32, variable string, value *godog.DocString) error {
	s := ctx.Value(variables{name: variable}).(*string)
	lines := strings.Split(*s, "\n")
	actual := strings.Join(lines[from-1:to], "\n")

	if strings.Compare(value.Content, actual) != 0 {
		return fmt.Errorf("Failed! Actual\n\n\"\"\"%s\"\"\"\n\nExpected\n\n\"\"\"%s\"\"\"", actual, value.Content)
	}

	return nil
}

func ends_with_newline(ctx context.Context, variable string) error {
	s := ctx.Value(variables{name: variable}).(*string)

	if !strings.HasSuffix(*s, "\n") {
		return fmt.Errorf("Failed! did not end with a newline")
	}

	return nil
}

func canvas_to_ppm(ctx context.Context, destination, canvas_var string) context.Context {
	canvas := ctx.Value(variables{name: canvas_var}).(*canvas)
	value := canvas.to_ppm()
	return context.WithValue(ctx, variables{name: destination}, value)
}

func TupleConstructors(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← tuple\(%s, %s, %s, %s\)$`, decimal, decimal, decimal, decimal)
	ctx.Step(regex, aTuple)

	regex = fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, aPoint)

	regex = fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, aVector)

	regex = fmt.Sprintf(`^(.+) ← color\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, aColor)

	ctx.Step(`^(.+) ← normalize\((.+)\)$`, aNormalized)
}

func TupleAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+)\.(x|y|z|w|red|green|blue) = %s$`, decimal)
	ctx.Step(regex, aComponentEquals)

	ctx.Step(`^(.) is (not )?a point$`, aPointCheck)
	ctx.Step(`^(.) is (not )?a vector$`, aVectorCheck)

	regex = fmt.Sprintf(`^(.+) = tuple\(%s, %s, %s, %s\)$`, decimal, decimal, decimal, decimal)
	ctx.Step(regex, compareTuple)

	regex = fmt.Sprintf(`^(.+) = vector\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, compareVector)

	regex = fmt.Sprintf(`^(.+) = point\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, comparePoint)

	regex = fmt.Sprintf(`^(.+) = color\(%s, %s, %s\)$`, decimal, decimal, decimal)
	ctx.Step(regex, compareColor)

	regex = fmt.Sprintf(`^magnitude\((.+)\) = (√)?%s$`, decimal)
	ctx.Step(regex, compareMag)

	ctx.Step(`^normalize\((.+)\) = (.+)$`, compareNormalize)

	regex = fmt.Sprintf(`^dot\((.+), (.+)\) = %s$`, decimal)
	ctx.Step(regex, compareDot)

	ctx.Step(`^cross\((.+), (.+)\) = (.+)$`, compareCross)
}

func TupleAssignments(ctx *godog.ScenarioContext) {
	ctx.Step(`^(.+) = (.+) \+ (.+)$`, add)
	ctx.Step(`^(.+) = (.+) \- (.+)$`, subtract)
	ctx.Step(`^(.+) = \-(.+)$`, negate)
	regex := fmt.Sprintf(`^(.+) = (.+) \* %s$`, decimal)
	ctx.Step(regex, scalar_mul)
	ctx.Step(`^(.+) = (.+) \* (.+)$`, mul)
	regex = fmt.Sprintf(`^(.+) = (.+) \/ %s$`, decimal)
	ctx.Step(regex, scalar_div)
}

func CanvasConstructors(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← canvas\(%s, %s\)$`, posint, posint)
	ctx.Step(regex, aCanvas)
}

func CanvasAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+)\.(width|height) = %s$`, posint)
	ctx.Step(regex, aCanvasComponentEquals)
	ctx.Step(`^every pixel of (.+) is (.+)$`, everyPixelCheck)
	regex = fmt.Sprintf(`^pixel_at\((.+), %s, %s\) = (.+)$`, posint, posint)
	ctx.Step(regex, pixel_at)
}

func CanvasAssignments(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^write_pixel\((.+), %s, %s, (.+)\)$`, posint, posint)
	ctx.Step(regex, write_pixel)
	ctx.Step(`^(.+) ← canvas_to_ppm\((.+)\)$`, canvas_to_ppm)
	ctx.Step(`^lines (\d+)-(\d+) of (.+) are$`, lines_are)
	ctx.Step(`^(.+) ends with a newline character$`, ends_with_newline)
	ctx.Step(`^set every pixel of (.+) to (.+)$`, everyPixelSet)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	TupleConstructors(ctx)
	TupleAssertions(ctx)
	TupleAssignments(ctx)

	CanvasConstructors(ctx)
	CanvasAssertions(ctx)
	CanvasAssignments(ctx)
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
