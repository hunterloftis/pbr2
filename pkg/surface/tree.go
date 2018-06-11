package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type Tree struct {
	surfs  []render.Surface
	lights []render.Object
}

func NewTree(ss ...render.Surface) *Tree {
	t := Tree{
		surfs: ss,
	}
	for _, s := range t.surfs {
		t.lights = append(t.lights, s.Lights()...)
	}
	return &t
}

func (t *Tree) Intersect(r *geom.Ray) (obj render.Object, dist float64) {
	dist = math.Inf(1)
	for _, s := range t.surfs {
		if o, d := s.Intersect(r); o != nil && d < dist {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (t *Tree) Lights() []render.Object {
	return t.lights
}
