package primitive

import (
	"fmt"
	"math"
	"math/rand"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vector3) String() string {
	return fmt.Sprintf("Vector3{X:%v,Y:%v,Z:%v}", v.X, v.Y, v.Z)
}

func (v Vector3) Len() float64 {
	return math.Sqrt(v.Len2())
}

func (v Vector3) Len2() float64 {
	return float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Norm() Vector3 {
	d := v.Len()
	return Vector3{
		X: float64(v.X) / d,
		Y: float64(v.Y) / d,
		Z: float64(v.Z) / d,
	}
}

func (v Vector3) Scale(f float64) Vector3 {
	return Vector3{
		X: f * float64(v.X),
		Y: f * float64(v.Y),
		Z: f * float64(v.Z),
	}
}

func (a Vector3) Project(b Vector3) Vector3 {
	return b.Scale(float64(a.Dot(b)) / (b.Len() * b.Len()))
}

func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vector3) Subtract(b Vector3) Vector3 {
	return Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vector3) Reflect(b Vector3) Vector3 {
	return a.Subtract(b.Scale(2 * a.Dot(b)))
}

func RandInUnitSphereVector3() Vector3 {
	for {
		v := Vector3{
			X: rand.Float64()*2 - 1,
			Y: rand.Float64()*2 - 1,
			Z: rand.Float64()*2 - 1,
		}
		if v.Len2() < 1 {
			return v
		}
	}
}

func RandUnitVector() Vector3 {
	return Vector3{
		X: rand.NormFloat64(),
		Y: rand.NormFloat64(),
		Z: rand.NormFloat64(),
	}.Norm()
}

func (a Vector3) Jitter(angle float64) Vector3 {
	d := math.Tan(angle)
	v := RandUnitVector().Scale(d)
	return a.Add(a.Norm()).Add(v)
}

func (a Vector3) Refract(b Vector3, etaRatio float64) Vector3 {
	cos := min(a.Scale(-1).Dot(b), 1)
	rayPerpendicular := a.Add(b.Scale(cos)).Scale(etaRatio)
	rayParallel := b.Scale(-math.Sqrt(math.Abs(1.0 - rayPerpendicular.Len2())))
	return rayPerpendicular.Add(rayParallel)
}

func (a Vector3) Dist(b Vector3) float64 {
	dx := (a.X - b.X)
	dy := (a.Y - b.Y)
	dz := (a.Z - b.Z)
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return d
}
