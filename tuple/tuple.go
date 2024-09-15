package tuple

import (
	"math"
	"rtt/shared"
)

type Tuple struct {
	X float64
	Y float64
	Z float64
	W float64
}

func (t *Tuple) IsPoint() bool {
	return t.W == 1
}

func (t *Tuple) IsVector() bool {
	return t.W == 0
}

func Point(x, y, z float64) *Tuple {
	return &Tuple{
		X: x,
		Y: y,
		Z: z,
		W: 1,
	}
}

func Vector(x, y, z float64) *Tuple {
	return &Tuple{
		X: x,
		Y: y,
		Z: z,
		W: 0,
	}
}

func Color(r, g, b float64) *Tuple {
	return &Tuple{
		X: r,
		Y: g,
		Z: b,
		W: 0,
	}
}

func (t *Tuple) Red() float64 {
	return t.X
}

func (t *Tuple) Green() float64 {
	return t.Y
}

func (t *Tuple) Blue() float64 {
	return t.Z
}

func (t *Tuple) Add(other *Tuple) *Tuple {
	return &Tuple{
		X: t.X + other.X,
		Y: t.Y + other.Y,
		Z: t.Z + other.Z,
		W: t.W + other.W,
	}
}

func (t *Tuple) Subtract(other *Tuple) *Tuple {
	return &Tuple{
		X: t.X - other.X,
		Y: t.Y - other.Y,
		Z: t.Z - other.Z,
		W: t.W - other.W,
	}
}

func (t *Tuple) Negate() *Tuple {
	return &Tuple{
		X: -t.X,
		Y: -t.Y,
		Z: -t.Z,
		W: -t.W,
	}
}

func (t *Tuple) ScalarMultiply(v float64) *Tuple {
	return &Tuple{
		X: t.X * v,
		Y: t.Y * v,
		Z: t.Z * v,
		W: t.W * v,
	}
}

func (t *Tuple) ScalarDiv(v float64) *Tuple {
	return &Tuple{
		X: t.X / v,
		Y: t.Y / v,
		Z: t.Z / v,
		W: t.W / v,
	}
}

func (t *Tuple) Magnitude() float64 {
	return math.Sqrt(math.Pow(t.X, 2) + math.Pow(t.Y, 2) + math.Pow(t.Z, 2) + math.Pow(t.W, 2))
}

func (t *Tuple) Normalize() *Tuple {
	mag := t.Magnitude()
	return &Tuple{
		X: t.X / mag,
		Y: t.Y / mag,
		Z: t.Z / mag,
		W: t.W / mag,
	}
}

func (t *Tuple) Dot(other *Tuple) float64 {
	return t.X*other.X +
		t.Y*other.Y +
		t.Z*other.Z +
		t.W*other.W
}

func (a *Tuple) Cross(b *Tuple) *Tuple {
	return Vector(a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X)
}

func (a *Tuple) Hadamard(b *Tuple) *Tuple {
	red := a.Red() * b.Red()
	green := a.Green() * b.Green()
	blue := a.Blue() * b.Blue()
	return Color(red, green, blue)
}

func CompareTuple(a, b *Tuple) bool {
	return shared.CompareFloat(a.X, b.X) && shared.CompareFloat(a.Y, b.Y) && shared.CompareFloat(a.Z, b.Z) && shared.CompareFloat(a.W, b.W)
}

var White = Color(1, 1, 1)
var Red = Color(1, 0, 0)
