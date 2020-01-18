package rvosimulator

import (
	"math"
	"fmt"
	"os"
)

// Vector2 :
type Vector2 struct {
	X float64
	Y float64
}

// NewVector2 :
func NewVector2(x float64, y float64) *Vector2 {
	v := &Vector2{
		X: x,
		Y: y,
	}
	return v
}

// Flip :
func Flip(vec *Vector2) *Vector2 {
	return &Vector2{X: -vec.X, Y: -vec.Y}
}


// Sub: 
func Sub(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return &Vector2{X: vec1.X - vec2.X, Y: vec1.Y - vec2.Y}
}

// Add :
func Add(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return &Vector2{X: vec1.X + vec2.X, Y: vec1.Y + vec2.Y}
}

// Mul :
func Mul(vec1 *Vector2, vec2 *Vector2) float64 {
	if math.IsNaN(vec1.X*vec2.X + vec1.Y*vec2.Y){
		fmt.Printf("NaN Mul\n")
		os.Exit(0)
	}
	return vec1.X*vec2.X + vec1.Y*vec2.Y
}

// MulOne :
func MulOne(vec *Vector2, s float64) *Vector2 {
	return &Vector2{X: vec.X * s, Y: vec.Y * s}
}

// Div :
func Div(vec *Vector2, s float64) *Vector2 {
	if math.IsNaN(vec.X / s) || math.IsNaN(vec.Y / s) {
		fmt.Printf("NaN Div\n")
		os.Exit(0)
	}
	return &Vector2{X: vec.X / s, Y: vec.Y / s}
}

// Equal :
func Equal(vec1 *Vector2, vec2 *Vector2) bool {
	return vec1.X == vec2.X && vec1.Y == vec2.Y
}

// NotEqual :
func NotEqual(vec1 *Vector2, vec2 *Vector2) bool {
	return vec1.X != vec2.X || vec1.Y != vec2.Y
}

// MulSum :
func MulSum(vec *Vector2, s float64) *Vector2 {
	return Add(vec, MulOne(vec, s))
}

// DivSum :
func DivSum(vec *Vector2, s float64) *Vector2 {

	fmt.Printf("DivSum\n")
	return Add(vec, Div(vec, s))
}

// AddSum :
func AddSum(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return Add(vec1, Add(vec1, vec2))
}

// SubSum :
func SubSum(vec1 *Vector2, vec2 *Vector2) *Vector2 {
	return Add(vec1, Sub(vec1, vec2))
}

// Sqr :
func Sqr(vec *Vector2) float64 {
	if math.IsNaN(Mul(vec, vec)) {
		fmt.Printf("NaN Sqr\n")
		os.Exit(0)
	}
	return Mul(vec, vec)
}

// Abs :
func Abs(vec *Vector2) float64 {
	if math.IsNaN(math.Sqrt(Mul(vec, vec))) {
		fmt.Printf("NaN Abs\n")
		os.Exit(0)
	}
	return math.Sqrt(Mul(vec, vec))
}

// Normalize :
func Normalize(vec *Vector2) *Vector2 {
	fmt.Printf("Normarize\n")
	return Div(vec, Abs(vec))
}

// Det :
func Det(vec1 *Vector2, vec2 *Vector2) float64 {
	if math.IsNaN(vec1.X*vec2.Y - vec1.Y*vec2.X) {
		fmt.Printf("NaN Det\n")
		os.Exit(0)
	}
	return vec1.X*vec2.Y - vec1.Y*vec2.X
}

// LeftOf :
func LeftOf(vec1 *Vector2, vec2 *Vector2, vec3 *Vector2) float64 {
	if math.IsNaN( Det(Sub(vec1, vec3), Sub(vec2, vec1))) {
		fmt.Printf("NaN LeftOf\n")
		os.Exit(0)
	}
	return Det(Sub(vec1, vec3), Sub(vec2, vec1))
}

// DistSqPointLineSegment :
func DistSqPointLineSegment(vec1 *Vector2, vec2 *Vector2, vec3 *Vector2) float64 {
	r := Mul(Sub(vec3, vec1), Sub(vec2, vec1)) / Sqr(Sub(vec2, vec1))
	if math.IsNaN(Sqr(Sub(vec3, vec1))) || math.IsNaN(Sqr(Sub(vec3, vec2))) || math.IsNaN(Sqr(Sub(vec3, Add(vec1, MulOne(Sub(vec2, vec1), r))))) {
		fmt.Printf("NaN Dist\n")
		os.Exit(0)
	}

	if r < 0 {
		return Sqr(Sub(vec3, vec1))
	} else if r > 1 {
		return Sqr(Sub(vec3, vec2))
	} else {
		return Sqr(Sub(vec3, Add(vec1, MulOne(Sub(vec2, vec1), r))))
	}
}
