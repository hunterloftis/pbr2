package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/phys"
)

type List struct {
}

func (l *List) Intersect(r *geom.Ray) (phys.Object, float64, bool) {
	return nil, 0, false
}

func (l *List) Lights() []phys.Object {
	return []phys.Object{}
}
