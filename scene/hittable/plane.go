package hittable

import "v0/primitive"

type Plane struct {
	Point  primitive.Vector3
	Normal primitive.Vector3
}

// Hits implements Hittable.
func (p *Plane) Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool) {
	// assuming vectors are all normalized
	denom := p.Normal.Dot(ray.Dest)
	if denom > 1e-6 {
		diff := p.Point.Subtract(ray.Src)
		t := diff.Dot(p.Normal) / denom
		return Hit{
			Point:    ray.Dest.Scale(t).Add(ray.Src),
			Normal:   p.Normal.Scale(-1),
			Distance: t,
			Outside:  true,
		}, t >= 0
	}

	return Hit{}, false
}
