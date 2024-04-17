package scene

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"v0/primitive"
	"v0/scene/hittable"
	"v0/scene/material"
)

type Object struct {
	UUID     uint64
	Material material.Material
	Hittable hittable.Hittable
}

func (o *Object) Equals(other *Object) bool {
	return o.UUID == other.UUID

}
func (o *Object) Hits(ray primitive.Ray3, interval primitive.IntervalF) (hittable.Hit, bool) {
	return o.Hittable.Hits(ray, interval)
}

func JSONUnmarshalObject(bs []byte, materials map[string]material.Material) (Object, error) {
	type unpacker struct {
		Type     string
		Material string
		Hittable map[string]any
	}

	u := &unpacker{}
	if err := json.Unmarshal(bs, &u); err != nil {
		return Object{}, err
	}

	material, ok := materials[u.Material]
	if !ok {
		return Object{}, fmt.Errorf("unknown material in object: %v", u)
	}

	o := Object{Material: material}
	switch u.Type {
	case "sphere":
		o.Hittable = &hittable.Sphere{}
	case "plane":
		o.Hittable = &hittable.Plane{}
	case "triangle":
		o.Hittable = &hittable.Triangle{}
	case "prism":
		o.Hittable = &hittable.Prism{}
	case "circle":
		o.Hittable = &hittable.Circle{}
	default:
		return o, fmt.Errorf("unknown object 'Type': %s", u.Type)
	}

	bytes, err := json.Marshal(u.Hittable)
	if err != nil {
		return Object{}, err
	}
	if err := json.Unmarshal(bytes, &o.Hittable); err != nil {
		return Object{}, err
	}

	o.UUID = rand.Uint64()

	return o, nil
}
