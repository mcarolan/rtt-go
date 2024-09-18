package canvas

import (
	"context"
	"errors"
	"fmt"
	"rtt/sharedtest"
	"rtt/tuple"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

func aCanvas(ctx context.Context, variable string, width, height int32) (context.Context, error) {
	c := NewCanvas(width, height)
	return context.WithValue(ctx, variables{name: variable}, c), nil
}

func aColor(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := tuple.Color(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func CanvasConstructors(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← canvas\(%s, %s\)$`, sharedtest.PosInt, sharedtest.PosInt)
	ctx.Step(regex, aCanvas)

	regex = fmt.Sprintf(`^(.+) ← color\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aColor)
}

func aCanvasComponentEquals(ctx context.Context, variable string, component string, value int32) error {
	canvas, ok := ctx.Value(variables{variable}).(*Canvas)

	if !ok {
		return fmt.Errorf("tuple [%s] is not set (will check component [%s] for value [%d])", variable, component, value)
	}

	actual := 0
	switch component {
	case "width":
		actual = int(canvas.Width)
	case "height":
		actual = int(canvas.Height)
	default:
		return fmt.Errorf("Unknown component '%s'", component)
	}

	if actual != int(value) {
		return fmt.Errorf("Expected value %d for component '%s', actual was %d", value, component, actual)
	}

	return nil
}

func everyPixelCheck(ctx context.Context, variable, color string) error {
	canvas := ctx.Value(variables{variable}).(*Canvas)
	c := ctx.Value(variables{name: color}).(*tuple.Tuple)

	for _, p := range canvas.Pixels {
		if !tuple.CompareTuple(&p, c) {
			return errors.New("pixel check failed")
		}
	}
	return nil
}

func everyPixelSet(ctx context.Context, variable, color string) error {
	canvas := ctx.Value(variables{variable}).(*Canvas)
	c := ctx.Value(variables{name: color}).(*tuple.Tuple)

	for x := 0; x < int(canvas.Width); x++ {
		for y := 0; y < int(canvas.Height); y++ {
			canvas.WritePixel(int32(x), int32(y), c)
		}
	}

	return nil
}

func writePixel(ctx context.Context, canvas_var string, x, y int32, color_var string) {
	canvas := ctx.Value(variables{name: canvas_var}).(*Canvas)
	color := ctx.Value(variables{name: color_var}).(*tuple.Tuple)
	canvas.WritePixel(x, y, color)
}

func pixelAt(ctx context.Context, canvas_var string, x, y int32, color_var string) error {
	canvas := ctx.Value(variables{name: canvas_var}).(*Canvas)
	color := ctx.Value(variables{name: color_var}).(*tuple.Tuple)
	actual := canvas.PixelAt(x, y)

	if !tuple.CompareTuple(color, actual) {
		return fmt.Errorf("pixel at %d, %d was %+v not %+v", x, y, actual, color)
	}
	return nil
}

func linesAre(ctx context.Context, from, to int32, variable string, value *godog.DocString) error {
	s := ctx.Value(variables{name: variable}).(*string)
	lines := strings.Split(*s, "\n")
	actual := strings.Join(lines[from-1:to], "\n")

	if strings.Compare(value.Content, actual) != 0 {
		return fmt.Errorf("Failed! Actual\n\n\"\"\"%s\"\"\"\n\nExpected\n\n\"\"\"%s\"\"\"", actual, value.Content)
	}

	return nil
}

func endsWithNewline(ctx context.Context, variable string) error {
	s := ctx.Value(variables{name: variable}).(*string)

	if !strings.HasSuffix(*s, "\n") {
		return fmt.Errorf("Failed! did not end with a newline")
	}

	return nil
}

func canvasToPPM(ctx context.Context, destination, canvas_var string) context.Context {
	canvas := ctx.Value(variables{name: canvas_var}).(*Canvas)
	value := canvas.ToPPM()
	return context.WithValue(ctx, variables{name: destination}, value)
}

func CanvasAssertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+)\.(width|height) = %s$`, sharedtest.PosInt)
	ctx.Step(regex, aCanvasComponentEquals)
	ctx.Step(`^every pixel of (.+) is (.+)$`, everyPixelCheck)
	regex = fmt.Sprintf(`^pixel_at\((.+), %s, %s\) = (.+)$`, sharedtest.PosInt, sharedtest.PosInt)
	ctx.Step(regex, pixelAt)
}

func CanvasAssignments(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^write_pixel\((.+), %s, %s, (.+)\)$`, sharedtest.PosInt, sharedtest.PosInt)
	ctx.Step(regex, writePixel)
	ctx.Step(`^(.+) ← canvas_to_ppm\((.+)\)$`, canvasToPPM)
	ctx.Step(`^lines (\d+)-(\d+) of (.+) are$`, linesAre)
	ctx.Step(`^(.+) ends with a newline character$`, endsWithNewline)
	ctx.Step(`^set every pixel of (.+) to (.+)$`, everyPixelSet)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	CanvasConstructors(ctx)
	CanvasAssertions(ctx)
	CanvasAssignments(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/canvas.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
