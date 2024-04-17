package material

import (
	"encoding/json"
	"fmt"
	"v0/primitive"
	"v0/scene/hittable"
	"v0/scene/light"
)

type Metal struct {
	Color primitive.ColorI
	Fuzz  float64
}

func (m *Metal) Scatters(ray primitive.Ray3, hit hittable.Hit, lights []light.Light) ([]Scatter, bool) {
	reflected := ray.Dest.Norm().Reflect(hit.Normal)
	scattered := primitive.NewRay3(hit.Point, reflected.Jitter(m.Fuzz))
	if scattered.Dest.Dot(hit.Normal) > 0 {
		return []Scatter{
			{
				Ray:    scattered,
				Albedo: m.Color.ToColor(),
				Weight: 0.5,
			},
			{
				Ray:    primitive.NewRay3(hit.Point, reflected.Jitter(m.Fuzz)),
				Albedo: m.Color.ToColor(),
				Weight: 0.5,
			},
		}, true
	}

	return []Scatter{}, false
}

func (m *Metal) UnmarshalJSON(bs []byte) error {
	type metal Metal
	m2 := &metal{}
	if err := json.Unmarshal(bs, &m2); err != nil {
		return err
	}

	*m = Metal(*m2)
	if m.Fuzz < 0 || m.Fuzz > 1 {
		return fmt.Errorf("metal material Fuzz must be between [0, 1]: %v", m)
	}

	return nil
}
