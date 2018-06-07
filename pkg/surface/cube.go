package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

// Cube describes the orientation and material of a unit cube
type Cube struct {
	Pos    *geom.Mat
	Mat    Material
	bounds geom.Bounds
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with material and optional transforms.
func UnitCube(m ...Material) *Cube {
	c := &Cube{
		Pos: geom.Identity(),
		Mat: DefaultMaterial{},
	}
	if len(m) > 0 {
		c.Mat = m[0]
	}
	return c.transform(geom.Identity())
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
	c.box = NewBox(min, max)
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

func (c *Cube) Intersect(ray *geom.Ray) (obj Object, dist float64, ok bool) {
	if ok, _, _ := c.bounds.Check(ray); !ok {
		return nil, 0, false
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
			return nil, 0, false
		}
	}
	if tmin > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmin)).Len(); dist >= bias {
			return c, dist, true
		}
	}
	if tmax > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmax)).Len(); dist >= bias {
			return c, dist, true
		}
	}
	return nil, 0, false
}

func (c *Cube) Center() geom.Vec {
	return c.Pos.MultPoint(geom.Vec{})
}

// At returns the normal geom.Vec at this point on the Surface
func (c *Cube) At(pt geom.Vec) (normal geom.Dir, mat Material) {
	normal = geom.Dir{}
	i := c.Pos.Inverse()  // global to local transform
	p1 := i.MultPoint(pt) // translate point into local space
	abs := p1.Abs()
	switch {
	case abs.X > abs.Y && abs.X > abs.Z:
		normal = geom.Dir{math.Copysign(1, p1.X), 0, 0}
	case abs.Y > abs.Z:
		normal = geom.Dir{0, math.Copysign(1, p1.Y), 0}
	default:
		normal = geom.Dir{0, 0, math.Copysign(1, p1.Z)}
	}
	return c.Pos.MultDir(normal), c.Mat.At(0, 0)
}

func (c *Cube) Box() *Box {
	return c.box
}

func (c *Cube) Emits() bool {
	return c.Mat.Emits()
}
