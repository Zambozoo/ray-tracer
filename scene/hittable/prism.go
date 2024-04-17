package hittable

import "v0/primitive"

type Prism struct {
	Position   primitive.Vector3
	Dimensions primitive.Vector3
}

// Hits implements Hittable.
func (p *Prism) Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool) {
	t1 := primitive.Vector3{
		X: (p.Position.X - ray.Src.X) / ray.Dest.X,
		Y: (p.Position.Y - ray.Src.Y) / ray.Dest.Y,
		Z: (p.Position.Z - ray.Src.Z) / ray.Dest.Z,
	}

	t2 := primitive.Vector3{
		X: (p.Position.X + p.Dimensions.X - ray.Src.X) / ray.Dest.X,
		Y: (p.Position.Y + p.Dimensions.Y - ray.Src.Y) / ray.Dest.Y,
		Z: (p.Position.Z + p.Dimensions.Z - ray.Src.Z) / ray.Dest.Z,
	}

	tNear := max(min(t1.X, t2.X), max(min(t1.Y, t2.Y), min(t1.Z, t2.Z)))
	tFar := min(max(t1.X, t2.X), min(max(t1.Y, t2.Y), max(t1.Z, t2.Z)))

	if tNear > tFar || tFar < 0 {
		return Hit{}, false
	}

	point := ray.Dest.Scale(tNear).Add(ray.Src)
	return Hit{
		Point:    point,
		Normal:   p.normal(point),
		Distance: tNear,
		Outside:  false,
	}, true
}

func (p *Prism) normal(point primitive.Vector3) primitive.Vector3 {
	center := p.Dimensions.Scale(0.5).Add(p.Position)
	n := primitive.Vector3{
		X: point.X - center.X,
		Y: point.Y - center.Y,
		Z: point.Z - center.Z,
	}
	return n.Norm()
}
