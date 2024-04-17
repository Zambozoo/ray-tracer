package light

import "v0/primitive"

type Directional struct {
	Direction primitive.Vector3
	Color     primitive.ColorI
}

func (d *Directional) ShadowVector(hitPoint primitive.Vector3) primitive.Vector3 {
	return d.Direction.Scale(-1)
}

func (d *Directional) ColorOf() primitive.Color {
	return d.Color.ToColor()
}
