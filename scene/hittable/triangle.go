package hittable

import (
	"encoding/json"
	"fmt"
	"math"
	"v0/primitive"
)

type Triangle struct {
	Points []primitive.Vector3
	Plane  Plane
}

func (t *Triangle) boundingVolume() *Prism {
	minPosition := primitive.Vector3{
		X: math.Inf(1),
		Y: math.Inf(1),
		Z: math.Inf(1),
	}
	maxPosition := primitive.Vector3{
		X: math.Inf(-1),
		Y: math.Inf(-1),
		Z: math.Inf(-1),
	}

	for _, p := range t.Points {
		if p.X > maxPosition.X {
			maxPosition.X = p.X
		} else if p.X < minPosition.X {
			minPosition.X = p.X
		}

		if p.Y > maxPosition.Y {
			maxPosition.Y = p.Y
		} else if p.X < minPosition.Y {
			minPosition.Y = p.Y
		}

		if p.Z > maxPosition.Z {
			maxPosition.Z = p.Z
		} else if p.Z < minPosition.Z {
			minPosition.Z = p.Z
		}
	}

	return &Prism{
		Position: minPosition,
		Dimensions: primitive.Vector3{
			X: maxPosition.X - minPosition.X,
			Y: maxPosition.Y - minPosition.Y,
			Z: maxPosition.Z - minPosition.Z,
		},
	}
}

// Hits implements Hittable.
func (tri *Triangle) Hits(ray primitive.Ray3, interval primitive.IntervalF) (Hit, bool) {
	if _, hits := tri.boundingVolume().Hits(ray, interval); !hits {
		return Hit{}, false
	}

	N := tri.Plane.Normal.Norm()

	NdotRayDirection := N.Dot(ray.Dest)
	if math.Abs(NdotRayDirection) < interval.Min {
		// Parallel
		return Hit{}, false
	}

	d := -N.Dot(tri.Points[0])
	t := -(N.Dot(ray.Src) + d) / NdotRayDirection
	if t < 0 {
		return Hit{}, false
	}

	P := ray.Src.Add(ray.Dest.Scale(t))
	edge0 := tri.Points[1].Subtract(tri.Points[0])
	vp0 := P.Subtract(tri.Points[0])
	C := edge0.Cross(vp0)
	if N.Dot(C) < 0 {
		return Hit{}, false
	}

	edge1 := tri.Points[2].Subtract(tri.Points[1])
	vp1 := P.Subtract(tri.Points[1])
	C = edge1.Cross(vp1)
	if N.Dot(C) < 0 {
		return Hit{}, false
	}

	edge2 := tri.Points[0].Subtract(tri.Points[2])
	vp2 := P.Subtract(tri.Points[2])
	C = edge2.Cross(vp2)
	if N.Dot(C) < 0 {
		return Hit{}, false
	}

	return Hit{
		Point:    P,
		Normal:   tri.Plane.Normal.Scale(-1).Norm(),
		Distance: t,
		Outside:  true,
	}, true
}

func (p *Triangle) UnmarshalJSON(b []byte) error {
	type unmarshaller struct {
		Points []primitive.Vector3
	}

	u := &unmarshaller{}
	if err := json.Unmarshal(b, &u); err != nil {
		return err
	}
	if len(u.Points) != 3 {
		return fmt.Errorf("triangles requires 3 points")
	}
	p.Points = u.Points

	// compute the plane's normal
	v0v1 := p.Points[1].Subtract(p.Points[0])
	v0v2 := p.Points[2].Subtract(p.Points[0])
	// no need to normalize
	N := v0v1.Cross(v0v2) // N

	// REQUIRES POINTS TO BE COUNTER CLOCKWISE
	p.Plane = Plane{
		Point:  p.Points[0],
		Normal: N,
	}

	return nil
}
