package surface

import (
	"sort"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

const maxDepth = 20
const leafTarget = 10

type SurfaceObject interface {
	render.Surface
	render.Object
}

// TODO: This is a very simple k-d tree and could probably be heavily optimized.
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
func (t *Tree) Intersect(ray *geom.Ray, maxDist float64) (obj render.Object, dist float64) {
	hit, min, max := t.bounds.Check(ray)
	if !hit || min >= maxDist {
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

func (t *Tree) IntersectSurfaces(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range t.surfaces {
		if o, d := s.Intersect(r, dist); o != nil {
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
	if len(t.surfaces) <= leafTarget || depth > maxDepth {
		t.leaf = true
		return t
	}
	// TODO: check each of 3 axes and evaluate how good the split would be
	// (balance, number in children vs number in parent, # of overlaps)
	// then choose the best one.
	t.axis = depth % 3
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
