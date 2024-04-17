package primitive

import "fmt"

type Ray2 struct{ ray[Vector2] }

func NewRay2(src Vector2, dest Vector2) Ray2 {
	return Ray2{ray: ray[Vector2]{Src: src, Dest: dest}}
}

type Ray3 struct{ ray[Vector3] }

func NewRay3(src Vector3, dest Vector3) Ray3 {
	return Ray3{ray: ray[Vector3]{Src: src, Dest: dest}}
}

type Ray4 struct{ ray[Vector4] }

func NewRay4(src Vector4, dest Vector4) Ray4 {
	return Ray4{ray: ray[Vector4]{Src: src, Dest: dest}}
}

type Vector[V any] interface {
	fmt.Stringer
	Subtract(V) V
	Len() float64
	Len2() float64
	Scale(float64) V
	Add(V) V
	Dot(V) float64
	Project(V) V
	Norm() V
}

type ray[V Vector[V]] struct {
	Src  V
	Dest V
}

func (r *ray[V]) FromOrigin() V {
	return r.Dest.Subtract(r.Src)
}

func (r *ray[V]) Len() float64 {
	return r.FromOrigin().Len()
}

func (r *ray[V]) Len2() float64 {
	return r.FromOrigin().Len2()
}

func (r *ray[V]) At(f float64) V {
	return r.Src.Add(r.Dest.Scale(f))
}

func (r *ray[V]) String() string {
	return fmt.Sprintf("Ray{Src:%v,Dest:%v}", r.Src, r.Dest)
}
