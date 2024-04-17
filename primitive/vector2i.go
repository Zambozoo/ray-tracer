package primitive

import "fmt"

type Vector2I struct {
	X int
	Y int
}

func (v *Vector2I) String() string {
	return fmt.Sprintf("Vector2{X:%v,Y:%v}", v.X, v.Y)
}

func (v *Vector2I) ToVector2() Vector2 {
	return Vector2{
		X: float64(v.X),
		Y: float64(v.Y),
	}
}
