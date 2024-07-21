package main

import "math"

type tuple struct {
	x float64
	y float64
	z float64
	w float64
}

func (t *tuple) is_point() bool {
	return t.w == 1
}

func (t *tuple) is_vector() bool {
	return t.w == 0
}

func point(x, y, z float64) *tuple {
	return &tuple{
		x: x,
		y: y,
		z: z,
		w: 1,
	}
}

func vector(x, y, z float64) *tuple {
	return &tuple{
		x: x,
		y: y,
		z: z,
		w: 0,
	}
}

func (t *tuple) add(other *tuple) *tuple {
	return &tuple{
		x: t.x + other.x,
		y: t.y + other.y,
		z: t.z + other.z,
		w: t.w + other.w,
	}
}

func (t *tuple) subtract(other *tuple) *tuple {
	return &tuple{
		x: t.x - other.x,
		y: t.y - other.y,
		z: t.z - other.z,
		w: t.w - other.w,
	}
}

func (t *tuple) negate() *tuple {
	return &tuple{
		x: -t.x,
		y: -t.y,
		z: -t.z,
		w: -t.w,
	}
}

func (t *tuple) scalar_multiply(v float64) *tuple {
	return &tuple{
		x: t.x * v,
		y: t.y * v,
		z: t.z * v,
		w: t.w * v,
	}
}

func (t *tuple) scalar_div(v float64) *tuple {
	return &tuple{
		x: t.x / v,
		y: t.y / v,
		z: t.z / v,
		w: t.w / v,
	}
}

func (t *tuple) magnitude() float64 {
	return math.Sqrt(math.Pow(t.x, 2) + math.Pow(t.y, 2) + math.Pow(t.z, 2) + math.Pow(t.w, 2))
}

func (t *tuple) normalize() *tuple {
	mag := t.magnitude()
	return &tuple{
		x: t.x / mag,
		y: t.y / mag,
		z: t.z / mag,
		w: t.w / mag,
	}
}

func (t *tuple) dot(other *tuple) float64 {
	return t.x*other.x +
		t.y*other.y +
		t.z*other.z +
		t.w*other.w
}

func (a *tuple) cross(b *tuple) *tuple {
	return vector(a.y*b.z-a.z*b.y,
		a.z*b.x-a.x*b.z,
		a.x*b.y-a.y*b.x)
}
