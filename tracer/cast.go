package tracer

import (
	"math"
	"sync"
	"v0/io"
	"v0/primitive"
	"v0/scene"
	"v0/scene/hittable"
	"v0/scene/light"
	"v0/scene/material"
)

var (
	rayInterval = primitive.NewIntervalF(0.0001, math.Inf(1))
)

func cast(input *io.Input, pixelQueue chan coloredPixel, pixelWaitGroup *sync.WaitGroup, rayWaitGroup *WaitGroup) func(weightedRay) {
	return func(wr weightedRay) {
		pixelWaitGroup.Add(1)
		pixelQueue <- castSerial(wr, input)
		rayWaitGroup.Done()
	}
}

func castSerial(wr weightedRay, input *io.Input) coloredPixel {
	// Base case, use background
	if wr.depth >= input.MaxDepth {
		unit := wr.ray.Dest.Norm()
		a := 0.5 * (unit.X + 1)
		return coloredPixel{
			pixel: wr.pixel,
			color: primitive.White.Scale((1 - a)).Add(input.Scene.Background.Scale(a)).Scale(wr.weight).Hadamard(wr.attenuation),
		}

	}

	closestHit, closestObject := collide(wr, input)
	if closestHit.Distance < rayInterval.Max {
		lightColor := lightColor(closestHit, closestObject, input)
		cp := coloredPixel{pixel: wr.pixel}
		if scatters, ok := closestObject.Material.Scatters(wr.ray, closestHit, input.Scene.Lights); ok {
			for _, scatter := range scatters {
				newRay := weightedRay{
					weight:      wr.weight * scatter.Weight / float64(len(scatters)),
					ray:         scatter.Ray,
					depth:       wr.depth + 1,
					pixel:       wr.pixel,
					attenuation: wr.attenuation.Hadamard(scatter.Albedo).Hadamard(lightColor),
				}
				newColoredPixel := castSerial(newRay, input)
				cp.color = cp.color.Add(newColoredPixel.color)
			}
		}
		return cp
	}

	// More Background
	unit := wr.ray.Dest.Norm()
	a := 0.5 * (unit.X + 1)
	return coloredPixel{
		pixel: wr.pixel,
		color: primitive.White.Scale((1 - a)).Add(input.Scene.Background.Scale(a)).Scale(wr.weight).Hadamard(wr.attenuation),
	}
}

func collide(wr weightedRay, input *io.Input) (hittable.Hit, scene.Object) {
	// Check collision
	closestHit := hittable.Hit{Distance: rayInterval.Max}
	var closestObject scene.Object
	for _, o := range input.Scene.Objects {
		if hit, ok := o.Hits(wr.ray, primitive.NewIntervalF(rayInterval.Min, closestHit.Distance)); ok && hit.Distance < closestHit.Distance {
			closestHit = hit
			closestObject = o
		}
	}

	// Adjust with bias
	if _, ok := closestObject.Material.(*material.Dielectric); !ok {
		if closestHit.Outside {
			closestHit.Point = closestHit.Point.Subtract(wr.ray.Dest.Scale(rayInterval.Min))
		} else {
			closestHit.Point = closestHit.Point.Add(closestHit.Normal.Scale(rayInterval.Min))
		}
	}

	return closestHit, closestObject
}

func lightColor(closestHit hittable.Hit, closestObject scene.Object, input *io.Input) primitive.Color {
	lightColor := primitive.Black
outer:
	for _, l := range input.Scene.Lights {
		if _, ok := l.(*light.Ambient); ok {
			lightColor = lightColor.Add(l.ColorOf())
			continue
		}
		sVector := l.ShadowVector(closestHit.Point)
		shadowRay := primitive.NewRay3(closestHit.Point, sVector.Jitter(30))

		for _, o := range input.Scene.Objects {
			if _, ok := o.Material.(*material.Dielectric); ok {
				continue
			} else if _, ok := o.Hits(shadowRay, rayInterval); ok {
				break outer
			}
		}

		shadowVector := shadowRay.Dest.Subtract(shadowRay.Src).Norm()
		n_dot_l := closestHit.Normal.Norm().Dot(shadowVector)
		if _, ok := closestObject.Material.(*material.Dielectric); ok {
			n_dot_l = math.Abs(n_dot_l)
		}

		if n_dot_l > 0 {
			c := l.ColorOf().Scale(n_dot_l * 255)
			lightColor = lightColor.Add(c)
		}

	}

	return lightColor.Clamp()
}
