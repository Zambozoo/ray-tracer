package scene

import (
	"encoding/json"
	"v0/primitive"
	"v0/scene/light"
	"v0/scene/material"
)

type Scene struct {
	Background primitive.Color

	Objects   []Object
	Lights    []light.Light
	Materials map[string]material.Material
}

func (s *Scene) UnmarshalJSON(b []byte) error {
	type sceneUnpacker struct {
		Ambient struct {
			Ka    float64
			Color primitive.ColorI
		}
		Background primitive.ColorI

		Objects   []any
		Lights    []any
		Materials map[string]any
	}

	su := &sceneUnpacker{}
	if err := json.Unmarshal(b, &su); err != nil {
		return err
	}

	s.Materials = map[string]material.Material{}
	for key, m := range su.Materials {
		bs, _ := json.Marshal(m)

		mat, err := material.JSONUnmarshalMaterial(bs)
		if err != nil {
			return err
		}

		s.Materials[key] = mat
	}

	s.Objects = make([]Object, 0, len(su.Objects))
	for _, o := range su.Objects {
		bs, _ := json.Marshal(o)

		obj, err := JSONUnmarshalObject(bs, s.Materials)
		if err != nil {
			return err
		}

		s.Objects = append(s.Objects, obj)
	}

	s.Lights = make([]light.Light, 0, len(su.Lights))
	for _, l := range su.Lights {
		bs, _ := json.Marshal(l)

		li, err := light.JSONUnmarshalLight(bs)
		if err != nil {
			return err
		}
		s.Lights = append(s.Lights, li)
	}

	s.Background = su.Background.ToColor()

	return nil
}
