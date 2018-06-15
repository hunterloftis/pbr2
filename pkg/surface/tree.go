package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

const (
	maxDepth   = 20
	leafTarget = 10
	xAxis      = 0
	yAxis      = 1
	zAxis      = 2
)

type SurfaceObject interface {
	render.Surface
	render.Object
}

// TODO: This is a very simple k-d tree and could probably be heavily optimized.
type Tree struct {
	Branch
	lights []render.Object
}

type Branch struct {
	surfaces []SurfaceObject
	bounds   *geom.Bounds
	left     *Branch
	right    *Branch
	axis     int
	wall     float64
	leaf     bool
}

func NewTree(ss ...SurfaceObject) *Tree {
	t := Tree{
		Branch: *newBranch(boundsAround(ss...), ss, 0),
	}
	for _, s := range t.Branch.surfaces {
		t.lights = append(t.lights, s.Lights()...)
	}
	return &t
}

// http://slideplayer.com/slide/7653218/
func (b *Branch) Intersect(ray *geom.Ray, maxDist float64) (obj render.Object, dist float64) {
	hit, min, max := b.bounds.Check(ray)
	if !hit || min >= maxDist {
		return nil, 0
	}
	if b.leaf {
		return b.IntersectSurfaces(ray, max)
	}
	var near, far *Branch
	if ray.DirArray[b.axis] >= 0 {
		near, far = b.left, b.right
	} else {
		near, far = b.right, b.left
	}
	split := (b.wall - ray.OrArray[b.axis]) * ray.InvArray[b.axis]
	if min >= split {
		return far.Intersect(ray, maxDist)
	}
	if max <= split {
		return near.Intersect(ray, maxDist)
	}
	if o, d := near.Intersect(ray, maxDist); o != nil {
		return o, d
	}
	return far.Intersect(ray, maxDist)
}

func (b *Branch) IntersectSurfaces(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range b.surfaces {
		if o, d := s.Intersect(r, dist); o != nil {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (t *Tree) Lights() []render.Object {
	return t.lights
}

func newBranch(bounds *geom.Bounds, surfaces []SurfaceObject, depth int) *Branch {
	b := Branch{
		surfaces: overlaps(bounds, surfaces),
		bounds:   bounds,
	}
	if len(b.surfaces) <= leafTarget || depth > maxDepth {
		b.leaf = true
		return &b
	}
	axis, min, max := extents(b.surfaces)
	b.axis = axis
	b.wall = (max + min) / 2
	lbounds, rbounds := bounds.Split(b.axis, b.wall)
	b.left = newBranch(lbounds, b.surfaces, depth+1)
	b.right = newBranch(rbounds, b.surfaces, depth+1)
	return &b
}

func extents(ss []SurfaceObject) (axis int, low, high float64) {
	min := geom.Vec{math.Inf(1), math.Inf(1), math.Inf(1)}
	max := geom.Vec{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for _, s := range ss {
		c := s.Bounds().Center
		min = c.Min(min)
		max = c.Max(max)
	}
	span := max.Minus(min).Abs()
	if span.X > span.Y && span.X > span.Z {
		return xAxis, min.X, max.X
	} else if span.Y > span.X && span.Y > span.Z {
		return yAxis, min.Y, max.Y
	}
	return zAxis, min.Z, max.Z
}

func overlaps(bounds *geom.Bounds, surfaces []SurfaceObject) []SurfaceObject {
	matches := make([]SurfaceObject, 0)
	for _, s := range surfaces {
		if s.Bounds().Overlaps(bounds) {
			matches = append(matches, s)
		}
	}
	return matches
}

func boundsAround(oo ...SurfaceObject) *geom.Bounds {
	if len(oo) == 0 {
		return geom.NewBounds(geom.Vec{}, geom.Vec{})
	}
	Bounds := oo[0].Bounds()
	for _, s := range oo {
		Bounds = geom.MergeBounds(Bounds, s.Bounds())
	}
	return Bounds
}
