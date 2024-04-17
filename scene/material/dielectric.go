package material

import (
	"math"
	"math/rand"
	"v0/primitive"
	"v0/scene/hittable"
	"v0/scene/light"
)

type Dielectric struct {
	RefractionIndex float64
	Color           primitive.ColorI
	SchlickScale    float64
	Fuzz            float64
}

func (d *Dielectric) Scatters(ray primitive.Ray3, hit hittable.Hit, lights []light.Light) ([]Scatter, bool) {
	var s []Scatter
	for i := 0; i < 2; i++ {
		refraction_ratio := d.RefractionIndex
		if hit.Outside {
			refraction_ratio = 1.0 / refraction_ratio
		}

		unit_direction := ray.Dest.Norm()
		cos_theta := min(unit_direction.Scale(-1).Dot(hit.Normal), 1.0)
		sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

		cannot_refract := refraction_ratio*sin_theta > 1.0
		var direction primitive.Vector3

		if cannot_refract || reflectance(cos_theta, refraction_ratio) > rand.Float64() {
			direction = d.reflect(unit_direction, hit.Normal)
		} else {
			direction = d.refract(unit_direction, hit.Normal, refraction_ratio)
		}

		s = append(s,
			Scatter{
				Ray:    primitive.NewRay3(hit.Point, direction),
				In:     !hit.Outside,
				Weight: 0.5,
				Albedo: d.Color.ToColor(),
			})
	}
	return s, true

}

func reflectance(cosine, ref_idx float64) float64 {
	// Use Schlick's approximation for reflectance.
	r0 := (1 - ref_idx) / (1 + ref_idx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

func (d *Dielectric) reflect(v, n primitive.Vector3) primitive.Vector3 {
	return v.Subtract(n.Scale(2 * v.Dot(n))).Jitter(d.Fuzz)
}

func (d *Dielectric) refract(uv, n primitive.Vector3, etai_over_etat float64) primitive.Vector3 {
	cos_theta := min(uv.Scale(-1).Dot(n), 1.0)
	r_out_perp := uv.Add(n.Scale(cos_theta)).Scale(etai_over_etat)
	r_out_parallel := n.Scale(-math.Sqrt(math.Abs(1.0 - r_out_perp.Len2())))
	return r_out_perp.Add(r_out_parallel).Jitter(d.Fuzz)
}
