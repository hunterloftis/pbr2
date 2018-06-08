package env

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Gradient struct {
	Up, Down rgb.Energy
}

func NewGradient(down, up rgb.Energy) *Gradient {
	return &Gradient{
		Down: down,
		Up:   up,
	}
}

func (g *Gradient) At(dir geom.Dir) rgb.Energy {
	cos := dir.Dot(geom.Up)
	vertical := (1 + cos) / 2
	return g.Down.Lerp(g.Up, math.Pow(vertical, 3))
}
