package main

import (
	"fmt"
	"math"
	"strings"
)

type canvas struct {
	pixels []tuple
	width  int32
	height int32
}

func new_canvas(width, height int32) *canvas {
	pixels := make([]tuple, width*height)
	return &canvas{
		pixels,
		width,
		height,
	}
}

func (c *canvas) write_pixel(x, y int32, color *tuple) {
	c.pixels[y*c.width+x] = *color
}

func (c *canvas) pixel_at(x, y int32) *tuple {
	return &c.pixels[y*c.width+x]
}

func component_to_255(c float64) byte {
	return byte(math.Max(0, math.Min(255, math.Round(c*255))))
}

func append_comp(builder, current_line *strings.Builder, comp string, force_newline bool) {
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

func (c *canvas) to_ppm() *string {
	var builder strings.Builder

	builder.WriteString("P3\n")
	builder.WriteString(fmt.Sprintf("%d %d\n", c.width, c.height))
	builder.WriteString("255\n")

	var current_line strings.Builder

	for i, p := range c.pixels {
		r := fmt.Sprintf("%d", component_to_255(p.red()))
		append_comp(&builder, &current_line, r, false)

		g := fmt.Sprintf("%d", component_to_255(p.green()))
		append_comp(&builder, &current_line, g, false)

		b := fmt.Sprintf("%d", component_to_255(p.blue()))
		x := i % int(c.width)
		append_comp(&builder, &current_line, b, x == int(c.width)-1)
	}
	builder.WriteString(current_line.String())
	builder.WriteString("\n")

	res := builder.String()
	return &res
}
