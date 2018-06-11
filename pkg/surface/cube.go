package surface

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Cube describes the orientation and material of a unit cube
type Cube struct {
	Pos    *geom.Mat
	Mat    Material
	bounds *geom.Bounds
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with material and optional transforms.
func UnitCube(m ...Material) *Cube {
	c := &Cube{
		Pos: geom.Identity(),
		Mat: &DefaultMaterial{},
	}
	if len(m) > 0 {
		c.Mat = m[0]
	}
	return c.transform(geom.Identity())
}

func (c *Cube) Intersect(ray *geom.Ray) (obj render.Object, dist float64) {
	if ok, _, _ := c.bounds.Check(ray); !ok {
		return nil, 0
	}
	inv := c.Pos.Inverse() // global to local transform
	r := inv.MultRay(ray)  // translate ray into local space
	tmin := 0.0
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		t0 := (-0.5 - r.OrArray[a]) * r.InvArray[a]
		t1 := (0.5 - r.OrArray[a]) * r.InvArray[a]
		if r.InvArray[a] < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}
		if tmax < tmin {
			return nil, 0
		}
	}
	if tmin > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmin)).Len(); dist >= bias {
			return c, dist
		}
	}
	if tmax > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmax)).Len(); dist >= bias {
			return c, dist
		}
	}
	return nil, 0
}

// At returns the normal geom.Vec at this point on the Surface
func (c *Cube) At(pt geom.Vec, in geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	normal = geom.Dir{}
	i := c.Pos.Inverse()  // global to local transform
	p1 := i.MultPoint(pt) // translate point into local space
	abs := p1.Abs()
	u, v := 0.0, 0.0
	switch {
	case abs.X > abs.Y && abs.X > abs.Z:
		normal = geom.Dir{math.Copysign(1, p1.X), 0, 0}
		u = p1.Z + 0.5
		v = p1.Y + 0.5
	case abs.Y > abs.Z:
		normal = geom.Dir{0, math.Copysign(1, p1.Y), 0}
		u = p1.Z + 0.5
		v = p1.X + 0.5
	default:
		normal = geom.Dir{0, 0, math.Copysign(1, p1.Z)}
		u = p1.X + 0.5
		v = p1.Y + 0.5
	}
	n := c.Pos.MultDir(normal)
	return n, c.Mat.At(u, v, in.Dot(n), rnd)
}

func (c *Cube) Bounds() *geom.Bounds {
	return c.bounds
}

func (c *Cube) Lights() []render.Object {
	if !c.Mat.Light().Zero() {
		return []render.Object{c}
	}
	return nil
}

func (c *Cube) Light() rgb.Energy {
	return c.Mat.Light()
}

func (c *Cube) Transmit() rgb.Energy {
	return c.Mat.Transmit()
}

func (c *Cube) transform(m *geom.Mat) *Cube {
	c.Pos = c.Pos.Mult(m)
	min := c.Pos.MultPoint(geom.Vec{})
	max := c.Pos.MultPoint(geom.Vec{})
	for x := -0.5; x <= 0.5; x += 1 {
		for y := -0.5; y <= 0.5; y += 1 {
			for z := -0.5; z <= 0.5; z += 1 {
				pt := c.Pos.MultPoint(geom.Vec{x, y, z})
				min = min.Min(pt)
				max = max.Max(pt)
			}
		}
	}
	c.bounds = geom.NewBounds(min, max)
	return c
}

func (c *Cube) Move(x, y, z float64) *Cube {
	return c.transform(geom.Trans(x, y, z))
}

func (c *Cube) Scale(x, y, z float64) *Cube {
	return c.transform(geom.Scale(x, y, z))
}

func (c *Cube) Rotate(x, y, z float64) *Cube {
	return c.transform(geom.Rot(geom.Vec{x, y, z}))
}

func (c *Cube) Center() geom.Vec {
	return c.Pos.MultPoint(geom.Vec{})
}
