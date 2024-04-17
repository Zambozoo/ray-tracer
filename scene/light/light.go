package light

import (
	"encoding/json"
	"fmt"
	"v0/primitive"
)

type Light interface {
	ShadowVector(hitPoint primitive.Vector3) primitive.Vector3
	ColorOf() primitive.Color
}

func JSONUnmarshalLight(bs []byte) (Light, error) {
	type unpacker struct {
		Type string
	}

	u := &unpacker{}
	if err := json.Unmarshal(bs, &u); err != nil {
		return nil, err
	}

	var l Light
	switch u.Type {
	case "directional":
		l = &Directional{}
	case "ambient":
		l = &Ambient{}
	case "sphere":
		l = &Sphere{}
	default:
		return nil, fmt.Errorf("unknown light 'Type': %s", u.Type)
	}

	if err := json.Unmarshal(bs, &l); err != nil {
		return nil, err
	}

	return l, nil
}
