package main

import (
	"fmt"
	"os"
	"rtt/canvas"
	"rtt/ray"
	"rtt/tuple"
)

// clock
// func main() {
// 	c := canvas.NewCanvas(800, 600)
// 	midX := c.Width / 2
// 	midY := c.Height / 2
// 	radius := float64(midY) / 2.0

// 	origin := tuple.Point(0, 0, 0)
// 	c.WritePixel(int32(origin.X)+midX, int32(origin.Z)+midY, tuple.Red)

// 	twelve := tuple.Point(0, 0, 1)

// 	for hour := 0.0; hour < 12; hour++ {
// 		r := transformations.RotationY(hour * math.Pi / 6.0)
// 		p := r.MultiplyTuple(twelve)
// 		c.WritePixel(int32(radius*p.X)+midX, int32(radius*p.Z)+midY, tuple.White)
// 	}

// 	ppm := c.ToPPM()

// 	if err := os.WriteFile("clock.ppm", []byte(*ppm), 0666); err != nil {
// 		fmt.Printf("Error writing result: %s", err)
// 		return
// 	}
// }

func main() {
	shape := ray.NewSphere()
	rayOrigin := tuple.Point(0, 0, -5)
	wallZ := 10.0
	wallSize := 7.0
	canvasPixels := 100.0
	pixelSize := wallSize / canvasPixels
	half := wallSize / 2.0

	c := canvas.NewCanvas(int32(canvasPixels), int32(canvasPixels))
	for y := 0; y < int(c.Height); y++ {
		worldY := half - pixelSize*float64(y)
		for x := 0; x < int(c.Width); x++ {
			worldX := -half + pixelSize*float64(x)
			position := tuple.Point(worldX, worldY, wallZ)

			r := ray.NewRay(*rayOrigin, *position.Subtract(rayOrigin).Normalize())
			intersections, err := shape.Intersect(r)

			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			if ray.Hit(intersections) != nil {
				c.WritePixel(int32(x), int32(y), tuple.Red)
			}
		}
	}
	ppm := c.ToPPM()

	if err := os.WriteFile("sphere.ppm", []byte(*ppm), 0666); err != nil {
		fmt.Printf("Error writing result: %s", err)
		os.Exit(1)
	}
}
