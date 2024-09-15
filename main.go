package main

import (
	"fmt"
	"math"
	"os"
	"rtt/canvas"
	"rtt/transformations"
	"rtt/tuple"
)

func main() {
	c := canvas.NewCanvas(800, 600)
	midX := c.Width / 2
	midY := c.Height / 2
	radius := float64(midY) / 2.0

	origin := tuple.Point(0, 0, 0)
	c.WritePixel(int32(origin.X)+midX, int32(origin.Z)+midY, tuple.Red)

	twelve := tuple.Point(0, 0, 1)

	for hour := 0.0; hour < 12; hour++ {
		r := transformations.RotationY(hour * math.Pi / 6.0)
		p := r.MultiplyTuple(twelve)
		c.WritePixel(int32(radius*p.X)+midX, int32(radius*p.Z)+midY, tuple.White)
	}

	ppm := c.ToPPM()

	if err := os.WriteFile("clock.ppm", []byte(*ppm), 0666); err != nil {
		fmt.Printf("Error writing result: %s", err)
		return
	}
}
