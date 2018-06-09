package surface

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Sphere describes a 3d sphere
// TODO: make all of these private, this is accessed through interfaces anyway
type Sphere struct {
	Pos    *geom.Mat
	Mat    Material
	bounds *geom.Bounds
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with a given material and optional transforms.
func UnitSphere(m ...Material) *Sphere {
	s := &Sphere{
		Pos: geom.Identity(),
		Mat: &DefaultMaterial{},
	}
	if len(m) > 0 {
		s.Mat = m[0]
	}
	return s.transform(geom.Identity())
}

// TODO: unify with cube.transform AABB calc
func (s *Sphere) transform(t *geom.Mat) *Sphere {
	s.Pos = s.Pos.Mult(t)
	min := s.Pos.MultPoint(geom.Vec{})
	max := s.Pos.MultPoint(geom.Vec{})
	for x := -0.5; x <= 0.5; x += 1 {
		for y := -0.5; y <= 0.5; y += 1 {
			for z := -0.5; z <= 0.5; z += 1 {
				pt := s.Pos.MultPoint(geom.Vec{x, y, z})
				min = min.Min(pt)
				max = max.Max(pt)
			}
		}
	}
	s.bounds = geom.NewBounds(min, max)
	return s
}

func (s *Sphere) Move(x, y, z float64) *Sphere {
	return s.transform(geom.Trans(x, y, z))
}

func (s *Sphere) Scale(x, y, z float64) *Sphere {
	return s.transform(geom.Scale(x, y, z))
}

func (s *Sphere) Rotate(x, y, z float64) *Sphere {
	return s.transform(geom.Rot(geom.Vec{x, y, z}))
}

func (s *Sphere) Center() geom.Vec {
	return s.Pos.MultPoint(geom.Vec{})
}

func (s *Sphere) Bounds() *geom.Bounds {
	return s.bounds
}

// Intersect tests whether the sphere intersects a given ray.
// http://tfpsly.free.fr/english/index.html?url=http://tfpsly.free.fr/english/3d/Raytracing.html
// TODO: http://kylehalladay.com/blog/tutorial/math/2013/12/24/Ray-Sphere-Intersection.html
func (s *Sphere) Intersect(ray *geom.Ray) (obj render.Object, dist float64) {
	if ok, _, _ := s.bounds.Check(ray); !ok {
		return nil, 0
	}
	i := s.Pos.Inverse()
	r := i.MultRay(ray)
	op := geom.Vec{}.Minus(r.Origin)
	b := op.Dot(geom.Vec(r.Dir))
	det := b*b - op.Dot(op) + 0.5*0.5
	if det < 0 {
		return nil, 0
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t1)).Len()
		if dist > bias {
			return s, dist
		}
	}
	t2 := b + root
	if t2 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t2)).Len()
		if dist > bias {
			return s, dist
		}
	}
	return nil, 0
}

// At returns the surface normal given a point on the surface.
func (s *Sphere) At(pt geom.Vec, in geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	i := s.Pos.Inverse()
	p := i.MultPoint(pt)
	pu, _ := p.Unit()
	n := s.Pos.MultDir(pu)
	return n, s.Mat.At(0, 0, in.Dot(n), rnd)
}

func (s *Sphere) Light() rgb.Energy {
	return s.Mat.Light()
}

func (s *Sphere) Lights() []render.Object {
	if !s.Mat.Light().Zero() {
		return []render.Object{s}
	}
	return nil
}