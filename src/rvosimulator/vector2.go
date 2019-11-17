package rvosimulator

import (
	"math"
)

type Vector2 struct {
	X float64
	Y float64
}

func NewVector2(x float64, y float64) *Vector2 {
	v := &Vector2{
		X: x,
		Y: y,
	}
	return v
}

func Flip(vec *Vector2) *Vector2 {
	return &Vector2{X: -vec.X, Y: -vec.Y}
}

func Sub(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return &Vector2{X: vec1.X - vec2.X, Y: vec1.Y - vec2.Y}
}

func Add(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return &Vector2{X: vec1.X + vec2.X, Y: vec1.Y + vec2.Y}
}

func Mul(vec1 *Vector2, vec2 *Vector2) float64 {
	return vec1.X*vec2.X + vec1.Y*vec2.Y
}

func MulOne(vec *Vector2, s float64) *Vector2 {
	return &Vector2{X: vec.X * s, Y: vec.Y * s}
}

func Div(vec *Vector2, s float64) *Vector2 {
	return &Vector2{X: vec.X / s, Y: vec.Y / s}
}

func Equal(vec1 *Vector2, vec2 *Vector2) bool {
	return vec1.X == vec2.X && vec1.Y == vec2.Y
}

func NotEqual(vec1 *Vector2, vec2 *Vector2) bool {
	return vec1.X != vec2.X || vec1.Y != vec2.Y
}

func MulSum(vec *Vector2, s float64) *Vector2 {
	return Add(vec, MulOne(vec, s))
}

func DivSum(vec *Vector2, s float64) *Vector2 {
	return Add(vec, Div(vec, s))
}

func AddSum(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return Add(vec1, Add(vec1, vec2))
}

func SubSum(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return Add(vec1, Sub(vec1, vec2))
}

func Sqr(vec *Vector2) float64 {
	// vec * vec
	return Mul(vec, vec)
}

func Abs(vec *Vector2) float64 {
	// sqrt(vec * vec)
	return float64(math.Sqrt(float64(Mul(vec, vec))))
}

func Normalize(vec *Vector2) *Vector2 {
	// vector / abs(vector)
	return Div(vec, Abs(vec))
}

func Det(vec1 *Vector2, vec2 *Vector2) float64 {
	return vec1.X*vec2.Y - vec1.Y*vec2.X
}

func LeftOf(vec1 *Vector2, vec2 *Vector2, vec3 *Vector2) float64 {
	return Det(Sub(vec1, vec3), Sub(vec2, vec1))
}

func DistSqPointLineSegment(vec1 *Vector2, vec2 *Vector2, vec3 *Vector2) float64 {
	r := Mul(Sub(vec3, vec1), Sub(vec2, vec1)) / Sqr(Sub(vec2, vec1))

	if r < 0 {
		return Sqr(Sub(vec3, vec1))
	} else if r > 1 {
		return Sqr(Sub(vec3, vec2))
	} else {
		// absSq(c - (a + r * (b - a)))
		return Sqr(Sub(vec3, Add(vec1, MulOne(Sub(vec2, vec1), r))))
	}
}
