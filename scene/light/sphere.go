package light

import "v0/primitive"

type Sphere struct {
	Position primitive.Vector3
	Radius   float64
	Color    primitive.ColorI
}

func (s *Sphere) ShadowVector(hitPoint primitive.Vector3) primitive.Vector3 {
	return hitPoint.Subtract(s.Position)
}

func (s *Sphere) ColorOf() primitive.Color {
	return s.Color.ToColor()
}
