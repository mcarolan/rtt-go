package ray

import (
	"math"
	"rtt/matrix"
	"rtt/tuple"
)

type Ray struct {
	Origin    tuple.Tuple
	Direction tuple.Tuple
}

type Sphere struct {
	Id                int
	transformation    matrix.Matrix
	transformationInv matrix.Matrix
}

type Intersection struct {
	T      float64
	Object int
}

var objectCounter = 0

func NewSphere() *Sphere {
	objectCounter += 1
	return &Sphere{
		Id:                objectCounter,
		transformation:    *matrix.Identity,
		transformationInv: *matrix.Identity,
	}
}

func (s *Sphere) SetTransform(transform *matrix.Matrix) error {
	inverse, err := transform.Invert()

	if err != nil {
		return err
	}

	s.transformation = *transform
	s.transformationInv = *inverse
	return nil
}

func (s *Sphere) NormalAt(worldPoint tuple.Tuple) *tuple.Tuple {
	objectPoint := s.transformationInv.MultiplyTuple(&worldPoint)
	objectNormal := objectPoint.Subtract(tuple.ZeroPoint)
	worldNormal := s.transformationInv.Transpose().MultiplyTuple(objectNormal)
	worldNormal.W = 0
	return worldNormal.Normalize()
}

func (s *Sphere) Intersect(ray *Ray) []Intersection {
	ray2 := ray.Transform(&s.transformationInv)

	sphereToRay := ray2.Origin.Subtract(tuple.ZeroPoint)

	a := ray2.Direction.Dot(&ray2.Direction)
	b := 2 * ray2.Direction.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1

	discriminant := math.Pow(b, 2) - 4*a*c

	if discriminant < 0 {
		return []Intersection{}
	} else {
		t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
		t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
		return []Intersection{*s.Intersection(t1), *s.Intersection(t2)}
	}
}

func Hit(intersections []Intersection) *Intersection {
	var hit *Intersection

	for _, intersection := range intersections {
		if intersection.T < 0 {
			continue
		}

		if hit == nil || intersection.T < hit.T {
			hit = &intersection
		}
	}

	return hit
}

func (s *Sphere) Intersection(t float64) *Intersection {
	return &Intersection{
		T:      t,
		Object: s.Id,
	}
}

func NewRay(origin tuple.Tuple, direction tuple.Tuple) *Ray {
	return &Ray{
		Origin:    origin,
		Direction: direction,
	}
}

func (r *Ray) Position(t float64) *tuple.Tuple {
	return r.Origin.Add(r.Direction.ScalarMultiply(t))
}

func (r *Ray) Transform(m *matrix.Matrix) *Ray {
	return &Ray{
		Origin:    *m.MultiplyTuple(&r.Origin),
		Direction: *m.MultiplyTuple(&r.Direction),
	}
}
