package hittable

import (
	"math"
	"v0/primitive"
)

type Sphere struct {
	Radius   float64
	Position primitive.Vector3
}

func (s *Sphere) Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool) {
	boundingVolume := &Prism{
		Position: primitive.Vector3{
			X: s.Position.X - s.Radius,
			Y: s.Position.Y - s.Radius,
			Z: s.Position.Z - s.Radius,
		},
		Dimensions: primitive.Vector3{
			X: s.Radius * 2,
			Y: s.Radius * 2,
			Z: s.Radius * 2,
		},
	}

	if _, hits := boundingVolume.Hits(ray, interval); !hits {
		return Hit{}, false
	}

	oc := ray.Src.Subtract(s.Position)
	a := ray.Dest.Len2()
	half_b := oc.Dot(ray.Dest)
	c := oc.Len2() - s.Radius*s.Radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return Hit{}, false
	}

	sqrtd := math.Sqrt(discriminant)

	root := (-half_b - sqrtd) / a
	if interval.Excludes(root) {
		root = (-half_b + sqrtd) / a
		if interval.Excludes(root) {
			return Hit{}, false
		}
	}

	point := ray.At(root)
	h := Hit{
		Distance: root,
		Point:    point,
	}

	unitNormOut := point.Subtract(s.Position).Scale(1 / s.Radius)
	h.SetNormal(ray, unitNormOut)

	return h, true
}

func RandomVectorInHemisphere(normal primitive.Vector3) primitive.Vector3 {
	randomUnitVector := primitive.RandUnitVector()
	if randomUnitVector.Dot(normal) > 0 {
		return randomUnitVector
	} else {
		return randomUnitVector.Scale(-1)
	}
}
