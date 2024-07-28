package canvas

import (
	"fmt"
	"math"
	"rtt/tuple"
	"strings"
)

type Canvas struct {
	Pixels []tuple.Tuple
	Width  int32
	Height int32
}

func NewCanvas(width, height int32) *Canvas {
	pixels := make([]tuple.Tuple, width*height)
	return &Canvas{
		pixels,
		width,
		height,
	}
}

func (c *Canvas) WritePixel(x, y int32, color *tuple.Tuple) {
	c.Pixels[y*c.Width+x] = *color
}

func (c *Canvas) PixelAt(x, y int32) *tuple.Tuple {
	return &c.Pixels[y*c.Width+x]
}

func componentTo255(c float64) byte {
	return byte(math.Max(0, math.Min(255, math.Round(c*255))))
}

func appendComponent(builder, current_line *strings.Builder, comp string, force_newline bool) {
	if force_newline || current_line.Len()+len(comp)+1 == 70 {
		builder.WriteString(current_line.String())
		builder.WriteString(" ")
		builder.WriteString(comp)
		builder.WriteString("\n")
		current_line.Reset()
	} else if current_line.Len()+len(comp)+1 > 70 {
		builder.WriteString(current_line.String())
		builder.WriteString("\n")
		builder.WriteString(comp)
		builder.WriteString(" ")
		current_line.Reset()
	} else if current_line.Len() > 0 {
		current_line.WriteString(" ")
		current_line.WriteString(comp)
	} else {
		current_line.WriteString(comp)
	}
}

func (c *Canvas) ToPPM() *string {
	var builder strings.Builder

	builder.WriteString("P3\n")
	builder.WriteString(fmt.Sprintf("%d %d\n", c.Width, c.Height))
	builder.WriteString("255\n")

	var current_line strings.Builder

	for i, p := range c.Pixels {
		r := fmt.Sprintf("%d", componentTo255(p.Red()))
		appendComponent(&builder, &current_line, r, false)

		g := fmt.Sprintf("%d", componentTo255(p.Green()))
		appendComponent(&builder, &current_line, g, false)

		b := fmt.Sprintf("%d", componentTo255(p.Blue()))
		x := i % int(c.Width)
		appendComponent(&builder, &current_line, b, x == int(c.Width)-1)
	}
	builder.WriteString(current_line.String())
	builder.WriteString("\n")

	res := builder.String()
	return &res
}
