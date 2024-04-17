package hittable

import (
	"v0/primitive"
)

type Hittable interface {
	Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool)
}

type Hit struct {
	Point    primitive.Vector3
	Normal   primitive.Vector3
	Distance float64
	Outside  bool
}

func (h *Hit) SetNormal(ray primitive.Ray3, unitNormalOut primitive.Vector3) {
	h.Outside = ray.Dest.Dot(unitNormalOut) < 0
	h.Normal = unitNormalOut
	if !h.Outside {
		h.Normal = unitNormalOut.Scale(-1)
	}
}
