package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type List struct {
}

func (l *List) Intersect(r *geom.Ray) (render.Object, float64, bool) {
	return nil, 0, false
}

func (l *List) Lights() []render.Object {
	return []render.Object{}
}
