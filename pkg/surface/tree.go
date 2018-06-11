package surface

import (
	"math"
	"sort"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type SurfaceObject interface {
	render.Surface
	render.Object
}

type Tree struct {
	surfaces []SurfaceObject
	lights   []render.Object
	bounds   *geom.Bounds
	left     *Tree
	right    *Tree
	axis     int
	wall     float64
	leaf     bool
}

func NewTree(ss ...SurfaceObject) *Tree {
	t := newBranch(boundsAround(ss...), ss, 0)
	for _, s := range t.surfaces {
		t.lights = append(t.lights, s.Lights()...)
	}
	return t
}

// http://slideplayer.com/slide/7653218/
func (t *Tree) Intersect(ray *geom.Ray) (obj render.Object, dist float64) {
	hit, min, max := t.bounds.Check(ray)
	if !hit {
		return nil, 0
	}
	if t.leaf {
		return t.IntersectSurfaces(ray, max)
	}
	var near, far *Tree
	if ray.DirArray[t.axis] >= 0 {
		near, far = t.left, t.right
	} else {
		near, far = t.right, t.left
	}
	split := (t.wall - ray.OrArray[t.axis]) * ray.InvArray[t.axis]
	if min >= split {
		return far.Intersect(ray)
	}
	if max <= split {
		return near.Intersect(ray)
	}
	if o, d := near.Intersect(ray); o != nil {
		return o, d
	}
	return far.Intersect(ray)
}

func (t *Tree) IntersectSurfaces(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range t.surfaces {
		if o, d := s.Intersect(r); o != nil && d < dist {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (t *Tree) Lights() []render.Object {
	return t.lights
}

func newBranch(bounds *geom.Bounds, surfaces []SurfaceObject, depth int) *Tree {
	t := &Tree{
		surfaces: overlaps(bounds, surfaces),
		bounds:   bounds,
	}
	limit := int(math.Max(1, math.Pow(2, float64(depth-1))))
	if len(t.surfaces) < limit || depth > 12 {
		t.leaf = true
		return t
	}
	t.axis = 0
	max := -1.0
	for i := 0; i < 3; i++ {
		length := bounds.Max.Axis(i) - bounds.Min.Axis(i)
		if length > max {
			max = length
			t.axis = i
		}
	}
	t.wall = median(t.surfaces, t.axis)
	lbounds, rbounds := bounds.Split(t.axis, t.wall)
	t.left = newBranch(lbounds, t.surfaces, depth+1)
	t.right = newBranch(rbounds, t.surfaces, depth+1)
	return t
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

func median(surfaces []SurfaceObject, axis int) float64 {
	vals := make([]float64, 0)
	for _, s := range surfaces {
		b := s.Bounds()
		vals = append(vals, b.MinArray[axis], b.MaxArray[axis])
	}
	sort.Float64s(vals)
	return vals[len(vals)/2]
}