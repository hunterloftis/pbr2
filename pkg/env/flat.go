package env

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Flat struct {
	Light rgb.Energy
}

func (f *Flat) At(geom.Dir) rgb.Energy {
	return f.Light
}
