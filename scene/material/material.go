package material

import (
	"encoding/json"
	"fmt"
	"v0/primitive"
	"v0/scene/hittable"
	"v0/scene/light"
)

type Material interface {
	Scatters(ray primitive.Ray3, hit hittable.Hit, lights []light.Light) ([]Scatter, bool)
}

type Scatter struct {
	Ray    primitive.Ray3
	Albedo primitive.Color
	Weight float64
	In     bool
}

func JSONUnmarshalMaterial(bs []byte) (Material, error) {
	type unpacker struct {
		Type string
	}

	var u *unpacker
	if err := json.Unmarshal(bs, &u); err != nil {
		return nil, err
	}

	var m Material
	switch u.Type {
	case "dielectric":
		m = &Dielectric{}
	case "lambertian":
		m = &Lambertian{}
	case "metal":
		m = &Metal{Fuzz: 1}
	default:
		return nil, fmt.Errorf("unknown material 'Type': %s", u.Type)
	}

	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}

	switch mat := m.(type) {
	case *Metal:
		if mat.Fuzz < 0 || mat.Fuzz > 1 {
			return nil, fmt.Errorf("metal material Fuzz must be between [0, 1]: %v", m)
		}
	default:
	}

	return m, nil
}
