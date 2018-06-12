package surface

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Triangle describes a triangle
type Triangle struct {
	Points  [3]geom.Vec
	Normals [3]geom.Dir
	Texture [3]geom.Vec
	Mat     Material
	Pos     *geom.Mat // TODO: implement transformations
	edge1   geom.Vec
	edge2   geom.Vec
	bounds  *geom.Bounds
}

// NewTriangle creates a new triangle
func NewTriangle(a, b, c geom.Vec, m ...Material) *Triangle {
	edge1 := b.Minus(a)
	edge2 := c.Minus(a)
	n, _ := edge1.Cross(edge2).Unit()
	t := &Triangle{
		Points:  [3]geom.Vec{a, b, c},
		Normals: [3]geom.Dir{n, n, n},
		Mat:     &DefaultMaterial{},
		edge1:   edge1,
		edge2:   edge2,
	}
	if len(m) > 0 {
		t.Mat = m[0]
	}
	min := t.Points[0].Min(t.Points[1]).Min(t.Points[2])
	max := t.Points[0].Max(t.Points[1]).Max(t.Points[2])
	t.bounds = geom.NewBounds(min, max)
	return t
}

func (t *Triangle) Bounds() *geom.Bounds {
	return t.bounds
}

// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (t *Triangle) Intersect(ray *geom.Ray) (obj render.Object, dist float64) {
	v0 := t.Points[0]
	v1 := t.Points[1]
	v2 := t.Points[2]
	edge1 := v1.Minus(v0)
	edge2 := v2.Minus(v0)
	h := geom.Vec(ray.Dir).Cross(edge2)
	a := edge1.Dot(h)
	if a > -bias && a < bias {
		return nil, 0
	}
	f := 1 / a
	s := ray.Origin.Minus(v0)
	u := f * s.Dot(h)
	if u < 0 || u > 1 {
		return nil, 0
	}
	q := s.Cross(edge1)
	v := f * geom.Vec(ray.Dir).Dot(q)
	if v < 0 || u+v > 1 {
		return nil, 0
	}
	dist = f * edge2.Dot(q)
	if dist <= bias {
		return nil, 0
	}
	return t, dist
}

// Intersect determines whether or not, and where, a Ray intersects this Triangle
// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
// TODO: fix.
func (t *Triangle) Intersect2(ray *geom.Ray) (obj render.Object, dist float64) {
	if ok, _, _ := t.bounds.Check(ray); !ok {
		return nil, 0
	}
	h, _ := ray.Dir.Cross(geom.Dir(t.edge2))
	a := t.edge1.Dot(geom.Vec(h))
	if a > -bias && a < bias {
		return nil, 0
	}
	f := 1.0 / a
	s := ray.Origin.Minus(t.Points[0])
	u := f * s.Dot(geom.Vec(h))
	if u < 0 || u > 1 {
		return nil, 0
	}
	q := s.Cross(t.edge1)
	v := f * geom.Vec(ray.Dir).Dot(q)
	if v < 0 || (u+v) > 1 {
		return nil, 0
	}
	dist = f * t.edge2.Dot(q)
	if dist < bias {
		return nil, 0
	}
	return t, dist
}

func (t *Triangle) Center() geom.Vec {
	c := geom.Vec{}
	for _, p := range t.Points {
		c = c.Plus(p)
	}
	return c.Scaled(1.0 / 3.0)
}

// At returns the material at a point on the Triangle
func (t *Triangle) At(pt geom.Vec, in geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	u, v, w := t.Bary(pt)
	n := t.normal(u, v, w)
	texture := t.texture(u, v, w)
	m := t.Mat.At(texture.X, texture.Y, in.Dot(n), rnd)
	return n, m
}

func (t *Triangle) Lights() []render.Object {
	if !t.Mat.Light().Zero() {
		return []render.Object{t}
	}
	return nil
}

func (t *Triangle) Light() rgb.Energy {
	return t.Mat.Light()
}

func (t *Triangle) Transmit() rgb.Energy {
	return t.Mat.Transmit()
}

// SetNormals sets values for each vertex normal
func (t *Triangle) SetNormals(a, b, c *geom.Dir) {
	if a != nil {
		t.Normals[0] = *a
	}
	if b != nil {
		t.Normals[1] = *b
	}
	if c != nil {
		t.Normals[2] = *c
	}
}

func (t *Triangle) SetTexture(a, b, c geom.Vec) {
	t.Texture[0] = a
	t.Texture[1] = b
	t.Texture[2] = c
}

// Normal computes the smoothed normal
func (t *Triangle) normal(u, v, w float64) geom.Dir { // TODO: instead of separate u, v, w just use a Vec and multiply
	n0 := t.Normals[0].Scaled(u)
	n1 := t.Normals[1].Scaled(v)
	n2 := t.Normals[2].Scaled(w)
	n, _ := n0.Plus(n1).Plus(n2).Unit()
	return n
}

func (t *Triangle) texture(u, v, w float64) geom.Vec {
	tex0 := t.Texture[0].Scaled(u)
	tex1 := t.Texture[1].Scaled(v)
	tex2 := t.Texture[2].Scaled(w)
	return tex0.Plus(tex1).Plus(tex2)
}

// Bary returns the Barycentric coords of Vector p on Triangle t
// https://codeplea.com/triangular-interpolation
func (t *Triangle) Bary(p geom.Vec) (u, v, w float64) {
	v0 := t.Points[1].Minus(t.Points[0])
	v1 := t.Points[2].Minus(t.Points[0])
	v2 := p.Minus(t.Points[0])
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	d := d00*d11 - d01*d01
	v = (d11*d20 - d01*d21) / d
	w = (d00*d21 - d01*d20) / d
	u = 1 - v - w
	return
}
