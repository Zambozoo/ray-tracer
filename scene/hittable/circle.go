package hittable

import "v0/primitive"

type Circle struct {
	Position primitive.Vector3
	Normal   primitive.Vector3
	Radius   float64
}

func (c *Circle) boundingVolume() *Prism {
	return &Prism{
		Position: primitive.Vector3{
			X: c.Position.X - c.Radius,
			Y: c.Position.Y - c.Radius,
			Z: c.Position.Z - c.Radius,
		},
		Dimensions: primitive.Vector3{
			X: c.Radius * 2,
			Y: c.Radius * 2,
			Z: c.Radius * 2,
		},
	}
}

func (c *Circle) Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool) {
	if _, hits := c.boundingVolume().Hits(ray, interval); !hits {
		return Hit{}, false
	}

	if hit, ok := (&Plane{Point: c.Position, Normal: c.Normal}).Hits(ray, interval); ok && hit.Point.Dist(c.Position) < c.Radius {
		return hit, true
	}

	return Hit{}, false
}
