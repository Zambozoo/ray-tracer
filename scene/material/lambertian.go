package material

import (
	"v0/primitive"
	"v0/scene/hittable"
	"v0/scene/light"
)

type Lambertian struct {
	Color primitive.ColorI
}

func (l *Lambertian) Scatters(ray primitive.Ray3, hit hittable.Hit, lights []light.Light) ([]Scatter, bool) {
	scatterRay := primitive.NewRay3(hit.Point, hit.Normal.Jitter(30))
	if scatterRay.Dest.Len2() < 3e-8 {
		scatterRay.Dest = hit.Normal
	}

	return []Scatter{
		{
			Ray:    scatterRay,
			Albedo: l.Color.ToColor(),
			Weight: 0.5,
		},
		{
			Ray:    scatterRay,
			Albedo: l.Color.ToColor(),
			Weight: 0.5,
		},
	}, true
}
