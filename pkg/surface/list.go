package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type List struct {
	surfs []render.Surface
}

func NewList(ss ...render.Surface) *List {
	return &List{
		surfs: ss,
	}
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

func (l *List) Lights() (ll []render.Object) {
	for _, s := range l.surfs {
		ll = append(ll, s.Lights()...)
	}
	return ll
}
