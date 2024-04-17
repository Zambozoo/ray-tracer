package primitive

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X float64
	Y float64
}

func (v Vector2) String() string {
	return fmt.Sprintf("Vector2{X:%v,Y:%v}", v.X, v.Y)
}

func (v Vector2) Len() float64 {
	return math.Sqrt(v.Len2())
}

func (v Vector2) Len2() float64 {
	return float64(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) Norm() Vector2 {
	d := v.Len()
	return Vector2{
		X: float64(v.X) / d,
		Y: float64(v.Y) / d,
	}
}

func (v Vector2) Scale(f float64) Vector2 {
	return Vector2{
		X: f * float64(v.X),
		Y: f * float64(v.Y),
	}
}

func (a Vector2) Project(b Vector2) Vector2 {
	return b.Scale(float64(a.Dot(b)) / (b.Len() * b.Len()))
}

func (a Vector2) Dot(b Vector2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Vector2) Subtract(b Vector2) Vector2 {
	return Vector2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}
