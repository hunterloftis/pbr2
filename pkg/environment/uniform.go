package environment

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/phys"
)

type Uniform struct {
	Light phys.Energy
}

func (u *Uniform) At(geom.Dir) phys.Energy {
	return u.Light
}
