package camera

import "github.com/hunterloftis/pbr2/pkg/geom"

type Pinhole struct {
}

func NewPinhole() *Pinhole {
	return &Pinhole{}
}

func (p *Pinhole) Ray(u, v float64) geom.Ray3 {
	return geom.Ray3{}
}
