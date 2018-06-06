package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/phys"
)

type List struct {
}

func (l *List) Intersect(r *geom.Ray) (phys.Hit, bool) {
	return phys.Hit{}, false
}

func (l *List) Lights() []phys.Object {
	return []phys.Object{}
}
