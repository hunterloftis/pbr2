package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type List struct {
	surfs  []render.Surface
	lights []render.Object
}

func NewList(ss ...render.Surface) *List {
	l := List{
		surfs: ss,
	}
	for _, s := range l.surfs {
		l.lights = append(l.lights, s.Lights()...)
	}
	return &l
}

func (l *List) Intersect(r *geom.Ray) (obj render.Object, dist float64) {
	dist = math.Inf(1)
	for _, s := range l.surfs {
		if o, d := s.Intersect(r); o != nil && d < dist {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (l *List) Lights() []render.Object {
	return l.lights
}
