package matrix

import (
	"errors"
	"math"
	"rtt/shared"
	"rtt/tuple"
)

type Matrix struct {
	values []float64
	width  int
	height int
}

func FromValues(values []float64) *Matrix {
	size := int(math.Sqrt(float64(len(values))))
	return &Matrix{
		values: values,
		width:  size,
		height: size,
	}
}

func matrix(size int) *Matrix {
	values := make([]float64, size*size)
	return &Matrix{
		values: values,
		width:  size,
		height: size,
	}
}

func Matrix4() *Matrix {
	return matrix(4)
}

func Matrix3() *Matrix {
	return matrix(3)
}

func Matrix2() *Matrix {
	return matrix(2)
}

var Identity = &Matrix{
	values: []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	},
	width:  4,
	height: 4,
}

func (m Matrix) At(y, x int) float64 {
	return m.values[y*m.height+x]
}

func (m Matrix) set(y, x int, value float64) {
	m.values[y*m.height+x] = value
}

func (a Matrix) Multiply(b *Matrix) *Matrix {
	m := Matrix4()

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			value := a.At(row, 0)*b.At(0, col) + a.At(row, 1)*b.At(1, col) + a.At(row, 2)*b.At(2, col) + a.At(row, 3)*b.At(3, col)
			m.set(row, col, value)
		}
	}

	return m
}

func (a Matrix) MultiplyTuple(t *tuple.Tuple) *tuple.Tuple {
	result := make([]float64, 4)

	for y := 0; y < 4; y++ {
		row := tuple.Tuple{
			X: a.At(y, 0),
			Y: a.At(y, 1),
			Z: a.At(y, 2),
			W: a.At(y, 3),
		}
		result[y] = row.Dot(t)
	}

	return &tuple.Tuple{
		X: result[0],
		Y: result[1],
		Z: result[2],
		W: result[3],
	}
}

func (a Matrix) Transpose() *Matrix {
	m := matrix(a.width)

	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			m.set(x, y, a.At(y, x))
		}
	}

	return m
}

func (a Matrix) Determinant() float64 {
	result := 0.0
	if a.width == 2 {
		result = a.At(0, 0)*a.At(1, 1) - a.At(0, 1)*a.At(1, 0)
	} else {
		for x := 0; x < a.width; x++ {
			result = result + a.At(0, x)*a.Cofactor(0, x)
		}
	}
	return result
}

func (a Matrix) Equals(b *Matrix) bool {
	if a.width != b.width || a.height != b.height {
		return false
	}
	for x := 0; x < int(a.width); x++ {
		for y := 0; y < int(a.height); y++ {
			valueA := a.At(y, x)
			valueB := b.At(y, x)

			if !shared.CompareFloat(valueA, valueB) {
				return false
			}
		}
	}
	return true
}

func (a Matrix) Submatrix(yToDelete, xToDelete int) *Matrix {
	res := matrix(a.width - 1)

	for y := 0; y < a.height; y++ {
		if y == yToDelete {
			continue
		}

		for x := 0; x < a.width; x++ {
			if x == xToDelete {
				continue
			}

			xOffset := 0
			if x > xToDelete {
				xOffset = -1
			}
			yOffset := 0
			if y > yToDelete {
				yOffset = -1
			}

			res.set(y+yOffset, x+xOffset, a.At(y, x))
		}
	}
	return res
}

func (a Matrix) Minor(y, x int) float64 {
	return a.Submatrix(y, x).Determinant()
}

func (a Matrix) Cofactor(y, x int) float64 {
	minor := a.Minor(y, x)

	if (y+x)%2 == 0 {
		return minor
	} else {
		return -minor
	}
}

func (a Matrix) IsInvertible() bool {
	return a.Determinant() != 0
}

func (a Matrix) Invert() (*Matrix, error) {
	if !a.IsInvertible() {
		return nil, errors.New("input matrix not invertible")
	}

	result := matrix(a.width)

	det := a.Determinant()

	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			cofactor := a.Cofactor(y, x)
			result.set(x, y, cofactor/det)
		}
	}

	return result, nil
}
