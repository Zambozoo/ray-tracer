package scene

import (
	"encoding/json"
	"fmt"
	"v0/primitive"
)

type Camera struct {
	LookRay        primitive.Ray3
	Up             primitive.Vector3
	ViewPlaneWidth float64
}

func (c *Camera) FocalLength() float64 {
	return c.LookRay.Len()
}

func (c *Camera) NormX() primitive.Vector3 {
	return c.Up.Cross(c.LookRay.FromOrigin()).Norm()
}

func (c *Camera) NormY() primitive.Vector3 {
	return c.Up.Norm()
}

func (c *Camera) NormZ() primitive.Vector3 {
	return c.LookRay.FromOrigin().Norm()
}

func (c *Camera) ViewPlaneRect(aspectRatio float64) primitive.Rect {
	viewPlane := primitive.Vector2{
		X: c.ViewPlaneWidth,
		Y: c.ViewPlaneWidth / aspectRatio,
	}

	center := c.NormZ().Scale(c.FocalLength())
	xOffset := c.NormX().Scale(viewPlane.X / 2)
	yOffset := c.NormY().Scale(-viewPlane.Y / 2)
	return primitive.Rect{
		BottomRight: center.Add(yOffset).Add(xOffset),
		TopLeft:     center.Subtract(yOffset).Subtract(xOffset),
	}
}

func (c *Camera) UnmarshalJSON(b []byte) error {
	type cameraUnpacker struct {
		CameraLookFrom primitive.Vector3
		CameraLookAt   primitive.Vector3
		CameraUp       primitive.Vector3
		ViewPlaneWidth float64

		AspectRatio float64
	}
	cu := &cameraUnpacker{}
	if err := json.Unmarshal(b, cu); err != nil {
		return err
	}

	if cu.CameraLookFrom == cu.CameraLookAt {
		return fmt.Errorf("zero 'Camera' ray: %v", cu)
	}

	v1 := cu.CameraUp.Subtract(cu.CameraLookFrom)
	v2 := cu.CameraLookAt.Subtract(cu.CameraLookFrom)
	cos := v1.Dot(v2) / (v1.Len() * v2.Len())
	if cos < 0 {
		return fmt.Errorf("viewpoint is out of view of camera: %v", cu)
	}

	c.LookRay = primitive.NewRay3(cu.CameraLookFrom, cu.CameraLookAt)
	c.Up = cu.CameraUp
	c.ViewPlaneWidth = cu.ViewPlaneWidth

	return nil
}
