package light

import "v0/primitive"

type Ambient struct {
	Color primitive.ColorI
}

func (a *Ambient) ShadowVector(hitPoint primitive.Vector3) primitive.Vector3 {
	return hitPoint
}

func (a *Ambient) ColorOf() primitive.Color {
	return a.Color.ToColor()
}
