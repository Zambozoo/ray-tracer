package primitive

import (
	"fmt"
	"math"
)

type Vector4 struct {
	X float64
	Y float64
	Z float64
	W float64
}

func (v Vector4) String() string {
	return fmt.Sprintf("Vector4{X:%v,Y:%v,Z:%v,W:%v}", v.X, v.Y, v.Z, v.W)
}

func (v Vector4) Len() float64 {
	return math.Sqrt(v.Len2())
}

func (v Vector4) Len2() float64 {
	return float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

func (v Vector4) Norm() Vector4 {
	d := v.Len()
	return Vector4{
		X: float64(v.X) / d,
		Y: float64(v.Y) / d,
		Z: float64(v.Z) / d,
		W: float64(v.W) / d,
	}
}

func (v Vector4) Scale(f float64) Vector4 {
	return Vector4{
		X: f * float64(v.X),
		Y: f * float64(v.Y),
		Z: f * float64(v.Z),
		W: f * float64(v.W),
	}
}

func (a Vector4) Project(b Vector4) Vector4 {
	return b.Scale(float64(a.Dot(b)) / (b.Len() * b.Len()))
}

func (a Vector4) Dot(b Vector4) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func (a Vector4) Add(b Vector4) Vector4 {
	return Vector4{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
		W: a.W + b.W,
	}
}

func (a Vector4) Subtract(b Vector4) Vector4 {
	return Vector4{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
		W: a.W - b.W,
	}
}
