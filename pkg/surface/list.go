package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
)

type List struct {
}

func (l *List) Intersect(r *geom.Ray) (Object, float64, bool) {
	return nil, 0, false
}

func (l *List) Lights() []Object {
	return []Object{}
}
