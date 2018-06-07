package environment

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Uniform struct {
	Light rgb.Energy
}

func (u *Uniform) At(geom.Dir) rgb.Energy {
	return u.Light
}
