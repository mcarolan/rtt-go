package ray

import "rtt/tuple"

type PointLight struct {
	Position  tuple.Tuple
	Intensity tuple.Tuple
}

type Material struct {
	Color     tuple.Tuple
	Diffuse   float64
	Specular  float64
	Shininess float64
}

func NewMaterial() *Material {
	return &Material{
		Color:     *tuple.White,
		Diffuse:   0.9,
		Specular:  0.9,
		Shininess: 200,
	}
}

func NewPointLight(position, intensity tuple.Tuple) *PointLight {
	return &PointLight{
		Intensity: intensity,
		Position:  position,
	}
}
